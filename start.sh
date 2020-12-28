#!/bin/bash -
BASE_PATH=/somfy PORT=9070 DEVICES_CONFIG=/opt/somfy/somfy.yaml SIGNALDUINO_ADDRESS=/dev/ttyUSB0 HOMEKIT_CONFIG_PATH=/opt/somfy/homekit HOMEKIT_CONFIG_PORT=40123 HOMEKIT_CONFIG_PIN=12344321 ./somfy-rts-gateway &>> somfy.log