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
	ParentID  *int64      `json:"parent_id,omitempty"`
	Deleted   bool        `json:"-"`
	User      CommentUser `json:"user"`
	Replies   []*Comment  `json:"replies,omitempty"`
}

type CommentUser struct {
	Username string `json:"username"`
}

type SwaggerCommentResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PostID    int64     `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	// Deleted   bool             `json:"deleted"`
	ParentID *int64           `json:"parent_id,omitempty"`
	User     CommentUser      `json:"user"`
	Replies  []CommentShallow `json:"replies,omitempty"`
}

type CommentShallow struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
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

func (s *PostgresCommentStore) GetByPostID(ctx context.Context, postID int64) ([]*Comment, error) {
	query :=
		`
	WITH RECURSIVE comment_tree AS (
  SELECT
    c.id,
    c.user_id,
    c.content,
    c.created_at,
    c.parent_id,
    c.post_id,
	c.deleted,
    1 AS depth,
    LPAD(c.id::text, 10, '0') AS path
  FROM comments c
  WHERE c.post_id = $1 AND c.parent_id IS NULL

  UNION ALL

  SELECT
    child.id,
    child.user_id,
    child.content,
    child.created_at,
    child.parent_id,
    child.post_id,
	child.deleted,
    ct.depth + 1,
    ct.path || '.' || LPAD(child.id::text, 10, '0') AS path
  FROM comments child
  JOIN comment_tree ct ON ct.id = child.parent_id
)

SELECT
  c.id, c.user_id, u.username, c.post_id,
  c.content, c.created_at, c.parent_id, c.deleted
FROM comment_tree c
JOIN users u ON c.user_id = u.id
ORDER BY path;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*Comment{}
	commentMap := make(map[int64]*Comment)
	for rows.Next() {
		var comment Comment
		comment.User = CommentUser{}
		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.User.Username,
			&comment.PostID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.ParentID,
			&comment.Deleted,
		)
		if err != nil {
			return nil, err
		}
		c := &comment
		commentMap[c.ID] = c
		comments = append(comments, c)
	}
	for _, comment := range comments {
		if comment.Deleted {
			comment.Content = "[deleted]"
			comment.User = CommentUser{
				Username: "[deleted]",
			}
		}
	}

	var roots []*Comment
	for _, comment := range comments {
		if comment.ParentID == nil {
			roots = append(roots, comment)
		} else if parent, ok := commentMap[*comment.ParentID]; ok {
			parent.Replies = append(parent.Replies, comment)
		}
	}
	return roots, nil

	// return comments, nil
}

func (s *PostgresCommentStore) Delete(ctx context.Context, id int64) error {
	query := `UPDATE comments SET content=$1, deleted = true WHERE id = $2`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	res, err := s.db.ExecContext(ctx, query, DeletedContent, id)
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

func (s *PostgresCommentStore) GetByID(ctx context.Context, id int64) (*Comment, error) {
	query := `SELECT id,post_id, user_id, content, created_at, parent_id FROM comments WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	row := s.db.QueryRowContext(ctx, query, id)
	comment := &Comment{}
	err := row.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.ParentID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return comment, nil
}
