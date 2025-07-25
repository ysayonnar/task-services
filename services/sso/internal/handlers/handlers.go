package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sso/internal/config"
	"sso/internal/queue"
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
	Cfg     *config.Config
	Broker  *queue.Broker
}

func (s *SsoServer) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	const op = "handlers.Register"
	log := s.Log.With(slog.String("op", op))

	if !utils.IsEmailValid(req.GetEmail()) {
		return nil, status.Error(codes.InvalidArgument, "email is invalid")
	}

	if len(req.GetPassword()) > 72 || len(req.GetPassword()) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password length is invalid")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error while hashing password", "error", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	userId, err := s.Storage.InsertUser(ctx, req.GetEmail(), string(passwordHash))
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		log.Error("error while inserting user", "error", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	data, err := json.Marshal(queue.UserCreatedDto{UserId: userId, Email: req.GetEmail()})

	err = s.Broker.Publish(ctx, queue.KEY_USER_REGISTERED, data)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &sso.RegisterResponse{UserId: userId}, nil
}

func (s *SsoServer) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	const op = "handlers.Login"
	log := s.Log.With(slog.String("op", op))

	if !utils.IsEmailValid(req.GetEmail()) {
		return nil, status.Error(codes.InvalidArgument, "email is invalid")
	}

	if len(req.GetPassword()) > 72 || len(req.GetPassword()) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password length is invalid")
	}

	user, err := s.Storage.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user with this email was not found")
		}

		log.Error("error while finding user", "error", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.GetPassword()))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "password is invalid")
	}

	token, err := utils.GenerateJwt(user.UserId, s.Cfg.Secret, s.Cfg.TokenTTL)
	if err != nil {
		log.Error("error while generating jwt", "error", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &sso.LoginResponse{Token: token}, nil
}

func (s *SsoServer) Delete(ctx context.Context, req *sso.DeleteRequest) (*sso.DeleteResponse, error) {
	const op = "handlers.Delete"
	log := s.Log.With(slog.String("op", op))

	if !utils.IsEmailValid(req.GetEmail()) {
		return nil, status.Error(codes.InvalidArgument, "email is invalid")
	}

	if len(req.GetPassword()) > 72 || len(req.GetPassword()) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password length is invalid")
	}

	user, err := s.Storage.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user with this email was not found")
		}

		log.Error("error while finding user", "error", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.GetPassword()))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "password is invalid")
	}

	_, err = s.Storage.DeleteUser(ctx, req.GetEmail())
	if err != nil {
		// NOTE: no checking for storage.ErrUserNotFound, it was at line 99
		log.Error("error while deleting user", "error", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	data, err := json.Marshal(queue.UserDeletedDto{UserId: user.UserId})

	err = s.Broker.Publish(ctx, queue.KEY_USER_DELETED, data)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &sso.DeleteResponse{IsDeleted: true}, nil
}
