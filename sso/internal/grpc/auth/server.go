package auth

import (
	"context"
	ssov1 "github.com/Karaulkin/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, UserID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

const (
	EmptyValue  = 0
	EmptyString = ""
)

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (
	*ssov1.LoginResponce, error) {

	if err := validateLogin(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		// TODO: обработка в зависимости от ошибки (2:31:16)
		/*
			// Ошибку auth.ErrInvalidCredentials мы создадим ниже
			        if errors.Is(err, auth.ErrInvalidCredentials) {
			            return nil, status.Error(codes.InvalidArgument, "invalid email or password")
			        }
		*/
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponce{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		/*
			// TODO: обработка в зависимости от ошибки
				if errors.Is(err, storage.ErrUserExists) {
					return nil, status.Error(codes.FailedPrecondition, "user already exists")
				}
		*/
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponce, error) {

	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO:
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponce{
		IsAdmin: isAdmin,
	}, nil

}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == EmptyString {
		return status.Error(codes.InvalidArgument, "missing email")
	}

	if req.GetPassword() == EmptyString {
		return status.Error(codes.InvalidArgument, "missing password")
	}

	if req.GetAppId() == EmptyValue {
		return status.Error(codes.InvalidArgument, "missing app id")
	}

	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == EmptyString {
		return status.Error(codes.InvalidArgument, "missing email")
	}

	if req.GetPassword() == EmptyString {
		return status.Error(codes.InvalidArgument, "missing password")
	}

	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == EmptyValue {
		return status.Error(codes.InvalidArgument, "missing user id")
	}
	return nil
}
