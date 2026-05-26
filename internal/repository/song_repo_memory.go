package repository

import (
	"encoding/json"
	"errors"
	"music/internal/domain"
	"os"
	"sync"
)

type InMemorySongRepo struct {
	songs    map[string]domain.Song
	mu       sync.RWMutex
	filePath string
}

func NewInMemorySongRepo(dataFilePath string) *InMemorySongRepo {
	repo := &InMemorySongRepo{
		songs:    make(map[string]domain.Song),
		filePath: dataFilePath,
	}
	repo.loadFromFile()
	return repo
}

func (r *InMemorySongRepo) Save(song *domain.Song) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.songs[song.ID] = *song
	return r.saveToFile()
}

func (r *InMemorySongRepo) FindByID(id string) (*domain.Song, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	song, exists := r.songs[id]
	if !exists {
		return nil, errors.New("song not found")
	}
	return &song, nil
}

func (r *InMemorySongRepo) FindAll() ([]domain.Song, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]domain.Song, 0, len(r.songs))
	for _, s := range r.songs {
		list = append(list, s)
	}
	return list, nil
}

func (r *InMemorySongRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.songs[id]; !exists {
		return errors.New("song not found")
	}
	delete(r.songs, id)
	return r.saveToFile()
}

func (r *InMemorySongRepo) saveToFile() error {
	list := make([]domain.Song, 0, len(r.songs))
	for _, s := range r.songs {
		list = append(list, s)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}

func (r *InMemorySongRepo) loadFromFile() {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return // file might not exist yet
	}
	var list []domain.Song
	if err := json.Unmarshal(data, &list); err != nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, s := range list {
		r.songs[s.ID] = s
	}
}
