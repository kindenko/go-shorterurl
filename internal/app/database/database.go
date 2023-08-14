package database

import (
	"database/sql"
	"fmt"
	"log"

	"context"
	"time"

	"github.com/kindenko/go-shorterurl/config"
	"github.com/kindenko/go-shorterurl/internal/app/structures"
	"github.com/kindenko/go-shorterurl/internal/app/utils"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDB struct {
	db  *sql.DB
	cfg *config.AppConfig
}

func (p PostgresDB) Save(fullURL string) (string, error) {

	shortURL := utils.RandString(fullURL)
	query := "insert into shorterurl(short, long) values ($1, $2)"
	_, err := p.db.Exec(query, shortURL, fullURL)
	if err != nil {
		log.Println("Failed to save short link into DB")
		return "", err
	}
	return shortURL, nil
}

func (p PostgresDB) Get(shortURL string) (string, error) {
	var long string
	query := "select long from shorterurl where short=$1"
	row := p.db.QueryRow(query, shortURL)
	if err := row.Scan(&long); err != nil {
		log.Println("Failed to get link from db")
		return "Error in Get from db", err
	}
	return long, nil
}

func (p PostgresDB) Batch(entities []structures.BatchEntity) ([]structures.BatchEntity, error) {
	var resultEntities []structures.BatchEntity
	// var ResultURL string

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	tx, err := p.db.Begin()
	if err != nil {
		fmt.Println("Error while begin tx")
		return resultEntities, err
	}
	for _, v := range entities {
		short := utils.RandString(v.OriginalURL)
		_, err = tx.ExecContext(ctx, "insert into "+"shorterurl"+"(short, long) values ($1, $2)", short, v.OriginalURL)
		if err != nil {
			fmt.Println("Error while ExecContext", err)
			tx.Rollback()
			return resultEntities, nil
		}

		// if p.cfg.ResultURL == "" {
		// 	fmt.Println(p.cfg.ResultURL)
		// 	ResultURL = "http://localhost:8080"
		// } else {
		// 	ResultURL = p.cfg.ResultURL
		// }

		resultEntities = append(resultEntities, structures.BatchEntity{
			CorrelationID: v.CorrelationID,
			ShortURL:      p.cfg.ResultURL + "/" + short,
		})

	}
	return resultEntities, tx.Commit()
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
