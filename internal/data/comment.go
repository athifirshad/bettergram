package data

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Comment struct {
	ID        uuid.UUID `json:"id"`
	PhotoID   uuid.UUID `json:"photo_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentModel struct {
	DB *pgxpool.Pool
}

func (m CommentModel) Insert(comment *Comment) error {
	query := `
		INSERT INTO comments (photo_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	args := []interface{}{comment.PhotoID, comment.UserID, comment.Content}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRow(ctx, query, args...).Scan(&comment.ID, &comment.CreatedAt)
}

func (m CommentModel) GetByPhotoID(photoID uuid.UUID) ([]*Comment, error) {
	query := `
		SELECT c.id, c.photo_id, c.user_id, u.username, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.photo_id = $1
		ORDER BY c.created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query, photoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment

	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PhotoID,
			&comment.UserID,
			&comment.Username,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}