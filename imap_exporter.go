package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goccy/go-yaml"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace = "imap"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Whether scraping IMAP metrics was successful.",
		nil,
		nil,
	)

	fll = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "first_loading_latency_seconds"),
		"Latency first loading of MUA.",
		[]string{"user", "cmd"},
		nil,
	)
)

type User struct {
	Username string
	Password string
}

type IMAPCollector struct {
	address string
	users   []*User
	logger  log.Logger
	failure bool
	sync.Mutex
}

func (i *IMAPCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- fll
}

func (i *IMAPCollector) Collect(ch chan<- prometheus.Metric) {
	i.Lock()
	defer i.Unlock()

	for _, u := range i.users {
		t := StartTimer()
		cl, err := NewClient(i.address)
		if err != nil {
			level.Error(i.logger).Log("msg", "Error connecting for IMAP", "err", err)
			i.failure = true
			return
		}
		ch <- prometheus.MustNewConstMetric(fll, prometheus.GaugeValue, t.Lap(), u.Username, "conn")
		defer cl.Close()

		if err := cl.Send("LOGIN", u.Username, u.Password); err != nil {
			level.Error(i.logger).Log("msg", "Error LOGIN command for IMAP",
				"err", err, "sent", cl.LastSent, "received", cl.LastReceived)
			i.failure = true
			return
		}
		ch <- prometheus.MustNewConstMetric(fll, prometheus.GaugeValue, t.Lap(), u.Username, "login")

		if err = cl.Send("LIST", "\"\"", "*"); err != nil {
			level.Error(i.logger).Log("msg", "Error LIST command for IMAP",
				"err", err, "sent", cl.LastSent, "received", cl.LastReceived)
			i.failure = true
			return
		}
		ch <- prometheus.MustNewConstMetric(fll, prometheus.GaugeValue, t.Lap(), u.Username, "list")

		if err = cl.Send("SELECT", "\"INBOX\""); err != nil {
			level.Error(i.logger).Log("msg", "Error SELECT command for IMAP",
				"err", err, "sent", cl.LastSent, "received", cl.LastReceived)
			i.failure = true
			return
		}
		ch <- prometheus.MustNewConstMetric(fll, prometheus.GaugeValue, t.Lap(), u.Username, "select")

		if err = cl.Send("FETCH", "1:20", "RFC822.HEADER"); err != nil {
			level.Error(i.logger).Log("msg", "Error FETCH command for IMAP",
				"err", err, "sent", cl.LastSent, "received", cl.LastReceived)
			i.failure = true
			return
		}
		ch <- prometheus.MustNewConstMetric(fll, prometheus.GaugeValue, t.Lap(), u.Username, "fetch")

		if err = cl.Send("LOGOUT"); err != nil {
			level.Error(i.logger).Log("msg", "Error LOGOUT command for IMAP",
				"err", err, "sent", cl.LastSent, "received", cl.LastReceived)
			i.failure = true
			return
		}
		ch <- prometheus.MustNewConstMetric(fll, prometheus.GaugeValue, t.Lap(), u.Username, "logout")
	}

	if i.failure {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0.0)
	} else {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1.0)
	}
}

func main() {
	prometheus.MustRegister(version.NewCollector("imap_exporter"))

	var (
		configFile    = kingpin.Flag("config", "Path to config yaml file for IMAP login.").Required().String()
		imapAddress   = kingpin.Flag("imap-address", "Address to connect on for IMAP.").Default("localhost:110").String()
		listenAddress = kingpin.Flag("listen-address", "Address to listen on for web interface and telemetry.").Default(":9993").String()
		metricsPath   = kingpin.Flag("telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version.Print("imap_exporter"))
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting imap_exporter", "version", version.Info())
	level.Info(logger).Log("build_context", version.BuildContext())

	html := `<html>
	<head>
		<title>IMAP4rev1 Exporter</title>
	</head>
	<body>
		<h1>IMAP4rev1 Exporter</h1>
		<p><a href='` + *metricsPath + `'>Metrics</a></p>
		<h2>Settings</h2>
		<ul>
			<li>Config File: ` + *configFile + `</li>
			<li>IMAP Address: ` + *imapAddress + `</li>
		</ul>
		<h2>Build</h2>
		<pre>` + version.Info() + ` ` + version.BuildContext() + `</pre>
		<hr />
		<small><a href="https://github.com/linyows/imap_exporter">Repository</a></small>
	</body>
</html>
`

	buf, err := ioutil.ReadFile(*configFile)
	if err != nil {
		level.Error(logger).Log("msg", "Error reading yaml file", "err", err)
		os.Exit(1)
	}

	var users []*User
	if err := yaml.Unmarshal(buf, &users); err != nil {
		level.Error(logger).Log("msg", "Error unmarshaling yaml file", "err", err)
		os.Exit(1)
	}

	c := &IMAPCollector{
		address: *imapAddress,
		users:   users,
		logger:  logger,
		failure: false,
	}

	prometheus.MustRegister(c)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	})

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
