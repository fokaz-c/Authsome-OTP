package otp_repository

import (
	otp_model "authsome-otp/internal/model"
	"context"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *otp_model.OTPModel) error
	FindByID(ctx context.Context, id uint) (*otp_model.OTPModel, error)
	FindByOTP(ctx context.Context, otpCode string) (*otp_model.OTPModel, error)
	FindByParentID(ctx context.Context, parentID string) ([]*otp_model.OTPModel, error)
	FindByParentIDAndSource(ctx context.Context, parentID, parentSource string) ([]*otp_model.OTPModel, error)
	Update(ctx context.Context, otp *otp_model.OTPModel) error
	Delete(ctx context.Context, id uint) error
	DeleteByParentID(ctx context.Context, parentID string) error
	DeleteExpiredByTimestamp(ctx context.Context, currentTimestamp int64) error
}
