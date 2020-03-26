package main

import (
	"encoding/json"
	"flag"
	"github.com/golang/glog"
	"gtask/cmd/http"
	"gtask/internal/session"
	"gtask/internal/task"
	"gtask/pkg/restful"
	"io/ioutil"
	"os"
	"runtime"
)

type Params struct {
	cfgPath string
}

type Config struct {
	Port      uint32 `json:"listen_port"`
	SecretKey string `json:"secret_key"`
}

func readFiles(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func main() {
	// 获取命令行参数
	params := &Params{}
	flag.StringVar(&params.cfgPath, "cfg", "./config/config.json", "config file. default './config/config.json'")
	// 初始化日志信息
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Parse()
	defer func() {
		glog.Flush()
	}()
	runtime.GOMAXPROCS(runtime.NumCPU())

	content, err := readFiles(params.cfgPath)
	if err != nil {
		glog.Fatal(err)
	}
	cfg := &Config{}
	err = json.Unmarshal(content, cfg)
	if err != nil {
		glog.Fatal(err)
	}

	session.Ins().Init(cfg.SecretKey)
	// 随机生成一个token
	session.Ins().GetToken(cfg.SecretKey)
	task.Ins().Init()

	restfulServer := restful.NewRestful()
	restfulServer.SetDefOpt(&restful.SchedulerOpt{UseCORS: true})

	//restfulServer.Post("/task/start", http.StartTaskHandler)
	//restfulServer.Post("/task/stop", http.StopTaskHandler)
	restfulServer.Post("/task/job", http.CreateJobHandler)
	restfulServer.Delete("/task/job", http.DeleteJobHandler)
	restfulServer.Post("/task/job/run", http.RunJobHandler)
	restfulServer.Post("/task/job/processor", http.CreateProcessorHandler)
	restfulServer.Patch("/task", http.UpdateTaskHandler)
	restfulServer.Patch("/task/job", http.UpdateJobHandler)
	restfulServer.Get("/task/job", http.GetJobHandler)

	restfulServer.Get("/session/token", http.GetToken)

	glog.Infof("gtask running %d", cfg.Port)
	task.Ins().Start()
	restfulServer.Start(cfg.Port)
	task.Ins().Stop()
	task.Ins().Wait()
}
