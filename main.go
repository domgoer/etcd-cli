/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : main.go
#   Created       : 2019-01-29 11:53:54
#   Describe      :
#
# ====================================================*/
package main

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/peterh/liner"

	"github.com/domgoer/etcd-cli/cmd"

	cmd2 "github.com/domgoer/etcd-cli/pkg/cmd"
	"github.com/domgoer/etcd-cli/pkg/version"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          "etcd-cli",
		Long:         "",
		SilenceUsage: true,
		RunE:         rootE,
	}

	saveCmd = &cobra.Command{
		Use:          "upload",
		Short:        "Use to upload the configuration file in binary form to the legal path specified by etcd",
		RunE:         uploadE,
		SilenceUsage: true,
	}

	downloadCmd = &cobra.Command{
		Use:          "download",
		Short:        "Download the path-mapped file locally.",
		RunE:         downloadE,
		SilenceUsage: true,
	}

	cfg         cmd.Config
	historyPath = path.Join(os.Getenv("HOME"), ".etcdcli_history")
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfg.Host, "host", "s", "127.0.0.1", "Etcd connection host.")
	rootCmd.PersistentFlags().Int32VarP(&cfg.Port, "port", "p", 2379, "Etcd connection port.")
	rootCmd.PersistentFlags().StringVar(&cfg.Ca, "ca-path", "", "cafile path used to connect to an etcd.")
	rootCmd.PersistentFlags().StringVar(&cfg.Cert, "cert-path", "", "certfile path used to connect to an etcd.")
	rootCmd.PersistentFlags().StringVar(&cfg.Key, "key-path", "", "keyfile path used to connect to an etcd.")
	rootCmd.PersistentFlags().StringVar(&cfg.Username, "username", "", "username used to connect to an etcd.")
	rootCmd.PersistentFlags().StringVar(&cfg.Password, "password", "", "password used to connect to an etcd.")

	cmd2.AddFlags(rootCmd)
	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(version.Command())

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func rootE(_ *cobra.Command, _ []string) error {
	// set line
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	// load history from file
	loadHistory(line)
	defer saveHistory(line)

	r, err := cmd.NewRoot(cfg)
	if err != nil {
		return err
	}
	defer r.Close()

	reg, _ := regexp.Compile(`'.*?'|".*?"|\S+`)
	for {
		prompt := fmt.Sprintf("%s[%d]  %s>", cfg.Host, cfg.Port, cmd.PWD)

		c, err := line.Prompt(prompt)

		if err != nil {
			return err
		}
		line.AppendHistory(c)

		fields := reg.FindAllString(c, -1)

		err = r.DoScan(fields)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func uploadE(_ *cobra.Command, args []string) error {
	r, err := cmd.NewRoot(cfg)
	if err != nil {
		return err
	}
	defer r.Close()
	if len(args) != 2 {
		return cmd.ErrInvalidParamNum
	}
	return r.Upload(args[0], args[1])
}

func downloadE(_ *cobra.Command, args []string) error {
	r, err := cmd.NewRoot(cfg)
	if err != nil {
		return err
	}
	defer r.Close()
	if len(args) != 2 && len(args) != 1 {
		return cmd.ErrInvalidParamNum
	}
	var localP = "./"
	if len(args) == 2 {
		localP = args[1]
	}
	return r.Download(args[0], localP)
}

func saveHistory(line *liner.State) {
	if f, err := os.OpenFile(historyPath, os.O_CREATE|os.O_RDWR, 0666); err != nil {
		fmt.Errorf("Error writing history file: %s", err.Error())
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func loadHistory(line *liner.State) {
	if f, err := os.OpenFile(historyPath, os.O_RDONLY, 0444); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}
