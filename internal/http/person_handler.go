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

// Handler для управления моделями Person
type PersonHandler struct {
	service logic.PersonService
}

// Конструктор хэндлера
func NewPersonHandler(service logic.PersonService) *PersonHandler {
	return &PersonHandler{service: service}
}

// Получение всех пользователей с пагинацией и поиском
// @Summary Получение списка всех пользователей
// @Description Возвращает список всех пользователей с поддержкой пагинации и поиска
// @Tags Users
// @Accept json
// @Produce json
// @Param limit query integer false "Максимальное количество записей на одной странице"
// @Param offset query integer false "Позиция начальной записи"
// @Param search query string false "Строка поиска по имени или фамилии"
// @Success 200 {array} app.Person
// @Failure 500 {object} app.ErrorResponse
// @Router /persons [GET]
func (h *PersonHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context() // Получаем контекст из запроса

	// Получаем параметры запроса
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")
	searchTerm := c.QueryParam("search")

	var limit, offset int
	//var err error

	// Преобразуем строку в число
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid 'limit' parameter"})
		}
		limit = l
	}

	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid 'offset' parameter"})
		}
		offset = o
	}

	// Получаем список пользователей с учётом фильтров
	persons, err := h.service.ListPersons(ctx, limit, offset, searchTerm)
	if err != nil {
		log.Printf("Error listing persons: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.JSON(http.StatusOK, persons)
}

// Получение отдельного пользователя по ID
// @Summary Получение пользователя по идентификатору
// @Description Возвращает информацию о конкретном пользователе по его ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path integer true "Идентификатор пользователя"
// @Success 200 {object} app.Person
// @Failure 404 {object} app.ErrorResponse
// @Failure 500 {object} app.ErrorResponse
// @Router /persons/{id} [GET]
func (h *PersonHandler) GetOne(c echo.Context) error {
	ctx := c.Request().Context() // Получаем контекст из запроса

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid 'id' parameter"})
	}

	// Получаем пользователя по ID
	person, err := h.service.GetPerson(ctx, id)
	if err != nil {
		log.Printf("Error getting person by ID %d: %v\n", id, err)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Person not found"})
	}

	return c.JSON(http.StatusOK, person)
}

// Создание нового пользователя
// @Summary Создание нового пользователя
// @Description Создает новую запись пользователя в базе данных
// @Tags Users
// @Accept json
// @Produce json
// @Param person body app.Person true "Информация о новом пользователе"
// @Success 201 {object} app.Person
// @Failure 400 {object} app.ErrorResponse
// @Failure 500 {object} app.ErrorResponse
// @Router /persons [POST]
func (h *PersonHandler) Create(c echo.Context) error {
	ctx := c.Request().Context() // Получаем контекст из запроса

	newPerson := &app.Person{}

	// Декодируем тело запроса
	if err := json.NewDecoder(c.Request().Body).Decode(newPerson); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	// Создаем нового пользователя
	err := h.service.AddPerson(ctx, newPerson)
	if err != nil {
		log.Printf("Error creating person: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create person"})
	}

	return c.JSON(http.StatusCreated, newPerson)
}

// Обновление пользователя
// @Summary Обновление информации о пользователе
// @Description Обновляет информацию о существующем пользователе
// @Tags Users
// @Accept json
// @Produce json
// @Param id path integer true "Идентификатор пользователя"
// @Param person body app.Person true "Обновленная информация о пользователе"
// @Success 200 {object} app.Person
// @Failure 400 {object} app.ErrorResponse
// @Failure 500 {object} app.ErrorResponse
// @Router /persons/{id} [PUT]
func (h *PersonHandler) Update(c echo.Context) error {
	ctx := c.Request().Context() // Получаем контекст из запроса

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

	// Обновляем пользователя
	err = h.service.EditPerson(ctx, id, updatedPerson)
	if err != nil {
		log.Printf("Error updating person with ID %d: %v\n", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update person"})
	}

	return c.JSON(http.StatusOK, updatedPerson)
}

// Удаление пользователя
// @Summary Удаление пользователя
// @Description Полностью удаляет пользователя по его идентификатору
// @Tags Users
// @Accept json
// @Produce json
// @Param id path integer true "Идентификатор пользователя"
// @Success 204
// @Failure 404 {object} app.ErrorResponse
// @Failure 500 {object} app.ErrorResponse
// @Router /persons/{id} [DELETE]
func (h *PersonHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context() // Получаем контекст из запроса

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid 'id' parameter"})
	}

	// Удаляем пользователя
	err = h.service.RemovePerson(ctx, id)
	if err != nil {
		log.Printf("Error deleting person with ID %d: %v\n", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete person"})
	}

	return c.NoContent(http.StatusNoContent)
}
