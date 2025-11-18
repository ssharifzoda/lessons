# Занятие 20: Git и совместная разработка

`"Эффективная работа с Git в команде"`

---

## 📝 План на сегодня

1.  **Основы Git:** Коммиты, ветки, пуши
2.  **Рабочие процессы:** Feature branches, PR, code review
3.  **Git в Go проектах:** .gitignore, go.mod
4.  **Решение конфликтов:** Merge и rebase
5.  **Best practices:** Хорошие привычки

---

## 1. Основы Git

### Базовые команды
```bash
# Инициализация
git init
git remote add origin https://github.com/user/repo.git

# Работа с изменениями
git add .
git commit -m "feat: add user authentication"
git push origin main

# Получение обновлений
git pull origin main
git fetch origin
```

### Полезные алиасы
```bash
# Добавить в ~/.gitconfig
[alias]
    co = checkout
    br = branch
    ci = commit
    st = status
    last = log -1 HEAD
    unstage = reset HEAD --
```

---

## 2. Рабочие процессы

### Feature Branch Workflow
```bash
# Создание feature ветки
git checkout -b feat/user-auth

# Регулярные коммиты
git add .
git commit -m "feat: add jwt token generation"
git commit -m "test: add auth middleware tests"

# Пуш в remote
git push -u origin feat/user-auth

# Создание Pull Request через GitHub UI
```

### Структура коммитов
```
feat: add user registration endpoint
fix: resolve database connection leak
docs: update API documentation
test: add unit tests for user service
refactor: simplify error handling
chore: update dependencies
```

---

## 3. Git в Go проектах

### .gitignore для Go
```gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out

# IDE
.vscode/
.idea/
*.swp
*.swo

# Go
/bin/
/pkg/
/go.sum

# Environment
.env
.env.local
```

### Работа с go.mod
```bash
# Добавление зависимости
go get github.com/gorilla/mux@v1.8.0

# Обновление зависимостей
go mod tidy
go mod download

# Ветвление с зависимостями
git add go.mod go.sum
git commit -m "chore: add gorilla/mux dependency"
```

---

## 4. Решение конфликтов

### Merge конфликты
```bash
# Пулл с rebase
git pull --rebase origin main

# При конфликте
git status                    # Посмотреть конфликтующие файлы
# Редактируем файлы, убираем <<<<<<<, =======, >>>>>>>
git add .
git rebase --continue

# Или отмена
git rebase --abort
```

### Стратегии решения
```go
// Конфликт в Go файле
<<<<<<< HEAD
func NewService() *Service {
    return &Service{db: db, cache: redis}
}
=======
func NewService(logger *Logger) *Service {
    return &Service{db: db, logger: logger}
}
>>>>>>> feature/logging

// Решение - объединить изменения
func NewService(logger *Logger) *Service {
    return &Service{db: db, cache: redis, logger: logger}
}
```

---

## 5. Best practices

### Мелкие атомарные коммиты
```bash
# ХОРОШО - отдельные логические изменения
git add auth.go
git commit -m "feat: implement JWT authentication"

git add auth_test.go  
git commit -m "test: add JWT validation tests"

git add middleware.go
git commit -m "feat: add auth middleware"

# ПЛОХО - все в одном коммите
git add .
git commit -m "add user auth and tests and middleware"
```

### Регулярные пуши
```bash
# Работа в feature ветке
git checkout -b feat/payment-service
# ... пишем код ...

# Регулярные коммиты и пуши
git add .
git commit -m "feat: add payment processing"
git push origin feat/payment-service

git add .
git commit -m "test: add payment tests" 
git push origin feat/payment-service
```

---

## 🎯 Практика: Типичный рабочий день

### Утро - начало работы
```bash
# Получить последние изменения
git checkout main
git pull origin main

# Создать feature ветку
git checkout -b feat/user-profile

# Начать работу...
```

### День - разработка
```bash
# Регулярные коммиты
git add handlers/user.go
git commit -m "feat: add get user profile endpoint"

git add models/user.go
git commit -m "refactor: update user model fields"

git add tests/user_test.go  
git commit -m "test: add user profile tests"
```

### Вечер - завершение
```bash
# Пушим изменения
git push -u origin feat/user-profile

# Создаем Pull Request через GitHub
# Ждем code review...

# После approval - мерджим
git checkout main
git pull origin main
git merge --no-ff feat/user-profile
git push origin main

# Удаляем feature ветку
git branch -d feat/user-profile
git push origin --delete feat/user-profile
```

---

## ❓ Частые проблемы и решения

### Отмена изменений
```bash
# Отмена незакоммиченных изменений
git checkout -- file.go
git reset --hard HEAD

# Отмена последнего коммита
git reset --soft HEAD~1

# Отмена пуша (осторожно!)
git reset --hard HEAD~1
git push --force origin main
```

### Потерянные коммиты
```bash
# Поиск потерянных коммитов
git reflog
git checkout <commit-hash>
git branch recovery-branch
```

---

## 🏠 Домашнее задание

**Задача 1:** Создайте репозиторий для вашего проекта
**Задача 2:** Реализуйте feature branch workflow для новой функции
**Задача 3:** Создайте Pull Request и проведите code review с напарником

---

Теперь вы готовы к **Деплой приложения**! 🚀