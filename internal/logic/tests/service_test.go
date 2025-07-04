package tests

import (
	"MyProjects/person_service/internal/app"
	"MyProjects/person_service/internal/db"
	"MyProjects/person_service/internal/logic"
	"MyProjects/person_service/internal/tests"
	"testing"
)

// TestAddPerson проверяет успешность добавления нового пользователя
func TestAddPerson(t *testing.T) {
	// Подготовка тестовой базы данных
	testDb := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(testDb, t)
	repo := db.NewSqlitePersonRepo(testDb)
	service := logic.NewDefaultPersonService(repo)

	newPerson := &app.Person{
		Email:     "new_user@test.ru",
		Phone:     "+79876543210",
		FirstName: "Игорь",
		LastName:  "Семенов",
	}

	// Выполняем операцию добавления пользователя
	err := service.AddPerson(newPerson)
	tests.MustFailIfError(t, err)

	// Проверяем наличие пользователя в базе данных
	fetchedPerson, err := service.GetPerson(newPerson.ID)
	tests.MustFailIfError(t, err)

	// Убеждаемся, что информация совпадает
	if fetchedPerson.Email != newPerson.Email ||
		fetchedPerson.Phone != newPerson.Phone ||
		fetchedPerson.FirstName != newPerson.FirstName ||
		fetchedPerson.LastName != newPerson.LastName {
		t.Fatalf("Полученный пользователь отличается от ожидаемого:\n%+v\nожидаемый:%+v", fetchedPerson, newPerson)
	}
}

// TestGetPerson проверяет успешность получения пользователя по ID
func TestGetPerson(t *testing.T) {
	// Подготовка тестовой базы данных
	testDb := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(testDb, t)
	repo := db.NewSqlitePersonRepo(testDb)
	service := logic.NewDefaultPersonService(repo)

	// Создание фиктивного пользователя
	fakePerson := tests.InsertFakePerson(testDb, t)

	// Запрашиваем пользователя по ID
	gottenPerson, err := service.GetPerson(fakePerson.ID)
	tests.MustFailIfError(t, err)

	// Проверяем совпадение полей
	if gottenPerson.Email != fakePerson.Email ||
		gottenPerson.Phone != fakePerson.Phone ||
		gottenPerson.FirstName != fakePerson.FirstName ||
		gottenPerson.LastName != fakePerson.LastName {
		t.Fatalf("Данные пользователя отличаются от ожидаемых.\nОжидаемые значения:\n %+v\nФактические значения:\n %+v", fakePerson, gottenPerson)
	}
}

// TestUpdatePerson проверяет успешность обновления пользователя
func TestUpdatePerson(t *testing.T) {
	// Подготовка тестовой базы данных
	testDb := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(testDb, t)
	repo := db.NewSqlitePersonRepo(testDb)
	service := logic.NewDefaultPersonService(repo)

	// Создание фиктивного пользователя
	fakePerson := tests.InsertFakePerson(testDb, t)

	// Изменённые данные пользователя
	updatedData := &app.Person{
		ID:        fakePerson.ID,
		Email:     "updated_email@test.ru",
		Phone:     "+79876543210",
		FirstName: "Иван",
		LastName:  "Новиков",
	}

	// Обновляем пользователя
	err := service.EditPerson(fakePerson.ID, updatedData)
	tests.MustFailIfError(t, err)

	// Проверяем изменение данных
	updatedPerson, err := service.GetPerson(fakePerson.ID)
	tests.MustFailIfError(t, err)

	if updatedPerson.Email != updatedData.Email ||
		updatedPerson.Phone != updatedData.Phone ||
		updatedPerson.FirstName != updatedData.FirstName ||
		updatedPerson.LastName != updatedData.LastName {
		t.Fatalf("Данные пользователя не были успешно обновлены.\nОжидаемые значения:\n %+v\nФактические значения:\n %+v", updatedData, updatedPerson)
	}
}
