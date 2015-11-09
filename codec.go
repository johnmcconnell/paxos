package paxos

import (
	"fmt"
	"net"
)

// Masks ...
const (
	ByteMask = 0xFF
	U32Mask  = 0xFFFFFFFF
)

// Uint64ToAddr ...
func Uint64ToAddr(V uint64) (*net.UDPAddr, error) {
	A, B := Uint64ToUint32s(V)

	B1, B2, B3, B4 := Uint32ToBytes(A)

	IP := BytesToIPS(B1, B2, B3, B4)
	Port := Uint32ToS(B)

	return net.ResolveUDPAddr(
		"udp",
		IP+":"+Port,
	)
}

// Uint32ToS ...
func Uint32ToS(V uint32) string {
	return fmt.Sprintf(
		"%v",
		V,
	)
}

// AddrToUint64 ...
func AddrToUint64(Addr *net.UDPAddr) uint64 {
	I, P := AddrToUint32s(Addr)

	return Uint32sToUint64(I, P)
}

// AddrToUint32s ...
func AddrToUint32s(Addr *net.UDPAddr) (uint32, uint32) {
	IPv4 := Addr.IP.To4()

	var IP uint32

	for I := 0; I < 4; I++ {
		IP = IP << 8
		IP |= uint32(IPv4[I])
	}

	return IP, uint32(Addr.Port)
}

// Uint32ToBytes ...
func Uint32ToBytes(V uint32) (byte, byte, byte, byte) {
	B1 := byte(V >> 24)
	B2 := byte((V >> 16) & ByteMask)
	B3 := byte((V >> 8) & ByteMask)
	B4 := byte(V & ByteMask)

	return B1, B2, B3, B4
}

// BytesToIPS ...
func BytesToIPS(B1, B2, B3, B4 byte) string {
	return fmt.Sprintf(
		"%v.%v.%v.%v",
		B1,
		B2,
		B3,
		B4,
	)
}

// AddrSToUint64 ...
func AddrSToUint64(AddrS string) (uint64, error) {
	addr, err := net.ResolveUDPAddr(
		"udp",
		AddrS,
	)

	if err != nil {
		return 0, err
	}

	IP, Port := AddrToUint32s(addr)

	return Uint32sToUint64(IP, Port), nil
}

// Uint32sToUint64 ...
func Uint32sToUint64(A, B uint32) uint64 {
	C := uint64(A) << 32
	C |= uint64(B)

	return C
}

// Uint64ToUint32s ...
func Uint64ToUint32s(C uint64) (uint32, uint32) {
	A := uint32(C >> 32)
	B := uint32(C & U32Mask)

	return A, B
}

// EncodeUint ...
func EncodeUint(B *[]byte, V uint32) {
	var I uint32

	for I < 4 {
		Shift := I * 8
		Part := V >> Shift

		E := byte(Part & ByteMask)
		(*B) = append(*B, E)

		I++
	}
}

// DecodeUint ...
func DecodeUint(B []byte, V *uint32) {
	var I uint32

	for I < 4 {
		Shift := I * 8
		Part := uint32(B[I])

		(*V) |= uint32(Part << Shift)

		I++
	}
}
