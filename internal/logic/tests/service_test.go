package tests

import (
	"MyProjects/person_service/internal/app"
	"MyProjects/person_service/internal/db"
	"MyProjects/person_service/internal/logic"
	"MyProjects/person_service/internal/tests"
	"context"
	"testing"
)

// Тест AddPerson Добавить пользователя
func TestAddPerson(t *testing.T) {
	// Подготовка тестовой базы данных
	testDb := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(testDb, t)
	repo := db.NewPostgresPersonRepo(testDb)
	service := logic.NewDefaultPersonService(repo)

	// Новый экземпляр пользователя
	newPerson := &app.Person{
		Email:     "new_user@test.ru",
		Phone:     "+79876543210",
		FirstName: "Игорь",
		LastName:  "Семенов",
	}

	// Контекст необходим для вызова сервисов
	ctx := context.Background()

	// Добавляем пользователя
	err := service.AddPerson(ctx, newPerson)
	tests.MustFailIfError(t, err)

	// Получаем сохранённого пользователя по новому ID
	fetchedPerson, err := service.GetPerson(ctx, newPerson.ID)
	tests.MustFailIfError(t, err)

	// Проверка правильности сохранения данных
	if fetchedPerson.Email != newPerson.Email ||
		fetchedPerson.Phone != newPerson.Phone ||
		fetchedPerson.FirstName != newPerson.FirstName ||
		fetchedPerson.LastName != newPerson.LastName {
		t.Fatalf("Полученный пользователь отличается от ожидаемого:\n%+v\nожидаемый:%+v", fetchedPerson, newPerson)
	}
}

// Тест GetPerson Получить пользователя по ID
func TestGetPerson(t *testing.T) {
	// Подготовка тестовой базы данных
	testDb := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(testDb, t)
	repo := db.NewPostgresPersonRepo(testDb)
	service := logic.NewDefaultPersonService(repo)

	// Создание фиктивного пользователя
	fakePerson := tests.InsertFakePerson(testDb, t)

	// Запрашиваем пользователя по ID
	ctx := context.Background()
	gottenPerson, err := service.GetPerson(ctx, fakePerson.ID)
	tests.MustFailIfError(t, err)

	// Проверяем правильность полученных данных
	if gottenPerson.Email != fakePerson.Email ||
		gottenPerson.Phone != fakePerson.Phone ||
		gottenPerson.FirstName != fakePerson.FirstName ||
		gottenPerson.LastName != fakePerson.LastName {
		t.Fatalf("Данные пользователя отличаются от ожидаемых.\nОжидаемые значения:\n %+v\nФактические значения:\n %+v", fakePerson, gottenPerson)
	}
}

// Тест UpdatePerson Обновление пользователя
func TestUpdatePerson(t *testing.T) {
	// Подготовка тестовой базы данных
	testDb := tests.PrepareTestDatabase(t)
	defer tests.CleanUp(testDb, t)
	repo := db.NewPostgresPersonRepo(testDb)
	service := logic.NewDefaultPersonService(repo)

	// Создание фиктивного пользователя
	fakePerson := tests.InsertFakePerson(testDb, t)

	// Новые данные пользователя
	updatedData := &app.Person{
		ID:        fakePerson.ID,
		Email:     "updated_email@test.ru",
		Phone:     "+79876543210",
		FirstName: "Иван",
		LastName:  "Иванушка",
	}

	// Обновляем пользователя
	ctx := context.Background()
	err := service.EditPerson(ctx, fakePerson.ID, updatedData)
	tests.MustFailIfError(t, err)

	// Проверяем обновление данных
	updatedPerson, err := service.GetPerson(ctx, fakePerson.ID)
	tests.MustFailIfError(t, err)

	if updatedPerson.Email != updatedData.Email ||
		updatedPerson.Phone != updatedData.Phone ||
		updatedPerson.FirstName != updatedData.FirstName ||
		updatedPerson.LastName != updatedData.LastName {
		t.Fatalf("Данные пользователя не были успешно обновлены.\nОжидаемые значения:\n %+v\nФактические значения:\n %+v", updatedData, updatedPerson)
	}
}
