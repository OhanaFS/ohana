package kv

import (
	"context"
	"errors"
	"sync"
	"time"
)

type MemoryKV struct {
	data map[string]kvData
	lock sync.RWMutex
}

type kvData struct {
	value  string
	expire time.Time
}

func NewMemoryKV() KV {
	return &MemoryKV{
		data: make(map[string]kvData),
	}
}

func (kv *MemoryKV) Get(ctx context.Context, key string) (string, error) {
	kv.lock.RLock()
	defer kv.lock.RUnlock()
	if data, ok := kv.data[key]; ok {
		if data.expire.After(time.Now()) {
			return data.value, nil
		}
		delete(kv.data, key)
	}
	return "", errors.New("key not found")
}

func (kv *MemoryKV) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	kv.data[key] = kvData{
		value:  value,
		expire: time.Now().Add(ttl),
	}
	return nil
}

func (kv *MemoryKV) Delete(ctx context.Context, key string) error {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	delete(kv.data, key)
	return nil
}
