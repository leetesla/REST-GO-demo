package robot

//Robot interface
type Robot interface {
	Stop()
	Action() error
}
