package utils

import (
	"go-pos/internal/domain"
	"go-pos/pkg/response"
	"math"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ExtractQueryOptions(c echo.Context) domain.QueryOptions {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")
	sort := c.QueryParam("sort")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	var filters []domain.Filter
	filterParams := c.QueryParams()["filter"]

	for _, valString := range filterParams {
		parts := strings.SplitN(valString, ":", 3)
		if len(parts) == 3 {
			filters = append(filters, domain.Filter{
				Field:    parts[0],
				Operator: parts[1],
				Value:    parts[2],
			})
		}
	}

	return domain.QueryOptions{
		Page:    page,
		Limit:   limit,
		Search:  search,
		Sort:    sort,
		Filters: filters,
	}
}

func BuildMetaPage(page, limit int, totalItems int64) response.MetaPage {
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	return response.MetaPage{
		Page:       page,
		PerPage:    limit,
		TotalPages: totalPages,
		TotalItems: totalItems,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

func PaginationScope(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}
