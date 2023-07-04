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

type gatewayUploader struct {
	client      domain.StorageV1Client
	repo        domain.GatewayRepository
	concurrency int
	logger      *zap.Logger
}

func NewGatewayUploader(
	client domain.StorageV1Client,
	repo domain.GatewayRepository,
	concurrency int,
	logger *zap.Logger,
) domain.GatewayUploaderService {
	return &gatewayUploader{
		client:      client,
		repo:        repo,
		concurrency: concurrency,
		logger:      logger,
	}
}

func (s *gatewayUploader) Upload(ctx context.Context, chunks domain.Chunks, file domain.File, reader io.Reader) error {
	group, groupCtx := errgroup.WithContext(ctx)
	chunksWithDataChan := make(chan domain.ChunkWithData, len(chunks))
	for i := 0; i < s.concurrency; i++ {
		workerIndex := i
		group.Go(func() error {
		forloop:
			for {
				select {
				case <-groupCtx.Done():
					s.logger.Info("group context is finished")
					break forloop

				case chunk, ok := <-chunksWithDataChan:
					if !ok {
						s.logger.Info("chunks are processed, exiting")
						break forloop
					}

					if err := s.client.PostChunk(groupCtx, chunk); err != nil {
						if err == context.Canceled {
							return nil
						}

						s.logger.Info(
							"failed to post chunk",
							zap.Int64("size", chunk.Size),
							zap.Int("worker index", workerIndex),
							zap.Error(err),
						)
						return fmt.Errorf("failed to post chunk: %w", err)
					}

					if err := s.repo.UpdateUploadedChunk(groupCtx, chunk.ID.String(), chunk.Hash); err != nil {
						s.logger.Error(
							"failed to update chunk is_uploaded",
							zap.String("server url", chunk.ServerURL),
							zap.Int64("size", chunk.Size),
							zap.Int("chunk index", chunk.Index),
							zap.Int("size", workerIndex),
							zap.Error(err),
						)
						return fmt.Errorf("failed to update file is_uploaded %s: %w", file.ID.String(), err)
					}

					s.logger.Info(
						"uploaded chunk",
						zap.String("server url", chunk.ServerURL),
						zap.Int64("size", chunk.Size),
						zap.Int("chunk index", chunk.Index),
						zap.Int("size", workerIndex),
						zap.String("hash", chunk.Hash),
					)
				}
			}

			return nil
		})
	}

	for _, chunk := range chunks {
		chunkHash := md5.New()
		chunkTeeReader := io.TeeReader(reader, chunkHash)

		// Can't concurrently read from io.Reader. This part is sequential.
		buf := make([]byte, chunk.Size)
		if _, err := chunkTeeReader.Read(buf); err != nil {
			if err == io.EOF {
				s.logger.Info("finished reading input file")
				break
			}
			return fmt.Errorf("failed to read chunk: %w", err)
		}

		chunk.Hash = hex.EncodeToString(chunkHash.Sum(nil))
		chunksWithDataChan <- chunk.GetWithData(buf)

		s.logger.Info("read file chunk", zap.Int64("size", chunk.Size))
	}
	close(chunksWithDataChan)

	if err := group.Wait(); err != nil {
		return fmt.Errorf("failed to wait: %w", err)
	}

	s.logger.Info(
		"finished file uploading",
		zap.String("ID", file.ID.String()),
		zap.String("Name", file.Name),
		zap.Int64("Size", file.Size),
	)

	return nil
}
