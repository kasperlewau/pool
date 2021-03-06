package pool_test

import (
	"fmt"
	"github.com/kasperlewau/pool"
	"testing"
	"time"
)

var pooltests = []struct {
	workers int
	queue   int
}{
	{1, 10},
	{20, 200},
	{50, 500},
	{100, 1000},
	{500, 5000},
}

func TestPool(t *testing.T) {
	for _, tt := range pooltests {
		title := fmt.Sprintf("%v workers & %v tasks", tt.workers, tt.queue)
		t.Run(title, func(t *testing.T) {
			p := pool.New(tt.workers, tt.queue)
			for i := 0; i < tt.queue; i++ {
				p.Add(func() (interface{}, error) { return i, nil })
			}
			time.Sleep(100 * time.Millisecond)
			if len(p.Workers) != tt.workers {
				t.Errorf("Expected pool to have spun up %v workers. have %v", tt.workers, len(p.Workers))
			}
			if len(p.Tasks) != tt.queue {
				t.Errorf("Expected pool to have queued %v tasks. have %v", tt.queue, len(p.Tasks))
			}
			p.Start()
			p.Wait()
		})
	}
}

var poolbenchmarks = []struct {
	workers int
	queue   int
}{
	{10, 100},
}

func BenchmarkPool(b *testing.B) {
	for _, tt := range poolbenchmarks {
		title := fmt.Sprintf("%v workers & %v tasks", tt.workers, tt.queue)
		b.Run(title, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				p := pool.New(tt.workers, tt.queue)
				for i := 0; i < tt.queue; i++ {
					p.Add(func() (interface{}, error) { return i, nil })
				}
				p.Start()
				p.Wait()
			}
		})
	}
}
