package data

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Photo struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`  
	PhotoURL  string    `json:"photo_url"`
	Caption   string    `json:"caption,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
type PhotoModel struct {
	DB *pgxpool.Pool
}

func (m PhotoModel) Insert(photo *Photo) error {
	query := `
		INSERT INTO photos (id, user_id, photo_url, caption)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`

	photo.ID = uuid.New() 
	args := []interface{}{photo.ID, photo.UserID, photo.PhotoURL, photo.Caption}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRow(ctx, query, args...).Scan(&photo.CreatedAt)
}

func (m PhotoModel) GetByID(id uuid.UUID) (*Photo, error) {
    query := `
        SELECT p.id, p.user_id, u.username, p.photo_url, p.caption, p.created_at
        FROM photos p
        JOIN users u ON p.user_id = u.id
        WHERE p.id = $1`

    var photo Photo
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    err := m.DB.QueryRow(ctx, query, id).Scan(
        &photo.ID,
        &photo.UserID,
        &photo.Username,
        &photo.PhotoURL,
        &photo.Caption,
        &photo.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &photo, nil
}
func (m PhotoModel) GetByUserID(userID int64) ([]*Photo, error) {
    query := `
        SELECT p.id, p.user_id, u.username, p.photo_url, p.caption, p.created_at
        FROM photos p
        JOIN users u ON p.user_id = u.id
        WHERE p.user_id = $1
        ORDER BY p.created_at DESC`

    rows, err := m.DB.Query(context.Background(), query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var photos []*Photo
    for rows.Next() {
        var photo Photo
        err := rows.Scan(
            &photo.ID,
            &photo.UserID,
            &photo.Username,
            &photo.PhotoURL,
            &photo.Caption,
            &photo.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        photos = append(photos, &photo)
    }

    return photos, nil
}

func (m PhotoModel) Search(query string) ([]*Photo, error) {
    sqlQuery := `
        SELECT p.id, p.user_id, u.username, p.photo_url, p.caption, p.created_at
        FROM photos p
        JOIN users u ON p.user_id = u.id
        WHERE p.caption ILIKE $1 OR u.username ILIKE $1
        ORDER BY p.created_at DESC
    `

    args := []interface{}{"%" + query + "%"}
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := m.DB.Query(ctx, sqlQuery, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var photos []*Photo
    for rows.Next() {
        var photo Photo
        err := rows.Scan(
            &photo.ID,
            &photo.UserID,
            &photo.Username,
            &photo.PhotoURL,
            &photo.Caption,
            &photo.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        photos = append(photos, &photo)
    }

    return photos, nil
}

func (m PhotoModel) GetAll() ([]*Photo, error) {
    query := `
        SELECT p.id, p.user_id, u.username, p.photo_url, p.caption, p.created_at
        FROM photos p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC`

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := m.DB.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var photos []*Photo
    for rows.Next() {
        var photo Photo
        err := rows.Scan(
            &photo.ID,
            &photo.UserID,
            &photo.Username,
            &photo.PhotoURL,
            &photo.Caption,
            &photo.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        photos = append(photos, &photo)
    }

    return photos, nil
}