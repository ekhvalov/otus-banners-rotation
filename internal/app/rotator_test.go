package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ekhvalov/otus-banners-rotation/internal/app"
	"github.com/ekhvalov/otus-banners-rotation/internal/app/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	id               = "100500"
	emptyID          = ""
	emptyDescription = ""
	description      = "Some description"
	bannerID         = "100500"
	slotID           = "100600"
	socialGroupID    = "100700"
	errStorage       = errors.New("storage error")
)

var testsCreateX = map[string]struct {
	isMockExpected        bool
	mockExpectDescription string
	mockReturnID          string
	mockReturnErr         error
	description           string
	want                  string
	err                   error
}{
	"empty description": {
		description: emptyDescription,
		err:         app.ErrEmptyDescription,
	},
	"storage error": {
		isMockExpected:        true,
		mockExpectDescription: description,
		mockReturnErr:         errStorage,
		description:           description,
		err:                   errStorage,
	},
	"no error": {
		isMockExpected:        true,
		mockExpectDescription: description,
		mockReturnID:          id,
		mockReturnErr:         nil,
		description:           description,
		want:                  id,
		err:                   nil,
	},
}

func TestRotator_CreateBanner(t *testing.T) {
	for testName, tt := range testsCreateX {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			storage := mock.NewMockStorage(controller)
			if tt.isMockExpected {
				storage.EXPECT().
					CreateBanner(context.Background(), tt.mockExpectDescription).
					Return(tt.mockReturnID, tt.mockReturnErr)
			}
			rotator := app.NewRotator(storage)

			got, err := rotator.CreateBanner(context.Background(), tt.description)

			if tt.err == nil {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.ErrorIs(t, err, tt.err)
				require.Empty(t, got)
			}
		})
	}
}

func TestRotator_CreateSlot(t *testing.T) {
	for testName, tt := range testsCreateX {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			storage := mock.NewMockStorage(controller)
			if tt.isMockExpected {
				storage.EXPECT().
					CreateSlot(context.Background(), tt.mockExpectDescription).
					Return(tt.mockReturnID, tt.mockReturnErr)
			}
			rotator := app.NewRotator(storage)

			got, err := rotator.CreateSlot(context.Background(), tt.description)

			if tt.err == nil {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.ErrorIs(t, err, tt.err)
				require.Empty(t, got)
			}
		})
	}
}

func TestRotator_CreateSocialGroup(t *testing.T) {
	for testName, tt := range testsCreateX {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			storage := mock.NewMockStorage(controller)
			if tt.isMockExpected {
				storage.EXPECT().
					CreateSocialGroup(context.Background(), tt.mockExpectDescription).
					Return(tt.mockReturnID, tt.mockReturnErr)
			}
			rotator := app.NewRotator(storage)

			got, err := rotator.CreateSocialGroup(context.Background(), tt.description)

			if tt.err == nil {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.ErrorIs(t, err, tt.err)
				require.Empty(t, got)
			}
		})
	}
}

var testsDeleteX = map[string]struct {
	isMockExpected bool
	mockExpectID   string
	mockReturnErr  error
	id             string
	err            error
}{
	"empty id": {
		id:  emptyID,
		err: app.ErrEmptyID,
	},
	"storage error": {
		isMockExpected: true,
		mockExpectID:   id,
		mockReturnErr:  errStorage,
		id:             id,
		err:            errStorage,
	},
	"no error": {
		isMockExpected: true,
		mockExpectID:   id,
		mockReturnErr:  nil,
		id:             id,
		err:            nil,
	},
}

func TestRotator_DeleteBanner(t *testing.T) {
	for testName, tt := range testsDeleteX {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			storage := mock.NewMockStorage(controller)
			if tt.isMockExpected {
				storage.EXPECT().
					DeleteBanner(context.Background(), tt.mockExpectID).
					Return(tt.mockReturnErr)
			}
			rotator := app.NewRotator(storage)

			err := rotator.DeleteBanner(context.Background(), tt.id)

			if tt.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}

func TestRotator_DeleteSlot(t *testing.T) {
	for testName, tt := range testsDeleteX {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			storage := mock.NewMockStorage(controller)
			if tt.isMockExpected {
				storage.EXPECT().
					DeleteSlot(context.Background(), tt.mockExpectID).
					Return(tt.mockReturnErr)
			}
			rotator := app.NewRotator(storage)

			err := rotator.DeleteSlot(context.Background(), tt.id)

			if tt.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}

func TestRotator_DeleteSocialGroup(t *testing.T) {
	for testName, tt := range testsDeleteX {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			storage := mock.NewMockStorage(controller)
			if tt.isMockExpected {
				storage.EXPECT().
					DeleteSocialGroup(context.Background(), tt.mockExpectID).
					Return(tt.mockReturnErr)
			}
			rotator := app.NewRotator(storage)

			err := rotator.DeleteSocialGroup(context.Background(), tt.id)

			if tt.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}

var testsAttachDetachBanner = map[string]struct {
	isMockExpected     bool
	mockExpectSlotID   string
	mockExpectBannerID string
	mockReturnErr      error
	slotID             string
	bannerID           string
	err                error
}{
	"empty slot id": {
		slotID:   emptyID,
		bannerID: bannerID,
		err:      app.ErrEmptyID,
	},
	"empty banner id": {
		slotID:   slotID,
		bannerID: emptyID,
		err:      app.ErrEmptyID,
	},
	"storage error": {
		isMockExpected:     true,
		mockExpectSlotID:   slotID,
		mockExpectBannerID: bannerID,
		mockReturnErr:      errStorage,
		slotID:             slotID,
		bannerID:           bannerID,
		err:                errStorage,
	},
	"no error": {
		isMockExpected:     true,
		mockExpectSlotID:   slotID,
		mockExpectBannerID: bannerID,
		mockReturnErr:      nil,
		slotID:             slotID,
		bannerID:           bannerID,
		err:                nil,
	},
}

func TestRotator_AttachBanner(t *testing.T) {
	for testName, tt := range testsAttachDetachBanner {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			storage := mock.NewMockStorage(controller)
			if tt.isMockExpected {
				storage.EXPECT().
					AttachBanner(context.Background(), tt.mockExpectSlotID, tt.mockExpectBannerID).
					Return(tt.mockReturnErr)
			}
			rotator := app.NewRotator(storage)

			err := rotator.AttachBanner(context.Background(), tt.slotID, tt.bannerID)

			if tt.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}

func TestRotator_DetachBanner(t *testing.T) {
	for testName, tt := range testsAttachDetachBanner {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			storage := mock.NewMockStorage(controller)
			if tt.isMockExpected {
				storage.EXPECT().
					DetachBanner(context.Background(), tt.mockExpectSlotID, tt.mockExpectBannerID).
					Return(tt.mockReturnErr)
			}
			rotator := app.NewRotator(storage)

			err := rotator.DetachBanner(context.Background(), tt.slotID, tt.bannerID)

			if tt.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}

func TestRotator_SelectBanner(t *testing.T) {
	tests := map[string]struct {
		mockStorage   func(controller *gomock.Controller) app.Storage
		slotID        string
		socialGroupID string
		wantID        string
		err           error
	}{
		"empty slot id": {
			mockStorage: func(controller *gomock.Controller) app.Storage {
				return mock.NewMockStorage(controller)
			},
			slotID:        emptyID,
			socialGroupID: socialGroupID,
			err:           app.ErrEmptyID,
		},
		"empty social group id": {
			mockStorage: func(controller *gomock.Controller) app.Storage {
				return mock.NewMockStorage(controller)
			},
			slotID:        slotID,
			socialGroupID: emptyID,
			err:           app.ErrEmptyID,
		},
		"storage error": {
			mockStorage: func(controller *gomock.Controller) app.Storage {
				storage := mock.NewMockStorage(controller)
				storage.EXPECT().
					SelectBanner(context.Background(), slotID, socialGroupID).
					Return(emptyID, errStorage)
				return storage
			},
			slotID:        slotID,
			socialGroupID: socialGroupID,
			err:           errStorage,
		},
		"no error": {
			mockStorage: func(controller *gomock.Controller) app.Storage {
				storage := mock.NewMockStorage(controller)
				storage.EXPECT().
					SelectBanner(context.Background(), slotID, socialGroupID).
					Return(bannerID, nil)
				return storage
			},
			slotID:        slotID,
			socialGroupID: socialGroupID,
			wantID:        bannerID,
			err:           nil,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			rotator := app.NewRotator(tt.mockStorage(controller))

			gotID, err := rotator.SelectBanner(context.Background(), tt.slotID, tt.socialGroupID)

			if tt.err == nil {
				require.NoError(t, err)
				require.Equal(t, tt.wantID, gotID)
			} else {
				require.ErrorIs(t, err, tt.err)
				require.Empty(t, gotID)
			}
		})
	}
}

func TestRotator_ClickBanner(t *testing.T) {
	tests := map[string]struct {
		mockStorage   func(controller *gomock.Controller) app.Storage
		slotID        string
		bannerID      string
		socialGroupID string
		err           error
	}{
		"empty slot id": {
			mockStorage: func(controller *gomock.Controller) app.Storage {
				return mock.NewMockStorage(controller)
			},
			slotID:        emptyID,
			bannerID:      bannerID,
			socialGroupID: socialGroupID,
			err:           app.ErrEmptyID,
		},
		"empty banner id": {
			mockStorage: func(controller *gomock.Controller) app.Storage {
				return mock.NewMockStorage(controller)
			},
			slotID:        slotID,
			bannerID:      emptyID,
			socialGroupID: socialGroupID,
			err:           app.ErrEmptyID,
		},
		"empty social group id": {
			mockStorage: func(controller *gomock.Controller) app.Storage {
				return mock.NewMockStorage(controller)
			},
			slotID:        slotID,
			bannerID:      bannerID,
			socialGroupID: emptyID,
			err:           app.ErrEmptyID,
		},
		"storage error": {
			mockStorage: func(controller *gomock.Controller) app.Storage {
				storage := mock.NewMockStorage(controller)
				storage.EXPECT().
					ClickBanner(context.Background(), slotID, bannerID, socialGroupID).
					Return(errStorage)
				return storage
			},
			slotID:        slotID,
			bannerID:      bannerID,
			socialGroupID: socialGroupID,
			err:           errStorage,
		},
		"no error": {
			mockStorage: func(controller *gomock.Controller) app.Storage {
				storage := mock.NewMockStorage(controller)
				storage.EXPECT().
					ClickBanner(context.Background(), slotID, bannerID, socialGroupID).
					Return(nil)
				return storage
			},
			slotID:        slotID,
			bannerID:      bannerID,
			socialGroupID: socialGroupID,
			err:           nil,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			rotator := app.NewRotator(tt.mockStorage(controller))

			err := rotator.ClickBanner(context.Background(), tt.slotID, tt.bannerID, tt.socialGroupID)

			if tt.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}
