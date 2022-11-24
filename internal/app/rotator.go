package app

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrEmptyDescription = errors.New("description is empty")
	ErrEmptyID          = errors.New("id is empty")
)

func NewRotator(storage Storage, eventQueue EventQueue) Rotator {
	return Rotator{storage: storage, eventQueue: eventQueue}
}

type Rotator struct {
	storage    Storage
	eventQueue EventQueue
}

func (r Rotator) CreateBanner(ctx context.Context, description string) (string, error) {
	if description == "" {
		return "", ErrEmptyDescription
	}
	id, err := r.storage.CreateBanner(ctx, description)
	if err != nil {
		return "", fmt.Errorf("create banner error: %w", err)
	}
	return id, nil
}

func (r Rotator) DeleteBanner(ctx context.Context, id string) error {
	if id == "" {
		return ErrEmptyID
	}
	err := r.storage.DeleteBanner(ctx, id)
	if err != nil {
		return fmt.Errorf("delete banner error: %w", err)
	}
	return nil
}

func (r Rotator) CreateSlot(ctx context.Context, description string) (string, error) {
	if description == "" {
		return "", ErrEmptyDescription
	}
	id, err := r.storage.CreateSlot(ctx, description)
	if err != nil {
		return "", fmt.Errorf("create slot error: %w", err)
	}
	return id, nil
}

func (r Rotator) DeleteSlot(ctx context.Context, id string) error {
	if id == "" {
		return ErrEmptyID
	}
	if err := r.storage.DeleteSlot(ctx, id); err != nil {
		return fmt.Errorf("delete slot error: %w", err)
	}
	return nil
}

func (r Rotator) CreateSocialGroup(ctx context.Context, description string) (string, error) {
	if description == "" {
		return "", ErrEmptyDescription
	}
	id, err := r.storage.CreateSocialGroup(ctx, description)
	if err != nil {
		return "", fmt.Errorf("create social group error: %w", err)
	}
	return id, nil
}

func (r Rotator) DeleteSocialGroup(ctx context.Context, id string) error {
	if id == "" {
		return ErrEmptyID
	}
	if err := r.storage.DeleteSocialGroup(ctx, id); err != nil {
		return fmt.Errorf("delete social group error: %w", err)
	}
	return nil
}

func (r Rotator) AttachBanner(ctx context.Context, slotID, bannerID string) error {
	if slotID == "" {
		return fmt.Errorf("slot id error: %w", ErrEmptyID)
	}
	if bannerID == "" {
		return fmt.Errorf("banner id error: %w", ErrEmptyID)
	}
	if err := r.storage.AttachBanner(ctx, slotID, bannerID); err != nil {
		return fmt.Errorf("attach banner error: %w", err)
	}
	return nil
}

func (r Rotator) DetachBanner(ctx context.Context, slotID, bannerID string) error {
	if slotID == "" {
		return fmt.Errorf("slot id error: %w", ErrEmptyID)
	}
	if bannerID == "" {
		return fmt.Errorf("banner id error: %w", ErrEmptyID)
	}
	if err := r.storage.DetachBanner(ctx, slotID, bannerID); err != nil {
		return fmt.Errorf("detach banner error: %w", err)
	}
	return nil
}

func (r Rotator) SelectBanner(ctx context.Context, slotID, socialGroupID string) (string, error) {
	if slotID == "" {
		return "", fmt.Errorf("slot id error: %w", ErrEmptyID)
	}
	if socialGroupID == "" {
		return "", fmt.Errorf("social group id error: %w", ErrEmptyID)
	}
	bannerID, err := r.storage.SelectBanner(ctx, slotID, socialGroupID)
	if err != nil {
		return "", fmt.Errorf("select banner error: %w", err)
	}
	// TODO: Log queue error?
	err = r.eventQueue.Put(ctx, Event{
		Type:           EventSelect,
		SlotID:         slotID,
		BannerID:       bannerID,
		SocialGroupID:  socialGroupID,
		TimestampMicro: time.Now().UnixMicro(),
	})
	return bannerID, nil
}

func (r Rotator) ClickBanner(ctx context.Context, slotID, bannerID, socialGroupID string) error {
	if slotID == "" {
		return fmt.Errorf("slot id error: %w", ErrEmptyID)
	}
	if bannerID == "" {
		return fmt.Errorf("banner id error: %w", ErrEmptyID)
	}
	if socialGroupID == "" {
		return fmt.Errorf("social group id error: %w", ErrEmptyID)
	}
	if err := r.storage.ClickBanner(ctx, slotID, bannerID, socialGroupID); err != nil {
		return fmt.Errorf("click banner error: %w", err)
	}
	// TODO: Log queue error?
	_ = r.eventQueue.Put(ctx, Event{
		Type:           EventClick,
		SlotID:         slotID,
		BannerID:       bannerID,
		SocialGroupID:  socialGroupID,
		TimestampMicro: time.Now().UnixMicro(),
	})
	return nil
}
