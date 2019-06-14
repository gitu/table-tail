package commands

import (
	"database/sql"
	"fmt"
	"github.com/gitu/table-tail/pkg/tail"
	"github.com/gitu/table-tail/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"time"
)

var (
	connectionURI string
	driver        string
	cfgFile       string
	table         string
	id            string
	fields        string
	format        string
	interval      string
	rootCmd       = &cobra.Command{
		Use:   "table-tail",
		Short: "tails the specified table",
		Run: func(cmd *cobra.Command, args []string) {

			duration, err := time.ParseDuration(viper.GetString("interval"))
			checkErr(err)

			driver := viper.GetString("driver")
			db, err := sql.Open(driver, viper.GetString("connection-url"))
			checkErr(err)
			util, err := utils.Get(driver)
			checkErr(err)

			connectionInfo, err := util.ConnectionInfo(db)
			checkErr(err)

			fmt.Println(connectionInfo)

			target := make(chan string)
			go func() {
				for {
					msg := <-target
					fmt.Println(msg)
				}
			}()

			c, err := tail.Start(db, viper.GetString("table"), target,
				tail.Fields(viper.GetString("fields")),
				tail.ID(viper.GetString("id")),
				tail.Format(viper.GetString("format")),
				tail.Interval(duration),
				tail.Placeholder(util.PlaceHolderMarker()),
			)
			if err != nil {
				panic(err)
			}

			quit := make(chan os.Signal)
			signal.Notify(quit, os.Interrupt)
			<-quit

			c.Stop()
		},
	}
)

func checkErr(err error) {
	if err != nil {
		panic("err: " + err.Error())
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".tail-config")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
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
	rootCmd.PersistentFlags().StringVarP(&driver, "driver", "d", "", "driver to use")
	rootCmd.PersistentFlags().StringVarP(&table, "table", "t", "", "table to watch f.e. log_table")
	rootCmd.PersistentFlags().StringVarP(&fields, "fields", "s", "", "fields that are selected f.e. ID,MSG")
	rootCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "id for tailing f.e. ID")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "", "format - f.e. {{.msg}}")
	rootCmd.PersistentFlags().StringVarP(&interval, "interval", "b", "", "interval - if empty 500ms is used")

	checkErr(viper.BindPFlag("connection-url", rootCmd.PersistentFlags().Lookup("connection-url")))
	checkErr(viper.BindPFlag("driver", rootCmd.PersistentFlags().Lookup("driver")))
	checkErr(viper.BindPFlag("table", rootCmd.PersistentFlags().Lookup("table")))
	checkErr(viper.BindPFlag("fields", rootCmd.PersistentFlags().Lookup("fields")))
	checkErr(viper.BindPFlag("id", rootCmd.PersistentFlags().Lookup("id")))
	checkErr(viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format")))
	checkErr(viper.BindPFlag("interval", rootCmd.PersistentFlags().Lookup("interval")))
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
