package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

// Middleware based on https://github.com/toorop/gin-logrus/blob/master/logger.go
// Body recording based on
// - https://github.com/gin-gonic/gin/issues/1363
// - https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin

// 2016-09-27 09:38:21.541541811 +0200 CEST
// 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700]
// "GET /apache_pb.gif HTTP/1.0" 200 2326
// "http://www.example.com/start.html"
// "Mozilla/4.08 [en] (Win98; I ;Nav)"

var timeFormat = "02/Jan/2006:15:04:05 -0700"

// Logger is the logrus logger handler
func LoggerMiddleware(logger logrus.FieldLogger) gin.HandlerFunc {

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknow"
	}

	return func(c *gin.Context) {

		//clone the request body reader.
		var reqBody string
		if c.Request.Body != nil {
			buf, _ := ioutil.ReadAll(c.Request.Body)
			reqBodyReader1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			reqBodyReader2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because reqBodyReader1 will be read.
			c.Request.Body = reqBodyReader2
			reqBody = readBody(reqBodyReader1)
		}

		// other handler can change c.Path so:
		path := c.Request.URL.Path
		blw := &responseBodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = blw
		c.Set("LOGGER", logger)
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		respLength := c.Writer.Size()
		if respLength < 0 {
			respLength = 0
		}

		entry := logger.WithFields(logrus.Fields{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"referer":    referer,
			"respLength": respLength,
			"userAgent":  clientUserAgent,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%s - %s [%s] \"%s %s\" %d %d \"%s\" \"%s\" (%dms)", clientIP, hostname, time.Now().Format(timeFormat), c.Request.Method, path, statusCode, respLength, referer, clientUserAgent, latency)
			if statusCode >= http.StatusInternalServerError {
				entry.Error(msg)
			} else if statusCode >= http.StatusBadRequest {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
		if strings.HasPrefix(path, "/api/") {
			//only debug log request/response from api endpoint.
			if len(reqBody) > 0 {
				entry.WithField("bodyType", "request").Debugln(reqBody) // Print request body
			}
			entry.WithField("bodyType", "response").Debugln(blw.body.String())
		}
	}
}

// Response Logging

type responseBodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Request Logging

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
