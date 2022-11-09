package screen

import (
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

type HKL uintptr

var (
	libuser32 *windows.LazyDLL

	mapVirtualKey          *windows.LazyProc
	mapVirtualKeyEx        *windows.LazyProc
	getKeyboardLayout      *windows.LazyProc
	getKeyboardLayoutNameW *windows.LazyProc

	vkKeyScanEx *windows.LazyProc
	vkKeyScanW  *windows.LazyProc
)

func init() {
	libuser32 = windows.NewLazySystemDLL("user32.dll")

	mapVirtualKey = libuser32.NewProc("MapVirtualKeyW")
	mapVirtualKeyEx = libuser32.NewProc("MapVirtualKeyExW")
	getKeyboardLayout = libuser32.NewProc("GetKeyboardLayout")
	vkKeyScanEx = libuser32.NewProc("VkKeyScanExW")
	vkKeyScanW = libuser32.NewProc("VkKeyScanW")
	getKeyboardLayoutNameW = libuser32.NewProc("GetKeyboardLayoutNameW")
}

const (
	MAPVK_VK_TO_VSC    = 0x0
	MAPVK_VSC_TO_VK    = 0x1
	MAPVK_VK_TO_CHAR   = 0x2
	MAPVK_VSC_TO_VK_EX = 0x3
	MAPVK_VK_TO_VSC_EX = 0x4

	KL_NAMELENGTH = 9
)

// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-mapvirtualkeya
func MapVirtualKey(uCode, uMapType uint) uint {
	ret, _, _ := mapVirtualKey.Call(
		uintptr(uCode),
		uintptr(uMapType),
	)
	return uint(ret)
}

// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-mapvirtualkeyexw
func MapVirtualKeyEx(uCode, uMapType uint, dwhkl HKL) uint {
	ret, _, _ := mapVirtualKeyEx.Call(
		uintptr(uCode),
		uintptr(uMapType),
		uintptr(dwhkl),
	)
	return uint(ret)
}

// See https://docs.microsoft.com/en-us/windows/desktop/api/winuser/nf-winuser-vkkeyscanw
func VkKeyScanW(char uint16) int16 {
	ret, _, _ := vkKeyScanW.Call(
		uintptr(char),
	)
	return int16(ret)
}

// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-vkkeyscanexw
func VkKeyScanEx(char uint16, hkl HKL) int16 {
	ret, _, _ := vkKeyScanEx.Call(
		uintptr(char),
		uintptr(hkl),
	)
	return int16(ret)
}

// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeyboardlayout
func GetKeyboardLayout(idThread uint32) HKL {
	ret, _, _ := getKeyboardLayout.Call(
		uintptr(idThread),
	)
	return HKL(ret)
}

// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeyboardlayout
func GetKeyboardLayoutNameW() string {
	buf := make([]byte, 2*KL_NAMELENGTH)
	_, _, _ = getKeyboardLayoutNameW.Call(uintptr(unsafe.Pointer(&buf[0])))
	return Utf16zToString(buf)
}

func LayoutHasAltGr(aLayout HKL) bool {
	for i := 32; i < 256; i++ {
		s := VkKeyScanEx(uint16(i), aLayout)
		if s != -1 && (s&0x600) == 0x600 { // In this context, Ctrl+Alt means AltGr.
			return true
		}
	}
	return false
}

func Utf16zToString(in []byte) string {
	out := make([]uint16, 0, len(in)/2)
	x := uint16(0)
	for i, b := range in {
		if i%2 == 0 {
			x = uint16(b)
		} else {
			x += uint16(b) << 8
			if x == 0 {
				break
			}
			out = append(out, x)
		}
	}
	return string(utf16.Decode(out))
}
