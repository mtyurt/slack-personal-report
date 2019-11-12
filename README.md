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


```
$ ./slack-personal-report -help
Usage of slack-personal-report:
SLACK_TOKEN=your-token slack-personal-report [OPTIONAL ARGUMENTS]

  Optional Arguments:
  -daily
        To print only previous day's messages (default true)
  -days int
        Number of days to search for in daily mode. Because day search starts
		from midnight by Slack. (default 1)
  -extra-search string
        Default search mode is 'from:me', use this flag if you want extra
		conditions on top of it, e.g.: '-extra-search=in:#channel'; in the end
		the search filter will be: 'from:me in:#channel' (default " ")
  -short
        Print only short output (default false)
  -weekly
        To print only previous week's messages (default false)
```

Future Work
-----
- Replace mentions in text with actual usernames / channel names
- Colorful terminal output with optional flag
