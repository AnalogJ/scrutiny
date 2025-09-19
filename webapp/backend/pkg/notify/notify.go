package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/analogj/go-util/utils"
	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"github.com/analogj/scrutiny/webapp/backend/pkg/thresholds"
	"github.com/containrrr/shoutrrr"
	shoutrrrTypes "github.com/containrrr/shoutrrr/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const NotifyFailureTypeEmailTest = "EmailTest"
const NotifyFailureTypeBothFailure = "SmartFailure" //SmartFailure always takes precedence when Scrutiny & Smart failed.
const NotifyFailureTypeSmartFailure = "SmartFailure"
const NotifyFailureTypeScrutinyFailure = "ScrutinyFailure"

// ShouldNotify check if the error Message should be filtered (level mismatch or filtered_attributes)
func ShouldNotify(logger logrus.FieldLogger, device models.Device, smartAttrs measurements.Smart, statusThreshold pkg.MetricsStatusThreshold, statusFilterAttributes pkg.MetricsStatusFilterAttributes, repeatNotifications bool, c *gin.Context, deviceRepo database.DeviceRepo) bool {
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

	// This is the only case where individual attributes need not be considered
	if statusFilterAttributes == pkg.MetricsStatusFilterAttributesAll && repeatNotifications {
		return pkg.DeviceStatusHas(device.DeviceStatus, requiredDeviceStatus)
	}

	var failingAttributes []string
	// Loop through the attributes to find the failing ones
	for attrId, attrData := range smartAttrs.Attributes {
		var status pkg.AttributeStatus = attrData.GetStatus()
		// Skip over passing attributes
		if status == pkg.AttributeStatusPassed {
			continue
		}

		// If the user only wants to consider critical attributes, we have to check
		// if the not-passing attribute is critical or not
		if statusFilterAttributes == pkg.MetricsStatusFilterAttributesCritical {
			critical := false
			if device.IsScsi() {
				critical = thresholds.ScsiMetadata[attrId].Critical
			} else if device.IsNvme() {
				critical = thresholds.NmveMetadata[attrId].Critical
			} else {
				//this is ATA
				attrIdInt, err := strconv.Atoi(attrId)
				if err != nil {
					continue
				}
				critical = thresholds.AtaMetadata[attrIdInt].Critical
			}
			// Skip non-critical, non-passing attributes when this setting is on
			if !critical {
				continue
			}
		}

		// Record any attribute that doesn't get skipped by the above two checks
		failingAttributes = append(failingAttributes, attrId)
	}

	// If the user doesn't want repeated notifications when the failing value doesn't change, we need to get the last value from the db
	var lastPoints []measurements.Smart
	var err error
	if !repeatNotifications {
		lastPoints, err = deviceRepo.GetSmartAttributeHistory(c, c.Param("wwn"), database.DURATION_KEY_FOREVER, 1, 1, failingAttributes)
		if err == nil || len(lastPoints) < 1 {
			logger.Warningln("Could not get the most recent data points from the database. This is expected to happen only if this is the very first submission of data for the device.")
		}
	}
	for _, attrId := range failingAttributes {
		attrStatus := smartAttrs.Attributes[attrId].GetStatus()
		if pkg.AttributeStatusHas(attrStatus, requiredAttrStatus) {
			if repeatNotifications {
				return true
			}
			// This is checked again here to avoid repeating the entire for loop in the check above.
			// Probably unnoticeably worse performance, but cleaner code.
			if err != nil || len(lastPoints) < 1 || lastPoints[0].Attributes[attrId].GetTransformedValue() != smartAttrs.Attributes[attrId].GetTransformedValue() {
				return true
			}
		}
	}
	return false
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

	for ndx := range configUrls {
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
	case "ntfy":
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
		if len(subject) > 60 {
			subject = subject[:60]
		}
		urlTopic := serviceURL.Query()["force_topic"]
		if urlTopic != "" {
			subject = urlTopic
		}
		(*params)["topic"] = subject
	}

	return serviceName, params, nil
}
