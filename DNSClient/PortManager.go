package DNSClient

type portStack chan uint16

var _portStack portStack

func InitClientMod(minPort uint16, maxPort uint16, _servers []string) error {
	_portStack = make(portStack, maxPort-minPort)

	for i := uint16(0); i < maxPort-minPort; i++ {
		_portStack <- uint16(i) + minPort
	}

	servers = _servers

	return nil
}

func getPort() uint16 {
	return <-_portStack
}

func putPort(port uint16) {
	_portStack <- port
}
