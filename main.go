package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

const baseURL = "http"

func main() {
	baseURL := os.Getenv("JENKINS_BASE_URL")
	user := os.Getenv("JENKINS_USER")
	token := os.Getenv("JENKINS_TOKEN")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	if len(os.Args) < 2 || len(os.Args) > 2 {
		fmt.Println("invalid number of arguments")
		os.Exit(1)
	}

	jobName := os.Args[1]

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/job/%s/lastBuild/api/json", baseURL, jobName), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.SetBasicAuth(user, token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var payload status
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(payload)
}

type status struct {
	JobName   string      `json:"fullDisplayName"`
	Building  bool        `json:"building"`
	URL       string      `json:"url"`
	Result    *string     `json:"result"`
	Timestamp json.Number `json:"timestamp"`
}

func (s status) String() string {
	b := bytes.NewBuffer(nil)
	w := tabwriter.NewWriter(b, 8, 8, 1, '\t', 0)

	if s.Building {
		fmt.Fprintf(w, "Name\t%v\n", s.JobName)
		fmt.Fprintf(w, "URL\t%v\n", s.URL)
		fmt.Fprintf(w, "Building\t%v\n", s.Building)
	} else {
		fmt.Fprintf(w, "Name\t%v\n", s.JobName)
		fmt.Fprintf(w, "URL\t%v\n", s.URL)
		fmt.Fprintf(w, "Building\t%v\n", s.Building)
		fmt.Fprintf(w, "Result\t%v\n", *s.Result)

		// TODO: Jenkins returns a weird timestamp. Not quite the length of UnixNano, but not Unix either.
		v, _ := strconv.Atoi(s.Timestamp.String()[:10])
		t := time.Unix(int64(v), 0).UTC()
		fmt.Fprintf(w, "Time\t%v (%v)\n", t.Format("2006-01-02 3:04PM MST"), time.Since(t))
	}

	w.Flush()
	return b.String()
}
