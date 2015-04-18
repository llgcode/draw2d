// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/wingui"
	"github.com/llgcode/ps"
)

// some help functions

func abortf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
	os.Exit(1)
}

func abortErrNo(funcname string, err error) {
	abortf("%s failed: %d %s\n", funcname, err, err)
}

// global vars

func TestDrawCubicCurve(gc draw2d.GraphicContext) {
	// draw a cubic curve
	x, y := 25.6, 128.0
	x1, y1 := 102.4, 230.4
	x2, y2 := 153.6, 25.6
	x3, y3 := 230.4, 128.0

	gc.SetStrokeColor(color.NRGBA{0, 0, 0, 0xFF})
	gc.SetLineWidth(10)
	gc.MoveTo(x, y)
	gc.CubicCurveTo(x1, y1, x2, y2, x3, y3)
	gc.Stroke()

	gc.SetStrokeColor(color.NRGBA{0xFF, 0, 0, 0xFF})

	gc.SetLineWidth(6)
	// draw segment of curve
	gc.MoveTo(x, y)
	gc.LineTo(x1, y1)
	gc.LineTo(x2, y2)
	gc.LineTo(x3, y3)
	gc.Stroke()
}

func DrawTiger(gc draw2d.GraphicContext) {
	if postscriptContent == "" {
		src, err := os.OpenFile("../../../ps/samples/tiger.ps", 0, 0)
		if err != nil {
			fmt.Println("can't find postscript file.")
			return
		}
		defer src.Close()
		bytes, err := ioutil.ReadAll(src)
		postscriptContent = string(bytes)
	}
	interpreter := ps.NewInterpreter(gc)
	reader := strings.NewReader(postscriptContent)
	interpreter.Execute(reader)
}

var (
	mh                syscall.Handle
	hdcWndBuffer      syscall.Handle
	wndBufferHeader   syscall.Handle
	wndBuffer         wingui.BITMAP
	ppvBits           *uint8
	backBuffer        *image.RGBA
	postscriptContent string
)

// WinProc called by windows to notify us of all windows events we might be interested in.
func WndProc(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) (rc uintptr) {
	_ = make([]int, 100000)
	switch msg {
	case wingui.WM_CREATE:
		hdc := wingui.GetDC(hwnd)
		hdcWndBuffer = wingui.CreateCompatibleDC(hdc)
		wndBufferHeader = wingui.CreateCompatibleBitmap(hdc, 600, 800)
		wingui.GetObject(wndBufferHeader, unsafe.Sizeof(wndBuffer), uintptr(unsafe.Pointer(&wndBuffer)))
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

		pixel := (*[600 * 800 * 4]uint8)(unsafe.Pointer(ppvBits))
		pixelSlice := pixel[:]
		backBuffer = &image.RGBA{pixelSlice, 4 * 600, image.Rect(0, 0, 600, 800)}
		fmt.Println("Create windows")
		rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
	case wingui.WM_COMMAND:
		switch syscall.Handle(lparam) {
		default:
			rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
		}
	case wingui.WM_PAINT:
		var ps wingui.PAINTSTRUCT
		hdc := wingui.BeginPaint(hwnd, &ps)
		t := time.Now()
		gc := draw2d.NewGraphicContext(backBuffer)
		/*gc.SetFillColor(color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
		gc.Clear()*/
		for i := 0; i < len(backBuffer.Pix); i += 1 {
			backBuffer.Pix[i] = 0xff
		}
		gc.Save()
		//gc.Translate(0, -380)
		DrawTiger(gc)
		gc.Restore()
		// back buf in
		var tmp uint8
		for i := 0; i < len(backBuffer.Pix); i += 4 {
			tmp = backBuffer.Pix[i]
			backBuffer.Pix[i] = backBuffer.Pix[i+2]
			backBuffer.Pix[i+2] = tmp
		}
		wingui.BitBlt(hdc, 0, 0, int(wndBuffer.Width), int(wndBuffer.Height), hdcWndBuffer, 0, 0, wingui.SRCCOPY)
		wingui.EndPaint(hwnd, &ps)
		rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
		dt := time.Now().Sub(t)
		fmt.Printf("Redraw in : %f ms\n", float64(dt)*1e-6)
	case wingui.WM_CLOSE:
		wingui.DestroyWindow(hwnd)
	case wingui.WM_DESTROY:
		wingui.PostQuitMessage(0)
	default:
		rc = wingui.DefWindowProc(hwnd, msg, wparam, lparam)
	}
	return
}

func rungui() int {
	var e error

	// GetModuleHandle
	mh, e = wingui.GetModuleHandle(nil)
	if e != nil {
		abortErrNo("GetModuleHandle", e)
	}

	// Get icon we're going to use.
	myicon, e := wingui.LoadIcon(0, wingui.IDI_APPLICATION)
	if e != nil {
		abortErrNo("LoadIcon", e)
	}

	// Get cursor we're going to use.
	mycursor, e := wingui.LoadCursor(0, wingui.IDC_ARROW)
	if e != nil {
		abortErrNo("LoadCursor", e)
	}

	// Create callback
	wproc := syscall.NewCallback(WndProc)

	// RegisterClassEx
	wcname := syscall.StringToUTF16Ptr("myWindowClass")
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
	if _, e := wingui.RegisterClassEx(&wc); e != nil {
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
	if e != nil {
		abortErrNo("CreateWindowEx", e)
	}
	fmt.Printf("main window handle is %x\n", wh)

	// ShowWindow
	wingui.ShowWindow(wh, wingui.SW_SHOWDEFAULT)

	// UpdateWindow
	if e := wingui.UpdateWindow(wh); e != nil {
		abortErrNo("UpdateWindow", e)
	}

	// Process all windows messages until WM_QUIT.
	var m wingui.Msg
	for {
		r, e := wingui.GetMessage(&m, 0, 0, 0)
		if e != nil {
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
	rc := rungui()
	os.Exit(rc)
}
