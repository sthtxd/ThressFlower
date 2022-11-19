package socket

import (
	"encoding/binary"
	"fmt"
	"kmpi-go/config"
	"kmpi-go/log"
	"kmpi-go/service"
	"math"
	"net"
	"strconv"
	"strings"
	"time"
)

// func StartSocketServer(port string) error {
// 	service := ":" + port
// 	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
// 	if err != nil {
// 		log.Error("StartSocket error", err.Error())
// 		return err
// 	}
// 	listener, err := net.ListenTCP("tcp", tcpAddr)
// 	if err != nil {
// 		log.Error("ListenTCP error", err.Error())
// 		return err
// 	}
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			continue
// 		}
// 		go holdFurClientA(conn, 10)
// 	}

// }

func StartSocketClient() {

	ipListStr, _ := config.GetValue("device_ip", "ipList")
	connType, _ := config.GetValue("device_ip", "connType")
	//分割字符串
	arr := strings.Split(ipListStr, ",")
	arrType := strings.Split(connType, ",")

	if len(arr) == 0 {
		panic("read device ipList  length is 0")

	}

	if len(arrType) == 0 {
		panic("read device Type  length is 0")

	}

	chVesselSend := make(chan int, 1)
	chVesselRecv := make(chan int, 1)
	chCtrlBoxSend := make(chan int, 1)
	chCtrlBoxRecv := make(chan int, 1)
	chFur1Send := make(chan int, 1)
	chFur1Recv := make(chan int, 1)
	chFur2Send := make(chan int, 1)
	chFur2Recv := make(chan int, 1)

	//遍历ip，每个ip开启一个线程去连接

	for index, ip := range arr {
		nType, _ := strconv.Atoi(arrType[index])
		go startClient(ip, nType, chVesselSend, chVesselRecv, chCtrlBoxSend, chCtrlBoxRecv, chFur1Send, chFur1Recv, chFur2Send, chFur2Recv)
	}
	//	go startClient("Vessel", 2, chVesselSend, chVesselRecv, chCtrlBoxSend, chCtrlBoxRecv, chFur1Send, chFur1Recv, chFur2Send, chFur2Recv)
	//	go startClient("AGV", 8, chVesselSend, chVesselRecv, chCtrlBoxSend, chCtrlBoxRecv, chFur1Send, chFur1Recv, chFur2Send, chFur2Recv)
}

func startClient(ip string, nType int, chVSend chan int, chVRecv chan int, chBoxSend chan int, chBoxRecv chan int, chF1Send chan int, chF1Recv chan int, chF2Send chan int, chF2Recv chan int) {
	defer reconnectClient(ip, nType, chVSend, chVRecv, chBoxSend, chBoxRecv, chF1Send, chF1Recv, chF2Send, chF2Recv)
	fmt.Println(ip)
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Error("connect " + ip + " fail")
		time.Sleep(time.Second * 5)
			return
	}
	switch nType {
	case 1: //熔炉
		deviceSmelterClient(conn, nType)
	case 2: //汤包
		deviceVesselClient(conn, nType, chVRecv, chVSend)
	case 3: //除气机
		deviceDegasClient(conn, nType)
	case 4: //烤包器
		deviceStoveClient(conn, nType)
	case 5: //定量炉
		deviceFurnaceClient(conn, nType, chF1Recv, chF1Send)
	case 6: //定量炉
		deviceFurnaceClient(conn, nType, chF2Recv, chF2Send)
	case 7: //烤包除气加铝信号箱
		deviceSignalBoxClient(conn, chBoxRecv, chBoxSend)
	case 8: //agv 控制
		AgvDispatch(conn, chVSend, chVRecv,chBoxSend, chBoxRecv, chF1Send, chF1Recv, chF2Send, chF2Recv,ip )
	case 9: //熔炉2
		deviceSmelterO2Client(conn, nType)
	default:
		reconnectClient(ip, nType, chVSend, chVRecv, chBoxSend, chBoxRecv, chF1Send, chF1Recv, chF2Send, chF2Recv)
	}

	
}

func reconnectClient(ip string, nType int, chVSend chan int, chVRecv chan int, chBoxSend chan int, chBoxRecv chan int, chF1Send chan int, chF1Recv chan int, chF2Send chan int, chF2Recv chan int) {
	go startClient(ip, nType, chVSend, chVRecv, chBoxSend, chBoxRecv, chF1Send, chF1Recv, chF2Send, chF2Recv)
	fmt.Println("disconnect and start connecting")
}

func deviceSmelterClient(conn net.Conn, nDeviceId int) {
	defer conn.Close()

	readCommandTotalData := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x03, 0x03, 0x00, 0x3C, 0x00, 0x01}
	

	for {
		//读取熔炉
		writeLen, err := conn.Write(readCommandTotalData)
		if err != nil {
			log.Error("Send command err,", err.Error())
			return
		}
		fmt.Println(writeLen)

		resByte := make([]byte, 16)
		_, err = conn.Read(resByte)
		fmt.Println(resByte)
		if err != nil {
			log.Error("Read err,", err.Error())
			return
		}

		nTemp := int(resByte[9])*256+int(resByte[10])

//		nWeight := int(resByte[1])

		
		go deviceSmelterInsertSql(0, nTemp, 0, nDeviceId)

		time.Sleep(time.Second)
	}
}

func deviceSmelterO2Client(conn net.Conn, nDeviceId int){
	readO2 := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x01, 0x03, 0x00, 0x60, 0x00, 0x01}
	for{
	//读取O2 Sensor
		writeLen, err := conn.Write(readO2)
		if err != nil {
			log.Error("Send command err,", err.Error())
			return
		}
		fmt.Println(writeLen)
		resByte := make([]byte, 16)
		_, err = conn.Read(resByte)
		fmt.Println(resByte)
		if err != nil {
			log.Error("Read err,", err.Error())
			return
		}

		fCurrent := float64(int(resByte[9])*256+int(resByte[9])) / 1000.0
		fO2 := (fCurrent-4)/16.0*(0.25-0.001) + 0.001
		
		go deviceSmelterInsertSql(0, 0, fO2, nDeviceId)
	}
}

func deviceWeightSensor(conn net.Conn, nDeviceId int) {
	readSensor := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x01, 0x03, 0x00, 0x00, 0x00, 0x02}
	for{
	//读取称重传感器
	writeLen, err := conn.Write(readSensor)
	if err != nil {
		log.Error("Send Sensor command err,", err.Error())
		return
	}
	fmt.Println(writeLen)

	resByte := make([]byte, 13)
	_, err = conn.Read(resByte)
	fmt.Println(resByte)
	if err != nil {
		log.Error("Read err,", err.Error())
		return
	}
	nWeight := int(uint16(resByte[11])*256 + uint16(resByte[12]))

	go deviceWeightDataInsertSql(nWeight, 1)
	}
}

func deviceSmelterInsertSql(nWeight int, nTemp int, nO2 float64, nDeviceId int) {
	err := service.InsertSqlserver("熔炉", 1, nTemp, 0, nWeight, 0, 0.0, 1, nO2, nDeviceId)
	if err != nil {
		log.Error("Read err,", err.Error())
	}
}

func deviceWeightDataInsertSql(nWeight int, nDeviceId int) {
	err := service.InsertWeightSensor(nWeight, nDeviceId)
	if err != nil {
		log.Error("insert weight err,", err.Error())
	}
}

func deviceVesselClient(conn net.Conn, nDeviceId int, chVSend chan int, chVRecv chan int) {
	defer conn.Close()

	readCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x01, 0x03, 0x00, 0x00, 0x00, 0x015}
	writeCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x01, 0x06, 0x00, 0x00, 0x00, 0x00}
	for {
		var chVData int

		select {
		case chVData = <-chVRecv:

		default:
		}
		resByte := make([]byte, 39)
		if chVData == 1 {
			fmt.Println("开始加铝")
			writeCommand[11] = 1
			writeLen, err := conn.Write(writeCommand)
			if err != nil {
				log.Error("Send command err,", err.Error())
				return
			}
			fmt.Println(writeLen)

			_, err = conn.Read(resByte)
			fmt.Println(resByte)
			if err != nil {
				log.Error("Read err,", err.Error())
				return
			}
		}
		//读取汤包
		writeLen, err := conn.Write(readCommand)
		if err != nil {
			log.Error("Send command err,", err.Error())
			return
		}
		fmt.Println(writeLen)

		_, err = conn.Read(resByte)
		fmt.Println(resByte)
		if err != nil {
			log.Error("Read err,", err.Error())
			return
		}
		convByte := make([]byte, 4)
		convByte[0] = resByte[9]
		convByte[1] = resByte[10]
		convByte[2] = resByte[11]
		convByte[3] = resByte[12]
		fO2 := float64(byteToFloat(convByte))

		convByte[0] = resByte[13]
		convByte[1] = resByte[14]
		convByte[2] = resByte[15]
		convByte[3] = resByte[16]
		fTemp := float64(byteToFloat(convByte))
		nTemp := int(fTemp)

		convByte[0] = resByte[21]
		convByte[1] = resByte[22]
		convByte[2] = resByte[23]
		convByte[3] = resByte[24]
		fWeight := float64(byteToFloat(convByte))
		nWeight := int(fWeight)
		//加铝结束
		if resByte[25] == 1 {
			chVSend <- 1
			fmt.Println("结束加铝")
			writeCommand[11] = 0
			writeLen, err := conn.Write(writeCommand)
			if err != nil {
				log.Error("Send command err,", err.Error())
				return
			}
			fmt.Println(writeLen)

			resByte := make([]byte, 25)
			_, err = conn.Read(resByte)
			fmt.Println(resByte)
			if err != nil {
				log.Error("Read err,", err.Error())
				return
			}
		}
		go deviceVesselInsertSql(fO2, nTemp, nWeight, 2)
		time.Sleep(time.Second * 5)
	}
}

func deviceVesselInsertSql(fO2 float64, nTemp int, nWeight int, nDeviceId int) {
	err := service.InsertSqlserver("汤包", 1, nTemp, 1, nWeight, 0, 0.0, 1, fO2, nDeviceId)
	if err != nil {
		log.Error("insert vessel err,", err.Error())
	}
}

func deviceDegasClient(conn net.Conn, nDeviceId int) {
	defer conn.Close()

	readH2 := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x01, 0x03, 0x00, 0x60, 0x00, 0x01}
	for {
		//读取除气机
		writeLen, err := conn.Write(readH2)
		if err != nil {
			log.Error("Send command err,", err.Error())
			return
		}
		fmt.Println(writeLen)

		resByte := make([]byte, 16)
		_, err = conn.Read(resByte)
		fmt.Println(resByte)
		if err != nil {
			log.Error("Read err,", err.Error())
			return
		}

		fO2 := float64(resByte[0])

		go deviceDegasInsertSql(fO2, 3)
		time.Sleep(time.Second * 5)
	}
}

func deviceDegasInsertSql(fO2 float64, nDeviceId int) {
	err := service.InsertSqlserver("除气机", 0, 0, 0, 0, 0, 0.0, 1, fO2, nDeviceId)
	if err != nil {
		log.Error("insert vessel err,", err.Error())
	}
}

func deviceStoveClient(conn net.Conn, nDeviceId int) {
	defer conn.Close()

	readCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x01, 0x03, 0x20, 0x0, 0x00, 0x01}
	for {
		//读取烤包器
		writeLen, err := conn.Write(readCommand)
		if err != nil {
			log.Error("Send command err,", err.Error())
			return
		}
		fmt.Println(writeLen)

		resByte := make([]byte, 11)
		_, err = conn.Read(resByte)
		fmt.Println(resByte)
		if err != nil {
			log.Error("Read err,", err.Error())
			return
		}
		nTemp := int(uint16(resByte[9])*256 + uint16(resByte[10]))

		go deviceStoveInsertSql(nTemp, 3)
		time.Sleep(time.Second * 5)
	}
}

func deviceStoveInsertSql(nTemp int, nDeviceId int) {
	err := service.InsertSqlserver("烤包器", 1, nTemp, 0, 0, 0, 0.0, 0, 0.0, nDeviceId)
	if err != nil {
		log.Error("insert vessel err,", err.Error())
	}
}

func deviceFurnaceClient(conn net.Conn, nIndex int, chFurSend chan int, chFurRecv chan int) {
	defer conn.Close()

	readCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x01, 0x03, 0x00, 0x2D, 0x00, 0x9}
	times:=0
	for {
		nData := 0
		select {
		case nData = <-chFurRecv:

		default:
		}
		if nData != 0 {
			if furnaceOperateFillDoor(nData, conn) {
				chFurSend <- 1
				nData = 0
			}
		}

		times++
		//50 读取一次数据
		if times>=10 {
		times=0
		_, err := conn.Write(readCommand)
		if err != nil {
			log.Error("Send command err,", err.Error())
			return
		}

		resByte := make([]byte, 27)
		_, err = conn.Read(resByte)
		fmt.Println(resByte)
		if err != nil {
			log.Error("Read err,", err.Error())
			return
		}

		nTemp := int(uint16(resByte[9])*256 + uint16(resByte[10]))/10
		nWeight := int(uint16(resByte[11])*256 + uint16(resByte[12]))
		convByte := make([]byte, 4)
		convByte[0] = resByte[20]
		convByte[1] = resByte[19]
		convByte[2] = resByte[22]
		convByte[3] = resByte[21]
		fO2 := float64(byteToFloat(convByte))
		convByte[0] = resByte[24]
		convByte[1] = resByte[23]
		convByte[2] = resByte[26]
		convByte[3] = resByte[25]
		fH2 := float64(byteToFloat(convByte))

		fmt.Println("nTemp,nWeight,fO2, fH2,",nTemp, nWeight, fO2, fH2)
		go deviceFurnaceInsertSql(nTemp, nWeight, fO2, fH2, nIndex)
		}
		time.Sleep(time.Second * 5)
	}
}

func deviceFurnaceInsertSql(nTemp int, nWeight int, fO2 float64, fH2 float64, nDeviceId int) {
	sName := "定量炉" + strconv.Itoa(nDeviceId-4)
	err := service.InsertSqlserver(sName, 1, nTemp, 1, nWeight, 1, fH2, 1, fO2, nDeviceId)
	if err != nil {
		log.Error("insert vessel err,", err.Error())
	}
}

func furnaceOperateFillDoor(operateType int, conn net.Conn) bool {
	bOk := false
	//第十个字节0xFF表示DO输出1，0x00表示输出0
	writeCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x01, 0x05, 0x00, 0x00, 0xFF, 0x00}
	if operateType == OPERATE_OPEN {
		writeCommand[10] = 0XFF
	} else if operateType == OPERATE_CLOSE {
		writeCommand[10] = 0X0
	}
	writeLen, err := conn.Write(writeCommand)
	if err != nil {
		log.Error("operate fill door command err,", err.Error())
		return bOk
	}
	fmt.Println(writeLen)

	resByte := make([]byte, 12)
	_, err = conn.Read(resByte)
	fmt.Println(resByte)
	if err != nil {
		log.Error("operate fill door Read err,", err.Error())
		return bOk
	}

	readCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x01, 0x02, 0x00, 0x20, 0x00, 0x04}
	writeLen, err = conn.Write(readCommand)
	if err != nil {
		log.Error("Send  operate fill door command err,", err.Error())
		return bOk
	}
	fmt.Println(writeLen)

	_, err = conn.Read(resByte)
	fmt.Println(resByte)
	if err != nil {
		log.Error("Read  operate fill door err,", err.Error())
		return bOk
	}
	result := 0
	if operateType == OPERATE_OPEN {
		result = 1
	} else if operateType == OPERATE_CLOSE {
		result = 0
	}
	if int(resByte[9]) == result {
		bOk = true
	}
	return bOk
}

func deviceSignalBoxClient(conn net.Conn, chBoxSend chan int, chBoxRecv chan int) {
	defer conn.Close()

	readCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x11, 0x02, 0x00, 0x20, 0x00, 0x04}
	var nLastCmd byte
	nLastCmd = 0
	for {

		_, err := conn.Write(readCommand)
		if err != nil {
			log.Error("SignalBox Send command err,", err.Error())
			return
		}
		

		resByte := make([]byte, 10)
		_, err = conn.Read(resByte)
		
		if err != nil {
			log.Error("SignalBox Read err,", err.Error())
			return
		}
		nCmd := resByte[9]

		if nCmd != nLastCmd {
			nLastCmd = nCmd
			chBoxSend <- int(nCmd)
			fmt.Println("signalbox send",nCmd)
		}
		select {
		case <-chBoxRecv:
			if resetSignal(conn) {
				chBoxSend <- 11 //reset signal success
				fmt.Println("signalbox reset signal success")
			} else {
				chBoxSend <- 12 //reset signal fail
				fmt.Println("signalbox reset signal fail")
			}
		default:
		}

		time.Sleep(time.Second * 2)
	}
}

func deviceRfidInsertSql(nRfid int, nDeviceId int) {
	err := service.InsertRfid(nRfid, nDeviceId)
	if err != nil {
		log.Error("insert rfid err,", err.Error())
	}
}

func byteToFloat(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func resetSignal(conn net.Conn) bool {
	bReturn := false
	//第十个字节0xFF表示DO输出1，0x00表示输出0
	byteCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x11, 0x05, 0x00, 0x00, 0xFF, 0x00}

	_, err := conn.Write(byteCommand)
	if err != nil {
		log.Error("reset Box signal Write err,", err.Error())
		return false
	}
	resByte := make([]byte, 25)
	_, err = conn.Read(resByte)
	fmt.Println(resByte)
	if err != nil {
		log.Error("reset Box signal Read err,", err.Error())
		return false
	}

	byteCommand[10] = 0
	_, err = conn.Write(byteCommand)
	if err != nil {
		log.Error("reset Box signal Write err,", err.Error())
		return false
	}
	_, err = conn.Read(resByte)
	fmt.Println(resByte)
	if err != nil {
		log.Error("reset Box signal Read err,", err.Error())
		return false
	}
	if !checkFinishSignal(conn) {
		bReturn = true
	}
	return bReturn
}

// func holdFurClient(conn net.Conn, ntype int) {
// 	defer util.ReturnError()
// 	defer conn.Close()
// 	//diCommand := []byte{0x01, 0x01, 0x00, 0x00, 0x00, 0x08, 0x3d, 0xcc}
// 	//11 05 00 00 FF 00 8E AA
// 	//diCommand := []byte{0x11, 0x05, 0x00, 0x00, 0xff, 0x00, 0x8e, 0xaa}
// 	//00 00 00 00 00 06 01 06 00 01 00 5A

// 	//furtype==1 重量炉；furtype==0 保温炉

// 	if ntype == 1 {

// 		conCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x01, 0x03, 0x00, 0x01, 0x00, 0x02}
// 	} else {
// 		conCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x11, 0x03, 0x00, 0x60, 0x00, 0x01}
// 	}
// 	for {
// 		// 发送数据
// 		writeLen, err := conn.Write(conCommand)
// 		if err != nil {
// 			log.Error("Write err is ", err.Error())
// 			return
// 		}
// 		fmt.Println(writeLen)
// 		resByteHead := make([]byte, 8)
// 		_, err = conn.Read(resByteHead)
// 		if err != nil {
// 			log.Error("read error:", err.Error())
// 			return
// 		}
// 		bodyLen := resByteHead[7]
// 		resBody := make([]byte, bodyLen)
// 		//parse()
// 		log.Info("resbody:%v", resBody)

// 		if ntype == 1 {
// 			nWeight := resBody[9]*256 + resBody[10]
// 			nTemp := resBody[11]*256 + resBody[12]
// 		} else {
// 			nWeight := ((resBody[9]*256+resBody[10])/1000.0 - 4) * 4000.0 / 16.0
// 			nTemp := 0
// 		}
// 		//采集数据异常，重新采集
// 		if nWeight > 6000 || nTemp > 1000 {
// 			Sleep(1000)
// 			continue
// 		}
// 		//数据库插入类型，重量，温度

// 		//1分钟采集一次
// 		Sleep(60 * 1000)
// 	}
// }
