package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type EnvVars interface {
	GenEnvsByFile(fp string) (map[string]string, error)
}

func GenEnvsByFile(fileNames []string) (map[string]map[string]string, error) {
	env_files := make(map[string]map[string]string)
	for _, fileName := range fileNames {
		envVars, err := parseEnvFile(fileName)
		if err != nil {
			log.Printf("Unable to find env vars %v in %s", err, fileName)
		} else {
			env_files[fileName] = envVars
		}
	}
	return env_files, nil
}

func parseEnvFile(fp string) (map[string]string, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	envMap := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			envMap[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return envMap, nil
}
