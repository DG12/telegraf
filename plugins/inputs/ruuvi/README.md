# Ruuvi Tag Input Plugin

Reads Ruuvi tag basic firmware output via raw bluetooth sockets. Requires raw socket access to HCI device,
meaning that bluetooth services must be disabled for that particular Bluetooth Interface.

### Configuration 

```toml
[[inputs.ruuvi]]
## Mandatory device to use
hci_device = "hci0"
``` 

### Measurements
- ruuvi:
  - tags:
    - device
  - fields:
    - rssi
	- humidity
	- temperature
    - pressure
    - voltage
    - acceleration_x
    - acceleration_y
    - acceleration_z

### Reasoning

Originally designed to measure conditions in beer fermenter and existing solution used Java.
You should not contaminate beer with Java.
