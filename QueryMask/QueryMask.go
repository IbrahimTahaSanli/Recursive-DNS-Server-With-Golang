package QueryMask

type QueryFlag uint16

const QRMASK QueryFlag = 0x8000

const OPCODE QueryFlag = 0x7800
const OPCODESTD QueryFlag = 0x0000
const OPCODEINV QueryFlag = 0x0800
const OPCODESTA QueryFlag = 0x1000

const AAMASK QueryFlag = 0x0400

const TCMASK QueryFlag = 0x0200

const RDMASK QueryFlag = 0x0100

const RAMASK QueryFlag = 0x0050

const RESERVEDMASK QueryFlag = 0x0070

const RCODEMASK QueryFlag = 0x000F
const RCODESUC QueryFlag = 0x0000
const RCODEFORMAT QueryFlag = 0x0001
const RCODESERFAIL QueryFlag = 0x0002
const RCODENAMEERR QueryFlag = 0x0003
const RCODENONSUP QueryFlag = 0x0004
const RCODEPOLICY QueryFlag = 0x0005

func (flag QueryFlag) IsQuery() QueryFlag {
	return flag & QRMASK
}

func (flag QueryFlag) SetQueResp(i int) QueryFlag {
	if i == 0 {
		return ^QRMASK & flag
	}
	return QRMASK | flag
}

func (flag QueryFlag) GetOpCode() QueryFlag {
	return flag & OPCODE
}

func (flag QueryFlag) SetOpCode(code QueryFlag) QueryFlag {
	return flag & ^OPCODE | QueryFlag(code)
}

func (flag QueryFlag) GetAthAns() QueryFlag {
	return flag & AAMASK
}

func (flag QueryFlag) SetAthAns(i int) QueryFlag {
	if i == 0 {
		return ^AAMASK & flag
	}
	return AAMASK | flag
}

func (flag QueryFlag) GetTruncation() QueryFlag {
	return flag & TCMASK
}

func (flag QueryFlag) SetTruncation(i int) QueryFlag {
	if i == 0 {
		return ^TCMASK & flag
	}
	return TCMASK | flag
}

func (flag QueryFlag) IsRecurDes() QueryFlag {
	return flag & RDMASK
}

func (flag QueryFlag) SetRecurDes(i int) QueryFlag {
	if i == 0 {
		return ^RDMASK & flag
	}
	return RDMASK | flag
}

func (flag QueryFlag) IsRecurAvb() QueryFlag {
	return flag & RAMASK
}

func (flag QueryFlag) SetRecurAvb(i int) QueryFlag {
	if i == 0 {
		return ^RAMASK & flag
	}
	return RAMASK | flag
}

func (flag QueryFlag) GetZCode() QueryFlag {
	return flag & RESERVEDMASK
}

func (flag QueryFlag) SetZCode(code QueryFlag) QueryFlag {
	return flag & ^RESERVEDMASK | QueryFlag(code)
}

func (flag QueryFlag) GetRCode() QueryFlag {
	return flag & RCODEMASK
}

func (flag QueryFlag) SetRCode(code QueryFlag) QueryFlag {
	return flag & ^RCODEMASK | QueryFlag(code)
}
