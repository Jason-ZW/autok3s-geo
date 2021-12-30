package cmd

import (
	"fmt"
	"os"

	"github.com/Jason-ZW/autok3s-geo/pkg/common"

	"github.com/morikuni/aec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const ascIIStr = `
               ,        , 
  ,------------|'------'|             _        _    _____ 
 / .           '-'    |-             | |      | |  |____ | 
 \\/|             |    |   __ _ _   _| |_ ___ | | __   / / ___
   |   .________.'----'   / _  | | | | __/ _ \| |/ /   \ \/ __|
   |   |        |   |    | (_| | |_| | || (_) |   <.___/ /\__ \
   \\___/        \\___/   \__,_|\__,_|\__\___/|_|\_\____/ |___/

`

var (
	cmd = &cobra.Command{
		Use:              "autok3s-geo",
		Short:            "autok3s-geo is used to collects metrics about locates remote IP-address and exposes metrics to InfluxDB.",
		Long:             `autok3s-geo is used to collects metrics about locates remote IP-address and exposes metrics to InfluxDB.`,
		TraverseChildren: true,
	}
)

func init() {
	cmd.PersistentFlags().BoolVarP(&common.Debug, "debug", "d", common.Debug, "Enable log debug level")
}

// Command root command.
func Command() *cobra.Command {
	cmd.Run = func(cmd *cobra.Command, args []string) {
		printASCII()

		if err := cmd.Help(); err != nil {
			logrus.Errorln(err)
			os.Exit(1)
		}
	}
	return cmd
}

func printASCII() {
	fmt.Print(aec.Apply(ascIIStr))
}
