package main

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/roylee0704/gron"
)

func init() {
}

type gitFetcher struct {
	wg           sync.WaitGroup
	branch       string
	dir          string
	currentSHA   string
	pollInterval time.Duration
	c            *gron.Cron
}

func NewGitFetcher(dir, branch string, pollInterval time.Duration) *gitFetcher {
	return &gitFetcher{
		dir:          dir,
		branch:       branch,
		pollInterval: pollInterval,
		c:            gron.New(),
	}
}

func (g *gitFetcher) fetch() error {
	fetchCmd := exec.Command("git", "fetch")
	fetchCmd.Dir = g.dir
	if err := fetchCmd.Run(); err != nil {
		return err
	}
	mergeCmd := exec.Command("git", "merge", "FETCH_HEAD")
	mergeCmd.Dir = g.dir
	if err := mergeCmd.Run(); err != nil {
		return err
	}
	return nil
}

func (g *gitFetcher) revParse() (string, error) {
	cmd := exec.Command("git", "rev-parse", g.branch)
	cmd.Dir = g.dir
	out, err := cmd.Output()
	return string(out), err
}

func (g *gitFetcher) OnNewSHA(fn func(sha string)) {
	g.c.AddFunc(gron.Every(g.pollInterval), func() {
		if err := g.fetch(); err != nil {
			fmt.Printf("Error Fetching Git Repo: %v\n", err)
			return
		}

		sha, err := g.revParse()
		if err != nil {
			fmt.Printf("Error with rev-parse: %v\n", err)
			return
		}

		if sha != g.currentSHA {
			fmt.Println(sha)
			fn(sha)
			g.currentSHA = sha
		}
	})
}

func (g *gitFetcher) Start() {
	g.wg.Add(1)
	g.c.Start()
}

func (g *gitFetcher) Stop() {
	g.wg.Done()
	g.c.Stop()
}

func (g *gitFetcher) Wait() {
	g.wg.Wait()
}

func main() {
	g := NewGitFetcher("/Users/michaeldesa/go/src/github.com/influxdata/benchy-mcbenchface", "master", 1*time.Second)
	g.OnNewSHA(func(sha string) {
		fmt.Println(sha)
	})
	g.Start()

	g.Wait()
}
