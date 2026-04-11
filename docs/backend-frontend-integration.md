# Хитрый лис
## Актуальная интеграция backend и frontend

Версия документа: 2026-04-12

Этот документ описывает не целевой контракт "на будущее", а текущее рабочее
состояние backend, с которым фронтенд уже может интегрироваться.

Связанные материалы:

- Swagger UI: `http://localhost:8080/swagger`
- OpenAPI YAML: `http://localhost:8080/openapi.yaml`
- Исторический целевой документ: `docs/frontend-backend-guide.txt`

---

## 1. Что уже есть на backend

На текущий момент backend уже поддерживает:

- guest auth;
- register/login/refresh/me;
- создание игры;
- вход в игру;
- lobby snapshot;
- старт игры;
- game snapshot;
- игровой websocket;
- real-time broadcast обновлений всем подключенным игрокам;
- connection presence в `players[].connected`;
- ограничение: один пользователь может быть только в одной незавершенной игре
  (`waiting` или `active`).

Под незавершенной игрой здесь понимается игра со статусом:

- `waiting`
- `active`

После перехода игры в `finished` пользователь может создавать новую или входить
в другую.

---

## 2. Базовая модель взаимодействия frontend с backend

Правильная схема для клиента такая:

1. Пользователь получает access token.
2. Клиент создает игру или входит в уже существующую.
3. Клиент открывает lobby screen.
4. После старта игры клиент делает HTTP-запрос за snapshot.
5. Затем клиент открывает websocket.
6. Все игровые команды отправляются только по websocket.
7. Все изменения экрана берутся из серверного state.

Главный принцип:

- HTTP нужен для auth, lifecycle игры, initial load и recovery.
- WebSocket нужен для live updates и игровых команд.

---

## 3. Базовые URL и авторизация

HTTP base URL:

```text
http://localhost:8080
```

WebSocket URL:

```text
ws://localhost:8080/ws/games/{gameId}?token=<access_token>
```

Для защищенных HTTP endpoints нужен header:

```text
Authorization: Bearer <access_token>
```

Для websocket допускается:

- query param `token`
- или `Authorization: Bearer <access_token>`

На практике для фронта проще использовать query param `token`.

---

## 4. Auth flow

### 4.1. Guest login

```http
POST /api/v1/auth/guest
```

Response:

```json
{
  "user": {
    "id": "uuid",
    "username": "guest_123",
    "is_guest": true,
    "role": "player",
    "created_at": "2026-04-12T10:00:00Z",
    "updated_at": "2026-04-12T10:00:00Z"
  },
  "access_token": "jwt",
  "refresh_token": "jwt"
}
```

Что должен сделать frontend:

- сохранить `access_token`;
- сохранить `user.id`;
- использовать token во всех следующих запросах.

### 4.2. Register / login / refresh / me

Также доступны:

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `GET /api/v1/auth/me`

Для MVP фронта достаточно guest auth, если нет отдельной бизнес-задачи на
настоящие аккаунты.

---

## 5. Flow игры через HTTP

### 5.1. Создать игру

```http
POST /api/v1/games
Authorization: Bearer <token>
```

Body:

```json
{
  "mode": "online"
}
```

Response `201`:

```json
{
  "game": {
    "id": "uuid",
    "status": "waiting"
  },
  "player": {
    "user_id": "uuid",
    "seat": 0
  }
}
```

После успеха frontend должен:

1. сохранить `game.id`;
2. перейти на lobby screen;
3. сделать `GET /api/v1/games/{id}`.

### 5.2. Войти в игру

```http
POST /api/v1/games/{id}/join
Authorization: Bearer <token>
```

Body можно отправлять пустым:

```json
{}
```

Response `200`:

```json
{
  "game": {
    "id": "uuid",
    "status": "waiting"
  },
  "player": {
    "user_id": "uuid",
    "seat": 1
  }
}
```

Особенности:

- если пользователь уже состоит именно в этой игре, `join` идемпотентен;
- если пользователь уже в другой `waiting/active` игре, вернется `409`;
- если игра уже стартовала, тоже вернется `409`.

### 5.3. Получить lobby snapshot

```http
GET /api/v1/games/{id}
Authorization: Bearer <token>
```

Response `200`:

```json
{
  "game": {
    "id": "uuid",
    "status": "waiting",
    "players": [
      {
        "user_id": "uuid",
        "seat": 0,
        "display_name": "guest_123",
        "is_me": true
      }
    ],
    "can_start": false,
    "min_players": 2,
    "max_players": 4
  }
}
```

Что важно:

- `can_start` уже приходит с сервера;
- frontend не должен самостоятельно решать, можно ли стартовать;
- seat использовать как серверный порядок игроков.

### 5.4. Старт игры

```http
POST /api/v1/games/{id}/start
Authorization: Bearer <token>
```

Response `200`:

```json
{
  "game": {
    "id": "uuid",
    "status": "active"
  },
  "redirect": {
    "route": "/game/uuid"
  }
}
```

Кто может стартовать:

- только создатель игры;
- только если игроков минимум двое;
- только пока игра в статусе `waiting`.

### 5.5. Получить игровой snapshot

```http
GET /api/v1/games/{id}/state
Authorization: Bearer <token>
```

Response `200`:

```json
{
  "state": {
    "...": "GameView"
  }
}
```

Этот endpoint обязателен для:

- initial load game screen;
- page refresh;
- recovery после reconnect;
- fallback, если websocket временно недоступен.

---

## 6. Реальный Public Game State

Фронту нужно ориентироваться именно на эти поля. Это то, что backend реально
отдает сейчас.

Пример:

```json
{
  "id": "uuid",
  "status": "active",
  "phase": "choose_goal",
  "result": "none",
  "version": 1,
  "turn": 1,
  "activeSeat": 0,
  "me": {
    "userId": "uuid",
    "seat": 0,
    "name": "guest_123",
    "pawnCell": 0,
    "connected": true
  },
  "players": [
    {
      "userId": "uuid",
      "seat": 0,
      "name": "guest_123",
      "pawnCell": 0,
      "connected": true
    }
  ],
  "board": {
    "cells": [
      {
        "index": 0,
        "type": "start",
        "hasClue": false
      }
    ]
  },
  "fox": {
    "track": 0,
    "escapeAt": 15
  },
  "suspects": [
    {
      "id": "suspect_1",
      "revealed": false,
      "excluded": false
    }
  ],
  "clues": [
    {
      "id": "clue_1",
      "revealed": false,
      "boardCell": 3
    }
  ],
  "turnState": {
    "goal": {
      "set": false,
      "type": ""
    },
    "pending": ""
  },
  "availableActions": [
    "choose_goal",
    "accuse"
  ]
}
```

### Что важно знать фронту про state

- `availableActions` уже готовый список разрешенных действий;
- `players[].connected` приходит с websocket presence и обновляется при
  connect/disconnect;
- `me` уже вычислен для текущего пользователя;
- state надо полностью заменять каждым новым update, а не пытаться патчить
  локально;
- не надо опираться на скрытые серверные поля.

### Чего фронт не должен использовать

Нельзя строить UI-логику на внутренних серверных данных, даже если они случайно
попадут в state в будущем:

- секретный culprit;
- внутренние флаги дедукции;
- сырой JSON из БД;
- внутренние repo/model поля.

---

## 7. WebSocket интеграция

### 7.1. Подключение

URL:

```text
ws://localhost:8080/ws/games/{gameId}?token=<access_token>
```

Подключение допускается только если:

- токен валиден;
- пользователь состоит в этой игре;
- игра существует.

После успешного подключения:

1. пользователь добавляется в room;
2. сервер делает broadcast актуального state;
3. `connected` у подключенного игрока становится `true`;
4. остальные игроки тоже получают обновленный state.

### 7.2. Формат клиентского сообщения

```json
{
  "id": "req-1",
  "type": "command",
  "command": "choose_goal",
  "payload": {
    "goal": "clue"
  }
}
```

Обязательные поля:

- `id`
- `type`
- `command`
- `payload`

`type` всегда должен быть:

```json
"command"
```

### 7.3. Формат update от сервера

```json
{
  "id": "req-1",
  "type": "update",
  "payload": {
    "state": {
      "...": "GameView"
    },
    "events": [
      {
        "type": "goal_chosen",
        "data": {
          "goal": "clue"
        }
      }
    ]
  }
}
```

Семантика:

- `state` важнее, чем `events`;
- `events` нужны для лога и анимаций;
- `id` связывает ответ с пользовательским действием;
- если событие не связано с локальной командой, `id` может быть пустым.

### 7.4. Формат error от сервера

```json
{
  "id": "req-1",
  "type": "error",
  "payload": {
    "code": "not_your_turn",
    "message": "It is not your turn."
  }
}
```

Правильное поведение frontend:

- не закрывать websocket из-за такого сообщения;
- не мутировать state локально;
- показать пользователю понятную ошибку;
- оставить UI в согласованном серверном состоянии.

---

## 8. Игровые команды, которые уже можно слать

### choose_goal

```json
{
  "id": "req-1",
  "type": "command",
  "command": "choose_goal",
  "payload": {
    "goal": "clue"
  }
}
```

`goal`:

- `clue`
- `suspect`

### roll_auto

```json
{
  "id": "req-2",
  "type": "command",
  "command": "roll_auto",
  "payload": {}
}
```

### move_pawn

```json
{
  "id": "req-3",
  "type": "command",
  "command": "move_pawn",
  "payload": {
    "steps": 3
  }
}
```

### take_clue

```json
{
  "id": "req-4",
  "type": "command",
  "command": "take_clue",
  "payload": {}
}
```

### reveal_suspects

```json
{
  "id": "req-5",
  "type": "command",
  "command": "reveal_suspects",
  "payload": {
    "suspectIds": ["suspect_1", "suspect_2"]
  }
}
```

### end_turn

```json
{
  "id": "req-6",
  "type": "command",
  "command": "end_turn",
  "payload": {}
}
```

### accuse

```json
{
  "id": "req-7",
  "type": "command",
  "command": "accuse",
  "payload": {
    "suspectId": "suspect_1"
  }
}
```

---

## 9. Ошибки, которые должен корректно обрабатывать frontend

### HTTP

Типовые сценарии:

- `401`: нет токена или он битый;
- `403`: нет доступа к игре;
- `404`: игра не найдена;
- `409`: конфликт бизнес-логики;
- `500`: общая серверная ошибка.

Типовые причины `409`:

- игра уже стартовала;
- игра заполнена;
- недостаточно игроков для старта;
- пользователь уже в другой незавершенной игре.

### WebSocket

Минимальный список актуальных кодов:

- `not_your_turn`
- `invalid_phase`
- `game_not_active`
- `game_finished`
- `goal_already_set`
- `goal_not_set`
- `no_pending_action`
- `all_clues_collected`
- `suspect_not_revealed`
- `suspect_excluded`
- `forbidden`
- `game_not_found`
- `internal_error`

Фронтенд должен:

- показывать короткое уведомление;
- не ронять websocket session из-за `error` сообщения;
- не переписывать state локальными догадками после ошибки.

---

## 10. Reconnect и жизненный цикл сокета

Актуальное ожидаемое поведение frontend:

1. открыть game screen;
2. получить snapshot по HTTP;
3. открыть websocket;
4. при каждом `update` полностью заменять локальный `state`;
5. при disconnect перейти в `reconnecting`;
6. пробовать переподключение с backoff;
7. после успешного reconnect снова получить server state и продолжить работу.

Важно:

- backend больше не закрывает websocket из-за общего HTTP timeout;
- backend больше не режет сокет просто из-за минуты бездействия;
- после disconnect presence игроков пересчитывается и рассылается заново.

Что фронтенд не должен делать:

- не воспроизводить потерянные команды локально;
- не считать, что pending локальное действие точно применилось;
- не "достраивать" пропущенные события без нового snapshot.

---

## 11. Рекомендуемая структура frontend store

Минимально разумная схема:

- `authStore`
  - `token`
  - `user`
  - `sessionStatus`
- `lobbyStore`
  - `gameId`
  - `lobby`
  - `loading`
  - `error`
- `gameStore`
  - `state`
  - `connectionStatus`
  - `pendingRequestIds`
  - `eventLog`
  - `lastError`
- `apiClient`
- `wsClient`

Компоненты не должны:

- открывать websocket напрямую внутри template/ui-логики;
- решать бизнес-правила вместо сервера;
- мутировать `state` вручную после нажатия кнопки.

---

## 12. Практический flow интеграции frontend

### 12.1. Auth screen

Экран должен уметь:

- создать guest session;
- сохранить `access_token`;
- показать ошибку при неудаче;
- перевести пользователя на lobby screen.

### 12.2. Lobby screen

Экран должен уметь:

- создать игру;
- войти в игру по `gameId`;
- периодически или по действию получать lobby snapshot;
- рендерить список игроков;
- показывать кнопку `Start`, только если `can_start == true`.

### 12.3. Game screen

Экран должен уметь:

- получить `GET /state`;
- открыть websocket;
- показывать:
  - текущую фазу;
  - чей ход;
  - игроков;
  - presence;
  - доступные действия;
  - ошибки;
  - состояние соединения.

### 12.4. Поведение кнопок

UI должен ориентироваться на `availableActions`.

Если действия нет в `availableActions`, кнопку надо:

- скрыть;
- или сделать disabled.

Не нужно самому вычислять:

- можно ли сейчас бросать;
- можно ли завершать ход;
- можно ли брать улику;
- можно ли обвинять.

---

## 13. Практический flow ручного тестирования

### HTTP

Порядок:

1. `POST /api/v1/auth/guest` для player1
2. `POST /api/v1/auth/guest` для player2
3. `POST /api/v1/games` под player1
4. `POST /api/v1/games/{id}/join` под player2
5. `GET /api/v1/games/{id}` у обоих
6. `POST /api/v1/games/{id}/start` под player1
7. `GET /api/v1/games/{id}/state` у обоих

### WebSocket

Подключения:

```text
ws://localhost:8080/ws/games/<GAME_ID>?token=<PLAYER1_TOKEN>
ws://localhost:8080/ws/games/<GAME_ID>?token=<PLAYER2_TOKEN>
```

Проверки:

- initial update приходит сразу;
- `players[].connected` меняется при connect/disconnect;
- оба клиента получают одинаковые `update` после команды;
- неактивный игрок получает `not_your_turn`;
- reconnect возвращает консистентный state.

---

## 14. Что стоит сделать дальше по backend

Ниже список следующих задач по приоритету.

### Приоритет 1. Довести публичный API до стабильного контракта

Что сделать:

- зафиксировать публичные DTO отдельно от внутренних domain-структур;
- перестать отдавать фронту поля, которые завязаны на внутреннюю форму Go state;
- добавить явные response DTO для websocket update/error;
- сделать контрактную совместимость между `/state` и websocket.

Зачем:

- фронт перестанет зависеть от случайной сериализации доменных структур;
- любые изменения домена будут меньше ломать клиент.

### Приоритет 2. Добавить lobby live-updates

Сейчас lobby живет только через HTTP snapshot. Следующий логичный шаг:

- websocket или SSE для lobby;
- live-обновление списка игроков;
- live-обновление `can_start`;
- автоматический переход в game screen после старта.

### Приоритет 3. Сделать нормальный game presence lifecycle

Сейчас presence считается по активным websocket connections. Дальше можно
улучшить:

- heartbeat/ping-pong;
- last_seen;
- более явные presence events;
- корректная обработка нескольких вкладок одного пользователя.

### Приоритет 4. Ввести коды и JSON-формат для HTTP ошибок

Сейчас часть HTTP ошибок возвращается plain text. Для фронта удобнее:

```json
{
  "code": "already_in_another_game",
  "message": "User is already in another game."
}
```

Это упростит:

- локализацию;
- маппинг ошибок в UI;
- обработку `409` и `403` без парсинга текста.

### Приоритет 5. Добавить интеграционные тесты backend

Нужно покрыть:

- auth flow;
- create/join/start/state;
- ограничение "один пользователь = одна незавершенная игра";
- websocket broadcast двум игрокам;
- reconnect/presence;
- negative cases.

### Приоритет 6. Довести игровую доменную модель

Сейчас игра уже работает, но дальше полезно:

- вынести setup/seed данных в отдельный модуль;
- сделать более прозрачную конфигурацию board/suspects/clues;
- отделить public view от internal secret state;
- добавить устойчивые game fixtures для тестов.

### Приоритет 7. Добавить observability

Полезно добавить:

- structured logs по request id и game id;
- websocket connect/disconnect metrics;
- события команд и ошибок;
- latency по HTTP и WS command handling.

### Приоритет 8. Подготовить production-ready auth/game rules

Если проект пойдет дальше MVP:

- revoke/expire стратегии токенов;
- invite code или приватные комнаты;
- reconnect resume policy;
- завершение/cleanup старых waiting games;
- защита от гонок при одновременном join.

---

## 15. Краткая памятка фронтендеру

- Всегда доверяй серверному `state`.
- После каждого `update` полностью заменяй локальный state.
- Все игровые команды шли только по websocket.
- Если сокет умер, не паникуй: reconnect + `/state`.
- Не вычисляй бизнес-правила на клиенте, если уже есть `availableActions`.
- Не храни секреты игры на клиенте.

