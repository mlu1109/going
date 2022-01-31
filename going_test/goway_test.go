package going_test

import (
	"database/sql"
	"fmt"
	"going"
	"going/datasrc"
	"going/datasrc/postgres"
	"going/migrsrc"
	"going/migrsrc/slice"
	"log"

	_ "github.com/lib/pq"
)

const (
	host        = "localhost"
	port        = "5432"
	user        = "going_user"
	password    = "going_password"
	dbname      = "going_test"
	search_path = "going_schema"
)

var psqlInfo = fmt.Sprintf(
	"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
	host, port, user, password, dbname, search_path)
var db *sql.DB
var ds *postgres.DS

func init() {
	log.Print("Initializing tests...")
	log.Print("Setting up database...")
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	log.Print("Initializing postgres datasource for migrations with managed schema...")
	ds = postgres.New(
		postgres.WithDB(db),
		postgres.WithSchema(search_path, true),
	)
	log.Print("Initialized tests")
}

func NewTestgoing(migrations []*migrsrc.Migration) *going.G {
	g, err := going.New(slice.New(migrations), ds)
	if err != nil {
		log.Panic(err)
	}
	err = g.Clean()
	if err != nil {
		log.Panic(err)
	}
	return g
}

func getAppliedMigrations() ([]*datasrc.Migration, error) {
	err := ds.Lock()
	if err != nil {
		return nil, err
	}
	defer ds.Unlock(false)
	applied, err := ds.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}
	return applied, nil
}
