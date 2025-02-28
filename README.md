# pmon3
[![CI](https://github.com/joe-at-startupmedia/pmon3/actions/workflows/ci.yml/badge.svg)](https://github.com/joe-at-startupmedia/pmon3/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/joe-at-startupmedia/pmon3/graph/badge.svg?token=RAILGSJ6XQ)](https://codecov.io/gh/joe-at-startupmedia/pmon3)
[![Release](https://img.shields.io/github/v/release/joe-at-startupmedia/pmon3)](https://github.com/joe-at-startupmedia/pmon3/releases/latest)


`pmon3` is a process manager (currently supports Linux and MacOSX) which allows you to keep processes alive forever. Processes can be declared directly via the [CLI](#section_commands) or in a [configuration file](#section_processconfig). `pmon3` allows you to get started quickly while providing configuration granularity at both the [system](#section_config) and [process](#section_processconfig) levels.

<img width="537" alt="pmon3_ls" src="https://github.com/joe-at-startupmedia/pmon3/assets/13522698/5d79ad53-664d-4ee7-bfac-f3fc94c2b316">

* [Introduction](#section_intro)
* [Installation](#section_install)
* [CLI Commands](#section_commands)
* [System Configuration](#section_config)
* [Process Configuration](#section_processconfig)
* [Process Dependencies](#section_dependencies)
* [Groups](#section_groups)
* [Event Handling](#section_events)
* [Flap Detection/Prevention](#section_flapping)
* [Debugging](#section_debugging)
* [Performance](#section_performance)
* [Problems](#section_problems)

<a name="section_intro"></a>
## Introduction

Golang currently has no officially supported process management tools. For managing processes, some use built-in commands such as `nohup [process] &`, or the process management tools provided by the operating system such as systemd. Alternatively, third-party process management tools such as Supervisor (python) or PM2 (node.js) can also be utilized.

Unlike PM2, `pmon3` is managed directly by the OS process manager, so if the `pmon3` CLI abnormally terminates, it will not affect the `pmond` daemon process. This is currently achieved by separating the `pmond` daemon process from the `pmon3` CLI.

By default, if a process abnormally terminates, `pmond` will try to restart the process. If you don't want a process to restart automatically, you can provide the `--no-autorestart` parameter flag from the CLI.

<a name="section_install"></a>
## Installation

[Releases](https://github.com/joe-at-startupmedia/pmon3/releases) 

### Using Go

```bash
  git clone https://github.com/joe-at-startupmedia/pmon3/ && cd pmon3 
  go mod tidy
  go build -o bin/pmon3 cmd/pmon3/pmon3.go
  go build -o bin/pmond cmd/pmond/pmond.go
  cp -R bin/pmon* /usr/local/bin/
  #create the configuration directory
  mkdir -p /etc/pmon3/config/
  cp config.yml /etc/pmon3/config/
  #start the daemon
  sudo /usr/local/bin/pmond &
```

### Makefile
The systemd installation process entails the following steps:
1. create the configuration directory
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
### Release Installer

```bash
wget -O - https://raw.githubusercontent.com/joe-at-startupmedia/pmon3/master/release-installer.bash | bash
```

:exclamation: After installing `pmon3` for the first time, both installation methods provided above should automatically enable and start the service. if the `pmond` service does not start automatically, you need to manually start the service.

```bash
sudo systemctl start pmond

# Others
sudo /usr/local/bin/pmond &
```

<a name="section_commands"></a>
## CLI Commands

### Help

```
Usage:
  pmon3 [command]

Available Commands:
  completion  Generate completion script
  del         Delete process by id or name
  desc        Show process information by id or name
  dgraph      Show the process queue order
  drop        Delete all processes
  exec        Spawn a new process
  export      Export Process Configuration
  group       Group level commands
  help        Help about any command
  init        Initialize all stopped processes
  kill        Terminate all processes
  log         Display process logs by id or name
  logf        Tail process logs by id or name
  ls          List all processes
  reset       Reset the restart counter(s)
  restart     (Re)start a process by id or name
  stop        Stop a process by id or name
  topn        Shows processes with unix top cmd
  version

Flags:
  -h, --help   help for pmon3

Use "pmon3 [command] --help" for more information about a command.
```

<a name="pmon3_exec"></a>
### Creating a process [run/exec]

```bash
pmon3 exec [process_binary] [flags]
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
--dependencies parent-process-name [--dependencies parent-process-name2]...

//provide a list of group names this process is associated to
--groups group-name-one [--groups group-name-2]
```

#### Example：

```bash
pmon3 exec ./bin/gin --args "-prjHome=`pwd`" --user joe
```

:exclamation: Parameter arguments need to use the absolute path.

### View List  [list/ls]

```bash
pmon3 ls
```

### (re)start the process [restart/start]

```bash
pmon3 restart [id or name]
```

<a name="pmon3_stop"></a>
### Stop the process  [stop]

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

### Delete the process  [del/delete]

```bash
pmon3 del [id or name]
```

### View details [show/desc]

```bash
pmon3 show [id or name]
```

<img width="475" alt="pmon3_show" src="https://github.com/joe-at-startupmedia/pmon3/assets/13522698/6b564a1c-0e26-468c-bd01-6dabce0c7620">

<a name="pmon3_kill"></a>
### Terminate all running process [kill]
```bash
pmon3 kill [--force]
```

<a name="pmon3_init"></a>
### (re)start all stopped process [init]
```bash
#(re)start processes specified in the Process Config only
pmon3 init --process-config-only

#(re)start processes specified in the Process Config and those which already exist in the database
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
#processes specified in the Process Config only
pmon3 dgraph --process-config-only

#processes specified in the Process Config and the database
pmon3 dgraph
```

<a name="pmon3_reset"></a>
### Reset the restart counter(s) [reset]

This command will reset the restart counter back to 0. This is useful to establish a clean slate without having to restart the pmond process.

```bash
#reset the restart counter on all processes
pmon3 reset

#reset the restart counter on a specific process
pmon3 reset [process_id_or_name]
```

<a name="pmon3_export"></a>
### Export Process Configuration [export]

This command is useful when you want to generate [Process Configuration](#section_appconfig) to use for pmond initialization from the specified `process_config_file`.

```bash
pmon3 export

#specify json as a format (default)
pmon3 export json

#specify toml as a format
pmon3 export toml

#specify yaml as a format
pmon3 export yaml

#order by name ascending instead of id
pmon3 export yaml -o name
```

### Top Native [topn]

This will output the resource utilization of all processes using the native `top` command that is pre-installed on most unix-based operating systems. It will only show those processes managed by (and including) the `pmond` process.

```bash
pmon3 topn
```
<img width="559" alt="pmon3_topn" src="https://github.com/joe-at-startupmedia/pmon3/assets/13522698/a77cce0f-55b0-479f-8489-d6aaf9fcdd6b">

<a name="section_config"></a>
## System Configuration
The default path of the configuration file is `/etc/pmon3/config/config.yml`. This value can be overridden with the `PMON3_CONF` environment variable. 
The following configuration options are available:
```yaml
# log levels: debug/info/warn/error
log_level: info

# kill processes on termination
handle_interrupts: true

# whether to reload from the configuration before executing a command
disable_reloads: false

# poll processes every [n] milliseconds
process_monitor_interval: 500

# wait [n] seconds before monitoring process statuses
initialization_period: 30

# a configuration file to specify a list of processes to start on the first initialization (json, yaml or toml)
process_config_file: /etc/pmon3/config/process.config.json

wait:
  # wait [n] milliseconds before outputting list after running init/stop/restart/kill/drop/exec
  cmd_exec_response: 1500
  # wait [n] milliseconds after connecting to IPC client before issuing commands
  ipc_connection: 0
  # wait [n] milliseconds after enqueueing a dependent process
  dependent_process_enqueued: 1000

#permission settings applied to files and directories of data/logs/message_queue
permissions:
  # an OS user to access files in data/logs/message_queue directory
  user:
  # an OS group to access files in data/logs/message_queue directory (must also provide a user)
  group:
  # the mode to apply to the data/logs/message_queue directory
  directory_mode: 0770
  # the mode to apply to files within the data/logs/message_queue directory
  file_mode: 0660

#settings specific to the logs directory (inherits from permissions by default)
logs:
  # directory where the logs are stored
  directory: /var/log/pmond
  # an OS user to access files in logs directory
  user:
  # an OS group to access files in logs directory (must also provide a user)
  group:
  # the mode to apply to the logs directory
  directory_mode: 0775
  # the mode to apply to files within the logs directory
  file_mode:

#settings specific to the data directory (inherits from permissions by default)
data:
  # directory where the database is stored
  directory: /etc/pmon3/data
  # an OS user to access files in data directory
  user:
  # an OS group to access files in data directory (must also provide a user)
  group:
  # the mode to apply to the data directory
  directory_mode:
  # the mode to apply to files within the data directory
  file_mode:

#settings specific to the message_queue directory (inherits from permissions by default)
message_queue:
  # directory where the messages are stored
  directory:
    # custom shared memory directory
    shmem: /dev/shm/
    # custom posix_mq directory
    posix_mq: /dev/mqueue/
  # a string to append to the name of the queue
  name_suffix:
  # an OS user to access files in message_queue directory
  user:
  # an OS group to access files in message_queue directory (must also provide a user)
  group:
  # the mode to apply to the message_queue directory
  directory_mode: 0775
  # the mode to apply to files within the message_queue directory
  file_mode: 0666

event_handling:
  # a script to execute when a process is restarted which accepts the process details json as the first argument
  process_restart:
  # a script to execute when a process fails (--no-autorestart) which accepts the process details json as the first argument
  process_failure:
  # a script to execute when a process backs off when flap detection is enabled
  process_backoff:

flap_detection:
  # enable flap detection
  is_enabled: false
  # the amount of times a process can restart (within the countdown threshold) until backoff evaluation begins
  threshold_restarted: 5
  # the amount of process monitor intervals during a processes backoff period until process evaluation proceeds as normal
  threshold_countdown: 120
  # the amount of process monitor intervals during a processes backoff period until the processes restart counter is decremented. disabled with 0 value.
  threshold_decrement: 60
```

Restarting pmond is usually unnecessary: All configuration changes should take effect when the next command is issued.

<a name="section_config_envvars"></a>
### Environment Variables

The configuration values can be overridden using environment variables:

* `CONFIGOR_LOGLEVEL`
* `CONFIGOR_HANDLEINTERRUPTS`
* `CONFIGOR_PROCESSMONITORINTERVAL`
* `CONFIGOR_INITIALIZATIONPERIOD`
* `CONFIGOR_PROCESSCONFIGFILE`
* `CONFIGOR_WAIT_CMDEXECRESPONSE`
* `CONFIGOR_WAIT_IPCCONNECTION`
* `CONFIGOR_WAIT_DEPENDENTPROCESSENQUEUED`
* `CONFIGOR_PERMISSIONS_USER`
* `CONFIGOR_PERMISSIONS_GROUP`
* `CONFIGOR_PERMISSIONS_DIRECTORYMODE`
* `CONFIGOR_PERMISSIONS_FILEMODE`
* `CONFIGOR_LOGS_DIRECTORY`
* `CONFIGOR_LOGS_USER`
* `CONFIGOR_LOGS_GROUP`
* `CONFIGOR_LOGS_DIRECTORYMODE`
* `CONFIGOR_LOGS_FILEMODE`
* `CONFIGOR_DATA_DIRECTORY`
* `CONFIGOR_DATA_USER`
* `CONFIGOR_DATA_GROUP`
* `CONFIGOR_DATA_DIRECTORYMODE`
* `CONFIGOR_DATA_FILEMODE`
* `CONFIGOR_MESSAGEQUEUE_DIRECTORY`
* `CONFIGOR_MESSAGEQUEUE_DIRECTORY_SHMEM`
* `CONFIGOR_MESSAGEQUEUE_DIRECTORY_POSIXMQ`
* `CONFIGOR_MESSAGEQUEUE_NAMESUFFIX`
* `CONFIGOR_MESSAGEQUEUE_USER`
* `CONFIGOR_MESSAGEQUEUE_GROUP`
* `CONFIGOR_MESSAGEQUEUE_DIRECTORYMODE`
* `CONFIGOR_MESSAGEQUEUE_FILEMODE`
* `CONFIGOR_EVENTHANDLER_PROCESSRESTART`
* `CONFIGOR_EVENTHANDLER_PROCESSFAILURE`
* `CONFIGOR_FLAPDETECTION_ISENABLED`
* `CONFIGOR_FLAPDETECTION_THRESHOLDRESTARTED`
* `CONFIGOR_FLAPDETECTION_THRESHOLDCOUNTDOWN`
* `CONFIGOR_FLAPDETECTION_THRESHOLDDECREMENT`

<a name="section_processconfig"></a>
## Process Configuration

By default, when `pmond` is restarted from a previously stopped state, it will load all processes in the database that were: 
* previously running
* have been marked as stopped as a result of pmond closing 
* have `--no-autorestart` set to false (default value)

```yaml
# a configuration file to specify a list of processes to start on the first initialization (json, yaml or toml)
process_config_file: /etc/pmon3/config/process.config.json
```

### Supported Formats

#### /etc/pmon3/config/process.config.json
```json
{
  "processes": [
    {
      "file": "/usr/local/bin/happac",
      "name": "happac1",
      "args": "-h startup-patroni-1.node.consul -p 5555 -r 5000",
      "user": "vagrant",
      "log_dir": "/var/log/custom/",
      "dependencies": ["happac2"],
      "groups": ["happac"]
    },
    {
      "file": "/usr/local/bin/happab",
      "name": "happac2",
      "log": "/var/log/happac2.log",
      "args": "-h startup-patroni-1.node.consul -p 5556 -r 5001",
      "user": "vagrant",
      "no_auto_restart": true,
      "groups": ["happac"]
    },
    {
      "file": "/usr/local/bin/node",
      "name": "metabase-api",
      "args": "/var/www/vhosts/metabase-api/index.js",
      "env_vars": "NODE_ENV=prod",
      "user": "dw_user"
    }
  ]
}
```

#### /etc/pmon3/config/process.config.yaml
```yaml
processes:
  - file: "/usr/local/bin/happac"
    name: happac1
    args: "-h startup-patroni-1.node.consul -p 5555 -r 5000"
    user: vagrant
    log_dir: "/var/log/custom/"
    dependencies:
    - happac2
    groups:
    - happac
  - file: "/usr/local/bin/happab"
    name: happac2
    log: "/var/log/happac2.log"
    args: "-h startup-patroni-1.node.consul -p 5556 -r 5001"
    user: vagrant
    no_auto_restart: true
    groups:
    - happac
  - file: "/usr/local/bin/node"
    name: metabase-api
    args: "/var/www/vhosts/metabase-api/index.js"
    env_vars: NODE_ENV=prod
    user: dw_user
```

#### /etc/pmon3/config/process.config.toml
Unlike json and yaml, all fields are camel-cased:
```toml
[[Processes]]
  File = "/usr/local/bin/happac"
  Name = "happac1"
  Args = "-h startup-patroni-1.node.consul -p 5555 -r 5000"
  User = "vagrant"
  LogDir = "/var/log/custom/"
  Dependencies = [ "happac2" ]
  Groups = [ "happac" ]

[[Processes]]
  File = "/usr/local/bin/happab"
  Name = "happac2"
  Log = "/var/log/happac2.log"
  Args = "-h startup-patroni-1.node.consul -p 5556 -r 5001"
  User = "vagrant"
  NoAutoRestart = true
  Groups = [ "happac" ]

[[Processes]]
  File = "/usr/local/bin/node"
  Name = "metabase-api"
  Args = "/var/www/vhosts/metabase-api/index.js"
  EnvVars = "NODE_ENV=prod"
  User = "dw_user"
```

### Generation Utility

Instead of configuring this file from scratch you can use the [export](#pmon3_export) command to output the configuration from the current process list. This allows the administrator to build a process list using imperative commands followed by exporting the results to a configuration file.

### Flags

All possible `flags` values matching those specified in the [exec](#exec_flags) command:

* file
* user
* log
* log_dir
* no_auto_restart
* args
* env_vars
* name
* dependencies
* groups

<a name="section_dependencies"></a>
## Dependencies

Dependencies (provided as an array) determine the order in which the processes are started. They are sorted using a directed acyclic graph meaning that there cannot be cyclical dependencies between processes. Dependency resolution can be debugged using the [dgraph](#pmon3_dgraph) command. Parent processes can wait `[n]` amount of seconds between spawning dependent processes by utilizing the `dependent_process_enqueued` configuration variable (currently defaults to `2` seconds).

<a name="section_groups"></a>
## Groups

Groups are useful when dealing with a large amount of related processes. Like processes, they are stored in the database and provide many-to-many cardinality. This allows the ability to associate multiple groups to one or more processes and vice versa. Groups can be managed via [Process Configuration](#section_processconfig) and through the [CLI](#exec_flags).

### Commands
```
Group level commands

Usage:
  pmon3 group [command]

Aliases:
  group, groups

Available Commands:
  assign      Assign group(s) to process(es)
  create      Create a new group
  del         Delete a group
  desc        Show group details and associated processes
  drop        Delete all processes associated to a group
  ls          List all groups
  remove      Remove process(es) from group(s)
  restart     (Re)start processes by group id or name
  stop        Stop all processes associated to a group

Flags:
  -h, --help   help for group

Use "pmon3 group [command] --help" for more information about a command.
```

### Examples

First lets create a group
```
pmon3 group create happac
```

Next, lets assign it to processes with ids 3 and 4
```
pmon3 group assign happac 3,4
```

Now, lets confirm that it was associated to the correct processes:
```
pmon3 group desc happac
```
![Screenshot Groups Desc](https://github.com/user-attachments/assets/446dbb6d-0ae5-47d2-bf99-6b5d9c314428)

Oops, process 4 should not be associated with the happac group, lets remove the association:
```
pmon3 group remove happac 4
```

This group is no longer useful, lets delete it while keeping the processes intact:
```
pmon3 group del happac
```

### Reloading Process Configuration Changes

If you make a change to the group in the [Process Configuration](#section_processconfig) while pmond is running, you can make the changes take effect by running the `init` command. The `init` command should not restart processes which are already running but it will apply changes from the `process_config_file`.
```
pmon3 init
```

<a name="section_events"></a>
## Event Handling With Custom Scripts

```yaml
event_handling:
  # a script to execute when a process is restarted which accepts the process details json as the first argument
  process_restart:
  # a script to execute when a process fails (--no-autorestart) which accepts the process details json as the first argument
  process_failure:
  # a script to execute when a process backs off when flap detection is enabled
  process_backoff:
```

### 1. Specify the executable script to run for the `process_restart` value. `pmond` will pass a json-escaped string of the process details as the first argument.
#### /etc/pmond/config/config.yml
```yaml
event_handling:
  # a script to execute when a process is restarted which accepts the process details json as the first argument
  process_restart: "/etc/pmon3/bin/on_restart.bash"
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

<a name="section_flapping"></a>
## Flap Detection/Prevention

![Screenshot Backkoff](https://github.com/user-attachments/assets/ec0d482f-9afe-438d-a5fd-18293a5c9ba2)


Flap Detection provides the ability to detect processes which are in a perpetually-failed state resulting in excessive restarts which can affect system performance among other things. To prevent excessive restarts you can enable flap detection in the configuration file. You can control the behavior by setting various [threshold parameters](#section_config) which ultimately affect how often a perpetually-failed process can restart within a given time interval.


### Enabling
```yaml
flap_detection:
  # enable flap detection
  is_enabled: true
```

### Restart Threshold
Defaulted to `5`, is the amount of process restarts before the flap prevention process begins at which point the process will cease restarts and enter the backoff state.
```yaml
flap_detection:
  # the amount of times a process can restart (within the countdown threshold) until backoff evaluation begins
  threshold_restarted: 5
```

### Countdown Threshold
Defaulted to `120`, is the amount of process monitor intervals until the flap prevention process (backoff state) ends and the process (if still in a perpetually failed state) will resume restarting as normal until the restart threshold is met again. The process monitor interval can also be set in the [configuration file](#section_config) which would affect the time in which it would take to countdown back to zero.
```yaml
flap_detection:
  # the amount of process monitor intervals during a processes backoff period until process evaluation proceeds as normal
  threshold_countdown: 120
```

### Decrement Threshold
Defaulted to `60`, disabled with `0`, is the amount of process monitor intervals during the flap prevention process (backoff state) until the internal process restart counter is decremented. This can affect how the countdown is reached effectively staggering process restarts during the countdown process. This is useful when you don't want to completely back off and allow for intermittent restarts during the flap prevention process (backoff state).
```yaml
flap_detection:
  # the amount of process monitor intervals during a processes backoff period until the processes restart counter is decremented. disabled with 0 value.
  threshold_decrement: 60
```

### Example 
Using the defaults provided above: since the process monitor interval defaults to `500` milliseconds, when a process enters the perpetually failed state and restarts the 5th time, it will enter the backoff state. At this point the flap prevention process will begin the countdown from `120` to `0`. Since the process monitor interval is `500` milliseconds, it will take `120` multiplied by `500` milliseconds (`60` seconds).

caveat: because the decrement threshold is `60`, the internal restart counter will decrement from `5` to `4` after `60` multiplied by `500` milliseconds (`30` seconds). It will restart `30` seconds into the backoff state instead of `60` seconds and continue counting down to `0` until exiting the flap detection process. If the process remains in a perpetually failed state, it will take `5` restarts to repeat this process all over again.

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

You can set the log level in the yaml configuration file.

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

The `posix_mq` build tag can be provided to swap out the underlying [gipc](https://github.com/joe-at-startupmedia/gipc/) library with [posix_mq](https://github.com/joe-at-startupmedia/posix_mq). The `posix_mq` wrapper does require `CGO_ENABLED=1` and is considerably faster but also consumes slightly more CPU and Memory. To enable `posix_mq` during the build process:
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

You should only use `sudo` to start the `pmond` process which requires superuser privileges due to the required process forking commands. However, the `pmon3` cli should be used *without* `sudo` to ensure that the spawned processes are attached to the correct parent pid. When using `sudo`, the processes will be attached to ppid 1 and as a result, will become orphaned if the `pmond` process exits prematurely. Using `sudo` also prevents non-root users from being able to access the log files. 

#### Spawn a new process as the root user
You must have sudo privileges to do this for security reasons. The `--user root` flag is redundant because the process is spawned as the calling user by default and commands ran as sudo are called by the root user.
```bash
sudo pmon3 exec /usr/local/bin/happac --user root
```
