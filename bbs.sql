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

CREATE TABLE messages (
 id INTEGER PRIMARY KEY,
 basename TEXT NOT NULL,
 subject TEXT NOT NULL,
    author TEXT NOT NULL,
    date TEXT NOT NULL,
    message TEXT NOT NULL,
    postto TEXT NOT NULL,
)