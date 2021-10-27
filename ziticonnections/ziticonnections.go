package ziticonnections

import (
	"context"
	"net"
	"os"
	"strings"
	"time"

	"github.com/mwitkow/go-conntrack"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"github.com/sirupsen/logrus"
)

var zitiConfigFile string

func GetZitiDialContextFunction(ctx context.Context, network string, addr string) (net.Conn, error) {

	configFile, err := config.NewFromFile(zitiConfigFile)

	if err != nil {
		logrus.WithError(err).Error("Error loading ziti config file")
		os.Exit(1)
	}
	addressString := strings.Split(addr, ":")[0]
	zitiDialInfo := strings.Split(addressString, ".")

	if len(zitiDialInfo) > 1 && strings.ToLower(zitiDialInfo[1]) == "ziti" {
		zitiDial := strings.Split(zitiDialInfo[0], "-")
		serviceName := zitiDial[0]
		identity := zitiDial[1]

		dialOptions := &ziti.DialOptions{
			Identity:       identity,
			ConnectTimeout: time.Minute,
		}

		context := ziti.NewContextWithConfig(configFile)
		return context.DialWithOptions(serviceName, dialOptions)
	} else {
		dialContext := conntrack.NewDialContextFunc(
			conntrack.DialWithTracing(),
			conntrack.DialWithName(addr))

		return dialContext(ctx, network, addr)
	}
}

func SetZitiConfigFile(filePath string) {
	zitiConfigFile = filePath
}
