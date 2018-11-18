# apm

> apm = **A**nyone's **P**rocess **M**anager

# NOTICE: README.md is under construction now!

## Overview

`apm` is designed to be a multiple-interface, easy to use process and reliable manager insprired by `PM2`.

It supports: (TODO)

- Process status control - start, stop, etc.
- Real-time process update
- Start a process with specific user & group

## API

TODO
- CLI (of course!)

## Design Layer 

```

-- Master (control all instances and handle pipes & updates)

-- Instance (with concrete status)

-- actual process <--> user
```

apm start --config apm.config.yml
apm stop --id 1
apm stop --name api-server

apm launch background <-- send RPC command
