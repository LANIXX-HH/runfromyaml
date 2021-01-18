package exec

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/ionrock/procs"
)

//CommandDockerRun run a command in a docker container with wait parameter and print description to shell
func CommandDockerRun(dcommand string, container string, cmd []string, desc string, wg *sync.WaitGroup) {
	color.New(color.FgGreen).Println("==> " + desc)
	fmt.Println(cmd)
	var docker []string

	if dcommand == "run" {
		docker = append([]string{"docker", dcommand, "-it", "--rm", container, "bash", "-c"}, strings.Join(cmd, " "))

	}
	if dcommand == "exec" {
		docker = append([]string{"docker", dcommand, container, "bash", "-c"}, strings.Join(cmd, " "))

	}

	command := exec.Command(docker[0], docker[1:]...)
	command.Env = os.Environ()
	color.New(color.FgYellow).Println("Command:", docker, "\n")
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	fmt.Sprintln(docker)

	if err := command.Run(); err != nil {
		log.Fatalf("Start: %v", err)
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
		log.Fatalf("Start: %v", err)
	}
	wg.Done()
}

//CommandShell run a command in a shell with wait parameter and print description to shell
func CommandShell(cmd []string, desc string, wg *sync.WaitGroup) {
	color.New(color.FgGreen).Println("==> " + desc)
	fmt.Println(cmd)
	var bash []string
	bash = append([]string{"bash", "-c"}, strings.Join(cmd, " "))
	command := exec.Command(bash[0], bash[1:]...)
	command.Env = os.Environ()
	color.New(color.FgYellow).Println("Command:", bash, "\n")
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	fmt.Sprintln(bash)

	if err := command.Run(); err != nil {
		log.Fatalf("Start: %v", err)
	}
	wg.Done()
}

//Command run a commad form string array with wait parameted and print description
func Command(cmd []string, desc string, wg *sync.WaitGroup) {
	color.New(color.FgGreen).Println("==> " + desc)
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = os.Environ()
	color.New(color.FgYellow).Println("Command:", cmd, "\n")
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		log.Fatalf("Start: %v", err)
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
