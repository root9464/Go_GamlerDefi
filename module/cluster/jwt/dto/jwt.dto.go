package jwt_dto

type UserData struct {
	ID int64 `json:"id" validate:"required"`
}

type UserJwtPayload struct {
	Iss  string `json:"iss" validate:"required"`
	Sub  int64  `json:"sub" validate:"required"`
	Iat  int64  `json:"iat" validate:"required"`
	Exp  int64  `json:"exp" validate:"required"`
	Hash string `json:"user_hash" validate:"required"`
}
