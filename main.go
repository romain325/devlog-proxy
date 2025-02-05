package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
)

type Proxy struct{}

func copyHeader(dst http.Header, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func logIncoming(r *http.Request, logger CustomLogger) {
	logger("> In : ", "URL", r.URL, "Header", r.Header)
}

func logOutgoing(r *http.Response, body []byte, logger CustomLogger) {
	logger("< Out:", "Status", r.StatusCode, "Header", r.Header, "Body", string(body))
}

func handleFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readConf() ProxyConf {
	jsonFile, err := os.Open("conf.json")
	handleFatalError(err)

	defer jsonFile.Close()

	var proxyConf ProxyConf
	confBytes, err := io.ReadAll(jsonFile)
	handleFatalError(err)

	handleFatalError(json.Unmarshal(confBytes, &proxyConf))

	return proxyConf
}

func (s *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("log").(CustomLogger)

	logIncoming(r, logger)

	r.RequestURI = ""
	r.URL.Scheme = "http"
	r.URL.Host = "localhost:1080"

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Errorf("Error making http Request %s\n", err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger("Could not read response body", "error", err)
	}

	logOutgoing(res, resBody, logger)

	copyHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)
	w.Write(resBody)
}

func main() {
	conf := readConf()

	ctx, ctxCancel := context.WithCancel(context.Background())

	for _, serverConf := range conf {

		server := &http.Server{
			Addr:    ":" + strconv.Itoa(serverConf.ProxyPort),
			Handler: &Proxy{},
			BaseContext: func(listener net.Listener) context.Context {
				ctx = context.WithValue(ctx, "conf", serverConf)
				styles := log.DefaultStyles()
				styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
					SetString(serverConf.Name + "[:" + strconv.Itoa(serverConf.ProxyPort) + "]").
					Bold(true).
					Foreground(lipgloss.Color("#FAFAFA")).
					Background(lipgloss.Color(getRngHexColor())).
					Padding(1).
					Align(lipgloss.Center).
					Width(22)
				logger := log.New(os.Stdout)
				logger.SetStyles(styles)

				f, _ := os.OpenFile(serverConf.Name+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)

				fileLogger := log.New(f)
				fileLogger.SetFormatter(log.JSONFormatter)

				ctx = context.WithValue(ctx, "log", CustomLogger(func(mess string, val ...any) {
					logger.Info(mess, val...)
					fileLogger.Info(mess, val...)
				}))
				return ctx
			},
		}

		log.Infof("Server %s setted up", serverConf.Name)

		go func() {
			err := server.ListenAndServe()
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("Server " + serverConf.Name + " closed")
			} else {
				handleFatalError(err)
			}
			ctxCancel()
		}()
	}
	<-ctx.Done()

}
