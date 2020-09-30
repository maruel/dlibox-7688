// Copyright 2020 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"
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
			let s = "💣";
			s += (10 > d.getHours  () ? "0" : "") + d.getHours  () + ":";
			s += (10 > d.getMinutes() ? "0" : "") + d.getMinutes() + ":";
			s += (10 > d.getSeconds() ? "0" : "") + d.getSeconds();
			s += "❤️";
			textElem.textContent = s;
			setTimeout(updateClock, 1000 - d.getTime() % 1000 + 20);
		}
		function updateTextSize() {
			let curFontSize = 20;
			// Iterating 5 times is the sweet spot. The other option is to calculate
			// the relative update.
			for (let i = 0; i < 5; i++) {
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
	// Disable log annotation, systemd already decorates with timestamp.
	log.SetFlags(0)

	bind := flag.String("http", "127.0.0.1:80", "listening port")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Stop on SIGTERM / SIGINT.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()
	// Stop on executable file modification (e.g. replaced via rsync).
	if err := watchFile(ctx, cancel); err != nil {
		log.Fatal(err)
	}

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
	ln, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Serving on %s", ln.Addr())
	srv := &http.Server{}
	go srv.Serve(ln)
	<-ctx.Done()
}

// watchFile cancels a context when the process' executable is modified.
func watchFile(ctx context.Context, cancel func()) error {
	n, err := os.Executable()
	if err != nil {
		return err
	}
	fi, err := os.Stat(n)
	if err != nil {
		return err
	}
	mod0 := fi.ModTime()
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err = w.Add(n); err != nil {
		return err
	}
	go func() {
		defer w.Close()
		done := ctx.Done()
		for {
			select {
			case <-done:
				return
			case err := <-w.Errors:
				log.Printf("watching %s failed: %v", n, err)
				return
			case <-w.Events:
				if fi, err = os.Stat(n); err != nil || !fi.ModTime().Equal(mod0) {
					log.Printf("%s was modified, exiting", n)
					cancel()
				}
			}
		}
	}()
	return nil
}
