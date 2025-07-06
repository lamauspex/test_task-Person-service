package tests

import (
	"MyProjects/person_service/internal/app"
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Устанавливаем соединение с бд в и создаём таблицу Person
func PrepareTestDatabase(t *testing.T) *sql.DB {
	t.Helper()

	// Открываем подключение к базе данных в памяти
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Unable to open database connection: %v", err)
	}

	// Создаём таблицу для тестов
	createTableQuery := `
CREATE TABLE IF NOT EXISTS persons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    phone TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL
);
`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Unable to create table: %v", err)
	}

	return db
}

// Создаем фейковую запись в базу данных
func InsertFakePerson(db *sql.DB, t *testing.T) *app.Person {
	t.Helper()

	// Устанавливаем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Генерация уникального Email
	randomSuffix := rand.Intn(100_000)
	uniqueEmail := fmt.Sprintf("test-%d-%d@example.com", time.Now().UnixNano(), randomSuffix)

	fakePerson := &app.Person{
		Email:     uniqueEmail,
		Phone:     "+79991234567",
		FirstName: "John",
		LastName:  "Wick",
	}

	insertQuery := `
    INSERT INTO persons (email, phone, first_name, last_name)
    VALUES ($1, $2, $3, $4)
    RETURNING id
`

	var insertedID int
	err := db.QueryRow(insertQuery, fakePerson.Email, fakePerson.Phone, fakePerson.FirstName, fakePerson.LastName).Scan(&insertedID)
	if err != nil {
		t.Fatalf("Unable to insert fake person: %v", err)
	}

	fakePerson.ID = insertedID
	return fakePerson
}

// Очищаем бд после завершения тестов
func CleanUp(db *sql.DB, t *testing.T) {
	t.Helper()

	dropTableQuery := `
DROP TABLE IF EXISTS persons;
`

	_, err := db.Exec(dropTableQuery)
	if err != nil {
		t.Logf("Warning: unable to clean up database after testing: %v", err)
	}

	// Закрытие соединения с бд
	err = db.Close()
	if err != nil {
		t.Logf("Warning: unable to close database connection: %v", err)
	}
}

// Останавливает тест и выводит сообщение об ошибке
func MustFailIfError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
