package docker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ExecResult is a struct for Inspect Exec Response function
type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

// List docker containers
func List() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Unable to get new docker client: %v", err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Printf("Unable to list containers: %v", err)
	}
	if len(containers) > 0 {
		fmt.Printf("NAME\t\tDATE\n\n")
		for _, container := range containers {
			fmt.Printf("%s\t\t%s\n", container.Names[0], time.Unix(container.Created, 0))
		}
	} else {
		log.Println("There are no containers running")
	}
	return err
}

// Exec a command in docker container
func Exec(ctx context.Context, containerID string, command []string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Unable to get new docker client: %v", err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Printf("Unable to list containers: %v", err)
	}
	if len(containers) > 0 {
		for _, container := range containers {
			fmt.Println(container.Names)
		}
	} else {
		log.Println("There are no containers running")
	}

	return err
	//returnvalue, err := docker.ContainerExecCreate(ctx, containerID, c)
	//return returnvalue, err
}

// //ExecTest is a cool function
// func ExecTest() {
// 	ctx := context.Background()
// 	cli, err := client.NewEnvClient()
// 	if err != nil {
// 		panic(err)
// 	}

// 	reader, err := cli.ImagePull(ctx, "docker.io/library/centos:7", types.ImagePullOptions{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	io.Copy(os.Stdout, reader)

// 	resp, err := cli.ContainerCreate(ctx, &container.Config{
// 		Image: "centos",
// 		//Cmd:   []string{"cat", "/etc/hosts"},
// 		//Tty:   true,
// 	}, nil, nil, "")
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err := cli.ContainerStart(ctx, resp.ID,
// 		types.ContainerStartOptions{}); err != nil {
// 		panic(err)
// 	}

// 	time.Sleep(5 * time.Second)

// 	c := types.ExecConfig{AttachStdout: true, AttachStderr: true,
// 		Cmd: []string{"echo", "Hello Himanshu"}}
// 	execID, _ := cli.ContainerExecCreate(ctx, resp.ID, c)
// 	fmt.Println(execID)

// 	config := types.ExecConfig{}
// 	res, er := cli.ContainerExecAttach(ctx, execID.ID, config)
// 	if er != nil {
// 		fmt.Println("Some Error")
// 	}

// 	err = cli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{})
// 	if err != nil {
// 		fmt.Println("Kuchh error")
// 	}
// 	content, _, _ := res.Reader.ReadLine()
// 	fmt.Println(string(content))

// }

// // InspectExecResp is a function to print output of command triggered from Exec function
// func InspectExecResp(ctx context.Context, id string) (ExecResult, error) {
// 	var execResult ExecResult
// 	docker, err := client.NewEnvClient()
// 	if err != nil {
// 		return execResult, err
// 	}
// 	closer(docker)

// 	resp, err := docker.ContainerExecAttach(ctx, id, types.ExecConfig{})
// 	if err != nil {
// 		return execResult, err
// 	}
// 	response(resp)

// 	// read the output
// 	var outBuf, errBuf bytes.Buffer
// 	outputDone := make(chan error)

// 	go func() {
// 		// StdCopy demultiplexes the stream into two buffers
// 		_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
// 		outputDone <- err
// 	}()

// 	select {
// 	case err := <-outputDone:
// 		if err != nil {
// 			return execResult, err
// 		}
// 		break

// 	case <-ctx.Done():
// 		return execResult, ctx.Err()
// 	}

// 	stdout, err := ioutil.ReadAll(&outBuf)
// 	if err != nil {
// 		return execResult, err
// 	}
// 	stderr, err := ioutil.ReadAll(&errBuf)
// 	if err != nil {
// 		return execResult, err
// 	}

// 	res, err := docker.ContainerExecInspect(ctx, id)
// 	if err != nil {
// 		return execResult, err
// 	}

// 	execResult.ExitCode = res.ExitCode
// 	execResult.StdOut = string(stdout)
// 	execResult.StdErr = string(stderr)
// 	return execResult, nil
// }

// //MyTest is ok
// func MyTest() {
// 	ctx := context.Background()
// 	cli, err := client.NewEnvClient()
// 	if err != nil {
// 		panic(err)
// 	}

// 	reader, err := cli.ImagePull(ctx, "docker.io/library/centos:7", types.ImagePullOptions{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	io.Copy(os.Stdout, reader)

// 	resp, err := cli.ContainerCreate(ctx, &container.Config{
// 		Image: "centos",
// 		//Cmd:   []string{"cat", "/etc/hosts"},
// 		//Tty:   true,
// 	}, nil, nil, "")
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err := cli.ContainerStart(ctx, resp.ID,
// 		types.ContainerStartOptions{}); err != nil {
// 		panic(err)
// 	}

// 	time.Sleep(5 * time.Second)

// 	c := types.ExecConfig{AttachStdout: true, AttachStderr: true,
// 		Cmd: []string{"echo", "Hello Himanshu"}}
// 	execID, _ := cli.ContainerExecCreate(ctx, resp.ID, c)
// 	fmt.Println(execID)

// 	config := types.ExecConfig{}
// 	res, er := cli.ContainerExecAttach(ctx, execID.ID, config)
// 	if er != nil {
// 		fmt.Println("Some Error")
// 	}

// 	err = cli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{})
// 	if err != nil {
// 		fmt.Println("Kuchh error")
// 	}
// 	content, _, _ := res.Reader.ReadLine()
// 	fmt.Println(string(content))

// }

// func closer(docker *client.Client) {
// 	defer docker.Close()
// }

// func response(resp types.HijackedResponse) {
// 	defer resp.Close()
// }
