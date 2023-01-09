package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <IOKit/graphics/IOGraphicsLib.h>
#include <ApplicationServices/ApplicationServices.h>

extern int DisplayServicesSetBrightness(CGDirectDisplayID id, float brightness)
	__attribute__((weak_import));
const char *APP_NAME;
static void errexit(const char *fmt, ...) {
	va_list ap;
	va_start(ap, fmt);
	fprintf(stderr, "%s: ", APP_NAME);
	vfprintf(stderr, fmt, ap);
	fprintf(stderr, "\n");
	exit(1);
}



static void setBrightness(CGDirectDisplayID dspy, float brightness) {
	if ((DisplayServicesSetBrightness != NULL) &&
		!DisplayServicesSetBrightness(dspy, brightness)) {
	}
}


void setMacOSBrightness(float bright) {
	CGDirectDisplayID display[16];
	CGDisplayCount numDisplays;
	CGDisplayErr err;
	err = CGGetOnlineDisplayList(16, display, &numDisplays);
	if (err != CGDisplayNoErr)
		errexit("cannot get list of displays (error %d)\n", err);

	for (CGDisplayCount i = 0; i < numDisplays; ++i) {
		CGDirectDisplayID dspy = display[i];
		CGDisplayModeRef mode = CGDisplayCopyDisplayMode(dspy);
		if (mode == NULL)
			continue;

		CGDisplayModeRelease(mode);

		setBrightness(dspy, bright);
	}
}
#cgo LDFLAGS: -framework IOKit -framework ApplicationServices -framework CoreDisplay -F /System/Library/PrivateFrameworks -framework DisplayServices
*/
import "C"

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var addr = flag.String("addr", "192.168.10.5:8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		f, err := strconv.ParseFloat(string(message), 32)
		if err == nil {
			fmt.Println(string(message))
			C.setMacOSBrightness(C.float(f))
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
