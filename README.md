# APM

> apm = **A**nyone's **P**rocess **M**anager
## Overview

Inspired by `PM2` where popular in the Node.js world, `apm` is designed to be a multiple-interface, reliable process manager.

`apm` supports the following features:
- create, start, stop, restart a process
- auto-restart when process quitted accidentially
- list current process status

## Install

```
go get -u github.com/DemoHn/apm/cmd/apm
```

## Usage

#### start
```
# start a new instance
apm start --name [name] --cmd [command to start]

# start an existing instance
apm start --id [id]
```

#### stop
```
apm stop --id [id]
```

#### list
```
# list info of all instances
apm list

# list info of one instance
apm list --id [id]
```
#### kill 
```
# stop all instances and kill apm daemon
apm kill
```

## TODO

- [ ] Make testcase works and covers at least 80% of code
- [ ] Add command: delete Instances
- [ ] Add logic: stop all instances when `apm kill` 
- [ ] Use config file to preload instances