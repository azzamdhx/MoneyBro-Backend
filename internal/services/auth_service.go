package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	mathrand "math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/azzamdhx/moneybro/backend/internal/models"
	"github.com/azzamdhx/moneybro/backend/internal/repository"
	"github.com/azzamdhx/moneybro/backend/internal/utils"
)

type AuthService struct {
	userRepo          repository.UserRepository
	passwordResetRepo repository.PasswordResetTokenRepository
	twoFACodeRepo     repository.TwoFACodeRepository
	emailService      *EmailService
	jwtSecret         string
	frontendURL       string
	accountService    *AccountService
}

func NewAuthService(
	userRepo repository.UserRepository,
	passwordResetRepo repository.PasswordResetTokenRepository,
	twoFACodeRepo repository.TwoFACodeRepository,
	emailService *EmailService,
	jwtSecret string,
	frontendURL string,
	accountService *AccountService,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		twoFACodeRepo:     twoFACodeRepo,
		emailService:      emailService,
		jwtSecret:         jwtSecret,
		frontendURL:       frontendURL,
		accountService:    accountService,
	}
}

type AuthPayload struct {
	Token       string
	User        *models.User
	Requires2FA bool
	TempToken   string
}

type TwoFAPayload struct {
	Token string
	User  *models.User
}

func (s *AuthService) Register(email, password, name string) (*AuthPayload, error) {
	if err := utils.ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := utils.ValidatePassword(password); err != nil {
		return nil, err
	}
	if err := utils.ValidateName(name); err != nil {
		return nil, err
	}

	existing, err := s.userRepo.GetByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		Name:         name,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Create default account for user
	if _, err := s.accountService.CreateDefaultAccount(user.ID); err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWT(user.ID.String(), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthPayload{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthPayload, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// Check if 2FA is enabled
	if user.TwoFAEnabled {
		// Generate temp token for 2FA verification
		tempToken, err := utils.GenerateTempToken(user.ID.String(), s.jwtSecret)
		if err != nil {
			return nil, err
		}

		// Generate and send 2FA code
		if err := s.send2FACode(ctx, user); err != nil {
			return nil, err
		}

		return &AuthPayload{
			Requires2FA: true,
			TempToken:   tempToken,
		}, nil
	}

	token, err := utils.GenerateJWT(user.ID.String(), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthPayload{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) send2FACode(ctx context.Context, user *models.User) error {
	// Delete any existing codes for this user
	_ = s.twoFACodeRepo.DeleteByUserID(user.ID)

	// Generate 6-digit code
	code := fmt.Sprintf("%06d", mathrand.Intn(1000000))

	// Create 2FA code (expires in 10 minutes)
	twoFACode := &models.TwoFACode{
		ID:        uuid.New(),
		UserID:    user.ID,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := s.twoFACodeRepo.Create(twoFACode); err != nil {
		return err
	}

	// Send email with code
	return s.emailService.Send2FACodeEmail(ctx, user.Email, user.Name, code)
}

func (s *AuthService) Verify2FA(ctx context.Context, tempToken, code string) (*TwoFAPayload, error) {
	// Verify temp token
	userID, err := utils.VerifyTempToken(tempToken, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Verify 2FA code
	twoFACode, err := s.twoFACodeRepo.GetValidByUserIDAndCode(userUUID, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kode verifikasi tidak valid atau sudah kadaluarsa")
		}
		return nil, err
	}

	// Mark code as used
	if err := s.twoFACodeRepo.MarkAsUsed(twoFACode.ID); err != nil {
		return nil, err
	}

	// Get user
	user, err := s.userRepo.GetByID(userUUID)
	if err != nil {
		return nil, err
	}

	// Generate actual JWT token
	token, err := utils.GenerateJWT(user.ID.String(), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &TwoFAPayload{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Resend2FACode(ctx context.Context, tempToken string) error {
	// Verify temp token
	userID, err := utils.VerifyTempToken(tempToken, s.jwtSecret)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid token")
	}

	user, err := s.userRepo.GetByID(userUUID)
	if err != nil {
		return err
	}

	return s.send2FACode(ctx, user)
}

func (s *AuthService) Enable2FA(userID uuid.UUID, password string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify password
	if !utils.CheckPassword(password, user.PasswordHash) {
		return errors.New("password tidak valid")
	}

	user.TwoFAEnabled = true
	now := time.Now()
	user.UpdatedAt = &now

	return s.userRepo.Update(user)
}

func (s *AuthService) Disable2FA(userID uuid.UUID, password string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify password
	if !utils.CheckPassword(password, user.PasswordHash) {
		return errors.New("password tidak valid")
	}

	user.TwoFAEnabled = false
	now := time.Now()
	user.UpdatedAt = &now

	return s.userRepo.Update(user)
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	if err := utils.ValidateEmail(email); err != nil {
		return err
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal if email exists
			return nil
		}
		return err
	}

	// Delete any existing tokens for this user
	_ = s.passwordResetRepo.DeleteByUserID(user.ID)

	// Generate secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return err
	}
	token := hex.EncodeToString(tokenBytes)

	// Create password reset token (expires in 1 hour)
	resetToken := &models.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.passwordResetRepo.Create(resetToken); err != nil {
		return err
	}

	// Send email with reset link
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, token)
	if err := s.emailService.SendPasswordResetEmail(ctx, user.Email, user.Name, resetLink); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	if err := utils.ValidatePassword(newPassword); err != nil {
		return err
	}

	resetToken, err := s.passwordResetRepo.GetValidByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired reset token")
		}
		return err
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update user password
	user := &resetToken.User
	user.PasswordHash = hashedPassword
	now := time.Now()
	user.UpdatedAt = &now

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Mark token as used
	if err := s.passwordResetRepo.MarkAsUsed(resetToken.ID); err != nil {
		return err
	}

	return nil
}
