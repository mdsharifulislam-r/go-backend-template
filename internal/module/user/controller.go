package user

import (
	"net/http"

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

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) error {
	payload, ok := middleware.BodyFromContext[CreateUserRequest](r.Context())
	if !ok {
		return apperrors.BadRequest("Invalid request body")
	}

	user, err := c.service.Create(payload)
	if err != nil {
		return err
	}

	response.Created(w, "User created successfully. Please verify your email.", user)
	return nil
}

func (c *Controller) GetProfile(w http.ResponseWriter, r *http.Request) error {
	claims, ok := middleware.CurrentUser(r.Context())
	if !ok {
		return apperrors.Unauthorized("You are not authorized")
	}

	user, err := c.service.GetProfile(claims.ID)
	if err != nil {
		return err
	}

	response.OK(w, "Profile fetched successfully", user)
	return nil
}

func (c *Controller) UpdateProfile(w http.ResponseWriter, r *http.Request) error {
	claims, ok := middleware.CurrentUser(r.Context())
	if !ok {
		return apperrors.Unauthorized("You are not authorized")
	}

	payload := UpdateProfileRequest{
		Name:    r.FormValue("name"),
		Contact: r.FormValue("contact"),
		Image:   r.FormValue("image"),
	}

	files := middleware.GetUploadedFiles(r.Context())
	imagePath := ""
	if imgs := files["image"]; len(imgs) > 0 {
		imagePath = imgs[0]
	}

	user, err := c.service.UpdateProfile(claims.ID, payload, imagePath)
	if err != nil {
		return err
	}

	response.OK(w, "Profile updated successfully", user)
	return nil
}

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) error {
	query := map[string]string{}
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			query[key] = values[0]
		}
	}

	users, pagination, err := c.service.GetAll(query)
	if err != nil {
		return err
	}

	response.Send(w, response.Options{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "Users fetched successfully",
		Data:       users,
		Pagination: pagination,
	})
	return nil
}
