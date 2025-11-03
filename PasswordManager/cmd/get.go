package cmd

import (
	"passmana/config"
	database "passmana/dBaseActions"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Add a new vault entry",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Print("Enter master password: ")
		// mpassword, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		// mpassword = strings.TrimSpace(mpassword)
		database.Fetch(username, platform, config.MasterPassword)

	},
}

func init() {

	getCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	getCmd.Flags().StringVarP(&platform, "platform", "p", "", "Platform name")
	getCmd.MarkFlagRequired("platform")
	getCmd.MarkFlagRequired("username")

	rootCmd.AddCommand(getCmd)
}
