package types

// Dns record types
const (
	RecordTypeA     = "A"
	RecordTypeAAAA  = "AAAA"
	RecordTypeCNAME = "CNAME"
	RecordTypeMX    = "MX"
	RecordTypeNS    = "NS"
	RecordTypeTXT   = "TXT"
	RecordTypeSRV   = "SRV"
	RecordTypeSOA   = "SOA"
	RecordTypePTR   = "PTR"
)

// dns recorf enum
type RecordType int

const (
	A RecordType = iota
	AAAA
	CNAME
	MX
	NS
	TXT
	SRV
	SOA
	PTR
)

func (r RecordType) String() string {
	switch r {
	case A:
		return RecordTypeA
	case AAAA:
		return RecordTypeAAAA
	case CNAME:
		return RecordTypeCNAME
	case MX:
		return RecordTypeMX
	case NS:
		return RecordTypeNS
	case TXT:
		return RecordTypeTXT
	case SRV:
		return RecordTypeSRV
	case SOA:
		return RecordTypeSOA
	case PTR:
		return RecordTypePTR
	default:
		return ""
	}
}
