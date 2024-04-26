package tests

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"

	"github.com/EfimVelichkin/3rd_module_GO/03-04-umanager/internal/env"
	"github.com/EfimVelichkin/3rd_module_GO/03-04-umanager/internal/env/config"
)

type IntegrationTestSuite struct {
	suite.Suite
	conf       config.Config
	closer     *env.Closer
	pgPool     *dockertest.Pool
	pgRes      *dockertest.Resource
	mongoPool  *dockertest.Pool
	mongoRes   *dockertest.Resource
	rabbitPool *dockertest.Pool
	rabbitRes  *dockertest.Resource
}

func (s *IntegrationTestSuite) SetupSuite() {
	// Prepare Containers with required resources
	s.pgPool, s.pgRes = StartPG()
	s.mongoPool, s.mongoRes = StartMongo()
	s.rabbitPool, s.rabbitRes = StartRabbit()
	SetupEnv()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e, c, err := env.Setup(ctx)
	s.Assert().NoError(err)

	s.conf = e.Config
	s.closer = c

	go func() {
		lis, err := net.Listen("tcp", e.Config.LinksService.GRPCServer.Addr)
		s.Assert().NoError(err)

		err = e.LinksGRPCServer.Serve(lis)
		s.Assert().NoError(err)
	}()

	go func() {
		lis, err := net.Listen("tcp", e.Config.UsersService.GRPCServer.Addr)
		s.Assert().NoError(err)

		err = e.UsersGRPCServer.Serve(lis)
		s.Assert().NoError(err)
	}()

	go func() {
		defer e.APIGWHTTPServer.Close()
		err := e.APIGWHTTPServer.ListenAndServe()
		s.Assert().NoError(err)
	}()

	time.Sleep(time.Second) // Wait all goroutimes are running
}

func (s *IntegrationTestSuite) TearDownSuite() {
	defer Stop(s.pgPool, s.pgRes)
	defer Stop(s.mongoPool, s.mongoRes)
	defer Stop(s.rabbitPool, s.rabbitRes)
	s.closer.Close(context.Background())
}

func (s *IntegrationTestSuite) SetupTest() {
}

func TestIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	time.Sleep(5 * time.Second) // wait before previous containers will be stopped

	suite.Run(t, new(IntegrationTestSuite))
}
