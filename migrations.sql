DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id       SERIAL  PRIMARY KEY,
  email    TEXT    UNIQUE       NOT NULL,
  password TEXT    NOT NULL
);