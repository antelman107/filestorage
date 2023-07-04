package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/antelman107/filestorage/pkg/domain"
)

type gatewayDeleter struct {
	client      domain.StorageV1Client
	repo        domain.GatewayRepository
	concurrency int
	logger      *zap.Logger
}

func NewGatewayDeleter(
	client domain.StorageV1Client,
	repo domain.GatewayRepository,
	concurrency int,
	logger *zap.Logger,
) domain.GatewayDeleterService {
	return &gatewayDeleter{
		client:      client,
		repo:        repo,
		concurrency: concurrency,
		logger:      logger,
	}
}

// Delete method sends requests to delete all chunks related to file using concurrent worker pool.
func (s *gatewayDeleter) Delete(ctx context.Context, chunks domain.Chunks, fileID string) error {
	group, groupCtx := errgroup.WithContext(ctx)
	chunksChan := make(chan domain.Chunk, len(chunks))
	for i := 0; i < s.concurrency; i++ {
		workerIndex := i
		group.Go(func() error {
		forloop:
			for {
				select {
				case <-groupCtx.Done():
					s.logger.Info("group context is finished")
					break forloop

				case chunk, ok := <-chunksChan:
					if !ok {
						s.logger.Info("chunks are processed, exiting")
						break forloop
					}

					if err := s.client.DeleteChunk(groupCtx, chunk); err != nil {
						if err == context.Canceled {
							return nil
						}

						s.logger.Info(
							"failed to post chunk",
							zap.Int64("size", chunk.Size),
							zap.Int("worker index", workerIndex),
							zap.Error(err),
						)
						return fmt.Errorf("failed to delete chunk: %w", err)
					}

					if err := s.repo.DeleteChunk(groupCtx, chunk.ID.String()); err != nil {
						s.logger.Error(
							"failed to delete chunk",
							zap.String("server url", chunk.ServerURL),
							zap.Int64("size", chunk.Size),
							zap.Int("chunk index", chunk.Index),
							zap.Int("size", workerIndex),
							zap.Error(err),
						)
						return fmt.Errorf("failed to delete chunk %s: %w", fileID, err)
					}

					s.logger.Info(
						"deleted chunk",
						zap.String("server url", chunk.ServerURL),
						zap.Int64("size", chunk.Size),
						zap.Int("chunk index", chunk.Index),
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
		"finished chunks deleting",
		zap.String("ID", fileID),
	)

	return nil
}
