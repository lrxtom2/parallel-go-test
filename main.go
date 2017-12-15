package main

import (
  "bufio"
  "bytes"
  "flag"
  "fmt"
  "io"
  "os"
  "os/exec"
  "strings"
)

func usage() string {
  return `Usage:
  parallel-go-test -f <binary> [-p n]

  Test names must be listed one per line on stdin.
`
}

func main() {
  binaryPath := flag.String("f", "", "file path of test package")
  parallelism := flag.Int("p", 1, "number of tests to execute in parallel")
  flag.Parse()

  if binaryPath == nil || *binaryPath == "" {
    fmt.Fprint(os.Stderr, usage())
    os.Exit(1)
  }

  if _, err := os.Stat(*binaryPath); err != nil {
    fmt.Fprintf(os.Stderr, "Not valid file path: %s\n", *binaryPath)
    os.Exit(1)
  }

  testNames := make([]string, 0, 0)
  stdInReader := bufio.NewReader(os.Stdin)

  for {
    line, err := stdInReader.ReadString('\n')
    if err != nil {
      if err == io.EOF {
        if strings.TrimSpace(line) != "" {
          testNames = append(testNames, line)
        }
        break
      }
      fmt.Fprintf(os.Stderr, "error reading stdin: %s", err)
      os.Exit(1)
    }

    if strings.TrimSpace(line) != "" {
      testNames = append(testNames, line)
    }
  }

  testQueue := make(chan string)
  messages := make(chan string)
  completed := make(chan struct{})

  for i := 0; i < *parallelism; i++ {
    go runWorker(testQueue, messages, completed, *binaryPath)
  }

  go func() {
    for _, testName := range testNames {
      testQueue <- strings.TrimSpace(testName)
    }
  }()

  resultsCount := 0
  for {
    select {
    case message := <-messages:
      fmt.Printf("%s", message)
    case <-completed:
      resultsCount++
    }

    if resultsCount == len(testNames) {
      break
    }
  }
}

func runWorker(inputQueue <-chan string, messages chan<- string,
               done chan<- struct{}, binaryName string) {
  for {
    select {
    case testName := <-inputQueue:
      messages <- runTest(testName, binaryName)
      done <- struct{}{}
    }
  }
}

func runTest(testName string, binaryName string) string {
  var stdResult bytes.Buffer
  cmd := exec.Command(binaryName, "-test.v", "-test.run",
                      fmt.Sprintf("^%s$", testName))
  cmd.Stdout = &stdResult
  cmd.Stderr = &stdResult

  if err := cmd.Run(); err != nil {
    stdResult.WriteString(err.Error())
  }

  return stdResult.String()
}
