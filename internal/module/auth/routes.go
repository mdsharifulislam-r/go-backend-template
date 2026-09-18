package auth

import (
	"github.com/go-chi/chi/v5"
	apperrors "github.com/mdsharifulislam-r/go-backend-template/internal/errors"
	"github.com/mdsharifulislam-r/go-backend-template/internal/middleware"
)

func RegisterRoutes(r chi.Router, ctl *Controller) {
	r.Route("/auth", func(r chi.Router) {
		r.With(middleware.Validate[LoginRequest]()).Post("/login", apperrors.Handle(ctl.Login))
		r.With(middleware.Validate[VerifyEmailRequest]()).Post("/verify-otp", apperrors.Handle(ctl.VerifyOTP))
		r.With(middleware.Validate[ForgotPasswordRequest]()).Post("/forgot-password", apperrors.Handle(ctl.ForgotPassword))
		r.With(middleware.Validate[ResetPasswordRequest]()).Post("/reset-password", apperrors.Handle(ctl.ResetPassword))

		r.With(
			middleware.Auth(),
			middleware.Validate[ChangePasswordRequest](),
		).Post("/change-password", apperrors.Handle(ctl.ChangePassword))
	})
}
