package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/artrctx/quoin-core/internal/database"
	"github.com/artrctx/quoin-core/internal/database/repository"
)

func main() {
	db := database.Get().Conn()
	res, err := repository.New(db).ValidateToken(context.Background(), "test-key")
	fmt.Printf("%T\n", res)
	fmt.Printf("%v\n", reflect.TypeOf(res))
	fmt.Println(res, err)
}
