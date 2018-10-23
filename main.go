package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const ControlPrefix = "#$"

var (
	Usage         = "usage: " + os.Args[0] + " <script>"
	ErrUnknownCmd = errors.New("unknown command")
	ErrNoArgs     = errors.New("no arguments given to command")
	ErrBadArg     = errors.New("invalid command argument")
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		log.Fatal(Usage)
	}
	if exec.Command("asciinema", "-h").Run() != nil {
		log.Fatal("asciinema is not installed")
	}

	s, err := NewScript(os.Args[1])
	if err != nil {
		log.Fatal("parsing script failed:", err)
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

// Command is an action to be run.
type Command interface {
	Run(*Script) error
}

// Shell is a shell command to execute.
type Shell struct {
	Cmd string
}

// NewShell creates a new Shell.
func NewShell(cmd string) Shell {
	return Shell{Cmd: cmd}
}

// Run runs the shell command.
func (s Shell) Run(sc *Script) error {
	for _, c := range s.Cmd {
		fmt.Printf("%s", string(c))
		time.Sleep(sc.TypingDelay)
	}
	return nil
}

// Wait is a command to wait for some time.
type Wait struct {
	Duration time.Duration
}

// NewWait creates a new Wait.
func NewWait(opts []string) (Wait, error) {
	if len(opts) == 0 {
		return Wait{}, ErrNoArgs
	}

	ms, err := strconv.ParseInt(opts[0], 10, 64)
	if err != nil {
		return Wait{}, ErrBadArg
	}

	return Wait{Duration: time.Millisecond * time.Duration(ms)}, nil
}

// Run waits for the specified amount of time.
func (w Wait) Run(*Script) error {
	time.Sleep(w.Duration)
	return nil
}

// Delay is a command to change the typing speed of other commands.
type Delay struct {
	Interval time.Duration
}

// NewDelay creates a new Delay.
func NewDelay(opts []string) (Delay, error) {
	if len(opts) == 0 {
		return Delay{}, ErrNoArgs
	}

	ms, err := strconv.ParseInt(opts[0], 10, 64)
	if err != nil {
		return Delay{}, ErrBadArg
	}

	return Delay{Interval: time.Millisecond * time.Duration(ms)}, nil
}

// Run changes the speed for subsequent commands.
func (s Delay) Run(sc *Script) error {
	sc.TypingDelay = s.Interval
	return nil
}

// NewControlCommand creates a new control command.
func NewControlCommand(cmd string) (Command, error) {
	tokens := strings.Split(cmd, ":")
	switch tokens[0] {
	case "delay":
		return NewDelay(tokens[1:])
	case "wait":
		return NewWait(tokens[1:])
	default:
		return nil, ErrUnknownCmd
	}
}

// Script is a shell script to be run and recorded by asciinema.
type Script struct {
	Name        string
	Commands    []Command
	TypingDelay time.Duration
}

// NewScript parses a new Script from the file at path.
func NewScript(path string) (*Script, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	s := &Script{TypingDelay: time.Millisecond * 25}
	s.Name, _ = filepath.Abs(path)

	lines := strings.Split(string(b), "\n")
	curCmd := ""
	for i, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, ControlPrefix) {
			if curCmd != "" {
				s.Commands = append(s.Commands, NewShell(curCmd))
				curCmd = ""
			}
			ctrl, err := NewControlCommand(line[2:])
			if err != nil {
				return nil, fmt.Errorf("%v (line %d)", err, i)
			}
			s.Commands = append(s.Commands, ctrl)
		} else {
			curCmd += line + "\n"
		}
	}
	if curCmd != "" {
		s.Commands = append(s.Commands, NewShell(curCmd[:len(curCmd)-1]))
	}

	return s, nil
}

// Run runs the script.
func (s *Script) Run() error {
	for i, c := range s.Commands {
		if err := c.Run(s); err != nil {
			return fmt.Errorf("%v (command %d)", err, i)
		}
	}
	return nil
}
