package proxy

import (
	"bufio"
	"fmt"
	"os/exec"
)

var (
	stdOutChan = make(chan logLine)
)

type proxyService struct {
	Command []string `json:"command"`
	WorkDir string   `json:"workDir"`
}

func serviceWorker(service proxyService, logger Logger) error {
	// Repeatedly start the service when it ends (babysitting)
	for {
		// Start the service
		cmd := exec.Command(service.Command[0], service.Command[1:]...)
		cmd.Dir = service.WorkDir

		logger.debug("getting stdout pipe")
		stdOut, err := cmd.StdoutPipe()
		if err != nil {
			logger.error(fmt.Sprintf("error getting stdout pipe: %s", err))
			return err
		}
		logger.debug("got stdout pipe")

		logger.debug("starting service")
		err = cmd.Start()
		if err != nil {
			logger.error(fmt.Sprintf("error starting service: %s", err))
			return err
		}

		// Read stdout line by line
		scanner := bufio.NewScanner(stdOut)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			logger.info(scanner.Text())
		}

		err = cmd.Wait()
		if err != nil {
			logger.error(fmt.Sprintf("service exited with error: %s", err))
		}
	}
}

func startServiceWithChannelLogger(service proxyService, serviceName string) {
	logger := NewChannelLogger(serviceName)
	go serviceWorker(service, logger)
}
