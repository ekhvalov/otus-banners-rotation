package internalgrpc

import (
	"context"
	"errors"

	"github.com/ekhvalov/otus-banners-rotation/internal/app"
	grpcapi "github.com/ekhvalov/otus-banners-rotation/pkg/api/grpc"
	"google.golang.org/genproto/googleapis/rpc/code"
)

var statusOK = grpcapi.Status{Code: code.Code_OK}

type handler struct {
	grpcapi.UnimplementedRotatorServer
	rotator app.Rotator
}

func (h *handler) CreateBanner(
	ctx context.Context,
	request *grpcapi.CreateBannerRequest,
) (*grpcapi.CreateBannerResponse, error) {
	id, err := h.rotator.CreateBanner(ctx, request.GetDescription())
	if err != nil {
		if errors.Is(err, app.ErrEmptyDescription) {
			return &grpcapi.CreateBannerResponse{Status: makeStatus(code.Code_INVALID_ARGUMENT, err)}, nil
		}
		return nil, err
	}
	return &grpcapi.CreateBannerResponse{Status: &statusOK, Id: id}, nil
}

func (h *handler) DeleteBanner(
	ctx context.Context,
	request *grpcapi.DeleteBannerRequest,
) (*grpcapi.DeleteBannerResponse, error) {
	err := h.rotator.DeleteBanner(ctx, request.GetId())
	if err != nil {
		if errors.Is(err, app.ErrEmptyID) {
			return &grpcapi.DeleteBannerResponse{Status: makeStatus(code.Code_INVALID_ARGUMENT, err)}, nil
		}
		return nil, err
	}
	return &grpcapi.DeleteBannerResponse{Status: &statusOK}, nil
}

func (h *handler) CreateSlot(
	ctx context.Context,
	request *grpcapi.CreateSlotRequest,
) (*grpcapi.CreateSlotResponse, error) {
	id, err := h.rotator.CreateSlot(ctx, request.GetDescription())
	if err != nil {
		if errors.Is(err, app.ErrEmptyDescription) {
			return &grpcapi.CreateSlotResponse{Status: makeStatus(code.Code_INVALID_ARGUMENT, err)}, nil
		}
		return nil, err
	}
	return &grpcapi.CreateSlotResponse{Status: &statusOK, Id: id}, nil
}

func (h *handler) DeleteSlot(
	ctx context.Context,
	request *grpcapi.DeleteSlotRequest,
) (*grpcapi.DeleteSlotResponse, error) {
	err := h.rotator.DeleteSlot(ctx, request.GetId())
	if err != nil {
		if errors.Is(err, app.ErrEmptyID) {
			return &grpcapi.DeleteSlotResponse{Status: makeStatus(code.Code_INVALID_ARGUMENT, err)}, nil
		}
		return nil, err
	}
	return &grpcapi.DeleteSlotResponse{Status: &statusOK}, nil
}

func (h *handler) CreateSocialGroup(
	ctx context.Context,
	request *grpcapi.CreateSocialGroupRequest,
) (*grpcapi.CreateSocialGroupResponse, error) {
	id, err := h.rotator.CreateSocialGroup(ctx, request.GetDescription())
	if err != nil {
		if errors.Is(err, app.ErrEmptyDescription) {
			return &grpcapi.CreateSocialGroupResponse{Status: makeStatus(code.Code_INVALID_ARGUMENT, err)}, nil
		}
		return nil, err
	}
	return &grpcapi.CreateSocialGroupResponse{Status: &statusOK, Id: id}, nil
}

func (h *handler) DeleteSocialGroup(
	ctx context.Context,
	request *grpcapi.DeleteSocialGroupRequest,
) (*grpcapi.DeleteSocialGroupResponse, error) {
	err := h.rotator.DeleteSocialGroup(ctx, request.GetId())
	if err != nil {
		if errors.Is(err, app.ErrEmptyID) {
			return &grpcapi.DeleteSocialGroupResponse{Status: makeStatus(code.Code_INVALID_ARGUMENT, err)}, nil
		}
		return nil, err
	}
	return &grpcapi.DeleteSocialGroupResponse{Status: &statusOK}, nil
}

func (h *handler) AttachBanner(
	ctx context.Context,
	request *grpcapi.AttachBannerRequest,
) (*grpcapi.AttachBannerResponse, error) {
	err := h.rotator.AttachBanner(ctx, request.GetSlotId(), request.GetBannerId())
	if err != nil {
		var errNotFound *app.ErrNotFound
		if errors.As(err, &errNotFound) {
			return &grpcapi.AttachBannerResponse{Status: makeStatus(code.Code_NOT_FOUND, err)}, nil
		}
		return nil, err
	}
	return &grpcapi.AttachBannerResponse{Status: &statusOK}, nil
}

func (h *handler) DetachBanner(
	ctx context.Context,
	request *grpcapi.DetachBannerRequest,
) (*grpcapi.DetachBannerResponse, error) {
	err := h.rotator.DetachBanner(ctx, request.GetSlotId(), request.GetBannerId())
	if err != nil {
		var errNotFound *app.ErrNotFound
		if errors.As(err, &errNotFound) {
			return &grpcapi.DetachBannerResponse{Status: makeStatus(code.Code_NOT_FOUND, err)}, nil
		}
		return nil, err
	}
	return &grpcapi.DetachBannerResponse{Status: &statusOK}, nil
}

func (h *handler) ClickBanner(
	ctx context.Context,
	request *grpcapi.ClickBannerRequest,
) (*grpcapi.ClickBannerResponse, error) {
	err := h.rotator.ClickBanner(ctx, request.GetSlotId(), request.GetBannerId(), request.GetSocialGroupId())
	if err != nil {
		var errNotFound *app.ErrNotFound
		if errors.As(err, &errNotFound) {
			return &grpcapi.ClickBannerResponse{Status: makeStatus(code.Code_NOT_FOUND, err)}, nil
		}
		var errNotAttached *app.ErrBannerNotAttached
		if errors.As(err, &errNotAttached) {
			return &grpcapi.ClickBannerResponse{Status: makeStatus(code.Code_FAILED_PRECONDITION, err)}, nil
		}
		return nil, err
	}
	return &grpcapi.ClickBannerResponse{Status: &statusOK}, nil
}

func (h *handler) SelectBanner(
	ctx context.Context,
	request *grpcapi.SelectBannerRequest,
) (*grpcapi.SelectBannerResponse, error) {
	bannerID, err := h.rotator.SelectBanner(ctx, request.GetSlotId(), request.GetSocialGroupId())
	if err != nil {
		var errNotFound *app.ErrNotFound
		if errors.As(err, &errNotFound) {
			return &grpcapi.SelectBannerResponse{Status: makeStatus(code.Code_NOT_FOUND, err)}, nil
		}
		var errNotAttached *app.ErrBannerNotAttached
		if errors.As(err, &errNotAttached) {
			return &grpcapi.SelectBannerResponse{Status: makeStatus(code.Code_FAILED_PRECONDITION, err)}, nil
		}
		return nil, err
	}
	return &grpcapi.SelectBannerResponse{
		Status:   &statusOK,
		BannerId: bannerID,
	}, nil
}

func makeStatus(c code.Code, err error) *grpcapi.Status {
	return &grpcapi.Status{
		Code:    c,
		Message: err.Error(),
	}
}
