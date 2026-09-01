#!/bin/bash

# =============================================================================
# URL Shortener — Интерактивный установщик
# =============================================================================

set -e

# --- Цвета и оформление ---
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# --- Утилиты ---
ask() {
    local prompt="$1"
    local default="$2"
    local result=""

    if [ -n "$default" ]; then
        read -rp "${prompt} [$default]: " result
        result="${result:-$default}"
    else
        read -rp "${prompt}: " result
    fi

    echo "$result"
}

check_dependency() {
    local cmd="$1"
    local name="$2"

    if ! command -v "$cmd" &> /dev/null; then
        echo -e "${RED}[ERROR]${NC} $name не найден. Установите его и повторите попытку."
        exit 1
    fi
}

print_separator() {
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# =============================================================================
# Заголовок
# =============================================================================
clear
print_separator
echo -e "  ${BOLD}${BLUE}URL Shortener — Установщик${NC}"
print_separator
echo ""

# =============================================================================
# 1. Проверка зависимостей
# =============================================================================
echo -e "${YELLOW}[1/4] Проверка зависимостей...${NC}"

check_dependency docker Docker

# Docker Compose распространяется как плагин ("docker compose", с пробелом),
# а не отдельная программа ("docker-compose") — проверяем именно так.
if ! docker compose version &> /dev/null; then
    echo -e "${RED}[ERROR]${NC} Docker Compose не найден. Установите его и повторите попытку."
    exit 1
fi

# Проверка запущен ли Docker daemon
if ! docker info &> /dev/null; then
    echo -e "${RED}[ERROR]${NC} Docker daemon не запущен. Запустите Docker и повторите попытку."
    exit 1
fi

echo -e "  ${GREEN}✓${NC} Docker: $(docker --version | grep -oP '\d+\.\d+\.\d+' | head -1)"
echo -e "  ${GREEN}✓${NC} Docker Compose: $(docker compose version | grep -oP '\d+\.\d+\.\d+' | head -1)"
echo ""

# =============================================================================
# 2. Сбор параметров от пользователя
# =============================================================================
echo -e "${YELLOW}[2/4] Конфигурация проекта${NC}"
print_separator
echo ""

# --- База данных ---
echo -e "${BOLD}📦 База данных (PostgreSQL)${NC}"
POSTGRES_USER=$(ask "Имя пользователя PostgreSQL" "shortener")
POSTGRES_PASSWORD=$(ask "Пароль PostgreSQL" "shortener")
POSTGRES_DB=$(ask "Имя базы данных" "shortener")
POSTGRES_PORT=$(ask "Порт PostgreSQL (0 для запрета внешнего доступа)" "")

# Если порт 0 или пустой — не публикуем внешний порт
if [ -z "$POSTGRES_PORT" ] || [ "$POSTGRES_PORT" = "0" ]; then
    POSTGRES_PORT=""
fi

echo ""

# --- Бэкенд ---
echo -e "${BOLD}⚙️  Бэкенд (Go)${NC}"
BACKEND_PORT=$(ask "Порт бэкенда (для прямых запросов к API в обход nginx, необязательно для работы сайта)" "8080")
echo ""

# --- Фронтенд ---
echo -e "${BOLD}🎨 Фронтенд (Vue.js)${NC}"
FRONTEND_PORT=$(ask "Порт фронтенда" "80")

# BASE_URL — это публичный адрес САЙТА (фронтенда), а не бэкенда: именно на нём
# nginx проксирует и /shorten, и переходы по коротким ссылкам. Подставляем
# порт фронтенда в дефолт, чтобы не приходилось вводить его вручную для
# нестандартных портов.
if [ "$FRONTEND_PORT" = "80" ]; then
    default_base_url="http://localhost"
else
    default_base_url="http://localhost:$FRONTEND_PORT"
fi
echo -e "  ${CYAN}ℹ${NC} BASE_URL — это публичный адрес сайта (без https://, пока не настроен реальный TLS-сертификат)."
BASE_URL=$(ask "Публичный адрес сайта" "$default_base_url")
echo ""

# --- Окружение ---
echo -e "${BOLD}🌍 Окружение${NC}"
APP_ENV=$(ask "Окружение (development / production)" "production")
echo ""

# =============================================================================
# 3. Генерация .env файла
# =============================================================================
echo -e "${YELLOW}[3/4] Генерация .env файла...${NC}"

ENV_FILE=".env"

{
    # Database
    echo "# Database configuration"
    echo "POSTGRES_USER=$POSTGRES_USER"
    echo "POSTGRES_PASSWORD=$POSTGRES_PASSWORD"
    echo "POSTGRES_DB=$POSTGRES_DB"
    if [ -n "$POSTGRES_PORT" ]; then
        echo "POSTGRES_PORT=$POSTGRES_PORT"
    fi

    # Backend
    echo ""
    echo "# Backend configuration"
    echo "BACKEND_PORT=$BACKEND_PORT"
    echo "BASE_URL=$BASE_URL"

    # Frontend
    echo ""
    echo "# Frontend configuration"
    echo "FRONTEND_PORT=$FRONTEND_PORT"

    # Environment
    echo ""
    echo "# Application mode (development or production)"
    echo "APP_ENV=$APP_ENV"
} > "$ENV_FILE"

echo -e "  ${GREEN}✓${NC} Файл $ENV_FILE создан."
echo ""

# =============================================================================
# 4. Запуск проекта
# =============================================================================
echo -e "${YELLOW}[4/4] Запуск проекта...${NC}"
print_separator
echo ""

docker compose up -d --build

if [ $? -eq 0 ]; then
    echo ""
    print_separator
    echo -e "  ${BOLD}${GREEN}✓ Проект успешно запущен!${NC}"
    print_separator
    echo ""
    echo -e "  ${BOLD}🌐 Фронтенд:${NC}  $BASE_URL"
    echo -e "  ${BOLD}⚙️  Бэкенд (напрямую, необязательно):${NC}  http://localhost:$BACKEND_PORT"
    echo ""
    echo -e "  ${CYAN}Логи:${NC} docker compose logs -f"
    echo -e "  ${CYAN}Остановка:${NC} docker compose down"
    echo ""
else
    echo ""
    echo -e "  ${RED}[ERROR]${NC} Не удалось запустить проект. Проверьте логи:"
    echo "  docker compose logs"
    exit 1
fi

print_separator