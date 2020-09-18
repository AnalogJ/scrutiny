package handler

import (
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/analogj/scrutiny/webapp/backend/pkg/notify"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// Send test notification
func SendTestNotification(c *gin.Context) {
	appConfig := c.MustGet("CONFIG").(config.Interface)

	testNotify := notify.Notify{
		Config: appConfig,
		Payload: notify.Payload{
			Mailer:       os.Args[0],
			Subject:      fmt.Sprintf("Scrutiny SMART error (EmailTest) detected on disk: XXXXX"),
			FailureType:  "EmailTest",
			Device:       "/dev/sda",
			DeviceType:   "ata",
			DeviceString: "/dev/sda",
			Message:      "TEST EMAIL from smartd for device: /dev/sda",
		},
	}
	err := testNotify.Send()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
		})
	} else {
		c.JSON(http.StatusOK, dbModels.DeviceWrapper{
			Success: true,
		})
	}
}
