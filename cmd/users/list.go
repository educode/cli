package users

import (
	"encoding/csv"
	"fmt"
	"github.com/hhu-educode/cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	page int
	pageSize int
	shouldPrintIdsOnly bool
	shouldPrintCsv bool
	includeEmpty bool
)

var usersListCommand = &cobra.Command{
	Use:   "list",
	Short: "Lists all users",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new API client
		client, err := api.NewClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		members, err := client.GetMembers(page, pageSize)
		if err != nil {
			log.Fatal(err)
		}

		if shouldPrintCsv {
			printCsv(members)
		} else if shouldPrintIdsOnly {
			printIds(members)
		} else {
			printTable(members)
		}
	},
}

func printIds(members []api.MemberInfo) {
	for _, member := range members {
		fmt.Println(member.Id)
	}
}

func printCsv(members []api.MemberInfo)  {
	writer := csv.NewWriter(os.Stdout)
	writer.Write([]string{"Kennung", "Anrede", "Vorname", "Nachname", "E-Mail", "Rolle"})

	for _, member := range members {
		writer.Write(member.ToRow())
	}

	writer.Flush()
}

func printTable(members []api.MemberInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Kennung", "Anrede", "Vorname", "Nachname", "E-Mail", "Rolle"})

	for _, member := range members {
		table.Append(member.ToRow())
	}

	table.Render()
}

func init() {
	usersListCommand.Flags().IntVar(&page, "page", 0, "The requested page")
	usersListCommand.Flags().IntVar(&pageSize, "page-size", 0, "Entries per page")

	usersListCommand.Flags().BoolVar(&shouldPrintIdsOnly, "ids", false, "Prints only user ids")
	usersListCommand.Flags().BoolVar(&shouldPrintCsv, "csv", false, "Prints the table in csv format")
	usersListCommand.Flags().BoolVar(&includeEmpty, "empty", false, "Includes empty submissions")
}
