// oarb.go - Observe and Report Buddy
//
// Copyright (C) 2023 Anthony Green - green@redhat.com
//

package main

import (
    "bufio"
    "flag"
    "fmt"
    "gopkg.in/yaml.v3"
    "io/ioutil"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "sync"
    // Uncomment or add imports as you implement the respective features.
    //"net/http"
)

type Channel struct {
     Name string    `yaml:"name"`
     Type string    `yaml:"type"`
     Settings map[string]interface{} `yaml:"settings"`
}

type Signal struct {
     Regex string      `yaml:"regex"`
     Channel string    `yaml:"channel"`
}

type Config struct {
     Channel []Channel `yaml:"channels"`
     Signal []Signal `yaml:"signals"`
}

func compileSignals(signals []Signal) ([]CompiledSignal, error) {
    var compiledSignals []CompiledSignal
    for _, signal := range signals {
        re, err := regexp.Compile(signal.Regex)
        if err != nil {
            return nil, fmt.Errorf("failed to compile regex %v: %v", signal.Regex, err)
        }
        compiledSignals = append(compiledSignals, CompiledSignal{
            Regex:   re,
            Channel: signal.Channel,
        })
    }
    return compiledSignals, nil
}

type CompiledSignal struct {
    Regex   *regexp.Regexp
    Channel string
}

type Notification struct {
    Channel     Channel
    Message     string
}

func loadYAMLConfig(filename string, config *Config) error {
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("error reading config file: %v", err)
    }

    err = yaml.Unmarshal(bytes, config)
    if err != nil {
        return fmt.Errorf("error parsing config file: %v", err)
    }

    return nil
}

func startWorkers(notificationQueue <-chan Notification, numWorkers int, wg *sync.WaitGroup) {
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for notification := range notificationQueue {
                // Process the notification (send emails, HTTP requests, etc.)
                fmt.Printf("Worker %d processing notification for channel %s\n", workerID, notification.Channel.Name)
//                triggerNotification(notification.ChannelName, notification.Message)
            }
        }(i)
    }
}

func main() {

     // Define a string flag for the configuration file
    configFilePath := flag.String("c", "", "Path to the oar configuration file")

    // Parse the flags
    flag.Parse()

    config := Config{}
    err := loadYAMLConfig(*configFilePath, &config)
    if err != nil {
       fmt.Println("Failed to load config:", err)
       os.Exit(1)
    }

    for _, channel := range config.Channel {
        fmt.Printf("Loaded channel: %s of type %s with settings: %v\n", channel.Name, channel.Type, channel.Settings)
    }

    for _, signal := range config.Signal {
        fmt.Printf("Loaded signal: %s to %s\n", signal.Regex, signal.Channel)
    }

    compiledSignals, _ := compileSignals(config.Signal)

    // The remaining arguments after flags are parsed
    subprocessArgs := flag.Args()
    if len(subprocessArgs) == 0 {
        fmt.Println("Error: No command provided to run")
        os.Exit(1)
    }

    notificationQueue := make(chan Notification, 100)

    // Use a WaitGroup to wait for the reading goroutines to finish
    var wg sync.WaitGroup
    var nwg sync.WaitGroup

    startWorkers(notificationQueue, 5, &nwg)

    channelMap := make(map[string]Channel)
    for _, ch := range config.Channel {
        channelMap[ch.Name] = ch
    }

    // Prepare to run the subprocess
    cmd := exec.Command(subprocessArgs[0], subprocessArgs[1:]...)
    // Rest of your code to handle the subprocess...
    stdout, _ := cmd.StdoutPipe()
    stderr, _ := cmd.StderrPipe()

    if err := cmd.Start(); err != nil {
        // Handle error
    }

    // Increment WaitGroup and start reading stdout
    wg.Add(2)
    go func() {
        defer wg.Done()
        monitorOutput(bufio.NewScanner(stdout), compiledSignals, notificationQueue, channelMap)
    }()
    go func() {
        defer wg.Done()
        monitorOutput(bufio.NewScanner(stderr), compiledSignals, notificationQueue, channelMap)
    }()

    // Wait for all reading to be complete
    wg.Wait()

    // Wait for the command to finish
    err = cmd.Wait()

    close(notificationQueue)

    // Wait for all reading to be complete
    nwg.Wait()

    // Handle exit status
    if err != nil {
        // cmd.Wait() returns an error if the command exits non-zero
        if exitError, ok := err.(*exec.ExitError); ok {
            // Get the command's exit code
            os.Exit(exitError.ExitCode())
        } else {
            // Other error types (not non-zero exit)
            fmt.Println("Error waiting for Command:", err)
            os.Exit(1)
        }
    } else {
        // Success (exit code 0)
        os.Exit(0)
    }
}

func monitorOutput(scanner *bufio.Scanner, compiledSignals []CompiledSignal, notificationQueue chan Notification, channelMap map[string]Channel) {
    for scanner.Scan() {
        line := scanner.Text()

        for _, signal := range compiledSignals {
            if signal.Regex.MatchString(line) {
                fmt.Printf("Match found for channel %s: %s\n", signal.Channel, line)
                channel, _ := channelMap[signal.Channel]
                notificationQueue <- Notification{Channel: channel, Message: line}
            }
        }

        fmt.Println(line)
    }
    if err := scanner.Err(); err != nil {
       fmt.Fprintf(os.Stderr, "Error reading from pipe: %v\n", err)
    }
}

func makeAPICall(line string) {
    // Implement the REST API call here
}

// loadRegexps loads regexps from a specified file and returns a slice of compiled regexps.
func loadRegexps(filepath string) ([]*regexp.Regexp, error) {
    // Read the file content
    content, err := os.ReadFile(filepath)
    if err != nil {
        return nil, fmt.Errorf("error reading regex file: %w", err)
    }

    // Split the content into lines
    lines := strings.Split(string(content), "\n")

    // Compile each line into a regexp
    var regexps []*regexp.Regexp
    for _, line := range lines {
        // Skip empty lines
        if strings.TrimSpace(line) == "" {
            continue
        }

        re, err := regexp.Compile(line)
        if err != nil {
            return nil, fmt.Errorf("error compiling regexp '%s': %w", line, err)
        }
        regexps = append(regexps, re)
    }

    return regexps, nil
}
