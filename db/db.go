package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/one111eric/blogger-backend/models"
)

type Database struct {
	DB *sql.DB
}

// Initialize initializes the database and creates tables if necessary
func Initialize(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
                CREATE TABLE IF NOT EXISTS posts (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        title TEXT,
                        author TEXT,
                        content TEXT,
                        postTime DATETIME,
                        editTime DATETIME
                )
        `)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) GetPosts() ([]models.Post, error) {
	rows, err := d.DB.Query("SELECT id, title, author, content, postTime, editTime FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Author, &p.Content, &p.PostTime, &p.EditTime); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (d *Database) CreatePost(post *models.Post) (int, error) {
	post.PostTime = time.Now() // Set the current time for postTime
	result, err := d.DB.Exec(
		"INSERT INTO posts (title, author, content, postTime) VALUES (?, ?, ?, ?)",
		post.Title, post.Author, post.Content, post.PostTime,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}
