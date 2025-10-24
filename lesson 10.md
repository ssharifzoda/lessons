# Занятие 10: Обработка ошибок в Go

`"Предсказуемое поведение программы при ошибках"`

---

## 📝 План на сегодня

1.  **Философия обработки ошибок в Go:** Явное лучше неявного
2.  **Тип error:** Стандартный интерфейс для ошибок
3.  **Создание ошибок:** errors.New, fmt.Errorf, кастомные ошибки
4.  **Проверка ошибок:** Идиоматический подход Go
5.  **Panic и recover:** Экстренные ситуации
6.  **Практика:** Обработка ошибок в реальных сценариях

---

## 1. Философия обработки ошибок в Go

**В Go ошибки - это значения**, которые возвращаются явно.

### Сравнение с другими языками:
```go
// Go (явная обработка)
result, err := doSomething()
if err != nil {
    // обрабатываем ошибку
}

// vs Исключения (неявные)
try {
    result = doSomething()
} catch (Exception e) {
    // обработка
}
```

### Преимущества подхода Go:
- **Ясность** - видно, какие функции могут вернуть ошибку
- **Контроль** - программист решает, как обработать каждую ошибку
- **Производительность** - нет накладных расходов на исключения

---

## 2. Тип error

**error** - это встроенный интерфейс:
```go
type error interface {
    Error() string
}
```

Любой тип, у которого есть метод `Error() string`, реализует интерфейс error.

### Стандартное использование:
```go
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("деление на ноль")
    }
    return a / b, nil
}

func main() {
    result, err := Divide(10, 0)
    if err != nil {
        fmt.Println("Ошибка:", err)
        return
    }
    fmt.Println("Результат:", result)
}
```

---

## 3. Создание ошибок

### Способ 1: errors.New()
```go
import "errors"

func ValidateAge(age int) error {
    if age < 0 {
        return errors.New("возраст не может быть отрицательным")
    }
    if age > 150 {
        return errors.New("возраст не может быть больше 150")
    }
    return nil
}
```

### Способ 2: fmt.Errorf() (чаще используется)
```go
func ProcessUser(name string, age int) error {
    if name == "" {
        return fmt.Errorf("имя не может быть пустым")
    }
    if age < 18 {
        return fmt.Errorf("пользователь %s несовершеннолетний: %d лет", name, age)
    }
    return nil
}
```

### Способ 3: Кастомные ошибки
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("ошибка валидации поля '%s': %s", e.Field, e.Message)
}

func ValidateUser(name string, age int) error {
    if name == "" {
        return ValidationError{Field: "name", Message: "не может быть пустым"}
    }
    if age < 0 {
        return ValidationError{Field: "age", Message: "не может быть отрицательным"}
    }
    return nil
}
```

---

## 4. Проверка ошибок

### Базовая проверка:
```go
file, err := os.Open("config.txt")
if err != nil {
    fmt.Println("Не удалось открыть файл:", err)
    return
}
defer file.Close() // defer гарантирует выполнение при выходе из функции
```

### Проверка конкретных типов ошибок:
```go
func handleError(err error) {
    // Проверяем, является ли ошибка кастомной ValidationError
    if ve, ok := err.(ValidationError); ok {
        fmt.Printf("Ошибка валидации: %s - %s\n", ve.Field, ve.Message)
        return
    }
    
    // Проверяем стандартные ошибки
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("Файл не существует")
        return
    }
    
    // Общий случай
    fmt.Println("Неизвестная ошибка:", err)
}
```

### errors.Is и errors.As:
```go
// errors.Is - проверяет, является ли ошибка конкретным типом
if errors.Is(err, sql.ErrNoRows) {
    fmt.Println("Запись не найдена в БД")
}

// errors.As - извлекает конкретный тип ошибки
var valErr ValidationError
if errors.As(err, &valErr) {
    fmt.Printf("Ошибка в поле %s: %s\n", valErr.Field, valErr.Message)
}
```

---

## 5. Panic и recover

### Panic - экстренная остановка программы:
```go
func criticalOperation() {
    panic("критическая ошибка: данные повреждены")
}

// Лучше так:
func safeCriticalOperation() error {
    return fmt.Errorf("данные повреждены")
}
```

### Recover - перехват panic (используется в defer):
```go
func safeExecute(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("перехвачена panic: %v", r)
        }
    }()
    
    fn()
    return nil
}

func main() {
    err := safeExecute(func() {
        panic("что-то пошло не так!")
    })
    
    if err != nil {
        fmt.Println("Ошибка:", err) // Ошибка: перехвачена panic: что-то пошло не так!
    }
}
```

---

## 🎯 Практика 1: Валидация пользовательских данных

**Задача:** Создать систему валидации с детальными ошибками

```go
package main
import (
    "errors"
    "fmt"
    "strings"
)

// Кастомные ошибки
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("валидация %s: %s (значение: %v)", e.Field, e.Message, e.Value)
}

type User struct {
    Username string
    Email    string
    Age      int
}

// Функции валидации
func validateUsername(username string) error {
    if len(username) < 3 {
        return ValidationError{
            Field:   "username",
            Value:   username,
            Message: "должен содержать至少 3 символа",
        }
    }
    if strings.Contains(username, " ") {
        return ValidationError{
            Field:   "username", 
            Value:   username,
            Message: "не может содержать пробелы",
        }
    }
    return nil
}

func validateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return ValidationError{
            Field:   "email",
            Value:   email,
            Message: "должен содержать @",
        }
    }
    return nil
}

func validateAge(age int) error {
    if age < 0 {
        return ValidationError{
            Field:   "age",
            Value:   age,
            Message: "не может быть отрицательным",
        }
    }
    if age < 18 {
        return ValidationError{
            Field:   "age",
            Value:   age, 
            Message: "должен быть至少 18 лет",
        }
    }
    return nil
}

func ValidateUser(user User) error {
    var errs []error
    
    if err := validateUsername(user.Username); err != nil {
        errs = append(errs, err)
    }
    
    if err := validateEmail(user.Email); err != nil {
        errs = append(errs, err)
    }
    
    if err := validateAge(user.Age); err != nil {
        errs = append(errs, err)
    }
    
    if len(errs) > 0 {
        return fmt.Errorf("ошибки валидации: %v", errs)
    }
    
    return nil
}

func main() {
    testUsers := []User{
        {"al", "invalid-email", 15},
        {"validuser", "valid@email.com", 25},
        {"user with space", "test@example.com", 20},
    }
    
    for _, user := range testUsers {
        fmt.Printf("Валидация пользователя: %+v\n", user)
        
        if err := ValidateUser(user); err != nil {
            fmt.Printf("❌ Ошибки: %v\n\n", err)
        } else {
            fmt.Printf("✅ Валидация пройдена успешно\n\n")
        }
    }
}
```

---

## 🎯 Практика 2: Файловые операции с обработкой ошибок

**Задача:** Безопасная работа с файлами

```go
package main
import (
    "fmt"
    "io"
    "os"
    "path/filepath"
)

// SafeFileOperations оборачивает файловые операции с обработкой ошибок
type SafeFileOperations struct{}

func (s SafeFileOperations) ReadFile(filename string) (string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return "", fmt.Errorf("не удалось открыть файл %s: %w", filename, err)
    }
    defer file.Close()
    
    content, err := io.ReadAll(file)
    if err != nil {
        return "", fmt.Errorf("не удалось прочитать файл %s: %w", filename, err)
    }
    
    return string(content), nil
}

func (s SafeFileOperations) WriteFile(filename, content string) error {
    // Создаем директорию, если её нет
    dir := filepath.Dir(filename)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("не удалось создать директорию %s: %w", dir, err)
    }
    
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("не удалось создать файл %s: %w", filename, err)
    }
    defer file.Close()
    
    _, err = file.WriteString(content)
    if err != nil {
        return fmt.Errorf("не удалось записать в файл %s: %w", filename, err)
    }
    
    return nil
}

func (s SafeFileOperations) CopyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("не удалось открыть исходный файл %s: %w", src, err)
    }
    defer sourceFile.Close()
    
    // Создаем целевой файл
    if err := s.WriteFile(dst, ""); err != nil {
        return err
    }
    
    destFile, err := os.OpenFile(dst, os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("не удалось открыть целевой файл %s: %w", dst, err)
    }
    defer destFile.Close()
    
    _, err = io.Copy(destFile, sourceFile)
    if err != nil {
        return fmt.Errorf("ошибка копирования из %s в %s: %w", src, dst, err)
    }
    
    return nil
}

func main() {
    fileOps := SafeFileOperations{}
    
    // Тестируем операции
    operations := []struct {
        name string
        fn   func() error
    }{
        {
            name: "Запись файла",
            fn: func() error {
                return fileOps.WriteFile("test/data.txt", "Привет, мир!")
            },
        },
        {
            name: "Чтение файла", 
            fn: func() error {
                content, err := fileOps.ReadFile("test/data.txt")
                if err != nil {
                    return err
                }
                fmt.Printf("Прочитано: %s\n", content)
                return nil
            },
        },
        {
            name: "Копирование файла",
            fn: func() error {
                return fileOps.CopyFile("test/data.txt", "test/backup.txt")
            },
        },
    }
    
    for _, op := range operations {
        fmt.Printf("Выполняем: %s\n", op.name)
        if err := op.fn(); err != nil {
            fmt.Printf("❌ Ошибка: %v\n\n", err)
        } else {
            fmt.Printf("✅ Успешно\n\n")
        }
    }
}
```

---

## 🎯 Практика 3: API клиент с обработкой ошибок

**Задача:** Создать HTTP клиент с гранулярной обработкой ошибок

```go
package main
import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "time"
)

// Кастомные ошибки для API
var (
    ErrNetwork         = errors.New("сетевая ошибка")
    ErrTimeout         = errors.New("превышено время ожидания")
    ErrInvalidResponse = errors.New("неверный ответ от сервера")
    ErrNotFound        = errors.New("ресурс не найден")
    ErrServerError     = errors.New("ошибка сервера")
)

type APIClient struct {
    baseURL    string
    httpClient *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
    return &APIClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Username string `json:"username"`
}

func (c *APIClient) GetUser(userID int) (*User, error) {
    url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)
    
    resp, err := c.httpClient.Get(url)
    if err != nil {
        // Анализируем тип сетевой ошибки
        if errors.Is(err, http.ErrHandlerTimeout) {
            return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
        }
        return nil, fmt.Errorf("%w: %v", ErrNetwork, err)
    }
    defer resp.Body.Close()
    
    // Обрабатываем HTTP статусы
    switch resp.StatusCode {
    case http.StatusOK:
        // Продолжаем обработку
    case http.StatusNotFound:
        return nil, ErrNotFound
    case http.StatusInternalServerError:
        return nil, ErrServerError
    default:
        return nil, fmt.Errorf("неожиданный статус: %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
    }
    
    var user User
    if err := json.Unmarshal(body, &user); err != nil {
        return nil, fmt.Errorf("%w: %v", ErrInvalidResponse, err)
    }
    
    return &user, nil
}

// Обработчик ошибок
func handleAPIError(err error) {
    fmt.Printf("Ошибка API: ")
    
    switch {
    case errors.Is(err, ErrNetwork):
        fmt.Println("Проверьте подключение к интернету")
    case errors.Is(err, ErrTimeout):
        fmt.Println("Сервер не отвечает, попробуйте позже")
    case errors.Is(err, ErrNotFound):
        fmt.Println("Запрошенный ресурс не найден")
    case errors.Is(err, ErrServerError):
        fmt.Println("Проблемы на сервере, попробуйте позже")
    case errors.Is(err, ErrInvalidResponse):
        fmt.Println("Сервер вернул неверные данные")
    default:
        fmt.Printf("Неизвестная ошибка: %v\n", err)
    }
}

func main() {
    client := NewAPIClient("https://jsonplaceholder.typicode.com")
    
    // Тестируем разные сценарии
    testCases := []struct {
        name   string
        userID int
    }{
        {"Валидный пользователь", 1},
        {"Несуществующий пользователь", 9999},
        // Для теста сетевых ошибок можно использовать несуществующий URL
    }
    
    for _, tc := range testCases {
        fmt.Printf("Тест: %s\n", tc.name)
        
        user, err := client.GetUser(tc.userID)
        if err != nil {
            handleAPIError(err)
        } else {
            fmt.Printf("✅ Пользователь: %+v\n", user)
        }
        fmt.Println()
    }
}
```

---

## ❓ Важные моменты

### Не игнорируйте ошибки!
```go
// ПЛОХО:
file, _ := os.Open("file.txt") // Игнорируем ошибку!

// ХОРОШО:
file, err := os.Open("file.txt")
if err != nil {
    return fmt.Errorf("не удалось открыть файл: %w", err)
}
```

### Используйте %w для wrapping ошибок:
```go
if err != nil {
    return fmt.Errorf("контекст: %w", err) // Сохраняет оригинальную ошибку
}
```

---

## 🏠 Домашнее задание

**Задача 1: Калькулятор с полной обработкой ошибок**
Создайте калькулятор, который обрабатывает:
- Деление на ноль
- Неверные операторы
- Некорректный ввод чисел
- Переполнение чисел

**Задача 2: Валидатор конфигурации**
Напишите функцию, которая проверяет конфигурацию приложения:
- Обязательные поля
- Корректные диапазоны значений
- Существование файлов/директорий
  Возвращайте детальные ошибки для каждой проблемы.

**Задача 3: Retry механизм**
Создайте функцию с повторными попытками:
```go
func WithRetry(fn func() error, maxAttempts int) error
```
Она должна повторять операцию при временных ошибках.

---

## 🚀 Что ждет на следующем занятии?

*   **Работа с файлами:** Чтение, запись, манипуляции
*   **Пакеты os и io:** Фундаментальные операции
*   **Работа с JSON и CSV**

**Удачи в решении задач! 🎉**