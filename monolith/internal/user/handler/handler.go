package handler

import (
	"errors"
	"net/http"

	appErr "github.com/waves2k/go-wallet/monolith/internal/errors"
	"github.com/waves2k/go-wallet/monolith/internal/logger"
	"github.com/waves2k/go-wallet/monolith/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/waves2k/go-wallet/monolith/internal/user/model"
	"github.com/waves2k/go-wallet/monolith/internal/user/service"
)

var (
	ErrInvalidIdParam = errors.New("invalid id param")
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) InitRoutes(router fiber.Router) fiber.Router {
	users := router.Group("/users", middleware.ErrorHandler())

	users.Post("/", h.Register)
	users.Get("/:id", h.GetProfile)
	users.Patch("/:id", h.UpdateProfile)

	return users
}

func (h *UserHandler) Register(c fiber.Ctx) error {
	const op = "UserHandler.Register"

	var req model.CreateUserRequest

	if err := c.Bind().JSON(&req); err != nil {
		logger.Log.Warn(err.Error(), op)
		return appErr.ErrBadRequest
	}

	user, err := h.svc.Register(c.Context(), req)
	if err != nil {
		logger.Log.Warn(err.Error(), op)
		return appErr.ErrInternalFailure
	}

	return writeJSON(c, fiber.StatusCreated, user)
}

func (h *UserHandler) GetProfile(c fiber.Ctx) error {
	const op = "UserHandler.GetProfile"

	id := c.Params("id")
	if id == "" {
		return appErr.ErrBadRequest
	}

	user, err := h.svc.GetProfile(c.Context(), id)
	if err != nil {
		return appErr.ErrInternalFailure
	}

	return writeJSON(c, http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c fiber.Ctx) error {
	const op = "UserHandler.UpdateProfile"

	id := c.Params("id")
	if id == "" {
		return appErr.ErrBadRequest
	}

	var req model.UpdateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return appErr.ErrBadRequest
	}

	user, err := h.svc.UpdateProfile(c.Context(), id, req)
	if err != nil {
		return appErr.ErrInternalFailure
	}

	return writeJSON(c, fiber.StatusOK, user)
}
