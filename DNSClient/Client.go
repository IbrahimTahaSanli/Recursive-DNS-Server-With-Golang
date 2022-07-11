package DNSClient

import (
	"bufio"
	"errors"
	"fmt"
	"net"

	"localDNS/Query"
	"localDNS/QueryMask"
)

var servers []string

var RECVFAILED error = errors.New("Recursive Search Failed")

func SendMessage(data []byte) ([]byte, error) {
	var port uint16 = getPort()
	defer putPort(port)

	respData := make([]byte, 512)

	local, err := net.ResolveUDPAddr("udp", ":"+fmt.Sprintf("%d", port))
	if err != nil {
		return nil, err
	}
	for _, server := range servers {
		remote, err := net.ResolveUDPAddr("udp", server+":53")
		if err != nil {
			continue
		}
		conn, err := net.DialUDP("udp", local, remote)
		if err != nil {
			continue
		}

		conn.Write(data)

		_, err = bufio.NewReader(conn).Read(respData)
		if err != nil {
			return nil, err
		}

		if Query.GetRCodeFromQuery(respData) == QueryMask.RCODESUC {
			return respData, nil
		}
		conn.Close()
	}

	return data, RECVFAILED
}
