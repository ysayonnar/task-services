package handlers

import (
	"context"
	"errors"
	"log/slog"
	"sso/internal/storage"
	"sso/internal/utils"

	sso "github.com/ysayonnar/task-contracts/sso/gen/go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SsoServer struct {
	sso.UnimplementedAuthServiceServer
	Log     *slog.Logger
	Storage *storage.Storage
}

func (s *SsoServer) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	const op = "handlers.Register"
	log := s.Log.With(slog.String("op", op))

	if !utils.IsEmailValid(req.Email) {
		return nil, status.Error(codes.InvalidArgument, "email is invalid")
	}

	if len(req.Password) > 72 || len(req.Password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password length is invalid")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error while hashing password", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	userId, err := s.Storage.InsertUser(ctx, req.Email, string(passwordHash))
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		log.Error("error while inserting user", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &sso.RegisterResponse{UserId: userId}, nil
}

func (s *SsoServer) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	const op = "handlers.Login"
	log := s.Log.With(slog.String("op", op))

	if !utils.IsEmailValid(req.Email) {
		return nil, status.Error(codes.InvalidArgument, "email is invalid")
	}

	if len(req.Password) > 72 || len(req.Password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password length is invalid")
	}

	user, err := s.Storage.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user with this email was not found")
		}

		log.Error("error while finding user", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "password is invalid")
	}

	//TODO: jwt token

	return nil, nil
}

// TODO: implement
func (s *SsoServer) Delete(ctx context.Context, req *sso.DeleteRequest) (*sso.DeleteResponse, error) {
	return nil, nil
}
