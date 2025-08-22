package anilist

import (
	"discord-anime-bot/internal/config"
	"encoding/json"
	"os"
	"sync"
)

var watchlistFile = config.LoadConfig().WatchlistFile
var watchlistMu sync.Mutex

type Watchlist map[string][]int // userID -> []animeID

// loadWatchlist loads the watchlist from disk
func loadWatchlist() (Watchlist, error) {
	watchlistMu.Lock()
	defer watchlistMu.Unlock()
	file, err := os.Open(watchlistFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(Watchlist), nil
		}
		return nil, err
	}
	defer file.Close()
	var wl Watchlist
	if err := json.NewDecoder(file).Decode(&wl); err != nil {
		return make(Watchlist), nil // fallback to empty
	}
	return wl, nil
}

// saveWatchlist saves the watchlist to disk
func saveWatchlist(wl Watchlist) error {
	watchlistMu.Lock()
	defer watchlistMu.Unlock()
	file, err := os.Create(watchlistFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(wl)
}

// AddToWatchlist adds an anime to a user's watchlist
func AddToWatchlist(userID string, animeID int) (string, error) {
	wl, err := loadWatchlist()
	if err != nil {
		return "", err
	}
	for _, id := range wl[userID] {
		if id == animeID {
			return "Anime already in your watchlist.", nil
		}
	}
	wl[userID] = append(wl[userID], animeID)
	if err := saveWatchlist(wl); err != nil {
		return "", err
	}
	return "Anime added to your watchlist.", nil
}

// RemoveFromWatchlist removes an anime from a user's watchlist
func RemoveFromWatchlist(userID string, animeID int) (string, error) {
	wl, err := loadWatchlist()
	if err != nil {
		return "", err
	}
	ids := wl[userID]
	newIDs := []int{}
	found := false
	for _, id := range ids {
		if id == animeID {
			found = true
			continue
		}
		newIDs = append(newIDs, id)
	}
	if !found {
		return "Anime not found in your watchlist.", nil
	}
	wl[userID] = newIDs
	if err := saveWatchlist(wl); err != nil {
		return "", err
	}
	return "Anime removed from your watchlist.", nil
}

// GetUserWatchlist returns a user's watchlist
func GetUserWatchlist(userID string) ([]int, error) {
	wl, err := loadWatchlist()
	if err != nil {
		return nil, err
	}
	return wl[userID], nil
}
