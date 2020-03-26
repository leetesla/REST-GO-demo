package robot

//运行的机器人保存在这里
import (
	"errors"
	"sync"
)

//RobotMap save running or stop robot
var RobotMap map[int]Robot
var RobotMapLock sync.Mutex

func init() {
	RobotMap = make(map[int]Robot)
}

// Register 注册robot
func Register(Id int, rbt Robot) error {
	RobotMapLock.Lock()
	defer RobotMapLock.Unlock()

	_, ok := RobotMap[Id]
	if ok {
		return errors.New("robot exist")
	}
	RobotMap[Id] = rbt
	return nil
}

// Remove 移除robot
// func Remove(Id int) {
// 	RobotMapLock.Lock()
// 	defer RobotMapLock.Unlock()
// 	delete(RobotMap, Id)
// }

// Stop 停止robot
func Stop(Id int) {
	RobotMapLock.Lock()
	defer RobotMapLock.Unlock()
	rbt, ok := RobotMap[Id]
	if !ok {
		return
	}
	rbt.Stop()
	delete(RobotMap, Id)
	return
}

func Action(Id int) {
	RobotMapLock.Lock()
	defer RobotMapLock.Unlock()
	rbt, ok := RobotMap[Id]
	if !ok {
		return
	}
	rbt.Action()
	return
}
