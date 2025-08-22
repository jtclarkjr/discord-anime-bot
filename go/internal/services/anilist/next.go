package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"discord-anime-bot/internal/graphql"
	"discord-anime-bot/internal/types"
)

// GetAnimeByID gets anime details by ID including next airing episode
func GetAnimeByID(animeID int) (*types.AnimeDetails, error) {
	anilistAPI := os.Getenv("ANILIST_API")

	variables := types.GraphQLNextVariables{
		ID: animeID,
	}

	requestBody := types.GraphQLRequest[types.GraphQLNextVariables]{
		Query:     graphql.GetAnimeDetailsQuery,
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

	var result types.AnimeDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Data.Media, nil
}
