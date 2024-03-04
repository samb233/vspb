package vspb

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	MAKE_TOOL_MESON = iota + 1
	MAKE_TOOL_CMAKE
	MAKE_TOOL_AUTOMAKE
	MAKE_TOOL_AUTOMAKE_GEN
)

func MatchMakeTool(dir string) (int, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return -1, err
	}

	builder := 0
	for _, file := range files {
		switch file.Name() {
		case "meson.build":
			if builder > MAKE_TOOL_MESON {
				builder = MAKE_TOOL_MESON
			}
		case "CMakeLists.txt":
			if builder > MAKE_TOOL_CMAKE {
				builder = MAKE_TOOL_CMAKE
			}
		case "configure":
			if builder > MAKE_TOOL_AUTOMAKE {
				builder = MAKE_TOOL_AUTOMAKE
			}
		case "autogen.sh":
			if builder > MAKE_TOOL_AUTOMAKE_GEN {
				builder = MAKE_TOOL_AUTOMAKE_GEN
			}
		default:
			continue
		}
	}

	return builder, nil
}

type CmdRunner struct {
	Env  map[string]string
	Dir  string
	Cmds []string

	Verbose bool
}

func (c *CmdRunner) SetEnv(key, value string) {
	c.Env[key] = value
}

func (cr *CmdRunner) Run() error {
	// TODO: choose a shell
	shell := "bash"

	fmt.Println("exec cmds: ")
	for _, cmd := range cr.Cmds {
		fmt.Printf("  %s", cmd)

		// check if a command is a cd command
		// if it is, change the dirs
		if strings.HasPrefix(cmd, "cd") {
			args := strings.SplitN(cmd, " ", 2)
			if len(args) != 2 {
				return fmt.Errorf("'cd' error: '%s'", cmd)
			}
			cr.Dir = args[1]
			continue
		}

		c := exec.Command(shell, "-c", cmd)
		c.Env = os.Environ()
		for key, value := range cr.Env {
			c.Env = append(c.Env, fmt.Sprintf("%s=%s", key, value))
		}
		c.Dir = cr.Dir

		if cr.Verbose {
			fmt.Println(":")
			c.Stdout = os.Stdout

			// TODO:
			// maybe should collect stderr even it is not verbose
			c.Stderr = os.Stderr
		}

		if err := c.Run(); err != nil {
			return fmt.Errorf("Exec error: %w", err)
		}

		if !cr.Verbose {
			fmt.Println(" -OK")
		}
	}

	return nil
}

func DefaultMaker(maketool int) []string {
	var cmds []string
	if maketool == MAKE_TOOL_MESON {
		cmds = []string{
			"meson setup build",
			"ninja -C build",
		}
	}

	if maketool == MAKE_TOOL_CMAKE {
		cmds = []string{
			"cmake -S . -B build",
			"cmake --build build",
		}
	}

	if maketool == MAKE_TOOL_AUTOMAKE {
		cmds = []string{
			"mkdir build",
			"cd build",
			"../configure",
			"make",
		}
	}

	if maketool == MAKE_TOOL_AUTOMAKE_GEN {
		cmds = []string{
			"mkdir build",
			"cd build",
			"../autogen.sh",
			"../configure",
			"make",
		}
	}

	return cmds
}
