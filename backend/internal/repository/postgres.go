package repository

import (
	"database/sql"
	"errors"

	"github.com/DragonZero000/url-shortener/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connection(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}
	res := db.Ping()
	if res != nil {
		return nil, res
	}
	return db, nil
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(code, oriURL string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("INSERT INTO public.urls (code, original_url) VALUES ($1, $2)", code, oriURL)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetByCode(code string) (*model.URL, error) {
	var URL model.URL
	err := r.db.QueryRow("SELECT id, code, original_url, created_at FROM public.urls WHERE code = $1", code).Scan(&URL.Id, &URL.Code, &URL.OriginalURL, &URL.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &URL, nil
}

func (r *Repository) GetByOriginalURL(originalURL string) (*model.URL, error) {
	var URL model.URL
	err := r.db.QueryRow("SELECT id, code, original_url, created_at FROM public.urls WHERE original_url = $1", originalURL).Scan(&URL.Id, &URL.Code, &URL.OriginalURL, &URL.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &URL, nil
}
