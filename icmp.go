package main

import "encoding/binary"

const (
	icmpv4EchoRequest = 8
	icmpv4EchoReply   = 0
)

type icmp struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	ID       uint16
	Seq      uint16
	Data     []byte
}

func checkSum(b []byte) uint16 {
	csumcv := len(b) - 1
	s := uint32(0)
	for i := 0; i < csumcv; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if csumcv&1 == 0 {
		s += uint32(b[csumcv])
	}
	s = s>>16 + s&0xffff
	s = s + s>>16
	return ^uint16(s)
}

func (p *icmp) marshal() ([]byte, error) {
	wb := make([]byte, 8+len(p.Data))
	wb[0] = p.Type
	wb[1] = p.Code
	wb[2] = 0
	wb[3] = 0
	binary.BigEndian.PutUint16(wb[4:6], p.ID)
	binary.BigEndian.PutUint16(wb[6:8], p.Seq)
	copy(wb[8:], p.Data)
	s := checkSum(wb)
	wb[2] ^= byte(s)
	wb[3] ^= byte(s >> 8)
	return wb, nil
}
