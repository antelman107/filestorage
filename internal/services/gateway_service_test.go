package services

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antelman107/filestorage/internal/logger"
	"github.com/antelman107/filestorage/pkg/domain"
	"github.com/antelman107/filestorage/pkg/mocks"
)

type gatewayServiceTestSuite struct {
	suite.Suite
	repoMock       *mocks.GatewayRepository
	uploaderMock   *mocks.GatewayUploaderService
	downloaderMock *mocks.GatewayDownloaderService
	deleterMock    *mocks.GatewayDeleterService
}

func TestGatewayServiceTestSuite(t *testing.T) {
	suite.Run(t, new(gatewayServiceTestSuite))
}

func (s *gatewayServiceTestSuite) SetupTest() {
	s.repoMock = mocks.NewGatewayRepository(s.T())
	s.uploaderMock = mocks.NewGatewayUploaderService(s.T())
	s.downloaderMock = mocks.NewGatewayDownloaderService(s.T())
	s.deleterMock = mocks.NewGatewayDeleterService(s.T())
}

func (s *gatewayServiceTestSuite) TearDownTest() {
	s.repoMock.AssertExpectations(s.T())
	s.uploaderMock.AssertExpectations(s.T())
	s.downloaderMock.AssertExpectations(s.T())
	s.deleterMock.AssertExpectations(s.T())
}

func (s *gatewayServiceTestSuite) TestUploadFileSuccess() {
	ctx := context.TODO()

	zapLogger, err := logger.Get()
	require.NoError(s.T(), err)

	reader := bytes.NewBufferString("1234567")

	u := NewGatewayService(
		s.repoMock,
		s.uploaderMock,
		s.downloaderMock,
		s.deleterMock,
		6,
		2,
		zapLogger,
	)
	s.repoMock.On("InitTransaction").Return(nil, nil)
	s.repoMock.On("LockServers", ctx).Return(nil)
	s.repoMock.On("GetServersUsages", ctx).Return(domain.ServersSummary{
		Usages: domain.ServersUsages{
			{ServerID: testUUID1, ServerURL: "1", Usage: 0},
		},
		TotalUsage: 0,
	}, nil)

	s.repoMock.On("StoreFile", ctx, domain.File{
		ID: testUUID3,
	}).Return(nil)

	s.repoMock.On(
		"StoreChunks",
		ctx,
		mock.MatchedBy(func(chunks domain.Chunks) bool {
			require.Len(s.T(), chunks, 1)
			require.Equal(s.T(), testUUID1, chunks[0].ServerID)
			require.Equal(s.T(), "1", chunks[0].ServerURL)
			require.Empty(s.T(), chunks[0].Hash)
			return true
		}),
	).Return(nil)

	s.uploaderMock.On(
		"Upload",
		ctx,
		mock.Anything,
		domain.File{
			ID: testUUID3,
		},
		reader,
	).Return(nil)

	s.repoMock.On(
		"UpdateFileIsUploaded",
		ctx,
		testUUID3.String(),
	).Return(nil)

	err = u.UploadFile(ctx, domain.File{ID: testUUID3}, reader)
	require.NoError(s.T(), err)
}

func (s *gatewayServiceTestSuite) TestDownloadFileSuccess() {
	ctx := context.TODO()

	zapLogger, err := logger.Get()
	require.NoError(s.T(), err)

	writer := &bytes.Buffer{}

	u := NewGatewayService(
		s.repoMock,
		s.uploaderMock,
		s.downloaderMock,
		s.deleterMock,
		6,
		2,
		zapLogger,
	)
	file := domain.File{ID: testUUID3}
	chunks := domain.Chunks{
		{ID: testUUID1},
	}
	s.repoMock.On("GetUploadedFile", ctx, testUUID3.String()).Return(file, nil)
	s.repoMock.On("GetChunksByFileID", ctx, testUUID3.String()).Return(chunks, nil)

	s.downloaderMock.On(
		"Download",
		ctx,
		chunks,
		writer,
	).Return(nil)

	err = u.DownloadFile(ctx, testUUID3.String(), writer)
	require.NoError(s.T(), err)
}

func (s *gatewayServiceTestSuite) TestDeleteFileSuccess() {
	ctx := context.TODO()

	zapLogger, err := logger.Get()
	require.NoError(s.T(), err)

	u := NewGatewayService(
		s.repoMock,
		s.uploaderMock,
		s.downloaderMock,
		s.deleterMock,
		6,
		2,
		zapLogger,
	)
	file := domain.File{ID: testUUID3}
	chunks := domain.Chunks{
		{ID: testUUID1},
	}
	s.repoMock.On("InitTransaction").Return(nil, nil)
	s.repoMock.On("LockFile", ctx, testUUID3.String()).Return(nil)
	s.repoMock.On("GetChunksByFileID", ctx, testUUID3.String()).Return(chunks, nil)

	s.deleterMock.On(
		"Delete",
		ctx,
		chunks,
		testUUID3.String(),
	).Return(nil)

	s.repoMock.On("DeleteFile", ctx, testUUID3.String()).Return(file, nil)

	ret, err := u.DeleteFile(ctx, testUUID3.String())
	require.NoError(s.T(), err)
	require.Equal(s.T(), file, ret)
}
