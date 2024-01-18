package hosts

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const hostsFile = "/etc/hosts"

func readHostsFile(processLine func(string) (string, bool)) (string, bool, error) {
	file, err := os.Open(hostsFile)
	if err != nil {
		return "", false, err
	}
	defer file.Close()

	var buffer bytes.Buffer
	changeMade := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			buffer.WriteString("\n")
			continue
		}

		processedLine, changed := processLine(line)
		if changed {
			changeMade = true
		}
		buffer.WriteString(processedLine + "\n")
	}

	return buffer.String(), changeMade, scanner.Err()
}

func updateHostsFile(newContent string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' | sudo tee /etc/hosts", newContent))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func AddEntryToHosts(ip string, host string) error {
	fmt.Println("IP :", ip)
	fmt.Println("Host :", host)
	processLine := func(line string) (string, bool) {
		fields := strings.Fields(line)
		if fields[0] != ip || strings.Contains(line, host) {
			fmt.Println("Inside :", line)
			return line, false
		}
		return line + " " + host, true
	}

	newContent, changeMade, err := readHostsFile(processLine)
	if err != nil {
		return err
	}

	fmt.Println("New content :", newContent)

	if changeMade {
		if err := updateHostsFile(newContent); err != nil {
			return err
		}
		fmt.Println("Entrée mise à jour ou ajoutée avec succès.")
		return nil
	}

	fmt.Println("L'entrée existe déjà.")
	return nil
}

func RemoveEntryFromHosts(ip string, host string) error {
	processLine := func(line string) (string, bool) {
		fields := strings.Fields(line)
		if fields[0] != ip {
			return line, false
		}
		var newFields []string
		for _, field := range fields {
			if field != host {
				newFields = append(newFields, field)
			}
		}
		return strings.Join(newFields, " "), len(newFields) != len(fields)
	}

	newContent, changeMade, err := readHostsFile(processLine)
	if err != nil {
		return err
	}

	if changeMade {
		if err := updateHostsFile(newContent); err != nil {
			return err
		}
		fmt.Println("Entrée supprimée avec succès.")
		return nil
	}

	fmt.Println("L'entrée n'a pas été trouvée.")
	return nil
}
