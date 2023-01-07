package muchan_test

import (
	"sync"
	"testing"
)

// goos: darwin
// goarch: arm64
// BenchmarkMumap-8                15407637               345.2 ns/op            82 B/op          1 allocs/op
// BenchmarkChamap-8                3611944              3102 ns/op             356 B/op          2 allocs/op
// BenchmarkChamapMethods-8         7121342              4815 ns/op             193 B/op          1 allocs/op

type muMap struct {
	m  map[int]int
	mu *sync.Mutex
}

func BenchmarkMumap(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(b.N)
	m := muMap{
		m:  map[int]int{},
		mu: new(sync.Mutex),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			m.mu.Lock()
			m.m[i] = i
			m.mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
}

type chaMap struct {
	m  map[int]int
	ch chan struct{}
}

func BenchmarkChamap(b *testing.B) {
	m := chaMap{
		m:  map[int]int{},
		ch: make(chan struct{}, 1),
	}
	wg := sync.WaitGroup{}
	wg.Add(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			m.ch <- struct{}{}
			m.m[i] = i
			<-m.ch
			wg.Done()
		}()
	}
	wg.Wait()
}

func (c *chaMap) Lock() {
	c.ch <- struct{}{}
}

func (c *chaMap) Unlock() {
	<-c.ch
}

func BenchmarkChamapMethods(b *testing.B) {
	m := chaMap{
		m:  map[int]int{},
		ch: make(chan struct{}, 1),
	}
	wg := sync.WaitGroup{}
	wg.Add(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			m.Lock()
			m.m[i] = i
			m.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
}
