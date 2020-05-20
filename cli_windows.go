//
// go.cli :: cli_windows.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package cli

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func (ui *CLI) title(title string) error {
	p, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return err
	}
	return setConsoleTitle(p)
}

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	pSetConsoleTitle = kernel32.NewProc("SetConsoleTitleW")
)

func setConsoleTitle(title *uint16) (err error) {
	r1, _, e1 := pSetConsoleTitle.Call(uintptr(unsafe.Pointer(title)))
	if r1 == 0 {
		if e1.(syscall.Errno) != 0 {
			err = e1
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
