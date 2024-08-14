# pmon3
[![Testing](https://github.com/joe-at-startupmedia/pmon3/actions/workflows/testing.yml/badge.svg)](https://github.com/joe-at-startupmedia/pmon3/actions/workflows/testing.yml)


`pmon3` is a process manager for Golang applications. It allows you to keep applications alive forever and to reload them without downtime.

<img width="537" alt="pmon3_ls" src="https://github.com/joe-at-startupmedia/pmon3/assets/13522698/5d79ad53-664d-4ee7-bfac-f3fc94c2b316">

* [Introduction](#section_intro)
* [Installation](#section_install)
* [Commands](#section_commands)
* [Configuration](#section_config)
* [Application(s) Config](#section_appconfig)
* [Event Handling](#section_events)
* [Debugging](#section_debugging)
* [Performance](#section_performance)
* [Problems](#section_problems)

<a name="section_intro"></a>
## Introduction

Golang currently has no officially supported process management tools. For the deployment of Golang services, some use Linux built-in commands such as `nohup [process] &`, or the process management tools provided by the operating system such as SystemD. Alternatively, third-party process management tools such as: Python's Supervisor or Node.js PM2 can also be utilized

Each method has certain advantages and disadvantages. We hope to provide a convenient and easy-to-use tool for Golang process deployment. There is no need to install other dependencies besides `bash-completion` for ease of command line utilization.

Unlike PM2, `pmon3` is managed directly by the OS process manager, so even if the `pmon3` management tool abnormally terminates, it will not affect the parent `pmond` process itself. This is currently achieved by separating the `pmond` daemon process from the `pmon3` agent.

By default, if a process abnormally terminates, `pmond` will try to restart the process. If you don't want a process to restart automatically, then you can provide the `--no-autorestart` parameter flag.

<a name="section_install"></a>
## How To Install

[Releases](https://github.com/joe-at-startupmedia/pmon3/releases) 

### Using Makefile
The systemd installation process entails the following steps:
1. create the log, configuration and database directories
1. create the log rotation file
1. create the bash completion profile (requires the bash-completion package)
1. enable and start the `pmond` system service

```bash
#build the project
make build
#install on systemd-based system
make systemd_install
```

<a name="release_installer"></a>
### Using Release Installer

```bash
wget -O - https://raw.githubusercontent.com/joe-at-startupmedia/pmon3/master/release-installer.bash | bash -s 1.15.1
```

:exclamation::exclamation: Note :exclamation::exclamation:

After installing `pmon3` for the first time, both installation methods provided above should automatically enable and start the service. if the `pmond` service does not start automatically, you need to manually start the service.

```bash
sudo systemctl start pmond

# Others
sudo /usr/local/pmon3/bin/pmond &
```

<a name="section_commands"></a>
## Commands

### Help

```
Usage:
  pmon3 [command]

Available Commands:
  completion  Generate completion script
  del         Delete process by id or name
  desc        Show process extended details
  dgraph      Show the process queue order
  drop        Delete all processes
  exec        Spawn a new process
  help        Help about any command
  init        initialize all stopped processes
  kill        Terminate all processes
  log         Display process logs by id or name
  logf        Tail process logs by id or name
  ls          List all processes
  restart     (re)start a process by id or name
  stop        Stop a process by id or name
  topn        Shows processes using the native top
  version

Flags:
  -h, --help   help for pmon3

Use "pmon3 [command] --help" for more information about a command.
```

<a name="pmon3_exec"></a>
### Running process [run/exec]

```bash
pmon3 exec [application_binary] [flags]
```

<a name="exec_flags"></a>
The starting process accepts several parameters. The parameter descriptions are as follows:

```
// The process name. It will use the file name of the binary as the default name if not provided 
--name

// Where to store logs. It will override the confuration files `logs_dir` property
--log-dir

// The absolute path of a custom log file
--log  -l

// Provide parameters to be passed to the binary, multiple parameters are separated by spaces
--args  -a "-arg1=val1 -arg2=val2"

// Provide environment variables (appended to those already existing on the system `os.Environ()`)
--env-vars "ENV_VAR_1=env_var_1_value ENV_VAR_2=env_var_2_value"

// managing user
--user -u

// Do not restart automatically. It will automatically restart by default.
--no-autorestart  -n

// Provide a list of process names that this process will depend on
--dependency parent-process-name [--dependency parent-process-name2]...
```

#### Example：

```bash
pmon3 exec ./bin/gin --args "-prjHome=`pwd`" --user ntt360
```

:exclamation::exclamation: Note :exclamation::exclamation:

Parameter arguments need to use the absolute path.

### View List  [ list/ls ]

```bash
pmon3 ls
```

### Top Native [ topn ]

This will output the resource utilization of all processes using the native `top` command that is pre-installed on most unix-based operating systems. It will only show those processes managed by (and including) the `pmond` process. The output is updated every few seconds until the process is terminated using Ctrl+C.

```bash
pmon3 topn
```
<img width="559" alt="pmon3_topn" src="https://github.com/joe-at-startupmedia/pmon3/assets/13522698/a77cce0f-55b0-479f-8489-d6aaf9fcdd6b">

### (re)start the process [ restart/start ]

```bash
pmon3 restart [id or name]
```

<a name="pmon3_stop"></a>
### Stop the process  [ stop ]

```bash
pmon3 stop [id or name]
```

### Process logging

```bash
# view logs of the process specified
pmon3 log [id or name]

# view logs of the process specified including those previously rotated/archived
pmon3 log -a [id or name]

# Similar to using `tail -f xxx.log`
pmon3 logf [id or name]
```

### Delete the process  [ del/delete ]

```bash
pmon3 del [id or name]
```

### View details [ show/desc ]

```bash
pmon3 show [id or name]
```

<img width="475" alt="pmon3_show" src="https://github.com/joe-at-startupmedia/pmon3/assets/13522698/6b564a1c-0e26-468c-bd01-6dabce0c7620">

<a name="pmon3_kill"></a>
### Terminate all running process [ kill ]
```bash
pmon3 kill [--force]
```

<a name="pmon3_init"></a>
### (re)start all stopped process [ init ]
```bash
#(re)start processes specified in the Apps Config only
pmon3 init --apps-config-only

#(re)start processes specified in the Apps Config and those which already exist in the database
pmon3 init
```

<a name="pmon3_drop"></a>
### Terminate and delete all processes [drop]
```bash
pmon3 drop [--force]
```

<a name="pmon3_dgraph"></a>
### Display the dependency graph [dgraph/order]

This command is useful to debug dependency resolution without (re)starting processes

```bash
#processes specified in the Apps Config only
pmon3 dgraph --apps-config-only

#processes specified in the Apps Config and the database
pmon3 dgraph
```

<a name="section_config"></a>
## Configuration
The default path of the configuration file is `/etc/pmon3/config/config.yml`. This value can be overridden with the `PMON3_CONF` environment variable. 
The following configuration options are available:
```yaml
# All commented values are the default when empty or omitted
#
# directory where the database is stored
#data_dir: /etc/pmon3/data
# directory where the logs are stored
#logs_dir: /var/log/pmond
# log levels: debug/info/warn/error
#log_level: info
# kill processes on termination
#handle_interrupts: true
# wait [n] milliseconds before outputting list after running init/stop/restart/kill/drop/exec
#cmd_exec_response_wait: 1500
# wait [n] milliseconds after connecting to IPC client before issuing commands
#ipc_connection_wait: 0
# wait [n] milliseconds after enqueueing a dependent process
#dependent_process_enqueued_wait: 2000
# a script to execute when a process is restarted which accepts the process details json as the first argument
#on_process_restart_exec:
# a script to execute when a process fails (--no-autorestart) which accepts the process details json as the first argument
#on_process_failure_exec:
# specify a custom posix_mq directory in order to apply the appropriate permissions
#posix_mq_dir: /dev/mqueue/
# specify a custom shared memory directory in order to apply the appropriate permissions
#shmem_dir: /dev/shm/
# specify a custom user to access files in posix_mq_dir  or shmem_dir
#mq_user:
# specify a custom group to access files in posix_mq_dir or shmem_dir (must also provide a user)
#mq_group:
# a JSON configuration file to specify a list of apps to start on the first initialization
#apps_config_file: /etc/pmon3/config/app.config.json
```

All configuration changes are effective when the next command is issued - restarting pmond is unnecessary.

The configuration values can be overridden using environment variables:


* `CONFIGOR_DATADIR`
* `CONFIGOR_LOGSDIR`
* `CONFIGOR_LOGLEVEL`
* `CONFIGOR_HANDLEINTERRUPTS`
* `CONFIGOR_CMDEXECRESPONSEWAIT`
* `CONFIGOR_IPCCONNECTIONWAIT`
* `CONFIGOR_DEPENDENTPROCESSENQUEUEDWAIT`
* `CONFIGOR_ONPROCESSRESTARTEXEC`
* `CONFIGOR_ONPROCESSFAILUREEXEC`
* `CONFIGOR_POSIXMESSAGEQUEUEDIR`
* `CONFIGOR_SHMEMDIR`
* `CONFIGOR_MESSAGEQUEUEUSER`
* `CONFIGOR_MESSAGEQUEUEGROUP`
* `CONFIGOR_APPSCONFIGFILE`

<a name="section_appconfig"></a>
## Application(s) Config

By default, when `pmond` is restarted from a previously stopped state, it will load all processes in the database that were previously running, have been marked as stopped as a result of pmond closing and have `--no-autorestart` set to false (default value).
If applications are specified in the Apps Config, they will overwrite matching processes which already exist in the database.

### Configuration

#### /etc/pmon3/config/config.yml
```yaml
#default value when empty or omitted
apps_config_file: /etc/pmon3/config/apps.config.json
```

#### /etc/pmon3/config/app.config.json
```json
{
  "apps": [
    {
      "file": "/usr/local/bin/happac",
      "flags": {
        "name": "happac1",
        "args": "-h startup-patroni-1.node.consul -p 5555 -r 5000",
        "user": "vagrant",
        "log_dir": "/var/log/custom/",
        "dependencies": ["happac2"]
      }
    },
    {
      "file": "/usr/local/bin/happab",
      "flags": {
        "name": "happac2",
        "log": "/var/log/happac2.log",
        "args": "-h startup-patroni-1.node.consul -p 5556 -r 5001",
        "user": "vagrant",
        "no_auto_restart": true
      }
    },
    {
      "file": "/usr/local/bin/node",
      "flags": {
        "name": "metabase-api",
        "args": "/var/www/vhosts/metabase-api/index.js",
        "env_vars": "NODE_ENV=prod",
        "user": "dw_user"
      }
    }
  ]
}
```

### Dependencies

Depenencies can be provided as a json array and determine the order in which the processes are booted. They are sorted using a directed acyclic graph meaning that there cannot be cyclical dependencies between processes (for obvious reasons). Dependecy resolution can be debugged using the [dgraph](#pmon3_dgraph) command. Parent processes can wait [n] amount of seconds between spawning dependent processes by utilziing the `dependent_process_enqueued_wait` configuration variable which currently defaults to 2 seconds.


### Flags

All possible `flags` values matching those specified in the [exec](#exec_flags) command:

* user
* log
* log_dir
* no_auto_restart
* args
* env_vars
* name
* dependencies

<a name="section_events"></a>
## Event Handling With Custom Scripts

```yaml
# a script to execute when a process is restarted which accepts the process details json as the first argument
on_process_restart_exec: ""
# a script to execute when a process fails (--no-autorestart) which accepts the process details json as the first argument
on_process_failure_exec: ""
```

### 1. Specify the executable script to run for the `on_process_restart_exec` value. `pmond` will pass a json-escaped string of the process details as the first argument.
#### /etc/pmond/config/config.yml
```yaml
on_process_restart_exec: "/etc/pmon3/bin/on_restart.bash"
```

### 2. create the script to accept the json-escaped process details:
#### /etc/pmon3/bin/on_restart.bash
```bash
PROCESS_JSON="$1"
PROCESS_ID=$(echo "${PROCESS_JSON}" | jq '.id')
PROCESS_NAME=$(echo "${PROCESS_JSON}" | jq '.name')
echo "process restarted: ${PROCESS_ID} - ${PROCESS_NAME}" >> /var/log/pmond/output.log
```

### 3. start pmond in debug mode
```bash 
$ PMON3_DEBUG=true pmond
INFO/vagrant/go_src/pmon3/pmond/observer/observer.go:29 pmon3/pmond/observer.HandleEvent() Received event: &{restarted 0xc0001da630}
WARN/vagrant/go_src/pmon3/pmond/observer/observer.go:47 pmon3/pmond/observer.onRestartEvent() restarting process: happac3 (3)
DEBU/vagrant/go_src/pmon3/pmond/observer/observer.go:70 pmon3/pmond/observer.onEventExec() Attempting event executor(restarted): /etc/pmon3/bin/on_restart.bash "{\"id\":3,\"created_at\":\"2024-05-03T05:44:25.114957302Z\",\"updated_at\":\"2024-05-03T06:09:18.71222185Z\",\"pid\":4952,\"log\":\"/var/log/pmond/acf3f83.log\",\"name\":\"happac3\",\"process_file\":\"/usr/local/bin/happac\",\"args\":\"-h startup-patroni-1.node.consul -p 5557 -r 5002\",\"status\":2,\"auto_restart\":true,\"uid\":1000,\"username\":\"vagrant\",\"gid\":1000}"
```

### 4. confirm the script executed successfully
```bash
$ tail /var/log/pmond/output.log
process restarted: 4 - "happac4"
```

<a name="section_debugging"></a>
## Debugging

### Environment Variables

You can specify debug verbosity from both the `pmon3` client and the `pmond` daemon process using the `PMON3_DEBUG` environment variable.

```bash
PMON3_DEBUG=true pmond 
```

`PMON3_DEBUG` accepts the following values:
* `true`: sets the debug level to debug
* `debug`: has the same effect as true
* `info`: sets the debug level to info
* `warn`: sets the debug level to warn
* `error`: sets the debug level to error

You can also debug the underlying IPC library using `QOG_DEBUG=true`

```bash
XIPC_DEBUG=true PMON3_DEBUG=true pmon3 ls
```

### Configuration File

If you want more control over the verbosity you can set the loglevel in the yaml configuration file.

##### /etc/pmond/config/config.yml
```yaml
#possible values: debug/info/warn/error
#default value when empty or omitted
log_level: "info"
```

If you do not specify a value, `info` will be the default Logrus level.

<a name="section_performance"></a>
## Performance Prioritization

### CGO_ENABLED=0

By default, no underlying libraries require CGO. This allows for portability between machines using different versions of GLIBC and also provides easy installation using the [Release Installer](#release_installer) . Benchmarking results have confirmed less memory and CPU utilization compared to using the libraries which do require `CGO_ENABLED=1` provided below:

### Posix MQ

The `posix_mq` build tag can be provided to swap out the underlying [golang-ipc](https://github.com/joe-at-startupmedia/golang-ipc/) library with [posix_mq](https://github.com/joe-at-startupmedia/posix_mq). The `posix_mq` wrapper does require `CGO_ENABLED=1` and is considerably faster but also consumes slightly more CPU and Memory. To enable `posix_mq` during the build process:
```bash
BUILD_FLAGS="-tags posix_mq" make build_cgo
```

### CGO-based Sqlite

By default, `pmon3` utilizes an non-CGO version of sqlite which is unnoticably less performant in most circumstances. To enable the CGO version of sqlite:
```bash
BUILD_FLAGS="-tags cgo_sqlite" make build_cgo
```

It depends on your requirements whether you need one or both. To enable both of these CGO-dependent modules for maximizing overall performance:
```bash
BUILD_FLAGS="-tags posix_mq,cgo_sqlite" make build_cgo
```

Or without using the Makefile:
```bash
CGO_ENABLED=1 go build -tags "posix_mq,cgo_sqlite" -o bin/pmon3 cmd/pmon3/pmon3.go
CGO_ENABLED=1 go build -tags "posix_mq,cgo_sqlite" -o bin/pmond cmd/pmond/pmond.go
```

### Unix Sockets
Significantly less performant than the default shared memory implementation and posix_mq implementation. It also has the capability of utilizing TCP cockets with additional build flags (currently: `build -tags net,network`).

```bash
BUILD_FLAGS="-tags net" make build
```

Or without using the Makefile:
```bash
CGO_ENABLED=0 go build -tags net -o bin/pmon3 cmd/pmon3/pmon3.go
CGO_ENABLED=0 go build -tags net -o bin/pmond cmd/pmond/pmond.go
```

<a name="section_problems"></a>
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
```bash
autoload -U +X compinit && compinit
autoload -U +X bashcompinit && bashcompinit
sudo sh -c "pmon3 completion zsh > /etc/profile.d/pmon3.sh"
source /etc/profile.d/pmon3.sh
```

### 4. FATA/vagrant/go_src/pmon3/cmd/pmon3/pmon3.go:27 main.main() pmond must be running

If you encounter the error above, make sure the `pmond` service has started successfully.

```bash
sudo systemctl start pmond
```

### 5. Should I use `sudo` commands?

You should only use `sudo` to start the `pmond` process which requires superuser privileges due to the required process forking commands. However, the `pmon3` cli should be used *without* `sudo` to ensure that the spawned processes are attached to the correct parent pid. When using `sudo`, the processes will be attached to ppid 1 and as a result, will become orphaned if the `pmond` process exits prematurely. Using `sudo` also prevents non-root users from being able to access the log files. The following Makefile command applies the adequate non-root permissions to the log files.

#### Applying permissions
```bash
#This is automatically called by make systemd_install
make systemd_permissions
```

#### Spawn a new process as the root user without using `sudo`
```bash
pmon3 exec /usr/local/bin/happac --user root
```

