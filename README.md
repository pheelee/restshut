# Restshut
![GitHub](https://img.shields.io/github/license/pheelee/restshut)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/pheelee/restshut)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/pheelee/restshut/release)
![GitHub release (latest by date)](https://img.shields.io/github/downloads/pheelee/restshut/latest/total)



I built this small cross os/architecture utility to support the ```turn_off``` action of [homeassistant's wake_on_lan switch](https://www.home-assistant.io/integrations/wake_on_lan/#examples). It exposes a http endpoint secured by a secret and optionally a list of allowed hosts. This http endpoint can be integrated in homeassistant using the [RESTful command integration](https://www.home-assistant.io/integrations/rest_command/) which can be used as action in the ```turn_off``` parameter.

---
## Installation

If you run restshut the first time it creates a sample config in the same directory which you can adjust to your needs. A logfile is also created in the same directory as restshut.

### Windows

To run the utility on every start in the background use a scheduled task.

### Linux

Please read the documentation for your linux distro on how to create a service.

Here is a sample system.d unit file:

```bash
[Unit]
Description=restshut service
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=1
User=myUser
ExecStart=/opt/restshut/restshut-linux-amd64

[Install]
WantedBy=multi-user.target
```

---
## Homeassistant Setup

In homeassistant you have to create a RESTful command in the configuration.
For security reasons I store the api key in the secrets file. The URL parameter is a template to support multiple endpoints in 1 configuration.
```yaml
rest_command:
    restshut:
        url: 'http://{{ip}}:7000'
        method: POST
        headers:
        Authorization: !secret restshut_key
```

Then create the Wake On Lan switch(es)

```yaml
switch:
- platform: wake_on_lan
  name: myPC
  host: 192.168.1.1
  mac: "12-34-56-78-90-AB"
  turn_off:
    service: rest_command.restshut
    data:
      ip: 192.168.1.1

- platform: wake_on_lan
  name: myNAS
  host: 192.168.1.100
  mac: "12-34-56-78-90-CD"
  turn_off:
    service: rest_command.restshut
    data:
      ip: 192.168.1.100
```