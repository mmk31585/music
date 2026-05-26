package repository

import "music/internal/domain"

type SongRepository interface {
	Save(song *domain.Song) error
	FindByID(id string) (*domain.Song, error)
	FindAll() ([]domain.Song, error)
	Delete(id string) error
}
