CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    email      TEXT UNIQUE NOT NULL,
    name       TEXT        NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
