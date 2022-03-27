package ziticonnections

import (
	"context"
	"net"
	"os"
	"strings"
	"time"

	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"github.com/sirupsen/logrus"
)

func GetZitiDialContext(zitiConfigPath string) func(ctx context.Context, network string, addr string) (net.Conn, error) {
	configFile, err := config.NewFromFile(zitiConfigPath)

	if err != nil {
		logrus.WithError(err).Error("Error loading ziti config file")
		os.Exit(1)
	}

	zitiCtx := ziti.NewContextWithConfig(configFile)

	dialContextFunc := func(ctx context.Context, network string, addr string) (net.Conn, error) {
		zitiDialInfo := strings.Split(addr, ":")[0]

		zitiDial := strings.Split(zitiDialInfo, "-")
		serviceName := zitiDial[0]
		identity := zitiDial[1]

		dialOptions := &ziti.DialOptions{
			Identity:       identity,
			ConnectTimeout: time.Minute,
		}

		return zitiCtx.DialWithOptions(serviceName, dialOptions)

	}

	return dialContextFunc
}
