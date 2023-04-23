package snowflake

// tests from https://github.com/godruoyi/go-snowflake/blob/master/snowflake_test.go

import (
	"sync"
	"testing"
	"time"
)

var id uint64

func BenchmarkID(b *testing.B) {
	epoch := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	gen := New[uint64](epoch, 0, 1*time.Millisecond, 41, 3, 19)

	for i := 0; i < b.N; i++ {
		id = gen.ID()
	}
}

func TestID(t *testing.T) {
	epoch := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	gen := New[uint64](epoch, 0, 1*time.Millisecond, 41, 3, 19)

	mp := make(map[uint64]struct{})
	for i := 0; i < 100000; i++ {
		id := gen.ID()
		if _, ok := mp[id]; ok {
			t.Error("ID repeated", id, i)
			break
		}
		mp[id] = struct{}{}
	}
}

func TestID_Multi(t *testing.T) {
	epoch := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	gen := New[uint64](epoch, 0, 1*time.Millisecond, 41, 3, 19)

	le := 100000
	ch := make(chan uint64, le)
	var wg sync.WaitGroup
	for i := 0; i < le; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := gen.ID()
			ch <- id
		}()
	}
	wg.Wait()
	close(ch)

	mp := make(map[uint64]struct{})
	for id := range ch {
		if _, ok := mp[id]; ok {
			t.Error("ID repeated", id)
			break
		}
		mp[id] = struct{}{}
	}
	if len(mp) != le {
		t.Error("map length not equal", le)
	}
}
