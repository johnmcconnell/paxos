package paxos

// MessageTypes ...
const (
	VoteReq byte = iota
	VoteRes
	CheckReq
	CheckRes
	Confirmation
	LeaderReq
	LeaderRes
	JoinReq
	JoinRes
)

// Pretty Mappings ...
var (
	MessageStringMap = map[byte]string{
		VoteReq:      "VoteReq",
		VoteRes:      "VoteRes",
		CheckReq:     "CheckReq",
		CheckRes:     "CheckRes",
		Confirmation: "Confirmation",
		LeaderReq:    "LeaderReq",
		LeaderRes:    "LeaderRes",
		JoinReq:      "JoinReq",
		JoinRes:      "JoinRes",
	}
	MessageByteMap map[string]byte
)

func init() {
	MessageByteMap = make(map[string]byte)

	for k, v := range MessageStringMap {
		MessageByteMap[v] = k
	}
}

// Pretty ...
func Pretty(B byte) string {
	return MessageStringMap[B]
}

// Ugly ...
func Ugly(S string) byte {
	return MessageByteMap[S]
}
