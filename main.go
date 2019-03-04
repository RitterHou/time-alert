package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"os"
	"strconv"
	"strings"
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
				open.Run("https://github.com/RitterHou/time-alert")
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

	conf := getConf()
	alertTimePoint := 30
	if val, ok := conf["alert_time_point"]; ok {
		alertTimePoint, _ = strconv.Atoi(val)
	}
	disabledHours := make([]int, 0)
	if val, ok := conf["disabled_hours"]; ok {
		for _, v := range strings.Split(val, ",") {
			disabledHour, _ := strconv.Atoi(v)
			disabledHours = append(disabledHours, disabledHour)
		}
	}

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		currentMinute := time.Now().Minute()
		for t := range ticker.C {
			h := t.Hour()
			m := t.Minute()
			// 如果相等则意味着还在这一分钟没有变，则不需要任何处理
			if m != currentMinute {
				if m%alertTimePoint == 0 && !contains(disabledHours, h) {
					say(h, m)
				}
			}
			currentMinute = m
		}
	}()

	systray.Run(onReady, func() {})
}
