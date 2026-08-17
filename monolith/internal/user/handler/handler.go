package handler

import (
	"errors"
	"net/http"

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
	users := router.Group("/users")

	users.Post("/", h.Register)
	users.Get("/:id", h.GetProfile)
	users.Patch("/:id", h.UpdateProfile)

	return users
}

func (h *UserHandler) Register(c fiber.Ctx) error {
	var req model.CreateUserRequest

	if err := c.Bind().JSON(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}

	user, err := h.svc.Register(c.Context(), req)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, err)
	}

	return writeJSON(c, fiber.StatusCreated, user)
}

func (h *UserHandler) GetProfile(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return writeError(c, http.StatusBadRequest, ErrInvalidIdParam)
	}

	user, err := h.svc.GetProfile(c.Context(), id)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, err)
	}

	return writeJSON(c, http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return writeError(c, http.StatusBadRequest, ErrInvalidIdParam)
	}

	var req model.UpdateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return writeError(c, http.StatusBadRequest, err)
	}

	user, err := h.svc.UpdateProfile(c.Context(), id, req)
	if err != nil {
		return writeError(c, http.StatusInternalServerError, err)
	}

	return writeJSON(c, fiber.StatusOK, user)
}
