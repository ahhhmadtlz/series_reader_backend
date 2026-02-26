package userhandler

import (
    "context"
    userparam "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
    uservalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/validator"
)

type UserService interface {
    Register(ctx context.Context, req userparam.RegisterRequest) (userparam.RegisterResponse, error)
    Login(ctx context.Context, req userparam.LoginRequest) (userparam.LoginResponse, error)
    RefreshAccessToken(ctx context.Context, req userparam.RefreshAccessTokenRequest) (userparam.RefreshAccessTokenResponse, error)
    GetProfile(ctx context.Context, userID uint) (userparam.GetProfileResponse, error)
    UpdateProfile(ctx context.Context, userID uint, req userparam.UpdateProfileRequest) (userparam.UpdateProfileResponse, error)
    ChangePassword(ctx context.Context, userID uint, req userparam.ChangePasswordRequest) (userparam.ChangePasswordResponse, error)
}

type Handler struct {
    service   UserService
    validator uservalidator.Validator
}

func New(service UserService, validator uservalidator.Validator) Handler {
    return Handler{
        service:   service,
        validator: validator,
    }
}