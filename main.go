package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type TestCase struct {
	input  []string
	output []string
}

func getTests(fc string) []TestCase {
	var buffer []byte
	bufLen := 0
	newlines := 0

	testCases := make([]TestCase, 0)

	var put []string

	// true == take input, false == take output
	inPhase := true

	for i := range len(fc) {
		if fc[i] == '\n' {
			newlines++
			if newlines < 2 {
				put = append(put, string(buffer))
			}
			buffer = make([]byte, 0)
		} else {
			newlines = 0
			buffer = append(buffer, fc[i])
			bufLen++
		}

		if newlines == 2 {
			if inPhase {
				testCases = append(testCases, TestCase{input: put})
			} else {
				testCases[len(testCases)-1].output = put
			}
			inPhase = !inPhase
			put = make([]string, 0)
		}
	}

	return testCases
}

func start(exec *exec.Cmd) (stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	stdin, err = exec.StdinPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	stdout, err = exec.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	stderr, err = exec.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	if err := exec.Start(); err != nil {
		return nil, nil, nil, err
	}

	return
}

const red string = "\033[31m"
const reset string = "\033[0m"

func printlnRed(a string) {
	if a == "" {
		return
	}
	fmt.Print(red)
	fmt.Println(a)
	fmt.Print(reset)
}

func getArgs() (command []string, program string) {
	if len(os.Args) < 3 {
		printlnRed("Usage: go run main.go <command> <program>")
		os.Exit(1)
	}
	programFilename := os.Args[len(os.Args)-1]
	program = programFilename[:strings.Index(programFilename, ".")]
	return os.Args[1:len(os.Args)], program
}

func main() {
	command, program := getArgs()

	// open a file from the filename from call argument
	fc, err := os.ReadFile(program + ".txt")
	if err != nil {
		panic(err)
	}

	tests := getTests(string(fc))

	wg := sync.WaitGroup{}

	for i, test := range tests {
		wg.Add(1)
		func() {
			defer wg.Done()
			exec := exec.Command(command[0], command[1:]...)

			stdin, stdout, stderr, err := start(exec)
			if err != nil {
				panic(err)
			}
			defer stdin.Close()

			fmt.Printf("=== Test %d ===", i+1)

			for _, input := range test.input {
				stdin.Write([]byte(input + "\n"))
			}

			sout, err := io.ReadAll(stdout)
			serr, err := io.ReadAll(stderr)
			printlnRed(string(serr))

			soutLines := strings.Split(string(sout), "\n")

			bad := false
			for i, line := range test.output {
				if i >= len(soutLines) {
					bad = true
					fmt.Println()
					printlnRed("Output is too short")
				}
				if soutLines[i] != line {
					bad = true
					fmt.Println()
					printlnRed("Output does not match")
					printlnRed("Expected: " + line)
					printlnRed("Got: " + soutLines[i])

				}
			}
			if bad {
				fmt.Println("input:")
				for _, input := range test.input {
					printlnRed(input)
				}
				fmt.Println("output:")
				printlnRed(string(sout))
				fmt.Println()
			} else {
				fmt.Println(" ok")
			}
		}()
	}

	wg.Wait()

	// read the file into stdout
}
