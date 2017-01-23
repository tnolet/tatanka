package work

import (
	"log"
	"os/exec"
	"strings"
)

func runNoopScraper(done chan bool) {

	cmd := exec.Command("/usr/local/bin/docker", "run", "busybox", "sleep", "5")
	err := cmd.Run()

	if err != nil {
		log.Println(err.Error())
	}

	err = cmd.Wait()
	done <- true
}

func StopNoopScraper() error {

	log.Println("Stopping NoopScraper")

	out, err := exec.Command(
		"/usr/local/bin/docker",
		"ps",
		"-q",
	).Output()

	containerIDs := strings.Split(string(out), "\n")

	for _, containerID := range containerIDs {

		log.Println("Stopping container:", containerID)
		out, err = exec.Command(
			"/usr/local/bin/docker",
			"kill",
			containerID,
		).Output()

		if err != nil {
			log.Printf("Docker command unsuccessful: %v, %v", err.Error(), out)
			return err
		}

		log.Printf("Succesfully stopped container: %v", containerID)

	}

	return nil
}

func runScraper(item WorkItem, done chan bool) {

	log.Printf("Starting container...")

	// scrapy crawl abc -a subject="My Subject" -a query="My+Subject"

	crawl := "crawl"
	site := "abc"
	arg := "-a"
	subject := "subject='" + "subject" + "'"
	query := "query=subject"
	cmd := exec.Command(
		"/usr/local/bin/docker",
		"run",
		"-d",
		container,
		crawl,
		site,
		arg,
		subject,
		arg,
		query)

	err := cmd.Run()

	if err != nil {
		log.Printf("Docker command unsuccessful: %v", err.Error())
	}
	cmd.Wait()
}

func stopScraper(containerID string) error {

	log.Printf("Stopping container: %v", containerID)

	out, err := exec.Command(
		"/usr/local/bin/docker",
		"kill",
		containerID,
	).Output()

	if err != nil {
		log.Printf("Docker command unsuccessful: %v, %v", err.Error(), out)
		return err
	}

	log.Printf("Succesfully stopped container: %v", string(out))

	return nil
}
