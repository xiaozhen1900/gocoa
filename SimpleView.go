package main

/*
#cgo CFLAGS: -I/System/Library/Frameworks/CoreGraphics.framework/Versions/A/Headers/
#cgo LDFLAGS: -framework AppKit 
#include <objc/objc-runtime.h>
#include <CoreGraphics.h>
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	. "gocoa"
	"unsafe"
)

/*
* The following is an example of something that I don't understand well and represents a problem:
* the exported symbol names seem to determine whether the funtions can be found by dlsym.
* If all are named similarly, only one is found. If they are dissimilar enough in their first
* character, they apparently can be loaded.
*
* Possibly, there's some asynchronous operation that requires the time to iterate over the
* symbol table between calls, and maybe changing the name exposes that behavior.
*		solution: calling dlopen in a blocking goroutine with a wait() call?
*
* One possible solution is trying not to repeatedly call dlopen/dlclose
 */
 
//export BIsOpaque
func BIsOpaque(self C.id, op C.SEL) C.BOOL {
	fmt.Println("isOpaque")
	return (C.BOOL)(1)
}

//export IInitWithFrame
func IInitWithFrame(self C.id, op C.SEL, aRect C.CGRect) C.id {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, aRect)

	simpleView := ObjectForId((uintptr)(unsafe.Pointer(self)))
	simpleView = simpleView.Class().Instance("alloc")
	simpleView = simpleView.CallSuperR("initWithFrame:", buf.Bytes())
	return (C.id)(unsafe.Pointer(simpleView.Pointer))
}

//export VDrawRect
func VDrawRect(self C.id, op C.SEL, aRect C.CGRect) {
	fmt.Println("drawRect:")

	view := ObjectForId((uintptr)(unsafe.Pointer(self)))
	view.ListInstanceVariables()
	view.ListMethods()
	
	view.Class().Super().ListMethods()

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, aRect)
	fmt.Println("len(aRectBytes)", len(buf.Bytes()))
	fmt.Println("bytes:", buf.Bytes())

	NSColor(RedColor).Call("set")
	
	
//	"To kill you must know your enemy. And in this case, my enemy, is a varmint."
	
	
	bezier := ClassForName("NSBezierPath").Instance("init")
	bezier.ListMethods()
	
	bezier.CallR("fillRect:", buf.Bytes()) // terrible here

}

//export ZWindowResize
func ZWindowResize(self C.id, op C.SEL, notification C.id) {
	fmt.Println("windowDidResize:")
	controller := ObjectForId((uintptr)(unsafe.Pointer(self)))
	simpleView := controller.InstanceVariable("itsView")
	simpleView.Call("invalidateIntrinsicContentSize")
}

/*
* main()
* Main function for testing
 */
func main() {

	view := ClassForName("NSView").Subclass("SimpleView")
	view.AddMethod("initWithFrame:", IInitWithFrame)
	view.AddMethod("drawRect:", VDrawRect)
	view.AddMethod("isOpaque", BIsOpaque)		// XXX
	view.Register()

	controller := ClassForName("NSObject").Subclass("Controller")
	controller.AddMethod("windowDidResize:", ZWindowResize)
	controller.AddIvar("itsView", view)
	controller.Register()

	app := ClassForName("NSApplication").Instance("sharedApplication")
	bundle := ClassForName("NSBundle").Instance("alloc")
	path := NSString(".")
	dict := NSDictionary("NSOwner", app)

	bundle = bundle.Call("initWithPath:", path)
	bundle.Call("loadNibFile:externalNameTable:withZone:", NSString("SimpleView"), dict, app.Call("zone"))

	app.Call("run")

}
