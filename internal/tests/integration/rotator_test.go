//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"github.com/ekhvalov/otus-banners-rotation/internal/app"
	"github.com/ekhvalov/otus-banners-rotation/internal/environment/config"
	"github.com/ekhvalov/otus-banners-rotation/internal/environment/queue/rabbitmq"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"

	grpcapi "github.com/ekhvalov/otus-banners-rotation/pkg/api/grpc"
	"github.com/stretchr/testify/suite"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultGrpcServerHost    = "localhost"
	defaultGrpcServerPort    = "8081"
	defaultRabbitmqHost      = "localhost"
	defaultRabbitmqPort      = "5672"
	defaultRabbitmqUsername  = "guest"
	defaultRabbitmqPassword  = "guest"
	defaultRabbitmqQueueName = "events"
)

func TestRotator(t *testing.T) {
	suite.Run(t, new(rotatorSuite))
}

type rotatorSuite struct {
	suite.Suite
	ctx           context.Context
	cancel        context.CancelFunc
	tick          time.Duration
	waitFor       time.Duration
	clientGrpc    grpcapi.RotatorClient
	banners       map[string]struct{}
	slots         map[string]struct{}
	socialGroups  map[string]struct{}
	queueConsumer *rabbitmq.Consumer
}

func (s *rotatorSuite) SetupSuite() {
	s.tick = time.Millisecond * 100
	s.waitFor = s.tick * 10 * 30
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

	s.queueConsumer = makeQueueConsumer()
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

func (s *rotatorSuite) Test_AllBannersSelected() {
	slotID := s.createSlot()
	groupID := s.createSocialGroup()
	banners := make(map[string]int)
	bannersCount := 100
	for i := 0; i < bannersCount; i++ {
		bannerID := s.createBanner()
		s.attachBanner(slotID, bannerID)
		banners[bannerID] = 0
	}

	for i := 0; i < bannersCount; i++ {
		resp, err := s.clientGrpc.SelectBanner(s.ctx, &grpcapi.SelectBannerRequest{
			SlotId:        slotID,
			SocialGroupId: groupID,
		})
		s.Require().NoError(err)
		s.Require().Equal(code.Code_OK, resp.GetStatus().GetCode())
		s.Require().NotEmpty(resp.GetBannerId())
		banners[resp.GetBannerId()]++
	}

	for bannerID, count := range banners {
		s.Require().Greater(count, 0, fmt.Sprintf("banner was not selecred: %s", bannerID))
	}
}

func (s *rotatorSuite) Test_PopularBannerSelectedFrequently() {
	slotID := s.createSlot()
	socialGroupID := s.createSocialGroup()
	bannersCTR := make(map[string]float64)
	bannersCount := 100
	popularBannerID := s.createBanner()
	s.attachBanner(slotID, popularBannerID)
	bannersCTR[popularBannerID] = 0.7
	defaultCTR := 0.3
	bannersSelects := make(map[string]int)
	for i := 0; i < bannersCount-1; i++ {
		bannerID := s.createBanner()
		bannersCTR[bannerID] = defaultCTR
		bannersSelects[bannerID] = 0
	}

	selectsTotal := bannersCount * 10
	for i := 0; i < selectsTotal; i++ {
		selectResponse, err := s.clientGrpc.SelectBanner(s.ctx, &grpcapi.SelectBannerRequest{
			SlotId:        slotID,
			SocialGroupId: socialGroupID,
		})
		s.Require().NoError(err)
		s.Require().Equal(code.Code_OK, selectResponse.GetStatus().GetCode())
		bannerID := selectResponse.GetBannerId()
		s.Require().NotEmpty(bannerID)
		bannersSelects[bannerID]++
		if rand.Float64() < bannersCTR[bannerID] {
			clickResponse, err := s.clientGrpc.ClickBanner(s.ctx, &grpcapi.ClickBannerRequest{
				SlotId:        slotID,
				BannerId:      bannerID,
				SocialGroupId: socialGroupID,
			})
			s.Require().NoError(err)
			s.Require().Equal(code.Code_OK, clickResponse.GetStatus().GetCode())
		}
	}

	selectRatios := make([]float64, 0)
	for _, selectsCount := range bannersSelects {
		selectRatios = append(selectRatios, float64(selectsCount)/float64(selectsTotal))
	}
	sort.Float64s(selectRatios)
	popularBannerRatio := float64(bannersSelects[popularBannerID]) / float64(selectsTotal)
	s.Require().Equal(selectRatios[bannersCount-1], popularBannerRatio)
	medianRatio := selectRatios[bannersCount/2]
	s.Require().GreaterOrEqual(popularBannerRatio, medianRatio*15.0)
}

func (s *rotatorSuite) Test_QueueEvents() {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()
	var eventsCh <-chan app.Event
	var err error
	s.Require().Eventually(func() bool {
		eventsCh, err = s.queueConsumer.Subscribe(ctx)
		return err == nil
	}, s.waitFor, s.tick)
	s.drainChannel(eventsCh)

	slotID := s.createSlot()
	bannerID := s.createBanner()
	socialGroupID := s.createSocialGroup()
	s.attachBanner(slotID, bannerID)

	_, err = s.clientGrpc.SelectBanner(s.ctx, &grpcapi.SelectBannerRequest{
		SlotId:        slotID,
		SocialGroupId: socialGroupID,
	})

	s.Require().NoError(err)
	event := s.getEvent(eventsCh)
	s.Require().Equal(app.EventSelect, event.Type)
	s.Require().Equal(bannerID, event.BannerID)
	s.Require().Equal(slotID, event.SlotID)
	s.Require().Equal(socialGroupID, event.SocialGroupID)

	_, err = s.clientGrpc.ClickBanner(s.ctx, &grpcapi.ClickBannerRequest{
		SlotId:        slotID,
		BannerId:      bannerID,
		SocialGroupId: socialGroupID,
	})

	s.Require().NoError(err)
	event = s.getEvent(eventsCh)
	s.Require().Equal(app.EventClick, event.Type)
	s.Require().Equal(bannerID, event.BannerID)
	s.Require().Equal(slotID, event.SlotID)
	s.Require().Equal(socialGroupID, event.SocialGroupID)
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

func makeQueueConsumer() *rabbitmq.Consumer {
	v, err := config.NewViper("", "TESTS", config.DefaultEnvKeyReplacer)
	if err != nil {
		panic(fmt.Sprintf("create viper error: %v", err))
	}
	cfg := rabbitmq.NewConfig(v)
	defaults := map[string]string{
		"TESTS_RABBITMQ_HOST":       defaultRabbitmqHost,
		"TESTS_RABBITMQ_PORT":       defaultRabbitmqPort,
		"TESTS_RABBITMQ_USERNAME":   defaultRabbitmqUsername,
		"TESTS_RABBITMQ_PASSWORD":   defaultRabbitmqPassword,
		"TESTS_RABBITMQ_QUEUE_NAME": defaultRabbitmqQueueName,
	}
	for key, value := range defaults {
		if _, ok := os.LookupEnv(key); !ok {
			err = os.Setenv(key, value)
			if err != nil {
				panic(fmt.Errorf("set env '%s' error: %w", key, err))
			}
		}
	}
	return rabbitmq.NewConsumer(cfg)
}

func (s *rotatorSuite) drainChannel(ch <-chan app.Event) {
	for {
		select {
		case <-ch:
		default:
			timeout := time.NewTimer(s.tick * 10)
			select {
			case <-ch:
				timeout.Stop()
			case <-timeout.C:
				return
			}
		}
	}
}

func (s *rotatorSuite) getEvent(eventsCh <-chan app.Event) app.Event {
	var event *app.Event
	s.Require().Eventually(func() bool {
		select {
		case e := <-eventsCh:
			event = &e
		default:
		}
		return event != nil
	}, s.waitFor, s.tick, "can not get an event from the queue")
	return *event
}
