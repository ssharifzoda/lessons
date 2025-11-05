# Занятие 14: Введение в веб-разработку на Go

# Создаем первый веб-сервер и REST API

---

📝 План на сегодня

1. Пакет net/http: Основы HTTP в Go
2. Создание HTTP сервера: ListenAndServe
3. Обработчики (Handlers): Функции для обработки запросов
4. Роутинг: Маршрутизация запросов
5. Работа с HTTP методами: GET, POST, PUT, DELETE
6. Практика: Создаем REST API для управления задачами

---

## 1. Пакет net/http - основы

net/http - стандартный пакет Go для работы с HTTP.

Основные компоненты:
- http.Handler - интерфейс для обработчиков
- http.Server - структура сервера
- http.Request - входящий запрос
- http.ResponseWriter - для отправки ответа

#### Простой HTTP handler:
```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Привет, мир!")
}
```

---

## 2. Создание HTTP сервера

Базовый сервер:
```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Регистрируем обработчики
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/about", aboutHandler)

	// Запускаем сервер
	fmt.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
	`<h1> Добро пожаловать!</h1>
	<p> Это главная страница </p>
	<ul>
	    <li> <a href = "/hello">Приветствие</a></li>
        <li> <a href = "/about" > О нас </a></li>
	</ul>`)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Гость"
	}
	fmt.Fprintf(w, "<h1>Привет, %s!</h1>", name)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
	`<h1>О нас</h1>
<p> Мы изучаем веб - разработку на Go! </p>
<a href = "/" > На главную </a >`)
}
```


---

## 3. Обработчики (Handlers)

Типы обработчиков:

#### 1. Функция-обработчик (HandlerFunc):
```go
func simpleHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Простой ответ"))
}
```
#### 2. Структура-обработчик (реализует http.Handler)
```go
type MyHandler struct {
    message string
}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Сообщение: %s", h.message)
}
```

Использование:
```go
func main() {
    // Функция-обработчик
    http.HandleFunc("/simple", simpleHandler)

    // Структура-обработчик
    customHandler := MyHandler{message: "Привет из обработчика"}
    http.Handle("/custom", customHandler)
    
    http.ListenAndServe(":8080", nil)
}
```

---

## 4. Роутинг

#### Базовый роутинг с mux:

```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/users/", usersHandler)
	mux.HandleFunc("/posts/", postsHandler)
	mux.HandleFunc("/api/", apiHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Главная страница")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Страница пользователей")
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Страница постов")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API endpoint")
}
```

---

## 5. Работа с HTTP методами

#### Обработка разных HTTP методов:

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getUser(w, r)
    case http.MethodPost:
        createUser(w, r)
    case http.MethodPut:
        updateUser(w, r)
    case http.MethodDelete:
        deleteUser(w, r)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

func getUser(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    fmt.Fprintf(w, "Получение пользователя с ID: %s", id)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    // Чтение данных из тела запроса
    name := r.FormValue("name")
    email := r.FormValue("email")
    
    fmt.Fprintf(w, "Создание пользователя: %s (%s)", name, email)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Обновление пользователя")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Удаление пользователя")
}
```

---

## 🎯 Практика 1: Простой веб-сервер с API

#### Задача: Создать сервер с HTML страницами и JSON API

```go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type User struct {
	ID int    json:"id"
	Name  string json:"name"
	Email string json:"email"
}

var users = []User{
	{ID: 1, Name: "Анна", Email: "anna@example.com"},
	{ID: 2, Name: "Петр", Email: "petr@example.com"},
	{ID: 3, Name: "Мария", Email: "maria@example.com"},
}

var tmpl = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .user { border: 1px solid #ddd; padding: 10px; margin: 10px 0; }
        .api-link { background: #f0f0f0; padding: 10px; margin: 10px 0; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{.Content}}
</body>
</html>
))`

func renderTemplate(w http.ResponseWriter, title, content string) {
	data := struct {
		Title   string
		Content template.HTML
	}{
		Title:   title,
		Content: template.HTML(content),
	}

	tmpl.Execute(w, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	content := ` 
        <p>Добро пожаловать на наш сайт!</p>
        <div class="api-link">
            <h3>Доступные API endpoints:</h3>
            <ul>
                <li><a href="/api/users">GET /api/users</a> - список пользователей</li>
                <li><a href="/api/users/1">GET /api/users/{id}</a> - пользователь по ID</li>
            </ul>
        </div>
        <p><a href="/users">Посмотреть пользователей</a></p>`

	renderTemplate(w, "Главная страница", content)
}

func usersPageHandler(w http.ResponseWriter, r *http.Request) {
	var userList string
	for _, user := range users {
		userList += fmt.Sprintf(`
<div class="user">
<h3>%s</h3>
<p>Email: %s</p>
<p>ID: %d</p>
<a href="/api/users/%d">JSON</a>
</div>`, user.Name, user.Email, user.ID, user.ID)
	}

	content := fmt.Sprintf(`
        <h2>Список пользователей</h2>
        %s
        <p><a href="/">На главную</a></p>
    `, userList)

	renderTemplate(w, "Пользователи", content)
}

func apiUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(users)
	case http.MethodPost:
		var newUser User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)

			return
		}

		newUser.ID = len(users) + 1
		users = append(users, newUser)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func apiUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Извлекаем ID из URL
	idStr := r.URL.Path[len("/api/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	// Ищем пользователя
	var foundUser *User
	for _, user := range users {
		if user.ID == id {
			foundUser = &user
			break
		}
	}

	if foundUser == nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(foundUser)
	case http.MethodPut:
		var updatedUser User
		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}

		updatedUser.ID = id
		// В реальном приложении здесь было бы обновление в базе данных
		for i := range users {
			if users[i].ID == id {
				users[i] = updatedUser
				break
			}
		}

		json.NewEncoder(w).Encode(updatedUser)
	case http.MethodDelete:
		// Удаляем пользователя
		for i, user := range users {
			if user.ID == id {
				users = append(users[:i], users[i+1:]...)
				break
			}
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()

	// HTML страницы
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/users", usersPageHandler)

	// API endpoints
	mux.HandleFunc("/api/users", apiUsersHandler)
	mux.HandleFunc("/api/users/", apiUserHandler)

	// Статические файлы
	mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	fmt.Println("🚀 Сервер запущен на http://localhost:8080")
	fmt.Println("📄 Главная страница: http://localhost:8080")
	fmt.Println("👥 Пользователи: http://localhost:8080/users")
	fmt.Println("🔗 API: http://localhost:8080/api/users")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("❌ Ошибка запуска сервера: %v\n", err)
	}
}
```
---

## 🎯 Практика 2: REST API для управления задачами

Задача: Создать полноценное REST API для менеджера задач

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Task struct {
	ID int       json:"id"
	Title       string    json:"title"
	Description string    json:"description"
	Completed   bool      json:"completed"
	CreatedAt   time.Time json:"created_at"
	DueDate     time.Time json:"due_date,omitempty"
	Priority    string    json:"priority" // low, medium, high
}

type TaskStore struct {
	sync.RWMutex
	tasks  map[int]Task
	nextID int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

func (ts *TaskStore) Create(task Task) Task {
	ts.Lock()
	defer ts.Unlock()

	task.ID = ts.nextID
	task.CreatedAt = time.Now()
	ts.tasks[task.ID] = task
	ts.nextID++

	return task
}

func (ts *TaskStore) Get(id int) (Task, bool) {
	ts.RLock()
	defer ts.RUnlock()

	task, exists := ts.tasks[id]
	return task, exists
}

func (ts *TaskStore) GetAll() []Task {
	ts.RLock()
	defer ts.RUnlock()

	tasks := make([]Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (ts *TaskStore) Update(id int, updated Task) (Task, bool) {
	ts.Lock()
	defer ts.Unlock()

	if _, exists := ts.tasks[id]; !exists {
		return Task{}, false
	}

	updated.ID = id
	ts.tasks[id] = updated
	return updated, true
}

func (ts *TaskStore) Delete(id int) bool {
	ts.Lock()
	defer ts.Unlock()

	if _, exists := ts.tasks[id]; exists {
		delete(ts.tasks, id)
		return true
	}
	return false
}

func (ts *TaskStore) GetByStatus(completed bool) []Task {
	ts.RLock()
	defer ts.RUnlock()

	var filtered []Task
	for _, task := range ts.tasks {
		if task.Completed == completed {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

type TaskServer struct {
	store *TaskStore
}

func NewTaskServer() *TaskServer {
	return &TaskServer{
		store: NewTaskStore(),
	}
}

func (ts *TaskServer) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if task.Priority == "" {
		task.Priority = "medium"
	}

	created := ts.store.Create(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (ts *TaskServer) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, exists := ts.store.Get(id)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (ts *TaskServer) getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Поддержка фильтрации по статусу
	completedStr := r.URL.Query().Get("completed")
	var tasks []Task

	if completedStr != "" {
		completed, err := strconv.ParseBool(completedStr)
		if err != nil {
			http.Error(w, "Invalid completed parameter", http.StatusBadRequest)
			return
		}
		tasks = ts.store.GetByStatus(completed)
	} else {
		tasks = ts.store.GetAll()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (ts *TaskServer) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updated Task
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	task, exists := ts.store.Update(id, updated)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (ts *TaskServer) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if !ts.store.Delete(id) {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ts *TaskServer) completeTaskHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, exists := ts.store.Get(id)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	task.Completed = true
	ts.store.Update(id, task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	server := NewTaskServer()

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/tasks", server.getAllTasksHandler).Methods("GET")
	api.HandleFunc("/tasks", server.createTaskHandler).Methods("POST")
	api.HandleFunc("/tasks/{id}", server.getTaskHandler).Methods("GET")
	api.HandleFunc("/tasks/{id}", server.updateTaskHandler).Methods("PUT")
	api.HandleFunc("/tasks/{id}", server.deleteTaskHandler).Methods("DELETE")
	api.HandleFunc("/tasks/{id}/complete", server.completeTaskHandler).Methods("PATCH")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	// Documentation
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		docs := `
        <h1>Task Manager API</h1>
        <h2>Endpoints:</h2>
        <ul>
            <li><strong>GET /api/v1/tasks</strong> - Get all tasks</li>
            <li><strong>POST /api/v1/tasks</strong> - Create a task</li>
            <li><strong>GET /api/v1/tasks/{id}</strong> - Get a task by ID</li>
            <li><strong>PUT /api/v1/tasks/{id}</strong> - Update a task</li>
            <li><strong>DELETE /api/v1/tasks/{id}</strong> - Delete a task</li>
            <li><strong>PATCH /api/v1/tasks/{id}/complete</strong> - Mark task as completed</li>
        </ul>
        <p>Use query parameter ?completed=true/false to filter tasks</p>`

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, docs)
	})

	// Add CORS middleware
	handler := enableCORS(r)

	port := ":8080"
	log.Printf("🚀 Task Manager API server started on http://localhost%s", port)
	log.Printf("📚 API documentation: http://localhost%s", port)

	log.Fatal(http.ListenAndServe(port, handler))
}
```

```bash
bash
# Тестирование API с curl
curl -X GET http://localhost:8080/api/v1/tasks
curl -X POST http://localhost:8080/api/v1/tasks -H "Content-Type: application/json" -d '{"title":"Learn Go","description":"Study web development","priority":"high"}'
curl -X GET http://localhost:8080/api/v1/tasks/1
curl -X PUT http://localhost:8080/api/v1/tasks/1 -H "Content-Type: application/json" -d '{"title":"Learn Go Web","description":"Study web development with Go","completed":false,"priority":"high"}'
curl -X PATCH http://localhost:8080/api/v1/tasks/1/complete
curl -X DELETE http://localhost:8080/api/v1/tasks/1
```
---

❓ Важные моменты

Безопасность заголовков:

```go
// Всегда устанавливайте Content-Type
w.Header().Set("Content-Type", "application/json")

// Для JSON API
w.Header().Set("Content-Type", "application/json; charset=utf-8")
```

Обработка ошибок:
```go
func jsonError(w http.ResponseWriter, message string, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}
```

---

## 🏠 Домашнее задание

#### Задача 1: Блог API
Создайте REST API для блога с endpoints для:
- CRUD операций с постами
- Комментарии к постам
- Категории постов
- Поиск и фильтрация