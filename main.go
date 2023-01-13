package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Worker struct {
	duration int
	numtasks int
}

// Start the worker
func (w *Worker) Start() {
	log.Println("Background Worker Started Successfully - Waiting for Tasks")
	for t := range time.Tick(time.Duration(w.duration) * time.Second) {
		err := processTasks(t, w.numtasks)
		if err != nil {
			log.Fatal(err)
			continue
		}
	}
}

// process tasks every t seconds for all tasks in numtasks
func processTasks(t time.Time, numtasks int) error {
	var tasks []string
	for i := 0; i < numtasks; i++ {
		tasks = append(tasks, fmt.Sprintf("Subscription %d", i))
	}

	for _, task := range tasks {
		go sendEmail(task)
	}
	return nil
}

// send a single email
func sendEmail(t string) {
	time.Sleep(5 * time.Second)
	log.Println("Sending Notification Email for: ", t)
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
