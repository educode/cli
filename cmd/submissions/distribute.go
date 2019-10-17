package submissions

import (
	"encoding/json"
	"fmt"
	"github.com/hhu-educode/cli/api"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Tutor struct {
	Id string `json:"id"`
	Hours int `json:"hours"`
	Count int
}

const (
	optionYes = "Yes, do it!"
	optionNo = "No, stop!"
)

var (
	seed int64
	dryRun bool
	yes bool
)

var submissionsDistributeCommand = &cobra.Command{
	Use:   "distribute [challenge] [tutors] [deadline]",
	Short: "Distributes submissions on a set of tutors",
	SilenceErrors: true,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		distribute(args[0], args[1], args[2])
	},
}

func distribute(challenge string, tutorsFile string, deadline string) {
	tutors, err := readTutors(tutorsFile)
	if err != nil {
		log.Fatal(err)
	}

	client, err := api.NewClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	members, err := client.GetAllMembers()
	if err != nil {
		log.Fatal(err)
	}

	submissions, err := client.GetSubmissions("", challenge)
	if err != nil {
		log.Fatal(err)
	}

	names := make(map[string]string)
	for _, member := range members {
		names[member.Id] = fmt.Sprintf("%s %s", member.Firstname, member.Lastname)
	}

	checkTutors(tutors, names)

	assignments := assignSubmissions(tutors, filterStudents(submissions, members), names)

	if _, err := time.Parse(time.RFC3339, deadline); err != nil {
		log.Fatal(err)
	}

	var reviews []api.Review
	for tutor, students := range assignments {
		for _, student := range students {
			reviews = append(reviews, api.Review{
				Challenge:     challenge,
				Reviewer:      api.User{Id:tutor},
				Student:       api.User{Id:student},
				Deadline:      deadline,
				Content:       "",
				PointsRevoked: false,
			})
		}
	}

	if dryRun {
		return
	}

	if !yes {
		prompt := promptui.Select{
			Label: fmt.Sprintf("Distribute %d submissions from challenge %s? ", len(reviews), challenge),
			Items: []string{optionNo, optionYes},
			HideHelp: true,
			HideSelected: true,
		}

		_, result, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		if result == optionNo {
			fmt.Println("Distribution cancelled")
			return
		}
	}

	if err := client.CreateReviews(reviews); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sucessfully assigned %d submission\n", len(reviews))
}

func checkTutors(tutors []Tutor, names map[string]string) {
	for _, tutor := range tutors {
		if _, exists := names[tutor.Id]; !exists {
			log.Fatalf("Unrecognized tutor id %s", tutor.Id)
		}
	}
}

func assignSubmissions(tutors []Tutor, submissions []api.SubmissionInfo, names map[string]string) map[string][]string {
	rand.Seed(seed)
	submissionCount := len(submissions)
	totalHours := sumHours(tutors)

	var totalCount = 0
	for i, _ := range tutors {
		tutor := &tutors[i]
		tutor.Count = int(math.Round(float64(tutor.Hours) / float64(totalHours) * float64(submissionCount)))
		totalCount += tutor.Count
	}

	for totalCount != submissionCount {
		if totalCount < submissionCount {
			tutors[rand.Intn(len(tutors))].Count += 1
			totalCount += 1
		} else {
			tutors[rand.Intn(len(tutors))].Count -= 1
			totalCount -= 1
		}
	}

	printTutorTable(tutors, names)

	assignments := make(map[string][]string)
	for _, tutor := range tutors {
		for i := 0; i < tutor.Count; i++ {
			submission := submissions[0]
			submissions = submissions[1:]
			assignments[tutor.Id] = append(assignments[tutor.Id], submission.User)
		}
	}

	return assignments
}

func sumHours(tutors []Tutor) int {
	sum := 0
	for _, tutor := range tutors {
		sum += tutor.Hours
	}
	return sum
}

func readTutors(path string) ([]Tutor, error) {
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tutors []Tutor
	if err := json.Unmarshal(buffer, &tutors); err != nil {
		return nil, err
	}

	return tutors, nil
}

func filterStudents(submissions []api.SubmissionInfo, members []api.MemberInfo) []api.SubmissionInfo {
	isStudent := make(map[string]bool)
	for _, member := range members {
		if member.Role == "student" {
			isStudent[member.Id] = true
		}
	}

	var students []api.SubmissionInfo
	for _, submission := range submissions {
		if isStudent[submission.User] {
			students = append(students, submission)
		}
	}

	return students
}

func printTutorTable(tutors []Tutor, names map[string]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tutor", "Abgaben"})

	for _, tutor := range tutors {
		table.Append([]string{names[tutor.Id], strconv.Itoa(tutor.Count)})
	}

	table.Render()
}

func init() {
	submissionsDistributeCommand.Flags().BoolVar(&yes, "yes", false, "Skip confirmation prompt")
	submissionsDistributeCommand.Flags().BoolVar(&dryRun, "dry-run", false, "Perform a dry run")
	submissionsDistributeCommand.Flags().Int64Var(&seed, "seed", 0, "RNG seed")
}
