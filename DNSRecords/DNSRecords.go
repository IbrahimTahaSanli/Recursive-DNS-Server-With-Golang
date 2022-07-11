package DNSRecords

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"

	"localDNS/Query"
	. "localDNS/tools"
)

type _DNSRecords struct {
	mu      sync.RWMutex
	records []DNSRecord
}

type DNSRecord struct {
	Name  string
	Type  uint16
	Value string
	TTL   uint32
}

func valueToByteValueA(adr string) ([]byte, error) {
	var retVal [4]byte
	for ind, item := range strings.Split(adr, ".") {
		if ind == 4 {
			return nil, errors.New("Out Of Index")
		}
		val, err := strconv.ParseInt(item, 10, 0)
		if err != nil {
			return nil, err
		}
		retVal[ind] = uint8(val)
	}
	return retVal[:], nil
}
func valueToByteValueAAAA(adr string) ([]byte, error) {
	var retVal [16]byte
	for ind, item := range strings.Split(adr, ":") {
		if ind == 6 {
			return nil, errors.New("Out Of Index")
		}

		if item == "" {
			item = "00"
		}

		val, err := strconv.ParseInt(item, 16, 0)
		if err != nil {
			return nil, err
		}

		ByteArraySetArr(&retVal, ind*2, Uint16ToByteArr(uint16(val)))

	}

	return retVal[:], nil
}

func (rec DNSRecord) ToByteArr(indexes *map[string]uint8, currentIndex *uint16) ([]byte, error) {
	retBytes := new(bytes.Buffer)
	var err error

	err = binary.Write(retBytes, binary.BigEndian, uint8(0xc0))
	if err != nil {
		return nil, err
	}

	name, isName := (*indexes)[rec.Name]
	if !isName {
		return nil, err
	}

	err = binary.Write(retBytes, binary.BigEndian, name)
	if err != nil {
		return nil, err
	}

	err = binary.Write(retBytes, binary.BigEndian, rec.Type)
	if err != nil {
		return nil, err
	}

	err = binary.Write(retBytes, binary.BigEndian, uint16(1))
	if err != nil {
		return nil, err
	}

	err = binary.Write(retBytes, binary.BigEndian, rec.TTL)
	if err != nil {
		return nil, err
	}

	*currentIndex += 12
	(*indexes)[rec.Value] = uint8(*currentIndex)

	var tmp []byte
	switch rec.Type {
	case 1:
		tmp, err = valueToByteValueA(rec.Value)
	case 5:
		tmp, err = NameToByteArr(strings.Split(rec.Value, "."))
	case 16:
		tmp, err = TXTToByteArr(rec.Value)
	case 28:
		tmp, err = valueToByteValueAAAA(rec.Value)
	default:
		tmp = []byte(rec.Value)
		err = nil
	}
	if err != nil {
		return nil, err
	}

	err = binary.Write(retBytes, binary.BigEndian, uint16(len(tmp)))
	if err != nil {
		return nil, err
	}

	err = binary.Write(retBytes, binary.BigEndian, tmp)
	if err != nil {
		return nil, err
	}

	return retBytes.Bytes(), nil

}

var Records _DNSRecords
var _path *string

func InitDNSRecords(path *string) error {
	_path = path
	return readParseJson()
}

func readParseJson() error {
	cont, err := os.ReadFile(*_path)
	if err != nil {
		return err
	}

	json.Unmarshal(cont, &Records.records)

	return err
}

func RefreshRecords() error {
	Records.mu.Lock()
	defer Records.mu.Unlock()
	return readParseJson()
}

const A uint16 = 1
const NS uint16 = 2
const CNAME uint16 = 5
const MX uint16 = 15
const TXT uint16 = 16
const AAAA uint16 = 28

var NOTFOUNDERRROR = errors.New("NotFound")
var TBI = errors.New("TBI")

func GetRecords(ques Query.Question) (retRec []DNSRecord, err error) {
	switch ques.Type {
	case A:
		fallthrough
	case CNAME:
		fallthrough
	case TXT:
		fallthrough
	case AAAA:
		Records.mu.RLock()
		retRec = make([]DNSRecord, 0)

		addr := strings.Join(ques.Name, ".")

	CONTINUE:
		for _, record := range Records.records {
			if record.Name == addr {
				retRec = append(retRec, record)
				if record.Type == ques.Type {
					Records.mu.RUnlock()
					return
				} else if record.Type == CNAME {
					addr = record.Value
					goto CONTINUE
				}
				goto END
			}
		}
	default:
		return nil, TBI

	}
END:
	Records.mu.RUnlock()
	return nil, NOTFOUNDERRROR //notfound
}
