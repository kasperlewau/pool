package pool

import (
	"fmt"
	"testing"
)

var pooltests = []struct {
	workers int
	queue   int
}{
	{1, 10},
	{2, 20},
	{5, 50},
	{10, 100},
	{15, 150},
	{20, 200},
	{100, 1000},
}

func TestPool(t *testing.T) {
	for _, tt := range pooltests {
		title := fmt.Sprintf("%v workers & %v tasks", tt.workers, tt.queue)
		t.Run(title, func(t *testing.T) {
			p := New(tt.workers, tt.queue)
			for i := 0; i < tt.queue; i++ {
				p.Add(func() (int, error) { return i, nil })
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
	{100, 1000},
}

func BenchmarkPool(b *testing.B) {
	for _, tt := range poolbenchmarks {
		title := fmt.Sprintf("%v workers & %v tasks", tt.workers, tt.queue)
		b.Run(title, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p := New(tt.workers, tt.queue)
				for i := 0; i < tt.queue; i++ {
					p.Add(func() (int, error) { return i, nil })
				}
				p.Start()
				p.Wait()
			}
		})
	}
}
