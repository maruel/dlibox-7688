#!/bin/bash
# Copyright 2020 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

set -eu

xset s off
xset -dpms
xset s noblank

unclutter -idle 0 &

# Gross hack to remove the 'Chromium crashed' infobar.
sed -i 's/"exited_cleanly":false/"exited_cleanly":true/' '~/.config/chromium/Local State' || true
sed -i 's/"exited_cleanly":false/"exited_cleanly":true/; s/"exit_type":"[^"]\+"/"exit_type":"Normal"/' ~/.config/chromium/Default/Preferences || true

# --touch-events=enabled
chromium-browser --start-fullscreen --kiosk --no-first-run \
  --disable-tab-switcher --disable-translate --disable-infobars \
  --disable-session-crashed-bubble \
  http://127.0.0.1/
