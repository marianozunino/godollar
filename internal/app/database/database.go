package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/volatiletech/sqlboiler/boil"
	"go.uber.org/fx"
)

// interface for the ine service

var Module = fx.Options(
	fx.Provide(connectDatabase),
)

type DB struct {
	*sql.DB
}

func connectDatabase() DB {

	// Get a handle to the SQLite database, using mattn/go-sqlite3
	dbConnection, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		panic(err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "./",
	}

	n, err := migrate.Exec(dbConnection, "sqlite3", migrations, migrate.Up)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n===== Applied %d migrations! =====\n", n)

	// Configure SQLBoiler to use the sqlite database
	boil.SetDB(dbConnection)

	return DB{dbConnection}
}
