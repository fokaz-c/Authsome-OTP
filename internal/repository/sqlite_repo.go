package otp_repository

import (
	otp_model "authsome-otp/internal/model"
	"context"
	"time"

	"gorm.io/gorm"
)

type gormOTPRepository struct {
	db *gorm.DB
}

func NewGormOTPRepository(db *gorm.DB) OTPRepository {
	return &gormOTPRepository{db: db}
}

func (r *gormOTPRepository) Create(ctx context.Context, otp *otp_model.OTPModel) error {
	return r.db.WithContext(ctx).Create(otp).Error
}

func (r *gormOTPRepository) FindByID(ctx context.Context, id uint) (*otp_model.OTPModel, error) {
	var otp otp_model.OTPModel
	err := r.db.WithContext(ctx).First(&otp, id).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *gormOTPRepository) FindByOTP(ctx context.Context, otpCode string) (*otp_model.OTPModel, error) {
	var otp otp_model.OTPModel
	err := r.db.WithContext(ctx).Where("otp = ?", otpCode).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *gormOTPRepository) FindByParentID(ctx context.Context, parentID string) ([]*otp_model.OTPModel, error) {
	var otps []*otp_model.OTPModel
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&otps).Error
	if err != nil {
		return nil, err
	}
	return otps, nil
}

func (r *gormOTPRepository) FindByParentIDAndSource(ctx context.Context, parentID, parentSource string) ([]*otp_model.OTPModel, error) {
	var otps []*otp_model.OTPModel
	err := r.db.WithContext(ctx).
		Where("parent_id = ? AND parent_source = ?", parentID, parentSource).
		Find(&otps).Error
	if err != nil {
		return nil, err
	}
	return otps, nil
}

func (r *gormOTPRepository) FindActiveByParentID(ctx context.Context, parentID string, expiryTime time.Time) (*otp_model.OTPModel, error) {
	var otp otp_model.OTPModel
	err := r.db.WithContext(ctx).
		Where("parent_id = ? AND created_at > ?", parentID, expiryTime).
		Order("created_at DESC").
		First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *gormOTPRepository) Update(ctx context.Context, otp *otp_model.OTPModel) error {
	return r.db.WithContext(ctx).Save(otp).Error
}

func (r *gormOTPRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&otp_model.OTPModel{}, id).Error
}

func (r *gormOTPRepository) DeleteByParentID(ctx context.Context, parentID string) error {
	return r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Delete(&otp_model.OTPModel{}).Error
}

func (r *gormOTPRepository) DeleteExpired(ctx context.Context, expiryTime time.Time) error {
	return r.db.WithContext(ctx).
		Where("created_at < ?", expiryTime).
		Delete(&otp_model.OTPModel{}).Error
}
