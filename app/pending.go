package main

import (
	"sync"
	"time"
)

type StampData struct {
	Blockheight string
	Stamp       string
	CreatedAt   time.Time
}

type PendingStamps struct {
	mu    sync.RWMutex
	items map[string]StampData
}

func NewPendingStamps() *PendingStamps {
	ps := &PendingStamps{
		items: make(map[string]StampData),
	}
	go ps.cleanup() // Start cleanup routine
	return ps
}

func (ps *PendingStamps) Set(paymentHash string, data StampData) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	data.CreatedAt = time.Now()
	ps.items[paymentHash] = data
}

func (ps *PendingStamps) Get(paymentHash string) (StampData, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	data, exists := ps.items[paymentHash]
	return data, exists
}

func (ps *PendingStamps) Delete(paymentHash string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	delete(ps.items, paymentHash)
}

// Cleanup routine to remove expired pending payments
func (ps *PendingStamps) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		ps.mu.Lock()
		now := time.Now()
		for hash, data := range ps.items {
			if now.Sub(data.CreatedAt) > 15*time.Minute {
				delete(ps.items, hash)
			}
		}
		ps.mu.Unlock()
	}
}
