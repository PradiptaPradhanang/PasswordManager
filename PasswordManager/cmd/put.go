package cmd

import (
	"passmana/config"
	database "passmana/dBaseActions"

	"github.com/spf13/cobra"
)

var platform, username, password string

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Add a new vault entry",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Print("Enter master password: ")
		// mpassword, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		// mpassword = strings.TrimSpace(mpassword)
		database.Insert(username, platform, password, config.MasterPassword)

	},
}

func init() {

	putCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	putCmd.Flags().StringVarP(&platform, "platform", "p", "", "Platform name")
	putCmd.Flags().StringVarP(&password, "password", "w", "", "Password")
	putCmd.MarkFlagRequired("platform")
	putCmd.MarkFlagRequired("username")
	putCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(putCmd)
}
