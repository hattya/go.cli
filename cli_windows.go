//
// go.cli :: cli_windows.go
//
//   Copyright (c) 2014-2016 Akinori Hattori <hattya@gmail.com>
//
//   Permission is hereby granted, free of charge, to any person
//   obtaining a copy of this software and associated documentation files
//   (the "Software"), to deal in the Software without restriction,
//   including without limitation the rights to use, copy, modify, merge,
//   publish, distribute, sublicense, and/or sell copies of the Software,
//   and to permit persons to whom the Software is furnished to do so,
//   subject to the following conditions:
//
//   The above copyright notice and this permission notice shall be
//   included in all copies or substantial portions of the Software.
//
//   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//   EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
//   MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//   NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
//   BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
//   ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
//   CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//   SOFTWARE.
//

package cli

import (
	"syscall"
	"unsafe"
)

func (ui *CLI) title(title string) error {
	p, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return err
	}
	return setConsoleTitle(p)
}

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

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
