package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	logger "github.com/sirupsen/logrus"
)

// 钉钉json 解析
type dingMsg struct {
	At struct {
		AtMobiles []string `json:"atMobiles"`
	} `json:"at"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`

	MsgType string `json:"msgtype"`
}

// 钉钉钩子
type dingHook struct {
	apiUrl     string
	levels     []logger.Level
	atMobiles  []string
	atUserIds  []string
	appName    string
	jsonBodies chan []byte
	closeChan  chan bool
}

//- Levels 代表在哪几个级别下应用这个hook
func (dh *dingHook) Levels() []logger.Level {
	return dh.levels
}

//- Fire 代表 执行具体什么逻辑
func (dh *dingHook) Fire(e *logger.Entry) error {
	msg, _ := e.String()
	dh.DirectSend(msg)
	return nil
}

// 设置报警主体函数
func (dh *dingHook) DirectSend(msg string) {
	dm := dingMsg{
		MsgType: "text",
	}

	// 报警主体
	dm.Text.Content = fmt.Sprintf("[test log]\n[app = %s]\n"+
		"[log info:%s]", dh.appName, msg,
	)

	dm.At.AtMobiles = dh.atMobiles

	bs, err := json.Marshal(dm)
	if err != nil {
		logger.Errorf("[消息:json.Marshal 失败][error：%v][msg: %v]", err, msg)
		return
	}

	// 向 dingding token 提交请求
	res, err := http.Post(dh.apiUrl, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		logger.Errorf("[消息发送失败][error:%v][msg:%v]", err, msg)
		return
	}
	if res != nil && res.StatusCode != 200 {
		logger.Errorf("[钉钉请求报错][状态码：%v][msg:%v]", res.StatusCode, msg)
		return
	}

}

func test(dh *dingHook, wg *sync.WaitGroup) {
	defer wg.Done()
	levle := logger.InfoLevel
	logger.SetLevel(levle)
	logger.SetReportCaller(true)
	logger.SetFormatter(&logger.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.AddHook(dh)
	logger.Info("这是一个 test ")
}

func main() {
	wg := sync.WaitGroup{}
	dh := &dingHook{
		apiUrl:     "https://oapi.dingtalk.com/robot/send?access_token=xxxxxx",
		levels:     []logger.Level{logger.WarnLevel, logger.InfoLevel},
		atMobiles:  []string{"177xxxxx"},
		appName:    "容器安全",
		jsonBodies: make(chan []byte),
		closeChan:  make(chan bool),
	}
	wg.Add(1)
	go test(dh, &wg)
	wg.Wait()
}
