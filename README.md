# What is the goal of the project?

Actually it's a playground and an attempt to write a tool with which I can both create documentation of the steps automatically, record all the necessary configurations and all the necessary commands that must be executed to achieve the goal.

# Why didn't I go with ansible?

My goal was to rewrite an existing bash tool in go. An essential part was a collection of all commands that need to be executed both inside the Docker container and outside the container to build the project. At the time, ansible was too complicated for that. This project is only one part, which only reads the commands from YAML and can execute them both outside and inside the container. Also conceivable are commands via SSH (Todo).

# HowTo build

~~~shell
go mod download
go build
go install runfromyaml
~~~

# HowTo execute

simple run pick commands.yaml in current directory and run all defined commands from this yaml file with descriptions

~~~shell
runfromyaml
~~~

you can select different yaml command collection with -file=_myfile.yaml_ 

~~~shell
runfromyaml -file=./other/path/my-collection.yaml
~~~

with more debug output

~~~shell
runfromyaml -file=my-collection.yaml -debug
~~~
