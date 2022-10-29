package banner

import (
	"math"
	"sync"
)

func NewUCB1Selector() Selector {
	return &selector{
		unselectedBanners: make(map[string]*banner),
		selectedBanners:   make(map[string]*banner),
	}
}

type banner struct {
	selectsCount uint
	clicksCount  uint
	score        float64
}

// TODO: Optimize.
type selector struct {
	unselectedBanners map[string]*banner
	selectedBanners   map[string]*banner
	mu                sync.Mutex
	maxScoreBannerID  string
	maxScore          float64
	totalSelectsCount uint
}

func (s *selector) AddBanner(bannerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.selectedBanners[bannerID]; ok {
		return ErrBannerAlreadyExist
	}
	if _, ok := s.unselectedBanners[bannerID]; ok {
		return ErrBannerAlreadyExist
	}
	s.unselectedBanners[bannerID] = &banner{}
	return nil
}

func (s *selector) DeleteBanner(bannerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.unselectedBanners[bannerID]; ok {
		delete(s.unselectedBanners, bannerID)
		return nil
	}
	if _, ok := s.selectedBanners[bannerID]; ok {
		delete(s.selectedBanners, bannerID)
		if s.maxScoreBannerID == bannerID {
			s.maxScoreBannerID = s.getMaxScoreBannerID()
			if s.maxScoreBannerID != "" {
				s.maxScore = s.selectedBanners[s.maxScoreBannerID].score
			}
		}
		return nil
	}
	return ErrBannerNotFound
}

func (s *selector) SelectBanner() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.unselectedBanners) > 0 {
		var id string
		var b *banner
		for id, b = range s.unselectedBanners {
			break
		}
		delete(s.unselectedBanners, id)
		s.totalSelectsCount++
		b.selectsCount++
		b.score = s.calculateBannerScore(b)
		s.selectedBanners[id] = b
		s.updateMaxScore(id, b)
		return id, nil
	}
	if s.maxScoreBannerID == "" {
		return "", ErrEmptyBannersList
	}
	b := s.selectedBanners[s.maxScoreBannerID]
	b.selectsCount++
	s.totalSelectsCount++
	b.score = s.calculateBannerScore(b)
	s.updateMaxScore(s.maxScoreBannerID, b)
	return s.maxScoreBannerID, nil
}

func (s *selector) RegisterClickForBanner(bannerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if b, ok := s.selectedBanners[bannerID]; ok {
		b.clicksCount++
		b.score = s.calculateBannerScore(b)
		s.updateMaxScore(bannerID, b)
		return nil
	}
	if b, ok := s.unselectedBanners[bannerID]; ok {
		delete(s.unselectedBanners, bannerID)
		b.clicksCount++
		b.selectsCount++
		s.totalSelectsCount++
		b.score = s.calculateBannerScore(b)
		s.selectedBanners[bannerID] = b
		s.updateMaxScore(bannerID, b)
		return nil
	}
	return ErrBannerNotFound
}

func (s *selector) getMaxScoreBannerID() string {
	maxScore := 0.0
	maxScoreBannerID := ""
	for bannerID, b := range s.selectedBanners {
		if b.score > maxScore {
			maxScore = b.score
			maxScoreBannerID = bannerID
		}
	}
	return maxScoreBannerID
}

func (s *selector) calculateBannerScore(b *banner) float64 {
	bannerRatio := float64(b.clicksCount) / float64(b.selectsCount)
	return bannerRatio + math.Sqrt((2.0*math.Log(float64(s.totalSelectsCount)))/float64(b.selectsCount))
}

func (s *selector) updateMaxScore(id string, b *banner) {
	if id == s.maxScoreBannerID && b.score < s.maxScore { // Current max banner is no longer max.
		s.maxScoreBannerID = s.getMaxScoreBannerID()
		s.maxScore = s.selectedBanners[s.maxScoreBannerID].score
		return
	}
	if b.score > s.maxScore {
		s.maxScoreBannerID = id
		s.maxScore = b.score
	}
}
