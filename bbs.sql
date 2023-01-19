CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    level INTEGER NOT NULL,
    linefeeds INTEGER NOT NULL,
    translation INTEGER NOT NULL,
    active INTEGER NOT NULL,
    clearscreen INTEGER NOT NULL,
);
