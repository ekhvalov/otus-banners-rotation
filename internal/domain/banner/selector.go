package banner

import "errors"

var (
	ErrBannerNotFound     = errors.New("banner is not found")
	ErrBannerAlreadyExist = errors.New("banner is already exist")
	ErrEmptyBannersList   = errors.New("banners list is empty")
)

// Selector defines methods of banner selector algorithm.
type Selector interface {
	// AddBanner adds banner into rotation
	// ErrBannerAlreadyExist could be returned
	AddBanner(bannerID string) error

	// DeleteBanner removes from rotation
	// ErrBannerNotFound could be returned
	DeleteBanner(bannerID string) error

	// SelectBanner returns banner for impression
	// ErrEmptyBannersList could be returned
	SelectBanner() (string, error)

	// RegisterClickForBanner increments clicks counter for banner
	// ErrBannerNotFound could be returned
	RegisterClickForBanner(bannerID string) error
}
