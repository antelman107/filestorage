package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antelman107/filestorage/pkg/domain"
)

var (
	testUUID1 = uuid.New()
	testUUID2 = uuid.New()
	testUUID3 = uuid.New()
	testUUID4 = uuid.New()
)

type chunkStrategyTestSuite struct {
	suite.Suite
}

func TestChunkStrategyTestSuite(t *testing.T) {
	suite.Run(t, new(chunkStrategyTestSuite))
}

func (s *chunkStrategyTestSuite) Test_SmallFile() {
	chunks := getChunks(chunksAnalysisInput{
		stats: domain.ServersSummary{
			Usages: domain.ServersUsages{
				{
					ServerID: testUUID1,
					Usage:    0,
				},
				{
					ServerID: testUUID2,
					Usage:    0,
				},
			},
			TotalUsage: 0,
		},
		numChunks:          6,
		minFileSizeToSplit: 2,
		fileSize:           1,
		fileID:             testUUID3,
	})

	require.Len(s.T(), chunks, 1)
	require.Equal(s.T(), testUUID1, chunks[0].ServerID)
	require.Equal(s.T(), int64(1), chunks[0].Size)
	require.Equal(s.T(), 0, chunks[0].Index)
	require.Equal(s.T(), testUUID3, chunks[0].FileID)
	require.Equal(s.T(), false, chunks[0].IsUploaded)
}

func (s *chunkStrategyTestSuite) Test_Split_Round() {
	chunks := getChunks(chunksAnalysisInput{
		stats: domain.ServersSummary{
			Usages: domain.ServersUsages{
				{
					ServerID: testUUID1,
					Usage:    2,
				},
				{
					ServerID: testUUID2,
					Usage:    3,
				},
				{
					ServerID: testUUID3,
					Usage:    4,
				},
				{
					ServerID: testUUID3,
					Usage:    5,
				},
			},
			TotalUsage: 14,
		},
		numChunks:          6,
		minFileSizeToSplit: 5,
		fileSize:           6,
		fileID:             testUUID3,
	})

	require.Len(s.T(), chunks, 6)
	totalBytes := int64(0)
	for i := 0; i < 6; i++ {
		require.Equal(s.T(), int64(1), chunks[i].Size)
		totalBytes += chunks[i].Size
	}

	require.Equal(s.T(), int64(6), totalBytes)
}

func (s *chunkStrategyTestSuite) Test_Split_NotRound() {
	chunks := getChunks(chunksAnalysisInput{
		stats: domain.ServersSummary{
			Usages: domain.ServersUsages{
				{
					ServerID: testUUID1,
					Usage:    2,
				},
				{
					ServerID: testUUID2,
					Usage:    3,
				},
				{
					ServerID: testUUID3,
					Usage:    4,
				},
				{
					ServerID: testUUID3,
					Usage:    5,
				},
			},
			TotalUsage: 14,
		},
		numChunks:          6,
		minFileSizeToSplit: 5,
		fileSize:           7,
		fileID:             testUUID3,
	})

	require.Len(s.T(), chunks, 6)
	totalBytes := int64(0)
	for i := 0; i < 6; i++ {
		totalBytes += chunks[i].Size
	}

	require.Equal(s.T(), int64(7), totalBytes)
}

func (s *chunkStrategyTestSuite) Test_Split_NotRound_LargeFile() {
	chunks := getChunks(chunksAnalysisInput{
		stats: domain.ServersSummary{
			Usages: domain.ServersUsages{
				{
					ServerID: testUUID1,
					Usage:    2,
				},
				{
					ServerID: testUUID2,
					Usage:    3,
				},
				{
					ServerID: testUUID3,
					Usage:    4,
				},
				{
					ServerID: testUUID3,
					Usage:    5,
				},
			},
			TotalUsage: 14,
		},
		numChunks:          6,
		minFileSizeToSplit: 5,
		fileSize:           113,
		fileID:             testUUID3,
	})

	require.Len(s.T(), chunks, 6)
	totalBytes := int64(0)
	for i := 0; i < 6; i++ {
		totalBytes += chunks[i].Size
	}

	require.Equal(s.T(), int64(113), totalBytes)
}
