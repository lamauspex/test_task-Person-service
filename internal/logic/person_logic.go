package logic

import (
	"MyProjects/person_service/internal/app"
	"MyProjects/person_service/internal/db"
	"errors"
	"fmt"
)

type PersonService interface {
	ListPersons(limit, offset int) ([]*app.Person, error)
	GetPerson(id int) (*app.Person, error)
	AddPerson(*app.Person) error
	EditPerson(int, *app.Person) error
	RemovePerson(int) error
}

type DefaultPersonService struct {
	repo db.PersonRepository
}

func NewDefaultPersonService(repo db.PersonRepository) PersonService {
	return &DefaultPersonService{repo: repo}
}

// ListPersons возвращает список персон с указанным лимитом и смещением
func (s *DefaultPersonService) ListPersons(limit, offset int) ([]*app.Person, error) {
	if limit <= 0 || offset < 0 {
		return nil, errors.New("invalid pagination parameters") // Неправильные параметры постраничной разбивки
	}

	persons, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("database error when retrieving persons: %w", err) // Ошибка базы данных
	}

	// Конвертируем срез структур в срез указателей
	results := make([]*app.Person, len(persons))
	for i, p := range persons {
		results[i] = &p
	}
	return results, nil
}

// GetPerson ищет и возвращает одну модель Person по уникальному идентификатору
func (s *DefaultPersonService) GetPerson(id int) (*app.Person, error) {
	person, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) { // Специальная ошибка для случая отсутствия записи
			return nil, fmt.Errorf("person with ID %d does not exist", id)
		}
		return nil, fmt.Errorf("database error when finding person: %w", err) // Другая ошибка базы данных
	}
	return person, nil
}

// AddPerson добавляет новую модель Person в базу данных
func (s *DefaultPersonService) AddPerson(p *app.Person) error {
	err := s.repo.Create(p)
	if err != nil {
		return fmt.Errorf("database error when adding person: %w", err) // Ошибка базы данных
	}
	return nil
}

// EditPerson обновляет модель Person в базе данных
func (s *DefaultPersonService) EditPerson(id int, p *app.Person) error {
	err := s.repo.Update(id, p)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) { // Специальная ошибка для случая отсутствия записи
			return fmt.Errorf("person with ID %d does not exist", id)
		}
		return fmt.Errorf("database error when updating person: %w", err) // Другая ошибка базы данных
	}
	return nil
}

// RemovePerson удаляет модель Person из базы данных
func (s *DefaultPersonService) RemovePerson(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) { // Специальная ошибка для случая отсутствия записи
			return fmt.Errorf("person with ID %d does not exist", id)
		}
		return fmt.Errorf("database error when removing person: %w", err) // Другая ошибка базы данных
	}
	return nil
}
