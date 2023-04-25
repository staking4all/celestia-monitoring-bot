package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/staking4all/celestia-monitoring-bot/services/models"
	"github.com/staking4all/celestia-monitoring-bot/services/monitor"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Daemon to monitor validators",
	Long:  "Monitors validators and pushes alerts to Telegram",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile, _ := cmd.Flags().GetString("file")
		dat, err := os.ReadFile(configFile)
		if err != nil {
			zap.L().Error("Error reading config.yaml", zap.Error(err))
			return err
		}

		config := models.Config{}
		err = yaml.Unmarshal(dat, &config)
		if err != nil {
			zap.L().Error("Error parsing config.yaml", zap.Error(err))
			return err
		}
		config.LoadDefault()

		// TODO: implement telegram register and alert service
		m, err := monitor.NewMonitorService(config)
		if err != nil {
			zap.L().Error("error stating monitor", zap.Error(err))
			return err
		}

		err = m.Run()
		if err != nil {
			zap.L().Error("monitor running", zap.Error(err))
			return err
		}

		return nil
	},
}

func init() {

	rootCmd.AddCommand(monitorCmd)
	monitorCmd.Flags().StringP("file", "f", "./config.yaml", "File path to config yaml")
}
