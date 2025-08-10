package anilist

import (
	"fmt"
	"log"

	"discord-anime-bot/internal/services/openai"
	"discord-anime-bot/internal/types"
)

// FindAnimeWithDetails finds anime using AI description and returns AniList details
func FindAnimeWithDetails(description, openAIAPIKey string) ([]types.AnimeMatch, error) {
	if openAIAPIKey == "" {
		return nil, fmt.Errorf("OpenAI is not configured. Please set OPENAI_API_KEY environment variable to use AI-powered anime search")
	}

	// Get AI recommendations
	recommendations, err := openai.FindAnimeByDescription(description, openAIAPIKey)
	if err != nil {
		return nil, err
	}

	var matches []types.AnimeMatch

	// Search for each recommendation on AniList
	for _, rec := range recommendations {
		searchResults, err := SearchAnime(rec.Title, 1, 5)
		if err != nil {
			log.Printf("Could not search for anime %q: %v", rec.Title, err)
			continue
		}

		if len(searchResults.Data.Page.Media) > 0 {
			// Use the first (most relevant) result
			anime := searchResults.Data.Page.Media[0]
			matches = append(matches, types.AnimeMatch{
				Anime:      anime,
				Reason:     rec.Reason,
				Confidence: rec.Confidence,
			})
		}
	}

	// Sort by confidence score (higher first)
	for i := 0; i < len(matches)-1; i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[i].Confidence < matches[j].Confidence {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	return matches, nil
}
