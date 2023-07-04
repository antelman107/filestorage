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

type uploaderTestSuite struct {
	suite.Suite
	storageClientMock *mocks.StorageV1Client
	repoMock          *mocks.GatewayRepository
}

func TestUploaderTestSuite(t *testing.T) {
	suite.Run(t, new(uploaderTestSuite))
}

func (s *uploaderTestSuite) SetupTest() {
	s.storageClientMock = mocks.NewStorageV1Client(s.T())
	s.repoMock = mocks.NewGatewayRepository(s.T())
}

func (s *uploaderTestSuite) TearDownTest() {
	s.storageClientMock.AssertExpectations(s.T())
	s.repoMock.AssertExpectations(s.T())
}

func (s *uploaderTestSuite) TestSuccess() {
	ctx := context.TODO()

	s.storageClientMock.On("PostChunk", mock.Anything, domain.ChunkWithData{
		Chunk: domain.Chunk{
			ID:        testUUID1,
			ServerURL: "1",
			Size:      1,
			Hash:      "c4ca4238a0b923820dcc509a6f75849b",
		},
		Data: []byte("1"),
	}).Return(nil)

	s.storageClientMock.On("PostChunk", mock.Anything, domain.ChunkWithData{
		Chunk: domain.Chunk{
			ID:        testUUID2,
			ServerURL: "2",
			Size:      2,
			Hash:      "37693cfc748049e45d87b8c7d8b9aacd",
		},
		Data: []byte("23"),
	}).Return(nil)

	s.repoMock.On(
		"UpdateUploadedChunk",
		mock.Anything,
		testUUID1.String(),
		"c4ca4238a0b923820dcc509a6f75849b",
	).Return(nil)
	s.repoMock.On(
		"UpdateUploadedChunk",
		mock.Anything,
		testUUID2.String(),
		"37693cfc748049e45d87b8c7d8b9aacd",
	).Return(nil)

	zapLogger, err := logger.Get()
	require.NoError(s.T(), err)

	reader := bytes.NewBufferString("123")

	u := NewGatewayUploader(s.storageClientMock, s.repoMock, 2, zapLogger)
	err = u.Upload(ctx, domain.Chunks{
		{ID: testUUID1, ServerURL: "1", Size: 1},
		{ID: testUUID2, ServerURL: "2", Size: 2},
	}, domain.File{
		Name: "fileName",
		Size: 3,
	}, reader)

	require.NoError(s.T(), err)
}
