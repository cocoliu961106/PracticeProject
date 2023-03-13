package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	// go run code-user/main.go
	cmd := exec.Command("go", "run", "code/code-user/main.go") // 相当于命令行执行命令
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr // 用户代码若编译错误，错误信息会在stderr里 等同于 cmd.Stderr = new(bytes.Buffer)
	cmd.Stdout = &out
	stdinPipe, err := cmd.StdinPipe() // 用户执行代码的标准输入
	if err != nil {
		log.Fatalln(err)
	}
	io.WriteString(stdinPipe, "23 11\n") // 用户执行代码中标准输入的值

	// 运行输入案例，拿到输出结果和标准的输出结果进行比对
	if err = cmd.Run(); err != nil {
		log.Fatalln(err, stderr.String())
	}
	fmt.Println("result: ", out.String())

	println("the result is: ", out.String() == "34\n")
}
