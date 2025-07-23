package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Password struct {
	Password string `json:"password"`
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var password Password

	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, "failed to read request body")
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &password); err != nil {
		SendError(w, fmt.Sprintf("invalid JSON: %v", err))
		return
	}

	value := os.Getenv("TODO_PASSWORD")
	if value == "" {
		SendError(w, "Authentication is not configured")
		return
	}

	if password.Password != value {
		SendError(w, "Неверный пароль")
		return
	}

	secret := []byte("my_secret_key")

	claims := jwt.MapClaims{
		"password_hash": fmt.Sprintf("%x", sha256.Sum256([]byte(password.Password))),
		"exp":           time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		SendError(w, "failed to sign jwt")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    signedToken,
		MaxAge:   8 * 3600,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})

	response := map[string]string{"token": signedToken}
	json.NewEncoder(w).Encode(response)
}
