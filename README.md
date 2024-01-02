# `pmon3`
`pmon3` is a process manager for Golang applications with a built-in load balancer. It allows you to keep applications alive forever and to reload them without downtime.

<img src="http://p0.qhimg.com/t017d6cbb68aed4b693.png" style="max-width:680px" />

## Start Process

```go
pmon3 run [./application binary] [arg1]  ...
```

## Introduction

Golang currently has no officially supported process management tools. For the deployment of Golang services, some use Linux built-in commands such as `nohup [process] &`, or the process management tools provided by the operating system such as SystemD. Alternatively, third-party process management tools such as: Python's Supervisor or Nodejs PM2 can also be utilized

Each method has certain advantages and disadvantages. We hope to provide a convenient and easy-to-use tool for Golang process deployment. There is no need to install other dependencies besides `bash-completion` for ease of command line utilization.

Unlike PM2, `pmon3` is managed directly by the OS process manager, so even if the `pmon3` management tool abnormally terminates, it will not affect the parent `pmond` process itself. This is currently achieved by seperating the `pmond` deamon process from the `pmon3` agent.

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
#build the project
make build
#install on systemd-based system
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
pmon3 help

Usage:
  pmon3 [command]

Available Commands:
  completion  Generate completion script
  del         Delete process by id or name
  desc        Show process extended details
  drop        Delete all processes
  exec        Spawn a new process
  help        Help about any command
  kill        Terminate all processes
  log         Display process logs by id or name
  logf        Tail process logs by id or name
  ls          List all processes
  restart     Restart a process by id or name
  stop        Stop a process by id or name
  version

Flags:
  -h, --help   help for pmon3

Use "pmon3 [command] --help" for more information about a command.
```

#### Running process [run/exec]

```shell
pmon3 run [./application binary] [arg1] [arg2] ...
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
pmon3 run ./bin/gin --args "-prjHome=`pwd`" --user ntt360
```

:exclamation::exclamation: Note :exclamation::exclamation:

Parameter arguments need to use the absolute path。

#### View List  [ list/ls ]

```shell
pmon3 ls
```

#### (re)tart the process [ restart ]

```shell
pmon3 restart [id or name]
```

#### Stop the process  [ stop ]

```shell
pmon3 stop [id or name]
```

#### Process logging

```shell
# view logs of the process specified
pmon3 log [id or name]

# Similar to using `tail -f xxx.log`
pmon3 logf [id or name]
```

#### Delete the process  [ del/delete ]

```shell
pmon3 del [id or name]
```

#### View details [ show/desc ]

```shell
pmon3 show [id or name]
```

#### Terminate all running process [ kill ]
```shell
pmon3 kill [--force]
```

#### Terminate and delete all processes [drop]
```shell
pmon3 drop [--force]
```

![](https://jscssimg-img.oss-cn-beijing.aliyuncs.com/89c3f649a583a852.png?t=1506950494)

## Development

### Testing
```
make test
```

## Common Problems

### 1. Log Rotation?

`pmon3` comes with a logrotate configuration file, which by default utilizes the `/var/log/pmond` directory. If you require a custom log path, you can customize `config.yml` and `rpm/pmond.logrotate`

### 2. The process startup parameter must pass the absolute path?

If there is a path in the parameters you pass, please use the absolute path. The `pmon3` startup process will start a new sandbox environment to avoid environmental variable pollution.

### 3. Command line automation

`pmon3` provides Bash automation. If you find that the command cannot be automatically provided, please install `bash-completion` and exit the terminal to re-enter:

```bash
sudo yum install -y bash-completion
```

#### Using ZSH instead of Bash
```shell
autoload -U +X compinit && compinit
autoload -U +X bashcompinit && bashcompinit
sudo sh -c "pmon3 completion zsh > /etc/profile.d/pmon3.sh"
source /etc/profile.d/pmon3.sh
```

### 4. FATA stat /var/run/pmon3/pmon3.sock: no such file or directory

If you encounter the error above, make sure the pmond service has started sucessfully.

```bash
sudo systemctl start pmond
```

### 5. Should I use `sudo` commands?

You should only use `sudo` to start the `pmond` process which requires superuser privileges due to the required process forking commands. However, the `pmon3` cli should be used *without* `sudo` to ensure that the spawned processes are attached to the correct parent pid. When using `sudo`, the processes will be attached to ppid 1 and as a result, will become orphaned if the `pmond` process exits prematurely. Using `sudo` also prevents non-root users from being able to access the log files. The following Makefile command applies the adequate non-root permissions to the log files and the database file:

```shell
#This is automatically called by make systemd_install
make systemd_permissions
```
