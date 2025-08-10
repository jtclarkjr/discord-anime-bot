package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"discord-anime-bot/internal/types"
)

// GetReleasingAnime gets all currently releasing anime
func GetReleasingAnime(page, perPage int) (*types.ReleasingAnimeResponse, error) {
	query := `
	query ($page: Int, $perPage: Int) {
		Page(page: $page, perPage: $perPage) {
			media(type: ANIME, status: RELEASING, sort: [POPULARITY_DESC]) {
				id
				title { 
					romaji 
					english 
				}
				nextAiringEpisode {
					episode
					airingAt
				}
			}
			pageInfo { 
				total 
				currentPage 
				lastPage 
				hasNextPage 
			}
		}
	}`

	variables := types.GraphQLSearchVariables{
		Page:    page,
		PerPage: perPage,
	}

	requestBody := types.GraphQLRequest{
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

	var result types.ReleasingAnimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
