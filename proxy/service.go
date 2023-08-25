package proxy

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

var (
	stdOutChan = make(chan logLine)
)

type proxyService struct {
	Command []string `json:"command"`
	WorkDir string   `json:"workDir"`
}

type logLine struct {
	line        string
	serviceName string
}

func (l logLine) String() string {
	return fmt.Sprintf("[%s] %s", l.serviceName, l.line)
}

func stdOutReceiver() {
	for {
		select {
		case line := <-stdOutChan:
			log.Println(line)
		}
	}
}

func serviceWorker(service proxyService, serviceName string) error {
	// Start the service
	cmd := exec.Command(service.Command[0], service.Command[1:]...)
	cmd.Dir = service.WorkDir

	log.Println("getting stdout pipe for", serviceName)
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("error getting stdout pipe for service %s: %s", serviceName, err)
		return err
	}
	log.Printf("got stdout pipe for %s", serviceName)
	stdOutChan <- logLine{fmt.Sprintf("starting service %s", serviceName), serviceName}
	log.Printf("starting service %s", serviceName)
	err = cmd.Start()
	if err != nil {
		log.Printf("error starting service %s: %s", serviceName, err)
		return err
	}

	for {
		// Read stdout line by line
		scanner := bufio.NewScanner(stdOut)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			stdOutChan <- logLine{scanner.Text(), serviceName}
		}
	}

	log.Println("started service", serviceName)
	return cmd.Wait()
}

func startService(service proxyService, serviceName string) {
	log.Println("trying to start service", serviceName)
	go serviceWorker(service, serviceName)
}
