package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/mdsharifulislam-r/go-backend-template/config"
	"github.com/mdsharifulislam-r/go-backend-template/internal/email"
	apperrors "github.com/mdsharifulislam-r/go-backend-template/internal/errors"
	"github.com/mdsharifulislam-r/go-backend-template/internal/middleware"
	"github.com/mdsharifulislam-r/go-backend-template/internal/module/auth"
	"github.com/mdsharifulislam-r/go-backend-template/internal/module/user"
	"github.com/mdsharifulislam-r/go-backend-template/internal/response"
	"gorm.io/gorm"
)

type Server struct {
	cfg    *config.Config
	db     *gorm.DB
	router chi.Router
}

func New(cfg *config.Config, db *gorm.DB) *Server {
	s := &Server{
		cfg:    cfg,
		db:     db,
		router: chi.NewRouter(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.CORS)
	s.router.Use(middleware.Logger)
	s.router.Use(chimw.RealIP)
	s.router.Use(chimw.RequestID)
	s.router.Use(chimw.Timeout(60 * time.Second))

	uploadDir, _ := filepath.Abs("uploads")
	_ = os.MkdirAll(uploadDir, 0o755)
	s.router.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	emailSvc := email.NewService(s.cfg)
	userSvc := user.NewService(s.db, emailSvc)
	authSvc := auth.NewService(s.db, emailSvc)

	if err := userSvc.EnsureSuperAdmin(s.cfg.SuperAdminEmail, s.cfg.SuperAdminPassword); err != nil {
		log.Printf("warning: failed to seed super admin: %v", err)
	}

	userCtl := user.NewController(userSvc)
	authCtl := auth.NewController(authSvc)

	s.router.Get("/", apperrors.Handle(healthHandler))
	s.router.Get("/health", apperrors.Handle(healthHandler))

	s.router.Route("/api/v1", func(r chi.Router) {
		user.RegisterRoutes(r, userCtl)
		auth.RegisterRoutes(r, authCtl)
	})
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.IP, s.cfg.Port)
	log.Printf("Server running on http://%s", addr)
	log.Printf("API base path: /api/v1")
	return http.ListenAndServe(addr, s.router)
}

func healthHandler(w http.ResponseWriter, r *http.Request) error {
	response.OK(w, "Server is running", map[string]string{
		"status": "ok",
	})
	return nil
}
