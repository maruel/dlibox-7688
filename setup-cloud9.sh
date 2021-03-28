#!/bin/bash
# Copyright 2020 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

set -eu

# From https://docs.docker.com/engine/install/debian/
echo "Installing prerequisites"
if [ ! -f /usr/bin/docker ]; then
  curl -fsSL https://get.docker.com -o get-docker.sh
  sudo sh get-docker.sh
  sudo usermod -aG docker pi
  rm get-docker.sh
fi

echo "Testing"
sudo docker run hello-world

echo "Pulling Cloud9"
# Will not require sudo after reboot
sudo docker pull linuxserver/cloud9

sudo tee /etc/systemd/system/cloud9.service > /dev/null <<EOF
# https://github.com/maruel/dlibox-7688
[Unit]
Description=Runs Cloud9 automatically upon boot
Wants=network-online.target
After=network-online.target
[Service]
User=pi
Group=pi
KillMode=mixed
Restart=always
TimeoutStopSec=10s
ExecStart=docker run --name=cloud9 -e PUID=1000 -e PGID=1000 -e TZ=America/Toronto -p 8000:8000 -v /home/pi/dlibox-7688:/code linuxserver/cloud9
[Install]
WantedBy=default.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable cloud9.service
sudo systemctl start cloud9.service


#  -e GITURL=https://github.com/linuxserver/docker-cloud9.git
#  `#optional` \
#  -e USERNAME= `#optional` \
#  -e PASSWORD= `#optional` \
#  -v /var/run/docker.sock:/var/run/docker.sock
#  --restart unless-stopped
