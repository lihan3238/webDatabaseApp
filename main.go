package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Database configuration
const (
	dbUser     = "root"
	dbPassword = "123456"
	dbName     = "movie"
	dbHost     = "192.168.50.128:3306"
)

// Movie represents a movie record
type Movie struct {
	Title  string
	Rating float64
	Tags   []string
}

// User represents a user record
type User struct {
	UserID int
	Gender string
	Name   string
}

// connectDB initializes a connection to the database
// connectDB initializes a connection to the database
func connectDB() (*sql.DB, error) {
	// Replace 'localhost' with your database server's IP address or hostname
	// and '3306' with your MySQL server's port if it's not the default.
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName))
	if err != nil {
		println("Error connecting to the database: ", err.Error())
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			t, _ := template.ParseFiles("search.html")
			t.Execute(w, nil)
		} else if r.Method == "POST" {
			// Handle search here
		}
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		// Here you will need to parse the input from the user and then execute the appropriate SQL queries.
		// I will provide an example of searching movies by user ID.
		userID := r.URL.Query().Get("userid")
		if userID != "" {
			// Execute query
			rows, err := db.Query(`
				SELECT m.title, r.rating
				FROM ratings r
				INNER JOIN movies m ON m.movieId = r.movieId
				WHERE r.userId = ?
				ORDER BY r.timestamp DESC`, userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var movies []Movie
			for rows.Next() {
				var movie Movie
				if err := rows.Scan(&movie.Title, &movie.Rating); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// Fetch top 3 tags for each movie
				tagRows, err := db.Query(`
					SELECT gt.tag
					FROM genomescores gs
					INNER JOIN genometags gt ON gt.tagId = gs.tagId
					WHERE gs.movieId = (SELECT movieId FROM movies WHERE title = ?)
					ORDER BY gs.relevance DESC
					LIMIT 3`, movie.Title)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer tagRows.Close()

				for tagRows.Next() {
					var tag string
					if err := tagRows.Scan(&tag); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					movie.Tags = append(movie.Tags, tag)
				}

				movies = append(movies, movie)
			}

			// Convert the movies slice to JSON
			jsonData, err := json.Marshal(movies)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Set the content type to application/json
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
