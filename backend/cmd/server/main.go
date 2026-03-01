package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"github.com/zhuliliang/crazyLI/backend/internal/auth"
	"github.com/zhuliliang/crazyLI/backend/internal/config"
	"github.com/zhuliliang/crazyLI/backend/internal/db"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()
	store := db.NewStore(pool)

	googleOAuth := auth.NewGoogleOAuth(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Get("/auth/google/login", func(w http.ResponseWriter, r *http.Request) {
		state := auth.NewStateToken()
		auth.SetStateCookie(w, state)
		http.Redirect(w, r, googleOAuth.AuthURL(state), http.StatusTemporaryRedirect)
	})

	r.Get("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		if err := auth.ValidateStateCookie(r, r.URL.Query().Get("state")); err != nil {
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "missing code", http.StatusBadRequest)
			return
		}
		token, err := googleOAuth.ExchangeCode(r.Context(), code)
		if err != nil {
			http.Error(w, "exchange failed", http.StatusBadGateway)
			return
		}
		userInfo, err := googleOAuth.FetchUser(r.Context(), token)
		if err != nil {
			http.Error(w, "user info failed", http.StatusBadGateway)
			return
		}
		savedUser, err := store.UpsertUser(r.Context(), db.User{
			GoogleID: userInfo.ID,
			Email:    userInfo.Email,
			Name:     userInfo.Name,
			Picture:  userInfo.Picture,
		})
		if err != nil {
			http.Error(w, "persist failed", http.StatusInternalServerError)
			return
		}
		if err := issueSessionCookie(w, cfg, savedUser); err != nil {
			log.Printf("session error: %v", err)
			http.Error(w, "session error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, cfg.FrontendURL+"/login/success", http.StatusTemporaryRedirect)
	})

	r.Post("/auth/email/register", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Name     string `json:"name"`
		}
		if err := decodeJSON(r, &body); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		body.Email = strings.TrimSpace(strings.ToLower(body.Email))
		body.Name = strings.TrimSpace(body.Name)
		if body.Email == "" || body.Password == "" || body.Name == "" {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}
		_, err := store.GetUserByEmail(r.Context(), body.Email)
		if err == nil {
			http.Error(w, "email already registered", http.StatusConflict)
			return
		}
		if err != nil && !errors.Is(err, db.ErrNotFound) {
			log.Printf("register lookup error: %v", err)
			http.Error(w, "lookup error", http.StatusInternalServerError)
			return
		}
		hash, err := auth.HashPassword(body.Password)
		if err != nil {
			log.Printf("hashing error: %v", err)
			http.Error(w, "hashing error", http.StatusInternalServerError)
			return
		}
		newUser, err := store.CreateEmailUser(r.Context(), db.User{
			Email:        body.Email,
			Name:         body.Name,
			PasswordHash: hash,
		})
		if err != nil {
			log.Printf("register persist error: %v", err)
			log.Printf("register persist error: %v", err)
			http.Error(w, "persist error", http.StatusInternalServerError)
			return
		}
		if err := issueSessionCookie(w, cfg, newUser); err != nil {
			http.Error(w, "session error", http.StatusInternalServerError)
			return
		}
		respondJSON(w, map[string]string{"status": "registered"})
	})

	r.Post("/auth/email/login", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := decodeJSON(r, &body); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		body.Email = strings.TrimSpace(strings.ToLower(body.Email))
		if body.Email == "" || body.Password == "" {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}
		user, err := store.GetUserByEmail(r.Context(), body.Email)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		if user.Provider != "email" || user.PasswordHash == "" {
			http.Error(w, "use Google login", http.StatusBadRequest)
			return
		}
		if err := auth.ComparePassword(user.PasswordHash, body.Password); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		if err := issueSessionCookie(w, cfg, user); err != nil {
			http.Error(w, "session error", http.StatusInternalServerError)
			return
		}
		respondJSON(w, map[string]string{"status": "ok"})
	})

	r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "missing session", http.StatusUnauthorized)
			return
		}
		claims, err := auth.ParseSession(cfg.JWTSecret, cookie.Value)
		if err != nil {
			http.Error(w, "invalid session", http.StatusUnauthorized)
			return
		}
		user, err := store.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}
		respondJSON(w, map[string]interface{}{
			"id":       user.ID,
			"email":    user.Email,
			"name":     user.Name,
			"picture":  user.Picture,
			"provider": user.Provider,
		})
	})

	addr := ":" + cfg.Port
	log.Printf("server running on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func issueSessionCookie(w http.ResponseWriter, cfg config.Config, user db.User) error {
	token, err := auth.CreateSessionToken(cfg.JWTSecret, user.ID, user.Email, user.Name)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	return nil
}

func respondJSON(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
