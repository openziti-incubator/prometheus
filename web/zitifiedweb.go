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
	"golang.org/x/net/netutil"
	"net"
	"strings"

	zitiCfg "github.com/openziti/sdk-golang/ziti/config"
)

var ZitiIdServiceName string
var ZitiIdFile string

func init() {
	ZitiIdServiceName = "boundprometheus"
	ZitiIdFile = "/mnt/v/temp/prometheus/prometheusZitiIdentity.json"
}

// Listener creates the TCP listener for web requests.
func (h *Handler) Listener() (net.Listener, error) {
	level.Info(h.logger).Log("msg", "Start listening for connections", "address", h.options.ListenAddress)

	var listener net.Listener
	var err error
	if strings.HasPrefix(h.options.ListenAddress, "ziti") {
		zcfg, e := zitiCfg.NewFromFile(ZitiIdFile)
		if e != nil {
			return nil, errors.Wrap(e, "could not create Ziti config")
		}
		zitiContext := ziti.NewContextWithConfig(zcfg)
		zopts := &ziti.ListenOptions{
			BindUsingEdgeIdentity: true,
		}
		zopts.BindUsingEdgeIdentity = true
		listener, err = zitiContext.ListenWithOptions(ZitiIdServiceName, zopts)
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
