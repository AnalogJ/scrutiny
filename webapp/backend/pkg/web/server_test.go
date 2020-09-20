package web_test

import (
	"encoding/json"
	mock_config "github.com/analogj/scrutiny/webapp/backend/pkg/config/mock"
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/analogj/scrutiny/webapp/backend/pkg/web"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
)

func TestHealthRoute(t *testing.T) {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("web.database.location").Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").Return(parentPath)

	ae := web.AppEngine{
		Config: fakeConfig,
	}

	router := ae.Setup()

	//test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/health", nil)
	router.ServeHTTP(w, req)

	//assert
	require.Equal(t, 200, w.Code)
	require.Equal(t, "{\"success\":true}", w.Body.String())
}

func TestRegisterDevicesRoute(t *testing.T) {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("web.database.location").Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").Return(parentPath)
	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup()
	file, err := os.Open("testdata/register-devices-req.json")
	require.NoError(t, err)

	//test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/devices/register", file)
	router.ServeHTTP(w, req)

	//assert
	require.Equal(t, 200, w.Code)
}

func TestUploadDeviceMetricsRoute(t *testing.T) {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").AnyTimes().Return(parentPath)
	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup()
	devicesfile, err := os.Open("testdata/register-devices-single-req.json")
	require.NoError(t, err)

	metricsfile, err := os.Open("testdata/upload-device-metrics-req.json")
	require.NoError(t, err)

	//test
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/devices/register", devicesfile)
	router.ServeHTTP(wr, req)
	require.Equal(t, 200, wr.Code)

	mr := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/device/0x5000cca264eb01d7/smart", metricsfile)
	router.ServeHTTP(mr, req)
	require.Equal(t, 200, mr.Code)

	//assert
}

func TestPopulateMultiple(t *testing.T) {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	//fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return("testdata/scrutiny_test.db")
	fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").AnyTimes().Return(parentPath)
	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup()
	devicesfile, err := os.Open("testdata/register-devices-req.json")
	require.NoError(t, err)

	metricsfile, err := os.Open("../models/testdata/smart-ata.json")
	require.NoError(t, err)
	failfile, err := os.Open("../models/testdata/smart-fail2.json")
	require.NoError(t, err)
	nvmefile, err := os.Open("../models/testdata/smart-nvme.json")
	require.NoError(t, err)
	scsifile, err := os.Open("../models/testdata/smart-scsi.json")
	require.NoError(t, err)
	scsi2file, err := os.Open("../models/testdata/smart-scsi2.json")
	require.NoError(t, err)

	//test
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/devices/register", devicesfile)
	router.ServeHTTP(wr, req)
	require.Equal(t, 200, wr.Code)

	mr := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/device/0x5000cca264eb01d7/smart", metricsfile)
	router.ServeHTTP(mr, req)
	require.Equal(t, 200, mr.Code)

	fr := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/device/0x5000cca264ec3183/smart", failfile)
	router.ServeHTTP(fr, req)
	require.Equal(t, 200, fr.Code)

	nr := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/device/0x5002538e40a22954/smart", nvmefile)
	router.ServeHTTP(nr, req)
	require.Equal(t, 200, nr.Code)

	sr := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/device/0x5000cca252c859cc/smart", scsifile)
	router.ServeHTTP(sr, req)
	require.Equal(t, 200, sr.Code)

	s2r := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/device/0x5000cca264ebc248/smart", scsi2file)
	router.ServeHTTP(s2r, req)
	require.Equal(t, 200, s2r.Code)

	//assert
}

func TestSendTestNotificationRoute(t *testing.T) {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").AnyTimes().Return(parentPath)
	fakeConfig.EXPECT().GetStringSlice("notify.urls").AnyTimes().Return([]string{"https://scrutiny.requestcatcher.com/test"})
	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup()

	//test
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/health/notify", strings.NewReader("{}"))
	router.ServeHTTP(wr, req)

	//assert
	require.Equal(t, 200, wr.Code)
}

func TestGetDevicesSummaryRoute_Nvme(t *testing.T) {
	//setup
	parentPath, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("web.database.location").AnyTimes().Return(path.Join(parentPath, "scrutiny_test.db"))
	fakeConfig.EXPECT().GetString("web.src.frontend.path").AnyTimes().Return(parentPath)
	ae := web.AppEngine{
		Config: fakeConfig,
	}
	router := ae.Setup()
	devicesfile, err := os.Open("testdata/register-devices-req-2.json")
	require.NoError(t, err)

	metricsfile, err := os.Open("../models/testdata/smart-nvme2.json")
	require.NoError(t, err)

	//test
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/devices/register", devicesfile)
	router.ServeHTTP(wr, req)
	require.Equal(t, 200, wr.Code)

	mr := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/device/a4c8e8ed-11a0-4c97-9bba-306440f1b944/smart", metricsfile)
	router.ServeHTTP(mr, req)
	require.Equal(t, 200, mr.Code)

	sr := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/summary", nil)
	router.ServeHTTP(sr, req)
	require.Equal(t, 200, sr.Code)
	var device dbModels.DeviceWrapper
	json.Unmarshal(sr.Body.Bytes(), &device)

	//assert
	require.Equal(t, "a4c8e8ed-11a0-4c97-9bba-306440f1b944", device.Data[0].WWN)
	require.Equal(t, "passed", device.Data[0].SmartResults[0].SmartStatus)
}
