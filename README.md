# 🦊 Коварный лис / Outfoxed 

<div align="center">

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Angular](https://img.shields.io/badge/Angular-frontend-DD0031?logo=angular&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-compose-2496ED?logo=docker&logoColor=white)
![Status](https://img.shields.io/badge/status-MVP%20in%20progress-orange)
![License](https://img.shields.io/badge/license-Educational-blueviolet)

**Многопользовательская веб-адаптация настольной игры про поиск хитрого лиса.**  
Игроки вместе собирают улики, открывают подозреваемых и пытаются вычислить виновника раньше, чем лис доберётся до норки.

</div>

---

## ✨ О проекте

**Коварный лис** — это учебный full-stack проект.

Проект сделан как многопользовательская клиент-серверная игра.

---

## 🎮 Что уже умеет проект

- регистрация и авторизация игроков
- создание публичных и приватных комнат
- вход в игру по id и по коду
- lobby с живым обновлением состава игроков
- старт партии из комнаты
- игровое поле **16×16**
- пошаговая игра с фазами хода
- real-time обновления через **WebSocket**
- ручной бросок кубиков с выбором оставляемых кубиков
- перемещение по клеткам поля
- работа с уликами и подозреваемыми
- автообработка таймаута хода
- reconnect игроков в активную игру
- возможность вернуться в свою незавершённую партию
- базовая логика “бот подхватывает игрока”, если тот пропал из игры

---

## 🧩 Игровая идея

Игроки действуют как команда и пытаются найти виновника.

### Цель
Найти хитрого лиса раньше, чем он доберётся до конца следа.

### Как проходит ход
1. Игрок выбирает цель хода:
    - искать улику
    - проверять подозреваемых
2. Бросает кубики
3. При удачном результате:
    - двигается к улике
    - или открывает подозреваемых
4. Завершает ход, и очередь переходит следующему игроку

### Что важно
- в игре участвуют **16 подозреваемых**
- у каждого подозреваемого **3 предмета**
- в каждой партии используется **6 случайных активных предметов**
- на поле размещаются **6 улик**
- улики распределяются по карте равномерно, чтобы не скучиваться в одном месте

---

## ⚙️ Технологии

### Backend
- Go
- Chi
- PostgreSQL
- WebSocket
- Docker

### Frontend
- Angular
- RxJS
- Signals / reactive state
- WebSocket client

### Infra
- Docker Compose
- Caddy / reverse proxy
- Cloudflare Tunnel для демо-доступа наружу

---

## 📦 Структура проекта

```text
.
├── backend/
│   ├── cmd/
│   ├── config/
│   ├── internal/
│   │   ├── app/
│   │   ├── modules/
│   │   │   ├── auth/
│   │   │   └── game/
│   │   ├── transport/
│   │   └── repo/
│   └── pkg/
├── frontend/
├── deploy/
├── docker-compose.yml
└── README.md
```

> Точная структура может немного отличаться в зависимости от текущей ветки разработки.

---

## 🚀 Быстрый старт

### 1. Клонировать репозиторий

```bash
git clone <your-repo-url>
cd <your-project-folder>
```

### 2. Подготовить `.env`

Пример:

```env
HTTP_ADDR=:8080

POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=postgres

DB_DATA_SOURCE=postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable

JWT_SECRET=super-strong-secret

IMAGE_REPOSITORY=ghcr.io/blxssy/outfoxed
```

### 3. Запустить проект

```bash
docker compose up -d --build 
```

### 4. Открыть в браузере

- фронт: `http://localhost:80`
- backend API через proxy: `http://localhost:8080/api/...`
- WebSocket: `ws://localhost:8080/ws/...`

---

## 🖥 Локальная разработка

### Backend
```bash
cd backend
go run ./cmd/app/main.go
```

### Frontend
```bash
cd frontend
npm install
ng serve
```

Если фронт запускается отдельно через dev server, убедись, что в dev-конфиге разрешён нужный origin для API.

---

## 🌍 Деплой

Для демо проект можно развернуть на небольшом VPS:

- 1 vCPU
- 2 GB RAM
- 20–40 GB SSD

---

## 🤝 Для чего этот проект

Проект полезен как пример:

- проектирования stateful multiplayer backend
- real-time взаимодействия через WebSocket
- работы с транзакциями и консистентностью
- синхронизации клиента и сервера
- организации игрового домена на Go

---

## 📚 Примечание

Этот репозиторий создан в учебных целях.  
Он развивается как дипломный проект и постепенно доводится до полноценного демонстрационного состояния.

---

<div align="center">

**Если тебе понравился проект — поставь ⭐**

</div>
