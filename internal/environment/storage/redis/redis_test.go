//go:build integration

package redis

import (
	"context"
	"fmt"
	"github.com/ekhvalov/otus-banners-rotation/internal/app"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ekhvalov/otus-banners-rotation/internal/environment/storage/redis/mock"
	rediscli "github.com/go-redis/redis/v9"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

const (
	defaultRedisHost     = "localhost"
	defaultRedisPort     = "6379"
	defaultRedisUsername = ""
	defaultRedisPassword = ""
	defaultRedisDatabase = "10"
)

func TestRedisStorage(t *testing.T) {
	suite.Run(t, new(redisSuite))
}

type redisSuite struct {
	suite.Suite
	ctx     context.Context
	cancel  context.CancelFunc
	tick    time.Duration
	waitFor time.Duration
	cfg     Config
	client  *rediscli.Client
}

func (s *redisSuite) SetupSuite() {
	s.tick = time.Millisecond * 100
	s.waitFor = s.tick * 1000 * 30
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.cfg = makeConfig()
	s.client = rediscli.NewClient(&rediscli.Options{
		Addr:     s.cfg.GetAddress(),
		Username: s.cfg.GetUsername(),
		Password: s.cfg.GetPassword(),
		DB:       s.cfg.GetDatabase(),
	})
}

func (s *redisSuite) SetupTest() {
	s.flushDB()
}

func (s *redisSuite) TearDownTest() {
	s.flushDB()
}

func (s *redisSuite) Test_CreateBanner() {
	bannerID := "100500"
	description := "Banner description"
	r := Redis{cfg: s.cfg, client: s.client, idGenerator: s.createIDGeneratorMock(bannerID)}

	gotBannerID, err := r.CreateBanner(s.ctx, description)

	s.Require().NoError(err)
	s.Require().Equal(bannerID, gotBannerID)
	s.Require().Equal(s.hGet(keyBanners, bannerID), description)
}

func (s *redisSuite) Test_DeleteBanner() {
	bannerID := "100500"
	s.hSet(keyBanners, bannerID, "description")
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.DeleteBanner(s.ctx, bannerID)

	s.Require().NoError(err)
	err = s.client.HGet(s.ctx, keyBanners, bannerID).Err()
	s.Require().ErrorIs(rediscli.Nil, err)
}

func (s *redisSuite) Test_CreateSlot() {
	slotID := "100500"
	description := "slot description"
	r := NewRedis(s.cfg, s.createIDGeneratorMock(slotID))

	gotSlotID, err := r.CreateSlot(s.ctx, description)

	s.Require().NoError(err)
	s.Require().Equal(slotID, gotSlotID)
	s.Require().Equal(s.hGet(keySlots, slotID), description)
}

func (s *redisSuite) Test_DeleteSlot() {
	slotID := "100500"
	s.hSet(keySlots, slotID, "description")
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.DeleteSlot(s.ctx, slotID)

	s.Require().NoError(err)
	err = s.client.HGet(s.ctx, keySlots, slotID).Err()
	s.Require().ErrorIs(rediscli.Nil, err)
}

func (s *redisSuite) Test_CreateSocialGroup() {
	socialGroupID := "100500"
	description := "socialGroup description"
	r := NewRedis(s.cfg, s.createIDGeneratorMock(socialGroupID))

	gotSocialGroupID, err := r.CreateSocialGroup(s.ctx, description)

	s.Require().NoError(err)
	s.Require().Equal(socialGroupID, gotSocialGroupID)
	s.Require().Equal(s.hGet(keySocialGroups, socialGroupID), description)
}

func (s *redisSuite) Test_DeleteSocialGroup() {
	socialGroupID := "100500"
	s.seedSocialGroup(socialGroupID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.DeleteSocialGroup(s.ctx, socialGroupID)

	s.Require().NoError(err)
	err = s.client.HGet(s.ctx, keySocialGroups, socialGroupID).Err()
	s.Require().ErrorIs(rediscli.Nil, err)
}

func (s *redisSuite) Test_AttachBanner() {
	bannerID := "100500"
	s.seedBanner(bannerID)
	slotID := "100600"
	s.seedSlot(slotID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.AttachBanner(s.ctx, slotID, bannerID)

	s.Require().NoError(err)
	bannerScore, err := s.client.ZScore(s.ctx, makeSlotBannersKey(slotID), bannerID).Result()
	s.Require().NoError(err)
	s.Require().True(math.IsInf(bannerScore, 1))
}

func (s *redisSuite) Test_AttachBanner_Error_BannerNotFound() {
	bannerID := "100500"
	slotID := "100600"
	s.seedSlot(slotID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.AttachBanner(s.ctx, slotID, bannerID)

	s.Require().Error(err)
	var errNotFound *app.ErrNotFound
	s.Require().ErrorAs(err, &errNotFound)
}

func (s *redisSuite) Test_AttachBanner_Error_SlotNotFound() {
	bannerID := "100500"
	s.seedBanner(bannerID)
	slotID := "100600"
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.AttachBanner(s.ctx, slotID, bannerID)

	s.Require().Error(err)
	var errNotFound *app.ErrNotFound
	s.Require().ErrorAs(err, &errNotFound)
}

func (s *redisSuite) Test_DetachBanner() {
	bannerID := "100500"
	s.seedBanner(bannerID)
	slotID := "100600"
	s.seedSlot(slotID)
	s.attachBanner(slotID, bannerID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.DetachBanner(s.ctx, slotID, bannerID)

	s.Require().NoError(err)
	err = s.client.ZScore(s.ctx, makeSlotBannersKey(slotID), bannerID).Err()
	s.Require().ErrorIs(rediscli.Nil, err)
}

func (s *redisSuite) Test_DetachBanner_Error_SlotNotFound() {
	bannerID := "100500"
	s.seedBanner(bannerID)
	slotID := "100600"
	s.attachBanner(slotID, bannerID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.DetachBanner(s.ctx, slotID, bannerID)

	s.Require().Error(err)
	var errNotFound *app.ErrNotFound
	s.Require().ErrorAs(err, &errNotFound)
}

func (s *redisSuite) Test_DetachBanner_Error_BannerNotFound() {
	bannerID := "100500"
	//s.seedBanner(bannerID)
	slotID := "100600"
	s.seedSlot(slotID)
	s.attachBanner(slotID, bannerID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.DetachBanner(s.ctx, slotID, bannerID)

	s.Require().Error(err)
	var errNotFound *app.ErrNotFound
	s.Require().ErrorAs(err, &errNotFound)
}

func (s *redisSuite) Test_DetachBanner_Error_BannerNotAttached() {
	bannerID := "100500"
	s.seedBanner(bannerID)
	slotID := "100600"
	s.seedSlot(slotID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.DetachBanner(s.ctx, slotID, bannerID)

	s.Require().Error(err)
	var errBannerNotAttached *app.ErrBannerNotAttached
	s.Require().ErrorAs(err, &errBannerNotAttached)
}

func (s *redisSuite) Test_SelectBanner() {
	slotID := "100600"
	s.seedSlot(slotID)
	bannerIDs := []string{"100501", "100502", "100503"}
	for _, id := range bannerIDs {
		s.seedBanner(id)
		s.attachBanner(slotID, id)
	}
	socialGroupID := "100700"
	s.seedSocialGroup(socialGroupID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	selectedBannerIDs := make([]string, len(bannerIDs))
	for i := 0; i < len(selectedBannerIDs); i++ {
		bannerID, err := r.SelectBanner(s.ctx, slotID, socialGroupID)
		s.Require().NoError(err)
		selectedBannerIDs[i] = bannerID
	}

	s.Require().ElementsMatch(selectedBannerIDs, bannerIDs)
	totalSelects := s.getInt(makeSlotSocialGroupSelectsTotalKey(slotID, socialGroupID))
	s.Require().Equal(len(selectedBannerIDs), totalSelects)
	for _, bannerID := range selectedBannerIDs {
		selects := s.hGetInt(makeSlotSocialGroupSelectsKey(slotID, socialGroupID), bannerID)
		s.Require().Equal(1, selects)
		score := s.zScore(makeSlotSocialGroupScoresKey(slotID, socialGroupID), bannerID)
		s.Require().Less(score, math.Inf(1))
		s.Require().GreaterOrEqual(score, 0.0)
	}
}

func (s *redisSuite) Test_SelectBanner_Error_SlotNotFound() {
	slotID := "100600"
	socialGroupID := "100700"
	s.seedSocialGroup(socialGroupID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	bannerID, err := r.SelectBanner(s.ctx, slotID, socialGroupID)
	s.Require().Error(err)
	var errNotFound *app.ErrNotFound
	s.Require().ErrorAs(err, &errNotFound)
	s.Require().Empty(bannerID)
}

func (s *redisSuite) Test_SelectBanner_Error_SocialGroupNotFound() {
	slotID := "100600"
	s.seedSlot(slotID)
	socialGroupID := "100700"
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	bannerID, err := r.SelectBanner(s.ctx, slotID, socialGroupID)
	s.Require().Error(err)
	var errNotFound *app.ErrNotFound
	s.Require().ErrorAs(err, &errNotFound)
	s.Require().Empty(bannerID)
}

func (s *redisSuite) Test_SelectBanner_Error_NoBannersFound() {
	slotID := "100600"
	s.seedSlot(slotID)
	bannerID := "100500"
	s.seedBanner(bannerID)
	socialGroupID := "100700"
	s.seedSocialGroup(socialGroupID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	bannerID, err := r.SelectBanner(s.ctx, slotID, socialGroupID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, app.ErrNoBannersFound)
	s.Require().Empty(bannerID)
}

func (s *redisSuite) Test_ClickBanner() {
	slotID := "100600"
	s.seedSlot(slotID)
	bannerID := "100500"
	s.seedBanner(bannerID)
	s.attachBanner(slotID, bannerID)
	socialGroupID := "100700"
	s.seedSocialGroup(socialGroupID)
	selectsKey := makeSlotSocialGroupSelectsKey(slotID, socialGroupID)
	err := s.client.HIncrBy(s.ctx, selectsKey, bannerID, 1).Err()
	s.Require().NoError(err)
	clicksKey := makeSlotSocialGroupClicksKey(slotID, socialGroupID)
	err = s.client.HIncrBy(s.ctx, clicksKey, bannerID, 0).Err()
	s.Require().NoError(err)
	totalSelectsKey := makeSlotSocialGroupSelectsTotalKey(slotID, socialGroupID)
	err = s.client.IncrBy(s.ctx, totalSelectsKey, 1).Err()
	s.Require().NoError(err)
	scoresKey := makeSlotSocialGroupScoresKey(slotID, socialGroupID)
	err = s.client.ZIncrBy(s.ctx, scoresKey, 0.0, bannerID).Err()
	s.Require().NoError(err)

	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))
	err = r.ClickBanner(s.ctx, slotID, bannerID, socialGroupID)

	s.Require().NoError(err)
	clicks := s.hGetInt(clicksKey, bannerID)
	s.Require().Equal(1, clicks, "clicks are invalid")
	score := s.zScore(scoresKey, bannerID)
	s.Require().Equal(calculateBannerScore(1, 1, 1), score, "score is invalid")
}

func (s *redisSuite) Test_ClickBanner_Error_BannerNotFound() {
	slotID := "100600"
	s.seedSlot(slotID)
	bannerID := "100500"
	socialGroupID := "100700"
	s.seedSocialGroup(socialGroupID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.ClickBanner(s.ctx, slotID, bannerID, socialGroupID)

	s.Require().Error(err)
	var errNotFound *app.ErrNotFound
	s.Require().ErrorAs(err, &errNotFound)
}

func (s *redisSuite) Test_ClickBanner_Error_BannerNotAttached() {
	slotID := "100600"
	s.seedSlot(slotID)
	bannerID := "100500"
	s.seedBanner(bannerID)
	socialGroupID := "100700"
	s.seedSocialGroup(socialGroupID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.ClickBanner(s.ctx, slotID, bannerID, socialGroupID)

	s.Require().Error(err)
	var errBannerNotAttached *app.ErrBannerNotAttached
	s.Require().ErrorAs(err, &errBannerNotAttached)
}

func (s *redisSuite) Test_ClickBanner_Error_SlotNotFound() {
	slotID := "100600"
	bannerID := "100500"
	s.seedBanner(bannerID)
	socialGroupID := "100700"
	s.seedSocialGroup(socialGroupID)
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.ClickBanner(s.ctx, slotID, bannerID, socialGroupID)

	s.Require().Error(err)
	var errNotFound *app.ErrNotFound
	s.Require().ErrorAs(err, &errNotFound)
}

func (s *redisSuite) Test_ClickBanner_Error_SocialGroupNotFound() {
	slotID := "100600"
	s.seedSlot(slotID)
	bannerID := "100500"
	s.seedBanner(bannerID)
	s.attachBanner(slotID, bannerID)
	socialGroupID := "100700"
	r := NewRedis(s.cfg, s.createIDGeneratorMock(""))

	err := r.ClickBanner(s.ctx, slotID, bannerID, socialGroupID)

	s.Require().Error(err)
	var errNotFound *app.ErrNotFound
	s.Require().ErrorAs(err, &errNotFound)
}

func (s *redisSuite) Test_SelectClickBanner_Concurrent() {
	slotID := "100600"
	s.seedSlot(slotID)
	socialGroupID := "100700"
	s.seedSocialGroup(socialGroupID)
	bannersCount := 100
	bannersClicksRatio := make(map[string]float64, bannersCount)
	bannersClicks := make(map[string]int)
	bannersSelects := make(map[string]int)
	for i := 0; i < bannersCount; i++ {
		bannerID := fmt.Sprintf("banner-%d", i)
		s.seedBanner(bannerID)
		s.attachBanner(slotID, bannerID)
		bannersClicksRatio[bannerID] = 0.15
		bannersClicks[bannerID] = 0
		bannersSelects[bannerID] = 0
	}
	mostPopularBannerID := "banner-0"
	bannersClicksRatio[mostPopularBannerID] = 0.7
	subPopularBannerID := "banner-1"
	bannersClicksRatio[subPopularBannerID] = 0.5

	workersCount := 5
	selectsPerWorker := 1000

	bannersClicksCh := make(chan string, workersCount*selectsPerWorker)
	bannersSelectsCh := make(chan string, workersCount*selectsPerWorker)
	wg := sync.WaitGroup{}
	r := NewRedis(s.cfg, NewUUIDGenerator())
	randomSeed := time.Now().UnixNano()
	fmt.Println("random seed: ", randomSeed)
	rand.Seed(randomSeed)

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < selectsPerWorker; j++ {
				bannerID, err := r.SelectBanner(s.ctx, slotID, socialGroupID)
				s.Require().NoError(err)
				bannersSelectsCh <- bannerID
				if rand.Float64() < bannersClicksRatio[bannerID] {
					err := r.ClickBanner(s.ctx, slotID, bannerID, socialGroupID)
					s.Require().NoError(err)
					bannersClicksCh <- bannerID
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
	close(bannersClicksCh)
	close(bannersSelectsCh)

	for bannerId := range bannersClicksCh {
		bannersClicks[bannerId]++
	}
	for bannerId := range bannersSelectsCh {
		bannersSelects[bannerId]++
	}

	selectsTotal := s.getInt(makeSlotSocialGroupSelectsTotalKey(slotID, socialGroupID))
	s.Require().Equal(workersCount*selectsPerWorker, selectsTotal, "total selects mismatched")
	clicksKey := makeSlotSocialGroupClicksKey(slotID, socialGroupID)
	selectsKey := makeSlotSocialGroupSelectsKey(slotID, socialGroupID)
	scoresKey := makeSlotSocialGroupScoresKey(slotID, socialGroupID)
	selectsRatios := make([]float64, 0)
	for bannerID := range bannersClicksRatio {
		expectedClicks := bannersClicks[bannerID]
		actualClicks := s.hGetIntOrDefault(clicksKey, bannerID, 0)
		s.Require().Equal(expectedClicks, actualClicks, fmt.Sprintf("clicks mismatched (%s)", bannerID))

		expectedSelects := bannersSelects[bannerID]
		actualSelects := s.hGetInt(selectsKey, bannerID)
		s.Require().GreaterOrEqual(actualSelects, 1) // Every banner has selected at least once
		s.Require().Equal(expectedSelects, actualSelects, fmt.Sprintf("selects mismatched (%s)", bannerID))
		selectsRatios = append(selectsRatios, float64(actualSelects)/float64(selectsTotal))

		actualScore := s.zScore(scoresKey, bannerID)
		s.Require().NotEqual(math.Inf(1), actualScore)
		s.Require().GreaterOrEqual(actualScore, 0.0)
	}
	sort.Float64s(selectsRatios)
	medianRatio := selectsRatios[bannersCount/2]
	mostPopularBannerRatio := float64(bannersSelects[mostPopularBannerID]) / float64(selectsTotal)
	s.Require().Equal(selectsRatios[bannersCount-1], mostPopularBannerRatio)
	s.Require().GreaterOrEqual(mostPopularBannerRatio, medianRatio*15.0)

	subPopularBannerRatio := float64(bannersSelects[subPopularBannerID]) / float64(selectsTotal)
	s.Require().Equal(selectsRatios[bannersCount-2], subPopularBannerRatio)
	s.Require().GreaterOrEqual(subPopularBannerRatio, medianRatio*2.0)

}

func (s *redisSuite) flushDB() {
	err := s.client.FlushDB(s.ctx).Err()
	s.Require().NoError(err)
}

func (s *redisSuite) createIDGeneratorMock(id string) IDGenerator {
	controller := gomock.NewController(s.T())
	idGenerator := mock.NewMockIDGenerator(controller)
	if id != "" {
		idGenerator.EXPECT().GenerateID().Return(id)
	}
	return idGenerator
}

func (s *redisSuite) seedBanner(bannerID string) {
	s.hSet(keyBanners, bannerID, fmt.Sprintf("banner %s description", bannerID))
}

func (s *redisSuite) seedSlot(slotID string) {
	s.hSet(keySlots, slotID, fmt.Sprintf("slot %s description", slotID))
}

func (s *redisSuite) seedSocialGroup(socialGroupID string) {
	s.hSet(keySocialGroups, socialGroupID, fmt.Sprintf("social group %s description", socialGroupID))
}

func (s *redisSuite) attachBanner(slotID, bannerID string) {
	err := s.client.ZAdd(s.ctx, makeSlotBannersKey(slotID), rediscli.Z{
		Score:  math.Inf(1),
		Member: bannerID,
	}).Err()
	s.Require().NoError(err)
}

func (s *redisSuite) hSet(key, field, value string) {
	err := s.client.HSet(s.ctx, key, field, value).Err()
	s.Require().NoError(err)
}

func (s *redisSuite) hGet(key, field string) string {
	result, err := s.client.HGet(s.ctx, key, field).Result()
	s.Require().NoError(err)
	return result
}

func (s *redisSuite) hGetInt(key, field string) int {
	cmd := s.client.HGet(s.ctx, key, field)
	s.Require().NoError(cmd.Err())
	value, err := cmd.Int()
	s.Require().NoError(err)
	return value
}
func (s *redisSuite) hGetIntOrDefault(key, field string, defaultValue int) int {
	cmd := s.client.HGet(s.ctx, key, field)
	if cmd.Err() != nil && cmd.Err() == rediscli.Nil {
		return defaultValue
	}
	s.Require().NoError(cmd.Err())
	value, err := cmd.Int()
	s.Require().NoError(err)
	return value
}

func (s *redisSuite) hGetFloat64(key, field string) float64 {
	cmd := s.client.HGet(s.ctx, key, field)
	s.Require().NoError(cmd.Err())
	value, err := cmd.Float64()
	s.Require().NoError(err)
	return value
}

func (s *redisSuite) getInt(key string) int {
	cmd := s.client.Get(s.ctx, key)
	s.Require().NoError(cmd.Err())
	value, err := cmd.Int()
	s.Require().NoError(err)
	return value
}

func (s *redisSuite) zScore(key, field string) float64 {
	score, err := s.client.ZScore(s.ctx, key, field).Result()
	s.Require().NoError(err)
	return score
}

func makeConfig() Config {
	v := viper.New()
	v.SetEnvPrefix("TESTS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	defaults := map[string]string{
		"TESTS_REDIS_HOST":     defaultRedisHost,
		"TESTS_REDIS_PORT":     defaultRedisPort,
		"TESTS_REDIS_USERNAME": defaultRedisUsername,
		"TESTS_REDIS_PASSWORD": defaultRedisPassword,
		"TESTS_REDIS_DATABASE": defaultRedisDatabase,
	}
	for envName, value := range defaults {
		if _, ok := os.LookupEnv(envName); !ok {
			err := os.Setenv(envName, value)
			if err != nil {
				panic(fmt.Errorf("set env error: %w", err))
			}
		}
	}
	return NewConfig(v)
}
