package services

import (
	httpcommon "chi-mysql-boilerplate/internal/domain/http_common"
	"chi-mysql-boilerplate/internal/domain/models"
	"context"
	"database/sql"
	"time"
)

type PostService struct {
	db *sql.DB
}

func NewPostService(db *sql.DB) *PostService {
	return &PostService{db: db}
}

func (p *PostService) Create(userId uint64, content string) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpcommon.DbConstants.Timeout)
	defer cancel()

	query := `
		INSERT INTO posts (content, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	creationTime := time.Now()
	result, err := p.db.ExecContext(ctx, query, content, userId, creationTime, creationTime)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	newPost := models.Post{
		ID:        uint64(id),
		Content:   content,
		UserID:    userId,
		CreatedAt: creationTime,
		UpdatedAt: creationTime,
	}

	return &newPost, nil
}

func (p *PostService) GetAll() ([]*models.PostResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpcommon.DbConstants.Timeout)
	defer cancel()

	query := `
		SELECT posts.id, content, user_id, username ,created_at, updated_at
		FROM posts JOIN users
		ON posts.user_id = users.id
	`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var posts []*models.PostResponse
	for rows.Next() {
		var post models.PostResponse
		if err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.UserID,
			&post.UserName,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (p *PostService) GetByUserId(userId uint64) ([]*models.PostResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpcommon.DbConstants.Timeout)
	defer cancel()

	query := `
		SELECT posts.id, content, user_id, username, created_at, updated_at
		FROM posts JOIN users
		ON posts.user_id = users.id
		WHERE user_id = ?
	`
	rows, err := p.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	var posts []*models.PostResponse
	for rows.Next() {
		var post models.PostResponse
		if err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.UserID,
			&post.UserName,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (p *PostService) GetById(id uint64) (*models.PostResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpcommon.DbConstants.Timeout)
	defer cancel()

	query := `
		SELECT posts.id, content, user_id, username, created_at, updated_at
		FROM posts JOIN users ON posts.user_id = users.id
		WHERE posts.id = ?
	`
	row := p.db.QueryRowContext(ctx, query, id)

	var post models.PostResponse
	if err := row.Scan(
		&post.ID,
		&post.Content,
		&post.UserID,
		&post.UserName,
		&post.CreatedAt,
		&post.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *PostService) UpdateById(id uint64, updateContent string) error {
	ctx, cancel := context.WithTimeout(context.Background(), httpcommon.DbConstants.Timeout)
	defer cancel()

	query := `
		UPDATE posts
		SET
			content = ?,
			updated_at = ?
		WHERE id = ?
	`
	_, err := p.db.ExecContext(ctx, query, updateContent, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostService) DeleteById(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), httpcommon.DbConstants.Timeout)
	defer cancel()

	query := `
		DELETE FROM posts
		WHERE id = ?
	`
	_, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
