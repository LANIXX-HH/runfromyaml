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
)

//CommandDockerRun run a command in a docker container with wait parameter and print description to shell
func CommandDockerRun(dcommand string, container string, cmd []string, desc string, _envs []string, wg *sync.WaitGroup) {
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
		color.New(color.FgYellow).Print(strings.Trim(fmt.Sprint(docker), "[]"), "\n")
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr

		fmt.Sprintln(docker)

		if err := command.Run(); err != nil {
			fmt.Println("Command: ", command)
			fmt.Println("Docker command: ", dcommand)
			fmt.Println("Container: ", container)
			fmt.Println("Error: ", err)
		}
	}
	wg.Done()
}

//CommandDockerComposeExec run a command in a docker container with wait parameter and print description to shell
//docker-compose -p $PROJECT -f $_COMPOSE_FILE --project-directory $_PWD exec -u $_PHP_WEB_USER $DOCKER_EXEC_PARAM ${@:1:1} bash -c "${@:2}"
func CommandDockerComposeExec(command string, service string, cmdoptions []string, dcoptions []string, cmd []string, desc string, envs []string, wg *sync.WaitGroup) {
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
		cmds := exec.Command("docker-compose", _compose...)
		cmds.Env = append(os.Environ(), envs...)
		color.New(color.FgYellow).Print("docker-compose", strings.Trim(fmt.Sprint(_compose), "[]"), "\n")
		cmds.Stdout = os.Stdout
		cmds.Stdin = os.Stdin
		cmds.Stderr = os.Stderr
		if err := cmds.Run(); err != nil {
			color.New(color.FgRed).Println("Command: ", command)
			color.New(color.FgRed).Println("Service: ", service)
			color.New(color.FgRed).Println("Docker Compose Options: ", dcoptions)
			color.New(color.FgRed).Println("Command Options: ", cmdoptions)
			color.New(color.FgRed).Println("Values: ", onecmd)
			color.New(color.FgRed).Println("Full: ", _compose)
			color.New(color.FgRed).Println("Error: ", err)
		}
	}
	wg.Done()
}

//CommandSSH run a command in a shell with wait parameter and print description to shell
func CommandSSH(user string, port int, host string, options []string, cmd []string, desc string, _envs []string, wg *sync.WaitGroup) {
	var ssh []string
	temp_cmds := strings.Join(cmd, " ")
	cmds := strings.Split(temp_cmds, ";")
	for _, sshcmd := range cmds {
		ssh = append(append([]string{"ssh", "-p", strconv.Itoa(port), "-l", user, host}, options...), sshcmd)
		command := exec.Command(ssh[0], ssh[1:]...)
		command.Env = append(os.Environ(), _envs...)
		color.New(color.FgYellow).Print(strings.Trim(fmt.Sprint(ssh), "[]"), "\n")
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			color.New(color.FgRed).Println("Command: ", command)
			color.New(color.FgRed).Println("User: ", user)
			color.New(color.FgRed).Println("Host: ", host)
			color.New(color.FgRed).Println("Port: ", port)
			color.New(color.FgRed).Println("Error: ", err)
		}
	}
	wg.Done()
}

//CommandShell run a command in a shell with wait parameter and print description to shell
func CommandShell(cmd []string, desc string, wg *sync.WaitGroup, index int, _envs []string) {
	var bash []string
	bash = append([]string{"bash", "-c"}, strings.Join(cmd, " "))
	command := exec.Command(bash[0], bash[1:]...)
	command.Env = append(os.Environ(), _envs...)
	color.New(color.FgYellow).Print(strings.Trim(fmt.Sprint(bash), "[]"), "\n")
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		color.New(color.FgRed).Println("Command", command)
		color.New(color.FgRed).Println("Error: ", err)
	}
	wg.Done()
}

//Command run a commad form string array with wait parameted and print description
func Command(cmd []string, desc string, wg *sync.WaitGroup, _envs []string) {
	temp_cmds := strings.Join(cmd, " ")
	cmds := strings.Split(temp_cmds, ";")

	for _, _cmd := range cmds {
		onecmd := strings.Split(_cmd, " ")
		command := exec.Command(onecmd[0], onecmd[1:]...)
		command.Env = append(os.Environ(), _envs...)
		color.New(color.FgYellow).Print("exec", strings.Trim(fmt.Sprint(command), "[]"), "\n")
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			color.New(color.FgRed).Println("Command", command)
			color.New(color.FgRed).Println("Error: ", err)
		}
	}
	wg.Done()
}
