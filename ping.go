package main

import (
	"errors"
	"math"
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
	report("正在 Ping %s [%s] 具有 %d 字节的数据", host, p.Addr, len(p.Data))
	icmpMsg := &icmp{
		Type: icmpv4EchoRequest,
		Code: 0,
		ID:   uint16(os.Getpid() & 0xFFFF),
		Seq:  p.Seq,
		Data: p.Data,
	}
	wb, _ := icmpMsg.marshal()
	var rb = make([]byte, 20+len(wb))
	var (
		sendMsgNum  int
		recvMsgNum  int
		minDuration = time.Duration(math.MaxInt64)
		maxDuration time.Duration
		sumDuration time.Duration
	)
	for i := 0; i < p.Runloop; i++ {
		start := time.Now()
		if _, err = conn.Write(wb); err != nil {
			report("发送 %s 错误：%s", p.Addr, err.Error())
			return
		}
		sendMsgNum++
		size, err := conn.Read(rb)
		if err != nil {
			report("接收 %s 错误：%s", p.Addr, err.Error())
			return
		}
		recvMsgNum++
		duration := time.Now().Sub(start)
		if duration < minDuration {
			minDuration = duration
		}
		if duration > maxDuration {
			maxDuration = duration
		}
		sumDuration += duration
		ttl := uint8(rb[8])
		dataSize := size
		if dataSize >= 20 {
			dataSize -= int(rb[0]&0x0f) << 2
		}
		dataSize -= 8
		report("来自 %s 的回复：字节=%d 时间=%dms TTL=%d", p.Addr, dataSize, duration/time.Millisecond, ttl)
		time.Sleep(1000 * time.Millisecond)
		conn.SetDeadline(time.Now().Add(p.Timeout))
	}
	report("\n%s 的 Ping 统计信息:", p.Addr)
	loseMsgNum := sendMsgNum - recvMsgNum
	report("    数据包：已发送 = %d，已接收 = %d， 丢失 = %d (%d%% 丢失)，", sendMsgNum, recvMsgNum, loseMsgNum, loseMsgNum*100/sendMsgNum)
	report("往返行程的估计时间(以毫秒为单位)：")
	report("    最短 = %dms，最长 = %dms，平均 = %dms", minDuration/time.Millisecond, maxDuration/time.Millisecond, int(sumDuration/time.Millisecond)/recvMsgNum)
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
