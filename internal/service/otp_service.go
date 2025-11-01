package otp_service

import (
	otp_model "authsome-otp/internal/model"
	otp_repository "authsome-otp/internal/repository"
	otp_type "authsome-otp/internal/type"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type OtpType = int
type Otp = string
type Metadata = otp_type.Metadata

const (
	NUMERIC    OtpType = 1 << 0
	ALPHABETIC OtpType = 1 << 1
)

type Options struct {
	Type             OtpType
	OtpLength        int
	MinNumCount      int
	MinAlphabetCount int
	MaxNumCount      int
	MaxAlphabetCount int
	ExpiresInSeconds int64
}

const (
	numbers   = "0123456789"
	alphabets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type OTPService struct {
	repo otp_repository.OTPRepository
}

func NewOTPService(repo otp_repository.OTPRepository) *OTPService {
	return &OTPService{repo: repo}
}

func (s *OTPService) SaveToDB(ctx context.Context, otp string, parentType, parentID string, metadata Metadata, expiresInSeconds int64) (uint, error) {
	if expiresInSeconds <= 0 {
		expiresInSeconds = 300
	}

	otpModel := &otp_model.OTPModel{
		OTP:          otp,
		ParentID:     parentID,
		ParentSource: parentType,
		Metadata:     metadata,
		ExpiresAt:    time.Now().Unix() + expiresInSeconds,
	}

	if err := s.repo.Create(ctx, otpModel); err != nil {
		return 0, fmt.Errorf("failed to save OTP: %w", err)
	}

	return otpModel.ID, nil
}

func (s *OTPService) GetByID(ctx context.Context, id uint) (*otp_model.OTPModel, error) {
	otpModel, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if time.Now().Unix() > otpModel.ExpiresAt {
		s.repo.Delete(ctx, id)
		return nil, fmt.Errorf("OTP expired")
	}

	return otpModel, nil
}

func (s *OTPService) ValidateOTP(ctx context.Context, token uint, otp, expectedParentType, expectedParentID string) (*otp_model.OTPModel, error) {
	fetchedOTP, err := s.GetByID(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("OTP not found or expired")
	}

	if fetchedOTP.OTP != otp {
		return nil, fmt.Errorf("OTP mismatch")
	}

	if fetchedOTP.ParentSource != expectedParentType || fetchedOTP.ParentID != expectedParentID {
		return nil, fmt.Errorf("invalid parent mismatch")
	}

	return fetchedOTP, nil
}

func (s *OTPService) DeleteExpiredOTPs(ctx context.Context) error {
	currentTime := time.Now().Unix()
	return s.repo.DeleteExpiredByTimestamp(ctx, currentTime)
}

func GenerateOTP(op Options) (string, error) {
	if op.OtpLength <= 0 {
		op.OtpLength = 6
	}
	if op.MinNumCount == 0 {
		op.MinNumCount = 4
	}
	if op.MaxNumCount == 0 {
		op.MaxNumCount = 6
	}
	if op.MinAlphabetCount == 0 {
		op.MinAlphabetCount = 4
	}
	if op.MaxAlphabetCount == 0 {
		op.MaxAlphabetCount = 6
	}

	if op.Type&(NUMERIC|ALPHABETIC) == 0 {
		return "", fmt.Errorf("no valid OTP type flags provided")
	}

	numCount := 0
	alphaCount := 0

	if op.Type&NUMERIC != 0 {
		numCount = clamp(op.MinNumCount, 0, min(op.MaxNumCount, op.OtpLength))
	}

	if op.Type&ALPHABETIC != 0 {
		remaining := op.OtpLength - numCount
		alphaCount = clamp(op.MinAlphabetCount, 0, min(op.MaxAlphabetCount, remaining))
	}

	if numCount+alphaCount > op.OtpLength {
		return "", fmt.Errorf("minimum character requirements exceed OTP length")
	}

	otp := make([]byte, 0, op.OtpLength)

	if numCount > 0 {
		chars, err := randChars(numbers, numCount)
		if err != nil {
			return "", fmt.Errorf("failed to generate numeric characters: %w", err)
		}
		otp = append(otp, chars...)
	}

	if alphaCount > 0 {
		chars, err := randChars(alphabets, alphaCount)
		if err != nil {
			return "", fmt.Errorf("failed to generate alphabetic characters: %w", err)
		}
		otp = append(otp, chars...)
	}

	remaining := op.OtpLength - len(otp)
	if remaining > 0 {
		charset := buildCharset(op.Type)
		chars, err := randChars(charset, remaining)
		if err != nil {
			return "", fmt.Errorf("failed to generate remaining characters: %w", err)
		}
		otp = append(otp, chars...)
	}

	if err := shuffle(otp); err != nil {
		return "", fmt.Errorf("failed to shuffle OTP: %w", err)
	}

	return string(otp), nil
}

func buildCharset(otpType OtpType) string {
	var charset string
	if otpType&NUMERIC != 0 {
		charset += numbers
	}
	if otpType&ALPHABETIC != 0 {
		charset += alphabets
	}
	return charset
}

func randChars(source string, n int) ([]byte, error) {
	if len(source) == 0 || n <= 0 {
		return []byte{}, nil
	}

	result := make([]byte, n)
	sourceLen := big.NewInt(int64(len(source)))

	for i := range n {
		idx, err := rand.Int(rand.Reader, sourceLen)
		if err != nil {
			return nil, err
		}
		result[i] = source[idx.Int64()]
	}

	return result, nil
}

func shuffle(slice []byte) error {
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return err
		}
		jIdx := j.Int64()
		slice[i], slice[jIdx] = slice[jIdx], slice[i]
	}
	return nil
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
