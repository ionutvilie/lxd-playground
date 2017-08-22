package main
//
// simple worker pools example of concurrency in go
// the advantage of this version is that it limits the number of concurrent routines
//

import (
	"fmt"
	"net/http"
	"time"
	"flag"
	"strconv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	Requests int = 20
	//Jobs int = 10
)

var (
	addr = flag.String("listen-address", ":8081", "The address to listen on for HTTP requests.")

	httpResponsesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			//Namespace: "testvertical",
			//Subsystem: "http_server",
			Name:      "http_requests_total",
			Help:      "http_requests_total_help",
		},
		[]string{"code", "method", "url"},
	)

	httpResponseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:"http_request_duration_milliseconds",
			Help: "http_request_duration_milliseconds",
			Buckets: []float64{.005,.025,.05,.075, .1 ,.25, .5, .75, 1},
		},
		[]string{"code", "method", "url"},
	)


	urls = []string{
		"http://10.171.180.201/ServletApp/",
		"http://10.171.180.201/ServletApp/",
		"http://10.171.180.219:8080/ServletApp/",
		"http://10.171.180.219:8080/ServletApp/",
	}
)

func init() {
	flag.Parse()
	// ----  >>  change here  << ----
	prometheus.MustRegister(httpResponsesTotal,httpResponseDuration)
}

type WrapHTTPHandler struct {
	handler http.Handler
}

type LoggedResponse struct {
	http.ResponseWriter
	status int
}


func (wrappedHandler *WrapHTTPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	loggedWriter := &LoggedResponse{ResponseWriter: writer, status: 200}
	wrappedHandler.handler.ServeHTTP(loggedWriter, request)

	status := strconv.Itoa(loggedWriter.status)
	//log.SetPrefix("[Info]")
	//log.Printf("[RemoteAddr: %s] URL: %s, STATUS: %d, Time Elapsed: %dnanoseconds.\n",
	//	request.RemoteAddr, request.URL, loggedWriter.status, elapsed)

	httpResponsesTotal.WithLabelValues(status, request.Method, "/").Inc()
}




func urlProbeWorker(id int, url string, jobs <-chan int, results chan<- string) {
	for j := range jobs {
		start := time.Now()
		// define a timeout for the http client
		timeout := time.Duration(2 * time.Second)
		// define http status variable
		var probeHttpStatus int
		// initialize a client and set the timeout
		client := http.Client{
			Timeout: timeout,
		}

		//http_response, error := client.Head(url)
		http_response, error := client.Get(url)
		if error != nil {
			probeHttpStatus = 520 // status in case of timeout
		} else {
			probeHttpStatus = http_response.StatusCode

		}
		secs := time.Since(start).Seconds()
		httpResponsesTotal.WithLabelValues(fmt.Sprintf("%v",probeHttpStatus), "GET", url).Inc()
		httpResponseDuration.WithLabelValues(fmt.Sprintf("%v",probeHttpStatus), "GET", url).Observe(secs)

		//return data to channel
		results <- fmt.Sprintf("Worker: %2d Job: %3d %s %d %.6f", id, j, url, probeHttpStatus, secs)
	}
}






func mainHandler(writer http.ResponseWriter, request *http.Request) {
	// The "/" pattern matches everything, so we need to check that we're at the root here.
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}
	Jobs := len(urls) * 1  // assign n workers per url
	start := time.Now()

	// initialize a channel with max size the number of requests
	jobs := make(chan int, Requests)
	result := make(chan string, Requests)

	// start number Jobs
	for i := 1; i <= Jobs; i++ {
		// calculate modulo between i and url length and always get
		// a correct address inside the slice
		go urlProbeWorker(i, urls[i%len(urls)], jobs, result)
	}
	for j := 1; j <= Requests; j++ {
		jobs <- j
	}
	close(jobs)

	// print the response from the go routines
	for k := 1; k <= Requests; k++ {
		fmt.Fprintf(writer,"%3d - %s\n", k, <-result)
	}
	secs := time.Since(start).Seconds()

	fmt.Fprintf(writer,"Program executed in %.6f seconds\n", secs)
}



func main() {
	http.HandleFunc("/", mainHandler)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*addr, &WrapHTTPHandler{http.DefaultServeMux})

}
