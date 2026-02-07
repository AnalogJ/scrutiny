package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	mock_config "github.com/analogj/scrutiny/webapp/backend/pkg/config/mock"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/web"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

/*
All tests in this file require the existance of a influxDB listening on port 8086

docker run --rm -it -p 8086:8086 \
-e DOCKER_INFLUXDB_INIT_MODE=setup \
-e DOCKER_INFLUXDB_INIT_USERNAME=admin \
-e DOCKER_INFLUXDB_INIT_PASSWORD=password12345 \
-e DOCKER_INFLUXDB_INIT_ORG=scrutiny \
-e DOCKER_INFLUXDB_INIT_BUCKET=metrics \
-e DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=my-super-secret-auth-token \
influxdb:2.0
*/

//func TestMain(m *testing.M) {
//	setup()
//	code := m.Run()
//	shutdown()
//	os.Exit(code)
//}

// InfluxDB will throw an error/ignore any submitted data with a timestamp older than the
// retention period. Lets fix this by opening test files, modifying the timestamp and returning an io.Reader
func helperReadSmartDataFileFixTimestamp(t *testing.T, smartDataFilepath string) io.Reader {
	metricsfile, err := os.Open(smartDataFilepath)
	require.NoError(t, err)

	metricsFileData, err := ioutil.ReadAll(metricsfile)
	require.NoError(t, err)

	//unmarshal because we need to change the timestamp
	var smartData collector.SmartInfo
	err = json.Unmarshal(metricsFileData, &smartData)
	require.NoError(t, err)
	smartData.LocalTime.TimeT = time.Now().Unix()
	updatedSmartDataBytes, err := json.Marshal(smartData)

	return bytes.NewReader(updatedSmartDataBytes)
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type ServerTestSuite struct {
	suite.Suite
	Basepath string
}

func TestServerTestSuite_WithEmptyBasePath(t *testing.T) {
	emptyBasePathSuite := new(ServerTestSuite)
	emptyBasePathSuite.Basepath = ""
	suite.Run(t, emptyBasePathSuite)
}

func TestServerTestSuite_WithCustomBasePath(t *testing.T) {
	emptyBasePathSuite := new(ServerTestSuite)
	emptyBasePathSuite.Basepath = "/basepath"
	suite.Run(t, emptyBasePathSuite)
}

func (suite *ServerTestSuite) TestHealthRoute() {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(suite.T())
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().SetDefault(gomock.Any(), gomock.Any()).AnyTimes()
	fakeConfig.EXPECT().UnmarshalKey(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	fakeConfig.EXPECT().GetString("web.database.location").Return(path.Join(parentPath, "scrutiny_test.db")).AnyTimes()
	fakeConfig.EXPECT().GetString("web.src.frontend.path").Return(parentPath).AnyTimes()
	fakeConfig.EXPECT().GetString("web.listen.basepath").Return(suite.Basepath).AnyTimes()

	fakeConfig.EXPECT().GetString("web.influxdb.scheme").Return("http").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.port").Return("8086").AnyTimes()
	fakeConfig.EXPECT().IsSet("web.influxdb.token").Return(true).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.token").Return("my-super-secret-auth-token").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.tls.insecure_skip_verify").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.retention_policy").Return(false).AnyTimes()
	if _, isGithubActions := os.LookupEnv("GITHUB_ACTIONS"); isGithubActions {
		// when running test suite in github actions, we run an influxdb service as a sidecar.
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("influxdb").AnyTimes()
	} else {
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("localhost").AnyTimes()
	}

	ae := web.AppEngine{
		Config: fakeConfig,
	}

	router := ae.Setup(logrus.WithField("test", suite.T().Name()))

	//test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", suite.Basepath+"/api/health", nil)
	router.ServeHTTP(w, req)

	//assert
	require.Equal(suite.T(), 200, w.Code)
	require.Equal(suite.T(), "{\"success\":true}", w.Body.String())
}

func (suite *ServerTestSuite) TestRegisterDevicesRoute() {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(suite.T())
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().SetDefault(gomock.Any(), gomock.Any()).AnyTimes()
	fakeConfig.EXPECT().UnmarshalKey(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	fakeConfig.EXPECT().GetString("web.database.location").Return(path.Join(parentPath, "scrutiny_test.db")).AnyTimes()
	fakeConfig.EXPECT().GetString("web.src.frontend.path").Return(parentPath).AnyTimes()
	fakeConfig.EXPECT().GetString("web.listen.basepath").Return(suite.Basepath).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.scheme").Return("http").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.port").Return("8086").AnyTimes()
	fakeConfig.EXPECT().IsSet("web.influxdb.token").Return(true).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.token").Return("my-super-secret-auth-token").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.tls.insecure_skip_verify").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.retention_policy").Return(false).AnyTimes()
	if _, isGithubActions := os.LookupEnv("GITHUB_ACTIONS"); isGithubActions {
		// when running test suite in github actions, we run an influxdb service as a sidecar.
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("influxdb").AnyTimes()
	} else {
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("localhost").AnyTimes()
	}

	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup(logrus.WithField("test", suite.T().Name()))
	file, err := os.Open("testdata/register-devices-req.json")
	require.NoError(suite.T(), err)

	//test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", suite.Basepath+"/api/devices/register", file)
	router.ServeHTTP(w, req)

	//assert
	require.Equal(suite.T(), 200, w.Code)
}

func (suite *ServerTestSuite) TestUploadDeviceMetricsRoute() {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(suite.T())
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().SetDefault(gomock.Any(), gomock.Any()).AnyTimes()
	fakeConfig.EXPECT().UnmarshalKey(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").AnyTimes().Return(parentPath)
	fakeConfig.EXPECT().GetString("web.listen.basepath").Return(suite.Basepath).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.scheme").Return("http").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.port").Return("8086").AnyTimes()
	fakeConfig.EXPECT().IsSet("web.influxdb.token").Return(true).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.token").Return("my-super-secret-auth-token").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetBool("user.metrics.repeat_notifications").Return(true).AnyTimes()
	fakeConfig.EXPECT().GetBool("user.collector.discard_sct_temp_history").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.tls.insecure_skip_verify").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.retention_policy").Return(false).AnyTimes()
	if _, isGithubActions := os.LookupEnv("GITHUB_ACTIONS"); isGithubActions {
		// when running test suite in github actions, we run an influxdb service as a sidecar.
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("influxdb").AnyTimes()
	} else {
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("localhost").AnyTimes()
	}
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.notify_level", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsNotifyLevelFail))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_filter_attributes", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusFilterAttributesAll))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_threshold", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusThresholdBoth))

	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup(logrus.WithField("test", suite.T().Name()))
	devicesfile, err := os.Open("testdata/register-devices-single-req.json")
	require.NoError(suite.T(), err)

	metricsfile := helperReadSmartDataFileFixTimestamp(suite.T(), "testdata/upload-device-metrics-req.json")

	//test
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", suite.Basepath+"/api/devices/register", devicesfile)
	router.ServeHTTP(wr, req)
	require.Equal(suite.T(), 200, wr.Code)

	mr := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", suite.Basepath+"/api/device/0x5000cca264eb01d7/smart", metricsfile)
	router.ServeHTTP(mr, req)
	require.Equal(suite.T(), 200, mr.Code)

	//assert
}

func (suite *ServerTestSuite) TestPopulateMultiple() {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(suite.T())
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().SetDefault(gomock.Any(), gomock.Any()).AnyTimes()
	fakeConfig.EXPECT().UnmarshalKey(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	//fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return("testdata/scrutiny_test.db")
	fakeConfig.EXPECT().GetStringSlice("notify.urls").Return([]string{}).AnyTimes()
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.notify_level", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsNotifyLevelFail))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_filter_attributes", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusFilterAttributesAll))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_threshold", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusThresholdBoth))
	fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").AnyTimes().Return(parentPath)
	fakeConfig.EXPECT().GetString("web.listen.basepath").Return(suite.Basepath).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.scheme").Return("http").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.port").Return("8086").AnyTimes()
	fakeConfig.EXPECT().IsSet("web.influxdb.token").Return(true).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.token").Return("my-super-secret-auth-token").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetBool("user.metrics.repeat_notifications").Return(true).AnyTimes()
	fakeConfig.EXPECT().GetBool("user.collector.discard_sct_temp_history").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.tls.insecure_skip_verify").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.retention_policy").Return(false).AnyTimes()
	if _, isGithubActions := os.LookupEnv("GITHUB_ACTIONS"); isGithubActions {
		// when running test suite in github actions, we run an influxdb service as a sidecar.
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("influxdb").AnyTimes()
	} else {
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("localhost").AnyTimes()
	}

	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup(logrus.WithField("test", suite.T().Name()))
	devicesfile, err := os.Open("testdata/register-devices-req.json")
	require.NoError(suite.T(), err)

	metricsfile := helperReadSmartDataFileFixTimestamp(suite.T(), "../models/testdata/smart-ata.json")
	failfile := helperReadSmartDataFileFixTimestamp(suite.T(), "../models/testdata/smart-fail2.json")
	nvmefile := helperReadSmartDataFileFixTimestamp(suite.T(), "../models/testdata/smart-nvme.json")
	scsifile := helperReadSmartDataFileFixTimestamp(suite.T(), "../models/testdata/smart-scsi.json")
	scsi2file := helperReadSmartDataFileFixTimestamp(suite.T(), "../models/testdata/smart-scsi2.json")

	//test
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", suite.Basepath+"/api/devices/register", devicesfile)
	router.ServeHTTP(wr, req)
	require.Equal(suite.T(), 200, wr.Code)

	mr := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", suite.Basepath+"/api/device/0x5000cca264eb01d7/smart", metricsfile)
	router.ServeHTTP(mr, req)
	require.Equal(suite.T(), 200, mr.Code)

	fr := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", suite.Basepath+"/api/device/0x5000cca264ec3183/smart", failfile)
	router.ServeHTTP(fr, req)
	require.Equal(suite.T(), 200, fr.Code)

	nr := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", suite.Basepath+"/api/device/0x5002538e40a22954/smart", nvmefile)
	router.ServeHTTP(nr, req)
	require.Equal(suite.T(), 200, nr.Code)

	sr := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", suite.Basepath+"/api/device/0x5000cca252c859cc/smart", scsifile)
	router.ServeHTTP(sr, req)
	require.Equal(suite.T(), 200, sr.Code)

	s2r := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", suite.Basepath+"/api/device/0x5000cca264ebc248/smart", scsi2file)
	router.ServeHTTP(s2r, req)
	require.Equal(suite.T(), 200, s2r.Code)

	//assert
}

//TODO: this test should use a recorded request/response playback.
//func TestSendTestNotificationRoute(t *testing.T) {
//	//setup
//	parentPath, _ := ioutil.TempDir("", "")
//	defer os.RemoveAll(parentPath)
//	mockCtrl := gomock.NewController(t)
//	fakeConfig := mock_config.NewMockInterface(mockCtrl)
//	fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return(path.Join(parentPath, "scrutiny_test.db"))
//	fakeConfig.EXPECT().GetString("web.src.frontend.path").AnyTimes().Return(parentPath)
//	fakeConfig.EXPECT().GetStringSlice("notify.urls").AnyTimes().Return([]string{"https://scrutiny.requestcatcher.com/test"})
//	ae := web.AppEngine{
//		Config: fakeConfig,
//	}
//	router := ae.Setup(logrus.New())
//
//	//test
//	wr := httptest.NewRecorder()
//	req, _ := http.NewRequest("POST", "/api/health/notify", strings.NewReader("{}"))
//	router.ServeHTTP(wr, req)
//
//	//assert
//	require.Equal(t, 200, wr.Code)
//}

func (suite *ServerTestSuite) TestSendTestNotificationRoute_WebhookFailure() {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(suite.T())
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().SetDefault(gomock.Any(), gomock.Any()).AnyTimes()
	fakeConfig.EXPECT().UnmarshalKey(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").AnyTimes().Return(parentPath)
	fakeConfig.EXPECT().GetString("web.listen.basepath").Return(suite.Basepath).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.scheme").Return("http").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.port").Return("8086").AnyTimes()
	fakeConfig.EXPECT().IsSet("web.influxdb.token").Return(true).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.token").Return("my-super-secret-auth-token").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.tls.insecure_skip_verify").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.retention_policy").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetStringSlice("notify.urls").AnyTimes().Return([]string{"https://unroutable.domain.example.asdfghj"})
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.notify_level", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsNotifyLevelFail))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_filter_attributes", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusFilterAttributesAll))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_threshold", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusThresholdBoth))

	if _, isGithubActions := os.LookupEnv("GITHUB_ACTIONS"); isGithubActions {
		// when running test suite in github actions, we run an influxdb service as a sidecar.
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("influxdb").AnyTimes()
	} else {
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("localhost").AnyTimes()
	}

	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup(logrus.WithField("test", suite.T().Name()))

	//test
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", suite.Basepath+"/api/health/notify", strings.NewReader("{}"))
	router.ServeHTTP(wr, req)

	//assert
	require.Equal(suite.T(), 500, wr.Code)
}

func (suite *ServerTestSuite) TestSendTestNotificationRoute_ScriptFailure() {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(suite.T())
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().SetDefault(gomock.Any(), gomock.Any()).AnyTimes()
	fakeConfig.EXPECT().UnmarshalKey(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").AnyTimes().Return(parentPath)
	fakeConfig.EXPECT().GetString("web.listen.basepath").Return(suite.Basepath).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.scheme").Return("http").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.port").Return("8086").AnyTimes()
	fakeConfig.EXPECT().IsSet("web.influxdb.token").Return(true).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.token").Return("my-super-secret-auth-token").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.tls.insecure_skip_verify").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.retention_policy").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetStringSlice("notify.urls").AnyTimes().Return([]string{"script:///missing/path/on/disk"})
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.notify_level", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsNotifyLevelFail))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_filter_attributes", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusFilterAttributesAll))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_threshold", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusThresholdBoth))

	if _, isGithubActions := os.LookupEnv("GITHUB_ACTIONS"); isGithubActions {
		// when running test suite in github actions, we run an influxdb service as a sidecar.
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("influxdb").AnyTimes()
	} else {
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("localhost").AnyTimes()
	}

	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup(logrus.WithField("test", suite.T().Name()))

	//test
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", suite.Basepath+"/api/health/notify", strings.NewReader("{}"))
	router.ServeHTTP(wr, req)

	//assert
	require.Equal(suite.T(), 500, wr.Code)
}

func (suite *ServerTestSuite) TestSendTestNotificationRoute_ScriptSuccess() {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(suite.T())
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().SetDefault(gomock.Any(), gomock.Any()).AnyTimes()
	fakeConfig.EXPECT().UnmarshalKey(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").AnyTimes().Return(parentPath)
	fakeConfig.EXPECT().GetString("web.listen.basepath").Return(suite.Basepath).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.scheme").Return("http").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.port").Return("8086").AnyTimes()
	fakeConfig.EXPECT().IsSet("web.influxdb.token").Return(true).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.token").Return("my-super-secret-auth-token").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.tls.insecure_skip_verify").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.retention_policy").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetStringSlice("notify.urls").AnyTimes().Return([]string{"script:///usr/bin/env"})
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.notify_level", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsNotifyLevelFail))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_filter_attributes", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusFilterAttributesAll))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_threshold", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusThresholdBoth))

	if _, isGithubActions := os.LookupEnv("GITHUB_ACTIONS"); isGithubActions {
		// when running test suite in github actions, we run an influxdb service as a sidecar.
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("influxdb").AnyTimes()
	} else {
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("localhost").AnyTimes()
	}

	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup(logrus.WithField("test", suite.T().Name()))

	//test
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", suite.Basepath+"/api/health/notify", strings.NewReader("{}"))
	router.ServeHTTP(wr, req)

	//assert
	require.Equal(suite.T(), 200, wr.Code)
}

func (suite *ServerTestSuite) TestSendTestNotificationRoute_ShoutrrrFailure() {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(suite.T())
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().SetDefault(gomock.Any(), gomock.Any()).AnyTimes()
	fakeConfig.EXPECT().UnmarshalKey(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").AnyTimes().Return(parentPath)
	fakeConfig.EXPECT().GetString("web.listen.basepath").Return(suite.Basepath).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.scheme").Return("http").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.port").Return("8086").AnyTimes()
	fakeConfig.EXPECT().IsSet("web.influxdb.token").Return(true).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.token").Return("my-super-secret-auth-token").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.tls.insecure_skip_verify").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.retention_policy").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetStringSlice("notify.urls").AnyTimes().Return([]string{"discord://invalidtoken@channel"})
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.notify_level", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsNotifyLevelFail))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_filter_attributes", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusFilterAttributesAll))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_threshold", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusThresholdBoth))

	if _, isGithubActions := os.LookupEnv("GITHUB_ACTIONS"); isGithubActions {
		// when running test suite in github actions, we run an influxdb service as a sidecar.
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("influxdb").AnyTimes()
	} else {
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("localhost").AnyTimes()
	}
	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup(logrus.WithField("test", suite.T().Name()))

	//test
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", suite.Basepath+"/api/health/notify", strings.NewReader("{}"))
	router.ServeHTTP(wr, req)

	//assert
	require.Equal(suite.T(), 500, wr.Code)
}

func (suite *ServerTestSuite) TestGetDevicesSummaryRoute_Nvme() {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(suite.T())
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().SetDefault(gomock.Any(), gomock.Any()).AnyTimes()
	fakeConfig.EXPECT().UnmarshalKey(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").AnyTimes().Return(parentPath)
	fakeConfig.EXPECT().GetString("web.listen.basepath").Return(suite.Basepath).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.scheme").Return("http").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.port").Return("8086").AnyTimes()
	fakeConfig.EXPECT().IsSet("web.influxdb.token").Return(true).AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.token").Return("my-super-secret-auth-token").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.org").Return("scrutiny").AnyTimes()
	fakeConfig.EXPECT().GetString("web.influxdb.bucket").Return("metrics").AnyTimes()
	fakeConfig.EXPECT().GetBool("user.metrics.repeat_notifications").Return(true).AnyTimes()
	fakeConfig.EXPECT().GetBool("user.collector.discard_sct_temp_history").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.tls.insecure_skip_verify").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetBool("web.influxdb.retention_policy").Return(false).AnyTimes()
	fakeConfig.EXPECT().GetStringSlice("notify.urls").AnyTimes().Return([]string{})
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.notify_level", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsNotifyLevelFail))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_filter_attributes", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusFilterAttributesAll))
	fakeConfig.EXPECT().GetInt(fmt.Sprintf("%s.metrics.status_threshold", config.DB_USER_SETTINGS_SUBKEY)).AnyTimes().Return(int(pkg.MetricsStatusThresholdBoth))

	if _, isGithubActions := os.LookupEnv("GITHUB_ACTIONS"); isGithubActions {
		// when running test suite in github actions, we run an influxdb service as a sidecar.
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("influxdb").AnyTimes()
	} else {
		fakeConfig.EXPECT().GetString("web.influxdb.host").Return("localhost").AnyTimes()
	}

	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup(logrus.WithField("test", suite.T().Name()))
	devicesfile, err := os.Open("testdata/register-devices-req-2.json")
	require.NoError(suite.T(), err)

	metricsfile := helperReadSmartDataFileFixTimestamp(suite.T(), "../models/testdata/smart-nvme2.json")

	//test
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", suite.Basepath+"/api/devices/register", devicesfile)
	router.ServeHTTP(wr, req)
	require.Equal(suite.T(), 200, wr.Code)

	mr := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", suite.Basepath+"/api/device/a4c8e8ed-11a0-4c97-9bba-306440f1b944/smart", metricsfile)
	router.ServeHTTP(mr, req)
	require.Equal(suite.T(), 200, mr.Code)

	sr := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", suite.Basepath+"/api/summary", nil)
	router.ServeHTTP(sr, req)
	require.Equal(suite.T(), 200, sr.Code)
	var deviceSummary models.DeviceSummaryWrapper
	err = json.Unmarshal(sr.Body.Bytes(), &deviceSummary)
	require.NoError(suite.T(), err)

	//assert
	require.Equal(suite.T(), "a4c8e8ed-11a0-4c97-9bba-306440f1b944", deviceSummary.Data.Summary["a4c8e8ed-11a0-4c97-9bba-306440f1b944"].Device.WWN)
	require.Equal(suite.T(), pkg.DeviceStatusPassed, deviceSummary.Data.Summary["a4c8e8ed-11a0-4c97-9bba-306440f1b944"].Device.DeviceStatus)
}
