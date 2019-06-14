package commands

import (
	"database/sql"
	"fmt"
	"github.com/gitu/table-tail/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	connectionURI string
	cfgFile       string
	rootCmd       = &cobra.Command{
		Use:   "install-avq-change",
		Short: "Installs avaloq change from exported zip file",
		Run: func(cmd *cobra.Command, args []string) {

			driver := viper.GetString("driver")
			db, err := sql.Open(driver, viper.GetString("connection-url"))
			checkErr(err)
			util, err := utils.Get(driver)
			checkErr(err)

			connectionInfo, err := util.ConnectionInfo(db)
			checkErr(err)

			fmt.Println(connectionInfo)

		},
	}
)

func checkErr(err error) {
	if err != nil {
		panic("err: " + err.Error())
	}
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName(".tail-config")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		fmt.Println("Will create configfile: .tail-config.yaml")
	}
	if err := viper.SafeWriteConfigAs(".tail-config.yaml"); err != nil {
		if os.IsNotExist(err) {
			err = viper.WriteConfigAs(".tail-config.yaml")
		}
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&connectionURI, "connection-url", "c", "", "url to connect to")
	rootCmd.PersistentFlags().StringVarP(&connectionURI, "driver", "d", "", "driver to use")
	checkErr(viper.BindPFlag("connection-url", rootCmd.PersistentFlags().Lookup("connection-url")))
	checkErr(viper.BindPFlag("driver", rootCmd.PersistentFlags().Lookup("driver")))
}

// SetDefaultDriver is used to set the default driver when calling from reduced util
func SetDefaultDriver(driver string) {
	viper.SetDefault("driver", driver)
}

// Execute calls the root program
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
