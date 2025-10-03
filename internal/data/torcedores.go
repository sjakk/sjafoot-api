package data

import (
	"context"
	"database/sql"
	"time"
)

type Torcedor struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Nome      string    `json:"nome"`
	Email     string    `json:"email"`
	TimeClube string    `json:"time"`
}

type TorcedorModel struct {
	DB *sql.DB
}

func (m TorcedorModel) Insert(torcedor *Torcedor) error {
	query := `
        INSERT INTO torcedores (nome, email, time_clube)
        VALUES ($1, $2, $3)
        RETURNING id, created_at`

	args := []interface{}{torcedor.Nome, torcedor.Email, torcedor.TimeClube}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// The Scan populates the ID and CreatedAt fields of the passed-in struct.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&torcedor.ID, &torcedor.CreatedAt)
}

func (m TorcedorModel) GetAllForTeam(team string) ([]*Torcedor, error) {
	query := `
        SELECT id, created_at, nome, email, time_clube
        FROM torcedores
        WHERE time_clube = $1
        ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, team)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var torcedores []*Torcedor

	for rows.Next() {
		var torcedor Torcedor
		err := rows.Scan(
			&torcedor.ID,
			&torcedor.CreatedAt,
			&torcedor.Nome,
			&torcedor.Email,
			&torcedor.TimeClube,
		)
		if err != nil {
			return nil, err
		}
		torcedores = append(torcedores, &torcedor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return torcedores, nil
}
