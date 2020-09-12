package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/analogj/go-util/utils"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/containrrr/shoutrrr"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Payload struct {
	Mailer       string `json:"mailer"`
	Subject      string `json:"subject"`
	Date         string `json:"date"`
	FailureType  string `json:"failure_type"`
	Device       string `json:"device"`
	DeviceType   string `json:"device_type"`
	DeviceString string `json:"device_string"`
	Message      string `json:"message"`
}

type Notify struct {
	Config  config.Interface
	Payload Payload
}

func (n *Notify) Send() error {
	//validate that the Payload is populated
	sendDate := time.Now()
	n.Payload.Date = sendDate.Format(time.RFC3339)

	//retrieve list of notification endpoints from config file
	configUrls := n.Config.GetStringSlice("notify.urls")

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
	if len(notifyScripts) > 0 {
		wg.Add(1)
		go n.SendShoutrrrNotification(&wg, notifyShoutrrr)
	}

	//and wait for completion, error or timeout.
	if waitTimeout(&wg, time.Minute) { //wait for 1 minute
		fmt.Println("Timed out while sending notifications")
	} else {
		fmt.Println("Sent notifications. Check logs for more information.")
	}
	return nil
}

func (n *Notify) SendWebhookNotification(wg *sync.WaitGroup, webhookUrl string) {
	defer wg.Done()
	log.Infof("Sending Webhook to %s", webhookUrl)
	requestBody, err := json.Marshal(n.Payload)
	if err != nil {
		log.Errorf("An error occurred while sending Webhook to %s: %v", webhookUrl, err)
		return
	}

	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Errorf("An error occurred while sending Webhook to %s: %v", webhookUrl, err)
		return
	}
	defer resp.Body.Close()
	//we don't care about resp body content, but maybe we should log it?
}

func (n *Notify) SendScriptNotification(wg *sync.WaitGroup, scriptUrl string) {
	defer wg.Done()

	//check if the script exists.
	scriptPath := strings.TrimPrefix(scriptUrl, "script://")
	log.Infof("Executing Script %s", scriptPath)

	if !utils.FileExists(scriptPath) {
		log.Errorf("Script does not exist: %s", scriptPath)
		return
	}

	copyEnv := os.Environ()
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_MAILER=%s", n.Payload.Mailer))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_SUBJECT=%s", n.Payload.Subject))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_DATE=%s", n.Payload.Date))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_FAILURE_TYPE=%s", n.Payload.FailureType))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_DEVICE=%s", n.Payload.Device))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_DEVICE_TYPE=%s", n.Payload.DeviceType))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_DEVICE_STRING=%s", n.Payload.DeviceString))
	copyEnv = append(copyEnv, fmt.Sprintf("SCRUTINY_MESSAGE=%s", n.Payload.Message))
	err := utils.CmdExec(scriptPath, []string{}, "", copyEnv, "")
	if err != nil {
		log.Errorf("An error occurred while executing script %s: %v", scriptPath, err)
	}
	return
}

func (n *Notify) SendShoutrrrNotification(wg *sync.WaitGroup, shoutrrrUrls []string) {
	log.Infof("Sending notifications to %v", shoutrrrUrls)

	defer wg.Done()
	sender, err := shoutrrr.CreateSender(shoutrrrUrls...)
	if err != nil {
		log.Errorf("An error occurred while sending notifications %v: %v", shoutrrrUrls, err)
		return
	}

	errs := sender.Send(n.Payload.Subject, nil) //structs.Map(n.Payload).())
	if len(errs) > 0 {
		log.Errorf("One or more errors occurred  occurred while sending notifications %v:\n %v", shoutrrrUrls, errs)
	}
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
