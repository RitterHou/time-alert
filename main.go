package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"os"
	"time"
)

func say(h int, m int) {
	go func() {
		play(current)

		for _, b := range format(h) {
			play(b)
		}
		play(hour)

		if m == 0 {
			play(clock) // 整点
		} else {
			if m < 10 {
				play(numbers[0])
			}
			for _, b := range format(m) {
				play(b)
			}
			play(minute)
		}
	}()
}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTooltip("Time Alert")

	autoStart := false // 开启自动启动
	go func() {
		autoStartMenu := systray.AddMenuItem("开机自动启动", "Auto Start")
		if autoStart {
			autoStartMenu.Check()
		}
		aboutMenu := systray.AddMenuItem("关于", "About")
		systray.AddSeparator()
		quitMenu := systray.AddMenuItem("退出", "Quit Time Alert")

		for {
			select {
			case <-autoStartMenu.ClickedCh:
				if autoStartMenu.Checked() {
					autoStartMenu.Uncheck()
				} else {
					autoStartMenu.Check()
				}
			case <-aboutMenu.ClickedCh:
				open.Run("https://www.google.com")
			case <-quitMenu.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func main() {
	file := os.Args[0]
	fmt.Println(file)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for t := range ticker.C {
			say(t.Hour(), t.Minute())
		}
	}()

	systray.Run(onReady, func() {})
}
