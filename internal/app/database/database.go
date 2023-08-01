package database

import (
	"database/sql"
	"log"

	"github.com/kindenko/go-shorterurl/internal/app/utils"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDB struct {
	db *sql.DB
}

func (p PostgresDB) Save(fullURL string) (string, error) {

	shortURL := utils.RandString()
	query := "insert into shorterurl(short, long) values ($1, $2)"
	_, err := p.db.Exec(query, shortURL, fullURL)
	if err != nil {
		log.Println("Failed to save short link into DB")
		return "", nil
	}
	return shortURL, nil
}

func (p PostgresDB) Get(shortURL string) (string, error) {
	var long string
	query := "select long from shorterurl where short=$1"
	row := p.db.QueryRow(query, shortURL)
	if err := row.Scan(&long); err != nil {
		log.Println("Failed to get link from db")
		return "Error in Get from db", nil
	}
	return long, nil
}

func InitDB(path string, baseurl string) *PostgresDB {
	if path == "" {
		return nil
	}

	db, err := sql.Open("pgx", path)
	if err != nil {
		log.Println(err)
		return nil
	}

	_, err = db.Exec("create table if not exists shorterurl(id serial, short text not null, long text not null)")
	if err != nil {
		log.Println(err)
		return nil
	}

	return &PostgresDB{db: db}
}
