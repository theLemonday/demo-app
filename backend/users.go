package main

import (
	"context"
	"net/http"
)

type (
	ContextKey int
	Role       string
)

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

var users = map[string]string{
	"admin": "adminpass",
	"user":  "userpass",
}

var userRoles = map[string]Role{
	"admin": RoleAdmin,
	"user":  RoleUser,
}

const (
	userContextKey ContextKey = iota
)

func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			if users[username] == password {
				role := userRoles[username]
				ctx := context.WithValue(r.Context(), userContextKey, role)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}

func authorizeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(userContextKey).(Role)
		if role == RoleAdmin {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
