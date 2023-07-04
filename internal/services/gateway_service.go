package services

import (
	"context"
	"fmt"
	"io"

	"go.uber.org/zap"

	"github.com/antelman107/filestorage/internal/database/repositories"
	"github.com/antelman107/filestorage/pkg/domain"
)

type gatewayService struct {
	repo               domain.GatewayRepository
	uploader           domain.GatewayUploaderService
	downloader         domain.GatewayDownloaderService
	deleter            domain.GatewayDeleterService
	numChunks          int64
	minFileSizeToSplit int64
	logger             *zap.Logger
}

func NewGatewayService(
	repo domain.GatewayRepository,
	uploader domain.GatewayUploaderService,
	downloader domain.GatewayDownloaderService,
	deleter domain.GatewayDeleterService,
	numChunks int64,
	minFileSizeToSplit int64,
	logger *zap.Logger,
) domain.GatewayService {
	return &gatewayService{
		repo:               repo,
		uploader:           uploader,
		downloader:         downloader,
		deleter:            deleter,
		numChunks:          numChunks,
		minFileSizeToSplit: minFileSizeToSplit,
		logger:             logger,
	}
}

// UploadFile uses properties from file and data from reader to split data into chunks and upload to storage servers.
func (s *gatewayService) UploadFile(ctx context.Context, file domain.File, reader io.Reader) error {
	var chunks domain.Chunks
	if err := repositories.WithTransaction(ctx, s.repo, func(ctx context.Context) error {
		// We are locking all server rows to avoid any concurrent operations with server during this transaction
		if err := s.repo.LockServers(ctx); err != nil {
			return fmt.Errorf("failed to lock servers: %w", err)
		}

		serverSummary, err := s.repo.GetServersUsages(ctx)
		if err != nil {
			return fmt.Errorf("failed to get server stats: %w", err)
		}

		if len(serverSummary.Usages) == 0 {
			return domain.ErrNoStorageSpace
		}

		if err := s.repo.StoreFile(ctx, file); err != nil {
			return fmt.Errorf("failed to store file: %w", err)
		}

		chunks = getChunks(chunksAnalysisInput{
			stats:              serverSummary,
			numChunks:          s.numChunks,
			minFileSizeToSplit: s.minFileSizeToSplit,
			fileSize:           file.Size,
			fileID:             file.ID,
		})

		if err := s.repo.StoreChunks(ctx, chunks); err != nil {
			return fmt.Errorf("failed to store chunks: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to apply transaction: %w", err)
	}

	if err := s.uploader.Upload(ctx, chunks, file, reader); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	if err := s.repo.UpdateFileIsUploaded(ctx, file.ID.String()); err != nil {
		return fmt.Errorf("failed to update file is_uploaded %s: %w", file.ID.String(), err)
	}

	return nil
}

// DownloadFile checks if file with specified id is stored in DB as uploaded one.
// If yes, file chunks content is downloaded from DB and is written to writer.
func (s *gatewayService) DownloadFile(ctx context.Context, id string, writer io.Writer) error {
	if _, err := s.repo.GetUploadedFile(ctx, id); err != nil {
		return fmt.Errorf("failed to get uploaded file: %w", err)
	}

	chunks, err := s.repo.GetChunksByFileID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get chunks: %w", err)
	}

	if err := s.downloader.Download(ctx, chunks, writer); err != nil {
		return fmt.Errorf("failed to download chunks: %w", err)
	}

	return nil
}

// DeleteFile deletes all file chunks from storage servers and also local DB data related to chunks and file.
func (s *gatewayService) DeleteFile(ctx context.Context, id string) (domain.File, error) {
	res := domain.File{}
	if err := repositories.WithTransaction(ctx, s.repo, func(ctx context.Context) error {
		if err := s.repo.LockFile(ctx, id); err != nil {
			return fmt.Errorf("failed to lock file: %w", err)
		}

		chunks, err := s.repo.GetChunksByFileID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get chunks: %w", err)
		}

		if err := s.deleter.Delete(ctx, chunks, id); err != nil {
			return fmt.Errorf("failed to delete chunks: %w", err)
		}

		file, err := s.repo.DeleteFile(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete file %s: %w", id, err)
		}
		res = file

		return nil
	}); err != nil {
		return domain.File{}, fmt.Errorf("failed to apply transaction: %w", err)
	}

	return res, nil
}

func (s *gatewayService) AddServer(ctx context.Context, server domain.Server) error {
	if err := s.repo.StoreServer(ctx, server); err != nil {
		return fmt.Errorf("failed to store server: %w", err)
	}

	return nil
}
