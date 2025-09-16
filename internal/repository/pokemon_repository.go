package repository

import (
	"context"
	"database/sql"

	"github.com/chandrasitinjak/integrate-pokeapi/internal/logger"
	"github.com/chandrasitinjak/integrate-pokeapi/internal/model"
	"go.uber.org/zap"
)

type PokemonRepository interface {
	GetAll(ctx context.Context) ([]model.Pokemon, error)
	BulkInsert(ctx context.Context, pokemons []model.Pokemon) error
}

type pokemonRepository struct {
	db *sql.DB
}

func NewPokemonRepository(db *sql.DB) PokemonRepository {
	return &pokemonRepository{db: db}
}

func (r *pokemonRepository) GetAll(ctx context.Context) ([]model.Pokemon, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, url, gender, rate FROM pokemon")
	if err != nil {
		logger.Log.Error("Failed to get all pokemon's data", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var pokemons []model.Pokemon
	for rows.Next() {
		var p model.Pokemon
		if err := rows.Scan(&p.ID, &p.Name, &p.URL, &p.Gender, &p.Rate); err != nil {
			logger.Log.Error("error scanning rows data", zap.Error(err))
			return nil, err
		}
		pokemons = append(pokemons, p)
	}
	return pokemons, nil
}

func (r *pokemonRepository) BulkInsert(ctx context.Context, pokemons []model.Pokemon) error {
	if len(pokemons) == 0 {
		return nil
	}

	query := "INSERT INTO pokemon (name, url, gender, rate) VALUES "
	args := []interface{}{}

	placeholders := ""
	for i, p := range pokemons {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "(?, ?, ?, ?)"
		args = append(args, p.Name, p.URL, p.Gender, p.Rate)
	}

	query += placeholders

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Failed to bulk insert", zap.Error(err))
	}
	return err
}
