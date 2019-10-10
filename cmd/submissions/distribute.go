package submissions

import (
	"encoding/json"
	"github.com/hhu-educode/cli/api"
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

var (
	seed int64
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

	submissions, err := client.GetSubmissions("", challenge)
	if err != nil {
		log.Fatal(err)
	}

	assignments := assignSubmissions(tutors, submissions)

	if _, err := time.Parse(time.RFC3339, deadline); err != nil {
		log.Fatal(err)
	}

	var reviews []api.Review
	for tutor, student := range assignments {
		reviews = append(reviews, api.Review{
			Challenge:     challenge,
			Reviewer:      api.User{Id:tutor},
			Student:       api.User{Id:student},
			Deadline:      deadline,
			Content:       "",
			PointsRevoked: false,
		})
	}

	//if err := client.CreateReviews(reviews); err != nil {
	//	log.Fatal(err)
	//}
}

func assignSubmissions(tutors []Tutor, submissions []api.SubmissionInfo) map[string]string {
	rand.Seed(seed)
	submissionCount := len(submissions)
	totalHours := sumHours(tutors)

	println(strconv.Itoa(submissionCount))

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

	printTutorTable(tutors)

	assignments := make(map[string]string)
	for _, tutor := range tutors {
		for i := 0; i < tutor.Count; i++ {
			submission := submissions[0]
			submissions = submissions[1:]
			assignments[tutor.Id] = submission.User
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

func printTutorTable(tutors []Tutor) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tutor", "Abgaben"})

	for _, tutor := range tutors {
		table.Append([]string{tutor.Id, strconv.Itoa(tutor.Count)})
	}

	table.Render()
}

func init() {
	submissionsDistributeCommand.Flags().Int64Var(&seed, "seed", 0, "RNG seed")
}
