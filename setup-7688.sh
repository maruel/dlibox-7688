#!/bin/bash
# Copyright 2020 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

set -eu

echo "- Setting up dlibox as system service"
sudo tee /etc/systemd/system/dlibox-7688.service > /dev/null <<EOF
# https://github.com/maruel/dlibox-7688
[Unit]
Description=Runs dlibox-7688 automatically upon boot
Wants=network-online.target
After=network-online.target
[Service]
User=root
Group=root
KillMode=mixed
Restart=always
TimeoutStopSec=10s
ExecStart=/home/pi/go/bin/dlibox-7688
Environment=GOTRACEBACK=all
AmbientCapabilities=CAP_NET_BIND_SERVICE
[Install]
WantedBy=default.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable dlibox-7688.service

echo "Install unclutter"
sudo apt install unclutter

echo "Install emojis"
mkdir -p ~/.fonts
wget -O ~/.fonts/NotoColorEmoji.ttf https://github.com/googlefonts/noto-emoji/raw/master/fonts/NotoColorEmoji.ttf
fc-cache -f -v

echo "Setup chromium-browser"
sed -i 's/"exited_cleanly": false/"exited_cleanly": true/' \
      ~/.config/chromium/Default/Preferences

mkdir -p ~/.config/lxsession/LXDE-pi
cat > ~/.config/lxsession/LXDE-pi/autostart <<EOF
@/home/pi/on-start-7688.sh
EOF

echo "You can reboot now"
