package Query

import (
	"localDNS/QueryMask"
	. "localDNS/tools"
)

type Query struct {
	TransID       uint16
	Flags         QueryMask.QueryFlag
	QuestionCount uint16
	AnswerRR      uint16
	AuthorityRR   uint16
	AditionalRR   uint16
	Questions     []Question
}

func GetRCodeFromQuery(data []byte) QueryMask.QueryFlag {
	return QueryMask.QueryFlag(ByteTouint16(data[2:4])) & QueryMask.RCODEMASK
}

func ParseQuery(data []byte) (ret *Query) {
	ret = &Query{
		TransID:       ByteTouint16(data[0:2]),
		Flags:         QueryMask.QueryFlag(ByteTouint16(data[2:4])),
		QuestionCount: ByteTouint16(data[4:6]),
		AnswerRR:      ByteTouint16(data[6:8]),
		AuthorityRR:   ByteTouint16(data[8:10]),
		AditionalRR:   ByteTouint16(data[10:12]),
		Questions:     parseQuestions(data[12:], ByteTouint16(data[4:6])),
	}
	return
}
