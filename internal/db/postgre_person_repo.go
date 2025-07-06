package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"MyProjects/person_service/internal/app"
)

// Специальная ошибка для случая, когда запись не найдена
var ErrRecordNotFound = errors.New("record not found")

// Интерфейс PersonRepository описывает набор методов для взаимодействия с базой данных
type PersonRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]*app.Person, error)
	FindByID(ctx context.Context, id int) (*app.Person, error)
	Create(ctx context.Context, person *app.Person) error
	Update(ctx context.Context, id int, updatedPerson *app.Person) error
	Delete(ctx context.Context, id int) error
	SearchAndList(ctx context.Context, limit, offset int, searchTerm string) ([]*app.Person, error) // новый метод
}

// Репозиторий PostgePersonRepo реализует интерфейс PersonRepository с поддержкой PostgreSQL
type PostgresPersonRepo struct {
	DB *sql.DB
}

// Конструктор нового репозитория
func NewPostgresPersonRepo(db *sql.DB) PersonRepository {
	return &PostgresPersonRepo{
		DB: db,
	}
}

// Получение списка пользователей
func (r *PostgresPersonRepo) GetAll(ctx context.Context, limit, offset int) ([]*app.Person, error) {
	query := `
    SELECT id, email, phone, first_name, last_name
    FROM persons
    ORDER BY id
    LIMIT $1
    OFFSET $2`

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("database error in GetAll: %w", err)
	}
	defer rows.Close()

	var persons []*app.Person
	for rows.Next() {
		p := &app.Person{}
		if err := rows.Scan(&p.ID, &p.Email, &p.Phone, &p.FirstName, &p.LastName); err != nil {
			return nil, fmt.Errorf("scanning error in GetAll: %w", err)
		}
		persons = append(persons, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error in GetAll: %w", err)
	}

	return persons, nil
}

// Поиск пользователя по ID
func (r *PostgresPersonRepo) FindByID(ctx context.Context, id int) (*app.Person, error) {
	query := `SELECT id, email, phone, first_name, last_name FROM persons WHERE id = $1`

	row := r.DB.QueryRowContext(ctx, query, id)

	p := &app.Person{}
	err := row.Scan(&p.ID, &p.Email, &p.Phone, &p.FirstName, &p.LastName)
	if err != nil {
		log.Printf("Error scanning row: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("database error in FindByID: %w", err)
	}

	return p, nil
}

// Создание новой записи в БД
func (r *PostgresPersonRepo) Create(ctx context.Context, person *app.Person) error {
	query := `
    INSERT INTO persons (email, phone, first_name, last_name)
    VALUES ($1, $2, $3, $4)
    RETURNING id`

	err := r.DB.QueryRowContext(ctx, query, person.Email, person.Phone, person.FirstName, person.LastName).Scan(&person.ID)
	if err != nil {
		return fmt.Errorf("insert execution error in Create: %w", err)
	}

	return nil
}

// Обновление существующего пользователя
func (r *PostgresPersonRepo) Update(ctx context.Context, id int, updatedPerson *app.Person) error {
	query := `
    UPDATE persons
    SET email = $1, phone = $2, first_name = $3, last_name = $4
    WHERE id = $5`

	result, err := r.DB.ExecContext(ctx, query, updatedPerson.Email, updatedPerson.Phone, updatedPerson.FirstName, updatedPerson.LastName, id)
	if err != nil {
		return fmt.Errorf("update execution error in Update: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no record was updated: %w", ErrRecordNotFound)
	}

	return nil
}

// Удаление пользователя по ID
func (r *PostgresPersonRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM persons WHERE id = $1"

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete execution error in Delete: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no record was deleted: %w", ErrRecordNotFound)
	}

	return nil
}

// Поиск и получение списка пользователей с фильтрацией и пагинацией
func (r *PostgresPersonRepo) SearchAndList(ctx context.Context, limit, offset int, searchTerm string) ([]*app.Person, error) {
	query := `
    SELECT id, email, phone, first_name, last_name
    FROM persons
    WHERE first_name ILIKE '%' || $1 || '%'
       OR last_name ILIKE '%' || $1 || '%'
       OR email ILIKE '%' || $1 || '%'
    ORDER BY id
    LIMIT $2
    OFFSET $3`

	rows, err := r.DB.QueryContext(ctx, query, searchTerm, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("database error in SearchAndList: %w", err)
	}
	defer rows.Close()

	var persons []*app.Person
	for rows.Next() {
		p := &app.Person{}
		if err := rows.Scan(&p.ID, &p.Email, &p.Phone, &p.FirstName, &p.LastName); err != nil {
			return nil, fmt.Errorf("scanning error in SearchAndList: %w", err)
		}
		persons = append(persons, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error in SearchAndList: %w", err)
	}

	return persons, nil
}
