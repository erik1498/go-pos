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
package model

import "github.com/google/uuid"

type Category struct {
	ID       int64     `gorm:"primaryKey;autoIncrement;" json:"-"`
	PublicID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null;" json:"id"`
	Name     string    `gorm:"type:varchar(100);not null" json:"name"`
}
```