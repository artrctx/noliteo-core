package main

import (
	"context"
	"fmt"

	"github.com/artrctx/quoin-core/internal/database"
	"github.com/artrctx/quoin-core/internal/database/repository"
)

func main() {
	db := database.Get().Conn()
	res, err := repository.New(db).ValidateToken(context.Background(), "test-keywrong")
	fmt.Printf("%T\n -- type --", res)
	fmt.Printf("%v, value\n", res)
	fmt.Println(res, err)
}
