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

	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
)

func main() {

	//webapp/backend/pkg/web/testdata/register-devices-req.json
	devices := "webapp/backend/pkg/web/testdata/register-devices-req.json"

	smartData := map[string][]string{
		"0x5000cca264eb01d7": {"webapp/backend/pkg/models/testdata/smart-ata.json", "webapp/backend/pkg/models/testdata/smart-ata-date.json", "webapp/backend/pkg/models/testdata/smart-ata-date2.json"},
		"0x5000cca264ec3183": {"webapp/backend/pkg/models/testdata/smart-fail2.json"},
		"0x5002538e40a22954": {"webapp/backend/pkg/models/testdata/smart-nvme.json"},
		"0x5000cca252c859cc": {"webapp/backend/pkg/models/testdata/smart-scsi.json"},
		"0x5000cca264ebc248": {"webapp/backend/pkg/models/testdata/smart-scsi2.json"},
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

}

func SendPostRequest(url string, file io.Reader) ([]byte, error) {
	response, err := http.Post(url, "application/json", file)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	log.Printf("%v\n", response.Status)

	return io.ReadAll(response.Body)
}

// InfluxDB will throw an error/ignore any submitted data with a timestamp older than the
// retention period. Lets fix this by opening test files, modifying the timestamp and returning an io.Reader
func readSmartDataFileFixTimestamp(daysToSubtract int, smartDataFilepath string) (io.Reader, error) {
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
	smartData.LocalTime.TimeT = time.Now().Add(daysToSubtractInHours * time.Hour).Unix()
	updatedSmartDataBytes, err := json.Marshal(smartData)

	return bytes.NewReader(updatedSmartDataBytes), nil
}
