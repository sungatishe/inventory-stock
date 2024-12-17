package service

import (
	"auth-service/api/proto"
	"auth-service/internal/models"
	"auth-service/utils"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type AuthService struct {
	repo UserRepository
	proto.UnimplementedAuthServiceServer
}

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	existingUser, _ := s.repo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	err = s.repo.CreateUser(ctx, user)

	return &proto.RegisterResponse{Message: "User created successfully"}, nil
}

func (s *AuthService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials user is nil")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials with password")
	}

	token, err := utils.GenerateJWT(user.ID, "admin")
	if err != nil {
		return nil, err
	}

	return &proto.LoginResponse{Token: token}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	tokenString := req.GetToken()

	_, role, err := utils.ParseJwt(tokenString)
	if err != nil {
		log.Printf("Error parsing JWT: %v", err)
		return &proto.ValidateTokenResponse{
			Valid: false,
			Role:  "",
		}, nil
	}

	return &proto.ValidateTokenResponse{
		Valid: true,
		Role:  role,
	}, nil
}
