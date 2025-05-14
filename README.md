# archstatuspage

Generates an HTML status page for arch linux systems.

## Waybar example usage

Intended usage is with waybar. 

Example

```bash
    "custom/failed-systemd": {
        "exec": "~/path/to/failed-systemd",
        "interval": 605,
        "return-type": "json",
        "on-click": "~/repos/archstatuspage/bin/archstatuspage && xdg-open /tmp/system-status.html"
    }

```

`failed-systemd`:

```bash
#!/bin/env bash

# This is a waybar module that checks for failed systemd services.
# See [here](https://wiki.archlinux.org/title/System_maintenance#Failed_systemd_services) for more information.

set -euo pipefail

SERVICE_ICON="󰣇"

if ! command -v systemctl &> /dev/null; then
    echo "{\"text\": \"ERROR - systemctl not found\", \"class\": \"systemd-error\"}"
    exit 1
fi

# Get failed services
FAILED_SERVICES=$(systemctl --failed --plain --no-legend | wc --lines)
FAILED_LIST=$(systemctl --failed --plain --no-legend | awk '{print $1}' | tr '\n' ',' | sed 's/,$//')

if [[ $FAILED_SERVICES -eq 0 ]]; then
    echo "{\"text\": \"${SERVICE_ICON} ✓\", \"class\": \"systemd-ok\", \"tooltip\": \"No failed systemd services\"}"
else
    echo "{\"text\": \"${SERVICE_ICON} ${FAILED_SERVICES}\", \"class\": \"systemd-failed\", \"tooltip\": \"Failed services: ${FAILED_LIST}\"}"
fi

```
