package main

import (
	"errors"
	"fmt"
	"github.com/pborman/getopt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	artisanPath := getopt.StringLong("artisan", 'a', "artisan", "The path to artisan executable.")
	numberOfProcesses := getopt.IntLong("processes", 'p', 5, "The number of works to run.")

	queue := getopt.StringLong("queue", 'q', "default", "The queue to listen on.")
	delay := getopt.IntLong("delay", 'd', 0, "Amount of time to delay failed jobs.")
	memory := getopt.IntLong("memory", 'm', 128, "The memory limit in megabytes.")
	sleep := getopt.IntLong("sleep", 's', 3, "Number of seconds to sleep when no jobs are available.")
	timeout := getopt.IntLong("timeout", 't', 60, "The number of seconds a child process can run for.")
	tries := getopt.IntLong("tries", 'r', 0, "The number of times to attempt a job.")

	getopt.Parse()

	runCommands(*artisanPath, *numberOfProcesses, *queue, *delay, *memory, *sleep, *timeout, *tries)
}

func runCommands(artisanPath string, numProcs int, queue string, delay int, memory int, sleep int, timeout int, tries int) {
	fmt.Println(" ==> Starting processes...")

	var wg sync.WaitGroup

	wg.Add(numProcs)

	for i := 0; i < numProcs; i++ {
		go runCommand(i, wg, artisanPath, queue, delay, memory, sleep, timeout, tries)
	}

	wg.Wait()
}

func runCommand(id int, wg sync.WaitGroup, artisanPath string, queue string, delay int, memory int, sleep int, timeout int, tries int) {
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

	fmt.Println(fmt.Sprintf("\033[36m ==>\033[37m Starting process\033[32m [%d]\033[33m [%s]\033[37m", id, strings.Join(command.Args, " ")))

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		defer wg.Done()

		printError(errors.New(fmt.Sprintf("Process %d stopped due to an error.", id)))
	} else {
		fmt.Printf("\033[31m ==> Process [%d] stopped, possibly due to queue restart signal. Restarting...\n\033[37m", id)
		runCommand(id, wg, artisanPath, queue, delay, memory, sleep, timeout, tries)
	}
}

func printOutput(output string) {
	fmt.Printf(" ==> Output: %s\n", output)
}

func printError(err error) {
	fmt.Println(fmt.Sprintf(" ==> Error: %s\n", err.Error()))
}
