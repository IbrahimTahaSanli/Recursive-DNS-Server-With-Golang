package Query

import (
	"bytes"
	"encoding/binary"
	. "localDNS/tools"
)

type Question struct {
	Name  []string
	Type  uint16
	Class uint16
}

func parseQuestions(data []byte, questionCount uint16) []Question {
	var ret []Question = make([]Question, questionCount)

	var i uint16 = 0
	ind := 0
	for ; i < questionCount; i++ {
		ret[i].Name = make([]string, 0)
		for {
			if data[ind] == 0 {
				break
			}
			var dom string = string(data[ind+1 : ind+int(data[ind])+1])
			ret[i].Name = append(ret[i].Name, dom)
			ind = ind + int(data[ind]) + 1
		}

		ret[i].Type = ByteTouint16(data[1+ind : 3+ind])
		ret[i].Class = ByteTouint16(data[3+ind : 5+ind])
		ind += 5
	}

	return ret
}

func (ques Question) toByteArr() ([]byte, error) {
	retBytes := new(bytes.Buffer)
	var err error

	for _, item := range ques.Name {
		err = binary.Write(retBytes, binary.BigEndian, uint8(len(item)))
		if err != nil {
			return nil, err
		}
		err = binary.Write(retBytes, binary.BigEndian, []byte(item))
		if err != nil {
			return nil, err
		}
	}

	err = binary.Write(retBytes, binary.BigEndian, uint8(0))
	if err != nil {
		return nil, err
	}

	err = binary.Write(retBytes, binary.BigEndian, ques.Type)
	if err != nil {
		return nil, err
	}

	err = binary.Write(retBytes, binary.BigEndian, ques.Class)
	if err != nil {
		return nil, err
	}

	return retBytes.Bytes(), nil
}
