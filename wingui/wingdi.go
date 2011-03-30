package wingui

import (
	"syscall"
	"unsafe"
)

const (

	WM_PAINT = 15
	
	BI_RGB = 0
	
	DIB_PAL_COLORS = 1
	DIB_RGB_COLORS = 0
	
	SRCCOPY = 0xCC0020
)

type RECT struct{
	Left int32
	Top int32
	Right int32
	Bottom int32 
}

type PAINTSTRUCT struct{
  HDC  uint32
  Erase int32 // bool
  RcPaint RECT
  Restore int32 // bool
  IncUpdate int32 // bool
  rgbReserved [32]byte
}

type BITMAP struct{
  Type int32
  Width int32
  Height int32
  WidthBytes int32
  Planes uint16
  BitsPixel uint16
  Bits * byte
}

type BITMAPINFOHEADER struct{
  Size uint32
  Width int32
  Height int32
  Planes uint16
  BitCount uint16
  Compression uint32
  SizeImage uint32
  XPelsPerMeter int32
  YPelsPerMeter int32
  ClrUsed uint32
  ClrImportant uint32
}

type BITMAPINFO struct{
  Header BITMAPINFOHEADER
  Colors [1]RGBQUAD
}

type RGBQUAD struct{
  Blue byte
  Green byte
  Red byte
  Reserved byte
}

var (
	modgdi32 = loadDll("gdi32.dll")

	procGetDC = getSysProcAddr(moduser32, "GetDC")
	procCreateCompatibleDC = getSysProcAddr(modgdi32, "CreateCompatibleDC")
	procGetObject = getSysProcAddr(modgdi32, "GetObjectW")
	procSelectObject = getSysProcAddr(modgdi32, "SelectObject")
	procBeginPaint = getSysProcAddr(moduser32, "BeginPaint")
	procEndPaint = getSysProcAddr(moduser32, "EndPaint")
	procCreateCompatibleBitmap = getSysProcAddr(modgdi32, "CreateCompatibleBitmap")
	procCreateDIBSection = getSysProcAddr(modgdi32, "CreateDIBSection")
	procBitBlt = getSysProcAddr(modgdi32, "BitBlt")
	
)

func GetDC(hwnd uint32) (hdc uint32) {
	r0, _, _ := syscall.Syscall(procGetDC, 1, uintptr(hwnd), 0, 0)
	hdc = uint32(r0)
	return hdc
}

func CreateCompatibleDC(hwnd uint32) (hdc uint32) {
	r0, _, _ := syscall.Syscall(procCreateCompatibleDC, 1, uintptr(hwnd), 0, 0)
	hdc = uint32(r0)
	return hdc
}

func GetObject(hgdiobj uint32, cbBuffer int, object uintptr) (size uint32) {
	r0, _, _ := syscall.Syscall(procGetObject, 3, uintptr(hgdiobj), uintptr(cbBuffer), object)
	size = uint32(r0)
	return size
}

func SelectObject(hdc uint32, hgdiobj uint32) (uint32) {
	r0, _, _ := syscall.Syscall(procSelectObject, 2, uintptr(hdc), uintptr(hgdiobj), 0)
	return uint32(r0)
}

func BeginPaint(hwnd uint32, ps *PAINTSTRUCT) (hdc uint32) {
	r0, _, _ := syscall.Syscall(procBeginPaint, 2, uintptr(hwnd), uintptr(unsafe.Pointer(ps)), 0)
	hdc = uint32(r0)
	return hdc
}

func EndPaint(hwnd uint32, ps *PAINTSTRUCT) bool {
	syscall.Syscall(procEndPaint, 2, uintptr(hwnd), uintptr(unsafe.Pointer(ps)), 0)
	return true
}

func CreateCompatibleBitmap(hdc uint32, width ,height int) ( hbitmap uint32) {
	r0, _, _ := syscall.Syscall(procCreateCompatibleBitmap, 3, uintptr(hdc), uintptr(width), uintptr(height))
	return uint32(r0)
}

func CreateDIBSection(hdc uint32, pbmi *BITMAPINFO , iUsage uint, ppvBits uintptr, hSection uint32, dwOffset uint32) ( hbitmap uint32) {
	r0, _, _ := syscall.Syscall6(procCreateDIBSection, 6, uintptr(hdc), uintptr(unsafe.Pointer(pbmi)), uintptr(iUsage), ppvBits, uintptr(hSection), uintptr(dwOffset))
	return uint32(r0)
}

func BitBlt(hdc uint32, nXDest, nYDest, nWidth, nHeight int, hdcSrc uint32, nXSrc, nYSrc int, dwRop uint32) ( bool) {
	r0, _, _ := syscall.Syscall9(procBitBlt, 9, uintptr(hdc), uintptr(nXDest), uintptr(nYDest), uintptr(nWidth), uintptr(nHeight), uintptr(hdcSrc), uintptr(nXSrc), uintptr(nYSrc), uintptr(dwRop))
	return r0 != 0
}