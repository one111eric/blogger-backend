package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
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
	query := `
		SELECT id, title, author, content, postTime, editTime FROM posts
	`
	rows, err := d.DB.Query(query)
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

func (d *Database) GetPostById(id int) (*models.Post, error) {
	query := `
		SELECT id, title, author, content, postTime, editTime 
		FROM posts WHERE id = ?
	`
	// Use QueryRow for a single result
	row := d.DB.QueryRow(query, id)

	var post models.Post
	err := row.Scan(&post.ID, &post.Title, &post.Author, &post.Content, &post.PostTime, &post.EditTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no post found with id: " + strconv.Itoa(id))
		}
		return nil, err
	}

	return &post, nil
}

func (d *Database) CreatePost(post *models.Post) (int, error) {
	post.PostTime = time.Now() // Set the current time for postTime
	query := `
		INSERT INTO posts (title, author, content, postTime)
		VALUES (?, ?, ?, ?)
	`
	result, err := d.DB.Exec(
		query,
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

func (d *Database) EditPost(post *models.Post) (int, error) {
	*post.EditTime = time.Now()
	query := `
		UPDATE posts SET content = ?, editTime = ? 
		WHERE id = ? 
	`
	result, err := d.DB.Exec(query, post.Content, *post.EditTime, post.ID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, errors.New("no rows affected, post not found")
	}
	return int(rowsAffected), nil
}

func (d *Database) DeletePost(id int) (int, error) {
	query := `
		DELETE from posts 
		WHERE id = ? 
	`
	result, err := d.DB.Exec(query, id)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, errors.New("no rows affected, post not found")
	}
	return int(rowsAffected), nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}
