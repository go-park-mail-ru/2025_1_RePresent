package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	entityAuth "retarget/internal/auth-service/entity/auth"
	repoAuth "retarget/internal/auth-service/repo/auth"
	"retarget/pkg/utils/optiLog"

	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/argon2"
	"gopkg.in/inf.v0"
)

// HashConfig определяет параметры Argon2id
type HashConfig struct {
	Memory      uint32 // in KB
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultHashConfig = HashConfig{
	Memory:      16 * 1024, // 16 MB
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// SimpleRateLimiter - простой рейт-лимитер на основе времени последнего доступа
type SimpleRateLimiter struct {
	lastAccess map[string]time.Time
	interval   time.Duration
	mu         sync.RWMutex
}

func NewSimpleRateLimiter(interval time.Duration) *SimpleRateLimiter {
	return &SimpleRateLimiter{
		lastAccess: make(map[string]time.Time),
		interval:   interval,
	}
}

func (r *SimpleRateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if last, ok := r.lastAccess[key]; ok && now.Sub(last) < r.interval {
		return false
	}
	r.lastAccess[key] = now
	return true
}

type AuthUsecaseInterface interface {
	Login(ctx context.Context, email string, password string, role int, requestID string) (*entityAuth.User, error)
	Logout(sessionId string) error
	Register(ctx context.Context, username string, email string, password string, role int, requestID string) (*entityAuth.User, error)
	GetUser(ctx context.Context, userID int, requestID string) (*entityAuth.User, error)
	CheckCode(code int, userId int) error
	CreateCode(userId int) (int, error)
	AddSession(userID int, role int) (*entityAuth.Session, error)
}

type AuthUsecase struct {
	authRepository    *repoAuth.AuthRepository
	sessionRepository *repoAuth.SessionRepository
	hashCfg           HashConfig
	rateLimiter       *SimpleRateLimiter
	asyncLogger       *optiLog.AsyncLogger
}

func NewAuthUsecase(
	userRepo *repoAuth.AuthRepository,
	sessionRepo *repoAuth.SessionRepository,
	logger *optiLog.AsyncLogger,
) *AuthUsecase {
	return &AuthUsecase{
		authRepository:    userRepo,
		sessionRepository: sessionRepo,
		hashCfg:           DefaultHashConfig,
		rateLimiter:       NewSimpleRateLimiter(1 * time.Second),
		asyncLogger:       logger,
	}
}

// -----------------------------
// Методы авторизации
// -----------------------------

func (a *AuthUsecase) Login(ctx context.Context, email string, password string, role int, requestID string) (*entityAuth.User, error) {
	startTime := time.Now()

	if !a.rateLimiter.Allow(email) {
		a.asyncLogger.Log(zapcore.WarnLevel, requestID, "Too many login attempts",
			optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
				"email": email,
			}))
		return nil, errors.New("слишком много попыток")
	}

	var user *entityAuth.User
	err := withRetry(func() error {
		var fetchErr error
		user, fetchErr = a.authRepository.GetUserByEmail(email, requestID)
		return fetchErr
	})
	if err != nil {
		a.asyncLogger.Log(zapcore.WarnLevel, requestID, "User not found or DB error",
			optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
				"email": email,
				"error": err.Error(),
			}))
		return nil, errors.New("incorrect user data")
	}

	if err := compareHashAndPassword(string(user.Password), password, a.hashCfg); err != nil {
		a.asyncLogger.Log(zapcore.WarnLevel, requestID, "Password mismatch",
			optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
				"userID": user.ID,
				"email":  user.Email,
			}))
		return nil, errors.New("incorrect user data")
	}

	a.asyncLogger.Log(zapcore.DebugLevel, requestID, "Login successful",
		optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
			"userID": user.ID,
			"email":  user.Email,
		}))

	return user, nil
}

func (a *AuthUsecase) Logout(sessionId string) error {
	err := a.sessionRepository.DelSession(sessionId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUsecase) GetUser(ctx context.Context, userID int, requestID string) (*entityAuth.User, error) {
	startTime := time.Now()

	var user *entityAuth.User
	err := withRetry(func() error {
		var fetchErr error
		user, fetchErr = a.authRepository.GetUserByID(userID, requestID)
		return fetchErr
	})

	if err != nil {
		a.asyncLogger.Log(zapcore.WarnLevel, requestID, "Failed to get user",
			optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
				"userID": userID,
				"error":  err.Error(),
			}))
		return nil, err
	}

	a.asyncLogger.Log(zapcore.DebugLevel, requestID, "User fetched successfully",
		optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
			"userID": user.ID,
			"role":   user.Role,
		}))

	return user, nil
}

func (a *AuthUsecase) Register(ctx context.Context, username string, email string, password string, role int, requestID string) (*entityAuth.User, error) {
	startTime := time.Now()

	if len(password) < 8 {
		a.asyncLogger.Log(zapcore.WarnLevel, requestID, "Password too short",
			optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
				"username": username,
			}))
		return nil, errors.New("пароль должен быть не короче 8 символов")
	}

	existingUser, err := a.authRepository.CheckEmailOrUsernameExists(email, username, requestID)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		if existingUser.Email == email {
			a.asyncLogger.Log(zapcore.WarnLevel, requestID, "Email already exists",
				optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
					"email": email,
				}))
			return nil, errors.New("пользователь с таким email уже существует")
		}
		if existingUser.Username == username {
			a.asyncLogger.Log(zapcore.WarnLevel, requestID, "Username already exists",
				optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
					"username": username,
				}))
			return nil, errors.New("пользователь с таким username уже существует")
		}
	}

	hashedPassword, err := hashPassword(password, a.hashCfg)
	if err != nil {
		a.asyncLogger.Log(zapcore.WarnLevel, requestID, "Password hashing failed",
			optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
				"error": err.Error(),
			}))
		return nil, err
	}

	user := &entityAuth.User{
		Username:    username,
		Email:       email,
		Password:    hashedPassword,
		Description: "",
		Balance:     entityAuth.Decimal{Dec: inf.NewDec(0, 0)},
		Role:        role,
	}

	if err := entityAuth.ValidateUser(user); err != nil {
		a.asyncLogger.Log(zapcore.WarnLevel, requestID, "User validation failed",
			optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
				"username": user.Username,
				"error":    err.Error(),
			}))
		return nil, err
	}

	if err := a.authRepository.CreateNewUser(user, requestID); err != nil {
		a.asyncLogger.Log(zapcore.WarnLevel, requestID, "User creation failed",
			optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
				"username": user.Username,
				"error":    err.Error(),
			}))
		return nil, err
	}

	a.asyncLogger.Log(zapcore.DebugLevel, requestID, "User registered successfully",
		optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
			"userID": user.ID,
			"email":  user.Email,
		}))

	return user, nil
}

func (a *AuthUsecase) AddSession(userID int, role int) (*entityAuth.Session, error) {
	session, err := a.sessionRepository.AddSession(userID, role)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// -----------------------------
// Хэширование паролей
// -----------------------------

func hashPassword(password string, cfg HashConfig) ([]byte, error) {
	salt := make([]byte, cfg.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	hash := argon2.IDKey([]byte(password), salt, cfg.Iterations, cfg.Memory, cfg.Parallelism, cfg.KeyLength)

	b64Encode := func(b []byte) string {
		return base64.RawStdEncoding.EncodeToString(b)
	}

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		cfg.Memory, cfg.Iterations, cfg.Parallelism,
		b64Encode(salt),
		b64Encode(hash))

	return []byte(encodedHash), nil
}

func compareHashAndPassword(storedHash, password string, cfg HashConfig) error {
	parts := strings.Split(storedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return errors.New("invalid hash format")
	}

	var version int
	n, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if n != 1 || err != nil {
		return fmt.Errorf("failed to parse version: %w", err)
	}
	if version != argon2.Version {
		return errors.New("incompatible version")
	}

	var memory, iterations uint32
	var parallelism uint8

	n, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if n != 3 || err != nil {
		return fmt.Errorf("failed to parse m/t/p values: %w", err)
	}

	salt, _ := base64.RawStdEncoding.DecodeString(parts[4])
	expectedHash, _ := base64.RawStdEncoding.DecodeString(parts[5])

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(expectedHash)))

	if !bytes.Equal(hash, expectedHash) {
		return errors.New("incorrect password")
	}

	return nil
}

// -----------------------------
// Вспомогательные функции
// -----------------------------

func withRetry(fn func() error) error {
	for attempt := 0; attempt < 3; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		if isTransientError(err) {
			if attempt < 2 {
				delay := time.Millisecond * time.Duration(100*(1<<uint(attempt)))
				time.Sleep(delay)
			}
			continue
		}
		break
	}
	return errors.New("operation failed after retries")
}

func isTransientError(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "network error") ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "EOF") ||
		strings.Contains(err.Error(), "unexpected EOF") ||
		strings.Contains(err.Error(), "server closed the connection unexpectedly") {
		return true
	}
	return false
}
