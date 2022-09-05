CREATE TABLE messages
(
    id         integer primary key autoincrement,
    text       text    not null,
    tg_id      integer not null,
    user_id    integer not null,
    is_bot     boolean not null,
    command    text,
    created_at datetime,
    updated_at datetime,
    FOREIGN KEY (user_id) REFERENCES users (tg_id)
);