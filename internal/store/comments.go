package store

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
	PostID    int64       `json:"post_id"`
	Content   string      `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
	User      CommentUser `json:"user"`
}

type CommentUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type PostgresCommentStore struct {
	db *sql.DB
}

func (s *PostgresCommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content)
              VALUES ($1, $2, $3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	row := s.db.QueryRowContext(ctx, query, comment.PostID, comment.UserID, comment.Content)
	err := row.Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresCommentStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
	SELECT c.id,c.user_id,c.post_id,c.content,c.created_at, users.username, users.id
	FROM comments AS c
	JOIN users on users.id = c.user_id
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []Comment{}
	for rows.Next() {
		var comment Comment
		comment.User = CommentUser{}
		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.PostID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.User.Username,
			&comment.User.ID,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
