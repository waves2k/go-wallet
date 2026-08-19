package service

import (
	"context"
	"fmt"
	"time"

	"github.com/waves2k/go-wallet/monolith/internal/auth"

	uModel "github.com/waves2k/go-wallet/monolith/internal/user/model"
	userRepo "github.com/waves2k/go-wallet/monolith/internal/user/repository"
	walModel "github.com/waves2k/go-wallet/monolith/internal/wallet/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	accessTokenDefaultDuration  = time.Minute * 15
	refreshTokenDefaultDuration = time.Hour * 24 * 7
)

type UserService interface {
	Register(ctx context.Context, req uModel.CreateUserRequest) (*uModel.UserWithWallet, error)
	GetProfile(ctx context.Context, id string) (*uModel.User, error)
	UpdateProfile(ctx context.Context, id string, req uModel.UpdateUserRequest) (*uModel.User, error)
	Login(ctx context.Context, req uModel.LoginRequest) (*uModel.LoginResponse, error)
}

type userService struct {
	userRepo userRepo.UserRepository
}

func NewUserService(
	uRepo userRepo.UserRepository,
) UserService {
	return &userService{
		userRepo: uRepo,
	}
}

func (s *userService) Register(ctx context.Context, req uModel.CreateUserRequest) (*uModel.UserWithWallet, error) {
	const op = "userService.Register"

	user, err := uModel.New(req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	wallet, err := walModel.New(user.ID)
	if err != nil {
		return nil, err
	}

	userWithWal, err := s.userRepo.CreateWithWallet(ctx, user, wallet)
	if err != nil {
		return nil, err
	}

	return userWithWal, nil
}

// func (s *userService) Register(ctx context.Context, req uModel.CreateUserRequest) (*uModel.User, error) {
// 	const op = "userService.Register"

// 	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
// 	if existing != nil {
// 		return nil, errors.New("email already registerd")
// 	}

// 	fmt.Println(op, "existing", existing)

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Println(op, "hashed password", hashedPassword)

// 	user := uModel.User{
// 		ID:           uuid.New().String(),
// 		Username:     req.Username,
// 		Email:        req.Email,
// 		PasswordHash: string(hashedPassword),
// 	}

// 	fmt.Println(op, "user", user)

// 	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer tx.Rollback(ctx)

// 	if err = s.userRepo.CreateTx(ctx, tx, &user); err != nil {
// 		return nil, err
// 	}

// 	wallet := walletModel.Wallet{
// 		ID:       uuid.New().String(),
// 		UserID:   user.ID,
// 		Balance:  0.00,
// 		Currency: "RUB",
// 		Status:   "active",
// 	}

// 	if err = s.walletRepo.CreateTx(ctx, tx, &wallet); err != nil {
// 		return nil, err
// 	}

// 	if err := tx.Commit(ctx); err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }

func (s *userService) GetProfile(ctx context.Context, id string) (*uModel.User, error) {
	const op = "userService.GetProfile"

	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) UpdateProfile(ctx context.Context, id string, req uModel.UpdateUserRequest) (*uModel.User, error) {
	const op = "userService.UpdateProfile"

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Username = req.Username
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, req uModel.LoginRequest) (*uModel.LoginResponse, error) {
	const op = "userService.Login"

	fmt.Println(op, "before executing sql query")

	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	fmt.Println(op, user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password hashes aren't the same // wrong email or password")
	}

	accessToken, err := auth.GenerateToken(user.ID, user.Email, accessTokenDefaultDuration)
	if err != nil {
		return nil, err
	}

	fmt.Println(op, "accessToken", accessToken)

	refreshToken, err := auth.GenerateToken(user.ID, user.Email, refreshTokenDefaultDuration)
	if err != nil {
		return nil, err
	}

	fmt.Println(op, "refreshToken", refreshToken)

	return &uModel.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
