#!/usr/bin/env bash

BINARY_NAME="updemia-client"

clientConfiguration() {
    DEFAULT_SCREENSHOT_PATH="/tmp/updemia"
    read -p "    Screenshots absolute path: [$DEFAULT_SCREENSHOT_PATH] " SCREENSHOTS_PATH
    SCREENSHOTS_PATH=${SCREENSHOTS_PATH:-$DEFAULT_SCREENSHOT_PATH}
}


clientInstallation() {
    defaults write com.apple.screencapture location $SCREENSHOTS_PATH
    killall SystemUIServer
    LATEST_BINARY_FILE_URL=$(curl -s https://api.github.com/repos/ScullWM/updemia-client/releases | grep browser_download_url | head -n 1 | cut -d '"' -f 4)
    curl -sL $LATEST_BINARY_FILE_URL -o "/tmp/$BINARY_NAME"
    chmod a+x "/tmp/$BINARY_NAME"
    mv "/tmp/$BINARY_NAME" "/usr/local/bin/$BINARY_NAME"
}

cat <<EOF

   _/    _/            _/_/_/              _/      _/  _/
  _/    _/  _/_/_/    _/    _/    _/_/    _/_/  _/_/        _/_/_/
 _/    _/  _/    _/  _/    _/  _/_/_/_/  _/  _/  _/  _/  _/    _/
_/    _/  _/    _/  _/    _/  _/        _/      _/  _/  _/    _/
 _/_/    _/_/_/    _/_/_/      _/_/_/  _/      _/  _/    _/_/_/
        _/
       _/

EOF

read -p "WARNING: this action will be replace your current Mac OSX screenshots path. Are sure about continuing? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "\n🚧  Configuration..."
    clientConfiguration
    echo -e "\n✨  Binary download and installation in progress..."
    clientInstallation
    echo -e "\n✅  Done!\n"
fi
