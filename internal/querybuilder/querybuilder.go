package querybuilder

import (
	"fmt"
	"strings"

	"github.com/mdsharifulislam-r/go-backend-template/internal/response"
	"gorm.io/gorm"
)

type Builder struct {
	db    *gorm.DB
	query map[string]string
}

func New(db *gorm.DB, query map[string]string) *Builder {
	return &Builder{db: db, query: query}
}

func (b *Builder) Search(fields ...string) *Builder {
	term := strings.TrimSpace(b.query["searchTerm"])
	if term == "" || len(fields) == 0 {
		return b
	}

	parts := make([]string, 0, len(fields))
	args := make([]any, 0, len(fields))
	for _, field := range fields {
		parts = append(parts, fmt.Sprintf("%s ILIKE ?", field))
		args = append(args, "%"+term+"%")
	}

	b.db = b.db.Where(strings.Join(parts, " OR "), args...)
	return b
}

func (b *Builder) Filter(exclude ...string) *Builder {
	skip := map[string]struct{}{
		"searchTerm": {},
		"sort":       {},
		"page":       {},
		"limit":      {},
		"fields":     {},
	}
	for _, key := range exclude {
		skip[key] = struct{}{}
	}

	for key, value := range b.query {
		if _, ok := skip[key]; ok || value == "" {
			continue
		}
		b.db = b.db.Where(fmt.Sprintf("%s = ?", key), value)
	}
	return b
}

func (b *Builder) Sort(defaultSort string) *Builder {
	sort := b.query["sort"]
	if sort == "" {
		sort = defaultSort
	}

	order := "ASC"
	field := sort
	if strings.HasPrefix(sort, "-") {
		order = "DESC"
		field = strings.TrimPrefix(sort, "-")
	}

	b.db = b.db.Order(fmt.Sprintf("%s %s", field, order))
	return b
}

func (b *Builder) Paginate(dest any) (*response.Pagination, error) {
	limit := parseInt(b.query["limit"], 10)
	page := parseInt(b.query["page"], 1)
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	var total int64
	if err := b.db.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * limit
	if err := b.db.Offset(offset).Limit(limit).Find(dest).Error; err != nil {
		return nil, err
	}

	totalPage := int((total + int64(limit) - 1) / int64(limit))
	return &response.Pagination{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}, nil
}

func parseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	n := 0
	for _, c := range value {
		if c < '0' || c > '9' {
			return fallback
		}
		n = n*10 + int(c-'0')
	}
	return n
}
