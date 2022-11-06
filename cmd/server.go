//go:build server
// +build server

package cmd

import (
	"fmt"
	"net/http"
	"time"

	C "github.com/dmzlingyin/clipshare/pkg/constant"
	"github.com/dmzlingyin/clipshare/pkg/log"
	"github.com/dmzlingyin/clipshare/routers"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server command uses in server that transmit received message to other devices",
	Run: func(cmd *cobra.Command, args []string) {
		server()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func server() {
	router := routers.InitRouter()
	s := &http.Server{
		Addr:           ":" + fmt.Sprintf("%d", C.ServerConf.Port),
		Handler:        router,
		ReadTimeout:    time.Minute,
		WriteTimeout:   time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	log.InfoLogger.Println("server start")
	s.ListenAndServe()
}
