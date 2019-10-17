package submissions

import (
	"encoding/csv"
	"encoding/json"
	"github.com/hhu-educode/cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	challenge string
	user string
	shouldPrintCsv bool
	shouldPrintJson bool
)

var submissionsListCommand = &cobra.Command{
	Use:   "list",
	Short: "Lists submissions",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new API client
		client, err := api.NewClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		submissions, err := client.GetSubmissions(user, challenge)
		if err != nil {
			log.Fatal(err)
		}

		if shouldPrintCsv {
			printCsv(submissions)
		} else if shouldPrintJson {
			printJson(submissions)
		} else {
			printTable(submissions)
		}
	},
}



func printCsv(submissions []api.SubmissionInfo)  {
	writer := csv.NewWriter(os.Stdout)
	writer.Write([]string{"Kennung", "Nachname", "Vorname", "Abgabe"})

	for _, submission := range submissions {
		writer.Write(submission.ToRow())
	}

	writer.Flush()
}

func printTable(submissions []api.SubmissionInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Challenge", "Student", "Punkte"})

	for _, submission := range submissions {
		table.Append(submission.ToRow())
	}

	table.Render()
}

func printJson(submissions []api.SubmissionInfo) {
	buffer, _ := json.MarshalIndent(submissions, "", "  ")
	print(string(buffer))
}

func init() {
	submissionsListCommand.Flags().StringVar(&user, "user", "", "Searches for a specific user's submissions")
	submissionsListCommand.Flags().StringVar(&challenge, "challenge", "", "Searches for a specific challenge")
	submissionsListCommand.Flags().BoolVar(&shouldPrintCsv, "csv", false, "Prints the submissions in csv format")
	submissionsListCommand.Flags().BoolVar(&shouldPrintJson, "json", false, "Prints the submissions in json format")
}
