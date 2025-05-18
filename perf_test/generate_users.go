package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

func randStringFromCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func randUsername() string {
	return randStringFromCharset(15, letters)
}

func randPassword() string {
	return "123456gFzxcFXZCorrect"
}

func toBase64JSONString(data map[string]interface{}) string {
	raw, _ := json.Marshal(data)
	return base64.StdEncoding.EncodeToString(raw)
}

func main() {

	serverAddr := "http://localhost:8025/api/v1/auth"

	signupFile, err := os.Create("signup_targets.json")
	if err != nil {
		panic(err)
	}
	defer signupFile.Close()

	loginFile, err := os.Create("login_targets.json")
	if err != nil {
		panic(err)
	}
	defer loginFile.Close()

	for i := 0; i < 100_000; i++ {
		username := randUsername()
		email := fmt.Sprintf("%s@mail.ru", username)
		password := randPassword()

		// Signup JSON target
		signupBody := map[string]interface{}{
			"username": username,
			"email":    email,
			"password": password,
			"role":     1,
		}
		signupTarget := map[string]interface{}{
			"method": "POST",
			"url":    fmt.Sprintf("%s/signup", serverAddr),
			"header": map[string][]string{
				"Content-Type": {"application/json"},
			},
			"body": toBase64JSONString(signupBody),
		}
		writeJSONLine(signupFile, signupTarget)

		// Login JSON target
		loginBody := map[string]interface{}{
			"email":    email,
			"password": password,
			"role":     1,
		}
		loginTarget := map[string]interface{}{
			"method": "POST",
			"url":    fmt.Sprintf("%s/login", serverAddr),
			"header": map[string][]string{
				"Content-Type": {"application/json"},
			},
			"body": toBase64JSONString(loginBody),
		}
		writeJSONLine(loginFile, loginTarget)
	}

	fmt.Println("Файлы signup_targets.json и login_targets.json успешно сгенерированы.")
}

// writeJSONLine сериализует JSON и записывает его в одну строку
func writeJSONLine(f *os.File, target map[string]interface{}) {
	line, err := json.Marshal(target)
	if err != nil {
		panic(err)
	}
	f.WriteString(string(line) + "\n")
}
