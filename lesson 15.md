# Занятие 15: Работа с базами данных в Go

`"Подключаем SQL базы данных к нашему приложению"`

---

## 📝 План на сегодня

1.  **Пакет database/sql:** Стандартный интерфейс для работы с БД
2.  **SQL драйверы:** Подключение разных СУБД
3.  **Подключение к БД:** Создание пула соединений
4.  **CRUD операции:** Create, Read, Update, Delete
5.  **Подготовленные запросы:** Безопасность и производительность
6.  **Транзакции:** Группировка операций
7.  **Практика:** Создаем приложение с persistence слоем

---

## 1. Пакет database/sql

**database/sql** - стандартный пакет Go для работы с SQL базами данных.

### Основные типы:
- **sql.DB** - пул соединений с БД
- **sql.Stmt** - подготовленный запрос
- **sql.Tx** - транзакция
- **sql.Rows** - результат запроса с несколькими строками
- **sql.Row** - результат запроса с одной строкой

### Установка драйвера (пример для PostgreSQL):
```bash
go get github.com/lib/pq
```

---

## 2. Подключение к базе данных

### Создание подключения:
```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq" // Драйвер PostgreSQL
)

func main() {
    // Параметры подключения
    connStr := "host=localhost port=5432 user=myuser password=mypassword dbname=mydb sslmode=disable"
    
    // Открываем соединение
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Проверяем подключение
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("✅ Успешно подключились к базе данных!")
}
```

### Другие популярные драйверы:
```go
// MySQL
import _ "github.com/go-sql-driver/mysql"
connStr := "user:password@tcp(localhost:3306)/dbname"

// SQLite
import _ "modernc.org/sqlite"
connStr := "file:test.db"

// SQL Server
import _ "github.com/denisenkom/go-mssqldb"
connStr := "server=localhost;user id=sa;password=yourpassword;database=yourdb"
```

---

## 3. CRUD операции

### Создание таблицы:
```go
func createTables(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name VARCHAR(200) NOT NULL,
        price DECIMAL(10,2) NOT NULL,
        stock INTEGER DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
    
    _, err := db.Exec(query)
    return err
}
```

---

## 4. Create (Вставка данных)

### Простая вставка:
```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func createUser(db *sql.DB, user *User) error {
    query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
    
    err := db.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
    if err != nil {
        return err
    }
    
    return nil
}
```

### Множественная вставка:
```go
func createUsers(db *sql.DB, users []User) error {
    query := `INSERT INTO users (name, email) VALUES ($1, $2)`
    
    for _, user := range users {
        _, err := db.Exec(query, user.Name, user.Email)
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## 5. Read (Чтение данных)

### Чтение одной строки:
```go
func getUserByID(db *sql.DB, id int) (*User, error) {
    var user User
    query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
    
    err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }
    
    return &user, nil
}
```

### Чтение нескольких строк:
```go
func getAllUsers(db *sql.DB) ([]User, error) {
    var users []User
    query := `SELECT id, name, email, created_at FROM users ORDER BY created_at DESC`
    
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }
    
    return users, nil
}
```

---

## 6. Update (Обновление данных)

```go
func updateUser(db *sql.DB, user *User) error {
    query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
    
    result, err := db.Exec(query, user.Name, user.Email, user.ID)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user not found")
    }
    
    return nil
}
```

---

## 7. Delete (Удаление данных)

```go
func deleteUser(db *sql.DB, id int) error {
    query := `DELETE FROM users WHERE id = $1`
    
    result, err := db.Exec(query, id)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user not found")
    }
    
    return nil
}
```

---

## 🎯 Практика 1: Менеджер пользователей с БД

**Задача:** Создать полную систему управления пользователями с PostgreSQL

```go
package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"
    
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

type UserService struct {
    db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
    return &UserService{db: db}
}

func (s *UserService) Init() error {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
    
    _, err := s.db.Exec(query)
    return err
}

func (s *UserService) Create(user *User) error {
    query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at`
    
    err := s.db.QueryRow(query, user.Name, user.Email).Scan(&user.ID, &user.CreatedAt)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

func (s *UserService) GetByID(id int) (*User, error) {
    var user User
    query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
    
    err := s.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, nil
}

func (s *UserService) GetAll() ([]User, error) {
    query := `SELECT id, name, email, created_at FROM users ORDER BY created_at DESC`
    
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to get users: %w", err)
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, user)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }
    
    return users, nil
}

func (s *UserService) Update(user *User) error {
    query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
    
    result, err := s.db.Exec(query, user.Name, user.Email, user.ID)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user not found")
    }
    
    return nil
}

func (s *UserService) Delete(id int) error {
    query := `DELETE FROM users WHERE id = $1`
    
    result, err := s.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user not found")
    }
    
    return nil
}

type UserHandler struct {
    service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    if user.Name == "" || user.Email == "" {
        http.Error(w, "Name and email are required", http.StatusBadRequest)
        return
    }
    
    if err := h.service.Create(&user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    user, err := h.service.GetByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.service.GetAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    user.ID = id
    
    if err := h.service.Update(&user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    if err := h.service.Delete(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

func connectDB() (*sql.DB, error) {
    // Для демонстрации используем SQLite, но можно заменить на PostgreSQL
    // connStr := "host=localhost port=5432 user=myuser password=mypassword dbname=mydb sslmode=disable"
    // db, err := sql.Open("postgres", connStr)
    
    // Используем SQLite для простоты демонстрации
    db, err := sql.Open("sqlite", "file:users.db?cache=shared&mode=memory")
    if err != nil {
        return nil, err
    }
    
    if err = db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}

func main() {
    // Подключаемся к базе данных
    db, err := connectDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Инициализируем сервис
    userService := NewUserService(db)
    if err := userService.Init(); err != nil {
        log.Fatal("Failed to init database:", err)
    }
    
    // Создаем несколько тестовых пользователей
    sampleUsers := []User{
        {Name: "Анна", Email: "anna@example.com"},
        {Name: "Петр", Email: "petr@example.com"},
        {Name: "Мария", Email: "maria@example.com"},
    }
    
    for i := range sampleUsers {
        userService.Create(&sampleUsers[i])
    }
    
    // Настраиваем HTTP handlers
    userHandler := NewUserHandler(userService)
    
    r := mux.NewRouter()
    
    // User routes
    r.HandleFunc("/api/users", userHandler.GetAllUsers).Methods("GET")
    r.HandleFunc("/api/users", userHandler.CreateUser).Methods("POST")
    r.HandleFunc("/api/users/{id}", userHandler.GetUser).Methods("GET")
    r.HandleFunc("/api/users/{id}", userHandler.UpdateUser).Methods("PUT")
    r.HandleFunc("/api/users/{id}", userHandler.DeleteUser).Methods("DELETE")
    
    // Health check
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        if err := db.Ping(); err != nil {
            http.Error(w, "Database connection failed", http.StatusServiceUnavailable)
            return
        }
        json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
    }).Methods("GET")
    
    // Documentation
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        docs := `
        <h1>User Management API</h1>
        <h2>Endpoints:</h2>
        <ul>
            <li><strong>GET /api/users</strong> - Get all users</li>
            <li><strong>POST /api/users</strong> - Create a user</li>
            <li><strong>GET /api/users/{id}</strong> - Get user by ID</li>
            <li><strong>PUT /api/users/{id}</strong> - Update user</li>
            <li><strong>DELETE /api/users/{id}</strong> - Delete user</li>
        </ul>
        `
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprint(w, docs)
    })
    
    port := ":8080"
    log.Printf("🚀 User Management API started on http://localhost%s", port)
    log.Fatal(http.ListenAndServe(port, r))
}
```

---

## 🎯 Практика 2: Блог с категориями и комментариями

**Задача:** Создать систему блога с реляционными связями

```go
package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"
    
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type Category struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type Post struct {
    ID         int       `json:"id"`
    Title      string    `json:"title"`
    Content    string    `json:"content"`
    CategoryID int       `json:"category_id"`
    Category   *Category `json:"category,omitempty"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type Comment struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    Author    string    `json:"author"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

type BlogService struct {
    db *sql.DB
}

func NewBlogService(db *sql.DB) *BlogService {
    return &BlogService{db: db}
}

func (s *BlogService) Init() error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS categories (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) UNIQUE NOT NULL
        )`,
        
        `CREATE TABLE IF NOT EXISTS posts (
            id SERIAL PRIMARY KEY,
            title VARCHAR(200) NOT NULL,
            content TEXT NOT NULL,
            category_id INTEGER REFERENCES categories(id),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
        
        `CREATE TABLE IF NOT EXISTS comments (
            id SERIAL PRIMARY KEY,
            post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
            author VARCHAR(100) NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
        
        `INSERT INTO categories (name) VALUES 
            ('Technology'), ('Science'), ('Art'), ('Travel')
        ON CONFLICT (name) DO NOTHING`,
    }
    
    for _, query := range queries {
        _, err := s.db.Exec(query)
        if err != nil {
            return err
        }
    }
    
    return nil
}

// Category methods
func (s *BlogService) CreateCategory(category *Category) error {
    query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
    return s.db.QueryRow(query, category.Name).Scan(&category.ID)
}

func (s *BlogService) GetCategories() ([]Category, error) {
    rows, err := s.db.Query("SELECT id, name FROM categories ORDER BY name")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var categories []Category
    for rows.Next() {
        var cat Category
        if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
            return nil, err
        }
        categories = append(categories, cat)
    }
    return categories, nil
}

// Post methods
func (s *BlogService) CreatePost(post *Post) error {
    query := `INSERT INTO posts (title, content, category_id) 
              VALUES ($1, $2, $3) 
              RETURNING id, created_at, updated_at`
    
    return s.db.QueryRow(query, post.Title, post.Content, post.CategoryID).
        Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (s *BlogService) GetPostWithCategory(id int) (*Post, error) {
    query := `
        SELECT p.id, p.title, p.content, p.category_id, p.created_at, p.updated_at,
               c.id, c.name
        FROM posts p
        LEFT JOIN categories c ON p.category_id = c.id
        WHERE p.id = $1
    `
    
    var post Post
    var category Category
    
    err := s.db.QueryRow(query, id).Scan(
        &post.ID, &post.Title, &post.Content, &post.CategoryID, 
        &post.CreatedAt, &post.UpdatedAt,
        &category.ID, &category.Name,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("post not found")
        }
        return nil, err
    }
    
    post.Category = &category
    return &post, nil
}

func (s *BlogService) GetPostsWithCategories() ([]Post, error) {
    query := `
        SELECT p.id, p.title, p.content, p.category_id, p.created_at, p.updated_at,
               c.id, c.name
        FROM posts p
        LEFT JOIN categories c ON p.category_id = c.id
        ORDER BY p.created_at DESC
    `
    
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var posts []Post
    for rows.Next() {
        var post Post
        var category Category
        
        err := rows.Scan(
            &post.ID, &post.Title, &post.Content, &post.CategoryID,
            &post.CreatedAt, &post.UpdatedAt,
            &category.ID, &category.Name,
        )
        if err != nil {
            return nil, err
        }
        
        post.Category = &category
        posts = append(posts, post)
    }
    
    return posts, nil
}

// Comment methods
func (s *BlogService) CreateComment(comment *Comment) error {
    query := `INSERT INTO comments (post_id, author, content) 
              VALUES ($1, $2, $3) 
              RETURNING id, created_at`
    
    return s.db.QueryRow(query, comment.PostID, comment.Author, comment.Content).
        Scan(&comment.ID, &comment.CreatedAt)
}

func (s *BlogService) GetCommentsByPostID(postID int) ([]Comment, error) {
    rows, err := s.db.Query(`
        SELECT id, post_id, author, content, created_at 
        FROM comments 
        WHERE post_id = $1 
        ORDER BY created_at DESC
    `, postID)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var comments []Comment
    for rows.Next() {
        var comment Comment
        if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Author, 
            &comment.Content, &comment.CreatedAt); err != nil {
            return nil, err
        }
        comments = append(comments, comment)
    }
    
    return comments, nil
}

// Транзакционный пример: создание поста с комментарием
func (s *BlogService) CreatePostWithInitialComment(post *Post, comment *Comment) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Создаем пост
    err = tx.QueryRow(
        `INSERT INTO posts (title, content, category_id) VALUES ($1, $2, $3) 
         RETURNING id, created_at, updated_at`,
        post.Title, post.Content, post.CategoryID,
    ).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
    
    if err != nil {
        return err
    }
    
    // Создаем комментарий
    comment.PostID = post.ID
    err = tx.QueryRow(
        `INSERT INTO comments (post_id, author, content) VALUES ($1, $2, $3) 
         RETURNING id, created_at`,
        comment.PostID, comment.Author, comment.Content,
    ).Scan(&comment.ID, &comment.CreatedAt)
    
    if err != nil {
        return err
    }
    
    return tx.Commit()
}

type BlogHandler struct {
    service *BlogService
}

func NewBlogHandler(service *BlogService) *BlogHandler {
    return &BlogHandler{service: service}
}

func (h *BlogHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
    posts, err := h.service.GetPostsWithCategories()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

func (h *BlogHandler) GetPost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }
    
    post, err := h.service.GetPostWithCategory(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

func (h *BlogHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
    var post Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    if post.Title == "" || post.Content == "" {
        http.Error(w, "Title and content are required", http.StatusBadRequest)
        return
    }
    
    if err := h.service.CreatePost(&post); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}

func (h *BlogHandler) GetComments(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["postId"])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }
    
    comments, err := h.service.GetCommentsByPostID(postID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comments)
}

func (h *BlogHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["postId"])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }
    
    var comment Comment
    if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    comment.PostID = postID
    
    if comment.Author == "" || comment.Content == "" {
        http.Error(w, "Author and content are required", http.StatusBadRequest)
        return
    }
    
    if err := h.service.CreateComment(&comment); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(comment)
}

func main() {
    // Подключаемся к базе данных
    db, err := sql.Open("postgres", "host=localhost port=5432 user=bloguser dbname=blogdb sslmode=disable")
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Проверяем подключение
    if err := db.Ping(); err != nil {
        log.Fatal("Database ping failed:", err)
    }
    
    // Инициализируем сервис
    blogService := NewBlogService(db)
    if err := blogService.Init(); err != nil {
        log.Fatal("Failed to init database:", err)
    }
    
    // Настраиваем HTTP handlers
    blogHandler := NewBlogHandler(blogService)
    
    r := mux.NewRouter()
    
    // Post routes
    r.HandleFunc("/api/posts", blogHandler.GetPosts).Methods("GET")
    r.HandleFunc("/api/posts", blogHandler.CreatePost).Methods("POST")
    r.HandleFunc("/api/posts/{id}", blogHandler.GetPost).Methods("GET")
    
    // Comment routes
    r.HandleFunc("/api/posts/{postId}/comments", blogHandler.GetComments).Methods("GET")
    r.HandleFunc("/api/posts/{postId}/comments", blogHandler.CreateComment).Methods("POST")
    
    // Categories
    r.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
        categories, err := blogService.GetCategories()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(categories)
    }).Methods("GET")
    
    port := ":8080"
    log.Printf("🚀 Blog API started on http://localhost%s", port)
    log.Fatal(http.ListenAndServe(port, r))
}
```

---

## ❓ Важные моменты

### Безопасность:
```go
// Всегда используйте подготовленные запросы для предотвращения SQL инъекций
// ПЛОХО:
query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", userName)

// ХОРОШО:
query := "SELECT * FROM users WHERE name = $1"
db.QueryRow(query, userName)
```

### Управление соединениями:
```go
// Настройка пула соединений
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

---

## 🏠 Домашнее задание

**Задача 1: Интернет-магазин**
Создайте базу данных для интернет-магазина с таблицами:
- Products (товары)
- Customers (клиенты)
- Orders (заказы)
- OrderItems (позиции заказа)
  Реализуйте CRUD операции для каждой таблицы

**Задача 2: Система бронирования**
Разработайте систему бронирования отелей с:
- Отелями и номерами
- Клиентами
- Бронированиями
- Доступностью номеров по датам

**Задача 3: Оптимизация запросов**
Для вашего блога добавьте:
- Пагинацию для постов и комментариев
- Поиск по заголовкам и содержимому
- Статистику (количество постов по категориям)

---

## 🚀 Что ждет на следующем занятии?

*   **Аутентификация и авторизация:** JWT токены
*   **Middleware для аутентификации**
*   **Хеширование паролей:** bcrypt

**Удачи в работе с базами данных! 🎉**