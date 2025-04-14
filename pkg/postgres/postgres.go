package postgres

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool     *pgxpool.Pool
	initOnce sync.Once
)

func GetPool() *pgxpool.Pool {
	initOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Second)
		defer cancel()

		dsn := os.Getenv("DATABASE_URL")

		if dsn == "" {
			fmt.Println("DATABASE_URL não definida")
		}

		var err error
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			fmt.Printf("Erro ao conectar ao banco: %v", err)
		}

		fmt.Println("Conexão com o banco estabelecida")
	})
	return pool
}

func ClosePool() {
	if pool != nil {
		pool.Close()
		fmt.Println("Conexão com o banco encerrada")
	}
}
