package logic

import (
	"MyProjects/person_service/internal/app"
	"MyProjects/person_service/internal/db"
	"context"
	"errors"
	"fmt"
)

// Интерфейс PersonService объявляет методы для управления пользователями
type PersonService interface {
	ListPersons(ctx context.Context, limit, offset int, searchTerm string) ([]*app.Person, error)
	GetPerson(ctx context.Context, id int) (*app.Person, error)
	AddPerson(ctx context.Context, p *app.Person) error
	EditPerson(ctx context.Context, id int, p *app.Person) error
	RemovePerson(ctx context.Context, id int) error
}

// DefaultPersonService реализует интерфейс PersonService
type DefaultPersonService struct {
	repo db.PersonRepository
}

// Новый экземпляр сервиса
func NewDefaultPersonService(repo db.PersonRepository) PersonService {
	return &DefaultPersonService{repo: repo}
}

// Список пользователей с фильтрацией по поиску, лимиту и смещению
func (s *DefaultPersonService) ListPersons(ctx context.Context, limit, offset int, searchTerm string) ([]*app.Person, error) {
	// Проверка допустимости лимитов и офсетов
	if limit < 0 || offset < 0 {
		return nil, errors.New("invalid pagination parameters")
	}

	// Устанавливаем максимальный предел записей
	if limit > 100 {
		limit = 100
	}

	// Обращаемся к репозиторию с новым критерием поиска
	persons, err := s.repo.SearchAndList(ctx, limit, offset, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("database error when searching or listing persons: %w", err)
	}

	return persons, nil
}

// Получить пользователя по ID
func (s *DefaultPersonService) GetPerson(ctx context.Context, id int) (*app.Person, error) {
	person, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) { // Ловим специфичную ошибку
			return nil, fmt.Errorf("person with ID %d does not exist", id)
		}
		return nil, fmt.Errorf("database error when finding person: %w", err) // Генеральная ошибка
	}
	return person, nil
}

// Добавить пользователя
func (s *DefaultPersonService) AddPerson(ctx context.Context, p *app.Person) error {
	err := s.repo.Create(ctx, p)
	if err != nil {
		return fmt.Errorf("database error when adding person: %w", err) // Возвращаем ошибку создания
	}
	return nil
}

// Редактировать пользователя
func (s *DefaultPersonService) EditPerson(ctx context.Context, id int, p *app.Person) error {
	err := s.repo.Update(ctx, id, p)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) { // Специфичная ошибка, если запись отсутствует
			return fmt.Errorf("person with ID %d does not exist", id)
		}
		return fmt.Errorf("database error when updating person: %w", err) // Общая ошибка
	}
	return nil
}

// Удалить пользователя
func (s *DefaultPersonService) RemovePerson(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) { // Отлавливаем случай удаления несуществующей записи
			return fmt.Errorf("person with ID %d does not exist", id)
		}
		return fmt.Errorf("database error when removing person: %w", err) // Обычная ошибка базы данных
	}
	return nil
}
