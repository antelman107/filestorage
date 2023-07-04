package services

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type storageServiceTestSuite struct {
	suite.Suite
}

func TestStorageServiceTestSuite(t *testing.T) {
	suite.Run(t, new(storageServiceTestSuite))
}

func (s *storageServiceTestSuite) TestStoreChunkSuccess() {
	tempDir := os.TempDir()
	testFileName := "123456"
	_ = os.Remove(os.TempDir() + "12/34/56/" + testFileName)

	srv := NewStorageService(tempDir)

	reader := bytes.NewBufferString("123")
	err := srv.StoreChunk(context.TODO(), testFileName, reader)

	b, err := os.ReadFile(os.TempDir() + "12/34/56/" + testFileName)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "123", string(b))
}

func (s *storageServiceTestSuite) TestGetChunkSuccess() {
	tempDir := os.TempDir()
	testFileName := "123456"
	_ = os.Remove(os.TempDir() + "12/34/56/" + testFileName)

	srv := NewStorageService(tempDir)

	reader := bytes.NewBufferString("123")
	require.NoError(s.T(), srv.StoreChunk(context.TODO(), testFileName, reader))

	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	require.NoError(s.T(), srv.GetChunk(context.TODO(), testFileName, writer))
	require.NoError(s.T(), writer.Flush())

	require.Equal(s.T(), "123", buffer.String())
}
