package paxos

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {
	assert := assert.New(t)

	IP := uint32(127 << 24)
	IP += 10 << 16
	IP += 10 << 8
	IP += 2

	m := Message{
		Payload: NotOK,
		IP:      IP,
		Port:    8888,
		Round:   42,
	}

	r := DecodeMessage(
		EncodeMessage(
			m,
		),
	)

	assert.Equal(
		m,
		r,
		"The two messages match",
	)
}

func TestURL(t *testing.T) {
	assert := assert.New(t)

	IP := uint32(127 << 24)
	IP += 10 << 16
	IP += 10 << 8
	IP += 2

	m := Message{
		Payload: NotOK,
		IP:      IP,
		Port:    8888,
		Round:   42,
	}

	assert.Equal(
		"127.10.10.2:8888",
		m.URL(),
		"The urls match",
	)
}

func TestBuildS(t *testing.T) {
	assert := assert.New(t)

	IP := uint32(127 << 24)
	IP += 10 << 16
	IP += 10 << 8
	IP += 2

	M := Message{
		Payload: NotOK,
		IP:      IP,
		Port:    3000,
		Round:   67,
	}

	NewM, err := BuildMessageS(
		"NotOK",
		"127.10.10.2",
		"3000",
		"67",
	)

	assert.Nil(
		err,
		"Parsed Correctly",
	)

	assert.Equal(
		M,
		NewM,
		"they match",
	)
}

func TestBuild(t *testing.T) {
	assert := assert.New(t)

	IP := uint32(127 << 24)
	IP += 10 << 16
	IP += 10 << 8
	IP += 2

	M := Message{
		Payload: OK,
		IP:      IP,
		Port:    3000,
		Round:   67,
	}

	URL := M.URL()

	assert.Equal(
		"127.10.10.2:3000",
		URL,
		"URL match",
	)

	ExpectedAddr, err := net.ResolveUDPAddr(
		"udp",
		URL,
	)

	assert.Nil(
		err,
		"no error",
	)

	Addr, err := M.UDPAddr()

	assert.Nil(
		err,
		"no error",
	)

	assert.Equal(
		*ExpectedAddr,
		*Addr,
	)

	NewM := BuildMessage(
		M.Payload,
		Addr,
		M.Round,
	)

	assert.Equal(
		M,
		NewM,
		"they match",
	)
}

func TestEncode(t *testing.T) {
	assert := assert.New(t)

	IP := uint32(127 << 24)
	IP += 10 << 16
	IP += 10 << 8
	IP += 2

	m := Message{
		Payload: OK,
		IP:      IP,
		Port:    2,
		Round:   67,
	}

	Encoded := EncodeMessage(
		m,
	)

	assert.Equal(
		16,
		len(Encoded),
		"Only takes 16 bytes",
	)
}
