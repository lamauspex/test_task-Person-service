package main

import (
	"MyProjects/person_service/internal/db"
	"MyProjects/person_service/internal/http"
	"MyProjects/person_service/internal/logic"
	"MyProjects/person_service/internal/middleware"
	"database/sql"
	"fmt"
	"log"

	_ "MyProjects/person_service/docs"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"

	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

func main() {
	// Загрузка переменных из файла `.env`
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	// Чтение переменных
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	// Формирование строки подключения к базе данных
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)

	// Открываем соединение с базой данных
	dbConnection, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConnection.Close()

	// Инициализируем репозиторий и сервис
	repo := db.NewPostgresPersonRepo(dbConnection)
	service := logic.NewDefaultPersonService(repo)

	// Инициализируем обработчик HTTP
	handler := http.NewPersonHandler(service)

	// Инициализируем Echo-framework
	e := echo.New()

	// Создание конфигурации логгера
	config := zap.NewProductionConfig()
	config.Level.SetLevel(zap.DebugLevel)
	logger, _ := config.Build()
	defer logger.Sync()

	// Добавляем логгер в виде мидлвари
	e.Use(middleware.LoggerMiddleware(logger))

	// Регистрация маршрутов API
	e.GET("/persons", handler.GetAll)
	e.GET("/persons/:id", handler.GetOne)
	e.POST("/persons", handler.Create)
	e.PUT("/persons/:id", handler.Update)
	e.DELETE("/persons/:id", handler.Delete)

	// Роут для swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Запускаем веб-сервер
	log.Fatal(e.Start(":8080"))
}
