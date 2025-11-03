Занятие 13: Модули и управление зависимостями

"Создаем профессиональные проекты с Go Modules"

---

📝 План на сегодня

1. Что такое Go Modules? Эволюция управления зависимостями в Go
2. Создание модуля: go mod init
3. Работа с зависимостями: go get, go mod tidy
4. Версионирование: Семантическое версионирование
5. Публикация модулей: Создание собственных библиотек
6. Практика: Создаем реальный проект с внешними зависимостями

---

1. Что такое Go Modules?

Go Modules - система управления зависимостями, встроенная в Go начиная с версии 1.11.

Преимущества Go Modules:

· Версионирование зависимостей
· Воспроизводимость сборок
· Изоляция проектов
· Поддержка семантического версионирования

Ключевые файлы:

· go.mod - описание модуля и зависимостей
· go.sum - контрольные суммы для безопасности

---

2. Создание модуля

Инициализация нового модуля:

# Создаем директорию проекта
mkdir myproject
cd myproject

# Инициализируем модуль
go mod init github.com/username/myproject
Содержимое go.mod после инициализации:

```
module github.com/username/myproject

go 1.21
```

Структура проекта:


myproject/
├── go.mod
├── go.sum
├── main.go
├── pkg/
│   └── utils/
│       └── strings.go
└── internal/
    └── config/
        └── config.go

---

3. Работа с зависимостями

Добавление зависимостей:

### Добавляем конкретную версию
go get github.com/gin-gonic/gin@v1.9.1

### Добавляем последнюю версию
go get github.com/gorilla/mux

### Добавляем все зависимости из кода
go mod tidy
Пример go.mod с зависимостями:
```
module github.com/username/myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/gorilla/mux v1.8.0
)

require (
    github.com/bytedance/sonic v1.9.1 // indirect
    github.com/golang/protobuf v1.5.2 // indirect
    // ... другие indirect зависимости
)
```

---

## 4. Версионирование

Семантическое версионирование (SemVer):

· v1.2.3
· v0.1.0 (нестабильная версия)
· v2.0.0 (обратно несовместимые изменения)

Специальные версии:

### Последняя версия
go get package@latest

### Конкретная версия
go get package@v1.2.3

### Конкретный коммит
go get package@a1b2c3d

### Ветка
go get package@master

---

# 5. Публикация модулей

Подготовка модуля к публикации:

// go.mod
module github.com/username/my-utility-library

go 1.21

// Код библиотеки
package mylibrary

import "fmt"

// Greet возвращает приветствие
func Greet(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}
Тегирование версии:

git tag v1.0.0
git push origin v1.0.0
---

🎯 Практика 1: Создание проекта с веб-фреймворком

Задача: Создать веб-сервер с использованием внешних зависимостей

# Создаем проект
mkdir web-server
cd web-server
go mod init github.com/username/web-server
```go
// main.go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    
    "github.com/gorilla/mux"
)

type Config struct {
    Port string json:"port"
}

type User struct {
    ID    int    json:"id"
    Name  string json:"name"
    Email string json:"email"
}

var users = []User{
    {ID: 1, Name: "Alice", Email: "alice@example.com"},
    {ID: 2, Name: "Bob", Email: "bob@example.com"},
    {ID: 3, Name: "Charlie", Email: "charlie@example.com"},
}

func loadConfig() (*Config, error) {
    data, err := os.ReadFile("config.json")
    if err != nil {
        // Конфиг по умолчанию
        return &Config{Port: "8080"}, nil
    }
    
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    return &config, nil
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {

vars := mux.Vars(r)
    id := vars["id"]
    
    // В реальном приложении здесь был бы поиск в БД
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "User ID: " + id,
        "status":  "success",
    })
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
    })
}

func main() {
    config, err := loadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    r := mux.NewRouter()
    
    // Маршруты
    r.HandleFunc("/users", getUsersHandler).Methods("GET")
    r.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
    r.HandleFunc("/health", healthHandler).Methods("GET")
    
    // Статические файлы
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
    
    log.Printf("Server starting on port %s", config.Port)
    log.Fatal(http.ListenAndServe(":"+config.Port, r))
}
```

```json
// config.json
{
    "port": "3000"
}
```

bash
# Добавляем зависимости
go mod tidy

# Запускаем сервер
go run main.go

---

🎯 Практика 2: Создание утилитарной библиотеки

Задача: Создать библиотеку для работы с строками и опубликовать ее структуру

bash
# Создаем библиотеку
mkdir string-utils
cd string-utils
go mod init github.com/username/string-utils


```go
// stringutils.go
package stringutils

import (
    "strings"
    "unicode"
)

// Reverse возвращает перевернутую строку
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// CountWords подсчитывает количество слов в строке
func CountWords(s string) int {
    words := strings.Fields(s)
    return len(words)
}

// IsPalindrome проверяет, является ли строка палиндромом
func IsPalindrome(s string) bool {
    cleaned := strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) || unicode.IsPunct(r) {
            return -1
        }
        return unicode.ToLower(r)
    }, s)
    
    return cleaned == Reverse(cleaned)
}

// Truncate обрезает строку до указанной длины
func Truncate(s string, maxLength int) string {
    if len(s) <= maxLength {
        return s
    }
    
    runes := []rune(s)
    if len(runes) <= maxLength {
        return s
    }
    
    return string(runes[:maxLength]) + "..."
}
```

```go
// stringutils_test.go
package stringutils

import "testing"

func TestReverse(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"hello", "olleh"},
        {"", ""},
        {"a", "a"},
        {"привет", "тевирп"},
    }
    
    for _, test := range tests {
        result := Reverse(test.input)
        if result != test.expected {
            t.Errorf("Reverse(%q) = %q, expected %q", test.input, result, test.expected)
        }
    }
}

func TestCountWords(t *testing.T) {
    tests := []struct {
        input    string
        expected int
    }{
        {"hello world", 2},
        {"", 0},
        {"   multiple   spaces   ", 2},
        {"one", 1},
    }
    
    for _, test := range tests {
        result := CountWords(test.input)
        if result != test.expected {
            t.Errorf("CountWords(%q) = %d, expected %d", test.input, result, test.expected)
        }
    }
}
```

```go
// go.mod (финальный вид)
module github.com/username/string-utils

go 1.21
```

---

🎯 Практика 3: Проект с несколькими зависимостями

Задача: Создать CLI утилиту для работы с HTTP API

bash
# Создаем проект
mkdir api-cli
cd api-cli
go mod init github.com/username/api-cli


```go
// main.go
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/fatih/color"
    "github.com/urfave/cli/v2"
)

type Post struct {
    UserID int    json:"userId"
    ID     int    json:"id"
    Title  string json:"title"
    Body   string json:"body"
}

type APIClient struct {
    baseURL string
    client  *http.Client
}

func NewAPIClient() *APIClient {
    return &APIClient{
        baseURL: "https://jsonplaceholder.typicode.com",
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *APIClient) GetPosts() ([]Post, error) {
    resp, err := c.client.Get(c.baseURL + "/posts")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var posts []Post
    if err := json.Unmarshal(body, &posts); err != nil {
        return nil, err
    }
    
    return posts, nil
}

func (c *APIClient) GetPost(id int) (*Post, error) {
    resp, err := c.client.Get(fmt.Sprintf("%s/posts/%d", c.baseURL, id))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var post Post
    if err := json.Unmarshal(body, &post); err != nil {
        return nil, err
    }
    
    return &post, nil
}

func displayPosts(posts []Post, limit int) {
    green := color.New(color.FgGreen).SprintFunc()
    blue := color.New(color.FgBlue).SprintFunc()
    
    if limit > 0 && limit < len(posts) {
        posts = posts[:limit]
    }
    
    for _, post := range posts {
        fmt.Printf("%s %s\n", green("📝"), blue(post.Title))
        fmt.Printf("   User ID: %d, Post ID: %d\n", post.UserID, post.ID)
        fmt.Printf("   %s\n\n", post.Body[:50]+"...")
    }
}

func main() {
    var limit int
    var postID int
    
    flag.IntVar(&limit, "limit", 5, "Limit number of posts to display")
    flag.IntVar(&postID, "post", 0, "Get specific post by ID")
    flag.Parse()
    
    client := NewAPIClient()
    
    if postID > 0 {
        post, err := client.GetPost(postID)
        if err != nil {
            log.Fatalf("Error fetching post: %v", err)
        }
        
        yellow := color.New(color.FgYellow).SprintFunc()
        cyan := color.New(color.FgCyan).SprintFunc()
        
        fmt.Printf("%s %s\n", yellow("📄"), cyan(post.Title))
        fmt.Printf("User ID: %d\n", post.UserID)
        fmt.Printf("Post ID: %d\n", post.ID)
        fmt.Printf("\n%s\n", post.Body)
    } else {
        posts, err := client.GetPosts()
        if err != nil {
            log.Fatalf("Error fetching posts: %v", err)
        }
        
        fmt.Printf("📬 Found %d posts\n\n", len(posts))
        displayPosts(posts, limit)
    }
}
```

bash
# Добавляем зависимости
go get github.com/fatih/color
go get github.com/urfave/cli/v2

# Собираем утилиту
go build -o api-cli

# Запускаем
./api-cli -limit 3
./api-cli -post 1

---

🎯 Практика 4: Работа с версиями и обновлениями

Задача: Управление версиями зависимостей в проекте

bash
# Просмотр текущих зависимостей
go list -m all

# Просмотр доступных версий пакета
go list -m -versions github.com/gorilla/mux

# Обновление до последней версии
go get -u github.com/gorilla/mux

# Обновление до конкретной версии
go get github.com/gorilla/mux@v1.8.0

# Обновление всех зависимостей
go get -u ./...

# Удаление неиспользуемых зависимостей
go mod tidy

# Проверка зависимостей
go mod verify

# Загрузка зависимостей в vendor
go mod vendor

Пример управления версиями:

```go
// go.mod с конкретными версиями
module my-project

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/gorilla/mux v1.8.0
    github.com/lib/pq v1.10.9
)

// replace для локальной разработки
replace github.com/username/my-library => ../my-library
```

---

❓ Важные моменты

Indirect зависимости:

· Прямые - импортируются напрямую в ваш код
· Indirect - зависимости ваших зависимостей

Миграция с GOPATH:

# Для старых проектов
go mod init
go mod tidy
Работа в приватных репозиториях:

# Установка git конфигурации для приватных репозиториев
git config --global url."https://github.com/".insteadOf "https://github.com/"
---

🏠 Домашнее задание

Задача 1: Создание модуля для работы с датами
Создайте модуль с функциями:

· Форматирование дат
· Расчет разницы между датами
· Проверка рабочих дней
  Опубликуйте структуру на GitHub

Задача 2: Миграция существующего проекта
Возьмите один из ваших предыдущих проектов и:

· Инициализируйте go.mod
· Добавьте необходимые зависимости
· Настройте правильную структуру проекта

Задача 3: CLI утилита с внешними зависимостями
Создайте утилиту командной строки которая:

· Использует 3 внешние зависимости
· Имеет несколько команд (flags)
· Сохраняет конфигурацию в файл

---

🚀 Что ждет на следующем занятии?

· Введение в веб-разработку: HTTP серверы в Go
· Роутинг и обработчики: Создание REST API
· Пакет net/http: Основы веб-программирования

Удачи в создании модулей! 🎉