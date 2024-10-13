#!/bin/bash

# IMPORTANT: need to change the permissions of the log.txt if it's recreated
# chmod 666 "$output_file"

# Define the output file
output_file="/schej.it/deploy_scripts/reboot_log.txt"

# Run screen -ls and capture the output
screen_output=$(screen -ls)

# Get current timestamp
timestamp=$(date "+%Y-%m-%d %H:%M:%S")

# Check if the output contains "schej.it-server"
cd /schej.it/deploy_scripts
if echo "$screen_output" | grep -q "schej.it-server"; then
    echo "[$timestamp] Screen 'schej.it-server' is running." >> "$output_file"
else
    echo "[$timestamp] Screen 'schej.it-server' is not running! Rebooting..." >> "$output_file"
    cd ../server
    screen -dmS schej.it-server sudo ./server -release=true
fi
