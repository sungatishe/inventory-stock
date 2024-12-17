CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,                 -- Уникальный идентификатор
    email VARCHAR(255) NOT NULL UNIQUE,   -- Email пользователя
    password VARCHAR(255) NOT NULL,       -- Хэшированный пароль
    role VARCHAR(50) DEFAULT 'user',      -- Роль пользователя (по умолчанию 'user')
    created_at TIMESTAMP DEFAULT NOW(),   -- Дата создания записи
    updated_at TIMESTAMP DEFAULT NOW()    -- Дата последнего обновления
);
