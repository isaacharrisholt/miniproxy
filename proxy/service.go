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
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			logger.error(fmt.Sprintf("error getting stdout pipe: %s", err))
			return err
		}
		logger.debug("got stdout pipe")

		logger.debug("getting stderr pipe")
		stderr, err := cmd.StderrPipe()
		if err != nil {
			logger.error(fmt.Sprintf("error getting stderr pipe: %s", err))
			return err
		}
		logger.debug("got stderr pipe")

		logger.debug("starting service")
		err = cmd.Start()
		if err != nil {
			logger.error(fmt.Sprintf("error starting service: %s", err))
			return err
		}

		// Read stdout line by line
		stdoutScanner := bufio.NewScanner(stdout)
		stdoutScanner.Split(bufio.ScanLines)
		stderrScanner := bufio.NewScanner(stderr)

		go func() {
			for stdoutScanner.Scan() {
				logger.info(stdoutScanner.Text())
			}
		}()
		go func() {
			for stderrScanner.Scan() {
				logger.error(stderrScanner.Text())
			}
		}()

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
