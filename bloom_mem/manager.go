package bloom_mem

import (
	"context"
	"go-pkg/bloom_mem/snapshot"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
)

type BloomManager struct {
	filter *bloom.BloomFilter
	path   string
	cancel context.CancelFunc
}

func NewBloomManager(
	path string,
	n uint,
	fpRate float64,
) (*BloomManager, error) {
	filter, err := snapshot.LoadOrCreateBloom(path, n, fpRate)
	if err != nil {
		return nil, err
	}

	return &BloomManager{
		filter: filter,
		path:   path,
	}, nil
}

func (m *BloomManager) Filter() *bloom.BloomFilter {
	return m.filter
}

func (m *BloomManager) Start(ctx context.Context, snapshotInterval time.Duration) {
	ctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	// Snapshot
	go func() {
		ticker := time.NewTicker(snapshotInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_ = snapshot.SaveBloom(m.path, m.filter)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Listen SIGTERM
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		_ = snapshot.SaveBloom(m.path, m.filter)
		cancel()
	}()
}

func (m *BloomManager) Stop() {
	if m.cancel != nil {
		m.cancel()
	}
	_ = snapshot.SaveBloom(m.path, m.filter)
}
