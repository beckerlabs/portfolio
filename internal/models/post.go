package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID       int
	Title    string
	Content  string
	Category string
	Created  time.Time
	Tags     string
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content, category, tags string) (int, error) {
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

func (m *PostModel) Get(id int) (*Post, error) {
	stmt := `SELECT id, title, content, category, created, tags
	FROM posts
	WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	p := &Post{}

	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &p.Created, &p.Tags)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (m *PostModel) GetAll() ([]*Post, error) {
	stmt := `SELECT id, title, content, category, created, tags
	FROM posts
	ORDER BY created DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}

	for rows.Next() {
		p := &Post{}
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
