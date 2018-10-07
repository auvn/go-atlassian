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
	}{}
)

func fatal(err error) {
	log.Fatal(err)
}

func init() {
	flag.StringVar(&options.ConfigFile, "config", ".config", "configuration file")
	flag.DurationVar(&options.MaxCommentAge, "age", 0, "max comment age")
}

type Config struct {
	Auth string
	URL  string
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
			Auth:    cfg.Auth,
			BaseURL: cfg.URL,
		},
	}
	activity, err := practivity.List(client, options.MaxCommentAge)
	if err != nil {
		fatal(err)
	}

	activity.Fprintf(os.Stdout)
}
