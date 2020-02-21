package threescale

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type inMemoryLimitsStorage struct {
	internalStorage *gocache.Cache
}

func newInMemoryLimitsStorage() *inMemoryLimitsStorage {
	goCache := gocache.New(gocache.NoExpiration, time.Minute)
	return &inMemoryLimitsStorage{internalStorage: goCache}
}

func (storage *inMemoryLimitsStorage) get(key string) (int, bool) {
	val, exists := storage.internalStorage.Get(key)

	if !exists {
		return 0, false
	}

	return val.(int), exists
}

func (storage *inMemoryLimitsStorage) create(key string, value int, duration time.Duration) bool {
	alreadyExistsErr := storage.internalStorage.Add(key, value, duration)
	return alreadyExistsErr == nil
}

func (storage *inMemoryLimitsStorage) decrement(key string, value int) error {
	_, err := storage.internalStorage.DecrementInt(key, value)
	return err
}
