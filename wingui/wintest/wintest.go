// +build ignore

package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	. "github.com/llgcode/draw2d/wingui"
)

// some help functions

func abortf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
	os.Exit(1)
}

func abortErrNo(funcname string, err error) {
	errno, _ := err.(syscall.Errno)
	abortf("%s failed: %d %s\n", funcname, uint32(errno), err)
}

// global vars

var (
	mh syscall.Handle
	bh syscall.Handle
)

// WinProc called by windows to notify us of all windows events we might be interested in.
func WndProc(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) (rc uintptr) {
	_ = make([]int, 100000)
	switch msg {
	case WM_CREATE:
		var e error
		// CreateWindowEx
		bh, e = CreateWindowEx(
			0,
			syscall.StringToUTF16Ptr("button"),
			syscall.StringToUTF16Ptr("Quit"),
			WS_CHILD|WS_VISIBLE|BS_DEFPUSHBUTTON,
			75, 70, 140, 25,
			hwnd, 1, mh, 0)
		if e != nil {
			abortErrNo("CreateWindowEx", e)
		}
		fmt.Printf("button handle is %x\n", bh)
		rc = DefWindowProc(hwnd, msg, wparam, lparam)
	case WM_COMMAND:
		switch syscall.Handle(lparam) {
		case bh:
			e := PostMessage(hwnd, WM_CLOSE, 0, 0)
			if e != nil {
				abortErrNo("PostMessage", e)
			}
		default:
			rc = DefWindowProc(hwnd, msg, wparam, lparam)
		}
	case WM_CLOSE:
		DestroyWindow(hwnd)
	case WM_DESTROY:
		PostQuitMessage(0)
	default:
		rc = DefWindowProc(hwnd, msg, wparam, lparam)
	}
	//fmt.Printf("WndProc(0x%08x, %d, 0x%08x, 0x%08x) (%d)\n", hwnd, msg, wparam, lparam, rc)
	return
}

func rungui() int {
	var e error

	// GetModuleHandle
	mh, e = GetModuleHandle(nil)
	if e != nil {
		abortErrNo("GetModuleHandle", e)
	}

	// Get icon we're going to use.
	myicon, e := LoadIcon(0, IDI_APPLICATION)
	if e != nil {
		abortErrNo("LoadIcon", e)
	}

	// Get cursor we're going to use.
	mycursor, e := LoadCursor(0, IDC_ARROW)
	if e != nil {
		abortErrNo("LoadCursor", e)
	}

	// Create callback
	wproc := syscall.NewCallback(WndProc)

	// RegisterClassEx
	wcname := syscall.StringToUTF16Ptr("myWindowClass")
	var wc Wndclassex
	wc.Size = uint32(unsafe.Sizeof(wc))
	wc.WndProc = wproc
	wc.Instance = mh
	wc.Icon = myicon
	wc.Cursor = mycursor
	wc.Background = COLOR_BTNFACE + 1
	wc.MenuName = nil
	wc.ClassName = wcname
	wc.IconSm = myicon
	if _, e := RegisterClassEx(&wc); e != nil {
		abortErrNo("RegisterClassEx", e)
	}

	// CreateWindowEx
	wh, e := CreateWindowEx(
		WS_EX_CLIENTEDGE,
		wcname,
		syscall.StringToUTF16Ptr("My window"),
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT, 300, 200,
		0, 0, mh, 0)
	if e != nil {
		abortErrNo("CreateWindowEx", e)
	}
	fmt.Printf("main window handle is %x\n", wh)

	// ShowWindow
	ShowWindow(wh, SW_SHOWDEFAULT)

	// UpdateWindow
	if e := UpdateWindow(wh); e != nil {
		abortErrNo("UpdateWindow", e)
	}

	// Process all windows messages until WM_QUIT.
	var m Msg
	for {
		r, e := GetMessage(&m, 0, 0, 0)
		if e != nil {
			abortErrNo("GetMessage", e)
		}
		if r == 0 {
			// WM_QUIT received -> get out
			break
		}
		TranslateMessage(&m)
		DispatchMessage(&m)
	}
	return int(m.Wparam)
}

func main() {
	rc := rungui()
	os.Exit(rc)
}
