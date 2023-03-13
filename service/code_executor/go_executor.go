package code_executor

import (
	"PracticeProject/helper"
	"PracticeProject/models"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

func GoExecuteTestCases(problemBasic *models.ProblemBasic, submitBasic *models.SubmitBasic, path string) string {
	// 答案错误的channel
	WrongAnswerCh := make(chan int)
	// 超内存的channel
	OOMCh := make(chan int)
	// 编译错误的channel
	CompileErrCh := make(chan int)
	// 答案正确的channel
	AcceptCh := make(chan int)
	// 非法代码的channel
	IllegalCodeCh := make(chan struct{})

	passCount := 0 // 通过测试案例的个数
	var lock sync.Mutex
	var msg string // 提示信息

	// 检查代码的合法性
	v, err := helper.CheckGoCodeValid(path)
	if err != nil {
		//c.JSON(http.StatusOK, gin.H{
		//	"code": -1,
		//	"msg":  "Code Check Error:" + err.Error(),
		//})
		msg = "Code Check Error:" + err.Error()
		return msg
	}
	if !v {
		go func() {
			IllegalCodeCh <- struct{}{}
		}()
	} else {
		for _, testCase := range problemBasic.TestCases {
			// 每一个testcase用一个协程执行
			go func(testCase *models.TestCase) {
				cmd := exec.Command("go", "run", path)
				var out, stderr bytes.Buffer
				cmd.Stderr = &stderr
				cmd.Stdout = &out
				stdinPipe, err := cmd.StdinPipe()
				if err != nil {
					log.Fatalln(err)
				}
				io.WriteString(stdinPipe, testCase.Input+"\n") // 将测试案例的input作为用户代码的输入

				var bm runtime.MemStats
				runtime.ReadMemStats(&bm)
				if err := cmd.Run(); err != nil {
					log.Println(err, stderr.String())
					msg = stderr.String()
					CompileErrCh <- 1
					return
				}
				var em runtime.MemStats
				runtime.ReadMemStats(&em)

				// 答案错误
				fmt.Println("testCase.Output: " + testCase.Output)
				fmt.Println("out.String(): " + out.String())
				if testCase.Output != out.String() {
					WrongAnswerCh <- 1
					return
				}
				// 运行超内存
				if em.Alloc/1024-(bm.Alloc/1024) > uint64(problemBasic.MaxMem) {
					OOMCh <- 1
					return
				}
				lock.Lock()
				passCount++
				if passCount == len(problemBasic.TestCases) {
					AcceptCh <- 1
				}
				lock.Unlock()
			}(testCase)
		}
	}
	select {
	case <-IllegalCodeCh:
		msg = "无效代码"
		submitBasic.Status = 6
	case <-WrongAnswerCh:
		msg = "答案错误"
		submitBasic.Status = 2
	case <-OOMCh:
		msg = "运行超内存"
		submitBasic.Status = 4
	case <-CompileErrCh:
		submitBasic.Status = 5
	case <-AcceptCh:
		msg = "答案正确"
		submitBasic.Status = 1
	case <-time.After(time.Millisecond * time.Duration(problemBasic.MaxRuntime)):
		if passCount == len(problemBasic.TestCases) {
			submitBasic.Status = 1
			msg = "答案正确"
		} else {
			submitBasic.Status = 3
			msg = "运行超时"
		}
	}
	return msg
}
