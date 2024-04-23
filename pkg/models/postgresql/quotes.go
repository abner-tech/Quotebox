package postgresql

import (
	"database/sql"
	"errors"

	models "amencia.net/quotebox/pkg/models"
)

type QuoteModel struct {
	DB *sql.DB
}

func (m *QuoteModel) Insert(author, categry, body string) (int, error) {
	var id int

	s := `
	INSERT INTO quotations(author_name, category, quote)
	VALUES ($1, $2, $3)
	RETURNING quotations_id
	`
	err := m.DB.QueryRow(s, author, categry, body).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *QuoteModel) Read() ([]*models.Quote, error) {
	// SQL statment
	readQuote := `
	SELECT author_name, category, quote
	FROM quotations
	LIMIT 10
	`

	rows, err := m.DB.Query(readQuote)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quotes := []*models.Quote{}
	for rows.Next() {
		q := &models.Quote{}
		err = rows.Scan(&q.Author_name, &q.Category, &q.Quote)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, q)
	}
	if err != nil {
		return nil, err
	}
	return quotes, nil
}

func (m *QuoteModel) Getid(id int) (*models.Quote, error) {
	s := `
	SELECT author_name, category, quote
	FROM quotations
	WHERE quotations_id = $1
	`
	q := &models.Quote{}
	err := m.DB.QueryRow(s, id).Scan(&q.Author_name, &q.Category, &q.Quote)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrRecordNotFound
		} else {
			return nil, err
		}
	}
	return q, err

}
