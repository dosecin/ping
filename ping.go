package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"os"
	"time"
)

type pinger struct {
	Seq     uint16
	Runloop int
	Timeout time.Duration
	Addr    string
	Data    []byte
}

func (p *pinger) ping(host string) {
	err := p.lookup(host)
	checkError(err)
	conn, err := net.Dial("ip4:icmp", p.Addr)
	checkError(err)
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(p.Timeout))
	report("正在 ping %s [%s] 具有 %d 字节的数据", host, p.Addr, len(p.Data))
	icmpHeader := &icmp{
		Type: icmpv4EchoRequest,
		Code: 0,
		ID:   uint16(os.Getpid() & 0xFFFF),
		Seq:  p.Seq,
	}
	var wb bytes.Buffer
	binary.Write(&wb, binary.BigEndian, icmpHeader)
	icmpHeader.Checksum = checkSum(wb.Bytes())
	wb.Reset()
	binary.Write(&wb, binary.BigEndian, p.Data)
	var rb = make([]byte, 20+len(p.Data))
	for i := 0; i < p.Runloop; i++ {
		if _, err = conn.Write(wb.Bytes()); err != nil {
			report("发送 %s 错误：%s", p.Addr, err.Error())
			return
		}
		size, err := conn.Read(rb)
		if err != nil {
			report("接收 %s 错误：%s", p.Addr, err.Error())
			return
		}
		report("来自 %s 的回复：字节=%d", p.Addr, size)
		conn.SetDeadline(time.Now().Add(p.Timeout))
		p.Seq++
	}
}

func (p *pinger) lookup(host string) error {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return err
	}
	if len(addrs) < 1 {
		return errors.New("unknown host")
	}
	p.Addr = addrs[0]
	return nil
}

func sendPingMsg(conn net.Conn) {

}
