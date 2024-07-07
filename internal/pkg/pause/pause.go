package pause

var pause bool

func Pause() {
	pause = true
}

func Resume() {
	pause = false
}

func Paused() bool {
	return pause
}
