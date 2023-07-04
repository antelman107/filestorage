package services

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antelman107/filestorage/internal/logger"
	"github.com/antelman107/filestorage/pkg/domain"
	"github.com/antelman107/filestorage/pkg/mocks"
)

type downloaderTestSuite struct {
	suite.Suite
	storageClientMock *mocks.StorageV1Client
}

func TestDownloaderTestSuite(t *testing.T) {
	suite.Run(t, new(downloaderTestSuite))
}

func (s *downloaderTestSuite) SetupTest() {
	s.storageClientMock = mocks.NewStorageV1Client(s.T())
}

func (s *downloaderTestSuite) TearDownTest() {
	s.storageClientMock.AssertExpectations(s.T())
}

func (s *downloaderTestSuite) TestSuccess() {
	ctx := context.TODO()
	s.storageClientMock.On("GetChunk", mock.Anything, domain.Chunk{
		ServerURL: "1",
		Size:      1,
		Hash:      "c4ca4238a0b923820dcc509a6f75849b",
	}).Return(&http.Response{
		Body: io.NopCloser(bytes.NewBufferString("1")),
	}, nil)

	s.storageClientMock.On("GetChunk", mock.Anything, domain.Chunk{
		ServerURL: "2",
		Size:      2,
		Hash:      "37693cfc748049e45d87b8c7d8b9aacd",
	}).Return(&http.Response{
		Body: io.NopCloser(bytes.NewBufferString("23")),
	}, nil)

	zapLogger, err := logger.Get()
	require.NoError(s.T(), err)

	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)

	u := NewGatewayDownloader(s.storageClientMock, zapLogger)
	err = u.Download(ctx, domain.Chunks{
		{
			ServerURL: "1",
			Size:      1,
			Hash:      "c4ca4238a0b923820dcc509a6f75849b",
		},
		{
			ServerURL: "2",
			Size:      2,
			Hash:      "37693cfc748049e45d87b8c7d8b9aacd",
		},
	}, writer)
	require.NoError(s.T(), writer.Flush())

	require.NoError(s.T(), err)
	require.Equal(s.T(), "123", buffer.String())
}
