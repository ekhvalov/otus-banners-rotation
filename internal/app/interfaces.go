//go:generate mockgen -destination=./mock/storage.gen.go -package mock . Storage
//go:generate mockgen -destination=./mock/event_queue.gen.go -package mock . EventQueue

package app

import "context"

type Storage interface {
	CreateBanner(ctx context.Context, description string) (id string, err error)
	DeleteBanner(ctx context.Context, id string) error
	CreateSlot(ctx context.Context, description string) (id string, err error)
	DeleteSlot(ctx context.Context, id string) error
	CreateSocialGroup(ctx context.Context, description string) (id string, err error)
	DeleteSocialGroup(ctx context.Context, id string) error
	AttachBanner(ctx context.Context, slotID, bannerID string) error
	DetachBanner(ctx context.Context, slotID, bannerID string) error
	SelectBanner(ctx context.Context, slotID, socialGroupID string) (bannerID string, err error)
	ClickBanner(ctx context.Context, slotID, bannerID, socialGroupID string) error
}

type EventQueue interface {
	Put(ctx context.Context, event Event) error
}
