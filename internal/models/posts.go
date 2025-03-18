package models

import (
	"database/sql"
	"errors"
	"time"
)

type PostsModelInterface interface {
	Insert(title, content, category, tags string) (int, error)
	Get(id int) (Posts, error)
	Latest() ([]Posts, error)
}

type Posts struct {
	ID                int
	Title             string
	Content           string
	Category          string
	Created           time.Time
	Tags              string
	StructuredContent []string
}

type PostsModel struct {
	DB *sql.DB
}

func (m *PostsModel) Insert(title, content, category, tags string) (int, error) {
	stmt := `INSERT INTO posts (title, content, category, created, tags)
	VALUES(?, ?, ?, UTC_TIMESTAMP(), ?)`

	result, err := m.DB.Exec(stmt, title, content, category, tags)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PostsModel) Get(id int) (Posts, error) {
	stmt := `SELECT postId, title, content, category, created, tags
	FROM posts
	WHERE postId = ?`

	row := m.DB.QueryRow(stmt, id)

	var p Posts

	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &p.Created, &p.Tags)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Posts{}, ErrNoRecord
		} else {
			return Posts{}, err
		}
	}

	return p, nil
}

func (m *PostsModel) Latest() ([]Posts, error) {
	stmt := `SELECT postId, title, content, category, created, tags
	FROM posts
	ORDER BY created DESC
	LIMIT 5`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Posts{}

	for rows.Next() {
		var p Posts

		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &p.Created, &p.Tags)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
