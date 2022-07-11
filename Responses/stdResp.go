package Responses

import (
	"localDNS/DNSRecords"
	"localDNS/Query"
	"localDNS/QueryMask"
	"strings"
)

func CreateStdResponse(query Query.Query, records []DNSRecords.DNSRecord) ([]byte, error) {

	flags := query.Flags.SetQueResp(1)           //It is Response
	flags = flags.SetOpCode(QueryMask.OPCODESTD) //Its a Standart Request
	flags = flags.SetAthAns(1)                   //AuthAns
	flags = flags.SetTruncation(0)               //I am passing this for now
	//Recursion Avaliable comes from header
	flags = flags.SetRecurAvb(1) //Recursion is avaliable by the way it coded;
	//Z which is reserved
	flags = flags.SetRCode(QueryMask.RCODESUC) //RCODE for succes

	query.Flags = flags
	query.AnswerRR = uint16(len(records))

	var arr, tmpArr []byte
	var err error

	arr, err = query.QueryToByteArr()
	if err != nil {
		return nil, err
	}

	currentIndex := uint16(len(arr))

	indexes := make(map[string]uint8)
	indexes[strings.Join(query.Questions[0].Name, ".")] = uint8(12)

	for _, rec := range records {
		tmpArr, err = rec.ToByteArr(&indexes, &currentIndex)
		if err != nil {
			return nil, err
		}

		arr = append(arr, tmpArr...)

	}

	return arr, nil

}
