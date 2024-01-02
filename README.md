# `pmon3`
`pmon3` is a process manager for Golang applications with a built-in load balancer. It allows you to keep applications alive forever and to reload them without downtime.

<img src="http://p0.qhimg.com/t017d6cbb68aed4b693.png" style="max-width:680px" />

## Start Process

```go
sudo pmon3 run [./application binary] [arg1]  ...
```

## Introduction

Golang currently has no officially support process management tools. For the deployment of `Go` services we use the Linux built-in command `nohup [process] &`  combination, or use the process management tools provided by the operating system such as SystemD. Alternatively, third-party process management tools such as: Python's Supervisor or Nodejs PM2 can also be utilized

Each method has certain advantages and disadvantages. We hope to inherit the convenient and easy-to-use ideas of the GO language deployment. There is no need to install other dependencies besides `bash-completion` for ease of command line utilization.

Unlike PM2, `pmon3` is managed directly by the OS process manager, so even if the `pmon3` management tool abnormally terminates, it will not affect the parent `pmon3` process itself. This is currently achieved by seperating the `pmond` deamon process from the `pmon3` agent.

By default, if the `pmon3` agent abnormally terminates, `pmond` will try to restart the process. If you don't want a process to restart automatically, then you can provide the `--no-autorestart` parameter flag.


## How To Install

[Releases](https://github.com/ntt360/pmon3/releases) 

### Using Makefile
The systemd installation process entails the following steps:
1. create the log, configuration and database directories
1. create the log rotation file
1. create the bash completion profile (requires the bash-completion package)
1. enable and start the `pmond` system service

```shell
make systemd_install
```

### RPM packages
These are compatible with systemd-based Linux distributions.

```bash
sudo dnf install -y https://github.com/joe-at-startupmedia/pmon3/releases/download/v1.13.0/pmon3-1.13.0-1.el9.x86_64.rpm
```

:exclamation::exclamation: Note :exclamation::exclamation:

After installing `pmon3` for the first time, the `pmond` service does not start automatically. You need to manually start the service:

```shell
sudo systemctl start pmond

# Others
sudo /usr/local/pmon3/bin/pmond &
```

## Command Introduction

#### Help

```shell
# View global help documentation
sudo pmon3 help

# View a specific command help
sudo pmon3 [command] help
```

global help documentation provides the following output:

```
Usage:
  pmon3 [command]

Available Commands:
  del         del process by id or name
  desc        print the process detail message
  exec        run one binary golang process file
  help        Help about any command
  ls          list all processes
  reload      reload some process
  start       start some process by id or name
  stop        stop running process
  log         display process log by id or name
  logf        display process log dynamic by id or name
  version     show `pmon3` version
```

#### Running process [run/exec]

```shell
sudo pmon3 run [./application binary] [arg1] [arg2] ...
```
The starting process accepts several parameters. The parameter descriptions are as follows:

```shell
// The process name. It will use the file name of the binary as the default name if not provided 
--name

// Where to store logs. It will use the following as the default: /var/log/pmon3/
--log   -l

// Only custom log directory, priority is lower than --log
--log_dir -d

// Provide parameters to be passed to the binary, multiple parameters are separated by spaces
--args  -a "-arg1=val1 -arg2=val2"

// managing user
--user  -u

// Do not restart automatically. It will automatically restart by default.
--no-autorestart    -n
```

#### Example：

```shell
sudo pmon3 run ./bin/gin --args "-prjHome=`pwd`" --user ntt360
```

:exclamation::exclamation: Note :exclamation::exclamation:

Parameter arguments need to use the absolute path。

#### View List  [ list/ls ]

```shell
sudo pmon3 ls
```

#### Start the process [ start ]

```shell
sudo pmon3 start [id or name]
```

#### Stop the process  [ stop ]

```shell
sudo pmon3 stop [id or name]
```

#### Reload the process [ reload ]

```shell
sudo pmon3 reload [id or name]
```

#### Process logging

```shell
# view logs of the process specified
sudo pmon3 log [id or name]

# Similar to using `tail -f xxx.log`
sudo pmon3 logf [id or name]
```

After editing the configuration file, the command needs to be used in conjunction with the startup process, the `reload` command defaults to only send the `SIGUSR2` signal to the startup process

If you want to customize the signal when you want to recoad, then use the `--sig` parameter:

```shell
// currently supported signals：HUP, USR1, USR2
sudo pmon3 reload --sig HUP [id or name]
```

#### Delete the process  [ del/delete ]

```shell
sudo pmon3 del [id or name]
```

#### View details [ show/desc ]

```shell
sudo pmon3 show [id or name]
```
![](https://jscssimg-img.oss-cn-beijing.aliyuncs.com/89c3f649a583a852.png?t=1506950494)

## Development

### Testing
```
make test
```

## Common Problems

### 1. Log Cutting?

`pmon3` comes with a logrotate configuration file, which by default utilizes the `/var/log/pmon3` directory. If you require a custom log path, please implement the log rotation by yourself.

### 2. The process startup parameter must pass the absolute path?

If there is a path in the parameters you pass, please use the absolute path. The `pmon3` startup process will start a new sandbox environment to avoid environmental variable pollution.

### 3. Command line automation

`pmon3` provides Bash automation. If you find that the command cannot be automatically provided in the sudo mode, please install `bash-completion` and exit the terminal to re-enter:

```bash
sudo yum install -y bash-completion
```

### 4. FATAL stat /var/run/pmon3/pmon3.sock: no such file or directory

If you encounter the error above, make sure the pmond service has started sucessfully.

```bash
sudo systemctl start pmond
```
