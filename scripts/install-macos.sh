#!/bin/bash

TARGET_DIR="/opt/jcloud"
CONFIG_DIR="/opt/jcloud/config"
ARCHIVE_PATH="bin/jcloud.tar.gz"
CONFIG_PATH="config/client.yaml"

sudo mkdir -p "$TARGET_DIR"
sudo mkdir -p "$CONFIG_DIR"

sudo mv "$ARCHIVE_PATH" "$TARGET_DIR/"
sudo mv "$CONFIG_PATH" "$CONFIG_DIR"

cd "$TARGET_DIR" || exit
sudo tar -xzf jcloud.tar.gz
sudo rm jcloud.tar.gz

if ! echo "$PATH" | grep -q "$TARGET_DIR"; then
    if [ -f "$HOME/.bash_profile" ]; then
        echo "export PATH=\$PATH:$TARGET_DIR" >> "$HOME/.bash_profile"
    elif [ -f "$HOME/.zshrc" ]; then
        echo "export PATH=\$PATH:$TARGET_DIR" >> "$HOME/.zshrc"
    else
        echo "No known shell configuration file found. PATH update failed."
    fi
fi

CLIENT_CONFIG_PATH="/opt/jcloud/config/client.yaml"
export CLIENT_CONFIG_PATH

echo "Install successfully completed, jcloud installed in $TARGET_DIR"
