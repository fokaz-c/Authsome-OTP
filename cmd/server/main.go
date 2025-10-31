package main

import (
	otp_model "authsome-otp/internal/model"
	otp_repository "authsome-otp/internal/repository"
	otp_service "authsome-otp/internal/service"
	otp_type "authsome-otp/internal/type"
	"context"
	"encoding/json"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func prettyPrint(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func demoGenerateAndSaveOtp(ctx context.Context, service *otp_service.OTPService) (uint, error) {
	log.Println("=========================================================")
	log.Println("Function 1: demoGenerateAndSaveOtp")
	log.Println("=========================================================")
	log.Println("This function generates a random OTP, stores it with metadata,")
	log.Println("and returns the full, saved OTP structure.")
	log.Println()

	options := otp_service.Options{
		Type:      otp_service.NUMERIC,
		OtpLength: 6,
	}
	metadata := otp_type.Metadata{
		TenantUsername: "test_user_signup",
		IdentityType:   "EMAIL",
		IdentitySource: "user@example.com",
	}
	parentID := "user_abc_123"
	parentSource := "TENANT_SIGNUP"

	log.Println("[Input] Options: 6-digit NUMERIC")
	log.Println("[Input] Metadata:", prettyPrint(metadata))
	log.Println("[Input] ParentID:", parentID)
	log.Println("[Input] ParentSource:", parentSource)
	log.Println()
	log.Println("--- Processing ---")

	otpString, err := otp_service.GenerateOTP(options)
	if err != nil {
		log.Printf("Failed to generate OTP: %v", err)
		return 0, err
	}
	log.Printf("Step 1: Generated secure OTP string: %s", otpString)

	otpTokenID, err := service.SaveToDB(ctx, otpString, parentSource, parentID, metadata)
	if err != nil {
		log.Printf("Failed to save OTP: %v", err)
		return 0, err
	}
	log.Printf("Step 2: Saved to database. Received new token ID: %d", otpTokenID)

	savedModel, err := service.GetByID(ctx, otpTokenID)
	if err != nil {
		log.Printf("Failed to fetch saved model: %v", err)
		return 0, err
	}

	log.Println()
	log.Println("[Output] FetchedOtp Structure (as *otp_model.OTPModel):")
	log.Printf("\n%s\n", prettyPrint(savedModel))

	return savedModel.ID, nil
}

func demoGetOtpById(ctx context.Context, service *otp_service.OTPService, id uint) {
	log.Println("=========================================================")
	log.Println("Function 2: demoGetOtpById")
	log.Println("=========================================================")
	log.Println("This function retrieves a stored OTP by its unique ID (token).")
	log.Println()

	log.Println("--- Case 1: Success (ID exists) ---")
	log.Printf("[Input] ID: %d", id)
	log.Println()
	log.Println("--- Processing ---")

	fetchedOTP, err := service.GetByID(ctx, id)
	if err != nil {
		log.Printf("[Output] Error: %v\n", err)
	} else {
		log.Println("[Output] Returns: FetchedOtp structure:")
		log.Printf("\n%s\n", prettyPrint(fetchedOTP))
	}

	log.Println()
	log.Println("--- Case 2: Failure (ID does not exist) ---")
	log.Printf("[Input] ID: %d", 999)
	log.Println()
	log.Println("--- Processing ---")

	fetchedOTP, err = service.GetByID(ctx, 999)
	if err != nil {
		log.Printf("[Output] Returns: %v (as expected)\n", err)
	} else {
		log.Println("[Output] Returns: FetchedOtp structure (This should not happen):")
		log.Printf("\n%s\n", prettyPrint(fetchedOTP))
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.Exec("DROP TABLE IF EXISTS otps")
	log.Println("Cleaned old 'otps' table.")

	err = db.AutoMigrate(&otp_model.OTPModel{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Database connection and migration successful.")
	log.Println()

	repo := otp_repository.NewGormOTPRepository(db)
	service := otp_service.NewOTPService(repo)
	ctx := context.Background()

	log.Println("Starting service demos...")

	createdID, err := demoGenerateAndSaveOtp(ctx, service)
	if err != nil {
		log.Fatalf("Demo 1 Failed: %v", err)
	}

	log.Println()
	log.Println()

	demoGetOtpById(ctx, service, createdID)

	log.Println()
	log.Println("=========================================================")
	log.Println("All demos complete.")
	log.Println("=========================================================")
}
