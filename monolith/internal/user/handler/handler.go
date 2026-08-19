package handler

import (
	"errors"
	"fmt"
	"net/http"

	appErr "github.com/waves2k/go-wallet/monolith/internal/errors"
	"github.com/waves2k/go-wallet/monolith/internal/logger"
	"github.com/waves2k/go-wallet/monolith/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/waves2k/go-wallet/monolith/internal/user/model"
	"github.com/waves2k/go-wallet/monolith/internal/user/service"
	userSvc "github.com/waves2k/go-wallet/monolith/internal/user/service"
	walSvc "github.com/waves2k/go-wallet/monolith/internal/wallet/service"
)

var (
	ErrInvalidIdParam = errors.New("invalid id param")
)

type UserHandler struct {
	userSvc userSvc.UserService
	walSvc  walSvc.WalletService
}

func NewUserHandler(
	uSvc service.UserService,
	wSvc walSvc.WalletService,
) *UserHandler {
	return &UserHandler{
		userSvc: uSvc,
		walSvc:  wSvc,
	}
}

func (h *UserHandler) InitRoutes(router fiber.Router) fiber.Router {
	users := router.Group("/users" /*recovery(),*/, middleware.RequestID(), middleware.ErrorHandler())

	users.Post("/", h.handleRegister)
	users.Get("/:id", middleware.AuthMiddleware(), h.handleGetProfile)
	users.Post("/login/", h.handleLogin)
	users.Patch("/:id", middleware.AuthMiddleware(), h.handleUpdateProfile)
	users.Get("/:id/wallet", middleware.AuthMiddleware(), h.handleGetUserWallet)

	return users
}

func (h *UserHandler) handleRegister(c fiber.Ctx) error {
	const op = "UserHandler.handleRegister"

	var req model.CreateUserRequest

	if err := c.Bind().JSON(&req); err != nil {
		logger.Log.Warn(err.Error(), op)
		return appErr.ErrBadRequest
	}

	fmt.Println(op, req)

	user, err := h.userSvc.Register(c.Context(), req)
	if err != nil {
		logger.Log.Warn(err.Error(), op)
		return appErr.ErrInternalFailure
	}

	fmt.Println(op, user)

	return writeJSON(c, fiber.StatusCreated, user.ToResponse())
}

func (h *UserHandler) handleGetProfile(c fiber.Ctx) error {
	const op = "UserHandler.handleGetProfile"

	id := c.Params("id")
	if id == "" {
		logger.Log.Warn(ErrInvalidIdParam.Error(), op)
		return appErr.ErrBadRequest
	}

	user, err := h.userSvc.GetProfile(c.Context(), id)
	if err != nil {
		logger.Log.Warn(err.Error(), op)
		return appErr.ErrInternalFailure
	}

	return writeJSON(c, http.StatusOK, user.ToResponse())
}

func (h *UserHandler) handleUpdateProfile(c fiber.Ctx) error {
	const op = "UserHandler.handleUpdateProfile"

	id := c.Params("id")
	if id == "" {
		logger.Log.Warn(ErrInvalidIdParam.Error(), op)
		return appErr.ErrBadRequest
	}

	var req model.UpdateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		logger.Log.Warn(err.Error(), op)
		return appErr.ErrBadRequest
	}

	user, err := h.userSvc.UpdateProfile(c.Context(), id, req)
	if err != nil {
		logger.Log.Warn(err.Error(), op)
		return appErr.ErrInternalFailure
	}

	return writeJSON(c, fiber.StatusOK, user.ToResponse())
}

func (h *UserHandler) handleLogin(c fiber.Ctx) error {
	const op = "UserHandler.handleLogin"

	var req model.LoginRequest

	if err := c.Bind().JSON(&req); err != nil {
		logger.Log.Warn(err.Error(), op)
		return appErr.ErrBadRequest
	}

	fmt.Println(op, req)

	resp, err := h.userSvc.Login(c.Context(), req)
	if err != nil {
		logger.Log.Warn(err.Error(), op)
		return appErr.ErrInternalFailure
	}

	return writeJSON(c, fiber.StatusOK, resp)
}

func (h *UserHandler) handleGetUserWallet(c fiber.Ctx) error {
	const op = "UserHandler.handleGetUserWallet"

	id := c.Params("id")
	if id == "" {
		logger.Log.Warn(ErrInvalidIdParam.Error(), op)
		return appErr.ErrBadRequest
	}

	wallet, err := h.walSvc.GetByUserID(c.Context(), id)
	if err != nil {
		logger.Log.Warn(err.Error(), op)
		return appErr.ErrInternalFailure
	}

	return writeJSON(c, fiber.StatusOK, wallet.ToResponse())
}
