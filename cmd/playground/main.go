package main

import (
	"context"
	"fmt"

	"github.com/artrctx/noliteo-core/internal/database"
	"github.com/artrctx/noliteo-core/internal/database/repository"
	"github.com/artrctx/noliteo-core/internal/jwt"
)

func main() {
	db := database.Get().Conn()
	res, err := repository.New(db).ValidateToken(context.Background(), "test-key")
	fmt.Println("key check", res, err)
	r, err := jwt.GenerateToken(jwt.Token{TID: res.ID.String(), Ident: res.Ident.String})
	fmt.Println("jwt gen: ", r, err)

	if err != nil {
		return
	}

	rr, err := jwt.VerifyToken(r)
	fmt.Println("parse", rr, err)
}
