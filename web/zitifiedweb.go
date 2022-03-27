// Copyright 2013 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"github.com/go-kit/log/level"
	"github.com/influxdata/influxdb/kit/errors"
	conntrack "github.com/mwitkow/go-conntrack"
	"github.com/openziti/sdk-golang/ziti"
	zitiCfg "github.com/openziti/sdk-golang/ziti/config"
	"golang.org/x/net/netutil"
	"net"
	"net/url"
	"os"
	"time"
)

// Listener creates the TCP listener for web requests.
func (h *Handler) Listener() (net.Listener, error) {
	var scheme string
	u, err := url.Parse(h.options.ListenAddress)
	if err == nil {
		scheme = u.Scheme
	}

	level.Info(h.logger).Log("msg", "Start listening for connections", "address", h.options.ListenAddress)

	var listener net.Listener
	if scheme == "ziti" {

		serviceName := os.Getenv("ZITI_LISTENER_SERVICE_NAME")
		idFile := os.Getenv("ZITI_LISTENER_IDENTITY_FILE")
		idName := os.Getenv("ZITI_LISTENER_IDENTITY_NAME")

		zcfg, e := zitiCfg.NewFromFile(idFile)
		if e != nil {
			return nil, errors.Wrap(e, "could not create Ziti config")
		}
		zitiContext := ziti.NewContextWithConfig(zcfg)
		zopts := &ziti.ListenOptions{
			ConnectTimeout: 5 * time.Minute,
		}
		if idName != "" {
			zopts.Identity = idName
		}
		listener, err = zitiContext.ListenWithOptions(serviceName, zopts)
		//listener, err = zitiContext.Listen(serviceName)
		h.options.ListenAddress = ""
		level.Info(h.logger).Log("h.options.ListenAddress set to : ", h.options.ListenAddress)
	} else {
		listener, err = net.Listen("tcp", h.options.ListenAddress)
	}
	if err != nil {
		return nil, err
	}
	listener = netutil.LimitListener(listener, h.options.MaxConnections)

	// Monitor incoming connections with conntrack.
	listener = conntrack.NewListener(listener,
		conntrack.TrackWithName("http"),
		conntrack.TrackWithTracing())

	return listener, nil
}
