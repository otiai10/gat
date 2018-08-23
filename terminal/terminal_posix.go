// +build !windows

package terminal

import (
	"syscall"
	"unsafe"
)

// GetTerminal ...
func GetTerminal() Rect {
	t := new(Rect)
	retCode, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(t)),
	)

	if int(retCode) == -1 {
		panic(err)
	}
	return *t
}
