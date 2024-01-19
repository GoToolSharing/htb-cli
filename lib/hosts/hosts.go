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
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' | sudo tee /etc/hosts > /dev/null", newContent))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func AddEntryToHosts(ip string, host string) error {
	ipFound := false
	hostAdded := false

	processLine := func(line string) (string, bool) {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			return line, false
		}

		fields := strings.Fields(trimmedLine)
		if fields[0] == ip {
			ipFound = true
			for _, field := range fields[1:] {
				if field == host {
					return line, false
				}
			}
			return line + " " + host, true
		}
		return line, false
	}

	newContent, changeMade, err := readHostsFile(processLine)
	if err != nil {
		return err
	}

	if !ipFound {
		newContent = strings.TrimSpace(newContent) + "\n" + ip + " " + host
		hostAdded = true
	} else {
		hostAdded = changeMade
	}

	if hostAdded {
		if err := updateHostsFile(strings.TrimSpace(newContent)); err != nil {
			return err
		}
		fmt.Println("Entry successfully updated or added.")
		return nil
	}

	fmt.Println("Entry already exists.")
	return nil
}

func RemoveEntryFromHosts(ip string, host string) error {
	hostRemoved := false

	processLine := func(line string) (string, bool) {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			return line, false
		}

		fields := strings.Fields(trimmedLine)
		if fields[0] == ip {
			var newFields []string
			newFields = append(newFields, ip)

			for _, field := range fields[1:] {
				if field != host {
					newFields = append(newFields, field)
				}
			}

			if len(newFields) > 1 {
				return strings.Join(newFields, " "), true
			}
			return "", true
		}
		return line, false
	}

	newContent, changeMade, err := readHostsFile(processLine)
	if err != nil {
		return err
	}

	if changeMade {
		newContent = strings.TrimSpace(newContent)
		if err := updateHostsFile(newContent); err != nil {
			return err
		}
		fmt.Println("Entry successfully deleted.")
		return nil
	}

	if !hostRemoved {
		fmt.Println("Entry not found.")
	}
	return nil
}
