package model

// Untuk mapping response API PokeAPI
type PokeAPIGenderResponse struct {
	ID                    int                    `json:"id"`
	Name                  string                 `json:"name"`
	PokemonSpeciesDetails []PokemonSpeciesDetail `json:"pokemon_species_details"`
}

type PokemonSpeciesDetail struct {
	PokemonSpecies struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon_species"`
	Rate int `json:"rate"`
}

// Untuk entity yang disimpan di DB
type Pokemon struct {
	ID     int
	Name   string
	URL    string
	Rate   int
	Gender string
}
