package service

import (
	"fmt"
	"time"
	"url-shortener/models"
	"url-shortener/storage"
)

type URLShortener struct {
	links []models.Link
	users []models.User
	stats []models.Stats
}

func New() (*URLShortener, error) {
	links, err := storage.LoadLinks()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки ссылок: %w", err)
	}
	users, err := storage.LoadUsers()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки пользователей: %w", err)
	}
	stats, err := storage.LoadStats()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки статистики: %w", err)
	}
	return &URLShortener{
		links: links,
		users: users,
		stats: stats,
	}, nil
}

func (s *URLShortener) AddLink(originalURL, shortCode, userID string) error {
	for _, l := range s.links {
		if l.ShortCode == shortCode {
			return fmt.Errorf("код %s уже занят", shortCode)
		}
	}
	link := models.Link{
		ID:          fmt.Sprintf("link_%d", time.Now().UnixNano()),
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}
	s.links = append(s.links, link)
	return storage.SaveLinks(s.links)
}

func (s *URLShortener) UpdateLink(shortCode, newURL string) error {
	for i, l := range s.links {
		if l.ShortCode == shortCode {
			s.links[i].OriginalURL = newURL
			return storage.SaveLinks(s.links)
		}
	}
	return fmt.Errorf("ссылка с кодом %s не найдена", shortCode)
}

func (s *URLShortener) DeleteLink(shortCode string) error {
	for i, l := range s.links {
		if l.ShortCode == shortCode {
			s.links = append(s.links[:i], s.links[i+1:]...)
			// Удаляем статистику для этой ссылки
			for j, stat := range s.stats {
				if stat.LinkID == l.ID {
					s.stats = append(s.stats[:j], s.stats[j+1:]...)
					break
				}
			}
			storage.SaveStats(s.stats)
			return storage.SaveLinks(s.links)
		}
	}
	return fmt.Errorf("ссылка с кодом %s не найдена", shortCode)
}

func (s *URLShortener) ListLinks() {
	if len(s.links) == 0 {
		fmt.Println("Нет ссылок.")
		return
	}
	for _, l := range s.links {
		fmt.Printf("ID: %s | Код: %s → %s | Владелец: %s\n", l.ID, l.ShortCode, l.OriginalURL, l.UserID)
	}
}

func (s *URLShortener) AddUser(name, email string) error {
	for _, u := range s.users {
		if u.Email == email {
			return fmt.Errorf("пользователь с email %s уже существует", email)
		}
	}
	user := models.User{
		ID:        fmt.Sprintf("user_%d", time.Now().UnixNano()),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
	s.users = append(s.users, user)
	return storage.SaveUsers(s.users)
}

func (s *URLShortener) ListUsers() {
	if len(s.users) == 0 {
		fmt.Println("Нет пользователей.")
		return
	}
	for _, u := range s.users {
		fmt.Printf("ID: %s | Имя: %s | Email: %s\n", u.ID, u.Name, u.Email)
	}
}
