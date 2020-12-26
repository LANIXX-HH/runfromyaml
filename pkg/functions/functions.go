package functions

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/ionrock/procs"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func exeCommandWithinBash(cmd []string, desc string, wg *sync.WaitGroup) {
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

func exeCommand(cmd []string, desc string, wg *sync.WaitGroup) {
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

func exeCommandTest(cmd []string, desc string, wg *sync.WaitGroup) {

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
		return line
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

func writeFile(file string, path string, perm os.FileMode) {
	bytefile := []byte(file)
	err := ioutil.WriteFile(path, bytefile, perm)
	check(err)
}

func readFile(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File contents: %s", content)
}

func remove(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}
