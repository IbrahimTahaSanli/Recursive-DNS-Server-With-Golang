package Responses

import (
	"localDNS/Query"
	"localDNS/QueryMask"
)

func CreateNonSupportedQuery(query Query.Query) ([]byte, error) {
	flags := query.Flags.SetQueResp(1)
	flags = flags.SetRCode(QueryMask.RCODENONSUP)

	query.Flags = flags

	resp, err := query.QueryToByteArr()

	if err != nil {
		return nil, err
	}

	return resp, nil

}
