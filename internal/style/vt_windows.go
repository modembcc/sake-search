//go:build windows

package style

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var (
	kernel32                    = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode          = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode          = kernel32.NewProc("SetConsoleMode")
	procGetFileInformationByHEx = kernel32.NewProc("GetFileInformationByHandleEx")
)

const enableVirtualTerminalProcessing = 0x0004

// isColorTerminal reports whether stdout can render ANSI color. It handles
// two cases on Windows: a real console host (cmd.exe/PowerShell/Windows
// Terminal), where ANSI processing must be explicitly enabled, and an
// MSYS2/Cygwin pty (Git Bash, as used by VS Code's integrated terminal),
// where stdout is actually a named pipe and always renders ANSI already.
func isColorTerminal() bool {
	handle := syscall.Handle(os.Stdout.Fd())

	var mode uint32
	ret, _, _ := procGetConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	if ret != 0 {
		// Real Windows console. Try to turn on VT processing; modern hosts
		// (Windows Terminal, PowerShell 7) already have it on and ignore
		// a redundant set, older conhost.exe hosts need the explicit call.
		procSetConsoleMode.Call(uintptr(handle), uintptr(mode|enableVirtualTerminalProcessing))
		return true
	}

	return isMSYSPty(handle)
}

// isMSYSPty reports whether handle is an MSYS2/Cygwin pseudo-terminal pipe,
// which GetConsoleMode can't see but which does support ANSI passthrough.
func isMSYSPty(handle syscall.Handle) bool {
	const fileNameInfo = 2
	var buf [1024]byte

	ret, _, _ := procGetFileInformationByHEx.Call(
		uintptr(handle),
		uintptr(fileNameInfo),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	if ret == 0 {
		return false
	}

	nameLen := *(*uint32)(unsafe.Pointer(&buf[0]))
	if nameLen == 0 || int(nameLen) > len(buf)-4 {
		return false
	}

	u16 := make([]uint16, nameLen/2)
	for i := range u16 {
		off := 4 + i*2
		u16[i] = uint16(buf[off]) | uint16(buf[off+1])<<8
	}
	return looksLikeMSYSPtyName(syscall.UTF16ToString(u16))
}

// looksLikeMSYSPtyName reports whether a pipe name matches the pattern MSYS2
// (Git Bash/mintty) and Cygwin use for their pty pipes, e.g.
// "\msys-1888ae32e00d56aa-pty0-to-master".
func looksLikeMSYSPtyName(name string) bool {
	name = strings.ToLower(name)
	return (strings.Contains(name, "msys-") || strings.Contains(name, "cygwin-")) && strings.Contains(name, "-pty")
}
