CREATE TABLE categories
(
    id integer primary key autoincrement,
    name text,
    user_id integer not null,
    created_at datetime,
    updated_at datetime,
    FOREIGN KEY(user_id) REFERENCES users(id)
);
