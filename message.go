package paxos

import (
	"fmt"
	"net"
	"strconv"
)

// Paxos Payloads
const (
	OK = iota
	NotOK
	Elected
	Election
	Yes
	No
	MessageSize = 16
)

// Pretty Mappings ...
var (
	PayloadStringMap = map[uint32]string{
		OK:       "OK",
		NotOK:    "NotOK",
		Elected:  "Elected",
		Election: "Election",
		Yes:      "Yes",
		No:       "No",
	}
	PayloadUintMap map[string]uint32
)

// Message ...
type Message struct {
	Payload uint32
	IP      uint32
	Port    uint32
	Round   uint32
}

func init() {
	PayloadUintMap = make(map[string]uint32)

	for k, v := range PayloadStringMap {
		PayloadUintMap[v] = k
	}
}

// String ...
func (m Message) String() string {
	S := PayloadStringMap[m.Payload]

	return fmt.Sprintf(
		"URL[%v] Status[%v] Round[%v]",
		m.URL(),
		S,
		m.Round,
	)
}

// BuildMessageS ...
func BuildMessageS(s, h, p, i string) (Message, error) {
	Zero := Message{}

	S := PayloadUintMap[s]

	I, err := strconv.ParseUint(i, 10, 32)

	if err != nil {
		return Zero, err
	}

	return BuildMessage2(S, h, p, uint32(I))
}

// BuildMessage ...
func BuildMessage(s uint32, addr *net.UDPAddr, i uint32) Message {
	IP, Port := AddrToUint32s(addr)

	m := Message{
		Payload: s,
		IP:      IP,
		Port:    Port,
		Round:   i,
	}

	return m
}

// BuildMessage2 ...
func BuildMessage2(s uint32, host, port string, i uint32) (Message, error) {
	Zero := Message{}

	ID, err := AddrSToUint64(
		host + ":" + port,
	)

	if err != nil {
		return Zero, err
	}

	IP, Port := Uint64ToUint32s(ID)

	m := Message{
		Payload: s,
		IP:      IP,
		Port:    Port,
		Round:   i,
	}

	return m, nil
}

// URL ...
func (m Message) URL() string {
	return fmt.Sprintf(
		"%v:%v",
		m.IPv4S(),
		m.Port,
	)
}

// PortS ...
func (m Message) PortS() string {
	return Uint32ToS(m.Port)
}

// ID ...
func (m Message) ID() uint64 {
	return Uint32sToUint64(m.IP, m.Port)
}

// IPv4B ....
func (m Message) IPv4B() (byte, byte, byte, byte) {
	return Uint32ToBytes(m.IP)
}

// IPv4S ....
func (m Message) IPv4S() string {
	B1, B2, B3, B4 := m.IPv4B()
	return BytesToIPS(B1, B2, B3, B4)
}

// UDPAddr ...
func (m Message) UDPAddr() (*net.UDPAddr, error) {
	addr, err := net.ResolveUDPAddr(
		"udp",
		m.URL(),
	)

	return addr, err
}

// Encoded ...
func (m Message) Encoded() []byte {
	return EncodeMessage(m)
}

// EncodeMessage ...
func EncodeMessage(m Message) []byte {
	var Encoded []byte
	EP := &Encoded

	EncodeUint(EP, m.Payload)
	EncodeUint(EP, m.IP)
	EncodeUint(EP, m.Port)
	EncodeUint(EP, m.Round)

	return *EP
}

// DecodeMessage ...
func DecodeMessage(BS []byte) Message {
	if len(BS) < MessageSize {
		panic(
			fmt.Sprintf(
				"Need 16 bytes to decode message, instead of %v",
				len(BS),
			),
		)
	}

	M := Message{}

	DecodeUint(BS[0:4], &M.Payload)
	DecodeUint(BS[4:8], &M.IP)
	DecodeUint(BS[8:12], &M.Port)
	DecodeUint(BS[12:16], &M.Round)

	return M
}
