package middleware

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mdsharifulislam-r/go-backend-template/internal/response"
)

type UploadField struct {
	Name     string
	MaxFiles int
	Required bool
}

type fileContextKey string

const UploadedFilesKey fileContextKey = "uploadedFiles"

func UploadFiles(destDir string, maxSize int64, fields ...UploadField) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseMultipartForm(maxSize); err != nil {
				response.Fail(w, http.StatusBadRequest, "Invalid multipart form", nil)
				return
			}

			if err := os.MkdirAll(destDir, 0o755); err != nil {
				response.Fail(w, http.StatusInternalServerError, "Unable to prepare upload directory", nil)
				return
			}

			uploaded := make(map[string][]string)

			for _, field := range fields {
				headers := r.MultipartForm.File[field.Name]
				if len(headers) == 0 {
					if field.Required {
						response.Fail(w, http.StatusBadRequest, field.Name+" is required", nil)
						return
					}
					continue
				}

				maxFiles := field.MaxFiles
				if maxFiles <= 0 {
					maxFiles = 1
				}
				if len(headers) > maxFiles {
					response.Fail(w, http.StatusBadRequest, "Too many files for "+field.Name, nil)
					return
				}

				paths := make([]string, 0, len(headers))
				for _, header := range headers {
					path, err := saveUpload(destDir, header)
					if err != nil {
						response.Fail(w, http.StatusBadRequest, err.Error(), nil)
						return
					}
					paths = append(paths, path)
				}
				uploaded[field.Name] = paths
			}

			ctx := context.WithValue(r.Context(), UploadedFilesKey, uploaded)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUploadedFiles(ctx context.Context) map[string][]string {
	files, _ := ctx.Value(UploadedFilesKey).(map[string][]string)
	return files
}

func GetSingleFilePath(ctx context.Context, key string) string {
	files := GetUploadedFiles(ctx)
	if imgs := files[key]; len(imgs) > 0 {
		return imgs[0]
	}
	return ""
}

func GetMultipleFilePaths(ctx context.Context, key string) []string {
	files := GetUploadedFiles(ctx)
	if imgs := files[key]; len(imgs) > 0 {
		return imgs
	}
	return []string{}
}

func saveUpload(destDir string, header *multipart.FileHeader) (string, error) {
	src, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("unable to open uploaded file")
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.NewString()[:8], ext)
	fullPath := filepath.Join(destDir, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("unable to save uploaded file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("unable to write uploaded file")
	}

	return "/uploads/" + filename, nil
}
