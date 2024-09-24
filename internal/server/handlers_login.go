package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/canonical/notary/internal/db"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func expireAfter() int64 {
	return time.Now().Add(time.Hour * 1).Unix()
}

type jwtNotaryClaims struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Permissions int    `json:"permissions"`
	jwt.StandardClaims
}

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Helper function to generate a JWT
func generateJWT(id int, username string, jwtSecret []byte, permissions int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtNotaryClaims{
		ID:          id,
		Username:    username,
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAfter(),
		},
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Login(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginParams LoginParams
		if err := json.NewDecoder(r.Body).Decode(&loginParams); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if loginParams.Username == "" {
			writeError(w, http.StatusBadRequest, "Username is required")
			return
		}
		if loginParams.Password == "" {
			writeError(w, http.StatusBadRequest, "Password is required")
			return
		}
		userAccount, err := env.DB.RetrieveUserByUsername(loginParams.Username)
		if err != nil {
			log.Println(err)
			if errors.Is(err, db.ErrIdNotFound) {
				writeError(w, http.StatusUnauthorized, "The username or password is incorrect. Try again.")
				return
			}
			writeError(w, http.StatusInternalServerError, "Internal Error")
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(userAccount.Password), []byte(loginParams.Password)); err != nil {
			writeError(w, http.StatusUnauthorized, "The username or password is incorrect. Try again.")
			return
		}
		jwt, err := generateJWT(userAccount.ID, userAccount.Username, env.JWTSecret, userAccount.Permissions)
		if err != nil {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, "Internal Error")
			return
		}
		w.WriteHeader(http.StatusOK)
		loginResponse := LoginResponse{
			Token: jwt,
		}
		err = writeJSON(w, loginResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}
