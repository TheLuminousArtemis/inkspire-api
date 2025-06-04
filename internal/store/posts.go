package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/lib/pq"
)

type PostUser struct {
	ID       int64
	Username string
}

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	Version   int       `json:"version"`
	User      PostUser  `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comment_count"`
}

type PostgresPostStore struct {
	db *sql.DB
}

func (s *PostgresPostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (title, content, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	row := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.UserID, pq.Array(post.Tags))
	err := row.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresPostStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `SELECT p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at, p.tags, p.version, u.id, u.username FROM posts p JOIN users u ON p.user_id = u.id WHERE p.id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	row := s.db.QueryRowContext(ctx, query, postID)
	post := &Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt, pq.Array(&post.Tags), &post.Version, &post.User.ID, &post.User.Username)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return post, nil
}

func (s *PostgresPostStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresPostStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts SET title = $1, content = $2, version = version+1 WHERE id = $3 AND VERSION=$4 RETURNING VERSION`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	log.Printf("Updating post ID %d with title=%s content=%s", post.ID, post.Title, post.Content)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *PostgresPostStore) GetUserFeed(ctx context.Context, userID int64, fq PagintatedFeedQuery) ([]PostWithMetadata, error) {
	query := `
	SELECT
	    p.id,
	    p.user_id,
	    p.title,
	    p.content,
	    p.created_at,
	    p.version,
	    p.tags,
	    u.username,
	    COUNT(c.id) AS comments_count
	FROM posts p
	LEFT JOIN comments c ON c.post_id = p.id
	JOIN users u ON p.user_id = u.id
	WHERE
	    (p.user_id = $1 OR p.user_id IN (
	        SELECT user_id FROM followers WHERE follower_id = $1
	    ))
	    AND ($4 = '' OR p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%')
	    AND (p.tags @> $5 OR $5 = '{}')
	GROUP BY p.id, u.username
	ORDER BY p.created_at ` + fq.Sort + `
	LIMIT $2 OFFSET $3
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags)) //

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var feed []PostWithMetadata
	for rows.Next() {
		var post PostWithMetadata
		err = rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, post)
	}
	return feed, nil
}
