package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/analogj/go-util/utils"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/containrrr/shoutrrr"
	shoutrrrTypes "github.com/containrrr/shoutrrr/pkg/types"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const NotifyFailureTypeEmailTest = "EmailTest"
const NotifyFailureTypeSmartPrefail = "SmartPreFailure"
const NotifyFailureTypeSmartFailure = "SmartFailure"
const NotifyFailureTypeSmartErrorLog = "SmartErrorLog"
const NotifyFailureTypeSmartSelfTest = "SmartSelfTestLog"

// TODO: include host and/or user label for device.
type Payload struct {
	Date         string `json:"date"`          //populated by Send function.
	FailureType  string `json:"failure_type"`  //EmailTest, SmartFail, ScrutinyFail
	DeviceType   string `json:"device_type"`   //ATA/SCSI/NVMe
	DeviceName   string `json:"device_string"` //dev/sda
	DeviceSerial string `json:"device"`        //WDDJ324KSO
	Test         bool   `json:"-"`             // false
}

func (p *Payload) GenerateMessage() string {
	//generate a detailed failure message
	return fmt.Sprintf("Scrutiny SMART error (%s) detected on device: %s", p.FailureType, p.DeviceName)
}

func (p *Payload) GenerateSubject() string {
	//generate a detailed failure message
	message := fmt.Sprintf(
		`Scrutiny SMART error notification for device: %s
Failure Type: %s
Device Name: %s
Device Serial: %s
Device Type: %s

Date: %s`, p.DeviceName, p.FailureType, p.DeviceName, p.DeviceSerial, p.DeviceType, p.Date)

	if p.Test {
		message = "TEST NOTIFICATION:\n" + message
	}

	return message
}

type Notify struct {
	Logger  logrus.FieldLogger
	Config  config.Interface
	Payload Payload
}

func (n *Notify) Send() error {
	//validate that the Payload is populated
	sendDate := time.Now()
	n.Payload.Date = sendDate.Format(time.RFC3339)

	//retrieve list of notification endpoints from config file
	configUrls := n.Config.GetStringSlice("notify.urls")
	n.Logger.Debugf("Configured notification services: %v", configUrls)

	//remove http:// https:// and script:// prefixed urls
	notifyWebhooks := []string{}
	notifyScripts := []string{}
	notifyShoutrrr := []string{}

	for _, url := range configUrls {
		if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
			notifyWebhooks = append(notifyWebhooks, url)
		} else if strings.HasPrefix(url, "script://") {
			notifyScripts = append(notifyScripts, url)
		} else {
			notifyShoutrrr = append(notifyShoutrrr, url)
		}
	}

	n.Logger.Debugf("Configured scripts: %v", notifyScripts)
	n.Logger.Debugf("Configured webhooks: %v", notifyWebhooks)
	n.Logger.Debugf("Configured shoutrrr: %v", notifyShoutrrr)

	//run all scripts, webhooks and shoutrr commands in parallel
	var wg sync.WaitGroup

	for _, notifyWebhook := range notifyWebhooks {
		// execute collection in parallel go-routines
		wg.Add(1)
		go n.SendWebhookNotification(&wg, notifyWebhook)
	}
	for _, notifyScript := range notifyScripts {
		// execute collection in parallel go-routines
		wg.Add(1)
		go n.SendScriptNotification(&wg, notifyScript)
	}
	for _, shoutrrrUrl := range notifyShoutrrr {
		wg.Add(1)
		go n.SendShoutrrrNotification(&wg, shoutrrrUrl)
	}

	//and wait for completion, error or timeout.
	n.Logger.Debugf("Main: waiting for notifications to complete.")
	//wg.Wait()
	if waitTimeout(&wg, time.Minute) { //wait for 1 minute
		fmt.Println("Timed out while sending notifications")
	} else {
		fmt.Println("Sent notifications. Check logs for more information.")
	}
	return nil
}

func (n *Notify) SendWebhookNotification(wg *sync.WaitGroup, webhookUrl string) {
	defer wg.Done()
	n.Logger.Infof("Sending Webhook to %s", webhookUrl)
	requestBody, err := json.Marshal(n.Payload)
	if err != nil {
		n.Logger.Errorf("An error occurred while sending Webhook to %s: %v", webhookUrl, err)
		return
	}

	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		n.Logger.Errorf("An error occurred while sending Webhook to %s: %v", webhookUrl, err)
		return
	}
	defer resp.Body.Close()
	//we don't care about resp body content, but maybe we should log it?
}

func (n *Notify) SendScriptNotification(wg *sync.WaitGroup, scriptUrl string) {
	defer wg.Done()

	//check if the script exists.
	scriptPath := strings.TrimPrefix(scriptUrl, "script://")
	n.Logger.Infof("Executing Script %s", scriptPath)

	if !utils.FileExists(scriptPath) {
		n.Logger.Errorf("Script does not exist: %s", scriptPath)
		return
	}

	copyEnv := os.Environ()
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_SUBJECT=%s", n.Payload.GenerateSubject()))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_DATE=%s", n.Payload.Date))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_FAILURE_TYPE=%s", n.Payload.FailureType))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_DEVICE_NAME=%s", n.Payload.DeviceName))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_DEVICE_TYPE=%s", n.Payload.DeviceType))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_DEVICE_SERIAL=%s", n.Payload.DeviceSerial))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_MESSAGE=%s", n.Payload.GenerateMessage()))
	err := utils.CmdExec(scriptPath, []string{}, "", copyEnv, "")
	if err != nil {
		n.Logger.Errorf("An error occurred while executing script %s: %v", scriptPath, err)
	}
	return
}

func (n *Notify) SendShoutrrrNotification(wg *sync.WaitGroup, shoutrrrUrl string) {

	fmt.Printf("Sending Notifications to %v", shoutrrrUrl)
	n.Logger.Infof("Sending notifications to %v", shoutrrrUrl)

	defer wg.Done()
	sender, err := shoutrrr.CreateSender(shoutrrrUrl)
	if err != nil {
		n.Logger.Errorf("An error occurred while sending notifications %v: %v", shoutrrrUrl, err)
		return
	}

	//sender.SetLogger(n.Logger.)
	serviceName, params, err := n.GenShoutrrrNotificationParams(shoutrrrUrl)
	n.Logger.Debug("notification data for %s: (%s)\n%v", serviceName, shoutrrrUrl, params)

	if err != nil {
		n.Logger.Errorf("An error occurred  occurred while generating notification payload for %s:\n %v", serviceName, shoutrrrUrl, err)
	}

	errs := sender.Send(n.Payload.GenerateMessage(), params)
	if len(errs) > 0 {
		n.Logger.Errorf("One or more errors occurred  occurred while sending notifications for %s:\n %v", shoutrrrUrl, errs)
		for _, err := range errs {
			n.Logger.Error(err)
		}
	}
}

func (n *Notify) GenShoutrrrNotificationParams(shoutrrrUrl string) (string, *shoutrrrTypes.Params, error) {
	serviceURL, err := url.Parse(shoutrrrUrl)
	if err != nil {
		return "", nil, err
	}

	serviceName := serviceURL.Scheme
	params := &shoutrrrTypes.Params{}

	logoUrl := "https://raw.githubusercontent.com/AnalogJ/scrutiny/master/webapp/frontend/src/ms-icon-144x144.png"
	subject := n.Payload.GenerateSubject()
	switch serviceName {
	// no params supported for these services
	case "discord", "hangouts", "ifttt", "mattermost", "teams":
		break
	case "gotify":
		(*params)["title"] = subject
	case "join":
		(*params)["title"] = subject
		(*params)["icon"] = logoUrl
	case "pushbullet":
		(*params)["title"] = subject
	case "pushover":
		(*params)["subject"] = subject
	case "slack":
		(*params)["title"] = subject
		(*params)["thumb_url"] = logoUrl
	case "smtp":
		(*params)["subject"] = subject
	case "standard":
		(*params)["subject"] = subject
	case "telegram":
		(*params)["subject"] = subject
	case "zulip":
		(*params)["topic"] = subject
	}

	return serviceName, params, nil
}

//utility functions
// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
