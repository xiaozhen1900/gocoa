package main

/*
#include <stdlib.h>
#include <objc/objc-runtime.h>
#include <CoreGraphics.h>
*/
//#cgo CFLAGS: -I/usr/include -I/System/Library/Frameworks/Foundation.framework/Versions/C/Headers/ -I/System/Library/Frameworks/AppKit.framework/Versions/C/Headers/ -I/System/Library/Frameworks/ApplicationServices.framework/Versions/A/Frameworks/HIServices.framework/Versions/A/Headers/ -I/System/Library/Frameworks/CoreGraphics.framework/Versions/A/Headers/
//#cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit -framework ApplicationServices -framework CoreGraphics
import "C"

import (
	"fmt"
	"gocoa"
	"unsafe"
)

func init() {
	gocoa.InitMac()
}

/*
* main()
* Main function for testing
 */
func main() {

	hellow := gocoa.ClassForName("NSObject").Subclass("ApplicationController")
	hellow.AddMethod("applicationWillFinishLaunching:", BApplicationWillFinishLaunching)
	hellow.AddIvar("scrollTable1", gocoa.ClassForName("NSScrollView"))
	hellow.Register()

	app := gocoa.ClassForName("NSApplication").Instance("sharedApplication")
	bundle := gocoa.ClassForName("NSBundle").Instance("alloc")
	path := gocoa.NSString(".")
	dict := gocoa.NSDictionary("NSOwner", app)

	bundle = bundle.Call("initWithPath:", path)
	bundle.Call("loadNibFile:externalNameTable:withZone:", gocoa.NSString("Application"), dict, app.Call("zone"))

	app.Call("run")
}

//export BApplicationWillFinishLaunching
func BApplicationWillFinishLaunching(self C.id, op C.SEL, notification C.id) {
	fmt.Println("applicationWillFinishLaunching:")

	notify := gocoa.NewObject((uintptr)(unsafe.Pointer(notification)))
	application := notify.Call("object")

	windowsArray := application.Call("windows")
	windowsCount := (gocoa.NSUInteger)(windowsArray.Call("count").Pointer)
	var ix gocoa.NSUInteger
	for ix = 0; ix < windowsCount; ix++ {
		window := windowsArray.CallI("objectAtIndex:", ix)
		window.Call("setTitle:", gocoa.NSString("Form Loaded"))
	}

	me := gocoa.NewObject((uintptr)(unsafe.Pointer(self)))
	scrollTable1 := me.InstanceVariable("scrollTable1")
	fmt.Println("scrollTable1 class:", scrollTable1.Class().Name())
	//	textBox1.Call("setStringValue:", gocoa.NSString("Form Loaded"))
}