package somfy

import "fmt"

type Button uint16

const ButtonMy = 0x1
const ButtonUp = 0x2
const ButtonDown = 0x4
const ButtonProg = 0x8

func getButtonFromCommand(cmd string) (Button, error) {
	switch cmd {
	case "up":
		return ButtonUp, nil
	case "down":
		return ButtonDown, nil
	case "my":
		return ButtonMy, nil
	case "prog":
		return ButtonProg, nil
	}
	return ButtonMy, fmt.Errorf("could not map command %s", cmd)
}
