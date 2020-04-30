package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	pgUrl, err := pq.ParseURL(os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", pgUrl)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db.Ping())

	DB = db
}
