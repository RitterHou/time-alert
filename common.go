package main

import (
	"bytes"
	"encoding/base64"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"io/ioutil"
	"time"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"os/user"
	"log"
	"path"
	"os"
)

var (
	icon    []byte
	current []byte
	clock   []byte
	hour    []byte
	minute  []byte
	numbers = make([][]byte, 11)

	app  string
	link string
)

func init() {
	var base64ToByteArray = func(value string) []byte {
		data, _ := base64.StdEncoding.DecodeString(value)
		return data
	}

	icon = base64ToByteArray(Icon)
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
	s, format, _ := wav.Decode(ioutil.NopCloser(bytes.NewReader(data)))

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
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
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
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
	oleutil.PutProperty(iDispatch, "TargetPath", source)
	oleutil.CallMethod(iDispatch, "Save")
	return nil
}

func makeShortcut() {
	createShortcut(app, link)
}

func removeShortcut() {
	os.Remove(link)
}
