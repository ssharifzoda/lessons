# Занятие 19: Деплой приложений

`"Выкладываем Go приложения в production"`

---

## 📝 План на сегодня

1.  **Сборка приложения:** Go build и окружения
2.  **Docker:** Контейнеризация приложения
3.  **Environment variables:** Конфигурация
4.  **Деплой на сервер:** Процесс развертывания
5.  **HTTPS и домены:** Настройка веб-сервера

---

## 1. Сборка приложения

### Базовая сборка
```bash
# Для production
GOOS=linux GOARCH=amd64 go build -o app main.go

# С уменьшенным размером
go build -ldflags="-s -w" -o app main.go

# С информацией о версии
go build -ldflags="-X main.version=1.0.0" -o app main.go
```

### Мульти-стадийный Dockerfile
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o main .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

---

## 2. Docker композ

```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/app
      - JWT_SECRET=your-secret-key
    depends_on:
      - db

  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=app
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
    volumes:
      - db_data:/var/lib/postgresql/data

volumes:
  db_data:
```

---

## 3. Environment variables

```go
package config

import (
    "os"
    "strconv"
)

type Config struct {
    Port        int
    DatabaseURL string
    JWTSecret   string
    Debug       bool
}

func Load() *Config {
    return &Config{
        Port:        getEnvInt("PORT", 8080),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost/app"),
        JWTSecret:   getEnv("JWT_SECRET", "default-secret"),
        Debug:       getEnvBool("DEBUG", false),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if boolVal, err := strconv.ParseBool(value); err == nil {
            return boolVal
        }
    }
    return defaultValue
}
```

---

## 4. Деплой на сервер

### Простой скрипт деплоя
```bash
#!/bin/bash
# deploy.sh

echo "🚀 Starting deployment..."

# Останавливаем текущее приложение
sudo systemctl stop myapp

# Копируем новую версию
scp app user@server:/opt/myapp/

# Запускаем приложение
sudo systemctl start myapp

echo "✅ Deployment complete!"
```

### Systemd сервис
```ini
# /etc/systemd/system/myapp.service
[Unit]
Description=My Go Application
After=network.target

[Service]
Type=simple
User=appuser
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/app
Restart=always
Environment=DATABASE_URL=postgres://user:pass@localhost/app
Environment=JWT_SECRET=your-production-secret

[Install]
WantedBy=multi-user.target
```

---

## 5. Nginx конфиг

```nginx
# /etc/nginx/sites-available/myapp
server {
    listen 80;
    server_name yourdomain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

---

## 🎯 Практика: Production-ready приложение

```go
package main

import (
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/yourname/yourapp/config"
    "github.com/yourname/yourapp/database"
    "github.com/yourname/yourapp/router"
)

func main() {
    // Загружаем конфигурацию
    cfg := config.Load()
    
    // Подключаемся к БД
    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Создаем роутер
    r := router.New(db, cfg.JWTSecret)
    
    // Graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    
    server := &http.Server{
        Addr:    ":" + strconv.Itoa(cfg.Port),
        Handler: r,
    }
    
    go func() {
        log.Printf("Server starting on port %d", cfg.Port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Server failed:", err)
        }
    }()
    
    <-stop
    log.Println("Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server shutdown failed:", err)
    }
    
    log.Println("Server stopped")
}
```

---

## ❓ Важные моменты

### Безопасность:
- Никогда не коммитьте секреты в git
- Используйте разные секреты для dev/prod
- Настройте firewall

### Мониторинг:
```bash
# Проверка статуса
sudo systemctl status myapp

# Логи
journalctl -u myapp -f

# Health check
curl http://localhost:8080/health
```

---

## 🏠 Домашнее задание

**Задача 1:** Соберите Docker образ вашего приложения
**Задача 2:** Настройте конфигурацию через environment variables  
**Задача 3:** Создайте простой скрипт деплоя

---

Следующее занятие: **Финальный проект**