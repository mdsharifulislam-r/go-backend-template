package auth

import (
	"net/http"
	"strings"

	apperrors "github.com/mdsharifulislam-r/go-backend-template/internal/errors"
	"github.com/mdsharifulislam-r/go-backend-template/internal/middleware"
	"github.com/mdsharifulislam-r/go-backend-template/internal/response"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) error {
	payload, ok := middleware.BodyFromContext[LoginRequest](r.Context())
	if !ok {
		return apperrors.BadRequest("Invalid request body")
	}

	data, err := c.service.Login(payload)
	if err != nil {
		return err
	}

	response.OK(w, "Login successfully", data)
	return nil
}

func (c *Controller) VerifyOTP(w http.ResponseWriter, r *http.Request) error {
	payload, ok := middleware.BodyFromContext[VerifyEmailRequest](r.Context())
	if !ok {
		return apperrors.BadRequest("Invalid request body")
	}

	data, err := c.service.VerifyOTP(payload)
	if err != nil {
		return err
	}

	response.OK(w, "Email verified successfully", data)
	return nil
}

func (c *Controller) ForgotPassword(w http.ResponseWriter, r *http.Request) error {
	payload, ok := middleware.BodyFromContext[ForgotPasswordRequest](r.Context())
	if !ok {
		return apperrors.BadRequest("Invalid request body")
	}

	data, err := c.service.ForgotPassword(payload)
	if err != nil {
		return err
	}

	response.OK(w, "Email sent successfully", data)
	return nil
}

func (c *Controller) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	payload, ok := middleware.BodyFromContext[ResetPasswordRequest](r.Context())
	if !ok {
		return apperrors.BadRequest("Invalid request body")
	}

	token := strings.TrimSpace(r.Header.Get("Authorization"))
	token = strings.TrimPrefix(token, "Bearer ")

	data, err := c.service.ResetPassword(payload, token)
	if err != nil {
		return err
	}

	response.OK(w, "Password reset successfully", data)
	return nil
}

func (c *Controller) ChangePassword(w http.ResponseWriter, r *http.Request) error {
	claims, ok := middleware.CurrentUser(r.Context())
	if !ok {
		return apperrors.Unauthorized("You are not authorized")
	}

	payload, ok := middleware.BodyFromContext[ChangePasswordRequest](r.Context())
	if !ok {
		return apperrors.BadRequest("Invalid request body")
	}

	data, err := c.service.ChangePassword(payload, claims.ID)
	if err != nil {
		return err
	}

	response.OK(w, "Password changed successfully", data)
	return nil
}
