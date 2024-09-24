/*
Copyright (c) 2013-present Snowplow Analytics Ltd.
All rights reserved.
This software is made available by Snowplow Analytics, Ltd.,
under the terms of the Snowplow Limited Use License Agreement, Version 1.0
located at https://docs.snowplow.io/limited-use-license-1.0
BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY PORTION
OF THE SOFTWARE, YOU AGREE TO THE TERMS OF SUCH LICENSE AGREEMENT.
*/

package ds

import (
	"log/slog"
	"os"

	"github.com/snowplow-product/snowplow-cli/internal/config"
	. "github.com/snowplow-product/snowplow-cli/internal/logging"
	"github.com/spf13/cobra"
)

var DataStructuresCmd = &cobra.Command{
	Use:     "data-structures",
	Aliases: []string{"ds"},
	Short:   "Work with Snowplow data structures",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := InitLogging(cmd); err != nil {
			return err
		}

		if err := config.InitConsoleConfig(cmd); err != nil {
			slog.Error("config failure", "error", err)
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	config.InitConsoleFlags(DataStructuresCmd)
}