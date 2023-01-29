package exec

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	functions "github.com/lanixx/runfromyaml/pkg/functions"
)

//CommandDockerRun run a command in a docker container with wait parameter and print description to shell
func CommandDockerRun(dcommand string, container string, cmd []string, desc string, _envs []string, wg *sync.WaitGroup, _level string, _output string) {
	var docker []string

	temp_cmds := strings.Join(cmd, " ")
	cmds := strings.Split(temp_cmds, ";")

	for _, str := range cmds {
		str = os.ExpandEnv(str)
		if dcommand == "run" {
			docker = []string{"docker", dcommand, "-it", "--rm", container, "sh", "-c", string(str)}

		}
		if dcommand == "exec" {
			docker = []string{"docker", dcommand, container, "sh", "-c", string(str)}

		}

		command := exec.Command(docker[0], docker[1:]...)
		command.Env = append(os.Environ(), _envs...)
		functions.PrintSwitch(color.FgYellow, _level, _output, strings.Trim(fmt.Sprint(docker), "[]"), "\n")
		command.Stdin = os.Stdin
		if _output == "rest" {
			out, err := command.CombinedOutput()
			if err != nil {
				functions.PrintRest(color.FgRed, "error", "Error: ", err, string(out))
			}
			functions.PrintRest(color.FgHiWhite, _level, string(out))
		}

		if _output == "file" {
			out, err := command.CombinedOutput()
			if err != nil {
				functions.PrintFile("error", "Error: ", err, string(out))
			}
			functions.PrintFile(_level, string(out))
		}

		if _output == "stdout" {
			command.Stdout = os.Stdout
			command.Stdin = os.Stdin
			command.Stderr = os.Stderr

			if err := command.Run(); err != nil {
				functions.PrintColor(color.FgRed, "error", "Command: ", command)
				functions.PrintColor(color.FgRed, "error", "Error: ", err)
			}
		}

	}
	wg.Done()
}

//CommandDockerComposeExec run a command in a docker container with wait parameter and print description to shell
//docker-compose -p $PROJECT -f $_COMPOSE_FILE --project-directory $_PWD exec -u $_PHP_WEB_USER $DOCKER_EXEC_PARAM ${@:1:1} bash -c "${@:2}"
func CommandDockerComposeExec(command string, service string, cmdoptions []string, dcoptions []string, cmd []string, desc string, envs []string, wg *sync.WaitGroup, _level string, _output string) {
	var compose []string
	if !reflect.ValueOf(dcoptions).IsNil() {
		compose = append(compose, dcoptions...)
	}
	// command is required
	compose = append(compose, command)
	if !reflect.ValueOf(cmdoptions).IsNil() {
		compose = append(compose, cmdoptions...)
	}
	if service != "" {
		compose = append(compose, service)
	}

	temp_cmds := strings.Join(cmd, " ")
	cmds := strings.Split(temp_cmds, ";")

	for _, _cmd := range cmds {
		var _compose []string
		_cmd = os.ExpandEnv(_cmd)
		onecmd := strings.Fields(_cmd)
		if !reflect.ValueOf(onecmd).IsNil() {
			_compose = append(compose, onecmd...)
		}
		//compose = append(append(append(append(dcoptions, command), cmdoptions...), service), cmd...)
		command := exec.Command("docker-compose", _compose...)
		command.Env = append(os.Environ(), envs...)
		functions.PrintSwitch(color.FgYellow, _level, _output, "docker-compose", strings.Trim(fmt.Sprint(_compose), "[]"), "\n")
		command.Stdin = os.Stdin
		if _output == "rest" {
			command.StdoutPipe()
			out, err := command.CombinedOutput()
			if err != nil {
				functions.PrintRest(color.FgRed, "error", "Error: ", err, string(out))
			}
			functions.PrintRest(color.FgHiWhite, _level, string(out))
		}

		if _output == "file" {
			out, err := command.CombinedOutput()
			if err != nil {
				functions.PrintFile("error", "Error: ", err, string(out))
			}
			functions.PrintFile(_level, string(out))
		}

		if _output == "stdout" {
			command.Stdout = os.Stdout
			command.Stdin = os.Stdin
			command.Stderr = os.Stderr

			if err := command.Run(); err != nil {
				functions.PrintColor(color.FgRed, "error", "Command: ", command)
				functions.PrintColor(color.FgRed, "error", "Error: ", err)
			}
		}
	}
	wg.Done()
}

//CommandSSH run a command in a shell with wait parameter and print description to shell
func CommandSSH(user string, port int, host string, options []string, cmd []string, desc string, _envs []string, wg *sync.WaitGroup, _level string, _output string) {
	temp_cmds := strings.Join(cmd, " ")
	cmds := strings.Split(temp_cmds, ";")
	for _, sshcmd := range cmds {
		ssh := append(append([]string{"ssh", "-p", strconv.Itoa(port), "-l", user, host}, options...), sshcmd)
		command := exec.Command(ssh[0], ssh[1:]...)
		command.Env = append(os.Environ(), _envs...)
		functions.PrintSwitch(color.FgYellow, _level, _output, strings.Trim(fmt.Sprint(ssh), "[]"), "\n")
		command.Stdin = os.Stdin
		if _output == "rest" {
			out, err := command.CombinedOutput()
			if err != nil {
				functions.PrintRest(color.FgRed, "error", "Error: ", err, string(out))
			}
			functions.PrintRest(color.FgHiWhite, _level, string(out))
		}

		if _output == "file" {
			out, err := command.CombinedOutput()
			if err != nil {
				functions.PrintFile("error", "Error: ", err, string(out))
			}
			functions.PrintFile(_level, string(out))
		}

		if _output == "stdout" {
			command.Stdout = os.Stdout
			command.Stdin = os.Stdin
			command.Stderr = os.Stderr

			if err := command.Run(); err != nil {
				functions.PrintColor(color.FgRed, "error", "Command: ", command)
				functions.PrintColor(color.FgRed, "error", "Error: ", err)
			}
		}
	}
	wg.Done()
}

//CommandShell run a command in a shell with wait parameter and print description to shell
func CommandShell(cmd []string, desc string, wg *sync.WaitGroup, index int, _envs []string, _level string, _output string) {
	bash := append([]string{"bash", "-c"}, strings.Join(cmd, " "))
	command := exec.Command(bash[0], bash[1:]...)
	command.Env = append(os.Environ(), _envs...)
	functions.PrintSwitch(color.FgYellow, _level, _output, strings.Trim(fmt.Sprint(bash), "[]"), "\n")
	command.Stdin = os.Stdin
	if _output == "rest" {
		out, err := command.CombinedOutput()
		if err != nil {
			functions.PrintRest(color.FgRed, "error", "Error: ", err, string(out))
		}
		functions.PrintRest(color.FgHiWhite, _level, string(out))
	}

	if _output == "file" {
		out, err := command.CombinedOutput()
		if err != nil {
			functions.PrintFile("error", "Error: ", err, string(out))
		}
		functions.PrintFile(_level, string(out))
	}

	if _output == "stdout" {
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			functions.PrintColor(color.FgRed, "error", "Command: ", command)
			functions.PrintColor(color.FgRed, "error", "Error: ", err)
		}
	}
	wg.Done()
}

//Command run a commad form string array with wait parameted and print description
func Command(cmd []string, desc string, wg *sync.WaitGroup, _envs []string, _level string, _output string) {
	temp_cmds := strings.Join(cmd, " ")
	cmds := strings.Split(temp_cmds, ";")

	for _, _cmd := range cmds {
		onecmd := strings.Split(_cmd, " ")
		command := exec.Command(onecmd[0], onecmd[1:]...)
		command.Env = append(os.Environ(), _envs...)
		functions.PrintSwitch(color.FgYellow, _level, _output, "exec", strings.Trim(fmt.Sprint(command), "[]"), "\n")
		command.Stdin = os.Stdin
		if _output == "rest" {
			out, err := command.CombinedOutput()
			if err != nil {
				functions.PrintRest(color.FgRed, "error", "Error: ", err, string(out))
			}
			functions.PrintRest(color.FgHiWhite, _level, string(out))
		}

		if _output == "file" {
			out, err := command.CombinedOutput()
			if err != nil {
				functions.PrintFile("error", "Error: ", err, string(out))
			}
			functions.PrintFile(_level, string(out))
		}

		if _output == "stdout" {
			command.Stdout = os.Stdout
			command.Stdin = os.Stdin
			command.Stderr = os.Stderr

			if err := command.Run(); err != nil {
				functions.PrintColor(color.FgRed, "error", "Command: ", command)
				functions.PrintColor(color.FgRed, "error", "Error: ", err)
			}
		}
	}
	wg.Done()
}
