package user

import (
	"github.com/go-chi/chi/v5"
	apperrors "github.com/mdsharifulislam-r/go-backend-template/internal/errors"
	"github.com/mdsharifulislam-r/go-backend-template/internal/middleware"
)

func RegisterRoutes(r chi.Router, ctl *Controller) {
	r.Route("/user", func(r chi.Router) {
		r.With(middleware.Validate[CreateUserRequest]()).Post("/", apperrors.Handle(ctl.Create))
		r.Get("/all", apperrors.Handle(ctl.GetAll))

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth())
			r.Get("/profile", apperrors.Handle(ctl.GetProfile))
			r.With(
				middleware.UploadFiles("uploads", 10<<20, middleware.UploadField{
					Name:     "image",
					MaxFiles: 1,
					Required: false,
				}),
			).Patch("/profile", apperrors.Handle(ctl.UpdateProfile))
		})
	})
}
