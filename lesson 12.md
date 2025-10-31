A Sino Humo:
Занятие 12: Тестирование в Go

"Пишем надежный код с помощью unit-тестов"

---

📝 План на сегодня

1. Зачем нужно тестирование? Философия тестирования в Go
2. Пакет testing: Основы написания тестов
3. Table-driven tests: Эффективный подход к тестированию
4. Тестирование HTTP handlers: Mocking и тестовые серверы
5. Покрытие кода (coverage): Анализ качества тестов
6. Практика: Пишем тесты для реального кода

---

1. Зачем нужно тестирование?

Тестирование - это проверка, что код работает правильно в различных сценариях.

Преимущества тестирования:

· Обнаружение багов на ранней стадии
· Уверенность при рефакторинге
· Документация поведения функций
· Качество кода

Философия Go:

· Тесты рядом с кодом (файлы _test.go)
· Минималистичный встроенный фреймворк
· Простота написания тестов

---

2. Пакет testing - основы

Структура тестового файла:

// calculator.go
package calculator

func Add(a, b int) int {
    return a + b
}

func Multiply(a, b int) int {
    return a * b
}
```go
// calculator_test.go
package calculator

import "testing"

// Тест функции Add
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    
    if result != expected {
        t.Errorf("Add(2, 3) = %d; expected %d", result, expected)
    }
}

// Тест функции Multiply
func TestMultiply(t *testing.T) {
    result := Multiply(4, 5)
    expected := 20
    
    if result != expected {
        t.Errorf("Multiply(4, 5) = %d; expected %d", result, expected)
    }
}

Запуск тестов:

bash
# Запуск всех тестов в пакете
go test

# Запуск с подробным выводом
go test -v

# Запуск конкретного теста
go test -v -run TestAdd

---

3. Table-driven tests

Table-driven tests - подход, когда тестовые данные выносятся в таблицу.

go
func TestAdd_TableDriven(t *testing.T) {
    tests := []struct {
        name     string
        a        int
        b        int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -1, -2},
        {"mixed numbers", -5, 10, 5},
        {"zero", 0, 0, 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; expected %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

---

4. Тестирование ошибок

Тестирование функций, возвращающих ошибки:

go
// math.go
package math

import "errors"

var ErrDivisionByZero = errors.New("division by zero")

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, ErrDivisionByZero
    }
    return a / b, nil
}


go
// math_test.go
package math

import "testing"

func TestDivide(t *testing.T) {
    tests := []struct {
        name        string
        a           float64
        b           float64
        expected    float64
        expectError bool
    }{
        {"normal division", 10.0, 2.0, 5.0, false},
        {"division by zero", 10.0, 0.0, 0.0, true},
        {"fraction result", 1.0, 2.0, 0.5, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Divide(tt.a, tt.b)
            
            if tt.expectError {
                if err == nil {
                    t.Error("expected error, but got nil")
                }
                if err != ErrDivisionByZero {
                    t.Errorf("expected ErrDivisionByZero, got %v", err)
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if result != tt.expected {
                    t.Errorf("Divide(%f, %f) = %f; expected %f", 
                        tt.a, tt.b, result, tt.expected)
                }
            }
        })
    }
}

---

5. Покрытие кода (Coverage)

Анализ покрытия тестами:

bash
# Запуск тестов с анализом покрытия
go test -cover

# Создание детального отчета покрытия
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
`

Пример отчета покрытия:

PASS
coverage: 85.7% of statements
ok      calculator 0.002s
---

🎯 Практика 1: Тестирование утилит для работы со строками

Задача: Протестировать функции для работы со строками

// string_utils.go
package utils

import "strings"

func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func CountVowels(s string) int {
    vowels := "aeiouаеёиоуыэюяAEIOUАЕЁИОУЫЭЮЯ"
    count := 0
    for _, char := range s {
        if strings.ContainsRune(vowels, char) {
            count++
        }
    }
    return count
}

func IsPalindrome(s string) bool {
    cleaned := strings.ToLower(strings.ReplaceAll(s, " ", ""))
    return cleaned == Reverse(cleaned)
}
```go
// string_utils_test.go
package utils

import "testing"

func TestReverse(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"empty string", "", ""},
        {"single character", "a", "a"},
        {"ascii string", "hello", "olleh"},
        {"unicode string", "привет", "тевирп"},
        {"with spaces", "hello world", "dlrow olleh"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Reverse(tt.input)
            if result != tt.expected {
                t.Errorf("Reverse(%q) = %q; expected %q", 
                    tt.input, result, tt.expected)
            }
        })
    }
}

func TestCountVowels(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected int
    }{
        {"no vowels", "bcdfg", 0},
        {"english vowels", "hello world", 3},
        {"russian vowels", "привет мир", 4},
        {"mixed case", "Hello World", 3},
        {"empty string", "", 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CountVowels(tt.input)
            if result != tt.expected {
                t.Errorf("CountVowels(%q) = %d; expected %d", 
                    tt.input, result, tt.expected)
            }
        })
    }
}

func TestIsPalindrome(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        {"empty string", "", true},
        {"single character", "a", true},
        {"palindrome ascii", "racecar", true},
        {"palindrome russian", "казак", true},
        {"not palindrome", "hello", false},
        {"with spaces", "а роза упала на лапу азора", true},
        {"case insensitive", "Racecar", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := IsPalindrome(tt.input)
            if result != tt.expected {
                t.Errorf("IsPalindrome(%q) = %t; expected %t", 
                    tt.input, result, tt.expected)
            }
        })
    }
}

---

🎯 Практика 2: Тестирование структур и методов

Задача: Протестировать систему управления банковским счетом

go
// bank.go
package bank

import (
    "errors"
    "fmt"
)

var (
    ErrInsufficientFunds = errors.New("insufficient funds")
    ErrInvalidAmount     = errors.New("invalid amount")
)

type Account struct {
    owner   string
    balance float64
}

func NewAccount(owner string, initialBalance float64) (*Account, error) {
    if initialBalance < 0 {
        return nil, ErrInvalidAmount
    }
    return &Account{owner: owner, balance: initialBalance}, nil
}

func (a *Account) Deposit(amount float64) error {
    if amount <= 0 {
        return ErrInvalidAmount
    }
    a.balance += amount
    return nil
}

func (a *Account) Withdraw(amount float64) error {
    if amount <= 0 {
        return ErrInvalidAmount
    }
    if amount > a.balance {
        return ErrInsufficientFunds
    }
    a.balance -= amount
    return nil
}

func (a *Account) Balance() float64 {
    return a.balance
}

func (a *Account) Owner() string {
    return a.owner
}

func (a *Account) Transfer(amount float64, to *Account) error {
    if err := a.Withdraw(amount); err != nil {
        return err
    }
    return to.Deposit(amount)
}


go
// bank_test.go
package bank

import (
    "testing"
)

func TestNewAccount(t *testing.T) {
    tests := []struct {
        name          string
        owner         string
        initialBalance float64
        expectedError bool
    }{
        {"valid account", "Alice", 100.0, false},
        {"zero balance", "Bob", 0.0, false},
        {"negative balance", "Charlie", -50.0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            acc, err := NewAccount(tt.owner, tt.initialBalance)
            
            if tt.expectedError {
                if err == nil {
                    t.Error("expected error, but got nil")
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if acc.Owner() != tt.owner {
                    t.Errorf("expected owner %q, got %q", tt.owner, acc.Owner())
                }
                if acc.Balance() != tt.initialBalance {
                    t.Errorf("expected balance %.2f, got %.2f", 
                        tt.initialBalance, acc.Balance())
                }
            }
        })
    }
}

func TestAccount_Deposit(t *testing.T) {
    acc, _ := NewAccount("Test", 100.0)
    
    tests := []struct {
        name          string
        amount        float64
        expectedBalance float64
        expectedError bool
    }{
        {"valid deposit", 50.0, 150.0, false},
        {"zero deposit", 0.0, 100.0, true},
        {"negative deposit", -10.0, 100.0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            initialBalance := acc.Balance()
            err := acc.Deposit(tt.amount)
            
            if tt.expectedError {
                if err == nil {
                    t.Error("expected error, but got nil")
                }
                if acc.Balance() != initialBalance {
                    t.Error("balance should not change on error")
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if acc.Balance() != tt.expectedBalance {
                    t.Errorf("expected balance %.2f, got %.2f", 
                        tt.expectedBalance, acc.Balance())
                }
            }
        })
    }
}

func TestAccount_Withdraw(t *testing.T) {
    acc, _ := NewAccount("Test", 100.0)
    
    tests := []struct {
        name          string
        amount        float64
        expectedBalance float64
        expectedError bool
    }{
        {"valid withdraw", 30.0, 70.0, false},
        {"withdraw zero", 0.0, 100.0, true},
        {"withdraw negative", -10.0, 100.0, true},
        {"insufficient funds", 150.0, 100.0, true},
        {"withdraw all", 100.0, 0.0, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            initialBalance := acc.Balance()
            err := acc.Withdraw(tt.amount)
            
            if tt.expectedError {
                if err == nil {
                    t.Error("expected error, but got nil")
                }
                if acc.Balance() != initialBalance {
                    t.Error("balance should not change on error")
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if acc.Balance() != tt.expectedBalance {
                    t.Errorf("expected balance %.2f, got %.2f", 
                        tt.expectedBalance, acc.Balance())
                }
            }
        })
    }
}

func TestAccount_Transfer(t *testing.T) {
    t.Run("successful transfer", func(t *testing.T) {
        alice, _ := NewAccount("Alice", 100.0)
        bob, _ := NewAccount("Bob", 50.0)
        
        err := alice.Transfer(30.0, bob)
        if err != nil {
            t.Errorf("unexpected error: %v", err)
        }
        
        if alice.Balance() != 70.0 {
            t.Errorf("expected Alice balance 70.0, got %.2f", alice.Balance())
        }
        
        if bob.Balance() != 80.0 {
            t.Errorf("expected Bob balance 80.0, got %.2f", bob.Balance())
        }
    })
    
    t.Run("insufficient funds for transfer", func(t *testing.T) {
        alice, _ := NewAccount("Alice", 30.0)
        bob, _ := NewAccount("Bob", 50.0)
        
        err := alice.Transfer(50.0, bob)
        if err != ErrInsufficientFunds {
            t.Errorf("expected ErrInsufficientFunds, got %v", err)
        }
        
        // Balances should not change
        if alice.Balance() != 30.0 {
            t.Error("Alice balance should not change on failed transfer")
        }
        if bob.Balance() != 50.0 {
            t.Error("Bob balance should not change on failed transfer")
        }
    })
}

---

🎯 Практика 3: Интеграционное тестирование

Задача: Протестировать HTTP handler с помощью тестового сервера

go
// handler.go
package api

import (
    "encoding/json"
    "net/http"
    "strconv"
)

type User struct {
    ID    int    json:"id"
    Name  string json:"name"
    Email string json:"email"
}

type UserService struct {
    users map[int]User
    nextID int
}

func NewUserService() *UserService {
    return &UserService{
        users: make(map[int]User),
        nextID: 1,
    }
}

func (s *UserService) CreateUser(name, email string) User {
    user := User{
        ID:    s.nextID,
        Name:  name,
        Email: email,
    }
    s.users[s.nextID] = user
    s.nextID++
    return user
}

func (s *UserService) GetUser(id int) (User, bool) {
    user, exists := s.users[id]
    return user, exists
}

func (s *UserService) GetAllUsers() []User {
    users := make([]User, 0, len(s.users))
    for _, user := range s.users {
        users = append(users, user)
    }
    return users
}

type Handler struct {
    userService *UserService
}

func NewHandler() *Handler {
    return &Handler{
        userService: NewUserService(),
    }
}

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        Name  string json:"name"
        Email string json:"email"
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if request.Name == "" || request.Email == "" {
        http.Error(w, "Name and email are required", http.StatusBadRequest)
        return
    }

    user := h.userService.CreateUser(request.Name, request.Email)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    user, exists := h.userService.GetUser(id)
    if !exists {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}


go
// handler_test.go
package api

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHandler_CreateUserHandler(t *testing.T) {
    handler := NewHandler()
    
    tests := []struct {
        name           string
        requestBody    map[string]string
        expectedStatus int
        expectedUser   bool
    }{
        {
            name: "valid user creation",
            requestBody: map[string]string{
                "name":  "Alice",
                "email": "alice@example.com",
            },
            expectedStatus: http.StatusCreated,
            expectedUser:   true,
        },
        {
            name: "missing name",
            requestBody: map[string]string{
                "email": "alice@example.com",
            },
            expectedStatus: http.StatusBadRequest,
            expectedUser:   false,
        },
        {
            name: "missing email",
            requestBody: map[string]string{
                "name": "Alice",
            },
            expectedStatus: http.StatusBadRequest,
            expectedUser:   false,
        },
        {
            name:           "empty body",
            requestBody:    map[string]string{},
            expectedStatus: http.StatusBadRequest,
            expectedUser:   false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            body, _ := json.Marshal(tt.requestBody)
            req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
            req.Header.Set("Content-Type", "application/json")
            
            rr := httptest.NewRecorder()
            handler.CreateUserHandler(rr, req)
            
            if rr.Code != tt.expectedStatus {
                t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
            }
            
            if tt.expectedUser {
                var user User
                if err := json.Unmarshal(rr.Body.Bytes(), &user); err != nil {
                    t.Errorf("failed to parse response: %v", err)
                }
                if user.Name != tt.requestBody["name"] {
                    t.Errorf("expected user name %q, got %q", 
                        tt.requestBody["name"], user.Name)
                }
            }
        })
    }
}

func TestHandler_GetUserHandler(t *testing.T) {
    handler := NewHandler()
    
    // Сначала создаем пользователя для тестирования
    handler.userService.CreateUser("Bob", "bob@example.com")
    
    tests := []struct {
        name           string
        userID         string
        expectedStatus int
        expectedUser   bool
    }{
        {
            name:           "existing user",
            userID:         "1",
            expectedStatus: http.StatusOK,
            expectedUser:   true,
        },
        {
            name:           "non-existing user",
            userID:         "999",
            expectedStatus: http.StatusNotFound,
            expectedUser:   false,
        },
        {
            name:           "invalid user ID",
            userID:         "abc",
            expectedStatus: http.StatusBadRequest,
            expectedUser:   false,
        },
        {
            name:           "missing user ID",
            userID:         "",
            expectedStatus: http.StatusBadRequest,
            expectedUser:   false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/user?id="+tt.userID, nil)
            rr := httptest.NewRecorder()
            
            handler.GetUserHandler(rr, req)
            
            if rr.Code != tt.expectedStatus {
                t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
            }
            
            if tt.expectedUser {
                var user User
                if err := json.Unmarshal(rr.Body.Bytes(), &user); err != nil {
                    t.Errorf("failed to parse response: %v", err)
                }
                if us