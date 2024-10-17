package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image/png"
	"log"
	"net/http"
	"os"
	"time"
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type Response struct {
	DeviceName    string   `json:"device_name"`
	Timestamp     string   `json:"timestamp"`
	OSName        string   `json:"os_name"`
	ImageBlob     string   `json:"image_blob"`
	CPUUsage      float64  `json:"cpu_usage"`
	MemoryUsage   string   `json:"memory_usage"`
	DiskUsage     string   `json:"disk_usage"`
	NetworkStats  []NetStat `json:"network_stats"`
}

type NetStat struct {
	InterfaceName string `json:"interface_name"`
	BytesSent     uint64 `json:"bytes_sent"`
	BytesRecv     uint64 `json:"bytes_recv"`
}

func takeScreenshot() (string, error) {
	log.Println("Taking screenshot...")
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		log.Printf("Error capturing screenshot: %v", err)
		return "", err
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		log.Printf("Error encoding screenshot: %v", err)
		return "", err
	}

	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	log.Println("Screenshot taken and encoded successfully")
	return base64Image, nil
}

func getDeviceName() string {
	name, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting hostname: %v", err)
		return "unknown"
	}
	log.Printf("Device name: %s", name)
	return name
}

func getCPUUsage() (float64, error) {
	log.Println("Getting CPU usage...")
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
		return 0, err
	}
	log.Printf("CPU usage: %.2f%%", percentages[0])
	return percentages[0], nil
}

func getMemoryUsage() (string, error) {
	log.Println("Getting memory usage...")
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting memory usage: %v", err)
		return "", err
	}
	usage := formatBytes(v.Used) + " / " + formatBytes(v.Total)
	log.Printf("Memory usage: %s", usage)
	return usage, nil
}

func getDiskUsage() (string, error) {
	log.Println("Getting disk usage...")
	usage, err := disk.Usage("/")
	if err != nil {
		log.Printf("Error getting disk usage: %v", err)
		return "", err
	}
	diskUsage := formatBytes(usage.Used) + " / " + formatBytes(usage.Total)
	log.Printf("Disk usage: %s", diskUsage)
	return diskUsage, nil
}

func getNetworkStats() ([]NetStat, error) {
	log.Println("Getting network stats...")
	stats, err := net.IOCounters(true)
	if err != nil {
		log.Printf("Error getting network stats: %v", err)
		return nil, err
	}

	var netStats []NetStat
	for _, stat := range stats {
		netStats = append(netStats, NetStat{
			InterfaceName: stat.Name,
			BytesSent:     stat.BytesSent,
			BytesRecv:     stat.BytesRecv,
		})
		log.Printf("Network interface %s: Sent %d bytes, Received %d bytes", stat.Name, stat.BytesSent, stat.BytesRecv)
	}

	return netStats, nil
}

// Helper function to format bytes
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func handleScreenshot(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling screenshot request...")
	deviceName := getDeviceName()
	currentTime := time.Now().Format(time.RFC3339)
	osName := "Windows" // or dynamically retrieve using runtime.GOOS
	log.Printf("OS: %s", osName)

	screenshotBlob, err := takeScreenshot()
	if err != nil {
		log.Printf("Failed to capture screenshot: %v", err)
		http.Error(w, "Failed to capture screenshot", http.StatusInternalServerError)
		return
	}

	cpuUsage, err := getCPUUsage()
	if err != nil {
		log.Printf("Failed to fetch CPU usage: %v", err)
		http.Error(w, "Failed to fetch CPU usage", http.StatusInternalServerError)
		return
	}

	memoryUsage, err := getMemoryUsage()
	if err != nil {
		log.Printf("Failed to fetch memory usage: %v", err)
		http.Error(w, "Failed to fetch memory usage", http.StatusInternalServerError)
		return
	}

	diskUsage, err := getDiskUsage()
	if err != nil {
		log.Printf("Failed to fetch disk usage: %v", err)
		http.Error(w, "Failed to fetch disk usage", http.StatusInternalServerError)
		return
	}

	networkStats, err := getNetworkStats()
	if err != nil {
		log.Printf("Failed to fetch network stats: %v", err)
		http.Error(w, "Failed to fetch network stats", http.StatusInternalServerError)
		return
	}

	response := Response{
		DeviceName:   deviceName,
		Timestamp:    currentTime,
		OSName:       osName,
		ImageBlob:    screenshotBlob,
		CPUUsage:     cpuUsage,
		MemoryUsage:  memoryUsage,
		DiskUsage:    diskUsage,
		NetworkStats: networkStats,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	log.Println("Screenshot request handled successfully")
}

func main() {
	http.HandleFunc("/screenshot", handleScreenshot)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
