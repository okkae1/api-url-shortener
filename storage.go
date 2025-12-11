package storage

import (
	"encoding/json"
	"os"
	"url-shortener/models"
)

const (
	linksFile = "data/links.json"
	usersFile = "data/users.json"
	statsFile = "data/stats.json"
)

func init() {
	os.MkdirAll("data", 0755)
}

func LoadLinks() ([]models.Link, error) {
	return loadJSON[[]models.Link](linksFile)
}

func SaveLinks(links []models.Link) error {
	return saveJSON(linksFile, links)
}

func LoadUsers() ([]models.User, error) {
	return loadJSON[[]models.User](usersFile)
}

func SaveUsers(users []models.User) error {
	return saveJSON(usersFile, users)
}

func LoadStats() ([]models.Stats, error) {
	return loadJSON[[]models.Stats](statsFile)
}

func SaveStats(stats []models.Stats) error {
	return saveJSON(statsFile, stats)
}

func loadJSON[T any](filename string) (T, error) {
	var data T
	content, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return data, err
	}
	err = json.Unmarshal(content, &data)
	return data, err
}

func saveJSON(filename string, data any) error {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, content, 0644)
}
