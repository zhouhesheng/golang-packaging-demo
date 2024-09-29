package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/denisbrodbeck/machineid"
)

type BroadcastInfo struct {
	ServerName      string   `json:"serverName"`
	ProtocolVersion string   `json:"protocolVersion"`
	ServerVersion   string   `json:"serverVersion"`
	ServerId        string   `json:"serverId"`
	LocalIp         string   `json:"localIp"`
	TsIp            string   `json:"tsIp"`
	Port            int      `json:"port"`
	LocalWebURLs    []string `json:"localWebURLs"`
	TsWebURLs       []string `json:"tsWebURLs"`
}

const UDP4MulticastAddress = "224.0.0.167:53315"

func main() {
	bc, err := NewBroadcaster(UDP4MulticastAddress)
	if err != nil {
		panic(err)
	}

	i := 0
	for {
		i++
		fmt.Printf("sending %d\n", i)
		_, err = bc.Write(MakeBroadcastInfo(53315))
		if err != nil {
			fmt.Printf("Send %v\n", err)
		}
		time.Sleep(time.Second * 1)
	}
}

func NewBroadcaster(address string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func MakeBroadcastInfo(port int) []byte {
	/*
	   hostnamectl
	    Static hostname: truenasBeta
	          Icon name: computer-desktop
	            Chassis: desktop 🖥️
	         Machine ID: 78826101253f4d05a9d6f2519dacfa8d
	            Boot ID: bef4f89fa9634d76aac6e7164d2b6fae
	   Operating System: Debian GNU/Linux 12 (bookworm)
	             Kernel: Linux 6.6.44-production+truenas
	       Architecture: x86-64
	    Hardware Vendor: ASUS
	     Hardware Model: PRIME H610M-A D4
	   Firmware Version: 3001
	*/

	info := BroadcastInfo{Port: port}
	hostname, err := os.Hostname()
	if err == nil {
		info.ServerName = hostname
	} else {
		fmt.Printf("Hostname error %v\n", err)
	}

	id, err := machineid.ID()
	if err == nil {
		info.ServerId = id
	} else {
		fmt.Printf("ServerId error %v\n", err)
	}

	info.LocalIp = GetLocalIP()

	jsonBody, _ := json.Marshal(info)
	return jsonBody
}
