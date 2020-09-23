package handler

import (
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/analogj/scrutiny/webapp/backend/pkg/notify"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

// Send test notification
func SendTestNotification(c *gin.Context) {
	appConfig := c.MustGet("CONFIG").(config.Interface)
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)

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
		logger.Errorln("An error occurred while sending test notification", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
		})
	} else {
		c.JSON(http.StatusOK, dbModels.DeviceWrapper{
			Success: true,
		})
	}
}
