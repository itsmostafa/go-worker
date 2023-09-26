package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/itsmostafa/go-worker/aws"
)

type Worker struct {
	duration int
	numtasks int
}

// send a single email
func sendEmail(ctx context.Context, t string, complete chan<- bool) {
	log.Println("Sending Email for: ", t)
	time.Sleep(5 * time.Second)
	log.Println("Email Sent Successfully")
	complete <- true

}

// process tasks every t seconds for all tasks in numtasks
func processTasks(t time.Time, numtasks int) {
	var tasks []string
	for i := 0; i < numtasks; i++ {
		tasks = append(tasks, fmt.Sprintf("Subscription %d", i))
	}

	// create a new context with timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	complete := make(chan bool, 2)

	for _, task := range tasks {
		go sendEmail(ctx, task, complete)
	}
}

// Start the worker
func (w *Worker) Start() {
	log.Println("Background Worker Started Successfully - Waiting for Tasks")
	for t := range time.Tick(time.Duration(w.duration) * time.Second) {
		log.Printf("Polling database... found %d emails to send", w.numtasks)
		processTasks(t, w.numtasks)
		fmt.Println("Number of Goroutines: ", runtime.NumGoroutine())
	}
}

func main() {
	// // Initialize a worker with poll every second and execute 10 tasks
	email := aws.SESEmail{
		From:        "test@testing.xyz",
		To:          "test@gmail.com",
		Subject:     "test",
		CC:          "test@inltesting.xyz",
		HtmlBody:    "<h1>test</h1><br><p>Hello this is a test from my golang app</p>",
		TextBody:    "text only test",
		PdfFileName: "test.pdf",
		PdfFile:     []byte("test"),
	}
	email.Send()
}
