package main

import (
	"fmt"
	"time"

	"github.com/influxdata/benchy-mcbenchface"
)

func main() {
	g := benchy.NewGitFetcher("/Users/michaeldesa/go/src/github.com/influxdata/benchy-mcbenchface", 1*time.Second)
	g.OnNewSHA("master", func(sha string) {
		fmt.Println(sha)
	})
	g.Start()

	g.Wait()
}
