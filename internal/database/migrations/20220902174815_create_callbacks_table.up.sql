CREATE TABLE callbacks
(
    id         integer primary key autoincrement,
    json       text    not null,
    tg_id      integer not null,
    user_id    integer not null,
    created_at datetime,
    updated_at datetime,
    FOREIGN KEY (user_id) REFERENCES users (tg_id)
);