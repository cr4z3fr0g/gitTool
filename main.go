package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var repoName, userName, token, repoURL string
	var repos []string
	flag.StringVar(&repoName, "n", "", "repository name")
	flag.StringVar(&token, "t", "", "github token")
	flag.StringVar(&userName, "u", "", "repository userName")
	flag.Parse()
	var sign = true
	var option int
	for sign {
		ShowMenu()
		option = 0
		for option == 0 {
			fmt.Println("please input your option:")
			_, _ = fmt.Scanf("%d", &option)
		}
		_, _ = fmt.Scanln()
		if option != 6 {
			if repoName == "" {
				repoName = GetRepoName()
				_, _ = fmt.Scanln()
			}
			if userName == "" {
				userName = GetUserName()
				_, _ = fmt.Scanln()
			}
			repoURL = GetRepourl(repoName, userName)
			repos = append(repos, "repo:"+userName+"/"+repoName)

		}
		switch option {
		case 1:
			{
				SearchIssue(repos)
			}
		case 2:
			{
				ListAppointIssue(repoURL, repoName, token)
			}
		case 3:
			{
				CreateIssue(repoURL, repoName, token)
			}
		case 4:
			{
				UpdateIssue(repoURL, repoName, token)
			}
		case 5:
			{
				CloseIssue(repoURL, repoName, token)
			}
		case 6:
			{
				sign = false
			}
		}
	}
	fmt.Println("process completed")
	os.Exit(0)
}
