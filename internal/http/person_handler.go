package http

import (
	"MyProjects/person_service/internal/app"
	"MyProjects/person_service/internal/logic"
	"encoding/json"
	"net/http"
	"strconv"

	"log"

	"github.com/labstack/echo/v4"
)

type PersonHandler struct {
	service logic.PersonService
}

func NewPersonHandler(service logic.PersonService) *PersonHandler {
	return &PersonHandler{service: service}
}

// GetAll обрабатывает запросы типа GET /person/, возвращает список моделей Person
func (h *PersonHandler) GetAll(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	// Проверка допустимости параметров
	limit, err := strconv.Atoi(limitStr)
	if err != nil && limitStr != "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid 'limit' parameter"})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil && offsetStr != "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid 'offset' parameter"})
	}

	// Получаем список персон
	persons, err := h.service.ListPersons(limit, offset)
	if err != nil {
		log.Printf("Error listing persons: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.JSON(http.StatusOK, persons)
}

// GetOne обрабатывает запросы типа GET /person/{id}, возвращает одну модель Person
func (h *PersonHandler) GetOne(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid 'id' parameter"})
	}

	// Получаем конкретного человека
	person, err := h.service.GetPerson(id)
	if err != nil {
		log.Printf("Error getting person by ID %d: %v\n", id, err)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Person not found"})
	}

	return c.JSON(http.StatusOK, person)
}

// Create обрабатывает запросы типа POST /person/, создает новую модель Person
func (h *PersonHandler) Create(c echo.Context) error {
	newPerson := &app.Person{}

	// Декодируем тело запроса
	if err := json.NewDecoder(c.Request().Body).Decode(newPerson); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	// Создаем нового человека
	err := h.service.AddPerson(newPerson)
	if err != nil {
		log.Printf("Error creating person: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create person"})
	}

	return c.JSON(http.StatusCreated, newPerson)
}

// Update обрабатывает запросы типа PUT /person/{id}, обновляет модель Person
func (h *PersonHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid 'id' parameter"})
	}

	updatedPerson := &app.Person{}

	// Декодируем тело запроса
	if err := json.NewDecoder(c.Request().Body).Decode(updatedPerson); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	// Обновляем человека
	err = h.service.EditPerson(id, updatedPerson)
	if err != nil {
		log.Printf("Error updating person with ID %d: %v\n", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update person"})
	}

	return c.JSON(http.StatusOK, updatedPerson)
}

// Delete обрабатывает запросы типа DELETE /person/{id}, удаляет модель Person
func (h *PersonHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid 'id' parameter"})
	}

	// Удаляем человека
	err = h.service.RemovePerson(id)
	if err != nil {
		log.Printf("Error deleting person with ID %d: %v\n", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete person"})
	}

	return c.NoContent(http.StatusNoContent)
}
