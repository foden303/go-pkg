package snapshot

import (
	"os"
	"path/filepath"

	"github.com/bits-and-blooms/bloom/v3"
)

func SaveBloom(path string, bf *bloom.BloomFilter) error {
	data, err := bf.MarshalBinary()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadBloom(path string) (*bloom.BloomFilter, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	bf := bloom.New(1, 1) // dummy, will overwrite
	if err := bf.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	return bf, nil
}

func LoadOrCreateBloom(path string, n uint, fpRate float64) (*bloom.BloomFilter, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); err == nil {
		// File exists → load
		return LoadBloom(path)
	}
	// File does not exist → create new
	return bloom.NewWithEstimates(n, fpRate), nil
}
