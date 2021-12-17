package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func GetUserName() (username string) {
	fmt.Println("Please input repository creator username:")
	_, err := fmt.Scanf("%s", &username)
	if err != nil {
		log.Fatal()
	}
	return
}

func GetRepoName() (repoName string) {
	fmt.Println("Please input your git repository name:")
	_, err := fmt.Scanf("%s", &repoName)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func GetToken() (token string) {
	fmt.Println("Please input your github token:")
	_, err := fmt.Scanf("%s", &token)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func GetRepourl(reponame, username string) (repourl string) {
	repourl = "https://api.github.com/repos/" + username + "/" + reponame + "/issues"
	return
}

func GetMessage() (message string) {
	var cmd1, cmd2, cmd3, cmd4 *exec.Cmd
	var tempFile string
	sysType := runtime.GOOS
	if sysType == "linux" {
		cmd1 = exec.Command("touch", tempFile)
		cmd2 = exec.Command("gedit", tempFile)
		cmd3 = exec.Command("cat", tempFile)
		cmd4 = exec.Command("rm", tempFile)
	} else if sysType == "windows" {
		cmd1 = exec.Command("fsutil", "file", "createnew", "temp.txt", "0")
		cmd2 = exec.Command("notepad", "temp.txt")
		cmd3 = exec.Command("cat", "temp.txt")
		cmd4 = exec.Command("rm", "temp.txt")
	}
	err := cmd1.Start()
	if err != nil {
		log.Fatal(err)
	}
	// cmd1作为第一条命令，前面没有命令需要等待
	// err = cmd1.Wait()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd2.Start()
	if err != nil {
		log.Fatal(err)
	}
	// cmd2必须等待cmd1创建文件命令完成才能开始执行
	err = cmd2.Wait()
	if err != nil {
		log.Fatal(err)
	}
	var buffer bytes.Buffer
	cmd3.Stdout = &buffer
	err = cmd3.Start()
	if err != nil {
		log.Fatal(err)
	}
	// cmd3需要等待用户向notepad窗口中输入信息并保存，所以必须等待cmd2完成才能前进到下一步
	err = cmd3.Wait()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd4.Start()
	if err != nil {
		log.Fatal(err)
	}
	// cmd4是删除文件的命令，在cmd3中我们已经将文件内容输出到buffer中，所以无需等待
	//cmd4.Wait()
	message = buffer.String()
	return
}
