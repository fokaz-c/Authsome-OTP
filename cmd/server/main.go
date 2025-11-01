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

func demoGenerateAndSaveOtp(ctx context.Context, service *otp_service.OTPService) (uint, string, string, string, error) {
	log.Println("=========================================================")
	log.Println("Function 1: demoGenerateAndSaveOtp")
	log.Println("=========================================================")
	log.Println("This function generates a random OTP, stores it with metadata,")
	log.Println("and returns the full, saved OTP structure.")
	log.Println()

	options := otp_service.Options{
		Type:             otp_service.NUMERIC,
		OtpLength:        6,
		ExpiresInSeconds: 300,
	}
	metadata := otp_type.Metadata{
		TenantUsername:   "test_user_signup",
		IdentityType:     "EMAIL",
		IdentitySource:   "user@example.com",
		IdentityPassword: "password",
	}
	parentID := "user_abc_123"
	parentSource := "TENANT_SIGNUP"

	log.Println("[Input] Options: 6-digit NUMERIC, expires in 300s")
	log.Println("[Input] Metadata:", prettyPrint(metadata))
	log.Println("[Input] ParentID:", parentID)
	log.Println("[Input] ParentSource:", parentSource)
	log.Println()
	log.Println("--- Processing ---")

	otpString, err := otp_service.GenerateOTP(options)
	if err != nil {
		log.Printf("Failed to generate OTP: %v", err)
		return 0, "", "", "", err
	}
	log.Printf("Step 1: Generated secure OTP string: %s", otpString)

	otpTokenID, err := service.SaveToDB(ctx, otpString, parentSource, parentID, metadata, options.ExpiresInSeconds)
	if err != nil {
		log.Printf("Failed to save OTP: %v", err)
		return 0, "", "", "", err
	}
	log.Printf("Step 2: Saved to database. Received new token ID: %d", otpTokenID)

	savedModel, err := service.GetByID(ctx, otpTokenID)
	if err != nil {
		log.Printf("Failed to fetch saved model: %v", err)
		return 0, "", "", "", err
	}

	log.Println()
	log.Println("[Output] FetchedOtp Structure (as *otp_model.OTPModel):")
	log.Printf("\n%s\n", prettyPrint(savedModel))

	return savedModel.ID, otpString, parentID, parentSource, nil
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

func demoValidateOtp(ctx context.Context, service *otp_service.OTPService, id uint, otpString, parentID, parentSource string) {
	log.Println("=========================================================")
	log.Println("Function 3: demoValidateOtp")
	log.Println("=========================================================")
	log.Println("This function demonstrates validating an OTP with correct and incorrect inputs.")
	log.Println()

	log.Println("--- Case 1: Success (Correct Token, OTP, and Parent Info) ---")
	log.Printf("[Input] Token: %d, OTP: %s, ParentID: %s, ParentSource: %s\n", id, otpString, parentID, parentSource)
	log.Println("--- Processing ---")
	validatedModel, err := service.ValidateOTP(ctx, id, otpString, parentSource, parentID)
	if err != nil {
		log.Printf("[Output] Error: %v\n", err)
	} else {
		log.Println("[Output] Success! Validated model metadata:")
		log.Printf("\n%s\n", prettyPrint(validatedModel.Metadata))
	}
	log.Println()

	log.Println("--- Case 2: Failure (Correct Token, Bad OTP) ---")
	log.Printf("[Input] Token: %d, OTP: %s, ParentID: %s, ParentSource: %s\n", id, "BAD-OTP", parentID, parentSource)
	log.Println("--- Processing ---")
	_, err = service.ValidateOTP(ctx, id, "BAD-OTP", parentSource, parentID)
	if err != nil {
		log.Printf("[Output] Error: %v (as expected)\n", err)
	} else {
		log.Println("[Output] Success! (This should not happen)")
	}
	log.Println()

	log.Println("--- Case 3: Failure (Correct Token, Correct OTP, Bad ParentID) ---")
	log.Printf("[Input] Token: %d, OTP: %s, ParentID: %s, ParentSource: %s\n", id, otpString, "WRONG-PARENT-ID", parentSource)
	log.Println("--- Processing ---")
	_, err = service.ValidateOTP(ctx, id, otpString, parentSource, "WRONG-PARENT-ID")
	if err != nil {
		log.Printf("[Output] Error: %v (as expected)\n", err)
	} else {
		log.Println("[Output] Success! (This should not happen)")
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

	createdID, otpString, parentID, parentSource, err := demoGenerateAndSaveOtp(ctx, service)
	if err != nil {
		log.Fatalf("Demo 1 Failed: %v", err)
	}

	log.Println()
	log.Println()

	demoGetOtpById(ctx, service, createdID)

	log.Println()
	log.Println()

	demoValidateOtp(ctx, service, createdID, otpString, parentID, parentSource)

	log.Println()
	log.Println("=========================================================")
	log.Println("All demos complete.")
	log.Println("=========================================================")
}
