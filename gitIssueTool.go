package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopl.io/ch4/github"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func ShowMenu() {
	fmt.Println("Options:")
	fmt.Println("1. list public repository's issues")
	fmt.Println("2. list public or private repository's issues (use git token)")
	fmt.Println("3. create repository's issue")
	fmt.Println("4. update repository's issue")
	fmt.Println("5. close repository's issue")
	fmt.Println("6. exit process")
}

func SearchIssue(repoName []string) {
	for _, tempName := range repoName {
		fmt.Printf("Now handle repository %s Issues\n", tempName)
		result, err := github.SearchIssues(repoName[:])
		fmt.Printf("Repository %s has %d issues\n", tempName, len(result.Items))
		if err != nil {
			log.Fatal(err)
		}
		for _, item := range result.Items {
			fmt.Printf("#%-5d %9.9s %.55s %s %s\n", item.Number, item.User.Login, item.Title,
				item.State, item.CreatedAt)
		}
	}
}

func ListAppointIssue(repoUrl, username, token string) {
	var datas []Data
	// var data github.IssuesSearchResult
	resp, err := http.Get(repoUrl)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 404 {
		// fmt.Println("You don't have permission access this repository!")
		for token == "" {
			token = GetToken()
			_, _ = fmt.Scanln()
		}
		// resp.Body.Close()
		req, err := http.NewRequest("GET", repoUrl, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(username, token)
		req.Header.Add("Accept", "application/vnd.github.v3+json")
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	//使用ioutil.ReadAll函数可以将url的响应值并返回为[]byte格式，
	//该格式是多个字符组成的切片（动态数组），可以看作一个字符串
	bodyText, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyText, &datas)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range datas {
		fmt.Printf("#%-5d %9.9s %.55s %s %s\n", item.Number, item.User.LoginName,
			item.Title, item.State, item.CreateAt)
	}
}

func CreateIssue(repoUrl, username, token string) {
	fmt.Println("please input your issue's title:")
	title := GetMessage()
	time.Sleep(time.Second * 1)
	fmt.Println("please input your comment:")
	comment := GetMessage()
	// 这里先创建一个结构，该结构与提交给api的json数据结构相同
	data := PostData{Title: title, Comment: comment}
	// 使用json包中的MarshalIndent函数将结构转换为[]byte类型的json数据
	jsonData, _ := json.MarshalIndent(data, "", "    ")
	// fmt.Printf("%s", jsonData)
	// 最后使用bytes包的newreader函数将使用[]bytes类型作为参数返回一个读取器
	for token == "" {
		fmt.Println("It's not support anonymous submit issue, you need take your token:")
		token = GetToken()
	}
	_, _ = fmt.Scanln()
	req, err := http.NewRequest("POST", repoUrl, bytes.NewReader(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(username, token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	if resp.StatusCode == 201 {
		fmt.Println("Creating issue successful!")
	}
}

func UpdateIssue(repoUrl, username, token string) {
	var issueID, title, comment string
	var data PostData
	for token == "" {
		fmt.Println("It's not support anonymous edit issue, you need take your token:")
		token = GetToken()
	}
	_, _ = fmt.Scanln()
	for issueID == "" {
		fmt.Println("Please input the issue id what's you want to patch:")
		_, _ = fmt.Scanf("%s", &issueID)
	}
	_, _ = fmt.Scanln()
	repoUrl = repoUrl + "/" + issueID
	fmt.Println("Please input title:")
	title = GetMessage()
	fmt.Println("Please input comment:")
	comment = GetMessage()
	data.Title = title
	data.Comment = comment
	jsonData, _ := json.MarshalIndent(&data, "", "    ")
	req, err := http.NewRequest("PATCH", repoUrl, bytes.NewReader(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(username, token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		fmt.Println("Patch issue successful!")
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "Patch issue fail, status code: %d\t%s", resp.StatusCode, resp.Status)
	}
}

func CloseIssue(repoUrl, username, token string) {
	var issueID string
	var data lockData
	var opt int
	for token == "" {
		fmt.Println("It's not support anonymous edit issue, you need take your token:")
		token = GetToken()
	}
	_, _ = fmt.Scanln()
	for issueID == "" {
		fmt.Println("Please input the issue id what's you want to lock:")
		_, _ = fmt.Scanf("%s", &issueID)
	}
	_, _ = fmt.Scanln()
	repoUrl = repoUrl + "/" + issueID + "/lock"
	fmt.Println("Please chose the reason of lock this issue:")
	fmt.Println("1.* off-topic\n2.* too heated\n3.* resolved\n4.* spam")
	for opt == 0 {
		_, _ = fmt.Scanf("%d", &opt)
	}
	_, _ = fmt.Scanln()
	switch opt {
	case 1:
		{
			data.LockReason = "off-topic"
		}
	case 2:
		{
			data.LockReason = "too heated"
		}
	case 3:
		{
			data.LockReason = "resolved"
		}
	case 4:
		{
			data.LockReason = "spam"
		}
	}
	jsonData, err := json.MarshalIndent(&data, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("PUT", repoUrl, bytes.NewReader(jsonData))
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.SetBasicAuth(username, token)
	resp, err := http.DefaultClient.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 204 {
		fmt.Println("Lock issue successful!")
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "Lock issue failed, status code :%d\t%s", resp.StatusCode, resp.Status)
	}
}
