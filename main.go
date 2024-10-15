package main

import (
	"fmt"
	"os/exec"

	"github.com/MarinX/keylogger"
)

type ScreenConfig struct {
	KeyboardKey string
	Bus         string
	// sudo ddcutil --dis 1 cap --verbose => search fo Power Mode
	VcpOffCode  string
	VcpOffValue string
}

func main() {

	config := []ScreenConfig{
		{KeyboardKey: "F10", Bus: "5", VcpOffCode: "0xd6", VcpOffValue: "0x05"},
		{KeyboardKey: "F11", Bus: "8", VcpOffCode: "0xd6", VcpOffValue: "0x05"},
		{KeyboardKey: "F12", Bus: "6", VcpOffCode: "0xd6", VcpOffValue: "0x05"},
	}

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

			for _, config := range config {
				if e.KeyPress() && l_altPressed && e.KeyString() == config.KeyboardKey {
					cmd := exec.Command("ddcutil", "--bus", config.Bus, "setvcp", config.VcpOffCode, config.VcpOffValue)

					_, err := cmd.CombinedOutput()
					if err != nil {
						fmt.Print(err)
					}
				}
			}
		}
	}
}
