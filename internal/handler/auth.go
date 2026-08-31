package handler

import (
	"errors"
	"log"
	"net/http"

	"example.com/acg-go-demo/internal/dto"
	"example.com/acg-go-demo/internal/response"
	"example.com/acg-go-demo/internal/service"
	"github.com/gin-gonic/gin"
)

const authCookieMaxAge = 7 * 24 * 60 * 60

type AuthHandler struct {
	authService service.AuthService
	cookieName  string
	secure      bool
}

func NewAuthHandler(
	authService service.AuthService,
	cookieName string,
	secure bool,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cookieName:  cookieName,
		secure:      secure,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid request", nil)
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// The raw token is kept in an HttpOnly cookie and is not returned in JSON.
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		h.cookieName,
		result.Token,
		authCookieMaxAge,
		"/",
		"",
		h.secure,
		true,
	)

	response.Success(c, result.User)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	rawToken, err := c.Cookie(h.cookieName)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		log.Printf("read auth cookie failed, err=%v", err)
		response.Fail(c, http.StatusInternalServerError, 50000, "internal server error", nil)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), rawToken); err != nil {
		h.handleServiceError(c, err)
		return
	}

	// Clearing a missing or already-invalid cookie is intentionally idempotent.
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(h.cookieName, "", -1, "/", "", h.secure, true)
	response.Success(c, nil)
}

func (h *AuthHandler) handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		response.Fail(c, http.StatusUnauthorized, 40101, "invalid credentials", nil)
	default:
		log.Printf("auth service error: %v", err)
		response.Fail(c, http.StatusInternalServerError, 50000, "internal server error", nil)
	}
}
