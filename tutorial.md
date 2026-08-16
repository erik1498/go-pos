1. Install Echo V4 untuk jadi frameworknya, gorm untuk database ormnya, dan google uuid untuk uuid generatornya
```bash
    go get github.com/labstack/echo/v4 gorm.io/gorm gorm.io/driver/postgres github.com/google/uuid
```

2. Buatkan project dengan struktur berikut
```text
├── cmd
│   └── main.go
├── go.mod
├── go.sum
└── internal
    └── database
        └── db.go
```

3. Dalam cmd/main.go
```go
package main

import (
	"go-pos/internal/database"

	"github.com/labstack/echo/v4"
)

const (
	dbAddress = "host=localhost user=postgres password=postgres dbname=go-pos sslmode=disable"
)

func main() {
	e := echo.New()

	database.GetDB(dbAddress)

	e.Logger.Fatal(e.Start(":3000"))
}

```

4. Dalam internal/database/db.go
```go
package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB(dbAddress string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbAddress))

	if err != nil {
		panic(err)
	}

	return db
}
```

5. Kita akan memulai membuat CRUD untuk table categories, pertama kita akan membuat struktur folder seperti berikut
```text
├── cmd
│   └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── database
│   │   ├── db.go
│   │   └── seed.go
│   ├── delivery
│   │   └── rest
│   │       ├── category_handler.go
│   │       ├── handler.go
│   │       └── routes.go
│   ├── domain
│   │   ├── category.go
│   │   ├── global.go
│   │   └── query.go
│   ├── model
│   │   └── category.go
│   ├── repository
│   │   └── category_repository.go
│   └── usecase
│       └── category_usecase.go
└── pkg
    ├── response
    │   └── response.go
    └── utils
        ├── pagination.go
        └── query_validator.go
```

kita mulai dengan internal/model/category.go

```go
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;" json:"-"`
	PublicID  uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null;" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CategoryRequest struct {
	Name string `json:"name"`
}
```

pada pkg/response/response.go
```go
package repository

import (
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type categoryRepo struct {
	db *gorm.DB
}

func GetCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepo{
		db: db,
	}
}

func (cRepo *categoryRepo) GetAll(opts domain.QueryOptions) ([]model.Category, int64, error) {
	var categoryList = []model.Category{}
	var totalItems int64

	dbQuery := cRepo.db.Model(&model.Category{})

	if opts.Search != "" {
		dbQuery = dbQuery.Where("name ILIKE ?", "%"+opts.Search+"%")
	}

	for _, f := range opts.Filters {
		queryStr := f.Field + " " + f.Operator + " ?"
		dbQuery = dbQuery.Where(queryStr, f.Value)
	}

	if err := dbQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	err := dbQuery.Order(opts.Sort).Scopes(utils.PaginationScope(opts.Page, opts.Limit)).Find(&categoryList).Error

	return categoryList, totalItems, err
}

func (cRepo *categoryRepo) Create(category model.Category) (model.Category, error) {
	if err := cRepo.db.Create(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Category{}, domain.CategoryErrNotFound
		}
		return model.Category{}, err
	}

	return category, nil
}

func (cRepo *categoryRepo) GetByPublicID(id uuid.UUID) (model.Category, error) {
	var category model.Category

	if err := cRepo.db.Where(&model.Category{PublicID: id}).First(&category).Error; err != nil {
		return model.Category{}, err
	}

	return category, nil
}

func (cRepo *categoryRepo) UpdateCategoryByID(id uuid.UUID, category model.Category) (model.Category, error) {
	if err := cRepo.db.Where(&model.Category{PublicID: id}).Clauses(clause.Returning{}).Updates(&category).Error; err != nil {
		return model.Category{}, err
	}
	return category, nil
}

func (cRepo *categoryRepo) DeleteCategoryByID(id uuid.UUID) error {
	if err := cRepo.db.Where(&model.Category{PublicID: id}).Delete(&model.Category{}).Error; err != nil {
		return err
	}

	return nil
}

```
pada pkg/utils/pagination.go
```go
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
```
pada pkg/utils/query_validator.go
```go
package utils

import (
	"go-pos/internal/domain"
	"strings"
)

var operatorMap = map[string]string{
	"eq":   "=",
	"neq":  "!=",
	"gt":   ">",
	"gte":  ">=",
	"lt":   "<",
	"lte":  "<=",
	"like": "ILIKE",
}

func SanitizeQuery(opts domain.QueryOptions, allowedFields, allowedSorts map[string]bool, defaultSort string) domain.QueryOptions {
	var validSorts []string

	if opts.Sort != "" {
		sorts := strings.Split(opts.Sort, ",")
		for _, s := range sorts {
			s = strings.TrimSpace(s)
			parts := strings.Split(s, " ")
			if len(parts) > 0 {
				field := parts[0]
				dir := "asc"
				if len(parts) == 2 && strings.ToLower(parts[1]) == "desc" {
					dir = "desc"
				}

				if allowedSorts[field] {
					validSorts = append(validSorts, field+" "+dir)
				}
			}
		}
	}

	opts.Sort = strings.Join(validSorts, ", ")
	if opts.Sort == "" {
		opts.Sort = defaultSort
	}

	var validFilter []domain.Filter
	for _, f := range opts.Filters {
		if !allowedFields[f.Field] {
			continue
		}

		sqlOperator, ok := operatorMap[f.Operator]
		if !ok {
			continue
		}

		if sqlOperator == "ILIKE" {
			f.Value = "%" + f.Value.(string) + "%"
		}

		validFilter = append(validFilter, domain.Filter{
			Field:    f.Field,
			Operator: sqlOperator,
			Value:    f.Value,
		})
	}

	opts.Filters = validFilter

	return opts
}

```

lalu pada internal/repository/category_repository.go

```go
package repository

import (
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type categoryRepo struct {
	db *gorm.DB
}

func GetCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepo{
		db: db,
	}
}

func (cRepo *categoryRepo) GetAll(opts domain.QueryOptions) ([]model.Category, int64, error) {
	var categoryList = []model.Category{}
	var totalItems int64

	dbQuery := cRepo.db.Model(&model.Category{})

	if opts.Search != "" {
		dbQuery = dbQuery.Where("name ILIKE ?", "%"+opts.Search+"%")
	}

	for _, f := range opts.Filters {
		queryStr := f.Field + " " + f.Operator + " ?"
		dbQuery = dbQuery.Where(queryStr, f.Value)
	}

	if err := dbQuery.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	err := dbQuery.Order(opts.Sort).Scopes(utils.PaginationScope(opts.Page, opts.Limit)).Find(&categoryList).Error

	return categoryList, totalItems, err
}

func (cRepo *categoryRepo) Create(category model.Category) (model.Category, error) {
	if err := cRepo.db.Create(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Category{}, domain.CategoryErrNotFound
		}
		return model.Category{}, err
	}

	return category, nil
}

func (cRepo *categoryRepo) GetByPublicID(id uuid.UUID) (model.Category, error) {
	var category model.Category

	if err := cRepo.db.Where(&model.Category{PublicID: id}).First(&category).Error; err != nil {
		return model.Category{}, err
	}

	return category, nil
}

func (cRepo *categoryRepo) UpdateCategoryByID(id uuid.UUID, category model.Category) (model.Category, error) {
	if err := cRepo.db.Where(&model.Category{PublicID: id}).Clauses(clause.Returning{}).Updates(&category).Error; err != nil {
		return model.Category{}, err
	}
	return category, nil
}

func (cRepo *categoryRepo) DeleteCategoryByID(id uuid.UUID) error {
	if err := cRepo.db.Where(&model.Category{PublicID: id}).Delete(&model.Category{}).Error; err != nil {
		return err
	}

	return nil
}
```

pada internal/usecase/category_usecase.go
```go
package usecase

import (
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/utils"

	"github.com/google/uuid"
)

type categoryUsecase struct {
	cRepo domain.CategoryRepository
}

func GetCategoryUsecase(cRepo domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{
		cRepo: cRepo,
	}
}

func (pUsecase *categoryUsecase) GetAll(opts domain.QueryOptions) ([]model.Category, int64, error) {
	allowedFields := map[string]bool{
		"name":       true,
		"created_at": true,
	}

	allowedSorts := map[string]bool{
		"name":       true,
		"created_at": true,
	}

	cleanOpts := utils.SanitizeQuery(opts, allowedFields, allowedSorts, "created_at desc")

	return pUsecase.cRepo.GetAll(cleanOpts)
}

func (pUsecase *categoryUsecase) Create(req model.Category) (model.Category, error) {
	category := model.Category{
		PublicID: uuid.Must(uuid.NewV7()),
		Name:     req.Name,
	}

	return pUsecase.cRepo.Create(category)
}

func (pUsecase *categoryUsecase) GetByPublicID(id uuid.UUID) (model.Category, error) {
	return pUsecase.cRepo.GetByPublicID(id)
}

func (pUsecase *categoryUsecase) UpdateCategoryByID(id uuid.UUID, req model.Category) (model.Category, error) {
	return pUsecase.cRepo.UpdateCategoryByID(id, req)
}

func (pUsecase *categoryUsecase) DeleteCategoryByID(id uuid.UUID) error {
	return pUsecase.cRepo.DeleteCategoryByID(id)
}

```
pada internal/delivery/rest/handler.go
```go
package rest

import (
	"go-pos/internal/domain"
)

type handler struct {
	cUsecase domain.CategoryUsecase
}

func NewHandler(
	cUsecase domain.CategoryUsecase,
) *handler {
	return &handler{
		cUsecase: cUsecase,
	}
}
```
pada internal/delivery/category_handler.go
```go
package rest

import (
	"encoding/json"
	"errors"
	"go-pos/internal/domain"
	"go-pos/internal/model"
	"go-pos/pkg/response"
	"go-pos/pkg/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *handler) GetAll(c echo.Context) error {
	opts := utils.ExtractQueryOptions(c)

	categoryList, totalItems, err := h.cUsecase.GetAll(opts)
	if err != nil {
		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	meta := utils.BuildMetaPage(opts.Page, opts.Limit, totalItems)

	return response.SuccessWithMeta(c, http.StatusOK, domain.SuccessGetData, categoryList, meta)
}

func (h *handler) Create(c echo.Context) error {
	var req model.Category
	err := json.NewDecoder(c.Request().Body).Decode(&req)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	category, err := h.cUsecase.Create(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": category,
	})
}

func (h *handler) GetByPublicID(c echo.Context) error {
	idParam := c.Param("id")

	publicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	category, err := h.cUsecase.GetByPublicID(publicID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, err.Error())
		}

		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessGetDataByID, category)
}

func (h *handler) UpdateCategoryById(c echo.Context) error {
	idParam := c.Param("id")

	publicId, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	var req model.Category

	err = json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrBadRequest.Error())
	}

	category, err := h.cUsecase.UpdateCategoryByID(publicId, req)
	if err != nil {

		if errors.Is(err, domain.CategoryErrNotFound) {
			return response.ErrNotFound(c, domain.CategoryErrNotFound.Error())
		}

		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.Success(c, http.StatusOK, domain.SuccessUpdateData, category)
}

func (h *handler) DeleteCategoryByID(c echo.Context) error {
	idParam := c.Param("id")

	publicID, err := uuid.Parse(idParam)
	if err != nil {
		return response.ErrBadRequest(c, domain.ErrIDInvalid.Error())
	}

	err = h.cUsecase.DeleteCategoryByID(publicID)
	if err != nil {

		if errors.Is(err, domain.CategoryErrNotFound) {
			return response.ErrNotFound(c, domain.CategoryErrNotFound.Error())
		}

		return response.ErrInternalServer(c, domain.ErrInternalServer.Error())
	}

	return response.NoContent(c)
}

```
pada internal/delivery/routes.go
```go
package rest

import "github.com/labstack/echo/v4"

func LoadRoutes(e *echo.Echo, h *handler) {
	categoryGroup := e.Group("/categories")

	categoryGroup.GET("", h.GetAll)
	categoryGroup.POST("", h.Create)
	categoryGroup.GET("/:id", h.GetByPublicID)
	categoryGroup.PUT("/:id", h.UpdateCategoryById)
	categoryGroup.DELETE("/:id", h.DeleteCategoryByID)
}
```
pada internal/database/db.go
```go
package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB(dbAddress string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbAddress))

	if err != nil {
		panic(err)
	}

	seedDB(db)

	return db
}
```
pada internal/database/seed.go
```go
package database

import (
	"go-pos/internal/model"

	"gorm.io/gorm"
)

func seedDB(db *gorm.DB) {
	db.AutoMigrate(&model.Category{})
}
```
pada cmd/main.go

```go
package main

import (
	"go-pos/internal/database"
	"go-pos/internal/delivery/rest"
	"go-pos/internal/repository"
	"go-pos/internal/usecase"

	"github.com/labstack/echo/v4"
)

const (
	dbAddress = "host=localhost user=postgres password=postgres dbname=go-pos sslmode=disable"
)

func main() {
	e := echo.New()

	db := database.GetDB(dbAddress)

	cRepo := repository.GetCategoryRepository(db)

	pUsecase := usecase.GetCategoryUsecase(cRepo)

	handler := rest.NewHandler(pUsecase)

	rest.LoadRoutes(e, handler)

	e.Logger.Fatal(e.Start(":3000"))
}

```