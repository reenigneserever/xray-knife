package net

import (
	"github.com/reenigneserever/xray-knife/utils/customlog"
	"github.com/reenigneserever/xray-knife/xray"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// TcpCmd represents the tcp command
var TcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "Examine TCP Connection delay to config's host",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		parsed, err := xray.ParseXrayConfig(configLink)
		if err != nil {
			customlog.Printf(customlog.Failure, "Couldn't parse the config!\n")
			os.Exit(1)
		}
		generalDetails := parsed.ConvertToGeneralConfig()

		tcpAddr, err := net.ResolveTCPAddr("tcp", generalDetails.Address+":"+generalDetails.Port)
		if err != nil {
			customlog.Printf(customlog.Failure, "ResolveTCPAddr failed: %v\n", err)
			os.Exit(1)
		}
		start := time.Now()
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			customlog.Printf(customlog.Failure, "Couldn't establish tcp conn! : %v\n", err)
			os.Exit(1)
		}
		customlog.Printf(customlog.Success, "Established TCP connection in %dms\n", time.Since(start).Milliseconds())
		conn.Close()
	},
}

func init() {
	TcpCmd.Flags().StringVarP(&configLink, "config", "c", "", "The xray config link")
}
