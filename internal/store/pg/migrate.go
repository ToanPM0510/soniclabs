package pg

import (
	"context"
	"os"
)

func MustExecFile(ctx context.Context, db *DB, path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	_, err = db.Pool.Exec(ctx, string(b))
	if err != nil {
		panic(err)
	}
}
