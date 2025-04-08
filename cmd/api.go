package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"simple-api/internal/cache"
	"simple-api/internal/dao"
	"simple-api/internal/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	databaseName := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", username, password, host, databaseName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	cacheAddr := fmt.Sprintf("%s:%d", os.Getenv("REDIS_HOST"), 6379)

	db.AutoMigrate(&model.Message{})

	dao := dao.New(db)
	cache := cache.New(cacheAddr)

	if err := cache.Ping(); err != nil {
		panic(fmt.Errorf("failed to connect to cache: %w", err))
	}

	// Start API server
	startApi(dao, cache)
}

func startApi(dao *dao.Dao, cache *cache.Cache) {
	api := echo.New()
	api.Use(middleware.Logger())

	api.GET("/", func(c echo.Context) error {
		ip := c.RealIP()
		remoteAddr := c.Request().RemoteAddr
		host := c.Request().Host

		return c.String(200, fmt.Sprintf("Hello from %s!\nYour IP: %s\nRemote Address: %s",
			host, ip, remoteAddr))
	})

	api.GET("/message/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(400, "invalid id")
		}

		uid := uint(id)

		messageCached := cache.Get(c.Request().Context(), fmt.Sprintf("message:%d", uid))
		if messageCached != "" {
			return c.String(200, fmt.Sprintf("message from cache: %s", messageCached))
		}

		messageDb, err := dao.Get(uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				msg := model.Message{
					ID:      uint(uid),
					Content: fmt.Sprintf("some message with id=%d", uid),
				}

				if err := dao.Save(&msg); err != nil {
					return c.String(500, fmt.Sprintf("error saving message to database: %v", err))
				}

				if err := cache.Set(c.Request().Context(), fmt.Sprintf("message:%d", uid), msg.Content); err != nil {
					return c.String(500, fmt.Sprintf("error saving message to cache: %v", err))
				}

				return c.String(200, fmt.Sprintf("new message created: %s", msg.Content))
			}

			return c.String(500, fmt.Sprintf("error getting message from database: %v", err))
		}

		if err := cache.Set(c.Request().Context(), fmt.Sprintf("message:%d", uid), messageDb.Content); err != nil {
			return c.String(500, fmt.Sprintf("error saving message to cache: %v", err))
		}

		return c.String(200, fmt.Sprintf("message from database: %s", messageDb.Content))
	})

	if err := api.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
