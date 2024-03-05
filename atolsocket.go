package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	connHost = "localhost"
	connPort = "5555"
	connType = "tcp"
)

var pack int

func main() {
	pack = 0
	fmt.Println("сервис запущен")
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

func obrabotatattcppacket(tcppacket []byte) ([]byte, error) {
	var resanswer []byte
	//lendatafromtask := tcppacket[2]
	//rashrazmbuffera := 5 + lendatafromtask
	//answersended := false
	//if buffer8[rashrazmbuffera] == byte(0xfe) {
	//	fmt.Println("в пакете есть ещё задание")
	//	buffer8_2 := buffer8[rashrazmbuffera:]
	//	answerfortcppacket, _ := obrabotatattcppacket(buffer8_2)
	//	fmt.Printf("returncommandbytes=%x\n", answerfortcppacket)
	//	conn.Write(answerfortcppacket)
	//	answersended = true
	//}
	idtransp := tcppacket[3]
	commandbufer := tcppacket[4]
	command := byte(0x00)
	codeofanswer := byte(0xa3)
	//if commandbufer == byte(0xc3) { //Req
	//}
	if commandbufer == byte(0xc2) { //ack
		fmt.Println("подвтердить выполнения заданиря ack")
		crcresult := getCRC8([]byte{idtransp, codeofanswer})
		resanswer = []byte{0xfe, 0x01, 0x00, idtransp, codeofanswer, crcresult}
	}
	if commandbufer == byte(0xc1) { //новое задание
		fmt.Println("Новое задание")
		command = tcppacket[9]
	}
	if command == byte(0x3f) { //запрос состояния ККТ
		kassir := byte(0x00)
		nomvzale := byte(0x01)
		dateinkkt := []byte{0x23, 0x01, 0x01}
		timeinkkt := []byte{0x23, 0x02, 0x02}
		flagi := byte(0x4f)
		zavnom := []byte{0x23, 0x04, 0x54, 0x67}
		model := byte(0x43)
		rezerv := []byte{0x03, 0x00}
		regraboty := byte(0x00)
		nomcheck := []byte{0x23, 0x45}
		nomsmeny := []byte{0x00, 0x45}
		sostchecka := byte(0x00)
		summachaka := []byte{0x00, 0x00, 0x00, 0x00, 0x00}
		desyttoch := byte(0x00)
		portkkt := byte(0x06)
		dataanswer := append([]byte{0x44, kassir, nomvzale}, dateinkkt...)
		dataanswer = append(dataanswer, timeinkkt...)
		dataanswer = append(dataanswer, flagi)              //10
		dataanswer = append(dataanswer, zavnom...)          //14
		dataanswer = append(dataanswer, model)              //15
		dataanswer = append(dataanswer, rezerv...)          //17
		dataanswer = append(dataanswer, regraboty)          //18
		dataanswer = append(dataanswer, nomcheck...)        //20
		dataanswer = append(dataanswer, nomsmeny...)        //22
		dataanswer = append(dataanswer, sostchecka)         //23
		dataanswer = append(dataanswer, summachaka...)      //28
		dataanswer = append(dataanswer, desyttoch, portkkt) //30
		resanswer, _ = getanswerforkkt(idtransp, codeofanswer, dataanswer)
	}
	if command == byte(0xa5) { //получение параметров ККТ
		kodosh := byte(0x00)
		versprotocola := byte(0x02)
		typedevice := byte(0x01)
		modeldevice := byte(0x43)
		nazvankassy := []byte("Kass atol")
		dataanswer := append([]byte{kodosh, versprotocola, typedevice, modeldevice, 0x10, 0x00, 0x11, 0x23, 0x00, 0x23, 0x23}, nazvankassy...)
		resanswer, _ = getanswerforkkt(idtransp, codeofanswer, dataanswer)
	}
	if command == byte(0x91) { //получение значение регистра
		registrreq := tcppacket[10]
		fmt.Printf("Запрос значения регистра %v", registrreq)
		dataanswer, err := getdataforregistr(registrreq)
		if err != nil {
			fmt.Printf("ошибка (%v) получения данных по регситру %x\n", err, registrreq)
		}
		//kodotv := byte(0x55)
		//kodmist := byte(0x00)
		//valreg := []byte{0x20, 0x00, 0x64}
		//dataanswer := append([]byte{kodotv, kodmist}, valreg...)
		resanswer, _ = getanswerforkkt(idtransp, codeofanswer, dataanswer)
	}
	if command == byte(0x45) {
		fmt.Println("Запрос состояния ККТ")

		kodotv := byte(0x55)
		regim := byte(0x00)
		flagi := byte(0x00)
		dataanswer := []byte{kodotv, regim, flagi}
		resanswer, _ = getanswerforkkt(idtransp, codeofanswer, dataanswer)
	}
	return resanswer, nil
} //obrabotatattcppacket

func handleConnection(conn net.Conn) {
	//ar idask, idadd byte
	var buffer8 []byte
	//buffer, err := bufio.NewReader(conn).ReadBytes('\n')
	fmt.Printf("begin----------pack%v\n", pack)
	pack++
	fmt.Println("pack=", pack)
	fmt.Printf("%v-ый сценарий\n", pack)

	buffer8 = make([]byte, 100)
	razmtcpbuff, _ := bufio.NewReader(conn).Read(buffer8)
	if razmtcpbuff == 0 {
		fmt.Println("подтвердить ack")
		idtransploc := byte(0x0d)
		crcresult := getCRC8([]byte{idtransploc, 0xa3})
		ackreulst := []byte{0xfe, 0x01, 0x00, idtransploc, 0xa3, crcresult}
		fmt.Printf("answer=%x\n", ackreulst)
		conn.Write(ackreulst)
		return
		//handleConnection(conn)
		//return
	}
	fmt.Println("razmtcpbuff=", razmtcpbuff)
	fmt.Printf("buffer=%x\n", buffer8)
	lendatafromtask := buffer8[1]
	fmt.Println("lendatafromtask=", lendatafromtask)
	rashrazmbuffera := 5 + lendatafromtask
	fmt.Println("rashrazmbuffera", rashrazmbuffera)
	fmt.Println("buffer8[rashrazmbuffera]", buffer8[rashrazmbuffera])
	answersended := false
	if buffer8[rashrazmbuffera] == byte(0xfe) {
		fmt.Println("в пакете есть ещё задание")
		buffer8_2 := buffer8[rashrazmbuffera:]
		answerfortcppacket, _ := obrabotatattcppacket(buffer8_2)
		fmt.Printf("returncommandbytes=%x\n", answerfortcppacket)
		conn.Write(answerfortcppacket)
		answersended = true
	}
	idtransp := buffer8[3]
	commandbufer := buffer8[4]
	command := byte(0x00)
	codeofanswer := byte(0xa3)
	if commandbufer == byte(0xc4) { //очистка буфера
		crcbyteabort := getCRC8([]byte{0x00, codeofanswer})
		abortres := []byte{0xfe, 0x01, 0x00, 0x00, codeofanswer, crcbyteabort}
		fmt.Printf("answer=%x\n", abortres)
		conn.Write(abortres)
		answersended = true
	}
	if commandbufer == byte(0xc2) {
		fmt.Println("подтвердить ack")
		crcresult := getCRC8([]byte{idtransp, codeofanswer})
		ackreulst := []byte{0xfe, 0x01, 0x00, idtransp, codeofanswer, crcresult}
		fmt.Printf("answer=%x\n", ackreulst)
		conn.Write(ackreulst)
		answersended = true
	}
	if commandbufer == byte(0xc1) {
		fmt.Println("Новое задание")
		command = buffer8[9]
	}
	if command == byte(0xa5) { //получение параметров ККТ
		kodosh := byte(0x00)
		versprotocola := byte(0x02)
		typedevice := byte(0x01)
		modeldevice := byte(0x43)
		nazvankassy := []byte("Kass atol")
		dataanswer := append([]byte{kodosh, versprotocola, typedevice, modeldevice, 0x10, 0x00, 0x11, 0x23, 0x00, 0x23, 0x23}, nazvankassy...)
		returncommandbytes, _ := getanswerforkkt(idtransp, codeofanswer, dataanswer)
		fmt.Printf("returncommandbytes=%x\n", returncommandbytes)
		conn.Write(returncommandbytes)
	}
	if command == byte(0x91) { //получение значение регистра
		registrreq := buffer8[10]
		fmt.Printf("Запрос значения регистра %v", registrreq)
		dataanswer, err := getdataforregistr(registrreq)
		if err != nil {
			fmt.Printf("ошибка (%v) получения данных по регситру %x\n", err, registrreq)
		}
		//kodotv := byte(0x55)
		//kodmist := byte(0x00)
		//valreg := []byte{0x20, 0x00, 0x64}
		//dataanswer := append([]byte{kodotv, kodmist}, valreg...)
		returncommandbytes, _ := getanswerforkkt(idtransp, codeofanswer, dataanswer)
		fmt.Printf("returncommandbytes=%x\n", returncommandbytes)
		conn.Write(returncommandbytes)
		answersended = true
	}
	if command == byte(0x45) {
		fmt.Println("Запрос состояния ККТ")
		kodotv := byte(0x55)
		regim := byte(0x00)
		flagi := byte(0x00)
		dataanswer := []byte{kodotv, regim, flagi}
		returncommandbytes, _ := getanswerforkkt(idtransp, codeofanswer, dataanswer)
		//bytesforcrc := append([]byte{idtransp, codeofanswer}, dataanswer...)
		//databytes := append([]byte{codeofanswer}, dataanswer...)
		//crcidtranspwithdata := getCRC8(bytesforcrc)
		////crciddata2 := getCRC8([]byte{idask2, res2})
		//fmt.Println("len=", len(databytes))
		//lendata := byte(len(databytes))
		//fmt.Printf("lenbye=%x\n", lendata)
		//returncommandbytes := append([]byte{0xfe, lendata, 0x00}, idtransp)
		//returncommandbytes = append(returncommandbytes, databytes...)
		//returncommandbytes = append(returncommandbytes, crcidtranspwithdata)
		fmt.Printf("returncommandbytes=%x\n", returncommandbytes)
		conn.Write(returncommandbytes)
		answersended = true
	}
	if !answersended {
		fmt.Printf("неизвестнаыя кманда %x\n", command)
		return
	}
	fmt.Printf("end----------pack%v\n", pack)
	handleConnection(conn)
	fmt.Printf("end2----------pack%v\n", pack)
	return
}

func getdataforregistr(numbreg byte) ([]byte, error) {
	var dataofregistr []byte
	if numbreg == 24 { //ширина ленты
		fmt.Println("данные по регистру - ширина бумаги")
		kodotv := byte(0x55)
		kodmist := byte(0x00)
		valreg := []byte{0x20, 0x00, 0x64}
		dataofregistr = append([]byte{kodotv, kodmist}, valreg...)
	} else if numbreg == 54 { //
		fmt.Println("данные по регистру - ФФД")
		versffdKKT := byte(2) //ФФД1.05,
		versffdFN := byte(2)
		versffd := byte(2)
		dateFFD := []byte{0x17, 0x01, 0x01}
		maksverffdKKT := byte(2)
		maksverffdFN := byte(2)
		minimverffdKKT := byte(2)
		dataofregistr = append([]byte{versffdKKT, versffdFN, versffd}, dateFFD...)
		dataofregistr = append(dataofregistr, maksverffdKKT, maksverffdFN, minimverffdKKT)
	} else {
		return dataofregistr, fmt.Errorf("неизветсный (%x) формат регистра", numbreg)
	}
	return dataofregistr, nil
}

func getanswerforkkt(idtransp, codeofanswer byte, dataanswer []byte) ([]byte, error) {
	bytesforcrc := append([]byte{idtransp, codeofanswer}, dataanswer...)
	databytes := append([]byte{codeofanswer}, dataanswer...)
	crcidtranspwithdata := getCRC8(bytesforcrc)
	//crciddata2 := getCRC8([]byte{idask2, res2})
	//fmt.Println("len=", len(databytes))
	lendata := byte(len(databytes))
	//fmt.Printf("lenbye=%x\n", lendata)
	returncommandbytes := append([]byte{0xfe, lendata, 0x00}, idtransp)
	returncommandbytes = append(returncommandbytes, databytes...)
	returncommandbytes = append(returncommandbytes, crcidtranspwithdata)
	//fmt.Printf("returncommandbytes=%x\n", returncommandbytes)
	return returncommandbytes, nil
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

/*if pack == 1 {
	fmt.Println("1-ый сценарий очистка")
	buffer8 = make([]byte, 19)
	bufio.NewReader(conn).Read(buffer8)
	fmt.Printf("buffer=%x\n", buffer8)
	crcbyteabort := getCRC8([]byte{0x00, 0xa3})
	abortres := []byte{0xfe, 0x01, 0x00, 0x00, 0xa3, crcbyteabort}
	fmt.Printf("answer=%x\n", abortres)
	conn.Write(abortres)
	fmt.Printf("end----------pack%v", pack)
	handleConnection(conn)
	fmt.Printf("end2----------pack%v", pack)
	return
} else if pack == 2 {
	fmt.Println("2-ый сценарий принять задание")
	buffer8 = make([]byte, 11)
	bufio.NewReader(conn).Read(buffer8)
	fmt.Printf("buffer=%x\n", buffer8)
	crcinprogress := getCRC8([]byte{0x00, 0xa2})
	inprogressres := []byte{0xfe, 0x01, 0x00, 0x00, 0xa2, crcinprogress}
	fmt.Printf("answer=%x\n", inprogressres)
	conn.Write(inprogressres)
	fmt.Printf("end----------pack%v\n", pack)
	handleConnection(conn)
	fmt.Printf("end2----------pack%v\n", pack)
	return
} else if pack == 3 {
	fmt.Println("3-ый сценарий вернуть результат 1")
	buffer8 = make([]byte, 7)
	bufio.NewReader(conn).Read(buffer8)
	fmt.Printf("buffer=%x\n", buffer8)

	idtransp := byte(0x01)
	codeofanswer := byte(0xa3)
	kodosh := byte(0x00)
	versprotocola := byte(0x02)
	typedevice := byte(0x01)
	modeldevice := byte(0x43)
	//regim := []byte{0x10, 0x00}
	////verdevice:=[]byte{0x??,0x??,0x??,0x??,0x??}
	//verdevice := []byte{0x11, 0x23, 0x00, 0x23, 0x23}
	nazvankassy := []byte("Kass atol")
	dataanswer := append([]byte{kodosh, versprotocola, typedevice, modeldevice, 0x10, 0x00, 0x11, 0x23, 0x00, 0x23, 0x23}, nazvankassy...)
	bytesforcrc := append([]byte{idtransp, codeofanswer}, dataanswer...)
	databytes := append([]byte{codeofanswer}, dataanswer...)
	crcidtranspwithdata := getCRC8(bytesforcrc)
	//crciddata2 := getCRC8([]byte{idask2, res2})
	fmt.Println("len=", len(databytes))
	lendata := byte(len(databytes))
	fmt.Printf("lenbye=%x\n", lendata)
	returncommandbytes := append([]byte{0xfe, lendata, 0x00}, idtransp)
	returncommandbytes = append(returncommandbytes, databytes...)
	returncommandbytes = append(returncommandbytes, crcidtranspwithdata)
	fmt.Printf("returncommandbytes=%x\n", returncommandbytes)
	conn.Write(returncommandbytes)

	//crcinprogress := getCRC8([]byte{0x00, 0xa2})
	//abortres := []byte{0xfe, 0x01, 0x00, 0x00, 0xa2, crcinprogress}
	//conn.Write(abortres)
	fmt.Printf("end----------pack%v\n", pack)
	handleConnection(conn)
	fmt.Printf("end2----------pack%v\n", pack)
	return
} else if (pack == 4) || (pack == 5) || (pack == 7) || (pack == 8) || (pack == 9) {
	fmt.Printf("%v-ый сценарий\n", pack)
	buffer8 = make([]byte, 100)
	bufio.NewReader(conn).Read(buffer8)
	fmt.Printf("buffer=%x\n", buffer8)
	idtransp := buffer8[3]
	if buffer8[4] == byte(0xc3) {
		fmt.Println("вернуть результат")

		codeofanswer := byte(0xa3)
		kodosh := byte(0x00)
		versprotocola := byte(0x02)
		typedevice := byte(0x01)
		modeldevice := byte(0x43)
		//regim := []byte{0x10, 0x00}
		////verdevice:=[]byte{0x??,0x??,0x??,0x??,0x??}
		//verdevice := []byte{0x11, 0x23, 0x00, 0x23, 0x23}
		nazvankassy := []byte("Kass atol")
		dataanswer := append([]byte{kodosh, versprotocola, typedevice, modeldevice, 0x10, 0x00, 0x11, 0x23, 0x00, 0x23, 0x23}, nazvankassy...)
		bytesforcrc := append([]byte{idtransp, codeofanswer}, dataanswer...)
		databytes := append([]byte{codeofanswer}, dataanswer...)
		crcidtranspwithdata := getCRC8(bytesforcrc)
		//crciddata2 := getCRC8([]byte{idask2, res2})
		fmt.Println("len=", len(databytes))
		lendata := byte(len(databytes))
		fmt.Printf("lenbye=%x\n", lendata)
		returncommandbytes := append([]byte{0xfe, lendata, 0x00}, idtransp)
		returncommandbytes = append(returncommandbytes, databytes...)
		returncommandbytes = append(returncommandbytes, crcidtranspwithdata)
		fmt.Printf("returncommandbytes=%x\n", returncommandbytes)
		conn.Write(returncommandbytes)
	} else if buffer8[4] == byte(0xc2) {
		fmt.Println("подвтердить на ack")

		crcresult := getCRC8([]byte{idtransp, 0xa3})
		ackreulst := []byte{0xfe, 0x01, 0x00, idtransp, 0xa3, crcresult}
		fmt.Printf("answer=%x\n", ackreulst)
		conn.Write(ackreulst)
	} else {
		answerfortcppacket, _ := obrabotatattcppacket(buffer8)
		fmt.Printf("returncommandbytes=%x\n", answerfortcppacket)
		conn.Write(answerfortcppacket)
		//answersended = true
	}
	fmt.Printf("end----------pack%v\n", pack)
	handleConnection(conn)
	fmt.Printf("end2----------pack%v\n", pack)
	return
} else if pack == 6 {
	fmt.Printf("%v-ый сценарий опять задание\n", pack)
	buffer8 = make([]byte, 11)
	bufio.NewReader(conn).Read(buffer8)
	fmt.Printf("buffer=%x\n", buffer8)
	idtransp := buffer8[3]

	codeofanswer := byte(0xa2)
	crcinprogress := getCRC8([]byte{idtransp, codeofanswer})
	inprogressres := []byte{0xfe, 0x01, 0x00, idtransp, codeofanswer, crcinprogress}
	fmt.Printf("answer=%x\n", inprogressres)
	conn.Write(inprogressres)
	fmt.Printf("end----------pack%v\n", pack)
	handleConnection(conn)
	fmt.Printf("end2----------pack%v\n", pack)
	return
} else {*/

/*fmt.Printf("buffer=%x\n", buffer8)
crcbyteabort := getCRC8([]byte{0x00, 0xa3})
abortres := []byte{0xfe, 0x01, 0x00, 0x00, 0xa3, crcbyteabort}
conn.Write(abortres)
buffer9, err := bufio.NewReader(conn).ReadBytes('\x03')
fmt.Printf("buffer9=%x\n", buffer9)
//
//idask2 := buffer9[3]
idask2 := byte(0x03)
fmt.Printf("idask2=%x\n", idask2)
res2 := byte(0xa3)
res2asyn := byte(0xa6)
idask2asyn := byte(0xf0)
Tid2Asyn := byte(0x01)
kodosh := byte(0x00)
versprot := byte(0x02)
typedivece := byte(0x01)
modeldevice := byte(0x43)
//regim := []byte{0x10, 0x00}
////verdevice:=[]byte{0x??,0x??,0x??,0x??,0x??}
//verdevice := []byte{0x11, 0x23, 0x00, 0x23, 0x23}
nazvan := []byte("Kass atol")
darares := append([]byte{kodosh, versprot, typedivece, modeldevice, 0x10, 0x00, 0x11, 0x23, 0x00, 0x23, 0x23}, nazvan...)
byteforcrc := append([]byte{idask2, res2}, darares...)
databytes := append([]byte{res2}, darares...)
crciddata2 := getCRC8(byteforcrc)
//crciddata2 := getCRC8([]byte{idask2, res2})
fmt.Println("len=", len(databytes))
lendata2 := byte(len(databytes))
fmt.Printf("lenbye=%x\n", lendata2)
commandbytes2 := append([]byte{0xfe, lendata2, 0x00}, idask2)
commandbytes2 = append(commandbytes2, databytes...)
commandbytes2 = append(commandbytes2, crciddata2)
fmt.Printf("commandbytes1=%x\n", commandbytes2)
conn.Write(commandbytes2)
commandbytes3 := append([]byte{0xfe, lendata2, 0x00}, idask2asyn)
databytesasynch := append([]byte{res2asyn}, Tid2Asyn) //добавили Tid задания
databytesasynch = append(databytesasynch, darares...)
commandbytes3 = append(commandbytes3, databytesasynch...)
byteforcrcasynch := append([]byte{idask2asyn, res2asyn}, darares...)
crciddata2aync := getCRC8(byteforcrcasynch)
commandbytes3 = append(commandbytes3, crciddata2aync)
fmt.Printf("commandbytes3=%x\n", commandbytes3)
return
fmt.Println("----------")
//myfirstbuffer := make([]byte, 1)
//bufio.NewReader(conn).Read(myfirstbuffer)
//fmt.Printf("myfirstbyte=%x\n", myfirstbuffer)
//if byte(myfirstbuffer[0]) != byte(0xfe) {
//	fmt.Println("xeeeee")
//	return
//}
mybuffer := make([]byte, 5)
bufio.NewReader(conn).Read(mybuffer)
fmt.Printf("my=%x\n", mybuffer)
commandForKKT := mybuffer[len(mybuffer)-1:]
fmt.Printf("commandForKKT=%x\n", commandForKKT)
dl := byte(mybuffer[1])

dataforcom := []byte{}
if dl > 1 {
	dataforcom = make([]byte, dl-1+4)
	bufio.NewReader(conn).Read(dataforcom)
	fmt.Printf("dataforcom=%x\n", dataforcom)
	//dataforcom = byte(mybuffer[2])
}
fmt.Println("---dffff-----")
return

//b2, _ := bufio.NewReader(conn).ReadByte()
//fmt.Printf("b2=%x\n", b2)
//conn.Write([]byte("\x06"))
buffer, err := bufio.NewReader(conn).ReadBytes('\x05')
//buffer, err := bufio.NewReader(conn).ReadBytes('\x05')
fmt.Printf("buffer=%x\n", buffer)
fmt.Println("len(buffer)=", len(buffer))
idask = buffer[3]
fmt.Printf("idask=%x\n", idask)
res := byte(0xa3)
crciddata := getCRC8([]byte{idask, res})
lendata := byte(1)
commandbytes := append([]byte{0xfe, lendata, 0x00}, idask, res, crciddata)
fmt.Printf("commandbytes1=%x\n", commandbytes)
conn.Write(commandbytes)
//dataandid := []byte{}
//crs := getCRC8([]byte{idask})
//scom := []byte{0xfe, 0x00, 0x00, idask, crs}
if len(buffer) > 16 {
	idadd = buffer[9]
	fmt.Printf("idask=%x\n", idadd)
	crciddata = getCRC8([]byte{idadd, res})
	commandbytes := append([]byte{0xfe, lendata, 0x00}, idadd, res, crciddata)
	fmt.Printf("commandbytes2=%x\n", commandbytes)
	conn.Write(commandbytes)
	//dataandid = []byte{res} //inProgress
	//dataandid := append(dataandid, []byte())
	//fmt.Sprintf("%v\xa3%v\xa2", string(idask), string(idadd))
}
//dataandid := fmt.Sprintf("%v\xa3%v\xa2", string(idask), string(idadd))
//fmt.Printf("dataandid=%x\n", dataandid)
//crc := getCRC8([]byte(dataandid))
//fmt.Printf("crc=%x\n", crc)
////scom := fmt.Sprintf("\xfe\x02\x00%v\xa3%v\xa2%v", string(idask), string(idadd), string(crc))
//scom = fmt.Sprintf("\xfe\x01\x00%v%v", dataandid, string(crc))
////scom := fmt.Sprintf("\xfe\x02\x00%v\xa3%v\xa2%v", string(idask), string(idadd), string(crc))
//fmt.Printf("scom=%x\n", []byte(scom))
//conn.Write([]byte(scom))
//fmt.Printf("buffer=%x\n", buffer)
//conn.Write([]byte("\x06"))
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

handleConnection(conn)*/
