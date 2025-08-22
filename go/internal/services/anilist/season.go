package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"discord-anime-bot/internal/graphql"
	"discord-anime-bot/internal/types"
)

// GetSeasonAnime gets all anime from a specific season and year
func GetSeasonAnime(season string, seasonYear int, page, perPage int) (*types.SeasonAnimeResponse, error) {
	anilistAPI := os.Getenv("ANILIST_API")

	variables := types.GraphQLSeasonVariables{
		Season:     strings.ToUpper(season),
		SeasonYear: seasonYear,
		Type:       "ANIME",
		Page:       page,
		PerPage:    perPage,
	}

	requestBody := types.GraphQLRequest[types.GraphQLSeasonVariables]{
		Query:     graphql.GetSeasonalAnimeQuery,
		Variables: variables,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(anilistAPI, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var result types.SeasonAnimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
