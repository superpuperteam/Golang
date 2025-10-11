CREATE TABLE expenses
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER        NOT NULL REFERENCES users (id),
    category_id INTEGER        NOT NULL REFERENCES categories (id),
    amount      NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    currency    CHAR(3)        NOT NULL,
    spent_at    TIMESTAMPTZ    NOT NULL,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),
    note        TEXT
);

CREATE INDEX expenses_user_id_idx ON expenses (user_id);
CREATE INDEX expenses_user_spent_at_idx ON expenses (user_id, spent_at);
