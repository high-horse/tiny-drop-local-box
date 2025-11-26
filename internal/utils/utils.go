package utils

import (
	"log"
	"net"
	"net/http"
	"syscall"
)


func GetUserIp(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	var err error
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	} 
	if ip == "" {
		ip, _, err = net.SplitHostPort(ip)
		if err != nil {
			// log.Printf("Error Parsing IP address : %v", err)
			// return ""
		}
	}

	return ip
}

func CheckDiskSpace(fileSize int64) bool {
	statfs := syscall.Statfs_t{}
	err := syscall.Statfs("./storage", &statfs)
	if err != nil {
		log.Println("Error checking disk space:", err)
		return false
	}

	// Get the available free space (in bytes)
	freeSpace := statfs.Bavail * uint64(statfs.Bsize)

	// Check if the available space after the upload will be enough
	if freeSpace < uint64(fileSize) {
		log.Printf("Not enough disk space! Required: %d bytes, Available: %d bytes", fileSize, freeSpace)
		return false
	}

	// If there's enough space
	return true
}
