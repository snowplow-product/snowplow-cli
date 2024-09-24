/*
Copyright (c) 2013-present Snowplow Analytics Ltd.
All rights reserved.
This software is made available by Snowplow Analytics, Ltd.,
under the terms of the Snowplow Limited Use License Agreement, Version 1.0
located at https://docs.snowplow.io/limited-use-license-1.0
BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY PORTION
OF THE SOFTWARE, YOU AGREE TO THE TERMS OF SUCH LICENSE AGREEMENT.
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/snowplow-product/snowplow-cli/cmd/ds"
)

var RootCmd = &cobra.Command{
	Use:   "snowplow-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().String("config", "",
		`Config file. Defaults to $HOME/.config/snowplow/snowplow.yml
Then on:
  Unix $XDG_CONFIG_HOME/snowplow/snowplow.yml
  Darwin $HOME/Library/Application Support/snowplow/snowplow.yml
  Windows %AppData%\snowplow\snowplow.yml`,
	)
	RootCmd.PersistentFlags().Bool("debug", false, "Log output level to Debug")
	RootCmd.PersistentFlags().BoolP("quiet", "q", false, "Log output level to Warn")
	RootCmd.PersistentFlags().BoolP("silent", "s", false, "Disable output")
	RootCmd.PersistentFlags().Bool("json-output", false, "Log output as json")
	RootCmd.AddCommand(ds.DataStructuresCmd)
}
