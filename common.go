package main

import (
	"bytes"
	"encoding/base64"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"time"
)

var (
	icon      []byte
	blackIcon []byte
	current   []byte
	clock     []byte
	hour      []byte
	minute    []byte
	numbers   = make([][]byte, 11)

	app     string
	link    string
	rootDir = ""
	logName = "TimeAlert.log"
)

func init() {
	var base64ToByteArray = func(value string) []byte {
		data, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			log.Fatalln(err)
		}
		return data
	}

	icon = base64ToByteArray(Icon)
	blackIcon = base64ToByteArray(BlackIcon)
	current = base64ToByteArray(CurrentBase64)
	clock = base64ToByteArray(Clock)
	hour = base64ToByteArray(HourBase64)
	minute = base64ToByteArray(MinuteBase64)
	numbers[0] = base64ToByteArray(Num0Base64)
	numbers[1] = base64ToByteArray(Num1Base64)
	numbers[2] = base64ToByteArray(Num2Base64)
	numbers[3] = base64ToByteArray(Num3Base64)
	numbers[4] = base64ToByteArray(Num4Base64)
	numbers[5] = base64ToByteArray(Num5Base64)
	numbers[6] = base64ToByteArray(Num6Base64)
	numbers[7] = base64ToByteArray(Num7Base64)
	numbers[8] = base64ToByteArray(Num8Base64)
	numbers[9] = base64ToByteArray(Num9Base64)
	numbers[10] = base64ToByteArray(Num10Base64)

	app = os.Args[0] // 可执行文件的路径
	userInfo, err := user.Current()
	if nil != err {
		log.Fatalln(err)
	}
	// 快捷方式的路径
	link = path.Join(userInfo.HomeDir, LinkSuffix)

	rootDir = path.Join(userInfo.HomeDir, ".TimeAlert")
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		err = os.Mkdir(rootDir, os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// 把数字格式化为对应的声音
func format(num int) [][]byte {
	if num <= 10 {
		return [][]byte{numbers[num]}
	} else if num < 20 {
		return [][]byte{numbers[10], numbers[num%10]}
	} else {
		ten := num / 10
		one := num % 10
		if one == 0 {
			return [][]byte{numbers[ten], numbers[10]}
		} else {
			return [][]byte{numbers[ten], numbers[10], numbers[one]}
		}
	}
}

// 播放声音
func play(data []byte) {
	var err error

	s, format, err := wav.Decode(ioutil.NopCloser(bytes.NewReader(data)))
	if err != nil {
		log.Fatalln(err)
	}

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatalln(err)
	}
	playing := make(chan struct{})
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(playing)
	})))
	<-playing // 阻塞直到声音播放结束
}

// int是否包含在slice中
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// 创建快捷方式
func createShortcut(source string, target string) error {
	var err error
	err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	if err != nil {
		return err
	}
	defer ole.CoUninitialize()
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()
	wShell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wShell.Release()
	cs, err := oleutil.CallMethod(wShell, "CreateShortcut", target)
	if err != nil {
		return err
	}
	iDispatch := cs.ToIDispatch()
	_, err = oleutil.PutProperty(iDispatch, "TargetPath", source)
	if err != nil {
		return err
	}
	_, err = oleutil.CallMethod(iDispatch, "Save")
	if err != nil {
		return err
	}
	return nil
}

// 创建快捷方式
func makeShortcut() {
	err := createShortcut(app, link)
	if err != nil {
		log.Fatalln(err)
	}
}

// 删除快捷方式
func removeShortcut() {
	err := os.Remove(link)
	if err != nil {
		log.Fatalln(err)
	}
}

// 更新快捷方式
func updateShortcut() {
	removeShortcut()
	makeShortcut()
}

// 设置log的相关属性
func initLog() {
	logFile, _ := os.OpenFile(path.Join(rootDir, logName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	os.Stdout = logFile
	os.Stderr = logFile
}

// 检查激活状态
func checkActive() bool {
	inactiveFile := path.Join(rootDir, ".inactive")
	if _, err := os.Stat(inactiveFile); os.IsNotExist(err) {
		// 不存在，则激活
		return true
	} else {
		return false
	}
}

// 设置为未激活
func setInactive() {
	inactiveFile := path.Join(rootDir, ".inactive")
	f, err := os.OpenFile(inactiveFile, os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// 设置为激活
func setActive() {
	inactiveFile := path.Join(rootDir, ".inactive")
	err := os.Remove(inactiveFile)
	if err != nil {
		log.Fatal(err)
	}
}
