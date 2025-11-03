/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"passmana/database"

	"github.com/spf13/cobra" // ← 1 import
	_ "modernc.org/sqlite"   // driver
)

var dbConn *sql.DB

func DB() *sql.DB { return dbConn } // ← ADD THIS LINE

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "passmana",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		//Ask master passwrod once per run
		dbConn, err = sql.Open("sqlite", "vault.db")
		if err != nil {
			return fmt.Errorf("failed to connect to DB: %w", err)
		}
		database.Set(dbConn) //set global db
		// Auto-create table
		_, err = dbConn.Exec(`
			CREATE TABLE IF NOT EXISTS creds (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT NOT NULL,
				platform TEXT NOT NULL UNIQUE,
				salt BLOB NOT NULL,
				nonce BLOB NOT NULL,
				
				cipherpass BLOB NOT NULL 
				
			);
			
			CREATE INDEX IF NOT EXISTS idx_platform_username ON creds(platform, username);
		`)

		return err
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if dbConn != nil {
			dbConn.Close()
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.passmana.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
