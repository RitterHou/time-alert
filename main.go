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
	timePoint      *systray.MenuItem
	activeAlert    bool
)

func init() {
	initLog() // 初始化日志设置

	log.Printf("TimeAlert starting %s\r\n", time.Now().Format("2006-01-02 15:04:05"))

	confFile := path.Join(rootDir, fileName)
	updateConf(confFile)

	confStat, err := os.Stat(confFile)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		ticket := time.NewTicker(1 * time.Second)
		for range ticket.C {
			stat, err := os.Stat(confFile)
			if err != nil {
				log.Fatal(err)
			}
			if stat.Size() != confStat.Size() || stat.ModTime() != confStat.ModTime() {
				log.Printf("Update conf %v\r\n", stat)
				updateConf(confFile) // 文件发生变化则更新配置
				timePoint.SetTitle("时间触发点：" + strconv.Itoa(alertTimePoint))
				confStat = stat
			}
		}
	}()

	activeAlert = checkActive()
}

func updateConf(confFile string) {
	var err error
	conf := getConf(confFile)
	alertTimePoint = 30
	if val, ok := conf["alert_time_point"]; ok {
		alertTimePoint, err = strconv.Atoi(val)
		if err != nil {
			log.Println(err)
		}
	}
	disabledHours = make([]int, 0)
	if val, ok := conf["disabled_hours"]; ok {
		for _, v := range strings.Split(val, ",") {
			disabledHour, err := strconv.Atoi(v)
			if err != nil {
				log.Println(err)
			}
			disabledHours = append(disabledHours, disabledHour)
		}
	}
}

func say(h int, m int) {
	if !activeAlert {
		return
	}
	go func() {
		play(current)

		if h > 12 { // 12小时制
			h = h - 12
		}

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
	go func() {
		timeAlert := systray.AddMenuItem("启用 TimeAlert", "")
		if activeAlert {
			systray.SetIcon(icon)
			systray.SetTooltip("Time Alert 已启用")
			timeAlert.Check()
		} else {
			systray.SetIcon(blackIcon)
			systray.SetTooltip("Time Alert 未启用")
		}

		timePoint = systray.AddMenuItem("时间触发点："+strconv.Itoa(alertTimePoint), "")
		timePoint.Disable()

		autoStartMenu := systray.AddMenuItem("开机启动", "Auto Start")
		if _, err := os.Stat(link); !os.IsNotExist(err) {
			autoStartMenu.Check()
			updateShortcut() // 把快捷方式指向当前可执行文件的路径，防止因移动文件而产生错误
		}
		settingsMenu := systray.AddMenuItem("编辑本地 Settings 文件", "Settings")
		logMenu := systray.AddMenuItem("显示日志...", "log")
		aboutMenu := systray.AddMenuItem("关于...", "About")
		systray.AddSeparator()
		quitMenu := systray.AddMenuItem("退出", "Quit Time Alert")

		for {
			select {
			case <-timeAlert.ClickedCh:
				if timeAlert.Checked() {
					timeAlert.Uncheck()
					activeAlert = false
					setInactive()
					systray.SetIcon(blackIcon)
					systray.SetTooltip("Time Alert 未启用")
				} else {
					timeAlert.Check()
					activeAlert = true
					setActive()
					systray.SetIcon(icon)
					systray.SetTooltip("Time Alert 已启用")
				}
			case <-autoStartMenu.ClickedCh:
				if autoStartMenu.Checked() {
					autoStartMenu.Uncheck()
					removeShortcut()
				} else {
					autoStartMenu.Check()
					makeShortcut()
				}
			case <-logMenu.ClickedCh:
				err := open.Run(path.Join(rootDir, logName))
				if err != nil {
					log.Fatal(err)
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
