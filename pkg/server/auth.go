package server

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("my_secret_key")

type AuthRequest struct {
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AuthResponse{Error: "method not allowed"})
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthResponse{Error: "invalid request"})
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AuthResponse{Error: "authentication not configured"})
		return
	}

	if req.Password != expectedPassword {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AuthResponse{Error: "incorrect password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password_hash": fmt.Sprintf("%x", sha256.Sum256([]byte(expectedPassword))),
		"exp":           time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AuthResponse{Error: "failed to generate token"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(8 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	if err := json.NewEncoder(w).Encode(AuthResponse{Token: tokenString}); err != nil {
		fmt.Printf("error encoding response: %v\n", err)
	}
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		tokenString := ""
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		}

		if tokenString == "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				sendAuthError(w, r)
				return
			}
			tokenString = cookie.Value
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			sendAuthError(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			sendAuthError(w, r)
			return
		}

		currentHash := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))
		if claims["password_hash"] != currentHash {
			sendAuthError(w, r)
			return
		}

		next(w, r)
	}
}

func sendAuthError(w http.ResponseWriter, r *http.Request) {
	isAPI := false
	if len(r.URL.Path) >= 4 {
		isAPI = r.URL.Path[:4] == "/api"
	}

	if r.Header.Get("Content-Type") == "application/json" || isAPI {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"}); err != nil {
			fmt.Printf("Error encoding JSON: %v\n", err)
		}
	} else {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
	}
}
