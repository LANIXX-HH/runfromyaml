# RUN FROM YAML

## What is the goal of the project?

Actually it's a playground and an attempt to write a tool with which I can both create documentation of the steps automatically, record all the necessary configurations and all the necessary commands that must be executed to achieve the goal.

Мy main goal of this project is to write a utility with which it would be convenient and easy to collect documentation, configuration files and the commands themselves under one roof in the right sequence. basically, the way you try to achieve your goal. this is the most efficient way for this task and i have not found a sensible utility that does this.

## Why didn't I go with ansible?

Very heavyweight for such a task in my opinion. I have met people who saw everything in ansible and then only they understand how it works. very quickly the focus is lost - to write down the steps you just did quite simply and quickly and preferably document what is going on here in general.

My goal was to rewrite an existing bash tool in go. An essential part was a collection of all commands that need to be executed both inside the Docker container and outside the container to build the project. At the time, ansible was too complicated for that. This project is only one part, which only reads the commands from YAML and can execute them both outside and inside the container. Also conceivable are commands via SSH (Todo).

## Current state of Project

At the moment I am testing this on my android phone, on my windows machine, on my mec and on my linux server and it is working pretty well so far.

## TODO's

- [ ] write tests !!
- [ ] implement connection between blocks (artifacts or other way. i don't know)
- [ ] implement dependency between blocks

## HowTo build

~~~shell
make clean && make
~~~

## How To get the binary

~~~shell
curl --silent --location "https://github.com/lanixx-hh/runfromyaml/releases/latest/download/runfromyaml-$(uname -s)-$(uname -m).tar.gz" | tar xz
~~~

## HowTo execute

simple run pick commands.yaml in current directory and run all defined commands from this yaml file with descriptions

~~~shell
runfromyaml
~~~

you can select different yaml command collection with -file=_myfile.yaml_

~~~shell
runfromyaml --file ./other/path/my-collection.yaml
~~~

with more debug output

~~~shell
runfromyaml --file my-collection.yaml -debug
~~~

## Full example based on tooling image setup

~~~shell
curl --silent --location https://raw.githubusercontent.com/LANIXX-HH/runfromyaml/master/examples/tooling.sh | sh
./runfromyaml --file tooling.yaml
~~~

## Options

~~~bash
Usage of ./runfromyaml:
  --debug
    	debug - activate debug mode to print more informations
  --file string
    	file - file with all defined commands, descriptions and configuration blocks in yaml fromat (default "commands.yaml")
  --host string
    	host - set host for rest api mode (default host is localhost) (default "localhost")
  --no-auth
    	no-auth - disable rest auth
  --port int
    	port - set http port for rest api mode (default http port is 8080) (default 8080)
  --rest
    	restapi - start this instance in background mode in rest api mode
  --restout
    	rest output - activate output to http response
  --user string
    	user - set username for rest api authentication (default username is rest)  (default "rest")  
~~~

### Examples

* parse yaml file localy

~~~bash

~~~

* REST API Mode

~~~bash
runfromyaml --rest
~~~

* REST API Mode without Authentication ( !!! CAUTION: Do not use it in public networks !!! )

~~~bash
runfromyaml --rest --no-auth
~~~

* REST API Mode with redirected output to http response

~~~bash
runfromyaml --rest --restout
~~~

* Example CURL Call for REST API Mode

~~~bash 
PASS=<rest_api_generated_password>
CURLOPT_TIMEOUT=30 curl -X POST -H "Content-Type: application/x-yaml" -u rest:$PASS --data-binary @examples/windows.yaml http://192.168.0.100:8000/
~~~

## Syntax

### Logging Settings

all the logging setting should be defined as following example:

~~~yaml
logging:
  - output: stdout
  - level: info
~~~

  * `level` - the following levels are possible: info, warn, debug, error, trace, fatal, panic
  * `output` - define how the output should happen
    * NIL (nothing was set. missing output option) - it nothing is defined, no output will be created :)
    * `stdout` - should be default output
    * `file` - all the output will be redirected to json logfile (implemented with logrus module) in the current temp directory. by start of this program the logging json file will be shown. 
    * `rest` - this payload should be delivered only via http post request as YAML. by default, if the programm is running in rest api mode, output will be overwritten to `rest`

### Environment Variables

all the environment variables should be defined as following example:

~~~yaml
env:
  - key: FOO
    value: bar
    ...
~~~

### CMD Blocks

all the commands & configurations should be written as following example:

~~~yaml
cmd:
  - type: conf
    confdest: $HOME/.aws/config
    confperm: 0644
    confdata: |
      [default]
      region = eu-central-1
      sso_start_url = https://mycompany.awsapps.com/start#/
      sso_region = eu-central-1
      sso_account_id = 123456789000
      sso_role_name = MyPowerRole
  - type: exec
    expandenv: true
    name: first
    desc: first command
    values:
      - /bin/cat $HOME/.aws/config
  - type: shell
    name: aws
    desc: show s3 buckets
    values:
      - aws sso login;
      - aws s3 ls
~~~

every section should begin with `-`

#### CMD Syntax

* `type` - this section describes the type of the current section. possible values:
  * `exec` - executes all commands described in the Values section.
  * `conf` - creates a configuration file with permissions under a specified path. required values are:
    * `confdest` - destination path for this configuration.
    * `confperm` - permissions (in Unix format. e.g.: `0644`) to save this file
    * `confdata` - contained data for the current configuration block
  * `shell` - this section defines a set of commands to be executed in a bash session
  * `docker` - this section defines a set of commands to be executed in a started and/or running container. required values for this section are:
    * `container` - name of the running container where all commands should be executed.
    * `command` - this section contains 2 different options how the proposed command should be executed. required values:
      * `run` - start all commands in the new container
      * `exec` - execute the command in the currently running container
    * `values` - this section defines all the commands to be executed in a started or running container
  * `docker-compose` - this section describes all things that should be executed with docker-compose. why do we need this? experience has shown that it is easier for developers to write a yaml file and specify the necessary options in a similar order than to remember the order of commands to execute. :) you can skip settings (global and command) with empty map like `[]` and set empty service name with `""`. to run a command inside of container, you should define values. for multiple command just separate it with semicolon (`;`)
  The required values for this section are:
    * `dcoptions` - global docker-compose options like path directory or docker-compose file(s).
    * `command` - docker-compose command like `run`, `up` or `down`.
    * `cmdoptions` - options needed for the selected command like `-i` or/and `-t`.
    * `service` - name of the service defined in the docker-compose yaml file
    * `values` - commands to be executed within the selected service (when starting the container or in the currently running container)
  * `ssh` - in this section you can run a remote command on specified host via SSH Connection
    * `user` - username for SSH Connection
    * `host` - hostname for SSH Connection
    * `port` - ssh port for SSH Connection
    * `options` - additional options for SSH Connection like `-i <path/to/ssh/public_key>`
    * `values` - set of commands separated by semicolon (`;`) which should be executed on remote host via SSH Connection
* `name` - this is the name of the section
* `desc` - long description of this section. should contain the really necessary information, what happens in this section.
* `values` - this section generally contains all the steps that should be executed to implement the described workflow. Multiple commands should be separated by `;`.
* `expandenv` - this is an optional switch setting used to enable or disable the environment variable resolution in the corresponding block. default value is disabled

## Examples

### Set logging options

* logging - with this you can set how the output should happen and with which log level that should take place

~~~yaml
logging:
  - level: info
  - output: stdout
~~~

### Define environment variables

* env - in this block variables are set, the other blocks both in the data area and in the setting in the respective block, such as path.

~~~yaml
env:
  - key: LC_ALL
    value: POSIX
  - key: PATH
    value: /usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin:/opt/homebrew/bin
~~~
  
### Create configuration file

* confdata - here the content of the configuration file is written in.

~~~yaml
--
cmd:
  - type: "conf"
    confdata: | 
      test
        test
          test
      test
    confdest: /etc/myconfig.conf
    confperm: 0644
~~~

### Call a OS Command with EXEC

* exec - with this you can start an OS system call with exec 

~~~yaml
  - type: "exec"
    expandenv: true
    name: "envprint"
    desc: "show the output of %COMSPEC% variable on windows system"
    values:
      - cmd
      - /C
      - echo %COMSPEC%
~~~

### Call a Shell Command

* shell - This gives the possibility to call a command that is wrapped in bash. for example: `bash -c 'ls -lisa'`
  
~~~yaml
  - type: "shell"
    expandenv: true
    name: "list"
    desc: "list all files in current directory"
    values:
      - ls -lisa
~~~

### Run a Command via SSH connection on remote server
  
* ssh - in this section you can define a block with username, hostname, port and additional options to run a command set remotely via this SSH Connection

~~~yaml
  - type: "ssh"
    expandenv: true
    name: "ssh-run"
    desc: "run ls command via ssh connection"
    user: $USER
    host: localhost
    port: 22
    options:
      - -i $HOME/.ssh/id_rsa-localhost
    values:
      - uname
      - -a ;
      - pwd****
~~~

### Docker-Compose Block

* docker-compose - Here you have the possibility to compose a complete docker-compose command as a YAML structure with global docker-compose options and specific command options and optionally execute there specific collection of commands separated by semicolon.

~~~yaml
  - type: "docker-compose"
    expandenv: true
    name: "build"
    desc: "build tooling"
    dcoptions:
      - -f $HOME/.tmp/tooling/docker-compose.yaml
      - --project-directory $HOME/.tmp/tooling
    cmdoptions: []
    command: build
    service: ""
    values: []
  - type: "docker-compose"
    expandenv: true
    name: "run"
    desc: "run tooling"
    dcoptions:
      - -f $HOME/.tmp/tooling/docker-compose.yaml
      - --project-directory $HOME/.tmp/tooling
    command: run
    service: tooling
    cmdoptions: []
    values:
      - zsh
~~~
  

  
### Call a Command inside a running Docker container or run it once.

* docker - this section can be used to start some command or set of multiple command separated by semicolon in a running container or by starting new container and terminate it after run.

attention: you can use only `run` or `exec` command

~~~yaml
  - type: docker
    expandenv: true
    desc: "run command from docker container"
    name: "docker-run"
    command: run
    container: alpine
    values:
        - uname
        - -a;
        - pwd
~~~
