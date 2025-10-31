package otp_model

import (
	otp_type "authsome-otp/internal/type"

	"gorm.io/gorm"
)

type OTPModel struct {
	gorm.Model
	OTP          string            `gorm:"column:otp;type:varchar(20);not null;index"`
	ParentID     string            `gorm:"column:parent_id;type:varchar(255);not null;index"`
	ParentSource string            `gorm:"column:parent_source;type:varchar(100);not null"`
	Metadata     otp_type.Metadata `gorm:"column:metadata;type:jsonb"`
}

func (OTPModel) TableName() string {
	return "otps"
}
