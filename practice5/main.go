package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Movie struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Year       int    `json:"year"`
	ActorCount int    `json:"actor_count"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:pass@localhost:5433/moviesdb?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		handleGetMovies(w, r, db)
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func handleGetMovies(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	q := r.URL.Query()

	var yearMinPtr *int
	var yearMaxPtr *int
	limit := 50
	offset := 0

	if s := strings.TrimSpace(q.Get("year_min")); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			http.Error(w, "invalid year_min", 400)
			return
		}
		yearMinPtr = &v
	}
	if s := strings.TrimSpace(q.Get("year_max")); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			http.Error(w, "invalid year_max", 400)
			return
		}
		yearMaxPtr = &v
	}
	if s := strings.TrimSpace(q.Get("limit")); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v < 0 {
			http.Error(w, "invalid limit", 400)
			return
		}
		limit = v
	}
	if s := strings.TrimSpace(q.Get("offset")); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v < 0 {
			http.Error(w, "invalid offset", 400)
			return
		}
		offset = v
	}

	sqlStr, args := buildMoviesQuery(yearMinPtr, yearMaxPtr, limit, offset)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	start := time.Now()
	rows, err := db.QueryContext(ctx, sqlStr, args...)
	elapsed := time.Since(start)

	w.Header().Set("X-Query-Time", fmt.Sprintf("%dms", elapsed.Milliseconds()))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err != nil {
		http.Error(w, fmt.Sprintf("query error: %v", err), 500)
		return
	}
	defer rows.Close()

	var out []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Year, &m.ActorCount); err != nil {
			http.Error(w, fmt.Sprintf("scan error: %v", err), 500)
			return
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("rows error: %v", err), 500)
		return
	}

	_ = json.NewEncoder(w).Encode(out)
}

func buildMoviesQuery(yearMin, yearMax *int, limit, offset int) (string, []any) {
	var sb strings.Builder
	var args []any
	i := 1

	sb.WriteString(`
SELECT
  m.id,
  m.title,
  m.year,
  COUNT(a.id) AS actor_count
FROM movies AS m
LEFT JOIN actors AS a ON a.movie_id = m.id
WHERE 1=1
`)

	if yearMin != nil {
		sb.WriteString(fmt.Sprintf("  AND m.year >= $%d\n", i))
		args = append(args, *yearMin)
		i++
	}
	if yearMax != nil {
		sb.WriteString(fmt.Sprintf("  AND m.year <= $%d\n", i))
		args = append(args, *yearMax)
		i++
	}

	sb.WriteString(`
GROUP BY m.id
ORDER BY m.year DESC, m.id ASC
`)

	if limit > 0 {
		sb.WriteString(fmt.Sprintf("LIMIT $%d\n", i))
		args = append(args, limit)
		i++
	} else {
		sb.WriteString("LIMIT 0\n")
	}
	if offset > 0 {
		sb.WriteString(fmt.Sprintf("OFFSET $%d\n", i))
		args = append(args, offset)
	}
	return sb.String(), args
}
