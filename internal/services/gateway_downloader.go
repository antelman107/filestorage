package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/antelman107/filestorage/pkg/domain"
)

// gatewayDownloaderConcurrency > 1 only makes sense with io.WriterAt as output.
// For example, if we need to store chunks into OS file, we could use io.WriterAt.
// With http.ResponseWriter that implements io.Writer but not io.WriterAt
// it only makes sense to use consequent single-threaded writing.
const gatewayDownloaderConcurrency = 1

type gatewayDownloader struct {
	client domain.StorageV1Client
	logger *zap.Logger
}

func NewGatewayDownloader(
	client domain.StorageV1Client,
	logger *zap.Logger,
) domain.GatewayDownloaderService {
	return &gatewayDownloader{
		client: client,
		logger: logger,
	}
}

// Download downloads chunks from storage server one by one and writes the content to writer.
func (s *gatewayDownloader) Download(ctx context.Context, chunks domain.Chunks, writer io.Writer) error {
	group, groupCtx := errgroup.WithContext(ctx)
	chunksChan := make(chan domain.Chunk, len(chunks))
	for i := 0; i < gatewayDownloaderConcurrency; i++ {
		workerIndex := i
		group.Go(func() error {
		forloop:
			for {
				select {
				case <-groupCtx.Done():
					s.logger.Debug("group context is finished")
					break forloop

				case chunk, ok := <-chunksChan:
					if !ok {
						s.logger.Debug("chunks are processed, exiting")
						break forloop
					}

					response, err := s.client.GetChunk(groupCtx, chunk)
					if err != nil {
						if err == context.Canceled {
							return nil
						}

						s.logger.Error(
							"failed to post chunk",
							zap.Int64("size", chunk.Size),
							zap.Int("worker index", workerIndex),
							zap.Error(err),
						)
						return fmt.Errorf("failed to get chunk: %w", err)
					}

					chunkHash := md5.New()
					chunkTeeReader := io.TeeReader(response.Body, chunkHash)

					if _, err := io.Copy(writer, chunkTeeReader); err != nil {
						_ = response.Body.Close()
						return fmt.Errorf("failed to io.Copy chunk: %w", err)
					}
					_ = response.Body.Close()

					readerHash := hex.EncodeToString(chunkHash.Sum(nil))
					if chunk.Hash != readerHash {
						return fmt.Errorf("downloaded chunk hash is invalid %s != %s", chunk.Hash, readerHash)
					}

					s.logger.Debug(
						"downloaded chunk",
						zap.String("server url", chunk.ServerURL),
						zap.Int64("size", chunk.Size),
						zap.Int("chunk index", chunk.Index),
						zap.Int("size", workerIndex),
						zap.Int("size", workerIndex),
					)
				}
			}

			return nil
		})
	}

	for _, chunk := range chunks {
		chunksChan <- chunk
	}
	close(chunksChan)

	if err := group.Wait(); err != nil {
		return fmt.Errorf("failed to wait: %w", err)
	}

	s.logger.Info(
		"finished file downloading chunks",
		zap.Int("Size", len(chunks)),
	)

	return nil
}
