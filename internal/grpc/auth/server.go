package auth

import (
	"context"
	ssov1 "github.com/IvanSukhinin/sso-proto/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Auth interface {
	Register(ctx context.Context, email string, password string) error
	Login(ctx context.Context, email string, password string, appId int64) (token string, err error)
	IsAdmin(ctx context.Context, userId uuid.UUID) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

type RequestValidate struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type LoginRequestValidate struct {
	RequestValidate
	AppId int64 `validate:"number"`
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*emptypb.Empty, error) {
	res := RequestValidate{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	if err := validator.New().Struct(res); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.auth.Register(ctx, res.Email, res.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	res := LoginRequestValidate{
		RequestValidate: RequestValidate{
			Email:    req.GetEmail(),
			Password: req.GetPassword(),
		},
		AppId: req.GetAppId(),
	}
	if err := validator.New().Struct(res); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, res.Email, res.Password, res.AppId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	isAdmin, err := s.auth.IsAdmin(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
