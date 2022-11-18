//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	grpcapi "github.com/ekhvalov/otus-banners-rotation/pkg/api/grpc"
	"github.com/stretchr/testify/suite"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultGrpcServerHost = "localhost"
	defaultGrpcServerPort = "8081"
)

func TestRotator(t *testing.T) {
	suite.Run(t, new(rotatorSuite))
}

type rotatorSuite struct {
	suite.Suite
	ctx          context.Context
	cancel       context.CancelFunc
	tick         time.Duration
	waitFor      time.Duration
	clientGrpc   grpcapi.RotatorClient
	banners      map[string]struct{}
	slots        map[string]struct{}
	socialGroups map[string]struct{}
}

func (s *rotatorSuite) SetupSuite() {
	s.tick = time.Millisecond * 100
	s.waitFor = s.tick * 1000 * 30
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.banners = make(map[string]struct{})
	s.slots = make(map[string]struct{})
	s.socialGroups = make(map[string]struct{})

	var grpcConn *grpc.ClientConn
	var err error
	s.Require().Eventually(func() bool {
		grpcConn, err = grpc.DialContext(
			s.ctx, getGrpcServerAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		return err == nil
	}, s.waitFor, s.tick, fmt.Sprintf("grpc connection error: %v", err))
	s.clientGrpc = grpcapi.NewRotatorClient(grpcConn)
}

func (s *rotatorSuite) TearDownSuite() {
	s.cancel()
}

func (s *rotatorSuite) TearDownTest() {
	for id := range s.banners {
		_, err := s.clientGrpc.DeleteBanner(s.ctx, &grpcapi.DeleteBannerRequest{Id: id})
		s.Require().NoError(err)
	}
	s.banners = make(map[string]struct{})

	for id := range s.slots {
		_, err := s.clientGrpc.DeleteSlot(s.ctx, &grpcapi.DeleteSlotRequest{Id: id})
		s.Require().NoError(err)
	}
	s.slots = make(map[string]struct{})

	for id := range s.socialGroups {
		_, err := s.clientGrpc.DeleteSocialGroup(s.ctx, &grpcapi.DeleteSocialGroupRequest{Id: id})
		s.Require().NoError(err)
	}
	s.socialGroups = make(map[string]struct{})
}

func (s *rotatorSuite) Test_CreateBanner() {
	s.createBanner()
}

func (s *rotatorSuite) Test_DeleteBanner() {
	id := s.createBanner()
	s.deleteBanner(id)
}

func (s *rotatorSuite) Test_CreateSlot() {
	s.createSlot()
}
func (s *rotatorSuite) Test_DeleteSlot() {
	id := s.createSlot()
	s.deleteSlot(id)
}

func (s *rotatorSuite) Test_CreateSocialGroup() {
	s.createSocialGroup()
}

func (s *rotatorSuite) Test_DeleteSocialGroup() {
	id := s.createSocialGroup()
	s.deleteSocialGroup(id)
}

func (s *rotatorSuite) Test_AttachBanner() {
	slotID := s.createSlot()
	bannerID := s.createBanner()
	s.attachBanner(slotID, bannerID)
}

func (s *rotatorSuite) Test_DetachBanner() {
	slotID := s.createSlot()
	bannerID := s.createBanner()
	s.attachBanner(slotID, bannerID)

	resp, err := s.clientGrpc.DetachBanner(s.ctx, &grpcapi.DetachBannerRequest{
		SlotId:   slotID,
		BannerId: bannerID,
	})

	s.Require().NoError(err)
	s.Require().Equal(code.Code_OK, resp.GetStatus().GetCode())
}

func (s *rotatorSuite) Test_ClickBanner() {
	slotID := s.createSlot()
	bannerID := s.createBanner()
	groupID := s.createSocialGroup()
	s.attachBanner(slotID, bannerID)

	resp, err := s.clientGrpc.ClickBanner(s.ctx, &grpcapi.ClickBannerRequest{
		SlotId:        slotID,
		BannerId:      bannerID,
		SocialGroupId: groupID,
	})

	s.Require().NoError(err)
	s.Require().Equal(code.Code_OK, resp.GetStatus().GetCode())
}

func (s *rotatorSuite) Test_SelectBanner() {
	slotID := s.createSlot()
	bannerID := s.createBanner()
	groupID := s.createSocialGroup()
	s.attachBanner(slotID, bannerID)

	resp, err := s.clientGrpc.SelectBanner(s.ctx, &grpcapi.SelectBannerRequest{
		SlotId:        slotID,
		SocialGroupId: groupID,
	})

	s.Require().NoError(err)
	s.Require().Equal(code.Code_OK, resp.GetStatus().GetCode())
	s.Require().Equal(bannerID, resp.GetBannerId())
}

func (s *rotatorSuite) createBanner() string {
	description := generateDescription("Banner")
	resp, err := s.clientGrpc.CreateBanner(s.ctx, &grpcapi.CreateBannerRequest{Description: description})

	s.Require().NoError(err, "create banner error")
	s.Require().NotNil(resp)
	s.Require().NotEmpty(resp.GetId())
	s.Require().Equal(code.Code_OK, resp.GetStatus().GetCode())

	s.banners[resp.GetId()] = struct{}{}
	return resp.GetId()
}

func (s *rotatorSuite) deleteBanner(id string) {
	resp, err := s.clientGrpc.DeleteBanner(s.ctx, &grpcapi.DeleteBannerRequest{Id: id})

	s.Require().NoError(err, "delete banner error")
	s.Require().NotNil(resp)
	s.Require().Equal(code.Code_OK, resp.GetStatus().GetCode())

	delete(s.banners, id)
}

func (s *rotatorSuite) createSlot() string {
	description := generateDescription("Slot")
	resp, err := s.clientGrpc.CreateSlot(s.ctx, &grpcapi.CreateSlotRequest{Description: description})

	s.Require().NoError(err, "create slot error")
	s.Require().NotNil(resp)
	s.Require().NotEmpty(resp.GetId())
	s.Require().Equal(code.Code_OK, resp.GetStatus().GetCode())

	s.slots[resp.GetId()] = struct{}{}
	return resp.GetId()
}

func (s *rotatorSuite) deleteSlot(id string) {
	resp, err := s.clientGrpc.DeleteSlot(s.ctx, &grpcapi.DeleteSlotRequest{Id: id})

	s.Require().NoError(err, "delete slot error")
	s.Require().NotNil(resp)
	s.Require().Equal(code.Code_OK, resp.GetStatus().GetCode())

	delete(s.slots, id)
}

func (s *rotatorSuite) createSocialGroup() string {
	description := generateDescription("Social group")
	resp, err := s.clientGrpc.CreateSocialGroup(s.ctx, &grpcapi.CreateSocialGroupRequest{Description: description})

	s.Require().NoError(err, "create socialGroup error")
	s.Require().NotNil(resp)
	s.Require().NotEmpty(resp.GetId())
	s.Require().Equal(code.Code_OK, resp.GetStatus().GetCode())

	s.socialGroups[resp.GetId()] = struct{}{}
	return resp.GetId()
}

func (s *rotatorSuite) deleteSocialGroup(id string) *grpcapi.DeleteSocialGroupResponse {
	resp, err := s.clientGrpc.DeleteSocialGroup(s.ctx, &grpcapi.DeleteSocialGroupRequest{Id: id})

	s.Require().NoError(err, "delete socialGroup error")
	s.Require().NotNil(resp)

	delete(s.socialGroups, id)
	return resp
}

func (s *rotatorSuite) attachBanner(slotID, bannerID string) {
	resp, err := s.clientGrpc.AttachBanner(s.ctx, &grpcapi.AttachBannerRequest{
		SlotId:   slotID,
		BannerId: bannerID,
	})

	s.Require().NoError(err)
	s.Require().Equal(code.Code_OK, resp.GetStatus().GetCode())
}

func getGrpcServerAddress() string {
	host := getEnv("TESTS_GRPC_SERVER_HOST", defaultGrpcServerHost)
	port := getEnv("TESTS_GRPC_SERVER_PORT", defaultGrpcServerPort)
	return fmt.Sprintf("%s:%s", host, port)
}

func getEnv(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func generateDescription(prefix string) string {
	return fmt.Sprintf("%s %d", prefix, time.Now().UnixMicro())
}
