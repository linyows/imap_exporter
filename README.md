IMAP4rev1 Exporter
==

<img alt="GitHub Workflow" src="https://img.shields.io/github/workflow/status/linyows/imap_exporter/Go?style=for-the-badge">

Export IMAP4rev1 server health to Prometheus.

Usage
--

To run it:

```sh
$ make
$ ./imap_exporter --config=CONFIG [<flags>]
```

Exported Metrics
--

Metric | Meaning | Labels
--- | --- | ---
imap_up | Whether scraping IMAP metrics was successful |
first_loadin_latency_seconds | Latency first loading of MUA | user, cmd

First Loading Latency is:

- TCP connection time seconds
- Login command time seconds
- List command time seconds
- Select command time seconds
- Fetch command time seconds
- Logout command time seconds

Contribution
--

1. Fork (https://github.com/linyows/imap_exporter/fork)
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the go test ./... command and confirm that it passes
1. Run gofmt -s
1. Create a new Pull Request

Author
--

[linyows](https://github.com/linyows)

