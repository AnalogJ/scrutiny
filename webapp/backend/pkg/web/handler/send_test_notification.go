package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/analogj/scrutiny/webapp/backend/pkg/notify"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Send test notification
func SendTestNotification(c *gin.Context) {
	appConfig := c.MustGet("CONFIG").(config.Interface)
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)

	testNotify := notify.Notify{
		Logger: logger,
		Config: appConfig,
		Payload: notify.Payload{
			FailureType:  "EmailTest",
			DeviceSerial: "FAKEWDDJ324KSO",
			DeviceType:   dbModels.DeviceProtocolAta,
			DeviceName:   "/dev/sda",
			Test:         true,
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
