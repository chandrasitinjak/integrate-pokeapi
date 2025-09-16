package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/chandrasitinjak/integrate-pokeapi/config"
	"github.com/chandrasitinjak/integrate-pokeapi/internal/logger"
	"github.com/chandrasitinjak/integrate-pokeapi/internal/model"
	"github.com/chandrasitinjak/integrate-pokeapi/internal/repository"
	"go.uber.org/zap"
)

type PokemonService interface {
	Sync(ctx context.Context) error
	GetAll(ctx context.Context) ([]model.Pokemon, error)
}

type pokemonService struct {
	repo    repository.PokemonRepository
	rdb     RedisService
	baseURL string
}

const cacheKey = "items"

func NewPokemonService(repo repository.PokemonRepository, cfg *config.Config, redisSvc RedisService) PokemonService {
	return &pokemonService{repo: repo, baseURL: cfg.BaseURL, rdb: redisSvc}
}

func RandomGender() string {
	genders := []string{"female", "genderless", "male"}
	rand.Seed(time.Now().UnixNano())
	return genders[rand.Intn(len(genders))]
}
func (s *pokemonService) Sync(ctx context.Context) error {
	gender := RandomGender()
	url := s.baseURL + gender

	var apiResp model.PokeAPIGenderResponse
	if err := s.callHTTP(url, &apiResp); err != nil {
		logger.Log.Error("Failed to get data from pokemon api", zap.Error(err))
		return err
	}

	var pokemons []model.Pokemon
	for _, d := range apiResp.PokemonSpeciesDetails {
		pokemons = append(pokemons, model.Pokemon{
			Name:   d.PokemonSpecies.Name,
			URL:    d.PokemonSpecies.URL,
			Rate:   d.Rate,
			Gender: gender,
		})
	}

	if err := s.repo.BulkInsert(ctx, pokemons); err != nil {
		logger.Log.Error("Failed to insert data to db", zap.Error(err))
		return err
	}

	// Hapus cache
	err := s.rdb.Del(ctx, cacheKey)
	if err != nil {
		logger.Log.Error("failed to delete cache", zap.Error(err))
	}

	return nil
}

func (s *pokemonService) GetAll(ctx context.Context) ([]model.Pokemon, error) {
	// Cek di Redis
	cached, err := s.rdb.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	if cached != nil {
		var pokemons []model.Pokemon
		if err := json.Unmarshal(cached, &pokemons); err == nil {
			logger.Log.Error("failed unmarshal data", zap.Error(err))

			return pokemons, nil
		}
	}

	pokemons, err := s.repo.GetAll(ctx)
	if err != nil {
		logger.Log.Error("failed to get data from db", zap.Error(err))
		return nil, err
	}

	// Simpan ke Redis
	// Set manual TTL 5 menit
	ttl := 30 * time.Second
	bytes, _ := json.Marshal(pokemons)
	_ = s.rdb.Set(ctx, cacheKey, bytes, ttl)

	return pokemons, nil
}

func (s *pokemonService) callHTTP(url string, target any) error {
	client := &http.Client{Timeout: 5 * time.Second}

	var lastErr error
	backoff := time.Second

	for attempt := 1; attempt <= 3; attempt++ {
		resp, err := client.Get(url)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: %w", attempt, err)
		} else {
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				lastErr = fmt.Errorf("attempt %d: unexpected status %d", attempt, resp.StatusCode)
			} else {
				if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
					lastErr = fmt.Errorf("attempt %d: decode error: %w", attempt, err)
				} else {
					return nil
				}
			}
		}

		time.Sleep(backoff)
		backoff *= 2
	}

	return fmt.Errorf("all attempts failed: %w", lastErr)
}
