package handlers

import (
	"encoding/json"
	"errors"
	"file-analyzer/internals/handlers/dto"
	"file-analyzer/internals/services"
	"file-analyzer/internals/utils"
	"log"
	"net/http"
	"time"
)

type AuthHandler struct {
	service *services.AuthService
	l       *log.Logger
}

func NewAuthHandler(l *log.Logger, service *services.AuthService) *AuthHandler {
	return &AuthHandler{service, l}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	err := DecodeJSON(r, &req)
	if err != nil {
		utils.FAIL(w, http.StatusBadRequest, "Invalid Request")
		return
	}
	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		h.l.Printf("Login Error %v ", err)
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.FAIL(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	SetCookie(r, w, "refresh_token", token.RefreshToken, 7*24*time.Hour)
	utils.SUCCESS(w, "Login Successfully", dto.LoginResponse{
		AccessToken: token.AccessToken,
	})
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	err := DecodeJSON(r, &req)
	if err != nil {
		utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	err = h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		h.l.Printf("Register Error %v ", err)
		utils.FAIL(w, 500, "Registration failed")
		return
	}
	utils.SUCCESS(w, "Users Registered Successfully", nil)
}

func (h *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	incomingRefreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		h.l.Println("Error Reading Refresh Token")
		utils.FAIL(w, http.StatusBadRequest, "Internal Server Error")
	}
	token, err := h.service.Refresh(incomingRefreshToken.Name)
	SetCookie(r, w, "refresh_token", token.RefreshToken, 7*24*time.Hour)
	utils.SUCCESS(w, "Token Refreshed Successfully", dto.LoginResponse{
		AccessToken: token.AccessToken,
	})
}

func SetCookie(r *http.Request, w http.ResponseWriter, name string, value string, maxAge time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Secure:   r.TLS != nil,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(maxAge.Seconds()),
	})
}

func DecodeJSON[T *dto.LoginRequest | *dto.RegisterRequest](r *http.Request, dst T) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
