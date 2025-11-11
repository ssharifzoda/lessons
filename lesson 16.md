# Занятие 16: Аутентификация и авторизация в Go

`"Защищаем наше API с помощью JWT и middleware"`

---

## 📝 План на сегодня

1.  **Аутентификация vs Авторизация:** В чем разница?
2.  **JWT токены:** Структура и принцип работы
3.  **Хеширование паролей:** bcrypt для безопасности
4.  **Middleware для аутентификации:** Защита маршрутов
5.  **Сессии и куки:** Альтернативный подход
6.  **Практика:** Создаем полную систему аутентификации

---

## 1. Аутентификация vs Авторизация

### Аутентификация (Authentication):
**"Кто ты?"** - проверка личности пользователя
- Логин/пароль
- JWT токены
- OAuth

### Авторизация (Authorization):
**"Что тебе разрешено?"** - проверка прав доступа
- Роли пользователей
- Разрешения
- ACL (Access Control Lists)

---

## 2. JWT токены

**JWT (JSON Web Token)** - стандарт для создания токенов доступа.

### Структура JWT:
```
header.payload.signature
```

### Пример JWT:
```go
import "github.com/golang-jwt/jwt/v5"

// Секретный ключ (должен храниться в безопасном месте)
var jwtSecret = []byte("your-secret-key")

// Claims (данные в токене)
type Claims struct {
    UserID int `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

// Создание JWT токена
func GenerateToken(userID int, email string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    
    claims := &Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "your-app",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// Валидация JWT токена
func ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    
    return claims, nil
}
```

---

## 3. Хеширование паролей

### Использование bcrypt:
```go
import "golang.org/x/crypto/bcrypt"

// Хеширование пароля
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// Проверка пароля
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### Пример использования:
```go
func main() {
    password := "mysecretpassword"
    
    // Хешируем пароль
    hash, err := HashPassword(password)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Password:", password)
    fmt.Println("Hash:    ", hash)
    
    // Проверяем пароль
    match := CheckPasswordHash(password, hash)
    fmt.Println("Password match:", match)
    
    // Неверный пароль
    wrongMatch := CheckPasswordHash("wrongpassword", hash)
    fmt.Println("Wrong password match:", wrongMatch)
}
```

---

## 4. Middleware для аутентификации

### Middleware для проверки JWT:
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Получаем токен из заголовка Authorization
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }
        
        // Формат: "Bearer <token>"
        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
            http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
            return
        }
        
        tokenString := bearerToken[1]
        
        // Валидируем токен
        claims, err := ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        // Добавляем claims в контекст запроса
        ctx := context.WithValue(r.Context(), "userClaims", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Вспомогательная функция для получения claims из контекста
func GetUserFromContext(r *http.Request) *Claims {
    if claims, ok := r.Context().Value("userClaims").(*Claims); ok {
        return claims
    }
    return nil
}
```

---

## 🎯 Практика 1: Полная система аутентификации

**Задача:** Создать систему регистрации и входа с JWT

```go
package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
    _ "github.com/lib/pq"
)

type User struct {
    ID        int       `json:"id"`
    Email     string    `json:"email"`
    Password  string    `json:"password,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}

type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

type AuthService struct {
    db        *sql.DB
    jwtSecret []byte
}

func NewAuthService(db *sql.DB, jwtSecret string) *AuthService {
    return &AuthService{
        db:        db,
        jwtSecret: []byte(jwtSecret),
    }
}

func (s *AuthService) Init() error {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
    
    _, err := s.db.Exec(query)
    return err
}

func (s *AuthService) Register(email, password string) (*User, error) {
    // Проверяем, существует ли пользователь
    var existingID int
    err := s.db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&existingID)
    if err != nil && err != sql.ErrNoRows {
        return nil, fmt.Errorf("database error: %w", err)
    }
    if err == nil {
        return nil, fmt.Errorf("user already exists")
    }
    
    // Хешируем пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    // Создаем пользователя
    var user User
    query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`
    err = s.db.QueryRow(query, email, string(hashedPassword)).Scan(&user.ID, &user.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    user.Email = email
    return &user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
    var user User
    var passwordHash string
    
    // Ищем пользователя
    err := s.db.QueryRow(
        "SELECT id, email, password_hash FROM users WHERE email = $1",
        email,
    ).Scan(&user.ID, &user.Email, &passwordHash)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("invalid credentials")
        }
        return "", fmt.Errorf("database error: %w", err)
    }
    
    // Проверяем пароль
    err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
    if err != nil {
        return "", fmt.Errorf("invalid credentials")
    }
    
    // Генерируем JWT токен
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "auth-service",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return s.jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    
    return claims, nil
}

func (s *AuthService) GetUserByID(id int) (*User, error) {
    var user User
    err := s.db.QueryRow(
        "SELECT id, email, created_at FROM users WHERE id = $1",
        id,
    ).Scan(&user.ID, &user.Email, &user.CreatedAt)
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

type AuthHandler struct {
    service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
    return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    if req.Email == "" || req.Password == "" {
        http.Error(w, "Email and password are required", http.StatusBadRequest)
        return
    }
    
    if len(req.Password) < 6 {
        http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
        return
    }
    
    user, err := h.service.Register(req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    token, err := h.service.Login(req.Email, req.Password)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }
    
    response := map[string]interface{}{
        "token": token,
        "type":  "Bearer",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
    claims := GetUserFromContext(r)
    if claims == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    user, err := h.service.GetUserByID(claims.UserID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// AuthMiddleware для проверки JWT токена
func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }
        
        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
            http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
            return
        }
        
        tokenString := bearerToken[1]
        claims, err := h.service.ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        // Добавляем claims в контекст
        ctx := context.WithValue(r.Context(), "userClaims", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Вспомогательная функция для получения claims из контекста
func GetUserFromContext(r *http.Request) *Claims {
    if claims, ok := r.Context().Value("userClaims").(*Claims); ok {
        return claims
    }
    return nil
}

func connectDB() (*sql.DB, error) {
    // Для демонстрации используем PostgreSQL
    connStr := "host=localhost port=5432 user=authuser dbname=authdb sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    if err = db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}

func main() {
    db, err := connectDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Инициализируем сервис аутентификации
    authService := NewAuthService(db, "your-super-secret-jwt-key")
    if err := authService.Init(); err != nil {
        log.Fatal("Failed to init database:", err)
    }
    
    authHandler := NewAuthHandler(authService)
    
    r := mux.NewRouter()
    
    // Public routes
    r.HandleFunc("/api/register", authHandler.Register).Methods("POST")
    r.HandleFunc("/api/login", authHandler.Login).Methods("POST")
    
    // Protected routes
    protected := r.PathPrefix("/api").Subrouter()
    protected.Use(authHandler.AuthMiddleware)
    protected.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")
    
    // Health check
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
    }).Methods("GET")
    
    port := ":8080"
    log.Printf("🚀 Authentication API started on http://localhost%s", port)
    log.Fatal(http.ListenAndServe(port, r))
}
```

---

## 🎯 Практика 2: Система с ролями и разрешениями

**Задача:** Расширить систему аутентификации с ролевой моделью

```go
package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
    _ "github.com/lib/pq"
)

// Роли пользователей
const (
    RoleAdmin = "admin"
    RoleUser  = "user"
    RoleGuest = "guest"
)

type User struct {
    ID        int       `json:"id"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}

type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

type AuthService struct {
    db        *sql.DB
    jwtSecret []byte
}

func NewAuthService(db *sql.DB, jwtSecret string) *AuthService {
    return &AuthService{
        db:        db,
        jwtSecret: []byte(jwtSecret),
    }
}

func (s *AuthService) Init() error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            role VARCHAR(50) DEFAULT 'user',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
        
        `INSERT INTO users (email, password_hash, role) VALUES 
            ('admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin')
        ON CONFLICT (email) DO NOTHING`,
    }
    
    for _, query := range queries {
        _, err := s.db.Exec(query)
        if err != nil {
            return err
        }
    }
    
    return nil
}

func (s *AuthService) Register(email, password, role string) (*User, error) {
    if role == "" {
        role = RoleUser
    }
    
    // Проверяем допустимость роли
    if role != RoleAdmin && role != RoleUser && role != RoleGuest {
        return nil, fmt.Errorf("invalid role")
    }
    
    // Проверяем существование пользователя
    var existingID int
    err := s.db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&existingID)
    if err != nil && err != sql.ErrNoRows {
        return nil, fmt.Errorf("database error: %w", err)
    }
    if err == nil {
        return nil, fmt.Errorf("user already exists")
    }
    
    // Хешируем пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    // Создаем пользователя
    var user User
    query := `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) 
              RETURNING id, role, created_at`
    err = s.db.QueryRow(query, email, string(hashedPassword), role).Scan(&user.ID, &user.Role, &user.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    user.Email = email
    return &user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
    var user User
    var passwordHash string
    
    err := s.db.QueryRow(
        "SELECT id, email, password_hash, role FROM users WHERE email = $1",
        email,
    ).Scan(&user.ID, &user.Email, &passwordHash, &user.Role)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("invalid credentials")
        }
        return "", fmt.Errorf("database error: %w", err)
    }
    
    // Проверяем пароль
    err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
    if err != nil {
        return "", fmt.Errorf("invalid credentials")
    }
    
    // Генерируем JWT токен
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "auth-service",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return s.jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    
    return claims, nil
}

// Middleware для проверки ролей
func RequireRole(allowedRoles ...string) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims := GetUserFromContext(r)
            if claims == nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            // Проверяем, есть ли роль пользователя в разрешенных ролях
            hasAccess := false
            for _, role := range allowedRoles {
                if claims.Role == role {
                    hasAccess = true
                    break
                }
            }
            
            if !hasAccess {
                http.Error(w, "Insufficient permissions", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

type AdminHandler struct {
    service *AuthService
}

func NewAdminHandler(service *AuthService) *AdminHandler {
    return &AdminHandler{service: service}
}

func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := h.service.db.Query(`
        SELECT id, email, role, created_at 
        FROM users 
        ORDER BY created_at DESC
    `)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        users = append(users, user)
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    var req struct {
        Role string `json:"role"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Проверяем допустимость роли
    if req.Role != RoleAdmin && req.Role != RoleUser && req.Role != RoleGuest {
        http.Error(w, "Invalid role", http.StatusBadRequest)
        return
    }
    
    result, err := h.service.db.Exec(
        "UPDATE users SET role = $1 WHERE id = $2",
        req.Role, userID,
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    if rowsAffected == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "User role updated"})
}

func main() {
    db, err := connectDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    authService := NewAuthService(db, "your-super-secret-jwt-key")
    if err := authService.Init(); err != nil {
        log.Fatal("Failed to init database:", err)
    }
    
    authHandler := NewAuthHandler(authService)
    adminHandler := NewAdminHandler(authService)
    
    r := mux.NewRouter()
    
    // Public routes
    r.HandleFunc("/api/register", authHandler.Register).Methods("POST")
    r.HandleFunc("/api/login", authHandler.Login).Methods("POST")
    
    // User routes (требуют аутентификации)
    userRouter := r.PathPrefix("/api/user").Subrouter()
    userRouter.Use(authHandler.AuthMiddleware)
    userRouter.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")
    
    // Admin routes (требуют роль admin)
    adminRouter := r.PathPrefix("/api/admin").Subrouter()
    adminRouter.Use(authHandler.AuthMiddleware)
    adminRouter.Use(RequireRole(RoleAdmin))
    adminRouter.HandleFunc("/users", adminHandler.GetAllUsers).Methods("GET")
    adminRouter.HandleFunc("/users/{id}/role", adminHandler.UpdateUserRole).Methods("PUT")
    
    port := ":8080"
    log.Printf("🚀 Role-Based Authentication API started on http://localhost%s", port)
    log.Fatal(http.ListenAndServe(port, r))
}
```

---

## ❓ Важные моменты

### Безопасность JWT:
```go
// Никогда не храните секретный ключ в коде!
// Используйте переменные окружения:
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
    log.Fatal("JWT_SECRET environment variable is required")
}
```

### Защита от timing attacks:
```go
// Используйте constant time сравнение для чувствительных данных
if subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) != 1 {
    return false
}
```

---

## 🏠 Домашнее задание

**Задача 1: OAuth2 интеграция**
Добавьте поддержку OAuth2 провайдеров (Google/GitHub) в вашу систему аутентификации

**Задача 2: Система разрешений**
Создайте систему с детальными разрешениями:
- Пользователи могут иметь multiple roles
- Разрешения на уровне отдельных операций
- Динамическое управление правами

**Задача 3: Безопасность API**
Реализуйте дополнительные меры безопасности:
- Rate limiting
- IP whitelisting
- Audit log для чувствительных операций

---

## 🚀 Что ждет на следующем занятии?

*   **Продвинутая работа с БД:** Миграции, оптимизация запросов
*   **Контекст в Go:** Отмена операций, таймауты
*   **Подготовленные запросы:** Повышение производительности

**Удачи в создании безопасных приложений! 🎉**