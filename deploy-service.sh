#!/bin/sh


GO_CMD=$(which go)

if [ -z "$GO_CMD" ]; then
    echo "Go is not installed or not in PATH"
    exit 1
fi

$GO_CMD build

sudo systemctl stop screen-power-controller.service

sudo cp screen-power-controller /usr/local/bin/screen-power-controller
sudo cp screen-power-controller.service /etc/systemd/system/screen-power-controller.service

sudo systemctl daemon-reload
sudo systemctl enable screen-power-controller.service
sudo systemctl start screen-power-controller.service