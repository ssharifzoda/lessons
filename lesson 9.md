# Занятие 9: Интерфейсы в Go

`"Контракты для структур и полиморфизм"`

---

## 📝 План на сегодня

1.  **Что такое интерфейсы?** Концепция контрактов
2.  **Объявление и реализация** интерфейсов
3.  **Полиморфизм в Go:** Один интерфейс - много реализаций
4.  **Пустой интерфейс interface{}:** Универсальный тип
5.  **Type assertions и type switches:** Работа с пустыми интерфейсами
6.  **Практика:** Создание гибких систем с интерфейсами

---

## 1. Что такое интерфейсы?

**Интерфейс** - это набор методов, которые должен реализовать тип.

### Аналогия:
- **Интерфейс** = Контракт (требования к работе)
- **Структура** = Исполнитель (реализует контракт)

```go
// Интерфейс определяет "что делать"
type Speaker interface {
    Speak() string
}

// Структуры реализуют "как делать"
type Dog struct {
    Name string
}

type Cat struct {
    Name string
}

// Dog реализует интерфейс Speaker
func (d Dog) Speak() string {
    return "Гав! Меня зовут " + d.Name
}

// Cat реализует интерфейс Speaker
func (c Cat) Speak() string {
    return "Мяу! Я " + c.Name
}
```

---

## 2. Объявление и реализация интерфейсов

### Объявление интерфейса:
```go
// Геометрические фигуры
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Переводчик
type Translator interface {
    Translate(text string) string
    GetLanguage() string
}
```

### Реализация интерфейса:
```go
type Rectangle struct {
    Width, Height float64
}

// Rectangle реализует интерфейс Shape
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Теперь Rectangle можно использовать везде, где ожидается Shape
```

---

## 3. Полиморфизм в Go

**Полиморфизм** - возможность использовать разные типы через общий интерфейс.

```go
package main
import "fmt"
import "math"

// Интерфейс
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Структура Круг
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Структура Прямоугольник
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Функция, работающая с любым Shape
func PrintShapeInfo(s Shape) {
    fmt.Printf("Площадь: %.2f\n", s.Area())
    fmt.Printf("Периметр: %.2f\n", s.Perimeter())
    fmt.Println()
}

func main() {
    shapes := []Shape{
        Circle{Radius: 5},
        Rectangle{Width: 4, Height: 6},
        Circle{Radius: 3},
    }
    
    for _, shape := range shapes {
        PrintShapeInfo(shape)
    }
}
```

---

## 4. Пустой интерфейс interface{}

**Пустой интерфейс** не требует никаких методов - ему удовлетворяет ЛЮБОЙ тип.

```go
// interface{} - может хранить значение любого типа
func PrintAny(value interface{}) {
    fmt.Printf("Значение: %v, Тип: %T\n", value, value)
}

func main() {
    PrintAny("Привет")        // string
    PrintAny(42)              // int
    PrintAny(3.14)            // float64
    PrintAny([]int{1, 2, 3})  // []int
    PrintAny(true)            // bool
}
```

---

## 5. Type assertions и type switches

### Type assertion (утверждение типа):
```go
func processValue(val interface{}) {
    // Пытаемся преобразовать к string
    if str, ok := val.(string); ok {
        fmt.Printf("Это строка: %s\n", str)
    } else {
        fmt.Printf("Это не строка: %v\n", val)
    }
}

func main() {
    processValue("текст")  // Это строка: текст
    processValue(123)      // Это не строка: 123
}
```

### Type switch (переключатель типов):
```go
func describe(value interface{}) {
    switch v := value.(type) {
    case string:
        fmt.Printf("Строка: %s (длина: %d)\n", v, len(v))
    case int:
        fmt.Printf("Целое число: %d\n", v)
    case bool:
        fmt.Printf("Булево значение: %t\n", v)
    default:
        fmt.Printf("Неизвестный тип: %T\n", v)
    }
}

func main() {
    describe("Привет")     // Строка: Привет (длина: 6)
    describe(42)           // Целое число: 42
    describe(true)         // Булево значение: true
    describe(3.14)         // Неизвестный тип: float64
}
```

---

## 🎯 Практика 1: Система платежей

**Задача:** Создать систему, поддерживающую разные способы оплаты

```go
package main
import "fmt"

// Интерфейс платежной системы
type PaymentMethod interface {
    ProcessPayment(amount float64) bool
    GetName() string
}

// Кредитная карта
type CreditCard struct {
    Number string
    Owner  string
}

func (c CreditCard) ProcessPayment(amount float64) bool {
    fmt.Printf("Обрабатываем платеж %.2f через кредитную карту %s\n", amount, c.Number)
    return true
}

func (c CreditCard) GetName() string {
    return "Кредитная карта"
}

// PayPal
type PayPal struct {
    Email string
}

func (p PayPal) ProcessPayment(amount float64) bool {
    fmt.Printf("Обрабатываем платеж %.2f через PayPal %s\n", amount, p.Email)
    return true
}

func (p PayPal) GetName() string {
    return "PayPal"
}

// Банковский перевод
type BankTransfer struct {
    AccountNumber string
}

func (b BankTransfer) ProcessPayment(amount float64) bool {
    fmt.Printf("Обрабатываем платеж %.2f через банковский перевод %s\n", amount, b.AccountNumber)
    return true
}

func (b BankTransfer) GetName() string {
    return "Банковский перевод"
}

// Функция для обработки заказа
func ProcessOrder(amount float64, method PaymentMethod) {
    fmt.Printf("\nОбработка заказа на сумму %.2f\n", amount)
    fmt.Printf("Способ оплаты: %s\n", method.GetName())
    
    if method.ProcessPayment(amount) {
        fmt.Println("✅ Платеж успешно обработан")
    } else {
        fmt.Println("❌ Ошибка обработки платежа")
    }
}

func main() {
    paymentMethods := []PaymentMethod{
        CreditCard{Number: "1234-5678-9012-3456", Owner: "Анна"},
        PayPal{Email: "anna@example.com"},
        BankTransfer{AccountNumber: "RU1234567890"},
    }
    
    for _, method := range paymentMethods {
        ProcessOrder(1500.50, method)
    }
}
```

---

## 🎯 Практика 2: Универсальное хранилище

**Задача:** Создать хранилище, которое может работать с любыми типами данных

```go
package main
import "fmt"

type Storage struct {
    data map[string]interface{}
}

func NewStorage() *Storage {
    return &Storage{
        data: make(map[string]interface{}),
    }
}

func (s *Storage) Set(key string, value interface{}) {
    s.data[key] = value
}

func (s *Storage) Get(key string) (interface{}, bool) {
    value, exists := s.data[key]
    return value, exists
}

func (s *Storage) Delete(key string) {
    delete(s.data, key)
}

func (s *Storage) PrintAll() {
    fmt.Println("Содержимое хранилища:")
    for key, value := range s.data {
        fmt.Printf("  %s: %v (%T)\n", key, value, value)
    }
}

func main() {
    storage := NewStorage()
    
    // Сохраняем разные типы данных
    storage.Set("name", "Анна")
    storage.Set("age", 25)
    storage.Set("scores", []int{95, 87, 92})
    storage.Set("isStudent", false)
    storage.Set("height", 170.5)
    
    storage.PrintAll()
    
    // Получаем и обрабатываем данные
    if name, exists := storage.Get("name"); exists {
        if str, ok := name.(string); ok {
            fmt.Printf("\nИмя: %s\n", str)
        }
    }
    
    if age, exists := storage.Get("age"); exists {
        switch v := age.(type) {
        case int:
            fmt.Printf("Возраст: %d лет\n", v)
        default:
            fmt.Printf("Неожиданный тип для возраста: %T\n", v)
        }
    }
}
```

---

## 🎯 Практика 3: Система уведомлений

**Задача:** Создать систему, отправляющую уведомления разными способами

```go
package main
import "fmt"

// Интерфейс уведомления
type Notifier interface {
    Send(message string) bool
    GetType() string
}

// Email уведомление
type EmailNotifier struct {
    Address string
}

func (e EmailNotifier) Send(message string) bool {
    fmt.Printf("Отправляем email на %s: %s\n", e.Address, message)
    return true
}

func (e EmailNotifier) GetType() string {
    return "Email"
}

// SMS уведомление
type SMSNotifier struct {
    Phone string
}

func (s SMSNotifier) Send(message string) bool {
    fmt.Printf("Отправляем SMS на %s: %s\n", s.Phone, message)
    return true
}

func (s SMSNotifier) GetType() string {
    return "SMS"
}

// Push уведомление
type PushNotifier struct {
    DeviceID string
}

func (p PushNotifier) Send(message string) bool {
    fmt.Printf("Отправляем Push на устройство %s: %s\n", p.DeviceID, message)
    return true
}

func (p PushNotifier) GetType() string {
    return "Push"
}

// Система уведомлений
type NotificationSystem struct {
    notifiers []Notifier
}

func (ns *NotificationSystem) AddNotifier(notifier Notifier) {
    ns.notifiers = append(ns.notifiers, notifier)
}

func (ns *NotificationSystem) Broadcast(message string) {
    fmt.Printf("Рассылаем уведомление: %s\n", message)
    fmt.Println("---")
    
    for _, notifier := range ns.notifiers {
        fmt.Printf("Тип: %s -> ", notifier.GetType())
        notifier.Send(message)
    }
}

func main() {
    system := NotificationSystem{}
    
    // Добавляем разные способы уведомления
    system.AddNotifier(EmailNotifier{Address: "user@example.com"})
    system.AddNotifier(SMSNotifier{Phone: "+7-900-123-45-67"})
    system.AddNotifier(PushNotifier{DeviceID: "device-12345"})
    
    // Рассылаем уведомление
    system.Broadcast("Ваш заказ готов к выдаче!")
}
```

---

## ❓ Важные моменты

### Интерфейсы реализуются неявно:
```go
type Writer interface {
    Write([]byte) (int, error)
}

// Любой тип с методом Write автоматически реализует Writer
type MyWriter struct{}

func (mw MyWriter) Write(data []byte) (int, error) {
    // реализация
    return len(data), nil
}
// MyWriter автоматически реализует Writer!
```

### Композиция интерфейсов:
```go
type Reader interface {
    Read([]byte) (int, error)
}

type ReadWriter interface {
    Reader
    Writer
}
// ReadWriter требует и Read, и Write методы
```

---

## 🏠 Домашнее задание

**Задача 1: Система сортировки**
Создайте интерфейс `Sorter` с методами:
- `Sort([]int) []int`
- `GetName() string`

Реализуйте его для разных алгоритмов сортировки:
- `BubbleSorter`, `QuickSorter`, `MergeSorter`

**Задача 2: Универсальный кэш**
Создайте интерфейс `Cache` с методами:
- `Set(key string, value interface{})`
- `Get(key string) (interface{}, bool)`
- `Delete(key string)`
- `Clear()`

Реализуйте две версии:
- `MemoryCache` (хранит в map)
- `FileCache` (сохраняет в файл)

**Задача 3: Система отчетов**
Создайте интерфейс `ReportGenerator` с методом `Generate(data interface{}) string`.
Реализуйте генераторы отчетов в разных форматах:
- `JSONReport`, `XMLReport`, `TextReport`

---

## 🚀 Что ждет на следующем занятии?

*   **Обработка ошибок:** Идиоматический подход Go
*   **Создание собственных ошибок**
*   **Panic и recover:** Экстренные ситуации

**Удачи в решении задач! 🎉**