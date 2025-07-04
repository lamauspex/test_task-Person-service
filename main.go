package main

import (
	"MyProjects/person_service/internal/db"
	"MyProjects/person_service/internal/http"
	"MyProjects/person_service/internal/logic"
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Инициализация базы данных
	dbConnection, err := sql.Open("sqlite3", "./person.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConnection.Close()

	// Инициализация репозитория и сервиса
	repo := db.NewSqlitePersonRepo(dbConnection)
	service := logic.NewDefaultPersonService(repo)

	// Инициализация хендлеров
	handler := http.NewPersonHandler(service)

	e := echo.New()

	// Настройка маршрутов
	e.GET("/persons", handler.GetAll)        // Получить список всех пользователей
	e.GET("/persons/:id", handler.GetOne)    // Получить одного пользователя по ID
	e.POST("/persons", handler.Create)       // Создание нового пользователя
	e.PUT("/persons/:id", handler.Update)    // Обновление существующего пользователя
	e.DELETE("/persons/:id", handler.Delete) // Удаление пользователя по ID

	log.Fatal(e.Start(":8080")) // Запуск сервера
}
