#!/bin/bash

# Run shared script
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source $SCRIPT_DIR/shared.sh

# Define variables
ROOT_FOLDER_SERVER_LOCATION="/schej.it"

# SCP some miscellaneous gitignored files
# scp -i $AWS_KEY_LOCATION -r $SCRIPT_DIR/../server/schools/waldorf/allowed_emails.json $SERVER_HOST:~/
# ssh $SERVER_HOST -i $AWS_KEY_LOCATION "sudo mv ~/allowed_emails.json $ROOT_FOLDER_SERVER_LOCATION/server/schools/waldorf/"

# Git pull on server, npm install, and restart server process
echo "Deploying server..."
ssh $SERVER_HOST -i $AWS_KEY_LOCATION "cd $ROOT_FOLDER_SERVER_LOCATION && sudo git stash && sudo git pull && cd server && sudo go build -buildvcs=false && screen -XS schej.it-server quit; screen -dmS schej.it-server sudo ./server -release=true"
echo "Done!"
