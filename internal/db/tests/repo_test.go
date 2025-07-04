package db_test

import (
	"MyProjects/person_service/internal/db"
	"MyProjects/person_service/internal/tests"
	"testing"
)

// Тест GetAll (получение списка персон)
func TestGetAll(t *testing.T) {
	dbConn := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(dbConn, t)

	repo := db.NewSqlitePersonRepo(dbConn)

	// Вставка нескольких фиктивных записей
	fakePerson1 := tests.InsertFakePerson(dbConn, t)
	fakePerson2 := tests.InsertFakePerson(dbConn, t)

	// Получаем все записи
	persons, err := repo.GetAll(10, 0)
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

// Тест FindByID (поиск одной записи по ID)
func TestFindByID(t *testing.T) {
	dbConn := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(dbConn, t)

	repo := db.NewSqlitePersonRepo(dbConn)

	// Вставка фиктивной записи
	fakePerson := tests.InsertFakePerson(dbConn, t)

	// Поиск по ID
	foundPerson, err := repo.FindByID(fakePerson.ID)
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
	_, err = repo.FindByID(-1)
	if err == nil || err != db.ErrRecordNotFound {
		t.Errorf("Expected ErrRecordNotFound for non-existing ID")
	}
}
