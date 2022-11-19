package socket

import (
	"fmt"
	"kmpi-go/config"
	"kmpi-go/dao"
	"kmpi-go/log"
	"kmpi-go/util"
	"net"
	"strconv"
	"testing"
	"time"
)

const (
	DO_NOTHING = iota
	WARM_VESSEL
	ADD_AL
	MOVE_TO_DEGAS
	MOVE_TO_FUR
	OPEN_FILL_COVER
	START_ADD_AL
	CHECK_ADD_FINISH
	CLOSE_FILL_COVER
	MOVE_TO_ORIGENAL_STATION
)

const (
	_ = iota
	OPERATE_OPEN
	OPERATE_CLOSE
)

const SMELTER_NODE=1
const SMELTER_PRE_NODE=2

const STOVE_PRE_NODE =5
const STOVE_NODE = 6

const FURNACE_PRE_NODE=9
const FURNACE_NODE=10

const DEGAS_PRE_NODE=7
const DEGAS_NODE=8
func TestWriteExcel(t *testing.T) {
	conn, err := net.Dial("tcp", "192.168.10.113:4196")
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return
	}

	defer conn.Close()

}

func AgvDispatch(conn net.Conn, chVSend chan int, chVRecv chan int, chBoxSend chan int, chBoxRecv chan int, chF1Send chan int, chF1Recv chan int, chF2Send chan int, chF2Recv chan int,ip string) {
	defer util.ReturnError()
	defer conn.Close()
	readCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x01, 0x03, 0x00, 0x00, 0x00, 0xA}
	//11 05 00 00 FF 00 8E AA
	//doCommand := []byte{0x11, 0x05, 0x00, 0x00, 0xff, 0x00, 0x8e, 0xaa}
	//00 00 00 00 00 06 01 06 00 01 00 5A

	//furtype==1 重量炉；furtype==0 保温炉
//	conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
//	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	nCmd := 0
	bAgvBusy := false
	nDesPort := 0
	nStep := 1
	var chFurS chan int
	var chFurR chan int
	var chBoxSignal int
	nLastPort := 0
	bMutex := false
	nCmdType := 0
	nFurPort := 1 //待机位
	sFurLiftHeight, _ := config.GetValue("parameter", "furnaceLiftHeight")
	SmelterLiftHeight, _ := config.GetValue("parameter", "smelterLiftHeight")
	nFurLiftHeight, _ := strconv.Atoi(sFurLiftHeight)
	nSmelterLiftHeight, _ := strconv.Atoi(SmelterLiftHeight)
	for {

		//获取AGV信息
		_, err := conn.Write(readCommand)
		if err != nil {
			log.Error("Write agv err", err.Error())
		}
		resByte := make([]byte, 30)
		_, err = conn.Read(resByte)
		bTaskOver := false //AGV任务完成标志
		nCurPort := 0
		nWorkHeight := 15
		if err == nil {

			if resByte[10] != 0 {
				bTaskOver = true
			}
			nCurPort = int(uint16(resByte[11])*256 + uint16(resByte[12])) //AGV当前站点
			if nLastPort != nCurPort {
				nLastPort = nCurPort
				go deviceRfidInsertSql(nCurPort, 1)
			}
			nWorkHeight = int(resByte[16]) //AGV目标点执行work高度
			fmt.Println("agv bTaskOver curport height", bTaskOver, nCurPort, nWorkHeight)
		} else {
			log.Error("Read agv err", err.Error())
			fmt.Println("Read agv err")
			conn = reConnectAgv(ip)
		}
		//检查AGV上动作完成按钮信号
		bFinishSignal := false
		if checkFinishSignal(conn) {
			fmt.Println("agv finish button push")
			bFinishSignal = true
		}

		//检查是否有呼叫信号

		select {
		case chBoxSignal = <-chBoxRecv:
		fmt.Println("recv signalbox cmd",nCmdType)
		default:
		}
		if chBoxSignal != 0 && !bMutex {
			if chBoxSignal > 0 && chBoxSignal < 10 {
				nCmdType = chBoxSignal
				chBoxSend <- 1
			}

			if chBoxSignal == 11 {
			fmt.Println("recv signalbox cmd",nCmdType)
				//chBoxSignal:1==烤包，2==熔炉，4==除气，8==待机位
				switch nCmdType {
				case 1:
					nCmd = WARM_VESSEL
				case 2:
					nCmd = ADD_AL
				case 4:
					nCmd = MOVE_TO_DEGAS
				case 8:
					nCmd = MOVE_TO_ORIGENAL_STATION
				default:
					nCmd = MOVE_TO_ORIGENAL_STATION
				}
				nStep = 1
				bAgvBusy = true

				chBoxSignal = 0
			}
			if chBoxSignal == 12 {
				chBoxSend <- 1
				log.Error("reset box signal fail")
			}
		}

		//检查是否有缺铝信号
		if !bAgvBusy {
			nWeight, _ := checkFurnaceWeight("5")
			if nWeight < 300 {
				nCmd = ADD_AL
				nStep = 1
				chFurS = chF1Send
				chFurR = chF1Recv
				nFurPort = 10
			}
		}

		bOK := false
		//执行命令
		switch nCmd {
		case WARM_VESSEL: //烤包
			bAgvBusy = true
			bMutex = true
			switch nStep {
			case 1: //烤包开盖动作完成确认

				if bFinishSignal {
					bOK = resetFinishSignal(conn)

				} else {
					nDesPort = STOVE_PRE_NODE
					agvGoToPort(nDesPort, conn, 15)
				}
				if bOK {
					nDesPort = STOVE_NODE
					agvGoToPort(nDesPort, conn, 15)
					nStep = 2
					bOK = false
					bFinishSignal = false
				}

			case 2: //烤包完成确认
				if bFinishSignal {
					bOK = resetFinishSignal(conn)
				}
				if bOK {
					nDesPort = 8
					agvGoToPort(nDesPort, conn, 0)
					nStep = 3
					bOK = false
				}

			case 3: //烤包关盖动作完成确认
				if bFinishSignal {
					bOK = resetFinishSignal(conn)
				}
				if bOK {
					nCmd = ADD_AL
					nStep = 1
					bOK = false
					bMutex = false
				}

			}

		case ADD_AL: //加铝
			bAgvBusy = true
			bMutex = true
			switch nStep {
			case 1: //开盖完成

				if bFinishSignal {
					bOK = resetFinishSignal(conn)
				} else {
					nDesPort = SMELTER_PRE_NODE
					agvGoToPort(nDesPort, conn, 0)
				}
				if bOK {
					nDesPort = SMELTER_NODE
					agvGoToPort(nDesPort, conn, nSmelterLiftHeight)
					nStep = 2
					bOK = false
				}

			case 2: //动作完成
				if bFinishSignal {
					bOK = resetFinishSignal(conn)
				}
				if bOK {
					nDesPort = SMELTER_PRE_NODE
					agvGoToPort(nDesPort, conn, 0)
					nStep = 3
					bOK = false
				}
			case 3: //关盖完成
				if bFinishSignal {
					bOK = resetFinishSignal(conn)
				}
				if bOK {
					nCmd = MOVE_TO_DEGAS
					nStep = 1
					bOK = false
					bMutex = false
				}

			}

		case MOVE_TO_DEGAS: //除气
			bAgvBusy = true
			bMutex = true
			switch nStep {
			case 1: //开盖完成

				if bFinishSignal {
					bOK = resetFinishSignal(conn)
				} else {
					nDesPort = DEGAS_PRE_NODE
					agvGoToPort(nDesPort, conn, 0)
				}
				if bOK {
					nDesPort = DEGAS_NODE
					agvGoToPort(nDesPort, conn, 0)
					nStep = 2
					bOK = false
				}

			case 2: //动作完成
				if bFinishSignal {
					bOK = resetFinishSignal(conn)
				}
				if bOK {
					nDesPort = DEGAS_PRE_NODE
					agvGoToPort(nDesPort, conn, 0)
					nStep = 3
					bOK = false
				}
			case 3: //关盖完成
				if bFinishSignal {
					bOK = resetFinishSignal(conn)
				}
				if bOK {
					nCmd = MOVE_TO_FUR
					nStep = 1
					bOK = false
					bMutex = false
				}

			}

		case MOVE_TO_FUR: //转运

			if nDesPort == 1 {
				nCmd = 8
			} else {
				nDesPort = nFurPort
				bSuccess := agvGoToPort(nDesPort, conn, 0)
				if bSuccess {
					nCmd = OPEN_FILL_COVER
					nStep = 1
				}
			}

		case OPEN_FILL_COVER: //开盖

			switch nStep {
			case 1:
				chFurS <- OPERATE_OPEN
				nStep = 2
			case 2:
				select {
				case <-chFurR:
					nStep = 3
				default:
				}
			case 3:
				if nDesPort == nCurPort {
					nDesPort = nFurPort + 1
					bSuccess := agvGoToPort(nDesPort, conn, nFurLiftHeight)
					if bSuccess {
						nCmd = START_ADD_AL
						nStep = 1
					}
				}
			}
		case START_ADD_AL:
			//开始加铝
			bMutex = true
			if nDesPort == nCurPort && bTaskOver && nFurLiftHeight == int(nWorkHeight) {
				chVSend <- 1
				nCmd = CHECK_ADD_FINISH
			}
		case CHECK_ADD_FINISH:
			//检查加铝完成信号

			select {
			case chData := <-chVRecv:
				if chData == 1 {
					nDesPort = 10
					bSuccess := agvGoToPort(nDesPort, conn, 0)
					if bSuccess {
						nCmd = CLOSE_FILL_COVER
						bMutex = false
					}
				}
			default:
			}

		case CLOSE_FILL_COVER: //关盖

			switch nStep {
			case 1:
				chFurS <- OPERATE_CLOSE
				nStep = 2
			case 2:
				select {
				case <-chFurR:
					nStep = 3
				default:
				}
			case 3:
				if nDesPort == nCurPort {
					bSuccess := agvGoToPort(nDesPort, conn, 0)
					if bSuccess {
						nCmd = MOVE_TO_ORIGENAL_STATION
						nStep = 1
					}
				}
			}
		case MOVE_TO_ORIGENAL_STATION: //完成回原点
			nDesPort = 1
			bSuccess := agvGoToPort(nDesPort, conn, 0)
			if bSuccess {
				nCmd = DO_NOTHING
				bAgvBusy = false
			}
		default:
			bAgvBusy = false
			bMutex = false
		}

		time.Sleep(2 * time.Second)

	}
}

func checkFurnaceWeight(furnaceId string) (int, error) {
	var res dao.AllDeviceDataDao
	err := dao.CronDb.Where("device_id=?", furnaceId).First(&res).Error
	return res.Weight, err
}

func agvGoToPort(nPort int, conn net.Conn, nHeight int) bool {
	

	byteCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x01, 0x10, 0x00, 0x02, 0x00, 0x02, 0x04, 0x00, 0x00, 0x00, 0x00}
	byteCommand[14] = uint8(nPort)   //目标站点
	byteCommand[15] = uint8(nHeight) //执行高度
	byteCommand[16] = 1              //允许执行
	
	_, err := conn.Write(byteCommand)
	if err != nil {
		log.Error("agvGoToPort Write err,", err.Error())
		return false
	}
	resByte := make([]byte, 25)
	_, err = conn.Read(resByte)
	fmt.Println(resByte)
	if err != nil {
		log.Error("agvGoToPort Read err,", err.Error())
		return false
	}
	fmt.Println("command agv go to port",nPort)
	return true
}

func resetFinishSignal(conn net.Conn) bool {
	bReturn := false
	//第十个字节0xFF表示DO输出1，0x00表示输出0
	byteCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x11, 0x05, 0x00, 0x00, 0xFF, 0x00}

	_, err := conn.Write(byteCommand)
	if err != nil {
		log.Error("reset Box signal Write err,", err.Error())
		return false
	}
	resByte := make([]byte, 12)
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

func checkFinishSignal(conn net.Conn) bool {
	bReturn := false
	byteCommand := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x11, 0x02, 0x00, 0x20, 0x00, 0x04}

	_, err := conn.Write(byteCommand)
	if err != nil {
		log.Error("check finish signal Write err,", err.Error())
		return false
	}
	resByte := make([]byte, 10)
	_, err = conn.Read(resByte)

	if err != nil {
		log.Error("check finish signal Read err,", err.Error())
		return false
	}
	if resByte[9] != 0 {
		bReturn = true
	}
	return bReturn
}

func reConnectAgv(ip string) net.Conn{
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Error("connect agv" + ip + " fail")
		time.Sleep(time.Second * 5)
		return conn
	}
	conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	return conn
}