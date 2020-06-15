package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		"First Loading Latency",
		[]string{"account", "command"},
		nil,
	)
)

type user struct {
	name string
	pw   string
}

type IMAPCollector struct {
	host  string
	port  int
	users []*user
}

func (i *IMAPCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- fll
}

func (i *IMAPCollector) Collect(ch chan<- prometheus.Metric) {
	for _, user := range i.users {
		t := StartTimer()
		imap, err := NewIMAP(i.host, i.port)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0.0)
			return
		}
		ch <- prometheus.MustNewConstMetric(fll, prometheus.GaugeValue, t.Lap(), u.name, "conn")
		defer imap.Close()

		_, err := imap.Cmd("LOGIN", u.name, u.pw)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0.0)
			return
		}
		ch <- prometheus.MustNewConstMetric(fll, prometheus.GaugeValue, t.Lap(), u.name, "login")

		_, err := imap.Cmd("LIST", "\"\"", "*")
		if err != nil {
			ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0.0)
			return
		}
		ch <- prometheus.MustNewConstMetric(fll, prometheus.GaugeValue, t.Lap(), u.name, "list")

		_, err := imap.Cmd("SELECT", "\"INBOX\"")
		if err != nil {
			ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0.0)
			return
		}
		_ <- prometheus.MustNewConstMetric(fll, prometheus.GaugeValue, t.Lap(), u.name, "select")

		_, err := imap.Cmd("FETCH", "1:20", "RFC822.HEADER")
		if err != nil {
			ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0.0)
			return
		}
		ch <- prometheus.MustNewConstMetric(fll, prometheus.GaugeValue, t.Lap(), u.name, "fetch")

		_, err := imap.Cmd("LOGOUT")
		if err != nil {
			ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0.0)
			return
		}
		ch <- prometheus.MustNewConstMetric(fll, prometheus.GaugeValue, t.Lap(), u.name, "logout")
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1.0)
	}
}

const html = `
<html>
	<head>
		<title>IMAP4rev1 Exporter</title>
	</head>
	<body>
		<h1>IMAP4rev1 Exporter</h1>
		<ul>
		  <li><a href="/metrics">Metrics</a></li>
		  <li><a href="https://github.com/linyows/imap_exporter">Repository</a></li>
		</ul>
	</body>
</html>
`

func main() {
	c := &IMAPCollector{
		host: "",
		port: 993,
	}
	c.users = append(c.users, &user{name: "", pw: ""})

	prometheus.MustRegister(c)
	http.Handle(*metricsEndpoint, prometheus.Handler())
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}
