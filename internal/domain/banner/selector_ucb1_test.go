package banner

import (
	"context"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_selector_AddBanner(t *testing.T) {
	t.Run("when banned is new then no error should be returned", func(t *testing.T) {
		s := NewUCB1Selector()
		err := s.AddBanner("1")
		require.NoError(t, err)
	})
	t.Run("when banner is already added then error should be returned", func(t *testing.T) {
		s := NewUCB1Selector()
		err := s.AddBanner("1")
		require.NoError(t, err)
		err = s.AddBanner("1")
		require.Error(t, err)
		require.ErrorIs(t, ErrBannerAlreadyExist, err)
	})
	t.Run("when banner is selected then error should be returned", func(t *testing.T) {
		s := NewUCB1Selector()
		err := s.AddBanner("1")
		require.NoError(t, err)
		id, err := s.SelectBanner()
		require.NoError(t, err)
		require.Equal(t, "1", id)
		err = s.AddBanner("1")
		require.Error(t, err)
		require.ErrorIs(t, ErrBannerAlreadyExist, err)
	})
}

func Test_selector_DeleteBanner(t *testing.T) {
	t.Run("when banner is not exist then error should be returned", func(t *testing.T) {
		s := NewUCB1Selector()
		err := s.DeleteBanner("1")
		require.Error(t, err)
		require.ErrorIs(t, ErrBannerNotFound, err)
	})
	t.Run("when banner is exist then no error should be returned", func(t *testing.T) {
		s := NewUCB1Selector()
		err := s.AddBanner("1")
		require.NoError(t, err)
		err = s.DeleteBanner("1")
		require.NoError(t, err)
	})
	t.Run("when banner previously selected then no error should be returned", func(t *testing.T) {
		s := NewUCB1Selector()
		err := s.AddBanner("1")
		require.NoError(t, err)
		id, err := s.SelectBanner()
		require.NoError(t, err)
		require.Equal(t, "1", id)
		err = s.DeleteBanner("1")
		require.NoError(t, err)
	})
}

func Test_selector_SelectBanner(t *testing.T) {
	t.Run("when no banners added then error should be returned", func(t *testing.T) {
		s := NewUCB1Selector()
		id, err := s.SelectBanner()
		require.Error(t, err)
		require.ErrorIs(t, ErrEmptyBannersList, err)
		require.Empty(t, id)
	})
	t.Run("every banner should be selected at least once", func(t *testing.T) {
		selectedBanners := map[string]uint{"1": 0, "2": 0, "3": 0, "4": 0}
		s := NewUCB1Selector()
		for bannerID := range selectedBanners {
			err := s.AddBanner(bannerID)
			require.NoError(t, err)
		}
		for range selectedBanners {
			bannerID, err := s.SelectBanner()
			require.NoError(t, err)
			require.NotEmpty(t, bannerID)
			selectedBanners[bannerID]++
		}
		require.Equal(t, map[string]uint{"1": 1, "2": 1, "3": 1, "4": 1}, selectedBanners)
	})
}

func Test_selector_RegisterClickForBanner(t *testing.T) {
	t.Run("when banner is not exist then error should be returned", func(t *testing.T) {
		s := NewUCB1Selector()
		err := s.RegisterClickForBanner("1")
		require.Error(t, err)
		require.ErrorIs(t, ErrBannerNotFound, err)
	})
	t.Run("click on unselected banner is not an error", func(t *testing.T) {
		s := NewUCB1Selector()
		err := s.AddBanner("1")
		require.NoError(t, err)
		err = s.RegisterClickForBanner("1")
		require.NoError(t, err)
	})
	t.Run("popular banners should be selected frequently", func(t *testing.T) {
		banners := make(map[string]float64) // bannerID to click-through ratio
		bannersCount := 100
		defaultCTR := 0.15 // default click-through ratio
		s := NewUCB1Selector()
		for i := 1; i <= bannersCount; i++ {
			id := strconv.Itoa(i)
			banners[id] = defaultCTR
			err := s.AddBanner(id)
			require.NoError(t, err)
		}
		popularBannerID := "100"
		subPopularBannerID := "10"
		banners[subPopularBannerID] = 0.5 // popular banner 1
		banners[popularBannerID] = 0.7    // popular banner 2
		bannerSelectTimes := 20
		selectsCount := bannersCount * bannerSelectTimes
		workersCount := 5
		totalSelects := workersCount * selectsCount
		selectedBanners := make(map[string]uint)
		selectedBannersCh := make(chan string, totalSelects)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			for bannerID := range selectedBannersCh {
				selectedBanners[bannerID]++
			}
			cancel()
		}()
		wg := sync.WaitGroup{}
		wg.Add(workersCount)
		for i := 0; i < workersCount; i++ {
			go func(s *Selector, wg *sync.WaitGroup) {
				defer wg.Done()
				for j := 0; j < selectsCount; j++ {
					bannerID, err := (*s).SelectBanner()
					require.NoError(t, err)
					require.NotEmpty(t, bannerID)
					selectedBannersCh <- bannerID
					//nolint:gosec
					// This is not critical part, so weak numbers generator is allowed
					if rand.Float64() < banners[bannerID] {
						err = (*s).RegisterClickForBanner(bannerID)
						require.NoError(t, err)
					}
				}
			}(&s, &wg)
		}
		wg.Wait()
		close(selectedBannersCh)
		<-ctx.Done()

		selectsPercent := make([]float64, bannersCount)
		i := 0
		for _, count := range selectedBanners {
			selectsPercent[i] = (float64(count) / float64(totalSelects)) * 100.0
			i++
		}
		sort.Float64s(selectsPercent)
		medianPercent := selectsPercent[bannersCount/2]
		popularBannerPercent := (float64(selectedBanners["100"]) / float64(totalSelects)) * 100.0
		popularToMedianRatio := popularBannerPercent / medianPercent

		require.GreaterOrEqual(t, popularBannerPercent, medianPercent)
		require.GreaterOrEqual(t, popularToMedianRatio, float64(bannerSelectTimes))
	})
}

func TestNewUCB1Selector(t *testing.T) {
	s := NewUCB1Selector()
	require.NotNil(t, s)
}
