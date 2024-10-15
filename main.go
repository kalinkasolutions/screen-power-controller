package main

import (
	"fmt"
	"os/exec"

	"github.com/MarinX/keylogger"
)

func main() {
	screenKeyMap := make(map[string]string)

	// define keys which should turn of a screen paired with alt-left key.
	// key: keyboard-key
	// value: monitor bus (use: ddcutil detect)
	screenKeyMap["F10"] = "5"
	screenKeyMap["F11"] = "8"
	screenKeyMap["F12"] = "6"

	keyboard := keylogger.FindKeyboardDevice()

	if len(keyboard) <= 0 {
		fmt.Println("No keyboard found...you will need to provide manual input path")
		// set input if not found
		keyboard = "/dev/input/event7"
	}

	k, err := keylogger.New(keyboard)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer k.Close()

	events := k.Read()

	l_altPressed := false

	for e := range events {
		switch e.Type {
		case keylogger.EvKey:

			if e.KeyRelease() && e.KeyString() == "L_ALT" {
				l_altPressed = false
			}

			if e.KeyPress() && e.KeyString() == "L_ALT" {
				l_altPressed = true
			}

			for key, bus := range screenKeyMap {
				if e.KeyPress() && l_altPressed && e.KeyString() == key {
					cmd := exec.Command("ddcutil", "--bus", bus, "setvcp", "0xd6", "0x05")

					_, err := cmd.CombinedOutput()
					if err != nil {
						fmt.Print(err)
					}
				}
			}
		}
	}
}
