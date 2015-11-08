package paxos

import (
	"fmt"
	"net"
	"strconv"
	"strings"
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
	Mask        = 0xFF
)

// Message ...
type Message struct {
	Payload uint32
	IP      uint32
	Port    uint32
	Round   uint32
}

// String ...
func (m Message) String() string {
	S := ""

	switch m.Payload {
	case OK:
		S = "OK"
	case NotOK:
		S = "NotOK"
	case Elected:
		S = "Elected"
	}

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

	var S uint32

	switch s {
	case "OK":
		S = OK
	case "NotOK":
		S = NotOK
	case "Elected":
		S = Elected
	}

	I, err := strconv.ParseUint(i, 10, 32)

	if err != nil {
		return Zero, err
	}

	URL := h + ":" + p

	Addr, err := net.ResolveUDPAddr(
		"udp",
		URL,
	)

	if err != nil {
		return Zero, err
	}

	return BuildMessage(S, Addr, uint32(I)), nil
}

// BuildMessage ...
func BuildMessage(s uint32, addr *net.UDPAddr, i uint32) Message {
	IPv4 := addr.IP.To4()

	var IP uint32

	IP += uint32(IPv4[0])
	IP = IP << 8
	IP += uint32(IPv4[1])
	IP = IP << 8
	IP += uint32(IPv4[2])
	IP = IP << 8
	IP += uint32(IPv4[3])

	m := Message{
		Payload: s,
		IP:      IP,
		Port:    uint32(addr.Port),
		Round:   i,
	}

	return m
}

// BuildMessage2 ...
func BuildMessage2(s uint32, host, port string, i uint32) (Message, error) {
	Zero := Message{}

	Port, err := strconv.ParseUint(port, 10, 32)

	if err != nil {
		return Zero, err
	}

	Subs := strings.Split(host, ".")
	if len(Subs) != 4 {
		return Zero, fmt.Errorf(
			"Invalid length host[%v]",
			host,
		)
	}

	Part, err := strconv.Atoi(Subs[0])
	if err != nil {
		return Zero, err
	}
	IP := uint32(Part) << 24

	Part, err = strconv.Atoi(Subs[1])
	if err != nil {
		return Zero, err
	}
	IP += uint32(Part) << 16

	Part, err = strconv.Atoi(Subs[2])
	if err != nil {
		return Zero, err
	}
	IP += uint32(Part) << 8

	Part, err = strconv.Atoi(Subs[3])
	if err != nil {
		return Zero, err
	}
	IP += uint32(Part)

	m := Message{
		Payload: s,
		IP:      IP,
		Port:    uint32(Port),
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
	return fmt.Sprintf(
		"%v",
		m.Port,
	)
}

// ID ...
func (m Message) ID() uint64 {
	ID := uint64(m.IP) << 32
	ID += uint64(m.Port)

	return ID
}

// IPv4B ....
func (m Message) IPv4B() (byte, byte, byte, byte) {
	B1 := byte(m.IP >> 24)
	B2 := byte((m.IP >> 16) & Mask)
	B3 := byte((m.IP >> 8) & Mask)
	B4 := byte(m.IP & Mask)

	return B1, B2, B3, B4
}

// IPv4S ....
func (m Message) IPv4S() string {
	B1, B2, B3, B4 := m.IPv4B()

	return fmt.Sprintf(
		"%v.%v.%v.%v",
		B1,
		B2,
		B3,
		B4,
	)
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

// EncodeUint ...
func EncodeUint(B *[]byte, V uint32) {
	var I uint32

	for I < 4 {
		Shift := (I * 8)
		Part := (V >> Shift)

		E := byte(Part & Mask)
		(*B) = append(*B, E)

		I++
	}
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

// DecodeUint ...
func DecodeUint(B []byte, V *uint32) {
	var I uint32

	for I < 4 {
		Shift := (I * 8)
		Part := uint32(B[I])

		(*V) |= uint32(Part << Shift)

		I++
	}
}
