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
runfromyaml -f ./other/path/my-collection.yaml
~~~

with more debug output

~~~shell
runfromyaml -f my-collection.yaml -debug
~~~

## Full example based on tooling image setup

~~~shell
curl --silent --location https://raw.githubusercontent.com/LANIXX-HH/runfromyaml/master/examples/tooling.sh | sh
./runfromyaml -f tooling.yaml
~~~

## Syntax

### Logging Settings

all the logging setting should be defined as following example:

~~~yaml
logging:
  - output: stdout
  - level: info
~~~

* output - here you define how the output of the commands to be executed should happen: directly into the standard output or into a log file. without this option no output will be produced. possible values:
  * stdout - direct output into current shell session in stdout
  * file - JSON Output of all to temporary logfile. the logging file will be printed by selecting this option.

* level - logging level: info, warn, error, debug, trace, fatal, panic

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
  * `docker-compose` - this section describes all things that should be executed with docker-compose. why do we need this? experience has shown that it is easier for developers to write a yaml file and specify the necessary options in a similar order than to remember the order of commands to execute. :)
  The required values for this section are:
    * `dcoptions` - global docker-compose options like path directory or docker-compose file(s).
    * `command` - docker-compose command like `run`, `up` or `down`.
    * `cmdoptions` - options needed for the selected command like `-i` or/and `-t`.
    * `service` - name of the service defined in the docker-compose yaml file
    * `values` - commands to be executed within the selected service (when starting the container or in the currently running container)
  * `ssh` t.b.d.
* `name` - this is the name of the section
* `desc` - long description of this section. should contain the really necessary information, what happens in this section.
* `values` - this section generally contains all the steps that should be executed to implement the described workflow. Multiple commands should be separated by `;`.
* `expandenv` - this is an optional switch setting used to enable or disable the environment variable resolution in the corresponding block. default value is disabled

## Examples

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
