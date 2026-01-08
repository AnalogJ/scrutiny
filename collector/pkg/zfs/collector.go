package zfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/analogj/scrutiny/collector/pkg/config"
	"github.com/analogj/scrutiny/collector/pkg/errors"
	"github.com/analogj/scrutiny/collector/pkg/zfs/detect"
	"github.com/analogj/scrutiny/collector/pkg/zfs/models"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

// Collector handles ZFS pool collection
type Collector struct {
	config      config.Interface
	logger      *logrus.Entry
	apiEndpoint *url.URL
	httpClient  *http.Client
}

// CreateCollector creates a new ZFS collector
func CreateCollector(appConfig config.Interface, logger *logrus.Entry, apiEndpoint string) (*Collector, error) {
	apiEndpointUrl, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil, err
	}

	timeout := 60
	if appConfig != nil && appConfig.IsSet("api.timeout") {
		timeout = appConfig.GetAPITimeout()
	}

	collector := &Collector{
		config:      appConfig,
		logger:      logger,
		apiEndpoint: apiEndpointUrl,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}

	return collector, nil
}

// Run executes the ZFS collection
func (c *Collector) Run() error {
	c.logger.Infoln("Starting ZFS pool collection")

	// Detect pools
	detector := detect.Detect{
		Logger: c.logger,
		Config: c.config,
	}

	pools, err := detector.Start()
	if err != nil {
		return err
	}

	if len(pools) == 0 {
		c.logger.Infoln("No ZFS pools found")
		return nil
	}

	c.logger.Infof("Found %d ZFS pool(s)", len(pools))

	// Filter pools with empty GUID
	validPools := lo.Filter[models.ZFSPool](pools, func(pool models.ZFSPool, _ int) bool {
		return len(pool.GUID) > 0
	})

	// Register pools with API
	poolWrapper, err := c.RegisterPools(validPools)
	if err != nil {
		return err
	}

	if !poolWrapper.Success {
		c.logger.Errorln("An error occurred while registering pools")
		return errors.ApiServerCommunicationError("An error occurred while registering pools")
	}

	// Upload metrics for each registered pool
	for _, pool := range poolWrapper.Data {
		if err := c.UploadMetrics(pool); err != nil {
			c.logger.Errorf("Failed to upload metrics for pool %s: %v", pool.Name, err)
			// Continue with other pools
		}
	}

	c.logger.Infoln("ZFS collection completed")
	return nil
}

// RegisterPools registers detected pools with the API
func (c *Collector) RegisterPools(pools []models.ZFSPool) (*models.ZFSPoolWrapper, error) {
	c.logger.Infoln("Sending detected pools to API for registration")

	apiEndpoint, _ := url.Parse(c.apiEndpoint.String())
	apiEndpoint, _ = apiEndpoint.Parse("api/zfs/pools/register")

	wrapper := models.ZFSPoolWrapper{
		Data: pools,
	}

	jsonData, err := json.Marshal(wrapper)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pools: %w", err)
	}

	c.logger.Debugf("Registering pools: %s", string(jsonData))

	resp, err := c.httpClient.Post(apiEndpoint.String(), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Errorf("Failed to register pools: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var responseWrapper models.ZFSPoolWrapper
	if err := json.NewDecoder(resp.Body).Decode(&responseWrapper); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &responseWrapper, nil
}

// UploadMetrics uploads metrics for a specific pool
func (c *Collector) UploadMetrics(pool models.ZFSPool) error {
	c.logger.Infof("Uploading metrics for pool %s (%s)", pool.Name, pool.GUID)

	apiEndpoint, _ := url.Parse(c.apiEndpoint.String())
	apiEndpoint, _ = apiEndpoint.Parse(fmt.Sprintf("api/zfs/pool/%s/metrics", strings.ToLower(pool.GUID)))

	jsonData, err := json.Marshal(pool)
	if err != nil {
		return fmt.Errorf("failed to marshal pool metrics: %w", err)
	}

	c.logger.Debugf("Uploading pool metrics: %s", string(jsonData))

	resp, err := c.httpClient.Post(apiEndpoint.String(), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Errorf("Failed to upload metrics for pool %s: %v", pool.Name, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	c.logger.Infof("Successfully uploaded metrics for pool %s", pool.Name)
	return nil
}
