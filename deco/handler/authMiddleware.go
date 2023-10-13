package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var authUrl = "http://localhost:8081/users"

type user struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

func AuthMiddleware(w http.ResponseWriter, r *http.Request, h http.Handler) {
	token := r.URL.Query().Get("token")
	re, _ := http.NewRequest("POST", authUrl+"/verify?token="+token, nil)
	resp, err := http.DefaultClient.Do(re)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	if resp.Body == nil {
		fmt.Println("response body is nil")
		return
	}

	parts := strings.Split(token, ".")
	payload, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println(err)
	}
	u := user{}
	err = json.Unmarshal(payload, &u)
	fmt.Println(u)
	ctx := context.WithValue(r.Context(), "userId", u.UserID)
	h.ServeHTTP(w, r.WithContext(ctx))
}
