package token

import "database/sql"

type TokenService struct {
	DB *sql.DB
}
