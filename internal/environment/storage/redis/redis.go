package redis

import (
	"context"
	"errors"
	"fmt"
	"math"

	rediscli "github.com/go-redis/redis/v9"
)

const (
	keyBanners      = "banners"
	keySlots        = "slots"
	keySocialGroups = "social_groups"
)

func NewRedis(config Config, idGenerator IDGenerator) *Redis {
	cli := rediscli.NewClient(&rediscli.Options{
		Addr:     config.GetAddress(),
		Password: config.GetPassword(),
		DB:       config.GetDatabase(),
	})
	return &Redis{client: cli, idGenerator: idGenerator}
}

type Redis struct {
	client      *rediscli.Client
	idGenerator IDGenerator
}

func (r *Redis) CreateBanner(ctx context.Context, description string) (id string, err error) {
	id = r.idGenerator.GenerateID()
	if err = r.hSet(ctx, keyBanners, id, description); err != nil {
		return "", err
	}
	return id, nil
}

func (r *Redis) DeleteBanner(ctx context.Context, id string) error {
	return r.hDel(ctx, keyBanners, id)
}

func (r *Redis) CreateSlot(ctx context.Context, description string) (id string, err error) {
	id = r.idGenerator.GenerateID()
	if err = r.hSet(ctx, keySlots, id, description); err != nil {
		return "", err
	}
	return id, nil
}

func (r *Redis) DeleteSlot(ctx context.Context, id string) error {
	return r.hDel(ctx, keySlots, id)
}

func (r *Redis) CreateSocialGroup(ctx context.Context, description string) (id string, err error) {
	id = r.idGenerator.GenerateID()
	if err = r.hSet(ctx, keySocialGroups, id, description); err != nil {
		return "", err
	}
	return id, nil
}

func (r *Redis) DeleteSocialGroup(ctx context.Context, id string) error {
	return r.hDel(ctx, keySocialGroups, id)
}

func (r *Redis) AttachBanner(ctx context.Context, slotID, bannerID string) error {
	return r.zAdd(ctx, makeSlotBannersKey(slotID), bannerID, math.Inf(1))
}

func (r *Redis) DetachBanner(ctx context.Context, slotID, bannerID string) error {
	return r.zRem(ctx, makeSlotBannersKey(slotID), bannerID)
}

func (r *Redis) SelectBanner(ctx context.Context, slotID, socialGroupID string) (bannerID string, err error) {
	scoresKey := makeSlotSocialGroupScoresKey(slotID, socialGroupID)
	keyCount, err := r.client.Exists(ctx, scoresKey).Result()
	if err != nil {
		return "", fmt.Errorf("exists of '%s' error: %w", scoresKey, err)
	}
	if keyCount == 0 {
		slotBannersKey := makeSlotBannersKey(slotID)
		err = r.client.Copy(ctx, slotBannersKey, scoresKey, 0, false).Err() // TODO: Get DB from config
		if err != nil {
			return "", fmt.Errorf("copy of '%s' to '%s' error: %w", slotBannersKey, scoresKey, err)
		}
	}
	bannerIDs, err := r.client.ZRevRange(ctx, scoresKey, 0, 0).Result()
	if err != nil {
		return "", fmt.Errorf("zrevrange of '%s' error: %w", scoresKey, err)
	}
	if len(bannerIDs) == 0 {
		return "", fmt.Errorf("no banners found") // TODO: Define special error
	}
	bannerID = bannerIDs[0]
	selectsKey := makeSlotSocialGroupSelectsKey(slotID, socialGroupID)
	selects, err := r.client.HIncrBy(ctx, selectsKey, bannerID, 1).Result()
	if err != nil {
		return "", fmt.Errorf("hincrby of '%s' error: %w", selectsKey, err)
	}
	clicksKey := makeSlotSocialGroupClicksKey(slotID, socialGroupID)
	clicks, err := r.hGetInt64OrDefault(ctx, clicksKey, bannerID, 0)
	if err != nil {
		return "", err
	}
	totalSelectsKey := makeSlotSocialGroupSelectsTotalKey(slotID, socialGroupID)
	totalSelects, err := r.client.Incr(ctx, totalSelectsKey).Result()
	if err != nil {
		return "", fmt.Errorf("incrby of '%s' error: %w", totalSelectsKey, err)
	}
	score := calculateBannerScore(float64(selects), float64(clicks), float64(totalSelects))
	if err = r.zAdd(ctx, scoresKey, bannerID, score); err != nil {
		return "", fmt.Errorf("zincrby of '%s' error: %w", scoresKey, err)
	}
	return bannerID, nil
}

func (r *Redis) ClickBanner(ctx context.Context, slotID, bannerID, socialGroupID string) error {
	totalSelectsKey := makeSlotSocialGroupSelectsTotalKey(slotID, socialGroupID)
	totalSelects, err := r.getInt64OrDefault(ctx, totalSelectsKey, 1)
	if err != nil {
		return err
	}
	selectsKey := makeSlotSocialGroupSelectsKey(slotID, socialGroupID)
	selects, err := r.hGetInt64OrDefault(ctx, selectsKey, bannerID, 1)
	if err != nil {
		return err
	}
	clicksKey := makeSlotSocialGroupClicksKey(slotID, socialGroupID)
	clicks, err := r.client.HIncrBy(ctx, clicksKey, bannerID, 1).Result()
	if err != nil {
		return fmt.Errorf("hincrby of '%s' '%s' error: %w", clicksKey, bannerID, err)
	}
	scoresKey := makeSlotSocialGroupScoresKey(slotID, socialGroupID)
	score := calculateBannerScore(float64(selects), float64(clicks), float64(totalSelects))
	if err = r.zAdd(ctx, scoresKey, bannerID, score); err != nil {
		return fmt.Errorf("zadd of '%s' error: %w", scoresKey, err)
	}
	return nil
}

func (r *Redis) hGetInt64OrDefault(ctx context.Context, key, field string, defaultValue int64) (int64, error) {
	cmd := r.client.HGet(ctx, key, field)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), rediscli.Nil) {
			return defaultValue, nil
		}
		return 0, fmt.Errorf("hget of '%s' '%s' error: %w", key, field, cmd.Err())
	}
	value, err := cmd.Int64()
	if err != nil {
		return 0, fmt.Errorf("hget of '%s' '%s' parse int64 error: %w", key, field, err)
	}
	return value, nil
}

func (r *Redis) getInt64OrDefault(ctx context.Context, key string, defaultValue int64) (int64, error) {
	cmd := r.client.Get(ctx, key)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), rediscli.Nil) {
			return defaultValue, nil
		}
		return 0, fmt.Errorf("get of '%s' error: %w", key, cmd.Err())
	}
	value, err := cmd.Int64()
	if err != nil {
		return 0, fmt.Errorf("get of '%s' parse int error: %w", key, err)
	}
	return value, nil
}

func (r *Redis) hSet(ctx context.Context, key, field, value string) error {
	if err := r.client.HSet(ctx, key, field, value).Err(); err != nil {
		return fmt.Errorf("hset of '%s' '%s' error: %w", key, field, err)
	}
	return nil
}

func (r *Redis) hDel(ctx context.Context, key, field string) error {
	if err := r.client.HDel(ctx, key, field).Err(); err != nil {
		return fmt.Errorf("hdel of '%s' '%s' error: %w", key, field, err)
	}
	return nil
}

func (r *Redis) zAdd(ctx context.Context, key, field string, score float64) error {
	z := rediscli.Z{Score: score, Member: field}
	if err := r.client.ZAdd(ctx, key, z).Err(); err != nil {
		return fmt.Errorf("zadd of '%s' '%s' error: %w", key, field, err)
	}
	return nil
}

func (r *Redis) zRem(ctx context.Context, key, field string) error {
	if err := r.client.ZRem(ctx, key, field).Err(); err != nil {
		return fmt.Errorf("zrem of '%s' '%s' error: %w", key, field, err)
	}
	return nil
}

func makeSlotBannersKey(slotID string) string {
	return fmt.Sprintf("slot:%s:banners", slotID)
}

func makeSlotSocialGroupSelectsKey(slotID, socialGroupID string) string {
	return makeSlotSocialGroupKey(slotID, socialGroupID, "selects")
}

func makeSlotSocialGroupSelectsTotalKey(slotID, socialGroupID string) string {
	return makeSlotSocialGroupKey(slotID, socialGroupID, "selects_total")
}

func makeSlotSocialGroupClicksKey(slotID, socialGroupID string) string {
	return makeSlotSocialGroupKey(slotID, socialGroupID, "clicks")
}

func makeSlotSocialGroupScoresKey(slotID, socialGroupID string) string {
	return makeSlotSocialGroupKey(slotID, socialGroupID, "scores")
}

func makeSlotSocialGroupKey(slotID, socialGroupID, suffix string) string {
	return fmt.Sprintf("slot:%s:social_group:%s:%s", slotID, socialGroupID, suffix)
}

func calculateBannerScore(selects, clicks, totalSelects float64) float64 {
	bannerRatio := clicks / selects
	return bannerRatio + math.Sqrt((2.0*math.Log(totalSelects))/selects)
}
