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

// GetReleasingAnime gets all currently releasing anime
func GetReleasingAnime(page, perPage int) (*types.ReleasingAnimeResponse, error) {
	anilistAPI := os.Getenv("ANILIST_API")

	variables := types.GraphQLSearchVariables{
		Page:    page,
		PerPage: perPage,
	}

	requestBody := types.GraphQLRequest[types.GraphQLSearchVariables]{
		Query:     graphql.GetReleasingAnimeQuery,
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

	var result types.ReleasingAnimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
