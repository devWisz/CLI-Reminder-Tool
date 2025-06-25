package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

const (
	markName  = "CLI reminder"
	markValue = "1"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <hh:mm> <message>\n", os.Args[0])
		os.Exit(1)
	}

	now := time.Now()

	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	t, err := w.Parse(os.Args[1], now)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		os.Exit(2)
	}

	if t == nil {
		fmt.Println("Could not parse time:", os.Args[1])
		os.Exit(3)
	}

	if now.After(t.Time) {
		fmt.Println("Time is in the past:", t.Time)
		os.Exit(3)
	}

	diff := t.Time.Sub(now)

	if os.Getenv(markName) == markValue {
		time.Sleep(diff)
		err = beeep.Alert("CLI Reminder", strings.Join(os.Args[2:], " "), "assets/information.png")
		if err != nil {
			fmt.Println("Error sending alert:", err)
			os.Exit(4)
		}
	} else {
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", markName, markValue))

		if err := cmd.Start(); err != nil {
			fmt.Println("Error starting command:", err)
			os.Exit(5)
		}

		fmt.Println("Reminder will come after:", diff.Round(time.Second))
		os.Exit(0)
	}
}
