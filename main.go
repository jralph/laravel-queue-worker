package main

import (
	"errors"
	"fmt"
	"github.com/pborman/getopt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func main() {
	artisanPath := getopt.StringLong("artisan", 'a', "artisan", "The path to artisan executable.")
	numberOfProcesses := getopt.IntLong("processes", 'p', 5, "The number of works to run.")

	queue := getopt.StringLong("queue", 0, "default", "The queue to listen on.")
	delay := getopt.IntLong("delay", 0, 0, "Amount of time to delay failed jobs.")
	memory := getopt.IntLong("memory", 0, 128, "The memory limit in megabytes.")
	sleep := getopt.IntLong("sleep", 0, 3, "Number of seconds to sleep when no jobs are available.")
	timeout := getopt.IntLong("timeout", 0, 60, "The number of seconds a child process can run for.")
	tries := getopt.IntLong("tries", 0, 0, "The number of times to attempt a job.")

	staggered := getopt.BoolLong("staggered", 's', "", "Stagger the starting of processes.")

	getopt.Parse()

	runCommands(
		*artisanPath,
		*numberOfProcesses,
		*queue,
		*delay,
		*memory,
		*sleep,
		*timeout,
		*tries,
		*staggered,
	)
}

/**
 * Handle the asynchronous running of the worker commands.
 */
func runCommands(
	artisanPath string,
	numProcs int,
	queue string,
	delay int,
	memory int,
	sleep int,
	timeout int,
	tries int,
	staggered bool,
) {
	fmt.Println(" ==> Starting processes...")

	var waitGroup sync.WaitGroup

	waitGroup.Add(numProcs)

	for i := 0; i < numProcs; i++ {
		go runCommand(
			i,
			waitGroup,
			artisanPath,
			queue,
			delay,
			memory,
			sleep,
			timeout,
			tries,
			staggered,
		)
	}

	waitGroup.Wait()
}

/**
 * Handle the running of an artisan queue worker and
 * handle restarting and error logging.
 */
func runCommand(
	id int,
	waitGroup sync.WaitGroup,
	artisanPath string,
	queue string,
	delay int,
	memory int,
	sleep int,
	timeout int,
	tries int,
	staggered bool,
) {
	// If we have requested to stagger the process spawning
	// sleep for a random amount of time between 0 and 5 seconds.
	if staggered {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	}

	command := exec.Command(
		"php",
		artisanPath,
		"queue:work",
		fmt.Sprintf("--queue=%s", queue),
		fmt.Sprintf("--delay=%d", delay),
		fmt.Sprintf("--memory=%d", memory),
		fmt.Sprintf("--sleep=%d", sleep),
		fmt.Sprintf("--timeout=%d", timeout),
		fmt.Sprintf("--tries=%d", tries),
	)

	fmt.Println(
		fmt.Sprintf(
			"\033[36m ==>\033[37m Starting process\033[32m [%d]\033[33m [%s]\033[37m",
			id,
			strings.Join(command.Args, " "),
		),
	)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		// If the command has exited with an error
		// do not restart it automatically.
		defer waitGroup.Done()

		printError(errors.New(fmt.Sprintf("Process %d stopped due to an error.", id)))
	} else {
		// If the command has exited with success
		// restart the command.
		fmt.Printf("\033[31m ==> Process [%d] stopped, possibly due to queue restart signal. Restarting...\033[37m\n", id)

		runCommand(
			id,
			waitGroup,
			artisanPath,
			queue,
			delay,
			memory,
			sleep,
			timeout,
			tries,
			staggered,
		)
	}
}

/**
 * Print out a given error.
 */
func printError(err error) {
	fmt.Println(fmt.Sprintf("\033[31m ==> Error: %s \033[37m", err.Error()))
}
