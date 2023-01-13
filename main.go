package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

type Worker struct {
	duration int
	numtasks int
}

// send a single email
func sendEmail(ctx context.Context, t string) {
	log.Println("Sending Email for: ", t)
	time.Sleep(10 * time.Second)
	log.Printf("%s Email Sent!", t)
	<-ctx.Done()
}

// process tasks every t seconds for all tasks in numtasks
func processTasks(t time.Time, numtasks int) error {
	var tasks []string
	for i := 0; i < numtasks; i++ {
		tasks = append(tasks, fmt.Sprintf("Subscription %d", i))
	}

	// create a new context with timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, task := range tasks {
		go sendEmail(ctx, task)
	}
	return nil
}

// Start the worker
func (w *Worker) Start() {
	log.Println("Background Worker Started Successfully - Waiting for Tasks")
	for t := range time.Tick(time.Duration(w.duration) * time.Second) {
		log.Printf("Polling database... found %d emails to send", w.numtasks)
		err := processTasks(t, w.numtasks)
		if err != nil {
			log.Fatal(err)
			continue
		}
	}
}

func main() {
	// Initialize a worker with poll every second and execute 10 tasks
	w := Worker{duration: 1, numtasks: 10}

	// Start worker in the background
	go w.Start()

	// Stop the worker by pressing any key
	fmt.Println("Press any key to stop the worker")
	var input string
	fmt.Scanln(&input)
	fmt.Println("Stopping the worker")
	os.Exit(0)
}
