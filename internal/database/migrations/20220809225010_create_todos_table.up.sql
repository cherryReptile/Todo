CREATE TABLE todos
(
    id          integer primary key autoincrement,
    name        text    not null,
    category_id integer not null,
    created_at  datetime,
    updated_at  datetime,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);