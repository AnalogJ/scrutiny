package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/analogj/go-util/utils"
	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"github.com/analogj/scrutiny/webapp/backend/pkg/thresholds"
	"github.com/containrrr/shoutrrr"
	shoutrrrTypes "github.com/containrrr/shoutrrr/pkg/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const NotifyFailureTypeEmailTest = "EmailTest"
const NotifyFailureTypeBothFailure = "SmartFailure" //SmartFailure always takes precedence when Scrutiny & Smart failed.
const NotifyFailureTypeSmartFailure = "SmartFailure"
const NotifyFailureTypeScrutinyFailure = "ScrutinyFailure"

// ShouldNotify check if the error Message should be filtered (level mismatch or filtered_attributes)
func ShouldNotify(device models.Device, smartAttrs measurements.Smart, statusThreshold pkg.MetricsStatusThreshold, statusFilterAttributes pkg.MetricsStatusFilterAttributes) bool {
	// 1. check if the device is healthy
	if device.DeviceStatus == pkg.DeviceStatusPassed {
		return false
	}

	//TODO: cannot check for warning notifyLevel yet.

	// setup constants for comparison
	var requiredDeviceStatus pkg.DeviceStatus
	var requiredAttrStatus pkg.AttributeStatus
	if statusThreshold == pkg.MetricsStatusThresholdBoth {
		// either scrutiny or smart failures should trigger an email
		requiredDeviceStatus = pkg.DeviceStatusSet(pkg.DeviceStatusFailedSmart, pkg.DeviceStatusFailedScrutiny)
		requiredAttrStatus = pkg.AttributeStatusSet(pkg.AttributeStatusFailedSmart, pkg.AttributeStatusFailedScrutiny)
	} else if statusThreshold == pkg.MetricsStatusThresholdSmart {
		//only smart failures
		requiredDeviceStatus = pkg.DeviceStatusFailedSmart
		requiredAttrStatus = pkg.AttributeStatusFailedSmart
	} else {
		requiredDeviceStatus = pkg.DeviceStatusFailedScrutiny
		requiredAttrStatus = pkg.AttributeStatusFailedScrutiny
	}

	// 2. check if the attributes that are failing should be filtered (non-critical)
	// 3. for any unfiltered attribute, store the failure reason (Smart or Scrutiny)
	if statusFilterAttributes == pkg.MetricsStatusFilterAttributesCritical {
		hasFailingCriticalAttr := false
		var statusFailingCriticalAttr pkg.AttributeStatus

		for attrId, attrData := range smartAttrs.Attributes {
			//find failing attribute
			if attrData.GetStatus() == pkg.AttributeStatusPassed {
				continue //skip all passing attributes
			}

			// merge the status's of all critical attributes
			statusFailingCriticalAttr = pkg.AttributeStatusSet(statusFailingCriticalAttr, attrData.GetStatus())

			//found a failing attribute, see if its critical
			if device.IsScsi() && thresholds.ScsiMetadata[attrId].Critical {
				hasFailingCriticalAttr = true
			} else if device.IsNvme() && thresholds.NmveMetadata[attrId].Critical {
				hasFailingCriticalAttr = true
			} else {
				//this is ATA
				attrIdInt, err := strconv.Atoi(attrId)
				if err != nil {
					continue
				}
				if thresholds.AtaMetadata[attrIdInt].Critical {
					hasFailingCriticalAttr = true
				}
			}

		}

		if !hasFailingCriticalAttr {
			//no critical attributes are failing, and notifyFilterAttributes == "critical"
			return false
		} else {
			// check if any of the critical attributes have a status that we're looking for
			return pkg.AttributeStatusHas(statusFailingCriticalAttr, requiredAttrStatus)
		}

	} else {
		// 2. SKIP - we are processing every attribute.
		// 3. check if the device failure level matches the wanted failure level.
		return pkg.DeviceStatusHas(device.DeviceStatus, requiredDeviceStatus)
	}
}

// TODO: include user label for device.
type Payload struct {
	HostId       string `json:"host_id,omitempty"` //host id (optional)
	DeviceType   string `json:"device_type"`       //ATA/SCSI/NVMe
	DeviceName   string `json:"device_name"`       //dev/sda
	DeviceSerial string `json:"device_serial"`     //WDDJ324KSO
	Test         bool   `json:"test"`              // false

	//private, populated during init (marked as Public for JSON serialization)
	Date        string `json:"date"`         //populated by Send function.
	FailureType string `json:"failure_type"` //EmailTest, BothFail, SmartFail, ScrutinyFail
	Subject     string `json:"subject"`
	Message     string `json:"message"`
}

func NewPayload(device models.Device, test bool, currentTime ...time.Time) Payload {
	payload := Payload{
		HostId:       strings.TrimSpace(device.HostId),
		DeviceType:   device.DeviceType,
		DeviceName:   device.DeviceName,
		DeviceSerial: device.SerialNumber,
		Test:         test,
	}

	//validate that the Payload is populated
	var sendDate time.Time
	if currentTime != nil && len(currentTime) > 0 {
		sendDate = currentTime[0]
	} else {
		sendDate = time.Now()
	}

	payload.Date = sendDate.Format(time.RFC3339)
	payload.FailureType = payload.GenerateFailureType(device.DeviceStatus)
	payload.Subject = payload.GenerateSubject()
	payload.Message = payload.GenerateMessage()
	return payload
}

func (p *Payload) GenerateFailureType(deviceStatus pkg.DeviceStatus) string {
	//generate a failure type, given Test and DeviceStatus
	if p.Test {
		return NotifyFailureTypeEmailTest // must be an email test if "Test" is true
	}
	if pkg.DeviceStatusHas(deviceStatus, pkg.DeviceStatusFailedSmart) && pkg.DeviceStatusHas(deviceStatus, pkg.DeviceStatusFailedScrutiny) {
		return NotifyFailureTypeBothFailure //both failed
	} else if pkg.DeviceStatusHas(deviceStatus, pkg.DeviceStatusFailedSmart) {
		return NotifyFailureTypeSmartFailure //only SMART failed
	} else {
		return NotifyFailureTypeScrutinyFailure //only Scrutiny failed
	}
}

func (p *Payload) GenerateSubject() string {
	//generate a detailed failure message
	var subject string
	if len(p.HostId) > 0 {
		subject = fmt.Sprintf("Scrutiny SMART error (%s) detected on [host]device: [%s]%s", p.FailureType, p.HostId, p.DeviceName)
	} else {
		subject = fmt.Sprintf("Scrutiny SMART error (%s) detected on device: %s", p.FailureType, p.DeviceName)
	}
	return subject
}

func (p *Payload) GenerateMessage() string {
	//generate a detailed failure message

	messageParts := []string{}

	messageParts = append(messageParts, fmt.Sprintf("Scrutiny SMART error notification for device: %s", p.DeviceName))
	if len(p.HostId) > 0 {
		messageParts = append(messageParts, fmt.Sprintf("Host Id: %s", p.HostId))
	}

	messageParts = append(messageParts,
		fmt.Sprintf("Failure Type: %s", p.FailureType),
		fmt.Sprintf("Device Name: %s", p.DeviceName),
		fmt.Sprintf("Device Serial: %s", p.DeviceSerial),
		fmt.Sprintf("Device Type: %s", p.DeviceType),
		"",
		fmt.Sprintf("Date: %s", p.Date),
	)

	if p.Test {
		messageParts = append([]string{"TEST NOTIFICATION:"}, messageParts...)
	}

	return strings.Join(messageParts, "\n")
}

func New(logger logrus.FieldLogger, appconfig config.Interface, device models.Device, test bool) Notify {
	return Notify{
		Logger:  logger,
		Config:  appconfig,
		Payload: NewPayload(device, test),
	}
}

type Notify struct {
	Logger  logrus.FieldLogger
	Config  config.Interface
	Payload Payload
}

func (n *Notify) Send() error {

	//retrieve list of notification endpoints from config file
	configUrls := n.Config.GetStringSlice("notify.urls")
	n.Logger.Debugf("Configured notification services: %v", configUrls)

	if len(configUrls) == 0 {
		n.Logger.Infof("No notification endpoints configured. Skipping failure notification.")
		return nil
	}

	//remove http:// https:// and script:// prefixed urls
	notifyWebhooks := []string{}
	notifyScripts := []string{}
	notifyShoutrrr := []string{}

	for ndx, _ := range configUrls {
		if strings.HasPrefix(configUrls[ndx], "https://") || strings.HasPrefix(configUrls[ndx], "http://") {
			notifyWebhooks = append(notifyWebhooks, configUrls[ndx])
		} else if strings.HasPrefix(configUrls[ndx], "script://") {
			notifyScripts = append(notifyScripts, configUrls[ndx])
		} else {
			notifyShoutrrr = append(notifyShoutrrr, configUrls[ndx])
		}
	}

	n.Logger.Debugf("Configured scripts: %v", notifyScripts)
	n.Logger.Debugf("Configured webhooks: %v", notifyWebhooks)
	n.Logger.Debugf("Configured shoutrrr: %v", notifyShoutrrr)

	//run all scripts, webhooks and shoutrr commands in parallel
	//var wg sync.WaitGroup
	var eg errgroup.Group

	for _, url := range notifyWebhooks {
		// execute collection in parallel go-routines
		_url := url
		eg.Go(func() error { return n.SendWebhookNotification(_url) })
	}
	for _, url := range notifyScripts {
		// execute collection in parallel go-routines
		_url := url
		eg.Go(func() error { return n.SendScriptNotification(_url) })
	}
	for _, url := range notifyShoutrrr {
		// execute collection in parallel go-routines
		_url := url
		eg.Go(func() error { return n.SendShoutrrrNotification(_url) })
	}

	//and wait for completion, error or timeout.
	n.Logger.Debugf("Main: waiting for notifications to complete.")

	if err := eg.Wait(); err == nil {
		n.Logger.Info("Successfully sent notifications. Check logs for more information.")
		return nil
	} else {
		n.Logger.Error("One or more notifications failed to send successfully. See logs for more information.")
		return err
	}
	////wg.Wait()
	//if waitTimeout(&wg, time.Minute) { //wait for 1 minute
	//	fmt.Println("Timed out while sending notifications")
	//} else {
	//}
	//return nil
}

func (n *Notify) SendWebhookNotification(webhookUrl string) error {
	n.Logger.Infof("Sending Webhook to %s", webhookUrl)
	requestBody, err := json.Marshal(n.Payload)
	if err != nil {
		n.Logger.Errorf("An error occurred while sending Webhook to %s: %v", webhookUrl, err)
		return err
	}

	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		n.Logger.Errorf("An error occurred while sending Webhook to %s: %v", webhookUrl, err)
		return err
	}
	defer resp.Body.Close()
	//we don't care about resp body content, but maybe we should log it?
	return nil
}

func (n *Notify) SendScriptNotification(scriptUrl string) error {
	//check if the script exists.
	scriptPath := strings.TrimPrefix(scriptUrl, "script://")
	n.Logger.Infof("Executing Script %s", scriptPath)

	if !utils.FileExists(scriptPath) {
		n.Logger.Errorf("Script does not exist: %s", scriptPath)
		return errors.New(fmt.Sprintf("custom script path does not exist: %s", scriptPath))
	}

	copyEnv := os.Environ()
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_SUBJECT=%s", n.Payload.Subject))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_DATE=%s", n.Payload.Date))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_FAILURE_TYPE=%s", n.Payload.FailureType))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_DEVICE_NAME=%s", n.Payload.DeviceName))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_DEVICE_TYPE=%s", n.Payload.DeviceType))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_DEVICE_SERIAL=%s", n.Payload.DeviceSerial))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_MESSAGE=%s", n.Payload.Message))
	if len(n.Payload.HostId) > 0 {
		copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_HOST_ID=%s", n.Payload.HostId))
	}
	err := utils.CmdExec(scriptPath, []string{}, "", copyEnv, "")
	if err != nil {
		n.Logger.Errorf("An error occurred while executing script %s: %v", scriptPath, err)
		return err
	}
	return nil
}

func (n *Notify) SendShoutrrrNotification(shoutrrrUrl string) error {

	fmt.Printf("Sending Notifications to %v", shoutrrrUrl)
	n.Logger.Infof("Sending notifications to %v", shoutrrrUrl)

	sender, err := shoutrrr.CreateSender(shoutrrrUrl)
	if err != nil {
		n.Logger.Errorf("An error occurred while sending notifications %v: %v", shoutrrrUrl, err)
		return err
	}

	//sender.SetLogger(n.Logger.)
	serviceName, params, err := n.GenShoutrrrNotificationParams(shoutrrrUrl)
	n.Logger.Debugf("notification data for %s: (%s)\n%v", serviceName, shoutrrrUrl, params)

	if err != nil {
		n.Logger.Errorf("An error occurred  occurred while generating notification payload for %s:\n %v", serviceName, shoutrrrUrl, err)
		return err
	}

	errs := sender.Send(n.Payload.Message, params)
	if len(errs) > 0 {
		var errstrings []string

		for _, err := range errs {
			if err == nil || err.Error() == "" {
				continue
			}
			errstrings = append(errstrings, err.Error())
		}
		//sometimes there are empty errs, we're going to skip them.
		if len(errstrings) == 0 {
			return nil
		} else {
			n.Logger.Errorf("One or more errors occurred while sending notifications for %s:", shoutrrrUrl)
			n.Logger.Error(errs)
			return errors.New(strings.Join(errstrings, "\n"))
		}
	}
	return nil
}

func (n *Notify) GenShoutrrrNotificationParams(shoutrrrUrl string) (string, *shoutrrrTypes.Params, error) {
	serviceURL, err := url.Parse(shoutrrrUrl)
	if err != nil {
		return "", nil, err
	}

	serviceName := serviceURL.Scheme
	params := &shoutrrrTypes.Params{}

	logoUrl := "https://raw.githubusercontent.com/AnalogJ/scrutiny/master/webapp/frontend/src/ms-icon-144x144.png"
	subject := n.Payload.Subject
	switch serviceName {
	// no params supported for these services
	case "hangouts", "mattermost", "teams", "rocketchat":
		break
	case "discord":
		(*params)["title"] = subject
	case "gotify":
		(*params)["title"] = subject
	case "ifttt":
		(*params)["title"] = subject
	case "join":
		(*params)["title"] = subject
		(*params)["icon"] = logoUrl
	case "opsgenie":
		(*params)["title"] = subject
	case "pushbullet":
		(*params)["title"] = subject
	case "pushover":
		(*params)["title"] = subject
	case "slack":
		(*params)["title"] = subject
	case "smtp":
		(*params)["subject"] = subject
	case "standard":
		(*params)["subject"] = subject
	case "telegram":
		(*params)["title"] = subject
	case "zulip":
		(*params)["topic"] = subject
	}

	return serviceName, params, nil
}
