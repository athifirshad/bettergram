package data

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Like struct {
	ID        uuid.UUID `json:"id"`
	PhotoID   uuid.UUID `json:"photo_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LikeModel struct {
	DB *pgxpool.Pool
}

var ErrDuplicateLike = errors.New("post already liked")

func (m LikeModel) Insert(like *Like) error {
	query := `
		INSERT INTO likes (photo_id, user_id)
		SELECT $1, $2
		WHERE NOT EXISTS (
			SELECT 1 FROM likes WHERE photo_id = $1 AND user_id = $2
		)
		RETURNING id, created_at`

	args := []interface{}{like.PhotoID, like.UserID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, args...).Scan(&like.ID, &like.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrDuplicateLike
		}
		return err
	}
	return nil
}

func (m LikeModel) Delete(photoID uuid.UUID, userID int64) error {
	query := `
		DELETE FROM likes
		WHERE photo_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, photoID, userID)
	return err
}

func (m LikeModel) GetLikesCount(photoID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) FROM likes
		WHERE photo_id = $1`

	var count int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, photoID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
