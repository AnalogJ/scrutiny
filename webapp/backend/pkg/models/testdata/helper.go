package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
)

func main() {
	if err := waitForServer("http://localhost:8080/api/health", time.Minute); err != nil {
		log.Fatalf("ERROR %v", err)
	}

	//webapp/backend/pkg/web/testdata/register-devices-req.json
	devices := "webapp/backend/pkg/web/testdata/register-devices-req.json"

	smartData := map[string][]string{
		"3ea22b35-682b-49fb-a655-abffed108e48": {"webapp/backend/pkg/models/testdata/smart-ata.json", "webapp/backend/pkg/models/testdata/smart-ata-date.json", "webapp/backend/pkg/models/testdata/smart-ata-date2.json"},
		"42caca8a-9b95-5c75-b059-305771a2a193": {"webapp/backend/pkg/models/testdata/smart-fail2.json"},
		"ecfaaf20-d1f6-558b-b33a-3e8db19a6c2c": {"webapp/backend/pkg/models/testdata/smart-nvme.json"},
		"d8796fe7-2422-520c-8991-e970993dad3e": {"webapp/backend/pkg/models/testdata/smart-scsi.json"},
		"00328b73-9f8a-53ad-8f20-8d0b1be00f47": {"webapp/backend/pkg/models/testdata/smart-scsi2.json"},
	}

	// send a post request to register devices
	file, err := os.Open(devices)
	if err != nil {
		log.Fatalf("ERROR %v", err)
	}
	defer file.Close()
	_, err = SendPostRequest("http://localhost:8080/api/devices/register", file)
	if err != nil {
		log.Fatalf("ERROR %v", err)
	}
	//

	for diskId, smartDataFileNames := range smartData {
		for _, smartDataFileName := range smartDataFileNames {
			for daysToSubtract := 0; daysToSubtract <= 30; daysToSubtract++ { //add 4 weeks worth of data
				smartDataReader, err := readSmartDataFileFixTimestamp(daysToSubtract, smartDataFileName)
				if err != nil {
					log.Fatalf("ERROR %v", err)
				}

				_, err = SendPostRequest(fmt.Sprintf("http://localhost:8080/api/device/%s/smart", diskId), smartDataReader)
				if err != nil {
					log.Fatalf("ERROR %v", err)
				}
			}

		}

	}

	// Add a second daily reading with a visible temperature and power-on-hours
	// trend to one drive, so its detail charts contain richer sample data.
	for daysToSubtract := 0; daysToSubtract <= 30; daysToSubtract++ {
		smartDataReader, err := readSmartDataFileWithTrend(daysToSubtract, "webapp/backend/pkg/models/testdata/smart-ata.json")
		if err != nil {
			log.Fatalf("ERROR %v", err)
		}

		_, err = SendPostRequest("http://localhost:8080/api/device/3ea22b35-682b-49fb-a655-abffed108e48/smart", smartDataReader)
		if err != nil {
			log.Fatalf("ERROR %v", err)
		}
	}

	if err := waitForDummyData("http://localhost:8080/api/summary", "3ea22b35-682b-49fb-a655-abffed108e48", 10*time.Second); err != nil {
		log.Fatalf("ERROR %v", err)
	}
	log.Println("Dummy data loaded")

}

func waitForServer(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		response, err := http.Get(url)
		if err == nil {
			response.Body.Close()
			if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices {
				return nil
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timed out waiting for %s", url)
}

func waitForDummyData(url string, scrutinyUUID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		response, err := http.Get(url)
		if err == nil {
			var summary models.DeviceSummaryWrapper
			decodeErr := json.NewDecoder(response.Body).Decode(&summary)
			response.Body.Close()
			if decodeErr == nil && response.StatusCode == http.StatusOK {
				deviceSummary := summary.Data.Summary[scrutinyUUID]
				if deviceSummary != nil && deviceSummary.SmartResults != nil {
					return nil
				}
			}
		}

		time.Sleep(250 * time.Millisecond)
	}

	return fmt.Errorf("dummy SMART data was submitted but did not appear in %s", url)
}

func SendPostRequest(url string, file io.Reader) ([]byte, error) {
	response, err := http.Post(url, "application/json", file)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("%v\n", response.Status)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return responseBody, fmt.Errorf("POST %s returned %s: %s", url, response.Status, responseBody)
	}

	return responseBody, nil
}

// InfluxDB will throw an error/ignore any submitted data with a timestamp older than the
// retention period. Lets fix this by opening test files, modifying the timestamp and returning an io.Reader
func readSmartDataFileFixTimestamp(daysToSubtract int, smartDataFilepath string) (io.Reader, error) {
	return readSmartDataFile(daysToSubtract, smartDataFilepath, false)
}

func readSmartDataFileWithTrend(daysToSubtract int, smartDataFilepath string) (io.Reader, error) {
	return readSmartDataFile(daysToSubtract, smartDataFilepath, true)
}

func readSmartDataFile(daysToSubtract int, smartDataFilepath string, addTrend bool) (io.Reader, error) {
	metricsfile, err := os.Open(smartDataFilepath)
	if err != nil {
		return nil, err
	}

	metricsFileData, err := io.ReadAll(metricsfile)
	if err != nil {
		return nil, err
	}
	//unmarshal because we need to change the timestamp
	var smartData collector.SmartInfo
	err = json.Unmarshal(metricsFileData, &smartData)
	if err != nil {
		return nil, err
	}

	daysToSubtractInHours := time.Duration(-1 * 24 * daysToSubtract)
	timestamp := time.Now().Add(daysToSubtractInHours * time.Hour)
	if addTrend {
		timestamp = timestamp.Add(-12 * time.Hour)
		smartData.Temperature.Current = 30 + int64((30-daysToSubtract)%7)
		smartData.PowerOnTime.Hours += int64((30 - daysToSubtract) * 24)
	}
	smartData.LocalTime.TimeT = timestamp.Unix()
	updatedSmartDataBytes, err := json.Marshal(smartData)

	return bytes.NewReader(updatedSmartDataBytes), nil
}
