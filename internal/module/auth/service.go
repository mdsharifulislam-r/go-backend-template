package auth

import (
	"time"

	"github.com/mdsharifulislam-r/go-backend-template/internal/email"
	apperrors "github.com/mdsharifulislam-r/go-backend-template/internal/errors"
	"github.com/mdsharifulislam-r/go-backend-template/internal/helper"
	"github.com/mdsharifulislam-r/go-backend-template/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	db    *gorm.DB
	email *email.Service
}

func NewService(db *gorm.DB, emailSvc *email.Service) *Service {
	return &Service{db: db, email: emailSvc}
}

func (s *Service) findUserByEmail(emailAddr string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", emailAddr).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound("User not found")
		}
		return nil, apperrors.Internal("Failed to fetch user")
	}
	return &user, nil
}

func (s *Service) findUserByID(id string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound("User not found")
		}
		return nil, apperrors.Internal("Failed to fetch user")
	}
	return &user, nil
}

func (s *Service) Login(req LoginRequest) (map[string]any, error) {
	user, err := s.findUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if !user.Verified {
		return nil, apperrors.BadRequest("Email not verified")
	}
	if user.Status == models.StatusDelete {
		return nil, apperrors.BadRequest("Your account has been deleted or deactivated")
	}
	if !helper.ComparePassword(req.Password, user.Password) {
		return nil, apperrors.BadRequest("Invalid credentials")
	}

	token, err := helper.SignToken(user)
	if err != nil {
		return nil, apperrors.Internal("Failed to generate token")
	}

	return map[string]any{
		"accessToken": token,
		"role":        user.Role,
	}, nil
}

func (s *Service) VerifyOTP(req VerifyEmailRequest) (any, error) {
	user, err := s.findUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	auth := helper.AuthFromJSON(user.Authentication)
	if auth.OneTimeCode == nil || *auth.OneTimeCode != req.OneTimeCode {
		return nil, apperrors.BadRequest("Invalid one time code")
	}
	if auth.ExpireAt == nil || auth.ExpireAt.Before(time.Now().UTC()) {
		return nil, apperrors.BadRequest("One time code expired")
	}

	if !auth.IsResetPassword {
		cleared, err := helper.AuthToJSON(models.Authentication{IsResetPassword: false})
		if err != nil {
			return nil, apperrors.Internal("Failed to clear OTP")
		}
		if err := s.db.Model(user).Updates(map[string]any{
			"verified":       true,
			"authentication": cleared,
		}).Error; err != nil {
			return nil, apperrors.Internal("Failed to verify email")
		}

		user.Verified = true
		user.Password = ""
		return user, nil
	}

	token, err := helper.CryptoToken(32)
	if err != nil {
		return nil, apperrors.Internal("Failed to create reset token")
	}

	if err := s.db.Create(&models.ResetToken{Token: token, UserID: user.ID}).Error; err != nil {
		return nil, apperrors.Internal("Failed to save reset token")
	}

	return map[string]any{"token": token}, nil
}

func (s *Service) ForgotPassword(req ForgotPasswordRequest) (map[string]any, error) {
	user, err := s.findUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user.Status == models.StatusDelete {
		return nil, apperrors.BadRequest("Your account has been deleted or deactivated")
	}

	otp := helper.GenerateOTP()
	expireAt := time.Now().UTC().Add(3 * time.Minute)
	authJSON, err := helper.AuthToJSON(models.Authentication{
		IsResetPassword: true,
		OneTimeCode:     &otp,
		ExpireAt:        &expireAt,
	})
	if err != nil {
		return nil, apperrors.Internal("Failed to prepare OTP")
	}

	if err := s.db.Model(user).Update("authentication", authJSON).Error; err != nil {
		return nil, apperrors.Internal("Failed to save OTP")
	}

	_ = s.email.Send(email.ResetPassword(user.Email, otp))

	return map[string]any{"email": user.Email}, nil
}

func (s *Service) ResetPassword(req ResetPasswordRequest, token string) (map[string]any, error) {
	if token == "" {
		return nil, apperrors.Unauthorized("Reset token is required")
	}

	var resetToken models.ResetToken
	if err := s.db.Where("token = ?", token).First(&resetToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound("Reset token not found")
		}
		return nil, apperrors.Internal("Failed to fetch reset token")
	}

	user, err := s.findUserByID(resetToken.UserID)
	if err != nil {
		return nil, err
	}
	if user.Status == models.StatusDelete {
		return nil, apperrors.BadRequest("Your account has been deleted or deactivated")
	}
	if helper.ComparePassword(req.NewPassword, user.Password) {
		return nil, apperrors.BadRequest("New password and old password are same")
	}

	hashed, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return nil, apperrors.Internal("Failed to hash password")
	}

	cleared, err := helper.AuthToJSON(models.Authentication{IsResetPassword: false})
	if err != nil {
		return nil, apperrors.Internal("Failed to clear authentication")
	}

	if err := s.db.Model(user).Updates(map[string]any{
		"password":       hashed,
		"authentication": cleared,
	}).Error; err != nil {
		return nil, apperrors.Internal("Failed to reset password")
	}

	_ = s.db.Delete(&resetToken)

	return map[string]any{"email": user.Email}, nil
}

func (s *Service) ChangePassword(req ChangePasswordRequest, userID string) (map[string]any, error) {
	user, err := s.findUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user.Status == models.StatusDelete {
		return nil, apperrors.BadRequest("Your account has been deleted or deactivated")
	}
	if !helper.ComparePassword(req.CurrentPassword, user.Password) {
		return nil, apperrors.BadRequest("Invalid credentials")
	}

	hashed, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return nil, apperrors.Internal("Failed to hash password")
	}

	if err := s.db.Model(user).Update("password", hashed).Error; err != nil {
		return nil, apperrors.Internal("Failed to change password")
	}

	return map[string]any{"email": user.Email}, nil
}
