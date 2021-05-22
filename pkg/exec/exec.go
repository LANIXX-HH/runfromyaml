package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/ionrock/procs"
)

//CommandDockerRun run a command in a docker container with wait parameter and print description to shell
func CommandDockerRun(dcommand string, container string, cmd []string, desc string, wg *sync.WaitGroup) {
	color.New(color.FgGreen).Println("==> " + desc)
	var docker []string

	temp_cmds := strings.Join(cmd, " ")
	cmds := strings.Split(temp_cmds, ";")

	for ind, str := range cmds {
		if dcommand == "run" {
			docker = []string{"docker", dcommand, "-it", "--rm", container, "sh", "-c", string(str)}

		}
		if dcommand == "exec" {
			docker = []string{"docker", dcommand, container, "sh", "-c", string(str)}

		}

		command := exec.Command(docker[0], docker[1:]...)
		command.Env = os.Environ()
		color.New(color.FgYellow).Println("Command (", ind, ") :", docker, "\n")
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
func CommandDockerComposeExec(options []string, cmd []string, desc string, wg *sync.WaitGroup) {
	color.New(color.FgGreen).Println("==> " + desc)
	var compose []string

	compose = append(options, cmd...)
	command := exec.Command("docker-compose", compose...)
	command.Env = os.Environ()
	color.New(color.FgYellow).Println("Command:", "docker-compose", compose, "\n")
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		color.New(color.FgRed).Println("Command: ", command)
		color.New(color.FgRed).Println("Options: ", options)
		color.New(color.FgRed).Println("Error: ", err)
	}
	wg.Done()
}

//CommandSSH run a command in a shell with wait parameter and print description to shell
func CommandSSH(user string, port int, host string, options []string, cmd []string, desc string, wg *sync.WaitGroup) {
	color.New(color.FgGreen).Println("==> " + desc)

	temp_cmds := strings.Join(cmd, " ")
	cmds := strings.Split(temp_cmds, ";")

	for ind, sshcmd := range cmds {
		var ssh []string
		ssh = append(append([]string{"ssh", "-p", strconv.Itoa(port), "-l", user, host}, options...), sshcmd)
		command := exec.Command(ssh[0], ssh[1:]...)
		command.Env = os.Environ()
		color.New(color.FgYellow).Println("Command (", ind, "):", ssh, "\n")
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr
		fmt.Sprintln(ssh)

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
func CommandShell(cmd []string, desc string, wg *sync.WaitGroup, index int) {
	color.New(color.FgGreen).Println("==> " + desc)
	var bash []string
	bash = append([]string{"bash", "-c"}, strings.Join(cmd, " "))
	command := exec.Command(bash[0], bash[1:]...)
	command.Env = os.Environ()
	color.New(color.FgYellow).Println("Command(", index, "):", bash, "\n")
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	fmt.Sprintln(bash)

	if err := command.Run(); err != nil {
		color.New(color.FgRed).Println("Command", command)
		color.New(color.FgRed).Println("Error: ", err)
	}
	wg.Done()
}

//Command run a commad form string array with wait parameted and print description
func Command(cmd []string, desc string, wg *sync.WaitGroup) {
	color.New(color.FgGreen).Println("==> " + desc)

	temp_cmds := strings.Join(cmd, " ")
	cmds := strings.Split(temp_cmds, ";")
	fmt.Println(cmds)

	for ind, onecmds := range cmds {
		onecmd := strings.Split(onecmds, " ")
		command := exec.Command(onecmd[0], onecmd[1:]...)
		command.Env = os.Environ()
		color.New(color.FgYellow).Println("Command(", ind, "):", command, "\n")
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

// ExeCommandTest is a test exec function
func CommandTest(cmd []string, desc string, wg *sync.WaitGroup) {

	// define command set
	cmds := []*exec.Cmd{
		exec.Command(cmd[0], cmd[1:]...),
	}

	// init procs with command set
	p := procs.Process{Cmds: cmds}

	// parse environment variables
	env := procs.ParseEnv(os.Environ())
	p.Env = env

	// prepare output handler
	p.OutputHandler = func(line string) string {
		color.New(color.FgGreen).Println("==> " + desc)
		color.New(color.FgYellow).Println("Command: ", cmd)
		return line
	}

	// prepare error handler
	p.ErrHandler = func(line string) string {
		color.New(color.FgRed).Println("Command: ", cmd)
		fmt.Println(cmds)
		fmt.Println(p)
		fmt.Println(env)
		os.Exit(0)
		return line
	}

	color.New(color.FgGreen).Println("==> " + desc)
	color.New(color.FgYellow).Println("Command: ", cmd)

	p.Run()
	p.Wait()

	color.New(color.FgBlue).Println(cmds)
	color.New(color.FgBlue).Println(p)
	color.New(color.FgBlue).Println(env)

	out, _ := p.Output()
	fmt.Printf(string(out))
	wg.Done()
}
