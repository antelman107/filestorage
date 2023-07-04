//go:build integration
// +build integration

package integration_tests

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/antelman107/filestorage/internal/clients"
	"github.com/antelman107/filestorage/internal/config"
	"github.com/antelman107/filestorage/internal/providers"
	"github.com/antelman107/filestorage/internal/services"
)

const (
	testGatewayURL  = "http://localhost:9090"
	testStorage1URL = "http://localhost:9091"
	testStorage2URL = "http://localhost:9092"
)

type integrationTestSuite struct {
	suite.Suite
	db  *sqlx.DB
	wgs []*sync.WaitGroup
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(integrationTestSuite))
}

func (s *integrationTestSuite) SetupTest() {
	require.NoError(s.T(), clearStorage())
	s.db = nil
	s.wgs = nil
}

func (s *integrationTestSuite) Test_SmallFile_Success() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	s.initApps(ctx)
	require.NoError(s.T(), clearDB(s.db))

	fileName := "1B"
	file, err := getDataFile(fileName)
	require.NoError(s.T(), err)

	// Post servers
	gwClient := clients.NewGatewayV1Client(testGatewayURL)
	srv1, err := gwClient.PostServer(ctx, testStorage1URL)
	require.NoError(s.T(), err)
	require.Equal(s.T(), testStorage1URL, srv1.URL)

	srv2, err := gwClient.PostServer(ctx, testStorage2URL)
	require.NoError(s.T(), err)
	require.Equal(s.T(), testStorage2URL, srv2.URL)

	// Post file
	responseFile, err := gwClient.PostFile(ctx, fileName, file)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), responseFile.ID)
	require.Equal(s.T(), int64(1), responseFile.Size)
	require.Equal(s.T(), "1B", responseFile.Name)

	// Validate file storage
	files, err := getStorageFilesPaths()
	require.NoError(s.T(), err)
	require.Len(s.T(), files, 1)
	require.True(s.T(), strings.HasPrefix(files[0], "storage/storage_data1"), files[0])

	// Read single stored file
	storedFileBytes, err := os.ReadFile(files[0])
	require.NoError(s.T(), err)
	require.Equal(s.T(), "1", string(storedFileBytes))

	downloadedFile, err := os.CreateTemp(os.TempDir(), "downloaded")
	require.NoError(s.T(), err)

	// Get file
	require.NoError(s.T(), gwClient.GetFileContent(ctx, responseFile.ID.String(), downloadedFile))

	// Read downloaded file
	downloadedFileBytes, err := os.ReadFile(downloadedFile.Name())
	require.NoError(s.T(), err)
	require.Equal(s.T(), "1", string(downloadedFileBytes))

	deleteFileResponse, err := gwClient.DeleteFile(ctx, responseFile.ID.String())
	require.NoError(s.T(), err)
	require.Equal(s.T(), responseFile.ID, deleteFileResponse.ID)
	require.Equal(s.T(), responseFile.Size, deleteFileResponse.Size)

	cancelFunc()
	for _, wg := range s.wgs {
		wg.Wait()
	}
}

func (s *integrationTestSuite) Test_2MBFile_Success() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	s.initApps(ctx)

	fileName := "2MB"
	file, err := getDataFile(fileName)
	require.NoError(s.T(), err)

	// Post servers
	gwClient := clients.NewGatewayV1Client(testGatewayURL)
	srv1, err := gwClient.PostServer(ctx, testStorage1URL)
	require.NoError(s.T(), err)
	require.Equal(s.T(), testStorage1URL, srv1.URL)

	srv2, err := gwClient.PostServer(ctx, testStorage2URL)
	require.NoError(s.T(), err)
	require.Equal(s.T(), testStorage2URL, srv2.URL)

	// Post file
	responseFile, err := gwClient.PostFile(ctx, fileName, file)
	require.NoError(s.T(), err)
	require.Equal(s.T(), int64(2*1024*1024), responseFile.Size)
	require.Equal(s.T(), "2MB", responseFile.Name)

	// Validate file storage
	files, err := getStorageFilesPaths()
	require.NoError(s.T(), err)
	require.Len(s.T(), files, 6)

	downloadedFile, err := os.CreateTemp(os.TempDir(), "downloaded")
	require.NoError(s.T(), err)

	// Get file
	require.NoError(s.T(), gwClient.GetFileContent(ctx, responseFile.ID.String(), downloadedFile))
	fstat, err := downloadedFile.Stat()
	require.NoError(s.T(), err)
	require.Equal(s.T(), int64(2*1024*1024), fstat.Size())

	cancelFunc()
	for _, wg := range s.wgs {
		wg.Wait()
	}
}

func (s *integrationTestSuite) Test_GetFileContent_Failed_DownloadNonExistentFile() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	s.initApps(ctx)

	gwClient := clients.NewGatewayV1Client(testGatewayURL)

	// Get file
	err := gwClient.GetFileContent(ctx, uuid.New().String(), io.Discard)
	require.ErrorContains(s.T(), err, "http code is not OK: 404")

	cancelFunc()
	for _, wg := range s.wgs {
		wg.Wait()
	}
}

func (s *integrationTestSuite) Test_Post_Failed_NoServers() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	s.initApps(ctx)

	gwClient := clients.NewGatewayV1Client(testGatewayURL)

	fileName := "1B"
	file, err := getDataFile(fileName)
	require.NoError(s.T(), err)

	// Post file
	responseFile, err := gwClient.PostFile(ctx, fileName, file)
	require.Empty(s.T(), responseFile.ID)
	require.ErrorContains(s.T(), err, "http code is not OK: 412")

	cancelFunc()
	for _, wg := range s.wgs {
		wg.Wait()
	}
}

func (s *integrationTestSuite) Test_Delete_Failed_NonExistentFile() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	s.initApps(ctx)

	gwClient := clients.NewGatewayV1Client(testGatewayURL)

	// Delete file
	file, err := gwClient.DeleteFile(ctx, uuid.New().String())
	require.Empty(s.T(), file.ID)
	require.ErrorContains(s.T(), err, "http code is not OK: 404")

	cancelFunc()
	for _, wg := range s.wgs {
		wg.Wait()
	}
}

func (s *integrationTestSuite) initApps(ctx context.Context) {
	gwps := providers.DefaultGatewayAppProviders
	gwps.ConfigLoader = providers.NewDefaultConfigLoader("./configs/default")
	sqlxProv := &customSQLXProvider{}
	gwps.SqlxProvider = sqlxProv
	gwApp := services.NewGatewayApp(gwps)

	require.NoError(s.T(), gwApp.Init())

	wg1 := &sync.WaitGroup{}
	wg1.Add(1)
	go func() {
		defer wg1.Done()

		gwApp.Run(ctx)
	}()

	storageConfigLoader1 := customStorageConfigLoader{
		config: config.StorageConfig{
			StoragePath: "./storage/storage_data1",
			HTTP: config.HTTP{
				ListenPort: ":9091",
			},
		},
	}
	stps1 := providers.DefaultStorageAppProviders
	stps1.ConfigLoader = &storageConfigLoader1
	stApp1 := services.NewStorageApp(stps1)
	require.NoError(s.T(), stApp1.Init())

	wg2 := &sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		defer wg2.Done()

		stApp1.Run(ctx)
	}()

	storageConfigLoader2 := customStorageConfigLoader{
		config: config.StorageConfig{
			StoragePath: "./storage/storage_data2",
			HTTP: config.HTTP{
				ListenPort: ":9092",
			},
		},
	}
	stps2 := providers.DefaultStorageAppProviders
	stps2.ConfigLoader = &storageConfigLoader2
	stApp2 := services.NewStorageApp(stps2)
	require.NoError(s.T(), stApp2.Init())

	wg3 := &sync.WaitGroup{}
	wg3.Add(1)
	go func() {
		defer wg3.Done()

		stApp2.Run(ctx)
	}()

	require.True(s.T(), isURLHealthy(ctx, testGatewayURL))
	require.True(s.T(), isURLHealthy(ctx, testStorage1URL))
	require.True(s.T(), isURLHealthy(ctx, testStorage2URL))

	s.db = sqlxProv.db
	require.NoError(s.T(), clearDB(s.db))

	s.wgs = []*sync.WaitGroup{wg1, wg2, wg3}
}
