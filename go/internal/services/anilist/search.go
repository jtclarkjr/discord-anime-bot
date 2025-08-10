package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"discord-anime-bot/internal/types"
)

const anilistAPI = "https://graphql.anilist.co"

// SearchAnime searches for anime using the AniList GraphQL API
func SearchAnime(query string, page, perPage int) (*types.SearchResponse, error) {
	searchQuery := `
	query ($search: String, $page: Int, $perPage: Int) {
		Page(page: $page, perPage: $perPage) {
			pageInfo {
				total
				currentPage
				lastPage
				hasNextPage
			}
			media(search: $search, type: ANIME) {
				id
				title {
					romaji
					english
					native
				}
				format
				status
				coverImage {
					large
				}
				siteUrl
			}
		}
	}`

	variables := types.GraphQLSearchVariables{
		Search:  query,
		Page:    page,
		PerPage: perPage,
	}

	requestBody := types.GraphQLRequest[types.GraphQLSearchVariables]{
		Query:     searchQuery,
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

	var result types.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
