package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"discord-anime-bot/internal/types"
)

// GetSeasonAnime gets all anime from a specific season and year
func GetSeasonAnime(season string, seasonYear int, page, perPage int) (*types.SeasonAnimeResponse, error) {
	anilistAPI := os.Getenv("ANILIST_API")
	query := `
	query SeasonAnime($season: MediaSeason, $seasonYear: Int, $type: MediaType, $page: Int, $perPage: Int) {
		Page(page: $page, perPage: $perPage) {
			media(season: $season, seasonYear: $seasonYear, type: $type, sort: [POPULARITY_DESC]) {
				id
				title { 
					romaji 
					english 
				}
				coverImage {
					medium
					large
				}
				status
			}
			pageInfo { 
				total 
				currentPage 
				lastPage 
				hasNextPage 
			}
		}
	}`

	variables := types.GraphQLSeasonVariables{
		Season:     strings.ToUpper(season),
		SeasonYear: seasonYear,
		Type:       "ANIME",
		Page:       page,
		PerPage:    perPage,
	}

	requestBody := types.GraphQLRequest[types.GraphQLSeasonVariables]{
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

	var result types.SeasonAnimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
