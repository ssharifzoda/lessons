# Занятие 18: Middleware и обработка ошибок

`"Создаем надежные API с middleware"`

---

## 📝 План на сегодня

1.  **Middleware паттерн:** Цепочка обработчиков
2.  **Стандартные middleware:** Логирование, CORS, аутентификация
3.  **Кастомные middleware:** Создание своих обработчиков
4.  **Централизованная обработка ошибок:** Единый формат ответов
5.  **Panic recovery:** Защита от падений

---

## 1. Middleware паттерн

```go
type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
    return func(next http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            next = middlewares[i](next)
        }
        return next
    }
}
```

---

## 2. Стандартные middleware

### Логирование
```go
func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        next.ServeHTTP(w, r)
        
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}
```

### CORS
```go
func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

---

## 3. Кастомные middleware

### Аутентификация
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token != "Bearer valid-token" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### Rate Limiting
```go
func RateLimit(requestsPerMinute int) Middleware {
    return func(next http.Handler) http.Handler {
        limiter := rate.NewLimiter(rate.Every(time.Minute), requestsPerMinute)
        
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Too many requests", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 4. Централизованная обработка ошибок

```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func ErrorHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic: %v", err)
                writeError(w, http.StatusInternalServerError, "Internal server error")
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

func writeError(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(AppError{
        Code:    code,
        Message: message,
    })
}
```

---

## 5. Panic recovery

```go
func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Recovered from panic: %v", err)
                writeError(w, http.StatusInternalServerError, "Something went wrong")
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}
```

---

## 🎯 Практика: Полный пример API

```go
func main() {
    r := mux.NewRouter()
    
    // API routes
    api := r.PathPrefix("/api").Subrouter()
    api.HandleFunc("/users", getUsers).Methods("GET")
    api.HandleFunc("/users", createUser).Methods("POST")
    
    // Apply middleware chain
    handler := Chain(
        Logging,
        CORS,
        Recovery,
        ErrorHandler,
    )(r)
    
    log.Println("Server started on :8080")
    http.ListenAndServe(":8080", handler)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    users := []User{{ID: 1, Name: "John"}}
    json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        writeError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }
    
    // Simulate error
    if user.Name == "" {
        writeError(w, http.StatusBadRequest, "Name is required")
        return
    }
    
    user.ID = 2
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

---

## ❓ Важные моменты

### Порядок middleware важен:
```go
// Правильный порядок:
Chain(
    Logging,    // 1. Логируем все запросы
    CORS,       // 2. Обрабатываем CORS
    Auth,       // 3. Проверяем аутентификацию
    Recovery,   // 4. Ловим паники
)
```

### Всегда устанавливайте Content-Type:
```go
w.Header().Set("Content-Type", "application/json")
```

---

## 🏠 Домашнее задание

**Задача 1:** Создайте middleware для проверки API ключа
**Задача 2:** Добавьте rate limiting на 100 запросов в минуту
**Задача 3:** Реализуйте middleware для сжатия ответов (gzip)

---

Следующее занятие: **Деплой приложений**