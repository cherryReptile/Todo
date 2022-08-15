CREATE TABLE users(
    id integer primary key autoincrement,
    name text not null unique,
    tg_id integer not null unique
);