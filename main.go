package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
)

// getOpenshiftPlatformType uses the 'oc' command to get the infrastructure platform type.
// It handles errors and returns the platform type as a string.
func getOpenshiftPlatformType() (string, error) {
    // Execute the 'oc' command to get the infrastructure resource in YAML format.
    cmd := exec.Command("oc", "get", "infrastructure", "cluster", "-o", "yaml")
    output, err := cmd.CombinedOutput()
    if err != nil {
        // Improved error message: include the command and its output.
        return "", fmt.Errorf("error running 'oc get infrastructure cluster -o yaml': %v, output: %s", err, output)
    }

    // Convert the output to a string for easier processing.
    outputStr := string(output)

    // Extract the platformType using string manipulation.  This is basic YAML parsing.
    // A proper YAML parser would be more robust, but for this simple task, string
    // manipulation is sufficient and avoids adding a dependency.
    platformType, err := extractPlatformType(outputStr)
    if err != nil {
        return "", err
    }
    return platformType, nil
}

// extractPlatformType extracts the platform type from the YAML output string.
func extractPlatformType(yamlOutput string) (string, error) {
    lines := strings.Split(yamlOutput, "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "platformType:") {
            parts := strings.SplitN(line, ":", 2)
            if len(parts) > 1 {
                platformType := strings.TrimSpace(parts[1])
                if platformType != "" {
                    return platformType, nil
                }
                return "", fmt.Errorf("platformType is empty in the output")

            }
        }
    }
    return "", fmt.Errorf("platformType not found in the output")
}

func main() {
    // Call the function to get the platform type.
    platformType, err := getOpenshiftPlatformType()
    if err != nil {
        // Print the error to standard error and exit with a non-zero exit code.
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Print the platform type to standard output.
    fmt.Println(platformType)
}


