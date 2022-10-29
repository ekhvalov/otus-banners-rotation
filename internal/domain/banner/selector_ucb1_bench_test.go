package banner

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func Benchmark_selector_complex(b *testing.B) {
	rand.Seed(10)
	banners := make(map[string]float64)
	defaultCTR := 0.25
	bannersCount := 100
	for i := 0; i < bannersCount; i++ {
		id := strconv.Itoa(i)
		banners[id] = defaultCTR
	}
	banners["39"] = 0.4
	banners["59"] = 0.5
	banners["99"] = 0.7
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		s := NewUCB1Selector()
		for id := range banners {
			_ = s.AddBanner(id)
		}
		for j := 0; j < 5; j++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, s *Selector) {
				defer wg.Done()
				for w := 0; w < bannersCount*10; w++ {
					if id, _ := (*s).SelectBanner(); id != "" {
						if rand.Float64() < banners[id] {
							_ = (*s).RegisterClickForBanner(id)
						}
					}
				}
			}(&wg, &s)
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, s *Selector) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				for k := 0; k < 10; k++ {
					id := strconv.Itoa(k*10 + j)
					_ = (*s).DeleteBanner(id)
				}
			}
		}(&wg, &s)
		wg.Add(1)
		go func(wg *sync.WaitGroup, s *Selector) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				for k := 0; k < 10; k++ {
					id := strconv.Itoa(k*10 + j)
					_ = (*s).AddBanner(id)
				}
			}
		}(&wg, &s)
		wg.Wait()
	}
}
