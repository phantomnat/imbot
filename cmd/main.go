package main

import (
	"fmt"

	"github.com/lxn/win"
	"github.com/phantomnat/imbot/pkg/screen"
)

func main() {
	key := 'l'
	layout := screen.GetKeyboardLayout(win.GetWindowThreadProcessId(win.GetForegroundWindow(), nil))
	//layoutName := screen.GetKeyboardLayoutNameW()
	vk := screen.VkKeyScanW(uint16(key))
	scan := screen.MapVirtualKey(uint(vk)&0xFF, screen.MAPVK_VK_TO_VSC)

	usLayout := screen.HKL(0x4090409)

	usVk := screen.VkKeyScanEx(uint16(key), usLayout)
	//usScan := screen.MapVirtualKeyEx(uint(usVk), screen.MAPVK_VK_TO_VSC, usLayout)
	usScan := screen.MapVirtualKey(uint(usVk), screen.MAPVK_VK_TO_VSC)

	fmt.Printf("\nkey      : %c", key)
	fmt.Printf("\n   layout: %x, vk: %c (%x), scan: %c (%x)", layout, vk, vk, scan, scan)
	fmt.Printf("\nus layout: %x, vk: %c (%x), scan: %c (%x)", usLayout, usVk, usVk, usScan, usScan)
}
