//go:generate mockgen -destination=./mock/storage.gen.go -package mock . Storage
//go:generate mockgen -destination=./mock/event_queue.gen.go -package mock . EventQueue
//go:generate mockgen -destination=./mock/rotator.gen.go -package mock . Rotator

package app

import (
	"context"
	"fmt"
)

func NewErrNotFound(message string) ErrNotFound {
	return ErrNotFound{message: message}
}

type ErrNotFound struct {
	message string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("not found: %s", e.message)
}

func NewErrBannerNotAttached(slotID, bannerID string) ErrBannerNotAttached {
	return ErrBannerNotAttached{slotID: slotID, bannerID: bannerID}
}

type ErrBannerNotAttached struct {
	slotID   string
	bannerID string
}

func (e *ErrBannerNotAttached) Error() string {
	return fmt.Sprintf("banner id '%s' is not attached to slot id '%s'", e.bannerID, e.slotID)
}

type Storage interface {
	// CreateBanner creates new banner.
	// Returns id of the created banner or an error
	CreateBanner(ctx context.Context, description string) (id string, err error)
	// DeleteBanner deletes banner with specified ID.
	// Returns ErrNotFound in case of the banner with specified id is not found.
	DeleteBanner(ctx context.Context, id string) error
	// CreateSlot creates new slot.
	// Returns id of the created slot or an error
	CreateSlot(ctx context.Context, description string) (id string, err error)
	// DeleteSlot deletes slot with specified ID.
	// Returns ErrNotFound in case of the slot with specified id is not found.
	DeleteSlot(ctx context.Context, id string) error
	// CreateSocialGroup creates new socialGroup.
	// Returns id of the created socialGroup or an error
	CreateSocialGroup(ctx context.Context, description string) (id string, err error)
	// DeleteSocialGroup deletes social group with specified ID.
	// Returns ErrNotFound in case of the social group with specified id is not found.
	DeleteSocialGroup(ctx context.Context, id string) error
	// AttachBanner attaches a banner to a slot.
	// Returns ErrNotFound in case of a banner or a slot is not found.
	AttachBanner(ctx context.Context, slotID, bannerID string) error
	// DetachBanner attaches a banner to a slot.
	// Returns ErrNotFound in case of a banner or a slot is not found.
	DetachBanner(ctx context.Context, slotID, bannerID string) error
	// SelectBanner selects a banner from a slot for social group.
	// Returns ErrNotFound in case of a banner or slot or social group is not found.
	// Returns ErrBannerNotAttached in case of a banner is not attached to a slot.
	SelectBanner(ctx context.Context, slotID, socialGroupID string) (bannerID string, err error)
	// ClickBanner registers a click on a banner in a slot by social group.
	// Returns ErrNotFound in case of a banner or slot or social group is not found.
	// Returns ErrBannerNotAttached in case of a banner is not attached to a slot.
	ClickBanner(ctx context.Context, slotID, bannerID, socialGroupID string) error
}

type EventQueue interface {
	Put(ctx context.Context, event Event) error
}

type Rotator interface {
	Storage
}
