package utils

import (
	"log"
	"net"
	"net/http"
	"syscall"
	"tiny-drop/internal/config"
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

func CheckDiskSpace(fileSize uint64) bool {

	const minFreeSpace = config.MinFreeSpace

	log.Println(minFreeSpace)
	statfs := syscall.Statfs_t{}
	err := syscall.Statfs("./storage", &statfs)
	if err != nil {
		log.Println("Error checking disk space:", err)
		return false
	}

	// Get the available free space (in bytes)
	freeSpace := statfs.Bavail * uint64(statfs.Bsize)
	log.Println("checking free space", freeSpace < uint64(fileSize))

	// Check if the available space after the upload will be enough
	if freeSpace < uint64(fileSize) {
		log.Printf("Not enough disk space! Required: %d bytes, Available: %d bytes", fileSize, freeSpace)
		return false
	}
	// remaining := freeSpace - fileSize
	// if remaining < minFreeSpace {
    //     log.Printf(
    //         "Not enough disk space! File requires %d bytes. Available: %d bytes. Needed remaining free: %d bytes, will remain: %d bytes",
    //         fileSize, freeSpace, minFreeSpace, remaining,
    //     )
    //     return false
    // }

	return true
}
