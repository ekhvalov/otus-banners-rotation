package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ekhvalov/otus-banners-rotation/internal/app"
	"github.com/ekhvalov/otus-banners-rotation/internal/app/mock"
	grpcapi "github.com/ekhvalov/otus-banners-rotation/pkg/api/grpc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/code"
)

var (
	id               = "100500"
	emptyID          = ""
	emptyDescription = ""
	description      = "Some description"
	bannerID         = "100500"
	slotID           = "100600"
	socialGroupID    = "100700"
	errRotator       = errors.New("rotator error")
	errNotFound      = app.NewErrNotFound("something is not found")
	errNotAttached   = app.NewErrBannerNotAttached(slotID, bannerID)
)

func Test_handler_CreateX(t *testing.T) {
	type response interface {
		GetStatus() *grpcapi.Status
		GetId() string
	}

	methods := map[string]struct {
		mockRotator func(c *gomock.Controller, expectedDesc string, returnID string, returnError error) app.Rotator
		callCreateX func(h handler, description string) (response, error)
	}{
		"CreateBanner": {
			mockRotator: func(c *gomock.Controller, expectedDesc string, returnID string, returnErr error) app.Rotator {
				rotator := mock.NewMockRotator(c)
				rotator.EXPECT().
					CreateBanner(context.Background(), expectedDesc).
					Return(returnID, returnErr)
				return rotator
			},
			callCreateX: func(h handler, description string) (response, error) {
				return h.CreateBanner(context.Background(), &grpcapi.CreateBannerRequest{Description: description})
			},
		},
		"CreateSlot": {
			mockRotator: func(c *gomock.Controller, expectedDescription string, returnID string, returnErr error) app.Rotator {
				rotator := mock.NewMockRotator(c)
				rotator.EXPECT().
					CreateSlot(context.Background(), expectedDescription).
					Return(returnID, returnErr)
				return rotator
			},
			callCreateX: func(h handler, description string) (response, error) {
				return h.CreateSlot(context.Background(), &grpcapi.CreateSlotRequest{Description: description})
			},
		},
		"CreateSocialGroup": {
			mockRotator: func(c *gomock.Controller, expectedDescription string, returnID string, returnErr error) app.Rotator {
				rotator := mock.NewMockRotator(c)
				rotator.EXPECT().
					CreateSocialGroup(context.Background(), expectedDescription).
					Return(returnID, returnErr)
				return rotator
			},
			callCreateX: func(h handler, description string) (response, error) {
				return h.CreateSocialGroup(
					context.Background(),
					&grpcapi.CreateSocialGroupRequest{Description: description},
				)
			},
		},
	}

	tests := map[string]struct {
		argDescription       string
		rotatorReturnID      string
		rotatorReturnError   error
		expectedResponseCode code.Code
		expectedResponseID   string
		expectedError        error
	}{
		"empty description": {
			argDescription:       emptyDescription,
			rotatorReturnError:   app.ErrEmptyDescription,
			expectedResponseCode: code.Code_INVALID_ARGUMENT,
		},
		"rotator error": {
			argDescription:     description,
			rotatorReturnError: errRotator,
			expectedError:      errRotator,
		},
		"no error": {
			argDescription:       description,
			rotatorReturnID:      id,
			expectedResponseCode: code.Code_OK,
			expectedResponseID:   id,
		},
	}

	for methodName, m := range methods {
		for testName, tt := range tests {
			t.Run(fmt.Sprintf("%s_%s", methodName, testName), func(t *testing.T) {
				controller := gomock.NewController(t)
				defer controller.Finish()
				r := m.mockRotator(controller, tt.argDescription, tt.rotatorReturnID, tt.rotatorReturnError)
				h := handler{rotator: r}

				resp, err := m.callCreateX(h, tt.argDescription)

				if tt.expectedError == nil {
					require.NoError(t, err)
					require.Equal(t, tt.expectedResponseCode, resp.GetStatus().GetCode())
					require.Equal(t, tt.expectedResponseID, resp.GetId())
				} else {
					require.ErrorIs(t, err, tt.expectedError)
					require.Nil(t, resp)
				}
			})
		}
	}
}

func Test_handler_DeleteX(t *testing.T) {
	type response interface {
		GetStatus() *grpcapi.Status
	}

	methods := map[string]struct {
		mockRotator      func(controller *gomock.Controller, expectedID string, returnError error) app.Rotator
		callDeleteMethod func(h handler, deleteID string) (response, error)
	}{
		"DeleteBanner": {
			mockRotator: func(controller *gomock.Controller, expectedID string, returnError error) app.Rotator {
				rotator := mock.NewMockRotator(controller)
				rotator.EXPECT().
					DeleteBanner(context.Background(), expectedID).
					Return(returnError)
				return rotator
			},
			callDeleteMethod: func(h handler, deleteID string) (response, error) {
				return h.DeleteBanner(context.Background(), &grpcapi.DeleteBannerRequest{Id: deleteID})
			},
		},
		"DeleteSlot": {
			mockRotator: func(controller *gomock.Controller, expectedID string, returnError error) app.Rotator {
				rotator := mock.NewMockRotator(controller)
				rotator.EXPECT().
					DeleteSlot(context.Background(), expectedID).
					Return(returnError)
				return rotator
			},
			callDeleteMethod: func(h handler, deleteID string) (response, error) {
				return h.DeleteSlot(context.Background(), &grpcapi.DeleteSlotRequest{Id: deleteID})
			},
		},
		"DeleteSocialGroup": {
			mockRotator: func(controller *gomock.Controller, expectedID string, returnError error) app.Rotator {
				rotator := mock.NewMockRotator(controller)
				rotator.EXPECT().
					DeleteSocialGroup(context.Background(), expectedID).
					Return(returnError)
				return rotator
			},
			callDeleteMethod: func(h handler, deleteID string) (response, error) {
				return h.DeleteSocialGroup(context.Background(), &grpcapi.DeleteSocialGroupRequest{Id: deleteID})
			},
		},
	}

	tests := map[string]struct {
		argID                string
		rotatorReturnError   error
		expectedResponseCode code.Code
		expectedError        error
	}{
		"empty id": {
			argID:                emptyID,
			rotatorReturnError:   app.ErrEmptyID,
			expectedResponseCode: code.Code_INVALID_ARGUMENT,
		},
		"rotator error": {
			argID:              id,
			rotatorReturnError: errRotator,
			expectedError:      errRotator,
		},
		"no error": {
			argID:                id,
			expectedResponseCode: code.Code_OK,
		},
	}

	for methodName, m := range methods {
		for testName, tt := range tests {
			t.Run(fmt.Sprintf("%s_%s", methodName, testName), func(t *testing.T) {
				controller := gomock.NewController(t)
				defer controller.Finish()
				h := handler{rotator: m.mockRotator(controller, tt.argID, tt.rotatorReturnError)}

				resp, err := m.callDeleteMethod(h, tt.argID)

				if tt.expectedError == nil {
					require.NoError(t, err)
					require.Equal(t, tt.expectedResponseCode, resp.GetStatus().GetCode())
				} else {
					require.ErrorIs(t, err, tt.expectedError)
					require.Nil(t, resp)
				}
			})
		}
	}
}

//nolint:dupl // Ignore duplication with Test_handler_DetachBanner
func Test_handler_AttachBanner(t *testing.T) {
	tests := map[string]struct {
		slotID           string
		bannerID         string
		rotatorReturnErr error
		wantErr          error
		wantResponseCode code.Code
	}{
		"not found error": {
			slotID:           slotID,
			bannerID:         bannerID,
			rotatorReturnErr: errNotFound,
			wantErr:          nil,
			wantResponseCode: code.Code_NOT_FOUND,
		},
		"rotator error": {
			slotID:           slotID,
			bannerID:         bannerID,
			rotatorReturnErr: errRotator,
			wantErr:          errRotator,
		},
		"no error": {
			slotID:           slotID,
			bannerID:         bannerID,
			wantResponseCode: code.Code_OK,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			r := mock.NewMockRotator(controller)
			r.EXPECT().
				AttachBanner(context.Background(), tt.slotID, tt.bannerID).
				Return(tt.rotatorReturnErr)

			h := &handler{rotator: r}

			gotResponse, err := h.AttachBanner(
				context.Background(),
				&grpcapi.AttachBannerRequest{SlotId: tt.slotID, BannerId: bannerID},
			)

			if tt.wantErr == nil {
				require.NoError(t, err)
				require.Equal(t, tt.wantResponseCode, gotResponse.GetStatus().GetCode())
			} else {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, gotResponse)
			}
		})
	}
}

//nolint:dupl // Ignore duplication with Test_handler_AttachBanner
func Test_handler_DetachBanner(t *testing.T) {
	tests := map[string]struct {
		slotID           string
		bannerID         string
		rotatorReturnErr error
		wantErr          error
		wantResponseCode code.Code
	}{
		"not found error": {
			slotID:           slotID,
			bannerID:         bannerID,
			rotatorReturnErr: errNotFound,
			wantErr:          nil,
			wantResponseCode: code.Code_NOT_FOUND,
		},
		"rotator error": {
			slotID:           slotID,
			bannerID:         bannerID,
			rotatorReturnErr: errRotator,
			wantErr:          errRotator,
		},
		"no error": {
			slotID:           slotID,
			bannerID:         bannerID,
			wantResponseCode: code.Code_OK,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			r := mock.NewMockRotator(controller)
			r.EXPECT().
				DetachBanner(context.Background(), tt.slotID, tt.bannerID).
				Return(tt.rotatorReturnErr)

			h := &handler{rotator: r}

			gotResponse, err := h.DetachBanner(
				context.Background(),
				&grpcapi.DetachBannerRequest{SlotId: tt.slotID, BannerId: bannerID},
			)

			if tt.wantErr == nil {
				require.NoError(t, err)
				require.Equal(t, tt.wantResponseCode, gotResponse.GetStatus().GetCode())
			} else {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, gotResponse)
			}
		})
	}
}

func Test_handler_ClickBanner(t *testing.T) {
	tests := map[string]struct {
		slotID           string
		bannerID         string
		socialGroupID    string
		rotatorReturnErr error
		wantResponseCode code.Code
		wantErr          error
	}{
		"not found error": {
			slotID:           slotID,
			bannerID:         bannerID,
			socialGroupID:    socialGroupID,
			rotatorReturnErr: errNotFound,
			wantErr:          nil,
			wantResponseCode: code.Code_NOT_FOUND,
		},
		"not attached error": {
			slotID:           slotID,
			bannerID:         bannerID,
			socialGroupID:    socialGroupID,
			rotatorReturnErr: errNotAttached,
			wantErr:          nil,
			wantResponseCode: code.Code_FAILED_PRECONDITION,
		},
		"rotator error": {
			slotID:           slotID,
			bannerID:         bannerID,
			socialGroupID:    socialGroupID,
			rotatorReturnErr: errRotator,
			wantErr:          errRotator,
		},
		"no error": {
			slotID:           slotID,
			bannerID:         bannerID,
			socialGroupID:    socialGroupID,
			rotatorReturnErr: nil,
			wantErr:          nil,
			wantResponseCode: code.Code_OK,
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			r := mock.NewMockRotator(controller)
			r.EXPECT().
				ClickBanner(context.Background(), tt.slotID, tt.bannerID, tt.socialGroupID).
				Return(tt.rotatorReturnErr)
			h := &handler{rotator: r}

			gotResponse, err := h.ClickBanner(context.Background(), &grpcapi.ClickBannerRequest{
				SlotId:        tt.slotID,
				BannerId:      tt.bannerID,
				SocialGroupId: tt.socialGroupID,
			})

			if tt.wantErr == nil {
				require.NoError(t, err)
				require.Equal(t, tt.wantResponseCode, gotResponse.GetStatus().GetCode())
			} else {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, gotResponse)
			}
		})
	}
}

func Test_handler_SelectBanner(t *testing.T) {
	tests := map[string]struct {
		slotID           string
		bannerID         string
		socialGroupID    string
		rotatorReturnErr error
		wantBannerID     string
		wantResponseCode code.Code
		wantErr          error
	}{
		"not found error": {
			slotID:           slotID,
			bannerID:         bannerID,
			socialGroupID:    socialGroupID,
			rotatorReturnErr: errNotFound,
			wantErr:          nil,
			wantResponseCode: code.Code_NOT_FOUND,
		},
		"not attached error": {
			slotID:           slotID,
			bannerID:         bannerID,
			socialGroupID:    socialGroupID,
			rotatorReturnErr: errNotAttached,
			wantErr:          nil,
			wantResponseCode: code.Code_FAILED_PRECONDITION,
		},
		"rotator error": {
			slotID:           slotID,
			bannerID:         bannerID,
			socialGroupID:    socialGroupID,
			rotatorReturnErr: errRotator,
			wantErr:          errRotator,
		},
		"no error": {
			slotID:           slotID,
			bannerID:         bannerID,
			socialGroupID:    socialGroupID,
			wantBannerID:     bannerID,
			rotatorReturnErr: nil,
			wantErr:          nil,
			wantResponseCode: code.Code_OK,
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			r := mock.NewMockRotator(controller)
			r.EXPECT().
				SelectBanner(context.Background(), tt.slotID, tt.socialGroupID).
				Return(tt.bannerID, tt.rotatorReturnErr)
			h := &handler{rotator: r}

			gotResponse, err := h.SelectBanner(context.Background(), &grpcapi.SelectBannerRequest{
				SlotId:        tt.slotID,
				SocialGroupId: tt.socialGroupID,
			})

			if tt.wantErr == nil {
				require.NoError(t, err)
				require.Equal(t, tt.wantResponseCode, gotResponse.GetStatus().GetCode())
				require.Equal(t, tt.wantBannerID, gotResponse.GetBannerId())
			} else {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, gotResponse)
			}
		})
	}
}
