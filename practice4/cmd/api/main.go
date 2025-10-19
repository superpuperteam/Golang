package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ssss1131/go-practice4/internal/storage"
	"log"
	"time"
)

func openDB(dsn string) *sqlx.DB {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db
}

func main() {
	db := openDB("postgres://user:password@localhost:5430/mydatabase?sslmode=disable")
	repo := storage.NewUserRepo(db)

	user1 := storage.User{Email: "user1@gmail.com", Name: "user1", Balance: 500}
	err := repo.InsertUser(&user1)
	if err != nil {
		log.Fatal(err)
		return
	}

	user2 := storage.User{Email: "user2@gmail.com", Name: "user2", Balance: 500}
	err = repo.InsertUser(&user2)
	if err != nil {
		log.Fatal(err)
		return
	}
	users, err := repo.GetAllUsers()
	if err != nil {
		log.Fatal(err)
		return
	}
	for i := range users {
		fmt.Println(users[i].Email, users[i].Name)
	}

	err = repo.TransferBalance(1, 2, 100)
	if err != nil {
		log.Fatal(err)
		return
	}

	user, err := repo.GetUserByID(1)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(user.Email, user.Name, user.Balance)
}
