package pkg

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {

	// setting file env yang berisi informasi db
    dbEnv := []any{}
	dbEnv = append(dbEnv, os.Getenv("DBUSER"))
	dbEnv = append(dbEnv, os.Getenv("DBPASS"))
	dbEnv = append(dbEnv, os.Getenv("DBHOST"))
	dbEnv = append(dbEnv, os.Getenv("DBPORT"))
	dbEnv = append(dbEnv, os.Getenv("DBNAME"))

	// setup database connection
	dbString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbEnv...)

	var err error
	DB, err := pgxpool.New(context.Background(), dbString)

	if err != nil {
		return nil, err
	}

	err = DB.Ping(context.Background())

	if err != nil {
		return nil, fmt.Errorf("[CONNECTED FAILED]: %w", err)
	}

	// dbClient, err := pgxpool.New(context.Background(), dbString)

	log.Println("[CONNECTED SUCCESS]: Connected to postgresql")
	return DB, nil
}