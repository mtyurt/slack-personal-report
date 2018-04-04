slack-personal-report
==========

A small program that prepares a report of [Slack](https://slack.com) messages you have sent.

Installation
------------

#### Binary installation

[Download](https://github.com/mtyurt/slack-personal-report/releases) a
compatible binary for your system. For convenience, place `slack-personal-report` in a
directory where you can access it from the command line. Usually this is
`/usr/local/bin`.

```bash
$ mv slack-personal-report /usr/local/bin
```

#### Via Go

If you want you can also get `slack-personal-report` via Go:

```bash
$ go get -u github.com/mtyurt/slack-personal-report
```

Setup
-----

1. Get a slack token, click [here](https://api.slack.com/docs/oauth-test-tokens) 
2. Install `slack-personal-project` as mentioned above.

Usage
-----
This program can be run with simply:

```bash
$ slack-personal-report your-auth-token
```

