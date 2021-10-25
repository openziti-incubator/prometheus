package ziticonnections

import (
	"context"
	"net"
	"os"
	"strings"

	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"github.com/sirupsen/logrus"
)

var zitiConfigFile string

func GetZitiDialContextFunction(ctx context.Context, network string, addr string) (net.Conn, error) {

	//service := "nick service"
	configFile, err := config.NewFromFile(zitiConfigFile)

	if err != nil {
		logrus.WithError(err).Error("Error loading ziti config file")
		os.Exit(1)
	}

	serviceName := strings.Split(addr, ":")[0]

	context := ziti.NewContextWithConfig(configFile)
	return context.Dial(serviceName)
}

func SetZitiConfigFile(filePath string) {
	zitiConfigFile = filePath
}
