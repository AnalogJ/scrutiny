package collector

import (
	"fmt"
	"net/url"

	"github.com/analogj/scrutiny/collector/pkg/common/shell"
	"github.com/analogj/scrutiny/collector/pkg/config"
	"github.com/analogj/scrutiny/collector/pkg/detect"
	webModels "github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/sirupsen/logrus"
)

type ZfsCollector struct {
	config config.Interface
	BaseCollector
	apiEndpoint *url.URL
	shell       shell.Interface
}

func CreateZfsCollector(appConfig config.Interface, logger *logrus.Entry, apiEndpoint string) (ZfsCollector, error) {
	apiEndpointUrl, err := url.Parse(apiEndpoint)
	if err != nil {
		return ZfsCollector{}, err
	}

	zc := ZfsCollector{
		config:      appConfig,
		apiEndpoint: apiEndpointUrl,
		BaseCollector: BaseCollector{
			logger: logger,
		},
		shell: shell.Create(),
	}

	return zc, nil
}

func (zc *ZfsCollector) Run() error {
	err := zc.Validate()
	if err != nil {
		return err
	}

	// Collect ZFS pool data
	return zc.CollectZfs()
}

func (zc *ZfsCollector) Validate() error {
	zc.logger.Infoln("Verifying ZFS tools availability")
	
	zfsDetector := detect.ZfsDetect{
		Logger: zc.logger,
		Config: zc.config,
		Shell:  zc.shell,
	}

	if !zfsDetector.IsZfsAvailable() {
		return fmt.Errorf("ZFS tools are not available on this system")
	}

	return nil
}

func (zc *ZfsCollector) CollectZfs() error {
	zc.logger.Infoln("Collecting ZFS pool data")

	zfsDetector := detect.ZfsDetect{
		Logger: zc.logger,
		Config: zc.config,
		Shell:  zc.shell,
	}

	// Skip if ZFS is not available
	if !zfsDetector.IsZfsAvailable() {
		zc.logger.Debug("ZFS tools not available, skipping ZFS collection")
		return nil
	}

	pools, err := zfsDetector.DetectZfsPools()
	if err != nil {
		return fmt.Errorf("error detecting ZFS pools: %v", err)
	}

	if len(pools) == 0 {
		zc.logger.Debug("No ZFS pools detected")
		return nil
	}

	zc.logger.Infof("Detected %d ZFS pools, publishing to API", len(pools))

	// Collect ZFS datasets
	datasets, err := zfsDetector.DetectZfsDatasets()
	if err != nil {
		zc.logger.Errorf("Error detecting ZFS datasets: %v", err)
		// Don't fail the entire collection if datasets fail
	} else {
		zc.logger.Infof("Detected %d ZFS datasets", len(datasets))
	}

	// Publish ZFS pool data to the API
	err = zc.PublishZfsPools(pools)
	if err != nil {
		return err
	}

	// Publish ZFS dataset data to the API if we have any
	if len(datasets) > 0 {
		return zc.PublishZfsDatasets(datasets)
	}

	return nil
}

func (zc *ZfsCollector) PublishZfsPools(pools []webModels.ZfsPool) error {
	zc.logger.Infoln("Publishing ZFS pool data")

	apiEndpoint, _ := url.Parse(zc.apiEndpoint.String())
	apiEndpoint, _ = apiEndpoint.Parse("api/zfs/pools/register")

	poolWrapper := webModels.ZfsPoolWrapper{
		Data: pools,
	}

	var respWrapper webModels.ZfsPoolWrapper
	err := zc.postJson(apiEndpoint.String(), poolWrapper, &respWrapper)
	if err != nil {
		zc.logger.Errorf("An error occurred while publishing ZFS pool data: %v", err)
		return err
	}

	if !respWrapper.Success {
		zc.logger.Errorln("API server rejected ZFS pool data")
		return fmt.Errorf("API server rejected ZFS pool data")
	}

	zc.logger.Infof("Successfully published %d ZFS pools", len(pools))
	return nil
}

func (zc *ZfsCollector) PublishZfsDatasets(datasets []webModels.ZfsDataset) error {
	zc.logger.Infoln("Publishing ZFS dataset data")

	apiEndpoint, _ := url.Parse(zc.apiEndpoint.String())
	apiEndpoint, _ = apiEndpoint.Parse("api/zfs/datasets/register")

	datasetWrapper := webModels.ZfsDatasetWrapper{
		Data: datasets,
	}

	var respWrapper webModels.ZfsDatasetWrapper
	err := zc.postJson(apiEndpoint.String(), datasetWrapper, &respWrapper)
	if err != nil {
		zc.logger.Errorf("An error occurred while publishing ZFS dataset data: %v", err)
		return err
	}

	if !respWrapper.Success {
		zc.logger.Errorln("API server rejected ZFS dataset data")
		return fmt.Errorf("API server rejected ZFS dataset data")
	}

	zc.logger.Infof("Successfully published %d ZFS datasets", len(datasets))
	return nil
}