package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
	"time"
)

func parseGitHubCommitsResponse(body string) []GithubCommit {
	type ghCommit struct {
		Commit struct {
			Message   string `json:"message"`
			Committer struct {
				Date time.Time `json:"date"`
			} `json:"committer"`
		} `json:"commit"`
	}

	var commits []GithubCommit

	var raw []ghCommit
	if err := json.Unmarshal([]byte(body), &raw); err != nil {
		// return GotGitHubActivityFailed{err} // など、失敗の Msg にする
	}

	// pretty.Println(body)

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	for _, c := range raw {
		date := c.Commit.Committer.Date.In(jst)
		commits = append(commits, GithubCommit{date: date, message: c.Commit.Message})
	}

	// pretty.Println(commits)

	return commits
}

type GitHubRepository struct {
	name string
}

func listRepositoryForAUser(r Http) []GitHubRepository {

	res := r.Get(
		"https://api.github.com/users/hiragi-dev/repos",
		url.Values{"sort": []string{"updated"}},
	)
	// pretty.Println(res)

	type ghRepository struct {
		Name string `json:"name"`
	}

	repos := make([]GitHubRepository, 0)
	var raw []ghRepository
	if err := json.Unmarshal([]byte(res), &raw); err != nil {
		// return GotGitHubActivityFailed{err} // など、失敗の Msg にする
	}

	for _, r := range raw {
		repos = append(repos, GitHubRepository{name: r.Name})
	}

	return repos
}

func getGitHubActivity(r Http, repoName string) []GithubCommit {
	commits := make([]GithubCommit, 0)

	url := fmt.Sprintf("https://api.github.com/repos/hiragi-dev/%s/commits", repoName)

	res := r.Get(url, nil)
	// pretty.Println(res)

	_commits := parseGitHubCommitsResponse(res)
	for i := range _commits {
		_commits[i].repoName = repoName
	}

	commits = append(commits, _commits...)

	return commits
}

func GetGitHubActivity[E Http](r E) []GithubCommit {
	// Cmd[E, UserMsg] {
	// return AsyncDo(func(r E) UserMsg {
	repos := listRepositoryForAUser(r)

	var commits []GithubCommit

	for i, repo := range repos {
		if i < 7 {
			commits = append(commits, getGitHubActivity(r, repo.name)...)
			// pretty.Println(i, commits)
		}
	}

	slices.SortFunc(commits, func(a GithubCommit, b GithubCommit) int {
		if a.date.After(b.date) {
			return -1
		} else if b.date.After(a.date) {
			return 1
		} else {
			return 0
		}
	})

	// pretty.Println(commits)

	return commits
	// })
}
