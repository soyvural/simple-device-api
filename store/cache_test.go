package store

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/soyvural/simple-device-api/types"

	"github.com/google/go-cmp/cmp"
)

func TestCreateDevice(t *testing.T) {
	tests := []struct {
		desc       string
		itemCount  int
		wantedSize int
	}{
		{
			desc:       "add one item",
			itemCount:  1,
			wantedSize: 1,
		},
		{
			desc:       "add many items as much as limit",
			itemCount:  limit,
			wantedSize: limit,
		},
		{
			desc:       "add many items more than limit",
			itemCount:  limit + 1,
			wantedSize: limit,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			c := NewCache()
			for i := 0; i < tc.itemCount; i++ {
				c.Insert(types.Device{ID: fmt.Sprintf("%d", i)})
			}
			c.mu.RLock()
			gotSize := len(c.data)
			c.mu.RUnlock()

			if diff := cmp.Diff(tc.wantedSize, gotSize); diff != "" {
				t.Fatalf("Item size mismatch (-want +got): %s\n", diff)
			}
		})
	}
}

func TestGetDevice(t *testing.T) {
	tests := []struct {
		desc      string
		id        string
		devices   []types.Device
		wantedNil bool
	}{
		{
			desc:    "get already existing item",
			id:      "1",
			devices: []types.Device{{ID: "1"}},
		},
		{
			desc:      "get non existing item",
			id:        "2",
			devices:   []types.Device{{ID: "1"}},
			wantedNil: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			c := NewCache()
			for _, d := range tc.devices {
				c.Insert(d)
			}
			if d := c.Get(tc.id); d == nil && !tc.wantedNil {
				t.Fatalf("got nil but wanted an existing device, id %s.", tc.id)
			}
		})
	}
}

func TestDeleteDevice(t *testing.T) {
	tests := []struct {
		desc      string
		id        string
		devices   []types.Device
		wantedNil bool
	}{
		{
			desc:    "get already existing item",
			id:      "1",
			devices: []types.Device{{ID: "1"}},
		},
		{
			desc:      "get non existing item",
			id:        "2",
			devices:   []types.Device{{ID: "1"}},
			wantedNil: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			c := NewCache()
			for _, d := range tc.devices {
				c.Insert(d)
			}
			if d := c.Delete(tc.id); d == nil && !tc.wantedNil {
				t.Fatalf("got nil but wanted an existing device, id %s.", tc.id)
			}
		})
	}
}

func TestCacheConcurrencyWithProducerConsumer(t *testing.T) {
	c := NewCache()

	wg := sync.WaitGroup{}
	wg.Add(2)

	ids := make(map[string]int)
	for i := 1; i <= 100; i++ {
		ids[fmt.Sprintf("%d", i)] = 0
	}

	// producer
	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i++ {
			c.Insert(types.Device{ID: fmt.Sprintf("%d", i)})
		}
	}()

	// consumer
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				for id := range ids {
					if c.Delete(id) != nil {
						ids[id]++
					}
				}
			}
		}
	}()
	wg.Wait()

	c.mu.RLock()
	gotSize := len(c.data)
	c.mu.RUnlock()

	for id, count := range ids {
		if count > 1 {
			t.Fatalf("id %s was deleted %d times\n", id, count)
		}
	}

	if diff := cmp.Diff(0, gotSize); diff != "" {
		t.Fatalf("Item size mismatch (-want +got): %s\n", diff)
	}
}

func TestCacheConcurrencyWithMultiProducers(t *testing.T) {
	c := NewCache()
	wanted := 1000
	var wg sync.WaitGroup
	wg.Add(wanted)
	for i := 1; i <= wanted; i++ {
		i := i
		go func() {
			defer wg.Done()
			c.Insert(types.Device{ID: fmt.Sprintf("%d", i)})
		}()

	}
	wg.Wait()

	c.mu.RLock()
	gotSize := len(c.data)
	c.mu.RUnlock()

	if diff := cmp.Diff(wanted, gotSize); diff != "" {
		t.Fatalf("Item size mismatch (-want +got): %s\n", diff)
	}
}
