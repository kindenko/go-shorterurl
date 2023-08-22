package database

import (
	"database/sql"
	"errors"

	"log"

	"context"
	"time"

	"github.com/kindenko/go-shorterurl/config"
	e "github.com/kindenko/go-shorterurl/internal/app/errors"
	"github.com/kindenko/go-shorterurl/internal/app/structures"
	"github.com/kindenko/go-shorterurl/internal/app/utils"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDB struct {
	db  *sql.DB
	cfg config.AppConfig
}

func (p PostgresDB) Save(fullURL string, shortURL string) (string, error) {
	var short string

	query := "insert into shorterurl(shortURL, longURL) values ($1, $2)"
	_, err := p.db.Exec(query, shortURL, fullURL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			log.Println(err)
			short, err = p.GetShortURL(fullURL)
			if err != nil {
				log.Println("faled search previously saved url")
				return "", nil
			}
			return short, e.ErrUniqueValue
		}
	}
	return shortURL, nil
}

func (p PostgresDB) GetShortURL(fullURL string) (string, error) {
	var short string
	query := "select short from shorterurl where longURL=$1"
	row := p.db.QueryRow(query, fullURL)
	if err := row.Scan(&short); err != nil {
		return "", err
	}
	return short, nil
}

func (p PostgresDB) Get(shortURL string) (string, error) {
	var long string
	query := "select long from shorterurl where shortURL=$1"
	row := p.db.QueryRow(query, shortURL)
	if err := row.Scan(&long); err != nil {
		log.Println("Failed to get link from db")
		return "Error in Get from db", err
	}
	return long, nil
}

func (p PostgresDB) Batch(entities []structures.BatchEntity) ([]structures.BatchEntity, error) {
	var resultEntities []structures.BatchEntity
	var ResultURL string

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	tx, err := p.db.Begin()
	if err != nil {
		log.Println("Error while begin tx")
		return resultEntities, err
	}
	for _, v := range entities {
		short := utils.RandString(v.OriginalURL)
		_, err = tx.ExecContext(ctx, "insert into shorterurl (shortURL, longURL) values ($1, $2)", short, v.OriginalURL)
		if err != nil {
			log.Println("Error while ExecContext", err)
			tx.Rollback()
			return resultEntities, nil
		}
		// костылище, не смог исправить
		if p.cfg.ResultURL == "" {
			log.Println(p.cfg.ResultURL)
			ResultURL = "http://localhost:8080"
		} else {
			ResultURL = p.cfg.ResultURL
		}

		resultEntities = append(resultEntities, structures.BatchEntity{
			CorrelationID: v.CorrelationID,
			ShortURL:      ResultURL + "/" + short,
		})

	}
	return resultEntities, tx.Commit()
}

func (p PostgresDB) Ping() error {
	if err := p.db.Ping(); err != nil {
		return err
	}
	return nil
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

	_, err = db.Exec("create table if not exists shorterurl(id serial not null, shortURL text not null not null, longURL text not null); create unique index on shorterurl (long)")
	if err != nil {
		log.Println(err)
		return nil
	}

	return &PostgresDB{db: db}
}
