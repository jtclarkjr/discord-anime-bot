package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"discord-anime-bot/internal/graphql"
	"discord-anime-bot/internal/types"
)

// SearchAnime searches for anime using the AniList GraphQL API
// Supports both text search and ID lookup
// query: Either a text query to search for or a numeric anime ID
// page: Page number for text search results (ignored for ID search)
// perPage: Number of results per page for text search (ignored for ID search)
// Returns: Page containing matching anime with pagination info
func SearchAnime(query string, page, perPage int) (*types.SearchResponse, error) {
	anilistAPI := os.Getenv("ANILIST_API")

	// Check if the query is a numeric ID
	trimmedQuery := strings.TrimSpace(query)
	if numericID, err := strconv.Atoi(trimmedQuery); err == nil {
		// Search by ID - return single result in page format
		return searchAnimeByID(anilistAPI, numericID)
	}

	// Text search
	return searchAnimeByText(anilistAPI, query, page, perPage)
}

// searchAnimeByID searches for anime by ID and returns it in page format
func searchAnimeByID(anilistAPI string, animeID int) (*types.SearchResponse, error) {

	variables := types.GraphQLSearchByIDVariables{
		ID: animeID,
	}

	requestBody := types.GraphQLRequest[types.GraphQLSearchByIDVariables]{
		Query:     graphql.SearchAnimeByIDQuery,
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var singleResult types.AniListSingleResponse[types.AnimeMedia]
	if err := json.NewDecoder(resp.Body).Decode(&singleResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert single result to page format
	result := &types.SearchResponse{}
	if singleResult.Data.Media.ID != 0 {
		result.Data.Page.Media = []types.AnimeMedia{singleResult.Data.Media}
		result.Data.Page.PageInfo = types.PageInfo{
			Total:       1,
			CurrentPage: 1,
			LastPage:    1,
			HasNextPage: false,
		}
	} else {
		// No anime found with that ID
		result.Data.Page.Media = []types.AnimeMedia{}
		result.Data.Page.PageInfo = types.PageInfo{
			Total:       0,
			CurrentPage: 1,
			LastPage:    1,
			HasNextPage: false,
		}
	}

	return result, nil
}

// searchAnimeByText searches for anime by text query
func searchAnimeByText(anilistAPI, query string, page, perPage int) (*types.SearchResponse, error) {

	variables := types.GraphQLSearchVariables{
		Search:  query,
		Page:    page,
		PerPage: perPage,
	}

	requestBody := types.GraphQLRequest[types.GraphQLSearchVariables]{
		Query:     graphql.SearchAnimeByTextQuery,
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var result types.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
