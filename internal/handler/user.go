package handler

import (
	"errors"
	"log"
	"strconv"

	"example.com/acg-go-demo/internal/dto"
	"example.com/acg-go-demo/internal/response"
	"example.com/acg-go-demo/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		response.Fail(c, 404, 40401, "user not found", nil)
	case errors.Is(err, service.ErrAccountTaken):
		response.Fail(c, 409, 40901, "account already exists", nil)
	default:
		log.Printf("user service error: %v", err)
		response.Fail(c, 500, 50000, "internal server error", nil)
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 40001, "invalid reqeust", nil)
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, user)
}

func parseUserID(c *gin.Context) (uint64, bool) {
	rawID := c.Param("id")

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, 400, 40002, "invalid user id", nil)
		return 0, false
	}
	return id, true
}

func (h *UserHandler) Get(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}

	user, err := h.userService.Get(c.Request.Context(), id)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) UpdateNickname(c *gin.Context) {
	id, ok := parseUserID(c)

	if !ok {
		return
	}

	var req dto.UpdateNicknameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 40001, "invalid request", nil)
		return
	}

	err := h.userService.UpdateNickname(c.Request.Context(), id, req.Nickname)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, ok := parseUserID(c)

	if !ok {
		return
	}

	err := h.userService.Delete(c.Request.Context(), id)

	if err != nil {
		h.handleServiceError(c, err)

		return
	}

	response.Success(c, nil)
}
