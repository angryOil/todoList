package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type user struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

func AuthMiddleware(w http.ResponseWriter, r *http.Request, h http.Handler) {
	token := r.Header.Get("token")
	fmt.Println("token?", token)
	if !tokenCheck(token) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid token"))
		return
	}

	parts := strings.Split(token, ".")
	payload, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	u := user{}
	err = json.Unmarshal(payload, &u)
	if u.UserID == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("user id is not valid"))
		return
	}
	ctx := context.WithValue(r.Context(), "userId", u.UserID)

	h.ServeHTTP(w, r.WithContext(ctx))
}

func tokenCheck(token string) bool {
	if token == "" {
		return false
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}
	return true
}
