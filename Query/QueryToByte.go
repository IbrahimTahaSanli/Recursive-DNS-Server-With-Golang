package Query

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var COULDNTPARSEQUERY = errors.New("Couldnt Parse Query To Byte Array")

func (query Query) QueryToByteArr() ([]byte, error) {
	responseByte := new(bytes.Buffer)

	err := binary.Write(responseByte, binary.BigEndian, query.TransID)
	if err != nil {
		return nil, err
	}

	err = binary.Write(responseByte, binary.BigEndian, query.Flags)
	if err != nil {
		return nil, err
	}

	err = binary.Write(responseByte, binary.BigEndian, query.QuestionCount)
	if err != nil {
		return nil, err
	}

	err = binary.Write(responseByte, binary.BigEndian, query.AnswerRR)
	if err != nil {
		return nil, err
	}

	err = binary.Write(responseByte, binary.BigEndian, query.AuthorityRR)
	if err != nil {
		return nil, err
	}

	err = binary.Write(responseByte, binary.BigEndian, query.AditionalRR)
	if err != nil {
		return nil, err
	}

	for _, item := range query.Questions {
		ques, err := item.toByteArr()
		if ques == nil {
			return nil, err
		}
		err = binary.Write(responseByte, binary.BigEndian, ques)
		if err != nil {
			return nil, err
		}
	}

	return responseByte.Bytes(), nil
}
