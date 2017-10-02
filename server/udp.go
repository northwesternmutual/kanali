// Copyright (c) 2017 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package server

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/logging"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/spf13/viper"
)

const (
	k8sNameMaxSize = 253
)

// StartUDPServer will start the udp server that is used to comminute between
// all running Kanali instances.
func StartUDPServer() (e error) {

	logger := logging.WithContext(nil)

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", viper.GetInt(config.FlagServerPeerUDPPort.GetLong())))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("upd server listening on :%s", viper.GetString(config.FlagServerPeerUDPPort.GetLong())))
	defer func() {
		if err := conn.Close(); err != nil {
			if e != nil {
				e = err
			}
		}
	}()

	// [NAMESPACE],[PROXYNAME],[KEYNAME] <= 761
	buf := make([]byte, k8sNameMaxSize*3+2)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}
		if err := spec.TrafficStore.Set(string(buf[0:n])); err != nil {
			logger.Error(err.Error())
		}
	}

}

// Emit will send a message to all other Kanali instances.
func Emit(binding spec.APIKeyBinding, keyName string, currTime time.Time) {

	logger := logging.WithContext(nil)

	for _, addr := range spec.KanaliEndpoints.Subsets[0].Addresses {

		if os.Getenv("POD_IP") == addr.IP {
			if err := spec.TrafficStore.Set(encodeKanaliGram(binding.ObjectMeta.Namespace, binding.Spec.APIProxyName, keyName, ",")); err != nil {
				logger.Error(err.Error())
			}
			continue
		}

		go func(ip string) {

			serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, viper.GetInt(config.FlagServerPeerUDPPort.GetLong())))
			if err != nil {
				logger.Warn(fmt.Sprintf("error resolving UDP address for %s:%d", ip, viper.GetInt(config.FlagServerPeerUDPPort.GetLong())))
				return
			}

			conn, err := net.DialUDP("udp", nil, serverAddr)
			if err != nil {
				logger.Warn(fmt.Sprintf("error dialing %s:%d", ip, viper.GetInt(config.FlagServerPeerUDPPort.GetLong())))
				return
			}

			_, err = conn.Write([]byte(fmt.Sprintf("%s,%s,%s", binding.ObjectMeta.Namespace, binding.Spec.APIProxyName, keyName)))
			if err != nil {
				logger.Warn(fmt.Sprintf("error writing traffic to %s:%d", ip, viper.GetInt(config.FlagServerPeerUDPPort.GetLong())))
				return
			}

			if err := conn.Close(); err != nil {
				logger.Error(err.Error())
			}

		}(addr.IP)

	}

}

func encodeKanaliGram(nSpace, pName, keyName, delimiter string) string {
	var buffer bytes.Buffer
	buffer.WriteString(nSpace)
	buffer.WriteString(delimiter)
	buffer.WriteString(pName)
	buffer.WriteString(delimiter)
	buffer.WriteString(keyName)
	return buffer.String()
}
