package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"discord-anime-bot/internal/types"
)

// GraphQLNextVariables represents variables for GraphQL next episode query
type GraphQLNextVariables struct {
	ID int `json:"id"`
}

// GraphQLNextRequest represents a GraphQL request for next episode
type GraphQLNextRequest struct {
	Query     string               `json:"query"`
	Variables GraphQLNextVariables `json:"variables"`
}

// GetAnimeByID gets anime details by ID including next airing episode
func GetAnimeByID(animeID int) (*types.AnimeDetails, error) {
	query := `
	query ($id: Int!) {
		Media(id: $id, type: ANIME) {
			id
			title { 
				romaji 
				english 
				native 
			}
			status
			format
			episodes
			nextAiringEpisode {
				episode
				airingAt
				timeUntilAiring
			}
			coverImage { 
				large 
			}
			siteUrl
		}
	}`

	variables := GraphQLNextVariables{
		ID: animeID,
	}

	requestBody := GraphQLNextRequest{
		Query:     query,
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
