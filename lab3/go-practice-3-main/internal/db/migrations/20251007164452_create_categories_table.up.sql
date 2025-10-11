CREATE TABLE categories
(
    id      SERIAL PRIMARY KEY,
    name    TEXT NOT NULL,
    user_id INTEGER REFERENCES users (id),
    CONSTRAINT categories_user_name_uniq UNIQUE (user_id, name)
);