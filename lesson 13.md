# Занятие 13: Модули и управление зависимостями

Создаем профессиональные проекты с Go Modules

---

📝 План на сегодня

1. Что такое Go Modules? Эволюция управления зависимостями в Go
2. Создание модуля: go mod init
3. Работа с зависимостями: go get, go mod tidy
4. Private и Public функции в Go
5. Практика: Создаем реальный проект с внешними зависимостями

---

1. Что такое Go Modules?

Go Modules - система управления зависимостями, встроенная в Go начиная с версии 1.11.

Преимущества Go Modules:
- Управление версиями библиотек
- Изоляция зависимостей между проектами
- Воспроизводимость сборок (одни и те же зависимости на любой машине)

Ключевые файлы:
- go.mod - описание модуля и зависимостей
- go.sum - контрольные суммы для безопасности

---
### 2. Создание модуля

Инициализация нового модуля:

```bash
# Создаем директорию проекта
mkdir myproject
cd myproject

# Инициализируем модуль
go mod init github.com/username/myproject
```

Содержимое go.mod после инициализации:

```bash
module github.com/username/myproject
# или module lesson13

go 1.24.9
```

Структура проекта:
```
myproject/
├── go.mod
├── main.go
└── calculator/
    └── functions.go
```

```go
# main.go
package main

import (
"fmt"
"github.com/manuchehr0/calculator/calculator"
)

func main() {
    fmt.Println("2 + 3 =", calculator.Add(2, 3))
    fmt.Println("5 - 2 =", calculator.Subtract(5, 2))
    fmt.Println("3 * 4 =", calculator.Multiply(3, 4))
    fmt.Println("10 / 2 =", calculator.Divide(10, 2))
}
```

```go
# calculator/functions.go
package calculator

func Add(a, b int) int {
    return a + b
}

func Subtract(a, b int) int {
    return a - b
}

func Multiply(a, b int) int {
    return a * b
}

func Divide(a, b int) int {
if b == 0 {
    return 0
}
return a / b
}

```
---

## 3. Работа с зависимостями

### Добавление зависимостей:

```bash
go get github.com/fatih/color
```

#### Добавление все зависимости из кода
`go mod tidy`

Пример go.mod с зависимостями:
```
module github.com/username/calculator

go 1.24.9

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.25.0 // indirect
)
```

### 3. Использование внешных модулей
file `main.go`
```go
package main

import (
  "fmt"
  "github.com/fatih/color"
  "github.com/manuchehr0/calculator/calculator"
)


func main() {
  // Считаем и выводим результаты с цветом

  // 2 + 3
  addResult := calculator.Add(2, 3)
  fmt.Println("2 + 3 =", color.New(color.FgGreen).Sprint(addResult))

  // 5 - 2
  subResult := calculator.Subtract(5, 2)
  fmt.Println("5 - 2 =", color.New(color.FgRed).Sprint(subResult))

  // 3 * 4
  mulResult := calculator.Multiply(3, 4)
  fmt.Println("3 * 4 =", color.New(color.FgGreen).Sprint(mulResult))

  // 10 / 2
  divResult := calculator.Divide(10, 2)
  fmt.Println("10 / 2 =", color.New(color.FgRed).Sprint(divResult))
}
```

### 4. Private и Public функции в Go
#### Зачем это нужно?
- Безопасность: уменьшаем шанс случайного использования неправильной функции.
- Читаемость кода: другие разработчики видят только важные, публичные функции.


#### Экспортируемые функции (Public)
- Начинаются с большой буквы.
- Доступны из других пакетов.
- Используются, когда мы хотим, чтобы кто-то мог вызвать функцию из нашего пакета.

```go
package calculator

// Экспортируемая функция
func Add(a, b int) int {
    return a + b
}
```

#### Неэкспортируемые функции (Private)
- Начинаются с маленькой буквы.
- Доступны только внутри своего пакета.
- Используются, чтобы скрыть детали реализации от других пакетов.

`file: calculator/functions.go`
```go
package calculator

// Неэкспортируемая функция
func subtract(a, b int) int {
    return a - b
}

// Экспортируемая функция использует неэкспортируемую
func Subtract(a, b int) int {
    return subtract(a, b)
}
```


---

## 🎯 Практика 1: Создание утилитарной библиотеки

Задача: Создать библиотеку для работы с строками и опубликовать ее структуру

```bash
# Создаем библиотеку
mkdir string-utils
cd string-utils
go mod init github.com/username/string-utils
```

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

## 🎯 Практика 2: Проект с несколькими зависимостями

Задача: Создать CLI утилиту для работы с HTTP API

bash
### Создаем проект
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

```bash
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
```

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
· Имеет несколько команд 
· Сохраняет конфигурацию в файл

---

Удачи в создании модулей! 🎉