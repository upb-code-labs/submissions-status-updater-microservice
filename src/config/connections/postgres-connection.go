package connections

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/config"
)

var pgConnection *sql.DB

func GetPostgresConnection() *sql.DB {
	if pgConnection == nil {
		// Connect
		pgConnectionString := config.GetEnvironment().DbConnectionString
		db, err := sql.Open("postgres", pgConnectionString)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Check connection
		if err = db.Ping(); err != nil {
			log.Fatal("[Postgres]: Error connecting to the database: " + err.Error())
		}

		// Set connection
		log.Println("[Postgres]: Connected")
		pgConnection = db
	}

	return pgConnection
}

func ClosePostgresConnection() {
	if pgConnection != nil {
		pgConnection.Close()
	}
}
