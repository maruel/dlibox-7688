// Copyright 2020 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const root = `<!DOCTYPE html>
<html>
<style type="text/css">
html {
	height: 100%;
}
body {
	background: #000000;
	color: #FFFFFF;
	display: flex;
	font-family: sans-serif;
	font-weight: bold;
  align-items: center;
  height: 100%;
  justify-content: center;
  margin: 0;
  padding: 0;
}
.clock {
	font-kerning: none;
}
</style>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
<script>
	"use strict";
	function autoReload() {
		try {
			let ws = new WebSocket('ws://' + window.location.host + '/auto-reload');
			ws.addEventListener('close', function(event) {
				window.location.reload();
			})
		} catch(ex) {
			window.console('autoReload(): ', ex);
		}
	}
	window.addEventListener("DOMContentLoaded", (event) => {
		autoReload();
		let textElem = document.getElementById("clock");
		function updateClock() {
			let d = new Date();
			let s = "";
			s += (10 > d.getHours  () ? "0" : "") + d.getHours  () + ":";
			s += (10 > d.getMinutes() ? "0" : "") + d.getMinutes() + ":";
			s += (10 > d.getSeconds() ? "0" : "") + d.getSeconds();
			textElem.textContent = s;
			setTimeout(updateClock, 1000 - d.getTime() % 1000 + 20);
		}
		function updateTextSize() {
			let curFontSize = 20;
			// Iterating 3 times is the sweet spot.
			for (let i = 0; i < 3; i++) {
				// 90% width.
				curFontSize *= 0.9 / (textElem.offsetWidth / textElem.parentNode.offsetWidth);
				textElem.style.fontSize = curFontSize + "pt";
			}
		}
		
		updateClock();
		updateTextSize();
		window.addEventListener("resize", updateTextSize);
  });
</script>
<title>Clock</title>
<div class="clock" id="clock"></div>
`

func main() {
	bind := flag.String("http", "127.0.0.1:80", "listening port")
	log.SetFlags(log.Lmicroseconds)
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, root)
	})
	http.HandleFunc("/auto-reload", func(w http.ResponseWriter, req *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		if _, err := upgrader.Upgrade(w, req, nil); err != nil {
			if _, ok := err.(websocket.HandshakeError); !ok {
				log.Println(err)
			}
			return
		}
		// Hang the connection. It'll close when the server is restarted.
		select {}
	})
	log.Printf("Serving")
	log.Fatal(http.ListenAndServe(*bind, nil))
}
