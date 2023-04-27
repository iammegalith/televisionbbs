CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    level INTEGER NOT NULL,
    active INTEGER NOT NULL,
    created TIMESTAMP NOT NULL,
    last_login TIMESTAMP
);
