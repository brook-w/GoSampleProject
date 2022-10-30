package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"
)

// ICMP:https://baike.baidu.com/item/ICMP/572452

var (
	timeout      int64
	size         int
	count        int
	typ          uint8 = 8
	code         uint8 = 0
	sendCount    int
	successCount int
	failCount    int
	minTs        int64 = math.MaxInt32
	maxTs        int64
	totalTs      int64
)

type ICMP struct {
	Type     uint8
	Code     uint8
	CheckSum uint16
	ID       uint16
	Sequence uint16
}

func main() {
	getCommandArgs()
	desIP := os.Args[len(os.Args)-1]

	conn, err := net.DialTimeout("ip:icmp", desIP, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	fmt.Printf("正在 Ping %s %s 具有 %d 字节的数据: \n", desIP, conn.RemoteAddr(), size)

	for i := 0; i < count; i++ {
		sendCount++
		icmp := &ICMP{
			Type:     typ,
			Code:     code,
			CheckSum: 0,
			ID:       1,
			Sequence: 1,
		}

		data := make([]byte, size)
		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, icmp)
		buffer.Write(data)
		data = buffer.Bytes()

		checkSum := checkSum(data)
		data[2] = byte(checkSum >> 8)
		data[3] = byte(checkSum)

		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
		startTime := time.Now()
		_, err = conn.Write(data)
		if err != nil {
			failCount++
			log.Println(err)
			return
		}

		rtuBuf := make([]byte, 1<<16)
		readCount, err := conn.Read(rtuBuf)
		if err != nil {
			log.Println(err)
			return
		}
		successCount++

		endTime := time.Since(startTime).Milliseconds()
		if minTs > endTime {
			minTs = endTime
		}
		if maxTs < endTime {
			maxTs = endTime
		}
		totalTs += endTime

		fmt.Printf(
			"来自 %d.%d.%d.%d 的回复: 字节=%d 时间=%dms TTL=%d\n",
			rtuBuf[12],
			rtuBuf[13],
			rtuBuf[14],
			rtuBuf[15],
			readCount-28,
			endTime,
			rtuBuf[8],
		)
		time.Sleep(time.Millisecond * 500)
	}

	fmt.Println()
	fmt.Printf(`%s 的 Ping 统计信息:
	数据包: 已发送 = %d，已接收 = %d，丢失 = %d (%.2f%% 丢失)，
往返行程的估计时间(以毫秒为单位):    
	最短 = %dms，最长 = %dms，平均 = %dms`,
		conn.RemoteAddr(),
		sendCount, successCount, failCount, float64(failCount)/float64(sendCount),
		minTs, maxTs, totalTs/int64(sendCount),
	)
}

func getCommandArgs() {
	flag.Int64Var(&timeout, "w", 1000, "请求超时时长（ms）")
	flag.IntVar(&size, "l", 32, "缓冲区大小（byte）")
	flag.IntVar(&count, "n", 4, "请求次数")
	flag.Parse()
}

func checkSum(data []byte) uint16 {
	length := len(data)
	index := 0
	var sum uint32 = 0
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		length -= 2
		index += 2
	}
	if length != 0 {
		sum += uint32(data[index])
	}
	hi16 := sum >> 16
	for hi16 != 0 {
		sum = hi16 + uint32(uint16(sum))
		hi16 = sum >> 16
	}
	return uint16(^sum)
}
