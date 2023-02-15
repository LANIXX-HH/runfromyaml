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
	cmds := splitCommands(cmd)
	for _, str := range cmds {
		str = os.ExpandEnv(str)
		if dcommand == "run" {
			docker = []string{"docker", dcommand, "-it", "--rm", container, "sh", "-c", string(str)}

		}
		if dcommand == "exec" {
			docker = []string{"docker", dcommand, container, "sh", "-c", string(str)}

		}
		runCommand(_envs, docker, _level, _output)
	}
	wg.Done()
}

//CommandDockerComposeExec run a command in a docker container with wait parameter and print description to shell
//docker-compose -p $PROJECT -f $_COMPOSE_FILE --project-directory $_PWD exec -u $_PHP_WEB_USER $DOCKER_EXEC_PARAM ${@:1:1} bash -c "${@:2}"
func CommandDockerComposeExec(command string, service string, cmdoptions []string, dcoptions []string, cmd []string, desc string, envs []string, wg *sync.WaitGroup, _level string, _output string) {
	compose := buildDockerComposeCommand(dcoptions, command, cmdoptions, service)
	cmds := splitCommands(cmd)
	for _, _cmd := range cmds {
		var _compose []string
		_cmd = os.ExpandEnv(_cmd)
		onecmd := strings.Fields(_cmd)
		if !reflect.ValueOf(onecmd).IsNil() {
			_compose = append(compose, onecmd...)
		}
		//compose = append(append(append(append(dcoptions, command), cmdoptions...), service), cmd...)
		runCommand(envs, _compose, _level, _output)
	}
	wg.Done()
}

//CommandSSH run a command in a shell with wait parameter and print description to shell
func CommandSSH(user string, port int, host string, options []string, cmd []string, desc string, _envs []string, wg *sync.WaitGroup, _level string, _output string) {
	cmds := splitCommands(cmd)
	for _, sshcmd := range cmds {
		ssh := append(append([]string{"ssh", "-p", strconv.Itoa(port), "-l", user, host}, options...), sshcmd)
		runCommand(_envs, ssh, _level, _output)
	}
	wg.Done()
}

//CommandShell run a command in a shell with wait parameter and print description to shell
func CommandShell(cmd []string, desc string, wg *sync.WaitGroup, index int, _envs []string, _level string, _output string) {
	bash := append([]string{"bash", "-c"}, strings.Join(cmd, " "))
	runCommand(_envs, bash, _level, _output)
	wg.Done()
}

//Command run a commad form string array with wait parameted and print description
func Command(cmd []string, desc string, wg *sync.WaitGroup, _envs []string, _level string, _output string) {
	cmds := splitCommands(cmd)
	for _, _cmd := range cmds {
		onecmd := strings.Split(_cmd, " ")
		runCommand(_envs, onecmd, _level, _output)
	}
	wg.Done()
}

// internal commands

func buildDockerComposeCommand(dcoptions []string, command string, cmdoptions []string, service string) []string {
	var compose []string
	compose = append(compose, "docker-compose")
	if !reflect.ValueOf(dcoptions).IsNil() {
		compose = append(compose, dcoptions...)
	}
	compose = append(compose, command)
	if !reflect.ValueOf(cmdoptions).IsNil() {
		compose = append(compose, cmdoptions...)
	}
	if service != "" {
		compose = append(compose, service)
	}
	return compose
}

func splitCommands(cmd []string) []string {
	temp_cmds := strings.Join(cmd, " ")
	cmds := strings.Split(temp_cmds, ";")
	return cmds
}

func runCommand(envs []string, onecmd []string, _level string, _output string) {
	command := exec.Command(onecmd[0], onecmd[1:]...)
	command.Env = append(os.Environ(), envs...)
	functions.PrintSwitch(color.FgYellow, _level, _output, strings.Trim(fmt.Sprint(onecmd), "[]"), "\n")
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
		err := command.Run()
		if err != nil {
			functions.PrintColor(color.FgRed, "error", "Error: ", err)
		}
	}
}
