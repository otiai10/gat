// +build windows

package terminal

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

type wchar uint16
type short int16
type dword uint32
type word uint16

type coord struct {
	x short
	y short
}

type smallRect struct {
	left   short
	top    short
	right  short
	bottom short
}

type consoleScreenBufferInfo struct {
	size              coord
	cursorPosition    coord
	attributes        word
	window            smallRect
	maximumWindowSize coord
}

// GetTerminal ...
func GetTerminal() Rect {
	t := new(Rect)

	var csbi consoleScreenBufferInfo
	r := syscall.Handle(os.Stdout.Fd())
	r1, _, _ := procGetConsoleScreenBufferInfo.Call(uintptr(r), uintptr(unsafe.Pointer(&csbi)))
	if r1 == 0 {
		t.Col = 80
		t.Row = 25
	} else {
		t.Col = uint16(csbi.window.right - csbi.window.left)
		t.Row = uint16(csbi.window.bottom - csbi.window.top)
	}

	return *t
}
