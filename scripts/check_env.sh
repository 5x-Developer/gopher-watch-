#!/bin/bash

CONFIG_FILE="configs/targets.json"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "ERROR: Configuration file $CONFIG_FILE not found."
    exit 1
fi

if [ ! -s "$CONFIG_FILE" ]; then
    echo "ERROR: Configuration file $CONFIG_FILE is empty."
    exit 1
fi

echo "Environment Checks passed successfully."
exit 0