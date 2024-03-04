package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	connHost = "localhost"
	connPort = "5555"
	connType = "tcp"
)

func main() {
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected.")

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	//buffer, err := bufio.NewReader(conn).ReadBytes('\n')
	//buffer, err := bufio.NewReader(conn).ReadBytes('\x02')
	fmt.Println("----------")
	//b2, _ := bufio.NewReader(conn).ReadByte()
	//fmt.Printf("b2=%x\n", b2)
	//conn.Write([]byte("\x06"))
	buffer, err := bufio.NewReader(conn).ReadBytes('\x05')
	//buffer, err := bufio.NewReader(conn).ReadBytes('\x05')
	idask := buffer[3]
	idadd := buffer[16]
	dataandid := fmt.Sprintf("%v\xa3%v\xa2", idask, idadd)
	crc := getCRC8([]byte(dataandid))
	fmt.Println("crc=", crc)
	scom := fmt.Sprintf("\xfe\x02\x00%v\xa3%v\xa2%v", idask, idadd, crc)
	fmt.Println("scom=", scom)
	conn.Write([]byte(scom))
	fmt.Printf("buffer=%x\n", buffer)
	conn.Write([]byte("\x06"))
	buffer2, err := bufio.NewReader(conn).ReadBytes('\x03')
	fmt.Printf("buffer2=%x\n", buffer2)
	b, err := bufio.NewReader(conn).ReadByte()
	fmt.Printf("b=%x\n", b)
	conn.Write([]byte("\x06"))
	buffer3, err := bufio.NewReader(conn).ReadBytes('\x04')
	fmt.Printf("buffer3=%x", buffer3)
	return
	conn.Write([]byte("\x06"))
	for i := 0; i < 4; i++ {
		b, err := bufio.NewReader(conn).ReadByte()
		if err != nil {
			conn.Write([]byte("\x06"))
			fmt.Println("ошибка")
			fmt.Println(err)
			continue
		}
		conn.Write([]byte("\x06"))
		fmt.Printf("%v=%x\n", i, b)
		//fmt.Printf("%v=%v\n", i, string(b))
	}
	return
	//buffer, err := bufio.NewReader(conn).ReadByte()
	buffer, err = bufio.NewReader(conn).ReadBytes('\x03')
	//buffer, err := bufio.NewReader(conn).ReadBytes('\xfe')

	if err != nil {
		fmt.Println("Client left.")
		conn.Close()
		return
	}

	log.Println("Client message:", string(buffer[:len(buffer)-1]))
	log.Println("Client message in bytes:", buffer[:len(buffer)-1])
	log.Printf("Client message in bytes 2: %x", buffer[:len(buffer)])

	//conn.Write([]byte("Hello, client!\n"))
	if buffer[0] == 0xfe {
		//if len(buffer) > 0 {
		log.Println("InProgress")
		conn.Write([]byte("\xfe\x01\x00\xa2"))
		//} else {
		//log.Println("Result")
		//conn.Write([]byte("\xa3OK"))
		//}
	} else {
		log.Println("send ask")
		conn.Write([]byte("\x06"))
	}

	//conn.Write(buffer)

	handleConnection(conn)
}

func handleClient(conn net.Conn) {
	defer conn.Close() // закрываем сокет при выходе из функции

	buf := make([]byte, 32) // буфер для чтения клиентских данных
	for {
		conn.Write([]byte("Hello, what's your name?\n")) // пишем в сокет

		readLen, err := conn.Read(buf) // читаем из сокета
		if err != nil {
			fmt.Println(err)
			break
		}

		conn.Write(append([]byte("Goodbye, "), buf[:readLen]...)) // пишем в сокет
	}
}

func getCRC8(iddata []byte) byte {
	valby := byte(0xff)
	for _, i := range iddata {
		valby = valby ^ i
		for i8 := 0; i8 < 8; i8++ {
			if (valby & byte(0x80)) != 0 {
				valby = ((valby << 1) ^ byte(0x31))
			} else {
				valby = (valby << 1)
			}
		}
	}
	valby = valby ^ byte(0x00)
	return valby
}
