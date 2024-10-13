#!/bin/bash

# Run screen -ls and capture the output
screen_output=$(screen -ls)

# Check if the output contains "schej.it-server"
if echo "$screen_output" | grep -q "schej.it-server"; then
    echo "Screen 'schej.it-server' is running."
else
    echo "Screen 'schej.it-server' is not running! Rebooting..."
    cd ../server
    screen -dmS schej.it-server sudo ./server -release=true
fi
