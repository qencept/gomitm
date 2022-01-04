package destination

import (
	"encoding/binary"
	"net"
	"os"
	"syscall"
	"unsafe"
)

const (
	SO_ORIGINAL_DST = 0x50
)

func Detect(c *net.TCPConn) (*net.TCPAddr, error) {
	var level int
	var p unsafe.Pointer
	var raw4 syscall.RawSockaddrInet4
	var raw6 syscall.RawSockaddrInet6
	var size uintptr

	la := c.LocalAddr().(*net.TCPAddr)
	if la.IP.To4() != nil {
		level = syscall.IPPROTO_IP
		p = unsafe.Pointer(&raw4)
		size = unsafe.Sizeof(raw4)
	} else if la.IP.To16() != nil {
		level = syscall.IPPROTO_IPV6
		p = unsafe.Pointer(&raw6)
		size = unsafe.Sizeof(raw6)
	}

	rc, err := c.SyscallConn()
	if err != nil {
		return nil, err
	}

	fn := func(fd uintptr) {
		if _, _, errno := syscall.Syscall6(syscall.SYS_GETSOCKOPT, fd, uintptr(level), uintptr(SO_ORIGINAL_DST), uintptr(p), uintptr(unsafe.Pointer(&size)), 0); errno != 0 {
			err = os.NewSyscallError("getsockopt", err)
		}
	}

	if err := rc.Control(fn); err != nil {
		return nil, err
	}

	od := new(net.TCPAddr)
	switch p {
	case unsafe.Pointer(&raw4):
		od.IP = make(net.IP, net.IPv4len)
		copy(od.IP, raw4.Addr[:])
		od.Port = int(binary.BigEndian.Uint16((*[2]byte)(unsafe.Pointer(&raw4.Port))[:]))
	case unsafe.Pointer(&raw6):
		od.IP = make(net.IP, net.IPv6len)
		copy(od.IP, raw6.Addr[:])
		od.Port = int(binary.BigEndian.Uint16((*[2]byte)(unsafe.Pointer(&raw6.Port))[:]))
	}

	return od, nil
}
