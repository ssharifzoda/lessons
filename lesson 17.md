# Занятие 17: Продвинутая работа с БД и Context

`"Эффективное управление соединениями и отмена операций"`

---

## 📝 План на сегодня

1.  **Context в БД операциях:** Таймауты и отмена
2.  **Подготовленные запросы:** Безопасность и производительность
3.  **Транзакции:** Группировка операций
4.  **Connection pool:** Настройка пула соединений
5.  **Миграции:** Управление схемой БД

---

## 1. Context в БД операциях

```go
import (
    "context"
    "database/sql"
    "time"
)

func getUserWithTimeout(db *sql.DB, userID int) (*User, error) {
    // 3 секунды на выполнение запроса
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    var user User
    err := db.QueryRowContext(ctx, 
        "SELECT id, name FROM users WHERE id = $1", userID).
        Scan(&user.ID, &user.Name)
    
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Отмена по внешнему контексту
func processUser(ctx context.Context, db *sql.DB, userID int) error {
    var user User
    err := db.QueryRowContext(ctx, 
        "SELECT id, name FROM users WHERE id = $1", userID).
        Scan(&user.ID, &user.Name)
    
    if err != nil {
        return err
    }
    
    // Дальнейшая обработка...
    return nil
}
```

---

## 2. Подготовленные запросы

```go
type UserRepo struct {
    getUserStmt *sql.Stmt
}

func NewUserRepo(db *sql.DB) (*UserRepo, error) {
    // Подготавливаем запрос один раз
    getUserStmt, err := db.Prepare("SELECT id, name, email FROM users WHERE id = $1")
    if err != nil {
        return nil, err
    }
    
    return &UserRepo{getUserStmt: getUserStmt}, nil
}

func (r *UserRepo) GetUser(userID int) (*User, error) {
    var user User
    err := r.getUserStmt.QueryRow(userID).Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Не забудьте закрыть
func (r *UserRepo) Close() error {
    return r.getUserStmt.Close()
}
```

---

## 3. Транзакции

```go
func transferMoney(db *sql.DB, fromID, toID int, amount float64) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // Safe rollback if not committed

    // Снимаем деньги
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1",
        amount, fromID)
    if err != nil {
        return err
    }

    // Добавляем деньги
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2", 
        amount, toID)
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

---

## 4. Connection Pool настройки

```go
func setupDB() *sql.DB {
    db, err := sql.Open("postgres", "connection-string")
    if err != nil {
        log.Fatal(err)
    }

    // Настройка пула соединений
    db.SetMaxOpenConns(25)      // Максимум соединений
    db.SetMaxIdleConns(25)      // Максимум idle соединений  
    db.SetConnMaxLifetime(5 * time.Minute) // Время жизни соединения
    db.SetConnMaxIdleTime(2 * time.Minute) // Время idle соединения

    return db
}
```

---

## 5. Простые миграции

```go
// migrations.go
var migrations = []string{
    `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        email VARCHAR(100) UNIQUE
    )`,
    
    `CREATE TABLE IF NOT EXISTS posts (
        id SERIAL PRIMARY KEY, 
        title VARCHAR(200),
        user_id INTEGER REFERENCES users(id)
    )`,
}

func runMigrations(db *sql.DB) error {
    for i, migration := range migrations {
        _, err := db.Exec(migration)
        if err != nil {
            return fmt.Errorf("migration %d failed: %w", i+1, err)
        }
    }
    return nil
}
```

---

## 🎯 Практика: Оптимизированный UserService

```go
type UserService struct {
    db *sql.DB
    getUserStmt *sql.Stmt
}

func NewUserService(db *sql.DB) (*UserService, error) {
    getUserStmt, err := db.Prepare(`
        SELECT id, name, email FROM users WHERE id = $1
    `)
    if err != nil {
        return nil, err
    }

    return &UserService{db: db, getUserStmt: getUserStmt}, nil
}

func (s *UserService) GetUser(ctx context.Context, userID int) (*User, error) {
    var user User
    err := s.getUserStmt.QueryRowContext(ctx, userID).
        Scan(&user.ID, &user.Name, &user.Email)
    
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    err = tx.QueryRowContext(ctx,
        "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
        user.Name, user.Email).Scan(&user.ID)
    
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (s *UserService) Close() error {
    return s.getUserStmt.Close()
}
```

---

## ❓ Важные моменты

### Всегда используйте:
- `QueryRowContext()` вместо `QueryRow()`
- `ExecContext()` вместо `Exec()`
- `BeginTx()` вместо `Begin()`

### Не забывайте:
- `defer tx.Rollback()` в транзакциях
- Закрывать подготовленные запросы
- Проверять `Rows.Err()` после `Rows.Next()`

---

## 🏠 Домашнее задание

**Задача 1:** Добавьте таймауты ко всем БД операциям
**Задача 2:** Реализуйте транзакцию для создания заказа с товарами  
**Задача 3:** Настройте connection pool для вашего приложения

---

Следующее занятие: **Middleware и обработка ошибок**