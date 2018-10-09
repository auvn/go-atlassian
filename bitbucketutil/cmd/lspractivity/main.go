package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/auvn/go-atlassian/atlassian"
	"github.com/auvn/go-atlassian/bitbucketutil/practivity"
	"gopkg.in/yaml.v2"
)

var (
	options = struct {
		ConfigFile    string
		MaxCommentAge time.Duration
		All           bool
	}{}
)

func fatal(err error) {
	log.Fatal(err)
}

func init() {
	flag.StringVar(&options.ConfigFile, "config", ".config", "configuration file")
	flag.DurationVar(&options.MaxCommentAge, "age", 0, "max comment age")
	flag.BoolVar(&options.All, "all", false, "show activity for all pull requests")
}

type Config struct {
	AuthToken string
	URL       string
}

func config() (cfg Config) {
	bb, err := ioutil.ReadFile(options.ConfigFile)
	if err != nil {
		fatal(err)
	}

	if err := yaml.Unmarshal(bb, &cfg); err != nil {
		fatal(err)
	}

	return cfg
}

func main() {
	flag.Parse()

	cfg := config()

	client := &atlassian.DefaultClient{
		Client: atlassian.Client{
			Auth:    "Bearer " + cfg.AuthToken,
			BaseURL: cfg.URL,
		},
	}
	activity, err := practivity.List(client, practivity.ListParams{
		IsAuthor: !options.All,
		MaxAge:   options.MaxCommentAge,
	})
	if err != nil {
		fatal(err)
	}

	activity.Fprintf(os.Stdout)
}
