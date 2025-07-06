package db_test

import (
	"MyProjects/person_service/internal/app"
	"MyProjects/person_service/internal/db"
	"MyProjects/person_service/internal/tests"
	"context"
	"testing"
)

// Тест GetAll Получение списка пользователей
func TestGetAll(t *testing.T) {
	dbConn := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(dbConn, t)

	repo := db.NewPostgresPersonRepo(dbConn)

	// Вставляем нескольких фиктивных записей
	fakePerson1 := tests.InsertFakePerson(dbConn, t)
	fakePerson2 := tests.InsertFakePerson(dbConn, t)

	// Получаем все записи с передачей контекста
	persons, err := repo.GetAll(context.Background(), 10, 0)
	tests.MustFailIfError(t, err)

	// Проверяем количество полученных записей
	if len(persons) != 2 {
		t.Errorf("Expected two persons, but got %d", len(persons))
	}

	// Проверяем наличие наших фиктивных записей
	foundIDs := make(map[int]bool)
	for _, p := range persons {
		foundIDs[p.ID] = true
	}

	if !foundIDs[fakePerson1.ID] || !foundIDs[fakePerson2.ID] {
		t.Errorf("Inserted persons not found in result set")
	}
}

// Тест FindByID Поиск пользователя по ID
func TestFindByID(t *testing.T) {
	dbConn := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(dbConn, t)

	repo := db.NewPostgresPersonRepo(dbConn)

	// Вставляем фиктивную запись
	fakePerson := tests.InsertFakePerson(dbConn, t)

	// Поиск по ID с передачей контекста
	foundPerson, err := repo.FindByID(context.Background(), fakePerson.ID)
	tests.MustFailIfError(t, err)

	// Проверяем равенство полученной записи оригинальной
	if foundPerson.ID != fakePerson.ID ||
		foundPerson.Email != fakePerson.Email ||
		foundPerson.Phone != fakePerson.Phone ||
		foundPerson.FirstName != fakePerson.FirstName ||
		foundPerson.LastName != fakePerson.LastName {
		t.Errorf("Retrieved person does not match original")
	}

	// Проверяем случай несуществующего ID
	_, err = repo.FindByID(context.Background(), -1)
	if err == nil || err != db.ErrRecordNotFound {
		t.Errorf("Expected ErrRecordNotFound for non-existing ID")
	}
}

// Тест Create Создание новой записи в БД
func TestCreate(t *testing.T) {
	dbConn := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(dbConn, t)

	repo := db.NewPostgresPersonRepo(dbConn)

	newPerson := &app.Person{
		Email:     "john.doe@example.com",
		Phone:     "+79991234567",
		FirstName: "John",
		LastName:  "Doe",
	}

	err := repo.Create(context.Background(), newPerson)
	tests.MustFailIfError(t, err)

	// Проверяем, что новая запись появилась в базе
	retrievedPerson, err := repo.FindByID(context.Background(), newPerson.ID)
	tests.MustFailIfError(t, err)

	if retrievedPerson.Email != newPerson.Email ||
		retrievedPerson.Phone != newPerson.Phone ||
		retrievedPerson.FirstName != newPerson.FirstName ||
		retrievedPerson.LastName != newPerson.LastName {
		t.Errorf("Created person doesn't match the expected data")
	}
}

// Тест Update Обновление существующего пользователя
func TestUpdate(t *testing.T) {
	dbConn := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(dbConn, t)

	repo := db.NewPostgresPersonRepo(dbConn)

	// Вставляем фиктивную запись
	fakePerson := tests.InsertFakePerson(dbConn, t)

	// Обновляем
	updatedPerson := &app.Person{
		ID:        fakePerson.ID,
		Email:     "updated@yandex.ry",
		Phone:     "+79996965959",
		FirstName: "Лобовь",
		LastName:  "Прокофьевна",
	}

	// Обновляем запись
	err := repo.Update(context.Background(), fakePerson.ID, updatedPerson)
	tests.MustFailIfError(t, err)

	// Получаем обновлкнную запись
	retrievedPerson, err := repo.FindByID(context.Background(), fakePerson.ID)
	tests.MustFailIfError(t, err)

	// Сравниваем данные
	if retrievedPerson.Email != updatedPerson.Email ||
		retrievedPerson.Phone != updatedPerson.Phone ||
		retrievedPerson.FirstName != updatedPerson.FirstName ||
		retrievedPerson.LastName != updatedPerson.LastName {
		t.Errorf("Created person doesn't match the expected data")
	}

}

// Тест Delete Удаление пользователя по ID
func TestDelete(t *testing.T) {
	dbConn := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(dbConn, t)

	repo := db.NewPostgresPersonRepo(dbConn)

	// Вставляем фиктивную запись
	fakePerson := tests.InsertFakePerson(dbConn, t)

	// Удаляем запись
	err := repo.Delete(context.Background(), fakePerson.ID)
	tests.MustFailIfError(t, err)

	// Пробуем найти удалённую запись
	_, err = repo.FindByID(context.Background(), fakePerson.ID)
	if err == nil || err != db.ErrRecordNotFound {
		t.Errorf("Deleted record still exists or wrong error returned: actual error is %v", err)
	}
}
