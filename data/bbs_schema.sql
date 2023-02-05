CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    level INTEGER NOT NULL,
    linefeeds INTEGER NOT NULL,
    hotkeys INTEGER NOT NULL,
    active INTEGER NOT NULL,
    clearscreen INTEGER NOT NULL
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

CREATE TABLE fileareas (
    id INTEGER PRIMARY KEY,
    areaname TEXT NOT NULL,
    filename TEXT NOT NULL,
    description TEXT NOT NULL,
    uploadedby TEXT NOT NULL,
    date TEXT NOT NULL,
    size INTEGER NOT NULL,
)

CREATE TABLE sessions (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  active TINYINT(1) DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);