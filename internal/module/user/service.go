package user

import (
	"time"

	"github.com/mdsharifulislam-r/go-backend-template/internal/email"
	apperrors "github.com/mdsharifulislam-r/go-backend-template/internal/errors"
	"github.com/mdsharifulislam-r/go-backend-template/internal/helper"
	"github.com/mdsharifulislam-r/go-backend-template/internal/models"
	"github.com/mdsharifulislam-r/go-backend-template/internal/querybuilder"
	"github.com/mdsharifulislam-r/go-backend-template/internal/response"
	"gorm.io/gorm"
)

type Service struct {
	db    *gorm.DB
	email *email.Service
}

func NewService(db *gorm.DB, emailSvc *email.Service) *Service {
	return &Service{db: db, email: emailSvc}
}

func (s *Service) Create(req CreateUserRequest) (*models.User, error) {
	var existing models.User
	if err := s.db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, apperrors.BadRequest("User already exist")
	} else if err != gorm.ErrRecordNotFound {
		return nil, apperrors.Internal("Failed to check existing user")
	}

	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.Internal("Failed to hash password")
	}

	role := models.RoleUser
	if req.Role != "" {
		role = models.UserRole(req.Role)
	}

	authJSON, err := helper.AuthToJSON(models.Authentication{IsResetPassword: false})
	if err != nil {
		return nil, apperrors.Internal("Failed to prepare authentication data")
	}

	user := models.User{
		Name:           req.Name,
		Email:          req.Email,
		Password:       hashed,
		Contact:        req.Contact,
		Role:           role,
		Status:         models.StatusActive,
		Verified:       false,
		Authentication: authJSON,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, apperrors.BadRequest("User not created")
	}

	otp := helper.GenerateOTP()
	expireAt := time.Now().UTC().Add(3 * time.Minute)
	authJSON, err = helper.AuthToJSON(models.Authentication{
		IsResetPassword: false,
		OneTimeCode:     &otp,
		ExpireAt:        &expireAt,
	})
	if err != nil {
		return nil, apperrors.Internal("Failed to prepare OTP data")
	}

	if err := s.db.Model(&user).Update("authentication", authJSON).Error; err != nil {
		return nil, apperrors.Internal("Failed to save OTP")
	}

	_ = s.email.Send(email.CreateAccount(user.Name, user.Email, otp))

	user.Password = ""
	return &user, nil
}

func (s *Service) GetProfile(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound("User not found")
		}
		return nil, apperrors.Internal("Failed to fetch profile")
	}
	return &user, nil
}

func (s *Service) UpdateProfile(userID string, req UpdateProfileRequest, imagePath string) (*models.User, error) {
	user, err := s.GetProfile(userID)
	if err != nil {
		return nil, err
	}

	updates := map[string]any{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Contact != "" {
		updates["contact"] = req.Contact
	}
	if imagePath != "" {
		updates["image"] = imagePath
	} else if req.Image != "" {
		updates["image"] = req.Image
	}

	if len(updates) == 0 {
		return nil, apperrors.BadRequest("Nothing to update")
	}

	if err := s.db.Model(user).Updates(updates).Error; err != nil {
		return nil, apperrors.Internal("Failed to update profile")
	}

	return s.GetProfile(userID)
}

func (s *Service) GetAll(query map[string]string) ([]models.User, *response.Pagination, error) {
	var users []models.User
	qb := querybuilder.New(s.db.Model(&models.User{}), query).
		Search("name", "email").
		Filter().
		Sort("-created_at")

	pagination, err := qb.Paginate(&users)
	if err != nil {
		return nil, nil, apperrors.Internal("Failed to fetch users")
	}
	return users, pagination, nil
}

func (s *Service) EnsureSuperAdmin(emailAddr, password string) error {
	var count int64
	if err := s.db.Model(&models.User{}).Where("role = ?", models.RoleSuperAdmin).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hashed, err := helper.HashPassword(password)
	if err != nil {
		return err
	}

	authJSON, err := helper.AuthToJSON(models.Authentication{IsResetPassword: false})
	if err != nil {
		return err
	}

	admin := models.User{
		Name:           "Super Admin",
		Email:          emailAddr,
		Password:       hashed,
		Role:           models.RoleSuperAdmin,
		Status:         models.StatusActive,
		Verified:       true,
		Authentication: authJSON,
	}

	return s.db.Create(&admin).Error
}
