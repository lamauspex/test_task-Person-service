package app_test

import (
	"MyProjects/person_service/internal/app"
	"MyProjects/person_service/internal/db"
	"MyProjects/person_service/internal/tests"
	"reflect"
	"testing"
)

// Тест создания новой записи
func TestCreatePerson(t *testing.T) {
	dbConn := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(dbConn, t)

	repo := db.NewSqlitePersonRepo(dbConn)

	testPerson := &app.Person{
		Email:     "john.doe@example.com",
		Phone:     "+79991234567",
		FirstName: "John",
		LastName:  "Doe",
	}

	err := repo.Create(testPerson)
	tests.MustFailIfError(t, err)

	foundPerson, err := repo.FindByID(testPerson.ID)
	tests.MustFailIfError(t, err)

	if !reflect.DeepEqual(foundPerson, testPerson) {
		t.Errorf("Expected person %+v but got %+v", testPerson, foundPerson)
	}
}

// Тест обновления существующей записи
func TestUpdatePerson(t *testing.T) {
	dbConn := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(dbConn, t)

	repo := db.NewSqlitePersonRepo(dbConn)

	// Вставляем исходную запись
	fakePerson := tests.InsertFakePerson(dbConn, t)

	// Изменяем некоторые поля
	updatedPerson := &app.Person{
		ID:        fakePerson.ID,
		Email:     "updated@example.com",
		Phone:     "+79997654321",
		FirstName: "Jane",
		LastName:  "Smith",
	}

	err := repo.Update(fakePerson.ID, updatedPerson)
	tests.MustFailIfError(t, err)

	// Проверяем обновление
	foundPerson, err := repo.FindByID(fakePerson.ID)
	tests.MustFailIfError(t, err)

	if !reflect.DeepEqual(foundPerson, updatedPerson) {
		t.Errorf("Expected updated person %+v but got %+v", updatedPerson, foundPerson)
	}
}

// Тест удаления записи
func TestDeletePerson(t *testing.T) {
	dbConn := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(dbConn, t)

	repo := db.NewSqlitePersonRepo(dbConn)

	// Вставляем исходную запись
	fakePerson := tests.InsertFakePerson(dbConn, t)

	// Удаляем запись
	err := repo.Delete(fakePerson.ID)
	tests.MustFailIfError(t, err)

	// Пытаемся снова найти удалённую запись
	_, err = repo.FindByID(fakePerson.ID)
	if err == nil {
		t.Error("Deleted record was still found in the database.")
	}
}
