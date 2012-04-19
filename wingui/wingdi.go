package wingui

import (
	"syscall"
	"unsafe"
)

const (
	WM_PAINT = 15

	BI_RGB = 0
	BI_BITFIELDS = 3

	DIB_PAL_COLORS = 1
	DIB_RGB_COLORS = 0

	BLACKNESS = 0x42
	DSTINVERT = 0x550009
	MERGECOPY = 0xC000CA
	MERGEPAINT = 0xBB0226
	NOTSRCCOPY = 0x330008
	NOTSRCERASE = 0x1100A6
	PATCOPY = 0xF00021
	PATINVERT = 0x5A0049
	PATPAINT = 0xFB0A09
	SRCAND = 0x8800C6
	SRCCOPY = 0xCC0020
	SRCERASE = 0x440328
	SRCINVERT = 0x660046
	SRCPAINT = 0xEE0086
	WHITENESS = 0xFF0062
)

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type PAINTSTRUCT struct {
	HDC         syscall.Handle
	Erase       int32 // bool
	RcPaint     RECT
	Restore     int32 // bool
	IncUpdate   int32 // bool
	rgbReserved [32]byte
}

type BITMAP struct {
	Type       int32
	Width      int32
	Height     int32
	WidthBytes int32
	Planes     uint16
	BitsPixel  uint16
	Bits       *byte
}

type BITMAPINFOHEADER struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

type BITMAPINFO struct {
	Header BITMAPINFOHEADER
	Colors [1]RGBQUAD
}

type RGBQUAD struct {
	Blue     byte
	Green    byte
	Red      byte
	Reserved byte
}

var (
	modgdi32 = syscall.NewLazyDLL("gdi32.dll")

	procGetDC                  = moduser32.NewProc("GetDC")
	procCreateCompatibleDC     = modgdi32.NewProc("CreateCompatibleDC")
	procGetObject              = modgdi32.NewProc("GetObjectW")
	procSelectObject           = modgdi32.NewProc("SelectObject")
	procBeginPaint             = moduser32.NewProc("BeginPaint")
	procEndPaint               = moduser32.NewProc("EndPaint")
	procCreateCompatibleBitmap = modgdi32.NewProc("CreateCompatibleBitmap")
	procCreateDIBSection       = modgdi32.NewProc("CreateDIBSection")
	procBitBlt                 = modgdi32.NewProc("BitBlt")
)

func GetDC(hwnd syscall.Handle) (hdc syscall.Handle) {
	r0, _, _ := syscall.Syscall(procGetDC.Addr(), 1, uintptr(hwnd), 0, 0)
	hdc = syscall.Handle(r0)
	return hdc
}

func CreateCompatibleDC(hwnd syscall.Handle) (hdc syscall.Handle) {
	r0, _, _ := syscall.Syscall(procCreateCompatibleDC.Addr(), 1, uintptr(hwnd), 0, 0)
	hdc = syscall.Handle(r0)
	return hdc
}

func GetObject(hgdiobj syscall.Handle, cbBuffer uintptr, object uintptr) (size uint32) {
	r0, _, _ := syscall.Syscall(procGetObject.Addr(), 3, uintptr(hgdiobj), uintptr(cbBuffer), object)
	size = uint32(r0)
	return size
}

func SelectObject(hdc syscall.Handle, hgdiobj syscall.Handle) syscall.Handle {
	r0, _, _ := syscall.Syscall(procSelectObject.Addr(), 2, uintptr(hdc), uintptr(hgdiobj), 0)
	return syscall.Handle(r0)
}

func BeginPaint(hwnd syscall.Handle, ps *PAINTSTRUCT) (hdc syscall.Handle){
	r0, _, _ := syscall.Syscall(procBeginPaint.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(ps)), 0)
	hdc = syscall.Handle(r0)
	return
}

func EndPaint(hwnd syscall.Handle, ps *PAINTSTRUCT) bool {
	syscall.Syscall(procEndPaint.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(ps)), 0)
	return true
}

func CreateCompatibleBitmap(hdc syscall.Handle, width, height uintptr) (hbitmap syscall.Handle) {
	r0, _, _ := syscall.Syscall(procCreateCompatibleBitmap.Addr(), 3, uintptr(hdc), uintptr(width), uintptr(height))
	return syscall.Handle(r0)
}

func CreateDIBSection(hdc syscall.Handle, pbmi *BITMAPINFO, iUsage uint, ppvBits uintptr, hSection uint32, dwOffset uint32) (hbitmap syscall.Handle) {
	r0, _, _ := syscall.Syscall6(procCreateDIBSection.Addr(), 6, uintptr(hdc), uintptr(unsafe.Pointer(pbmi)), uintptr(iUsage), ppvBits, uintptr(hSection), uintptr(dwOffset))
	return syscall.Handle(r0)
}

func BitBlt(hdc syscall.Handle, nXDest, nYDest, nWidth, nHeight int, hdcSrc syscall.Handle, nXSrc, nYSrc int, dwRop uint32) bool {
	r0, _, _ := syscall.Syscall9(procBitBlt.Addr(), 9, uintptr(hdc), uintptr(nXDest), uintptr(nYDest), uintptr(nWidth), uintptr(nHeight), uintptr(hdcSrc), uintptr(nXSrc), uintptr(nYSrc), uintptr(dwRop))
	return r0 != 0
}
