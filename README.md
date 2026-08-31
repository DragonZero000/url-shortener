# URL Shortener

[![Docker](https://img.shields.io/badge/Docker-Production-ready-brightgreen)](https://www.docker.com/)
[![Go](https://img.shields.io/badge/Backend-Go%201.26-blue)](https://go.dev/)
[![Vue.js](https://img.shields.io/badge/Frontend-Vue.js%203-green)](https://vuejs.org/)

> Production-ready URL shortener service с контейнерной архитектурой и автоматической установкой.

---

## 📋 Содержание

- [Возможности](#-возможности)
- [Архитектура](#-архитектура)
- [Требования](#-требования)
- [Быстрый старт](#-быстрый-старт)
- [Ручная установка](#-ручная-установка)
- [Конфигурация](#-конфигурация)
- [Команды](#-команды)
- [API](#-api)
- [Структура проекта](#-структура-проекта)

---

## ✨ Возможности

| Компонент | Описание |
|-----------|----------|
| 🚀 **Go Backend** | Высокопроизводительный API на Go 1.26 с PostgreSQL |
| 🎨 **Vue.js Frontend** | Современный SPA интерфейс на Vue 3 + Vite |
| 🐳 **Docker Production** | Полная контейнеризация для быстрого развёртывания |
| 🔧 **Автоматическая установка** | Интерактивный скрипт `install.sh` с проверкой зависимостей |
| 🗄️ **PostgreSQL 16** | Надёжное хранение данных с healthcheck и volume |
| ⚡ **Healthchecks** | Автоматическая проверка готовности сервисов в Docker Compose |

---

## 🏗️ Архитектура

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   Frontend      │────▶│    Backend       │────▶│  PostgreSQL     │
│   (Nginx + Vue) │     │  (Go API)        │     │  (Postgres:16)  │
│   :80             │     │  :8080           │     │  :5432          │
└─────────────────┘     └──────────────────┘     └─────────────────┘
```

---

## 📦 Требования

| Компонент | Версия | Назначение |
|-----------|--------|------------|
| Docker | 20+ | Запуск контейнеров |
| Docker Compose | v2+ | Оркестрация сервисов |

> **Примечание:** `docker compose` (с дефисом) — это плагин Docker, а не отдельная утилита. Убедитесь, что используете правильную версию: `docker compose version`.

---

## 🚀 Быстрый старт

### Вариант 1: Автоматическая установка (рекомендуется)

```bash
# Запуск интерактивного установщика
chmod +x install.sh && ./install.sh
```

Скрипт автоматически:
- Проверит наличие Docker и Docker Compose
- Запросит параметры конфигурации
- Создаст `.env` файл
- Соберёт и запустит проект

### Вариант 2: Ручная установка

```bash
# 1. Копируем переменные окружения
cp .env-example .env

# 2. Редактируем .env под свои нужды

# 3. Запускаем проект
docker compose up -d --build
```

---

## ⚙️ Конфигурация

Все настройки проекта управляются через `.env` файл. Ниже описаны все доступные переменные:

### База данных (PostgreSQL)

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `POSTGRES_USER` | `shortener` | Имя пользователя БД |
| `POSTGRES_PASSWORD` | `shortener` | Пароль пользователя БД |
| `POSTGRES_DB` | `shortener` | Имя базы данных |
| `POSTGRES_PORT` | `5432` | Внешний порт (установите `0` для запрета внешнего доступа) |

### Backend (Go API)

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `BACKEND_PORT` | `8080` | Порт API сервера |
| `BASE_URL` | `http://localhost:8080` | Базовый URL бэкенда (используется для генерации коротких ссылок) |

### Frontend (Vue.js + Nginx)

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `FRONTEND_PORT` | `80` | Порт веб-сервера |
| `VITE_API_BASE_URL` | `http://localhost:8080` | URL API для фронтенда |
| `VITE_WEB_BASE_URL` | `http://localhost` | Базовый URL веб-приложения |

### Окружение

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `APP_ENV` | `production` | Режим приложения (`development` или `production`) |

---

## 🛠️ Команды

### Запуск и управление

```bash
# Запустить проект в фоновом режиме
docker compose up -d --build

# Остановить все сервисы
docker compose down

# Просмотр логов (с хвостом)
docker compose logs -f

# Перезапустить конкретный сервис
docker compose restart backend

# Очистка контейнеров и volumes (⚠️ удалит данные БД!)
docker compose down -v
```

### Проверка состояния

```bash
# Статус запущенных контейнеров
docker ps

# Healthcheck PostgreSQL
docker exec url-shortener-postgres-1 pg_isready -U shortener
```

---

## 📡 API

Backend предоставляет REST API для работы с короткими ссылками.

### Основные эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/api/shorten` | Создать новую короткую ссылку |
| `GET` | `/api/stats/{id}` | Получить статистику переходов |

---

## 📁 Структура проекта

```
url-shortener/
├── .env-example              # Шаблон переменных окружения
├── docker-compose.yml        # Конфигурация Docker Compose
├── install.sh                # Интерактивный установщик
│
├── backend/                  # Go API сервер
│   ├── Dockerfile            # Multi-stage сборка (golang → alpine)
│   ├── go.mod / go.sum       # Зависимости Go
│   ├── cmd/server/main.go    # Точка входа приложения
│   ├── internal/             # Бизнес-логика
│   ├── migrations/           # SQL миграции БД
│   └── postman/              # Коллекция для тестирования API
│
├── frontend/                 # Vue.js SPA + Nginx
│   ├── Dockerfile            # Сборка (bun → nginx:alpine)
│   ├── nginx.conf            # Конфигурация веб-сервера
│   ├── src/                  # Исходный код Vue компонентов
│   └── public/               # Статические ресурсы
│
└── README.md                 # Документация проекта
```

---

## 🐳 Production Deployment

### Рекомендации для продакшена:

1. **Измените пароли по умолчанию** в `.env` файле на надёжные значения.
2. **Отключите внешний порт PostgreSQL**, установив `POSTGRES_PORT=0`.
3. **Настройте HTTPS** — добавьте reverse proxy (Nginx/Traefik) перед фронтендом.
4. **Используйте внешние volumes** для хранения данных PostgreSQL вместо Docker volumes.
5. **Настройте логирование и мониторинг** (Prometheus + Grafana, ELK).

### Пример production конфигурации:

```bash
# .env (production)
POSTGRES_USER=prod_user
POSTGRES_PASSWORD=<strong-password-here>
POSTGRES_DB=url_shortener_db
POSTGRES_PORT=0                    # Запрещаем внешний доступ к БД
BACKEND_PORT=8080
BASE_URL=https://api.yourdomain.com
FRONTEND_PORT=443
VITE_API_BASE_URL=https://api.yourdomain.com
VITE_WEB_BASE_URL=https://yourdomain.com
APP_ENV=production
```
