## Описание таблиц

### auth_user
Хранит информацию о пользователях системы: логин, почта, пароль, описание, баланс и роль.

### transaction
Фиксирует все транзакции пользователя: уникальный ID, внешний ID транзакции, сумму, тип, статус и дату.

### payment
Отражает платежи пользователей, включая сумму, баланс после платежа, статус и дату создания.

### banner
Хранит рекламные баннеры, их описание, содержимое, ссылку, баланс, статус и владельца.

---

## Функциональные зависимости

### Relation: auth_user

{id} -> {username, email, password, description, balance, created_at, updated_at, role}

### Relation: transaction

{id} -> {transaction_id, user_id, amount, type, status, created_at}

{transaction_id} -> {id, user_id, amount, type, status, created_at}

### Relation: payment

{id} -> {owner_id, amount, created_at, status, balance}

### Relation: banner

{id} -> {title, description, balance, content, link, deleted, status, owner_id}


---

## Нормализация

### Первая нормальная форма (1НФ)
Все атрибуты атомарны, отсутствуют повторяющиеся группы. Нет составных типов данных (array, json и т.д.). Все отношения удовлетворяют 1НФ.

### Вторая нормальная форма (2НФ)
В 2НФ отношение должно быть в 1НФ, и все неключевые атрибуты должны зависеть от всего ключа. Во всех таблицах ключи простые (id), и все неключевые атрибуты зависят от них полностью.

### Третья нормальная форма (3НФ)
Требует, чтобы все неключевые атрибуты зависели только от ключа и не зависели друг от друга (транзитивных зависимостей нет). Все зависимости в таблицах этому соответствуют.

### Нормальная форма Бойса-Кодда (НФБК)
Для всех функциональных зависимостей X -> Y, X является суперключом. В каждой таблице определён один кандидатный ключ (id), от которого функционально зависят все остальные поля.

Следовательно, все отношения находятся в НФБК.

---

## ER-диаграмма (mermaid)

```mermaid
erDiagram
    auth_user {
        integer id
        text username
        text email
        bytea password
        text description
        numeric balance
        timestamp created_at
        timestamp updated_at
        integer role
    }

    banner {
        integer id
        integer owner_id
        text title
        text description
        integer balance
        text content
        text link
        boolean deleted
        integer status
        numeric max_price
    }

    payment {
        integer id
        integer owner_id
        integer amount
        timestamp created_at
        integer status
    }

    transaction {
        integer id
        text transaction_id
        integer user_id
        double precision amount
        text type
        text status
        timestamp created_at
    }

    auth_user ||--o| banner : "owns"
    auth_user ||--o| payment : "owns"
    auth_user ||--o| transaction : "made"
    banner }|--|| auth_user : "belongs to"
    payment }|--|| auth_user : "belongs to"
    transaction }|--|| auth_user : "belongs to"
```

