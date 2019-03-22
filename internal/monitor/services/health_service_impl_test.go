package services

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/alekns/tinyrstats/pkg/helpers/network"
	"github.com/go-kit/kit/log"
	. "github.com/smartystreets/goconvey/convey"
)

const testRequestURL = "http://127.0.0.1:12234/"

var server *http.Server

func bringUpServer(status, waitMs int) {
	var wait = make(chan struct{})
	var mux = http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(time.Duration(waitMs) * time.Millisecond)
		w.Header().Add("TestHeader", "TestValue")
		if req.Method == "POST" {
			body, _ := ioutil.ReadAll(req.Body)
			w.Header().Add("BodyValue", string(body))
		}
		w.WriteHeader(status)
		io.WriteString(w, "Hello, world!\n")
	})
	server = &http.Server{
		Addr:    "127.0.0.1:12234",
		Handler: mux,
	}
	go func() {
		go func() { close(wait) }()
		server.ListenAndServe()
	}()
	<-wait
}

func tearDownServer() {
	server.Close()
}

const httpSuccessStatus = 200
const httpSuccessError = 500

func TestHTTPResourceServiceSpec(t *testing.T) {
	Convey("Given created http health service", t, func(c C) {
		svc := NewHTTPHealthService(&config.TasksSettings{
			DefaultTimeout: 2000,
		}, log.NewNopLogger())

		c.Convey("When request is valid", func(c C) {
			bringUpServer(httpSuccessStatus, 10)
			c.Reset(func() {
				tearDownServer()
			})

			c.Convey("The response should has success values", func(c C) {

				status, err := svc.CheckStatus(context.Background(), &monitor.HealthTask{
					URL:     testRequestURL,
					Timeout: 1000,
					Method:  "GET",
				})

				c.So(err, ShouldBeNil)
				c.So(status, ShouldNotBeNil)
				c.So(status.StatusCode, ShouldEqual, httpSuccessStatus)
			})
		})

		c.Convey("When server return error status code", func(c C) {
			bringUpServer(httpSuccessError, 10)
			c.Reset(func() {
				tearDownServer()
			})

			c.Convey("The response should has error values", func(c C) {

				status, err := svc.CheckStatus(context.Background(), &monitor.HealthTask{
					URL:     testRequestURL,
					Timeout: 1000,
					Method:  "GET",
				})

				c.So(err, ShouldBeNil)
				c.So(status, ShouldNotBeNil)
				c.So(status.Error, ShouldNotBeNil)
				c.So(status.Error.Text, ShouldContainSubstring, network.ErrHTTPServerError.Error())
				c.So(status.StatusCode, ShouldEqual, httpSuccessError)
			})
		})

		c.Convey("When request has too low timeout value", func(c C) {
			bringUpServer(httpSuccessStatus, 100)
			c.Reset(func() {
				tearDownServer()
			})

			c.Convey("The response should has error values with true of isTimeout", func(c C) {
				status, err := svc.CheckStatus(context.Background(), &monitor.HealthTask{
					URL:     testRequestURL,
					Timeout: 15,
					Method:  "GET",
				})

				c.So(err, ShouldBeNil)
				c.So(status, ShouldNotBeNil)
				c.So(status.Error, ShouldNotBeNil)
				c.So(status.Error.IsTimeout, ShouldBeTrue)
				c.So(status.Error.IsDNSError, ShouldBeFalse)
				c.So(status.Error.Text, ShouldContainSubstring, "timeout")
			})
		})

	})
}
