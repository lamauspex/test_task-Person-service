package db

import (
	"MyProjects/person_service/internal/app"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Определение специальной ошибки для случаев, когда запись не найдена
var ErrRecordNotFound = errors.New("record not found")

// Интерфейс PersonRepository описывает контракт для работы с моделью Person
type PersonRepository interface {
	GetAll(limit, offset int) ([]app.Person, error)
	FindByID(id int) (*app.Person, error)
	Create(person *app.Person) error
	Update(id int, updatedPerson *app.Person) error
	Delete(id int) error
}

// Struct для реализации интерфейса PersonRepository
type SqlitePersonRepo struct {
	DB *sql.DB
}

// Фабричный метод для создания экземпляра репозитория
func NewSqlitePersonRepo(db *sql.DB) PersonRepository {
	return &SqlitePersonRepo{DB: db}
}

// GetAll извлекает список персон с применением фильтра (лимит и смещение)
func (r *SqlitePersonRepo) GetAll(limit, offset int) ([]app.Person, error) {
	rows, err := r.DB.Query(
		"SELECT id, email, phone, first_name, last_name FROM persons LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("database error in GetAll: %w", err)
	}
	defer rows.Close()

	var persons []app.Person
	for rows.Next() {
		var person app.Person
		if err := rows.Scan(&person.ID, &person.Email, &person.Phone, &person.FirstName, &person.LastName); err != nil {
			return nil, fmt.Errorf("scanning error in GetAll: %w", err)
		}
		persons = append(persons, person)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error in GetAll: %w", err)
	}

	return persons, nil
}

// FindByID находит и возвращает конкретную запись по ID
func (r *SqlitePersonRepo) FindByID(id int) (*app.Person, error) {
	row := r.DB.QueryRow("SELECT id, email, phone, first_name, last_name FROM persons WHERE id=?", id)

	var p app.Person
	err := row.Scan(&p.ID, &p.Email, &p.Phone, &p.FirstName, &p.LastName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound // Возвращаем нашу спецошибку
		}
		return nil, fmt.Errorf("database error in FindByID: %w", err)
	}

	return &p, nil
}

// Create сохраняет новую запись в базе данных
func (r *SqlitePersonRepo) Create(person *app.Person) error {
	stmt, err := r.DB.Prepare("INSERT INTO persons(email, phone, first_name, last_name) VALUES(?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare statement error in Create: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(person.Email, person.Phone, person.FirstName, person.LastName)
	if err != nil {
		return fmt.Errorf("execution error in Create: %w", err)
	}

	return nil
}

// Update обновляет существующие данные в базе данных
func (r *SqlitePersonRepo) Update(id int, updatedPerson *app.Person) error {
	stmt, err := r.DB.Prepare("UPDATE persons SET email=?, phone=?, first_name=?, last_name=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("prepare statement error in Update: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(updatedPerson.Email, updatedPerson.Phone, updatedPerson.FirstName, updatedPerson.LastName, id)
	if err != nil {
		return fmt.Errorf("execution error in Update: %w", err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil || affectedRows == 0 {
		return fmt.Errorf("no records were updated in Update: %w", ErrRecordNotFound)
	}

	return nil
}

// Delete удаляет запись из базы данных по её ID
func (r *SqlitePersonRepo) Delete(id int) error {
	stmt, err := r.DB.Prepare("DELETE FROM persons WHERE id=?")
	if err != nil {
		return fmt.Errorf("prepare statement error in Delete: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("execution error in Delete: %w", err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil || affectedRows == 0 {
		return fmt.Errorf("no records were deleted in Delete: %w", ErrRecordNotFound)
	}

	return nil
}
