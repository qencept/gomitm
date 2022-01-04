package destination

import (
	"encoding/binary"
	"net"
	"os"
	"syscall"
	"unsafe"
)

const (
	PF_IN       = 1
	PF_OUT      = 2
	DIOCNATLOOK = 0xc0544417
)

type pfiocNatlook struct {
	Saddr     [16]byte
	Daddr     [16]byte
	Rsaddr    [16]byte
	Rdaddr    [16]byte
	Sxport    [4]byte
	Dxport    [4]byte
	Rsxport   [4]byte
	Rdxport   [4]byte
	Af        uint8
	Proto     uint8
	Variant   uint8
	Direction uint8
}

func Detect(c *net.TCPConn) (*net.TCPAddr, error) {
	la := c.LocalAddr().(*net.TCPAddr)
	ra := c.RemoteAddr().(*net.TCPAddr)

	f, err := os.Open("/dev/pf")
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	fd := f.Fd()
	nl := pfiocNatlook{}
	if ra.IP.To4() != nil {
		copy(nl.Saddr[:net.IPv4len], ra.IP.To4())
		copy(nl.Daddr[:net.IPv4len], la.IP.To4())
		nl.Af = syscall.AF_INET
	}
	if ra.IP.To16() != nil && ra.IP.To4() == nil {
		copy(nl.Saddr[:], ra.IP)
		copy(nl.Daddr[:], la.IP)
		nl.Af = syscall.AF_INET6
	}
	binary.BigEndian.PutUint16((*[2]byte)(unsafe.Pointer(&nl.Sxport))[:2], uint16(ra.Port))
	binary.BigEndian.PutUint16((*[2]byte)(unsafe.Pointer(&nl.Dxport))[:2], uint16(la.Port))
	nl.Proto = syscall.IPPROTO_TCP
	ioc := uintptr(DIOCNATLOOK)
	for _, dir := range []byte{PF_OUT, PF_IN} {
		nl.Direction = dir
		err = ioctl(fd, int(ioc), unsafe.Pointer(&nl))
		if err == nil || err != syscall.ENOENT {
			break
		}
	}
	if err != nil {
		return nil, os.NewSyscallError("ioctl", err)
	}
	od := new(net.TCPAddr)
	od.Port = int(binary.BigEndian.Uint16(nl.Rdxport[:2]))
	switch nl.Af {
	case syscall.AF_INET:
		od.IP = make(net.IP, net.IPv4len)
		copy(od.IP, nl.Rdaddr[:net.IPv4len])
	case syscall.AF_INET6:
		od.IP = make(net.IP, net.IPv6len)
		copy(od.IP, nl.Rdaddr[:])
	}

	return od, nil
}

func ioctl(s uintptr, ioc int, b unsafe.Pointer) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, s, uintptr(ioc), uintptr(b)); errno != 0 {
		return error(errno)
	}
	return nil
}
