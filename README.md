### Bitbucket utils
#### List Pull Requests Activity

Install:

``` shell
$ get -u github.com/auvn/go-atlassian/bitbucketutil/cmd/lspractivity
```

Create configuration file:

``` shell
$ cat <<EOT>> ~/.work/.lspractivity.config
authtoken: <your generated bitbucket token>
url: https://<your bitbucket host>/rest
EOT
```

Define an alias to simplify usage of the command:

``` shell
$ alias lspractivity="lspractivity -config ~/.work/.lspractivity.config"
```

Use to see comments activity in your open pull requests:

``` shell
$ lspractivity -age 24h | less # show comments for the last 24hours
```
