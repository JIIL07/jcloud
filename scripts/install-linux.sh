#!/bin/bash

TARGET_DIR="/opt"
BIN_DIR="/opt/jcloud/bin"
ARCHIVE_PATH="./jcloud.tar.gz"
PROFILE_FILE="/etc/profile"

sudo mv $ARCHIVE_PATH /opt/

cd /opt
sudo tar -xzvf jcloud.tar.gz -C $TARGET_DIR

sudo rm jcloud.tar.gz

sudo sh -c "echo 'export CLIENT_CONFIG_PATH=\"$CONFIG_DIR/client.yaml\"' >> $PROFILE_FILE"
sudo sh -c "echo 'export PATH=\$PATH:$BIN_DIR' >> $PROFILE_FILE"

echo "Install successfully completed, jcloud installed in $TARGET_DIR/jcloud"
source etc/profile
echo "Environment variables updated in $PROFILE_FILE"