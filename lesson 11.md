# Занятие 11: Работа с файлами и строками

`"Чтение, запись и обработка данных из файлов"`

---

## 📝 План на сегодня

1.  **Пакеты os и io:** Основы работы с файловой системой
2.  **Чтение файлов:** Разные способы чтения данных
3.  **Запись файлов:** Создание и обновление файлов
4.  **Работа с JSON:** Маршалинг и демаршалинг данных
5.  **Обработка CSV файлов:** Работа с табличными данными
6.  **Практика:** Реальные кейсы работы с файлами

---

## 1. Пакеты os и io - основы

### Основные типы и функции:
```go
import (
    "os"
    "io"
    "fmt"
)

// Открытие файла
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close() // Важно закрывать файлы!

// Создание файла
newFile, err := os.Create("newfile.txt")
if err != nil {
    return err
}
defer newFile.Close()

// Проверка существования файла
if _, err := os.Stat("file.txt"); os.IsNotExist(err) {
    fmt.Println("Файл не существует")
}
```

---

## 2. Чтение файлов

### Способ 1: Чтение всего файла сразу (ioutil.ReadFile)
```go
import "os"

func readEntireFile(filename string) (string, error) {
    content, err := os.ReadFile(filename)
    if err != nil {
        return "", err
    }
    return string(content), nil
}
```

### Способ 2: Построчное чтение (bufio.Scanner)
```go
import (
    "bufio"
    "os"
)

func readFileLineByLine(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    
    return lines, nil
}
```

### Способ 3: Чтение с буфером
```go
func readWithBuffer(filename string) (string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return "", err
    }
    defer file.Close()

    buffer := make([]byte, 1024) // Буфер 1KB
    var content []byte
    
    for {
        n, err := file.Read(buffer)
        if err == io.EOF {
            break
        }
        if err != nil {
            return "", err
        }
        content = append(content, buffer[:n]...)
    }
    
    return string(content), nil
}
```

---

## 3. Запись файлов

### Способ 1: Запись всего содержимого
```go
func writeEntireFile(filename, content string) error {
    return os.WriteFile(filename, []byte(content), 0644)
}
```

### Способ 2: Построчная запись
```go
func writeLines(filename string, lines []string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    defer writer.Flush()

    for _, line := range lines {
        _, err := writer.WriteString(line + "\n")
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### Способ 3: Дозапись в файл
```go
func appendToFile(filename, content string) error {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(content + "\n")
    return err
}
```

---

## 4. Работа с JSON

### Структуры в JSON (маршалинг):
```go
import "encoding/json"

type Person struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string `json:"email,omitempty"` // omitempty - пропускать если пустое
    City    string `json:"city"`
}

func savePersonToJSON(filename string, person Person) error {
    // Преобразуем структуру в JSON
    jsonData, err := json.MarshalIndent(person, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(filename, jsonData, 0644)
}
```

### JSON в структуры (демаршалинг):
```go
func loadPersonFromJSON(filename string) (*Person, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var person Person
    if err := json.Unmarshal(data, &person); err != nil {
        return nil, err
    }

    return &person, nil
}
```

---

## 5. Обработка CSV файлов

### Чтение CSV:
```go
import (
    "encoding/csv"
    "os"
)

func readCSV(filename string) ([][]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    return records, nil
}
```

### Запись CSV:
```go
func writeCSV(filename string, records [][]string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, record := range records {
        if err := writer.Write(record); err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## 🎯 Практика 1: Система управления конфигурацией

**Задача:** Создать систему для работы с конфигурационными файлами

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type Config struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    Logging  LoggingConfig  `json:"logging"`
}

type ServerConfig struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Name     string `json:"name"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoggingConfig struct {
    Level string `json:"level"`
    File  string `json:"file"`
}

type ConfigManager struct {
    filename string
}

func NewConfigManager(filename string) *ConfigManager {
    return &ConfigManager{filename: filename}
}

func (cm *ConfigManager) Load() (*Config, error) {
    data, err := os.ReadFile(cm.filename)
    if err != nil {
        return nil, fmt.Errorf("не удалось прочитать конфиг: %w", err)
    }

    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
    }

    return &config, nil
}

func (cm *ConfigManager) Save(config *Config) error {
    data, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return fmt.Errorf("ошибка создания JSON: %w", err)
    }

    if err := os.WriteFile(cm.filename, data, 0644); err != nil {
        return fmt.Errorf("не удалось записать конфиг: %w", err)
    }

    return nil
}

func (cm *ConfigManager) CreateDefault() error {
    defaultConfig := &Config{
        Server: ServerConfig{
            Host: "localhost",
            Port: 8080,
        },
        Database: DatabaseConfig{
            Host:     "localhost",
            Port:     5432,
            Name:     "mydb",
            Username: "admin",
            Password: "secret",
        },
        Logging: LoggingConfig{
            Level: "info",
            File:  "app.log",
        },
    }

    return cm.Save(defaultConfig)
}

func main() {
    configManager := NewConfigManager("config.json")
    
    // Создаем конфиг по умолчанию если файла нет
    if _, err := os.Stat("config.json"); os.IsNotExist(err) {
        fmt.Println("Создаем конфиг по умолчанию...")
        if err := configManager.CreateDefault(); err != nil {
            fmt.Printf("Ошибка: %v\n", err)
            return
        }
    }
    
    // Загружаем конфиг
    config, err := configManager.Load()
    if err != nil {
        fmt.Printf("Ошибка загрузки конфига: %v\n", err)
        return
    }
    
    fmt.Printf("Конфигурация загружена:\n")
    fmt.Printf("Сервер: %s:%d\n", config.Server.Host, config.Server.Port)
    fmt.Printf("База данных: %s@%s:%d/%s\n", 
        config.Database.Username, config.Database.Host, 
        config.Database.Port, config.Database.Name)
    fmt.Printf("Логирование: уровень=%s, файл=%s\n", 
        config.Logging.Level, config.Logging.File)
}
```

---

## 🎯 Практика 2: Анализатор логов

**Задача:** Создать инструмент для анализа лог-файлов

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"
)

type LogEntry struct {
    Timestamp time.Time
    Level     string
    Message   string
}

type LogAnalyzer struct {
    entries []LogEntry
}

func NewLogAnalyzer() *LogAnalyzer {
    return &LogAnalyzer{
        entries: make([]LogEntry, 0),
    }
}

func (la *LogAnalyzer) ParseLogFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineNumber := 0
    
    for scanner.Scan() {
        lineNumber++
        line := scanner.Text()
        
        entry, err := la.parseLogLine(line)
        if err != nil {
            fmt.Printf("Ошибка парсинга строки %d: %v\n", lineNumber, err)
            continue
        }
        
        la.entries = append(la.entries, entry)
    }
    
    return scanner.Err()
}

func (la *LogAnalyzer) parseLogLine(line string) (LogEntry, error) {
    parts := strings.SplitN(line, " ", 3)
    if len(parts) < 3 {
        return LogEntry{}, fmt.Errorf("неверный формат лога")
    }
    
    timestamp, err := time.Parse("2006-01-02T15:04:05", parts[0])
    if err != nil {
        return LogEntry{}, fmt.Errorf("неверный формат времени: %w", err)
    }
    
    return LogEntry{
        Timestamp: timestamp,
        Level:     parts[1],
        Message:   parts[2],
    }, nil
}

func (la *LogAnalyzer) CountByLevel() map[string]int {
    counts := make(map[string]int)
    
    for _, entry := range la.entries {
        counts[entry.Level]++
    }
    
    return counts
}

func (la *LogAnalyzer) FindErrors() []LogEntry {
    var errors []LogEntry
    
    for _, entry := range la.entries {
        if entry.Level == "ERROR" {
            errors = append(errors, entry)
        }
    }
    
    return errors
}

func (la *LogAnalyzer) GenerateReport() string {
    levelCounts := la.CountByLevel()
    errors := la.FindErrors()
    
    var report strings.Builder
    
    report.WriteString("=== ОТЧЕТ ПО ЛОГАМ ===\n\n")
    report.WriteString("Статистика по уровням:\n")
    for level, count := range levelCounts {
        report.WriteString(fmt.Sprintf("  %s: %d\n", level, count))
    }
    
    report.WriteString(fmt.Sprintf("\nВсего ошибок: %d\n", len(errors)))
    if len(errors) > 0 {
        report.WriteString("\nПоследние ошибки:\n")
        for i := 0; i < len(errors) && i < 5; i++ {
            report.WriteString(fmt.Sprintf("  %s: %s\n", 
                errors[i].Timestamp.Format("15:04:05"), 
                errors[i].Message))
        }
    }
    
    return report.String()
}

func main() {
    // Создаем тестовый лог-файл
    testLog := []string{
        "2024-01-15T10:00:01 INFO Сервер запущен",
        "2024-01-15T10:00:15 DEBUG Подключение к базе данных",
        "2024-01-15T10:01:23 WARN Медленный запрос",
        "2024-01-15T10:02:45 ERROR Ошибка подключения к БД",
        "2024-01-15T10:03:12 INFO Восстановление подключения",
        "2024-01-15T10:04:33 ERROR Таймаут операции",
        "2024-01-15T10:05:47 INFO Операция завершена",
    }
    
    // Записываем тестовый лог
    if err := os.WriteFile("app.log", []byte(strings.Join(testLog, "\n")), 0644); err != nil {
        fmt.Printf("Ошибка создания лог-файла: %v\n", err)
        return
    }
    
    // Анализируем логи
    analyzer := NewLogAnalyzer()
    if err := analyzer.ParseLogFile("app.log"); err != nil {
        fmt.Printf("Ошибка анализа логов: %v\n", err)
        return
    }
    
    // Генерируем отчет
    report := analyzer.GenerateReport()
    fmt.Println(report)
    
    // Сохраняем отчет в файл
    if err := os.WriteFile("log_report.txt", []byte(report), 0644); err != nil {
        fmt.Printf("Ошибка сохранения отчета: %v\n", err)
        return
    }
    
    fmt.Println("Отчет сохранен в log_report.txt")
}
```

---

## 🎯 Практика 3: Система управления задачами с сохранением в CSV

**Задача:** Создать менеджер задач с сохранением в CSV формате

```go
package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "time"
)

type Task struct {
    ID          int
    Title       string
    Description string
    Priority    string
    Completed   bool
    CreatedAt   time.Time
    DueDate     time.Time
}

type TaskManager struct {
    tasks    []Task
    filename string
    nextID   int
}

func NewTaskManager(filename string) *TaskManager {
    return &TaskManager{
        tasks:    make([]Task, 0),
        filename: filename,
        nextID:   1,
    }
}

func (tm *TaskManager) LoadTasks() error {
    file, err := os.Open(tm.filename)
    if os.IsNotExist(err) {
        return nil // Файл не существует - это нормально
    }
    if err != nil {
        return err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return err
    }

    for _, record := range records {
        if len(record) != 7 {
            continue // Пропускаем некорректные строки
        }

        id, _ := strconv.Atoi(record[0])
        completed, _ := strconv.ParseBool(record[4])
        createdAt, _ := time.Parse(time.RFC3339, record[5])
        dueDate, _ := time.Parse(time.RFC3339, record[6])

        task := Task{
            ID:          id,
            Title:       record[1],
            Description: record[2],
            Priority:    record[3],
            Completed:   completed,
            CreatedAt:   createdAt,
            DueDate:     dueDate,
        }

        tm.tasks = append(tm.tasks, task)
        
        if id >= tm.nextID {
            tm.nextID = id + 1
        }
    }

    return nil
}

func (tm *TaskManager) SaveTasks() error {
    file, err := os.Create(tm.filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, task := range tm.tasks {
        record := []string{
            strconv.Itoa(task.ID),
            task.Title,
            task.Description,
            task.Priority,
            strconv.FormatBool(task.Completed),
            task.CreatedAt.Format(time.RFC3339),
            task.DueDate.Format(time.RFC3339),
        }

        if err := writer.Write(record); err != nil {
            return err
        }
    }

    return nil
}

func (tm *TaskManager) AddTask(title, description, priority string, dueDate time.Time) {
    task := Task{
        ID:          tm.nextID,
        Title:       title,
        Description: description,
        Priority:    priority,
        Completed:   false,
        CreatedAt:   time.Now(),
        DueDate:     dueDate,
    }

    tm.tasks = append(tm.tasks, task)
    tm.nextID++
}

func (tm *TaskManager) CompleteTask(id int) bool {
    for i := range tm.tasks {
        if tm.tasks[i].ID == id {
            tm.tasks[i].Completed = true
            return true
        }
    }
    return false
}

func (tm *TaskManager) GetPendingTasks() []Task {
    var pending []Task
    for _, task := range tm.tasks {
        if !task.Completed {
            pending = append(pending, task)
        }
    }
    return pending
}

func (tm *TaskManager) PrintTasks() {
    fmt.Println("=== ЗАДАЧИ ===")
    for _, task := range tm.tasks {
        status := "❌"
        if task.Completed {
            status = "✅"
        }
        
        fmt.Printf("%s [%s] %s (ID: %d)\n", 
            status, task.Priority, task.Title, task.ID)
        fmt.Printf("   Описание: %s\n", task.Description)
        fmt.Printf("   Срок: %s\n", task.DueDate.Format("02.01.2006"))
        fmt.Println()
    }
}

func main() {
    taskManager := NewTaskManager("tasks.csv")
    
    // Загружаем существующие задачи
    if err := taskManager.LoadTasks(); err != nil {
        fmt.Printf("Ошибка загрузки задач: %v\n", err)
        return
    }
    
    // Добавляем тестовые задачи если их нет
    if len(taskManager.tasks) == 0 {
        fmt.Println("Добавляем тестовые задачи...")
        
        taskManager.AddTask(
            "Изучить Go", 
            "Пройти курс по Go программированию", 
            "high", 
            time.Now().AddDate(0, 1, 0),
        )
        
        taskManager.AddTask(
            "Сделать ДЗ", 
            "Выполнить домашние задания по файлам", 
            "medium", 
            time.Now().AddDate(0, 0, 7),
        )
        
        // Сохраняем задачи
        if err := taskManager.SaveTasks(); err != nil {
            fmt.Printf("Ошибка сохранения задач: %v\n", err)
            return
        }
    }
    
    // Показываем задачи
    taskManager.PrintTasks()
    
    // Показываем pending задачи
    pending := taskManager.GetPendingTasks()
    fmt.Printf("Незавершенных задач: %d\n", len(pending))
    
    // Сохраняем отчет по pending задачам
    var report strings.Builder
    report.WriteString("НЕЗАВЕРШЕННЫЕ ЗАДАЧИ:\n\n")
    for _, task := range pending {
        report.WriteString(fmt.Sprintf("- %s [%s]\n", task.Title, task.Priority))
        report.WriteString(fmt.Sprintf("  Срок: %s\n", task.DueDate.Format("02.01.2006")))
        report.WriteString(fmt.Sprintf("  ID: %d\n\n", task.ID))
    }
    
    if err := os.WriteFile("pending_tasks.txt", []byte(report.String()), 0644); err != nil {
        fmt.Printf("Ошибка сохранения отчета: %v\n", err)
    } else {
        fmt.Println("Отчет по незавершенным задачам сохранен в pending_tasks.txt")
    }
}
```

---

## ❓ Важные моменты

### Всегда закрывайте файлы:
```go
// ХОРОШО:
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close() // Гарантированное закрытие

// ПЛОХО:
file, _ := os.Open("file.txt")
// Забыли закрыть - утечка файлового дескриптора
```

### Обработка больших файлов:
```go
// Для больших файлов используйте буферизованное чтение
// чтобы не загружать весь файл в память
```

---

## 🏠 Домашнее задание

**Задача 1: Конвертер форматов**
Напишите программу, которая:
- Читает CSV файл с данными пользователей
- Конвертирует его в JSON формат
- Сохраняет результат в новый файл

**Задача 2: Анализатор текста**
Создайте программу, которая:
- Читает текстовый файл
- Подсчитывает статистику (количество слов, строк, символов)
- Находит самые частые слова
- Сохраняет отчет в отдельный файл

**Задача 3: Резервное копирование**
Напишите утилиту, которая:
- Копирует все файлы из указанной директории
- Сохраняет их в zip-архив с датой в имени
- Вести лог операций копирования

---

**Удачи в решении задач! 🎉**