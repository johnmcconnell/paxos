package paxos

// MessageTypes ...
const (
	VoteReq byte = iota
	VoteRes
	CheckReq
	CheckRes
	Confirmation
)

// Pretty ...
func Pretty(B byte) string {
	switch B {
	case VoteReq:
		return "VoteReq"
	case VoteRes:
		return "VoteRes"
	case CheckReq:
		return "CheckReq"
	case CheckRes:
		return "CheckRes"
	case Confirmation:
		return "Confirmation"
	}

	return "<unfound>"
}
