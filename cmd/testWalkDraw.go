// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"syscall"
	"os"
	"unsafe"
	"image"
	"io/ioutil"
	"strings"
	"draw2d.googlecode.com/hg/draw2d"
	"draw2d.googlecode.com/hg/postscript"
	"draw2d.googlecode.com/hg/wingui"
)

// some help functions

func abortf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
	os.Exit(1)
}

func abortErrNo(funcname string, err int) {
	abortf("%s failed: %d %s\n", funcname, err, syscall.Errstr(err))
}

// global vars


func TestDrawCubicCurve(gc draw2d.GraphicContext) {
	// draw a cubic curve
	x, y := 25.6, 128.0
	x1, y1 := 102.4, 230.4
	x2, y2 := 153.6, 25.6
	x3, y3 := 230.4, 128.0

	gc.SetFillColor(image.NRGBAColor{0xAA, 0xAA, 0xAA, 0xFF})
	gc.SetLineWidth(10)
	gc.MoveTo(x, y)
	gc.CubicCurveTo(x1, y1, x2, y2, x3, y3)
	gc.Stroke()

	gc.SetStrokeColor(image.NRGBAColor{0xFF, 0x33, 0x33, 0x88})

	gc.SetLineWidth(6)
	// draw segment of curve
	gc.MoveTo(x, y)
	gc.LineTo(x1, y1)
	gc.LineTo(x2, y2)
	gc.LineTo(x3, y3)
	gc.Stroke()
}

var (
	mh                uint32
	wndBufferHeader   uint32
	wndBuffer         wingui.BITMAP
	hdcWndBuffer      uint32
	ppvBits           *image.RGBAColor
	backBuffer        *image.RGBA
	postscriptContent string
)

// WinProc called by windows to notify us of all windows events we might be interested in.
func WndProc(hwnd, msg uint32, wparam, lparam int32) uintptr {
	var rc int32

	switch msg {
	case wingui.WM_CREATE:
		hdc := wingui.GetDC(hwnd)
		wndBufferHeader = wingui.CreateCompatibleBitmap(hdc, 600, 800)
		wingui.GetObject(wndBufferHeader, unsafe.Sizeof(wndBuffer), uintptr(unsafe.Pointer(&wndBuffer)))
		hdcWndBuffer = wingui.CreateCompatibleDC(hdc)
		wingui.SelectObject(hdcWndBuffer, wndBufferHeader)

		var bmp_header wingui.BITMAPINFOHEADER
		bmp_header.Size = uint32(unsafe.Sizeof(bmp_header))
		bmp_header.Width = 600
		bmp_header.Height = 800
		bmp_header.SizeImage = 0 // the api says this must be 0 for BI_RGB images
		bmp_header.Compression = wingui.BI_RGB
		bmp_header.BitCount = 32
		bmp_header.Planes = 1
		bmp_header.XPelsPerMeter = 0
		bmp_header.YPelsPerMeter = 0
		bmp_header.ClrUsed = 0
		bmp_header.ClrImportant = 0
		//bitmap info
		var bmpinfo wingui.BITMAPINFO
		bmpinfo.Colors[0].Blue = 0
		bmpinfo.Colors[0].Green = 0
		bmpinfo.Colors[0].Red = 0
		bmpinfo.Colors[0].Reserved = 0
		bmpinfo.Header = bmp_header
		wndBufferHeader = wingui.CreateDIBSection(hdc, &bmpinfo, wingui.DIB_RGB_COLORS, uintptr(unsafe.Pointer(&ppvBits)), 0, 0)
		wingui.GetObject(wndBufferHeader, unsafe.Sizeof(wndBufferHeader), uintptr(unsafe.Pointer(&wndBuffer)))
		hdcWndBuffer = wingui.CreateCompatibleDC(hdc)
		wingui.SelectObject(hdcWndBuffer, wndBufferHeader)

		pixel := (*[600 * 800]image.RGBAColor)(unsafe.Pointer(ppvBits))
		pixelSlice := pixel[:]
		backBuffer = &image.RGBA{pixelSlice, 600, image.Rect(0, 0, 600, 800)}
		fmt.Println("Create windows")
		rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
	case wingui.WM_COMMAND:
		switch uint32(lparam) {
		default:
			rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
		}
	case wingui.WM_PAINT:
		var ps wingui.PAINTSTRUCT
		hdc := wingui.BeginPaint(hwnd, &ps)
		gc := draw2d.NewImageGraphicContext(backBuffer)
		gc.SetFillColor(image.RGBAColor{0xFF, 0xFF, 0xFF, 0xFF})
		gc.Clear()
		gc.Save()
		//gc.Translate(0, -380)
		interpreter := postscript.NewInterpreter(gc)
		reader := strings.NewReader(postscriptContent)
		interpreter.Execute(reader)
		gc.Restore()
		wingui.BitBlt(hdc, 0, 0, int(wndBuffer.Width), int(wndBuffer.Height), hdcWndBuffer, 0, 0, wingui.SRCCOPY)
		wingui.EndPaint(hwnd, &ps)
		rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
	case wingui.WM_CLOSE:
		wingui.DestroyWindow(hwnd)
	case wingui.WM_DESTROY:
		wingui.PostQuitMessage(0)
	default:
		rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
	}
	return uintptr(rc)
}

func rungui() int {
	var e int

	// GetModuleHandle
	mh, e = wingui.GetModuleHandle(nil)
	if e != 0 {
		abortErrNo("GetModuleHandle", e)
	}

	// Get icon we're going to use.
	myicon, e := wingui.LoadIcon(0, wingui.IDI_APPLICATION)
	if e != 0 {
		abortErrNo("LoadIcon", e)
	}

	// Get cursor we're going to use.
	mycursor, e := wingui.LoadCursor(0, wingui.IDC_ARROW)
	if e != 0 {
		abortErrNo("LoadCursor", e)
	}

	// Create callback
	wproc := syscall.NewCallback(WndProc)

	// RegisterClassEx
	wcname := syscall.StringToUTF16Ptr("Test Draw2d")
	var wc wingui.Wndclassex
	wc.Size = uint32(unsafe.Sizeof(wc))
	wc.WndProc = wproc
	wc.Instance = mh
	wc.Icon = myicon
	wc.Cursor = mycursor
	wc.Background = wingui.COLOR_BTNFACE + 1
	wc.MenuName = nil
	wc.ClassName = wcname
	wc.IconSm = myicon
	if _, e := wingui.RegisterClassEx(&wc); e != 0 {
		abortErrNo("RegisterClassEx", e)
	}

	// CreateWindowEx
	wh, e := wingui.CreateWindowEx(
		wingui.WS_EX_CLIENTEDGE,
		wcname,
		syscall.StringToUTF16Ptr("My window"),
		wingui.WS_OVERLAPPEDWINDOW,
		wingui.CW_USEDEFAULT, wingui.CW_USEDEFAULT, 600, 800,
		0, 0, mh, 0)
	if e != 0 {
		abortErrNo("CreateWindowEx", e)
	}
	fmt.Printf("main window handle is %x\n", wh)

	// ShowWindow
	wingui.ShowWindow(wh, wingui.SW_SHOWDEFAULT)

	// UpdateWindow
	if e := wingui.UpdateWindow(wh); e != 0 {
		abortErrNo("UpdateWindow", e)
	}

	// Process all windows messages until WM_QUIT.
	var m wingui.Msg
	for {
		r, e := wingui.GetMessage(&m, 0, 0, 0)
		if e != 0 {
			abortErrNo("GetMessage", e)
		}
		if r == 0 {
			// WM_QUIT received -> get out
			break
		}
		wingui.TranslateMessage(&m)
		wingui.DispatchMessage(&m)
	}
	return int(m.Wparam)
}

func main() {
	src, err := os.Open("../resource/postscript/tiger.ps", 0, 0)
	if err != nil {
		fmt.Println("can't find postscript file.")
		return
	}
	defer src.Close()
	bytes, err := ioutil.ReadAll(src)
	postscriptContent = string(bytes)
	rc := rungui()
	os.Exit(rc)
}
