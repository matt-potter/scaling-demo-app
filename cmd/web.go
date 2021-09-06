package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/matt-potter/scaling-demo-app/logger"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Starts http listener",
	Run:   runWebHandler,
}

var (
	webDownstream string
	webPort       string
)

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringVarP(&webDownstream, "downstream", "d", "", "Downstream service host")
	webCmd.Flags().StringVarP(&webPort, "port", "p", "", "Listening port")
	sqsCmd.MarkFlagRequired("port")
}

func handler(w http.ResponseWriter, r *http.Request) {

	req, err := ioutil.ReadAll(r.Body)

	h := http.Client{
		Timeout: time.Second * 10,
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error: %s", err.Error())))
	}

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

	logger.Log.Infof("from upstream: msg: %s", req)

	if webDownstream == "" {

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("received %s", req)))

	} else {

		resp, err := h.Post(webDownstream, "text/plain", bytes.NewBuffer([]byte(req)))

		if err != nil {
			logger.Log.Fatal(err.Error())
		}

		rBody, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			logger.Log.Fatal(err.Error())
		}

		logger.Log.Infof("from downstream: status: %d, msg: %s", resp.StatusCode, rBody)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("received %s", req)))
	}

}

func runWebHandler(cmd *cobra.Command, args []string) {
	http.HandleFunc("/", handler)
	logger.Log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", webPort), nil))
}
