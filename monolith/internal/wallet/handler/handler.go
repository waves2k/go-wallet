package handler

import (
	"github.com/waves2k/go-wallet/monolith/internal/wallet/service"
)

type WalletHandler struct {
	svc service.WalletService
}

func NewWalletService(svc service.WalletService) *WalletHandler {
	return &WalletHandler{
		svc: svc,
	}
}

// func (h *WalletHandler) InitRoutes(router fiber.Router) fiber.Router {
// 	router.Group()
// }
