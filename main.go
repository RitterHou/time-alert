package main

import (
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	alertTimePoint int
	disabledHours  []int
)

func init() {
	initLog() // 初始化日志设置

	var err error
	conf := getConf()
	alertTimePoint = 30
	if val, ok := conf["alert_time_point"]; ok {
		alertTimePoint, err = strconv.Atoi(val)
		if err != nil {
			log.Fatalln(err)
		}
	}
	disabledHours = make([]int, 0)
	if val, ok := conf["disabled_hours"]; ok {
		for _, v := range strings.Split(val, ",") {
			disabledHour, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalln(err)
			}
			disabledHours = append(disabledHours, disabledHour)
		}
	}
}

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

	go func() {
		systray.AddMenuItem("时间触发点："+strconv.Itoa(alertTimePoint), "").Disable()
		autoStartMenu := systray.AddMenuItem("开机自动启动", "Auto Start")
		if _, err := os.Stat(link); !os.IsNotExist(err) {
			autoStartMenu.Check()
			updateShortcut() // 把快捷方式指向当前可执行文件的路径，防止因移动文件而产生错误
		}
		settingsMenu := systray.AddMenuItem("编辑配置文件", "Settings")
		aboutMenu := systray.AddMenuItem("关于", "About")
		systray.AddSeparator()
		quitMenu := systray.AddMenuItem("退出", "Quit Time Alert")

		for {
			select {
			case <-autoStartMenu.ClickedCh:
				if autoStartMenu.Checked() {
					autoStartMenu.Uncheck()
					removeShortcut()
				} else {
					autoStartMenu.Check()
					makeShortcut()
				}
			case <-settingsMenu.ClickedCh:
				err := open.Run(path.Join(rootDir, fileName))
				if err != nil {
					log.Fatal(err)
				}
			case <-aboutMenu.ClickedCh:
				err := open.Run("https://github.com/RitterHou/time-alert")
				if err != nil {
					log.Fatal(err)
				}
			case <-quitMenu.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func main() {
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
