package services

import (
	"github.com/google/uuid"

	"github.com/antelman107/filestorage/pkg/domain"
)

type chunksAnalysisInput struct {
	stats              domain.ServersSummary
	numChunks          int64
	minFileSizeToSplit int64
	fileSize           int64
	fileID             uuid.UUID
}

// getChunks analyzes chunksAnalysisInput and returns list of chunks.
func getChunks(in chunksAnalysisInput) domain.Chunks {
	numChunks := in.numChunks
	if in.fileSize < in.minFileSizeToSplit {
		return domain.Chunks{
			{
				ID:        uuid.New(),
				FileID:    in.fileID,
				ServerID:  in.stats.Usages[0].ServerID,
				ServerURL: in.stats.Usages[0].ServerURL,
				Index:     0,
				Size:      in.fileSize,
			},
		}
	}

	chunkSize := in.fileSize / numChunks

	totalBytes := int64(0)
	chunks := make(domain.Chunks, numChunks)

	for i := range chunks {
		server := in.stats.Usages[i%len(in.stats.Usages)]

		chunks[i] = domain.Chunk{
			ID:        uuid.New(),
			FileID:    in.fileID,
			ServerID:  server.ServerID,
			ServerURL: server.ServerURL,
			Index:     i,
			Size:      chunkSize,
		}

		totalBytes += chunkSize
	}

	// Distribute remaining bytes across servers
	bytesDiff := in.fileSize - totalBytes
	if bytesDiff > 0 {
		for i := 0; int64(i) < bytesDiff; i++ {
			chunks[i%len(chunks)].Size += 1
		}
	}

	return chunks
}
