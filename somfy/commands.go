package somfy

import "fmt"

type Button uint16

const ButtonMy = 0x1
const ButtonUp = 0x2
const ButtonDown = 0x4
const ButtonProg = 0x8

const CmdMy = "my"
const CmdUp = "up"
const CmdDown = "down"
const CmdProg = "prog"

func getButtonFromCommand(cmd string) (Button, error) {
	switch cmd {
	case CmdUp:
		return ButtonUp, nil
	case CmdDown:
		return ButtonDown, nil
	case CmdMy:
		return ButtonMy, nil
	case CmdProg:
		return ButtonProg, nil
	}
	return ButtonMy, fmt.Errorf("could not map command %s", cmd)
}
