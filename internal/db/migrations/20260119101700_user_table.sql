-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id       UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    email    VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255)        NOT NULL
);

CREATE TABLE IF NOT EXISTS message
(
    id         UUID                  DEFAULT gen_random_uuid() PRIMARY KEY,
    type       VARCHAR(20)  NOT NULL DEFAULT 'message',
    msg        VARCHAR(255) NOT NULL,
    channel    VARCHAR(255) NOT NULL,
    username   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_message_channel ON message(channel);
CREATE INDEX IF NOT EXISTS idx_message_created_at ON message(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_message_channel_created_at ON message(channel, created_at DESC);


INSERT INTO message (type, msg, channel, username, created_at)
VALUES ('system', 'Добро пожаловать в общий чат!', 'general', 'System', now() - interval '1 hour'),
       ('message', 'Привет всем! Как дела?', 'general', 'Алексей', now() - interval '45 minutes'),
       ('message', 'Всем привет! Все отлично!', 'general', 'Мария', now() - interval '30 minutes'),
       ('system', 'Добро пожаловать в чат помощи', 'help', 'System', now() - interval '2 hours'),
       ('message', 'Как пользоваться этим чатом?', 'help', 'Новичок', now() - interval '1 hour');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_message_channel;
DROP INDEX IF EXISTS idx_message_created_at;
DROP INDEX IF EXISTS idx_message_channel_created_at;
DROP TABLE IF EXISTS message;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
