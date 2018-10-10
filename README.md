[![Build Status](https://travis-ci.com/auvn/go-atlassian.svg)](https://travis-ci.com/auvn/go-atlassian)

### Bitbucket utils
#### List Pull Requests Activity

Install:

``` shell
$ go get -u github.com/auvn/go-atlassian/bitbucketutil/cmd/lspr
```

Create configuration file:

``` shell
$ cat <<EOT>> ~/.work/.lspr.config
authtoken: <your generated bitbucket token>
url: https://<your bitbucket host>/rest
EOT
```

Define an alias to simplify usage of the command:

``` shell
$ alias lspr="lspr -config ~/.work/.lspr.config"
```

Use to see comments activity in your open pull requests:

``` shell
$ lspr -age 24h | less # show comments for the last 24hours
```
