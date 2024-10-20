# webshare-cli

Simple CLI utility to manage files on Webshare.cz

```
go build .
```


## Before you start using it

- Copy `wscli.sample` to `~/.wscli` and fill `username` and `password` fields
- Copy `wscli` binary to any path stored in `$PATH`
- Ensure `wscli status` returns correct values for your Webshare.cz account


## Usage

```
$ wscli -h
$ wscli [command] -h
```