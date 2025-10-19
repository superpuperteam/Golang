package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type User struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	Email   string  `db:"email"`
	Balance float64 `db:"balance"`
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) InsertUser(user *User) error {
	const q = `INSERT INTO users (name, email, balance) VALUES ($1, $2, $3)`
	_, err := ur.db.Exec(q, user.Name, user.Email, user.Balance)
	return err
}

func (ur *UserRepo) GetUserByID(id int) (User, error) {
	const q = `SELECT * FROM users WHERE id = $1`
	var u User
	err := ur.db.Get(&u, q, id)
	if err != nil {
		log.Fatal(err)
		return User{}, err
	}
	return u, nil
}

func (ur *UserRepo) GetAllUsers() ([]User, error) {
	const q = `SELECT * FROM users`
	var users []User
	err := ur.db.Select(&users, q)
	if err != nil {
		log.Fatal(err)
		return []User{}, err
	}
	return users, nil
}

func (ur *UserRepo) TransferBalance(fromId int, toId int, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if fromId == toId {
		return fmt.Errorf("sender and receiver must differ")
	}

	tx, err := ur.db.BeginTxx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	res, err := tx.Exec(`
		UPDATE users
		SET balance = balance - $1
		WHERE id = $2 AND balance >= $1
	`, amount, fromId)
	if err != nil {
		return fmt.Errorf("debit sender: %w", err)
	}
	aff, _ := res.RowsAffected()
	if aff != 1 {
		return fmt.Errorf("sender not found or insufficient funds")
	}

	res, err = tx.Exec(`
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2
	`, amount, toId)
	if err != nil {
		return fmt.Errorf("credit receiver: %w", err)
	}
	aff, _ = res.RowsAffected()
	if aff != 1 {
		return fmt.Errorf("receiver not found")
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}
