package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type JobClient struct {
	Host      string
	Port      int
	SecretKey string
	Token     string
	CurJob    string
}

type JobInfo struct {
	ProcessorCount int64 `json:"processorsCount"`
	State          int64 `json:"state"`
}

func (c *JobClient) Login(host string, port int, secretKey string) string {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/session/token?secretKey=%s", host, port, secretKey))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return ""
	}
	if status, ok := data["status"]; ok && int64(status.(float64)) == 0 {
		if token, ok := data["body"].(map[string]interface{})["token"]; ok {
			c.Token = token.(string)
			c.SecretKey = secretKey
			c.Host = host
			c.Port = port
			return c.Token
		}
	}
	return ""
}

func (c *JobClient) GetJob(key string) *JobInfo {
	request, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/task/job?key=%s", c.Host, c.Port, key), nil)
	request.Header.Set("token", c.Token)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	jobInfo := &JobInfo{}
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil
	}
	if status, ok := data["status"]; ok && int64(status.(float64)) == 0 {
		jobInfo.ProcessorCount = int64(data["body"].(map[string]interface{})["processorsCount"].(float64))
		jobInfo.State = int64(data["body"].(map[string]interface{})["state"].(float64))
		return jobInfo
	}
	return nil
}

func (c *JobClient) CreateJob(key string) (b bool) {
	r, err := json.Marshal(map[string]interface{}{"key": key})
	if err != nil {
		return
	}
	request, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/task/job", c.Host, c.Port), bytes.NewReader(r))
	request.Header.Set("token", c.Token)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}
	if status, ok := data["status"]; ok && int64(status.(float64)) == 0 {
		b = true
	}
	return
}

func (c *JobClient) RunJob(key string) (b bool) {
	r, err := json.Marshal(map[string]interface{}{"key": key})
	if err != nil {
		return
	}
	request, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/task/job/run", c.Host, c.Port), bytes.NewReader(r))
	request.Header.Set("token", c.Token)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}
	if status, ok := data["status"]; ok && int64(status.(float64)) == 0 {
		b = true
	}
	return
}

func (c *JobClient) DeleteJob(key string) (b bool) {
	r, err := json.Marshal(map[string]interface{}{"key": key})
	if err != nil {
		return
	}
	request, _ := http.NewRequest("DELETE", fmt.Sprintf("http://%s:%d/task/job", c.Host, c.Port), bytes.NewReader(r))
	request.Header.Set("token", c.Token)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}
	if status, ok := data["status"]; ok && int64(status.(float64)) == 0 {
		b = true
	}
	return
}

func (c *JobClient) CreateProcessor(key, code string, trigger int64, bReset, bLoop, bExit bool) (b bool) {
	r, err := json.Marshal(map[string]interface{}{
		"key":     key,
		"code":    code,
		"trigger": trigger,
		"bReset":  bReset,
		"bLoop":   bLoop,
		"bExit":   bExit,
	})
	if err != nil {
		return
	}
	request, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/task/job/processor", c.Host, c.Port), bytes.NewReader(r))
	request.Header.Set("token", c.Token)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}
	if status, ok := data["status"]; ok && int64(status.(float64)) == 0 {
		b = true
	}
	return
}

type Params struct {
	Host string
	Port int
}

func main() {
	params := &Params{}
	flag.StringVar(&params.Host, "h", "", "host")
	flag.IntVar(&params.Port, "p", 1126, "port default 1126")
	flag.Parse()
	if params.Host == "" || params.Port == 0 {
		fmt.Println("host or port is nil")
		return
	}
	client := &JobClient{}
	for i := 0; i < 3; i++ {
		fmt.Print("secretKey:")
		secretKey, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Print("\n")
		if err == nil {
			client.Token = client.Login(params.Host, params.Port, string(secretKey))
			if client.Token != "" {
				break
			}
		}

	}
	if client.Token == "" {
		glog.Info("login failed")
		return
	}
	for {
		fmt.Printf("%s > ", client.CurJob)
		inputReader := bufio.NewReader(os.Stdin)
		cmd, _ := inputReader.ReadString('\n')
		cmd = cmd[0 : len(cmd)-1]
		args := strings.Split(cmd, " ")
		if len(args) < 1 {
			fmt.Println("cmd error")
			continue
		}
		if args[0] == "use" {
			if len(args) < 2 {
				fmt.Println("cmd error")
				continue
			}

			job := client.GetJob(args[1])
			if job != nil {
				client.CurJob = args[1]
				fmt.Printf("select job [%s]\n", args[1])
			} else {
				fmt.Printf("use job [%s] failed\n", args[1])
			}
		}
		if args[0] == "run" {
			if client.CurJob == "" {
				fmt.Println("please use job")
				continue
			}
			if client.RunJob(client.CurJob) {
				fmt.Printf("run job [%s] success\n", client.CurJob)
			} else {
				fmt.Printf("run job [%s] failed\n", client.CurJob)
			}
		}
		if args[0] == "delete" {
			if client.CurJob == "" {
				fmt.Println("please use job")
				continue
			}
			client.DeleteJob(client.CurJob)
			client.CurJob = ""
		}
		if args[0] == "create" {
			if args[1] == "job" {
				if len(args) < 3 {
					fmt.Println("cmd error")
					continue
				}
				if client.CreateJob(args[2]) {
					fmt.Printf("create job [%s] success\n", args[2])
				} else {
					fmt.Printf("create job [%s] failed\n", args[2])
				}
			}
			if args[1] == "processor" {
				if len(args) < 7 {
					fmt.Println("cmd error")
					continue
				}
				if client.CurJob == "" {
					fmt.Println("please use job")
					continue
				}
				code, err := ioutil.ReadFile(args[2])
				if err != nil {
					fmt.Println("code file error")
					continue
				}
				var (
					trigger int64
					bReset  bool
					bLoop   bool
					bExit   bool
				)
				trigger, _ = strconv.ParseInt(args[3], 10, 64)
				if args[4] == "1" {
					bReset = true
				}
				if args[5] == "1" {
					bLoop = true
				}
				if args[6] == "1" {
					bExit = true
				}
				if client.CreateProcessor(client.CurJob, string(code), trigger, bReset, bLoop, bExit) {
					fmt.Printf("create processor success\n")
				} else {
					fmt.Printf("create processor failed\n")
				}
			}
		}
	}
}
