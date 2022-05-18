module github.com/analogj/scrutiny

go 1.13

require (
	github.com/analogj/go-util v0.0.0-20190301173314-5295e364eb14
	github.com/citilinkru/libudev v1.0.0
	github.com/containrrr/shoutrrr v0.4.4
	github.com/fatih/color v1.10.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-gormigrate/gormigrate/v2 v2.0.0
	github.com/golang/mock v1.4.3
	github.com/google/uuid v1.2.0 // indirect
	github.com/influxdata/influxdb-client-go/v2 v2.8.2
	github.com/jaypipes/ghw v0.6.1
	github.com/jinzhu/gorm v1.9.16
	github.com/klauspost/compress v1.12.1 // indirect
	github.com/kvz/logstreamer v0.0.0-20150507115422-a635b98146f0 // indirect
	github.com/mattn/go-sqlite3 v1.14.4 // indirect
	github.com/mitchellh/mapstructure v1.2.2
	github.com/onsi/ginkgo v1.16.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.5.1
	github.com/urfave/cli/v2 v2.2.0
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect
	gorm.io/driver/sqlite v1.1.3
	gorm.io/gorm v1.20.2
	nhooyr.io/websocket v1.8.7 // indirect
)

// Remove once the following PR/Issues have been merged:
// - https://github.com/influxdata/influxdb-client-go/pull/328
// - https://github.com/influxdata/influxdb-client-go/issues/327
replace github.com/influxdata/influxdb-client-go/v2 => github.com/analogj/influxdb-client-go/v2 v2.8.2-jk
