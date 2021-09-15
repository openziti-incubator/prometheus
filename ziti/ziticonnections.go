package ziticonnections

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"github.com/sirupsen/logrus"
)

func GetZitiDialContextFunction(context.Context, string, string) (net.Conn, error) {
	service := "nick service"
	configFile, err := config.NewFromFile("/home/npieros/Downloads/nickendpoint03.json")

	if err != nil {
		logrus.WithError(err).Error("Error loading ziti config file")
		os.Exit(1)
	}

	context := ziti.NewContextWithConfig(configFile)
	return context.Dial(service)
}

func GetZitiListener(context.Context, string, string) net.Listener {
	service := "nick service"
	configFile, err := config.NewFromFile("/home/npieros/Downloads/nickendpoint03.json")

	if err != nil {
		logrus.WithError(err).Error("Error loading ziti config file")
		os.Exit(1)
	}

	context := ziti.NewContextWithConfig(configFile)

	listener, err := context.Listen(service)

	if err != nil {
		panic(err)
	}
	return listener
}

func GetZitiConn(timeout time.Duration) (net.Conn, error) {
	service := "nick service"
	configFile, err := config.NewFromFile("/home/npieros/Downloads/nickendpoint03.json")

	if err != nil {
		logrus.WithError(err).Error("Error loading ziti config file")
		os.Exit(1)
	}

	context := ziti.NewContextWithConfig(configFile)
	return context.DialWithOptions(service, &ziti.DialOptions{ConnectTimeout: timeout})
}
