package whitelist

import (
	"encoding/json"
	"github.com/hhu-educode/cli/api"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var whitelistSyncCommand = &cobra.Command{
	Use:   "sync [members]",
	Short: "Synchronizes the whitelist using the specified members",
	SilenceErrors: true,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		syncWhiteList(args[0])
	},
}

func syncWhiteList(membersPath string) {
	buffer, err := ioutil.ReadFile(membersPath)
	if err != nil {
		log.Fatal(err)
	}

	var members api.Members
	if err := json.Unmarshal(buffer, &members); err != nil {
		log.Fatal(err)
	}

	client, err := api.NewClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.SyncWhitelist(members); err != nil {
		log.Fatal(err)
	}

	log.Printf("Synced %d members", len(members.Usernames))
}

func init() {

}
