
# Simple Person Service 

Простой RESTful CRUD-сервис для управления моделью Person с использованием PostgeSQL базы данных и фреймворка Echo.



## 🟢 Описание проекта 
Этот проект реализует полный цикл операций над сущностью Person: 
создание, чтение, обновление и удаление записей. Использует четкую архитектуру с выделением слоев, 
что позволяет легко поддерживать и расширять приложение.



### 🟢 Технологии

- **Framework:** Echo (версия 4.13.4)
- **Logging:** Uber Zap (версия 1.27.0)
- **Database:** PostgreSQL
- **Validation:** Validator (versión 10.27.0)
- **Documentation:** Swaggo (Swagger UI интеграция для Echo, версия 1.16.4)
- **Environment Variables Handling:** Godotenv (версия 1.5.1)
- **JSON Serialization:** ByteDance Sonic (быстрое JSON кодирование, версия 1.13.3)
- **Error Management:** MultiErr (версия 1.10.0)
- **Concurrency Tools:** Modern-Go Concurrent (версия bacd9c7ef1dd)
- **Performance Optimization:** Valyala bytebufferpool (версия 1.0.0)


### 🟢 Структура приложения 

- **Model:** Представлена структурой Person, содержащей ID, Email, Phone, First Name и Last Name.
- **Handlers:** HTTP-обработчики реализованы с использованием фреймворка Echo.
- **Logic:** Реализована логика для каждой операции CRUD.
- **Repository:** Используется слой доступа к данным, позволяющий абстрагироваться от конкретной СУБД (PostgeSQL).



### 🟢 Доступные endpoints 

- GET	/person/	Получение списка пользователей
- GET	/person/{id}	Получение конкретного пользователя
- POST	/person/	Создание новой записи пользователя
- PUT	/person/{id}	Обновление существующей записи пользователя
- DELETE	/person/{id}	Удаление записи пользователя



### 🟢 Установка

Чтобы запустить проект на своем компьютере, выполните следующие шаги:

1. **Клонирование репозитория:**
```shell
git clone https://github.com/lamauspex/Person_service
```

2. **В файле .env**
 
Замените **DB_PASSWORD=YOUR_DB_PASSWORD_HERE** на действующий пароль


4. **Установка зависимостей:**
```shell
go mod download
```


4. **Запуск**
```shell
go run main.go
```
```shell
Swagger API http://localhost:8080/swagger/index.html
```



Ваш вклад в проект приветствуется! Если вы хотите внести изменения или улучшения, создайте pull request или откройте issue на GitHub.

#### Контакты

Если у вас есть вопросы или предложения, не стесняйтесь связаться со мной:

- Имя: Резник Кирилл
- Email: lamauspex@yandex.ru
- GitHub: https://github.com/lamauspex
- Telegram: @lamauspex

Спасибо за интерес к проекту! Надеюсь, он будет полезен в вашей работе в цифровом маркетинге.


