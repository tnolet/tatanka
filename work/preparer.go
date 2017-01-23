package work

import (
// "log"
// "os/exec"
)

/*
  Executes arbitrary stuff to prepare for the actual work,
  for instance installing packages or containers.
*/
func Prepare() error {

	// log.Println("Preparing work...")

	// commands := [][]string{
	// 	[]string{"yum", "update", "-y"},
	// 	[]string{"yum", "install", "-y", "docker"},
	// 	[]string{"service", "docker", "start"},
	// 	[]string{"docker", "pull", "tnolet/scraper:0.1.0"},
	// }

	// for _, command := range commands {
	// 	out, err := exec.Command(command).Output()

	// 	if err != nil {
	// 		log.Printf("Prepare command unsuccessful: %v, %v", err.Error(), out)
	// 		return err
	// 	}

	// 	log.Printf("Succesfully executed preparation: %v", out)
	// }
	return nil

}
