package scheduler

import (
	"ecommerce/internal/sync"
	"time"
)

type SyncScheduler struct {
	odooSync *sync.OdooSync
	stop     chan struct{}
}

func NewSyncScheduler(odooSync *sync.OdooSync) *SyncScheduler {
	return &SyncScheduler{
		odooSync: odooSync,
		stop:     make(chan struct{}),
	}
}

func (s *SyncScheduler) Start() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := s.odooSync.SyncProducts(); err != nil {
					// Log error
				}
				if err := s.odooSync.SyncOrders(); err != nil {
					// Log error
				}
			case <-s.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *SyncScheduler) Stop() {
	close(s.stop)
}
