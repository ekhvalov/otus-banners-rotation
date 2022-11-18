package internalgrpc

import (
	"context"
	"fmt"

	grpcapi "github.com/ekhvalov/otus-banners-rotation/pkg/api/grpc"
)

type handler struct {
	grpcapi.UnimplementedRotatorServer
}

func (r *handler) CreateBanner(ctx context.Context, request *grpcapi.CreateBannerRequest) (*grpcapi.CreateBannerResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (r *handler) DeleteBanner(ctx context.Context, request *grpcapi.DeleteBannerRequest) (*grpcapi.DeleteBannerResponse, error) {
	//TODO implement me
	return nil, fmt.Errorf("unimplemented")
}

func (r *handler) CreateSlot(ctx context.Context, request *grpcapi.CreateSlotRequest) (*grpcapi.CreateSlotResponse, error) {
	//TODO implement me
	return nil, fmt.Errorf("unimplemented")
}

func (r *handler) DeleteSlot(ctx context.Context, request *grpcapi.DeleteSlotRequest) (*grpcapi.DeleteSlotResponse, error) {
	//TODO implement me
	return nil, fmt.Errorf("unimplemented")
}

func (r *handler) CreateSocialGroup(ctx context.Context, request *grpcapi.CreateSocialGroupRequest) (*grpcapi.CreateSocialGroupResponse, error) {
	//TODO implement me
	return nil, fmt.Errorf("unimplemented")
}

func (r *handler) DeleteSocialGroup(ctx context.Context, request *grpcapi.DeleteSocialGroupRequest) (*grpcapi.DeleteSocialGroupResponse, error) {
	//TODO implement me
	return nil, fmt.Errorf("unimplemented")
}

func (r *handler) AttachBanner(ctx context.Context, request *grpcapi.AttachBannerRequest) (*grpcapi.AttachBannerResponse, error) {
	//TODO implement me
	return nil, fmt.Errorf("unimplemented")
}

func (r *handler) DetachBanner(ctx context.Context, request *grpcapi.DetachBannerRequest) (*grpcapi.DetachBannerResponse, error) {
	//TODO implement me
	return nil, fmt.Errorf("unimplemented")
}

func (r *handler) ClickBanner(ctx context.Context, request *grpcapi.ClickBannerRequest) (*grpcapi.ClickBannerResponse, error) {
	//TODO implement me
	return nil, fmt.Errorf("unimplemented")
}

func (r *handler) SelectBanner(ctx context.Context, request *grpcapi.SelectBannerRequest) (*grpcapi.SelectBannerResponse, error) {
	//TODO implement me
	return nil, fmt.Errorf("unimplemented")
}
