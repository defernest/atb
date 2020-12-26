package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"log"
	"strings"
)

func GitlabClient(token string) *gitlab.Client {
	git, err := gitlab.NewClient(token, gitlab.WithBaseURL("http://192.168.20.51/api/v4"))
	if err != nil {
		log.Fatal(err)
	}
	return git
}

func GetListIssues(sprint string, token string) ([]*gitlab.Issue, error) {
	sprint = strings.TrimSpace(sprint)
	issueOpt := &gitlab.ListIssuesOptions{
		Milestone: gitlab.String(sprint),
		Scope:     gitlab.String("all"),
	}
	issues, resp, err := GitlabClient(token).Issues.ListIssues(issueOpt)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("ошибка при выполнении запроса к gitlab API: %s", err)
	}
	if resp.TotalItems < 1 {
		return nil, fmt.Errorf("не нашел такого спринта, либо в спринте '%s' нет задач", sprint)
	}
	return issues, nil
}