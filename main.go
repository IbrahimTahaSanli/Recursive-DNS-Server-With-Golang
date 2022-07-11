package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"localDNS/DNSClient"
	"localDNS/DNSRecords"
	"localDNS/Logging"
	"localDNS/Query"
	"localDNS/Responses"
)

var listen *net.UDPConn

type connection struct {
	count int
	addr  *net.UDPAddr
	err   error
	data  []byte
}

type Config struct {
	UDPINPORT               uint16
	DATABUFFER              uint16
	DNSRECORDPATH           string
	ISRECURSIVE             bool
	OUTDNSSERVERS           []string
	UDPOUTPORTMIN           uint16
	UDPOUTPORTMAX           uint16
	LOGGINGENABLED          bool
	DATABASEPATH            string
	DATABASEUSERNAME        string
	DATABASEPASSWORD        string
	DATABASENAME            string
	DATABASECOLNAME         string
	DNSRECORDSRENEWINTERVAL int64
}

var CONFIG Config

func (conn connection) String() string {
	return fmt.Sprintf("{ \"Count\": %d, \"addr\": \"%s\", \"err\": \"%s\", \"data\": \"% x\"  }", conn.count, conn.addr, conn.err, conn.data)
}

func byteTouint16(bytes []byte) uint16 {
	return (uint16(bytes[0]) << 8) | uint16(bytes[1])
}

func readUdpData(conn connection) {
	start := time.Now()

	var resp []byte
	var err error

	if conn.err != nil {
		fmt.Println("Error Occured: %e", conn.err)
		return
	}

	que := *Query.ParseQuery(conn.data)

	if que.QuestionCount > 1 { // More than one question makes some complications so it doesnt a must
		if CONFIG.ISRECURSIVE {
			resp, err = DNSClient.SendMessage(conn.data)
		} else {
			resp, err = Responses.CreateNonSupportedQuery(que)
		}

		if err != DNSClient.RECVFAILED && err != nil {
			return
		}

		listen.WriteToUDP(resp, conn.addr)
		fmt.Println(time.Since(start))
		return
	}

	rec, err := DNSRecords.GetRecords(que.Questions[0])
	switch err {
	case nil:
		resp, err = Responses.CreateStdResponse(que, rec)
		if err != nil {
			return
		}

	case DNSRecords.NOTFOUNDERRROR:
		if CONFIG.ISRECURSIVE {
			resp, err = DNSClient.SendMessage(conn.data)
		} else {
			resp, err = Responses.CreateNotFoundResponse(que)
		}

		if err != DNSClient.RECVFAILED && err != nil {
			return
		}
	case DNSRecords.TBI:
		if CONFIG.ISRECURSIVE {
			resp, err = DNSClient.SendMessage(conn.data)
		} else {
			resp, err = Responses.CreateNonSupportedQuery(que)
		}

		if err != DNSClient.RECVFAILED && err != nil {
			return
		}
	default:
		resp, err = DNSClient.SendMessage(conn.data)
		if err != DNSClient.RECVFAILED && err != nil {
			return
		}

	}

	listen.WriteToUDP(resp, conn.addr)
	if CONFIG.LOGGINGENABLED {
		go Logging.Log(conn.data, resp, time.Since(start).Nanoseconds())
	}
}

func renewDNSRecords() {
	for {
		time.Sleep(time.Duration(CONFIG.DNSRECORDSRENEWINTERVAL))
		DNSRecords.RefreshRecords()
	}
}

func main() {
	var err error
	if len(os.Args) == 1 {
		conf, err := os.ReadFile("./config.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		json.Unmarshal(conf, &CONFIG)
	} else if len(os.Args) == 2 {
		conf, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		json.Unmarshal(conf, &CONFIG)
	} else {
		fmt.Println("So Many Arg")
		return
	}

	err = DNSRecords.InitDNSRecords(&CONFIG.DNSRECORDPATH)
	if err != nil {
		fmt.Println("%e", err)
		return
	}
	go renewDNSRecords()

	if CONFIG.ISRECURSIVE {
		err = DNSClient.InitClientMod(CONFIG.UDPOUTPORTMIN, CONFIG.UDPOUTPORTMAX, CONFIG.OUTDNSSERVERS)
		if err != nil {
			fmt.Println("%e", err)
			return
		}
	}
	if CONFIG.LOGGINGENABLED {
		Logging.InitLogger(CONFIG.DATABASEPATH, CONFIG.DATABASEUSERNAME, CONFIG.DATABASEPASSWORD, CONFIG.DATABASENAME, CONFIG.DATABASECOLNAME)
	}

	listen, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 53,
	})

	if err != nil {
		fmt.Println("Failed to establish monitoring!, err:", err)
		return
	}
	defer listen.Close()
	for {
		{
			_data := make([]byte, CONFIG.DATABUFFER)
			n, tmp, err := listen.ReadFromUDP(_data)
			var conn connection = connection{
				count: n,
				addr:  tmp,
				err:   err,
				data:  _data,
			}

			go readUdpData(conn)
		}
	}
}
