package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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
	emailService      *EmailService
	jwtSecret         string
	frontendURL       string
}

func NewAuthService(
	userRepo repository.UserRepository,
	passwordResetRepo repository.PasswordResetTokenRepository,
	emailService *EmailService,
	jwtSecret string,
	frontendURL string,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		emailService:      emailService,
		jwtSecret:         jwtSecret,
		frontendURL:       frontendURL,
	}
}

type AuthPayload struct {
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

	token, err := utils.GenerateJWT(user.ID.String(), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthPayload{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(email, password string) (*AuthPayload, error) {
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

	token, err := utils.GenerateJWT(user.ID.String(), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthPayload{
		Token: token,
		User:  user,
	}, nil
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
