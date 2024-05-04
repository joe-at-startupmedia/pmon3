# `pmon3`
`pmon3` is a process manager for Golang applications. It allows you to keep applications alive forever and to reload them without downtime.

<img width="726" alt="pmon3 help" src="https://github.com/joe-at-startupmedia/pmon3/assets/13522698/6f266ccf-7a14-4afc-b626-cbd53e245bc1">

## Introduction

Golang currently has no officially supported process management tools. For the deployment of Golang services, some use Linux built-in commands such as `nohup [process] &`, or the process management tools provided by the operating system such as SystemD. Alternatively, third-party process management tools such as: Python's Supervisor or Nodejs PM2 can also be utilized

Each method has certain advantages and disadvantages. We hope to provide a convenient and easy-to-use tool for Golang process deployment. There is no need to install other dependencies besides `bash-completion` for ease of command line utilization.

Unlike PM2, `pmon3` is managed directly by the OS process manager, so even if the `pmon3` management tool abnormally terminates, it will not affect the parent `pmond` process itself. This is currently achieved by seperating the `pmond` deamon process from the `pmon3` agent.

By default, if a process abnormally terminates, `pmond` will try to restart the process. If you don't want a process to restart automatically, then you can provide the `--no-autorestart` parameter flag.

## How To Install

[Releases](https://github.com/joe-at-startupmedia/pmon3/releases) 

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

### Using Release Installer

```bash
wget -O - https://raw.githubusercontent.com/joe-at-startupmedia/pmon3/master/release-installer.bash | bash -s 1.13.2
```

:exclamation::exclamation: Note :exclamation::exclamation:

After installing `pmon3` for the first time, both installation methods provided above should automatically enable and start the service. if the `pmond` service does not start automatically, you need to manually start the service.

```shell
sudo systemctl start pmond

# Others
sudo /usr/local/pmon3/bin/pmond &
```

## Command Introduction

#### Help

```
Usage:
  pmon3 [command]

Available Commands:
  completion  Generate completion script
  del         Delete process by id or name
  desc        Show process extended details
  drop        Delete all processes
  exec        Spawn a new process
  help        Help about any command
  init        Restart all stopped processes
  kill        Terminate all processes
  log         Display process logs by id or name
  logf        Tail process logs by id or name
  ls          List all processes
  restart     Restart a process by id or name
  stop        Stop a process by id or name
  version

Flags:
  -h, --help   help for pmon3
```

#### Running process [run/exec]

```
pmon3 exec [application_binary] [flags]
```
The starting process accepts several parameters. The parameter descriptions are as follows:

```
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

```
pmon3 exec ./bin/gin --args "-prjHome=`pwd`" --user ntt360
```

:exclamation::exclamation: Note :exclamation::exclamation:

Parameter arguments need to use the absolute path。

#### View List  [ list/ls ]

```
pmon3 ls
```

#### (re)start the process [ restart/start ]

```
pmon3 restart [id or name]
```

#### Stop the process  [ stop ]

```
pmon3 stop [id or name]
```

#### Process logging

```
# view logs of the process specified
pmon3 log [id or name]

# Similar to using `tail -f xxx.log`
pmon3 logf [id or name]
```

#### Delete the process  [ del/delete ]

```
pmon3 del [id or name]
```

#### View details [ show/desc ]

```
pmon3 show [id or name]
```

<img width="516" alt="pmon3 desc" src="https://github.com/joe-at-startupmedia/pmon3/assets/13522698/881af535-9a83-472c-b2a5-b0caff5596f3">

#### Terminate all running process [ kill ]
```
pmon3 kill [--force]
```

#### Restart all steopped process [ init ]
```
pmon3 init
```

#### Terminate and delete all processes [drop]
```
pmon3 drop [--force]
```

## Event Handling With Custom Scripts

```yaml
# a script to execute when a process is restarted which accepts the process details json as the first argument
on_process_restart_exec: ""
# a script to execute when a process fails (--no-autorestart) which accepts the process details json as the first argument
on_process_failure_exec: ""
```

#### 1. Specify the executable script to run for the `on_process_restart_exec` value. pmond will pass a json-escaped string of the process details as the first argument.
#### /etc/pmond/config/config.yaml
```yaml
on_process_restart_exec: "/etc/pmon3/bin/on_restart.bash"
```

#### 2. create the script to accept the json-escaped process details:
#### /etc/pmon3/bin/on_restart.bash
```bash
PROCESS_JSON="$1"
PROCESS_ID=$(echo "${PROCESS_JSON}" | jq '.id')
PROCESS_NAME=$(echo "${PROCESS_JSON}" | jq '.name')
echo "process restarted: ${PROCESS_ID} - ${PROCESS_NAME}" >> /var/log/pmond/output.log
```

#### 3. start pmond in debug mode
```bash 
$ PMON3_DEBUG=true pmond
INFO/vagrant/go_src/pmon3/pmond/observer/observer.go:29 pmon3/pmond/observer.HandleEvent() Received event: &{restarted 0xc0001da630}
WARN/vagrant/go_src/pmon3/pmond/observer/observer.go:47 pmon3/pmond/observer.onRestartEvent() restarting process: happac3 (3)
DEBU/vagrant/go_src/pmon3/pmond/observer/observer.go:70 pmon3/pmond/observer.onEventExec() Attempting event executor(restarted): /etc/pmon3/bin/on_restart.bash "{\"id\":3,\"created_at\":\"2024-05-03T05:44:25.114957302Z\",\"updated_at\":\"2024-05-03T06:09:18.71222185Z\",\"pid\":4952,\"log\":\"/var/log/pmond/acf3f83.log\",\"name\":\"happac3\",\"process_file\":\"/usr/local/bin/happac\",\"args\":\"-h startup-patroni-1.node.consul -p 5557 -r 5002\",\"status\":2,\"auto_restart\":true,\"uid\":1000,\"username\":\"vagrant\",\"gid\":1000}"
```

#### 4. confirm the script executed successfully
```bash
$ tail /var/log/pmond/output.log
process restarted: 4 - "happac4"
```

## Debugging

### Environment Variables

You can specify debug verbosity from both the pmon3 client and the pmond daemon process using `PMON3_DEBUG=true` as an environment vairable. This will set the Logrus level to `debug`

```bash
PMON3_DEBUG=true pmond 
```

You can also debug the underlying IPC library using `QOG_DEBUG=true`

```bash
QOG_DEBUG=true PMON3_DEBUG=true pmon3 ls
```

### Configuration File

If you want more control over the verbosity you can set the loglevel in the yaml configuration file.

##### /etc/pmond/config/config.yaml
```
# log levels: debug/info/warn/error
log_level: "info"
```

If you do not specify a value, `info` will be the default Logrus level.

## Performance Prioritization

### CGO_ENABLED=0

By default, no underlying libraries require CGO. This allows for portability between different machines using different versions of GLIBC and also provide easy installation using the (release-installer)[https://github.com/joe-at-startupmedia/pmon3/blob/master/release-installer.bash] . Benchmarking results have confirmed less memory and CPU utilization compared to using the libraries which do require `CGO_ENABLED=1`.

### posix_mq

The `posix_mq` build tag can be provided to swap out the underlying [golang-ipc](https://github.com/joe-at-startupmedia/golang-ipc/) library with (posix_mq)[https://github.com/joe-at-startupmedia/posix_mq]. The `posix_mq` wrapper does require `CGO_ENABLED=1` and is considerably faster but it also consumes slightly more CPU and Memory. To enable `posix_mq` during the build process:
```bash
BUILD_TAGS="posix_mq" make build-cgo
```

### CGO-based Sqlite

By default, `pmon3` utilizes an non-CGO version of sqlite which is unnoticably less performant in most circumstances. To enable the CGO version of sqlite:
```bash
BUILD_TAGS="cgo_sqlite" make build-cgo
```

If you enable `posix_mq`, you might as well enable `cgo_sqlite` since `CGO_ENABLED=1` is required for either module. It depends on your requirements whether or not you need one or both. To enable both of these CGO-dependent modules for maximizing overall performance:

```bash
BUILD_TAGS="posix_mq,cgo_sqlite" make build-cgo
```

Or without using the Makefile:
```bash
CGO_ENABLED=1 go build -tags "posix_mq,cgo_sqlite" -o bin/pmon3 cmd/pmon3/pmon3.go
CGO_ENABLED=1 go build -tags "posix_mq,cgo_sqlite" -o bin/pmond cmd/pmond/pmond.go
```

## Common Problems

### 1. Log Rotation?

`pmon3` comes with a logrotate configuration file, which by default utilizes the `/var/log/pmond` directory. If you require a custom log path, you can customize `config.yml` and `rpm/pmond.logrotate`

### 2. The process startup parameter must pass the absolute path?

If there is a path in the parameters you pass, please use the absolute path. The `pmon3` startup process will start a new sandbox environment to avoid environmental variable pollution.

### 3. Command line automation

`pmon3` provides Bash automation. If you find that the command cannot be automatically provided, please install `bash-completion` and exit the terminal to re-enter:

```shell
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

```shell
sudo systemctl start pmond
```

### 5. Should I use `sudo` commands?

You should only use `sudo` to start the `pmond` process which requires superuser privileges due to the required process forking commands. However, the `pmon3` cli should be used *without* `sudo` to ensure that the spawned processes are attached to the correct parent pid. When using `sudo`, the processes will be attached to ppid 1 and as a result, will become orphaned if the `pmond` process exits prematurely. Using `sudo` also prevents non-root users from being able to access the log files. The following Makefile command applies the adequate non-root permissions to the log files.

#### Applying permissions
```shell
#This is automatically called by make systemd_install
make systemd_permissions
```

#### Spawn a new process as the root user without using `sudo`
```
pmon3 exec /usr/local/bin/happac --user root
```
