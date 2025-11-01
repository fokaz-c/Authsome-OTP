package otp_model

import (
	otp_type "authsome-otp/internal/type"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OTPModel struct {
	ID           uint              `gorm:"primarykey" json:"-"`
	IDString     string            `gorm:"-" json:"id"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
	ExpiresAt    int64             `gorm:"column:expires_at;not null;index" json:"expiresAt"`
	OTP          string            `gorm:"column:otp;type:varchar(20);not null;index" json:"code"`
	ParentID     string            `gorm:"column:parent_id;type:varchar(255);not null;index" json:"context"`
	ParentSource string            `gorm:"column:parent_source;type:varchar(100);not null" json:"-"`
	Metadata     otp_type.Metadata `gorm:"column:metadata;type:jsonb" json:"metadata"`
}

func (OTPModel) TableName() string {
	return "otps"
}

func (o *OTPModel) AfterFind(tx *gorm.DB) error {
	o.IDString = fmt.Sprintf("otp-%s", o.generateID())
	return nil
}

func (o *OTPModel) AfterCreate(tx *gorm.DB) error {
	o.IDString = fmt.Sprintf("otp-%s", o.generateID())
	return nil
}

func (o *OTPModel) generateID() string {
	return "otp-" + uuid.NewString()
}
