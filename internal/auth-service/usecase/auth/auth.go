package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	entityAuth "retarget/internal/auth-service/entity/auth"
	repoAuth "retarget/internal/auth-service/repo/auth"

	"golang.org/x/crypto/argon2"
	"gopkg.in/inf.v0"
)

type HashConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var (
	DefaultHashConfig = HashConfig{
		Memory:      64 * 1024, // 64 MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
)

type AuthUsecaseInterface interface {
	Login(email string, password string, role int, requestID string) (*entityAuth.User, error)
	Logout(sessionId string) error
	Register(username string, email string, password string, role int, requestID string) (*entityAuth.User, error)

	GetUser(userId int, requestID string) (*entityAuth.User, error)
	CheckCode(code int, userId int) error
	CreateCode(userId int) (int, error)

	AddSession(userId int, role int) (*entityAuth.Session, error)
}

type AuthUsecase struct {
	authRepository    *repoAuth.AuthRepository
	sessionRepository *repoAuth.SessionRepository
	hashCfg           HashConfig
}

func NewAuthUsecase(userRepo *repoAuth.AuthRepository, sessionRepo *repoAuth.SessionRepository) *AuthUsecase {
	return &AuthUsecase{
		authRepository:    userRepo,
		sessionRepository: sessionRepo,
		hashCfg:           DefaultHashConfig,
	}
}

func (a *AuthUsecase) Login(email string, password string, role int, requestID string) (*entityAuth.User, error) {
	user, err := a.authRepository.GetUserByEmail(email, requestID)
	if err != nil {
		return nil, errors.New("incorrect user data")
	}

	if err := compareHashAndPassword(string(user.Password), password, a.hashCfg); err != nil {
		return nil, errors.New("incorrect user data")
	}

	return user, nil
}

func (a *AuthUsecase) Logout(sessionId string) error {
	err := a.sessionRepository.DelSession(sessionId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUsecase) GetUser(userID int, requestID string) (*entityAuth.User, error) {
	user, err := a.authRepository.GetUserByID(userID, requestID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *AuthUsecase) Register(username string, email string, password string, role int, requestID string) (*entityAuth.User, error) {
	existingUser, err := a.authRepository.CheckEmailOrUsernameExists(email, username, requestID)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		if existingUser.Email == email {
			return nil, errors.New("пользователь с таким email уже существует")
		}
		if existingUser.Username == username {
			return nil, errors.New("пользователь с таким username уже существует")
		}
	}

	hashedPassword, err := hashPassword(password, a.hashCfg)
	if err != nil {
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
		return nil, err
	}

	if err := a.authRepository.CreateNewUser(user, requestID); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthUsecase) AddSession(userID int, role int) (*entityAuth.Session, error) {
	session, err := a.sessionRepository.AddSession(userID, role)
	if err != nil {
		return nil, err
	}
	return session, nil
}

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
	fmt.Sscanf(parts[2], "v=%d", &version)
	if version != argon2.Version {
		return errors.New("incompatible version")
	}

	var memory, iterations uint32
	var parallelism uint8
	fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)

	salt, _ := base64.RawStdEncoding.DecodeString(parts[4])
	expectedHash, _ := base64.RawStdEncoding.DecodeString(parts[5])

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(expectedHash)))

	if !bytes.Equal(hash, expectedHash) {
		return errors.New("incorrect password")
	}

	return nil
}
