package data

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	Users    UserModel
	Tokens   TokenModel
	Photos   PhotoModel
	Likes    LikeModel
	Comments CommentModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Users:    UserModel{DB: db},
		Tokens:   TokenModel{DB: db},
		Photos:   PhotoModel{DB: db},
		Likes:    LikeModel{DB: db},
		Comments: CommentModel{DB: db},
	}
}