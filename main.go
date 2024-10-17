package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kbinani/screenshot"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"golang.org/x/net/context"
)

type Response struct {
	DeviceName   string    `json:"device_name"`
	Timestamp    string    `json:"timestamp"`
	OSName       string    `json:"os_name"`
	ImageBlob    string    `json:"image_blob"`
	CPUUsage     float64   `json:"cpu_usage"`
	MemoryUsage  string    `json:"memory_usage"`
	DiskUsage    string    `json:"disk_usage"`
	NetworkStats []NetStat `json:"network_stats"`
}

type NetStat struct {
	InterfaceName string `json:"interface_name"`
	BytesSent     uint64 `json:"bytes_sent"`
	BytesRecv     uint64 `json:"bytes_recv"`
}

var (
	deviceName string
	osName     string
)

func init() {
	var err error
	deviceName, err = os.Hostname()
	if err != nil {
		deviceName = "unknown"
	}
	osName = "Windows" // Or use runtime.GOOS for dynamic OS detection
}

func takeScreenshot() (string, error) {
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 70})
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func getSystemInfo() (Response, error) {
	cpuUsage, _ := cpu.Percent(0, false)
	v, _ := mem.VirtualMemory()
	d, _ := disk.Usage("/")
	netStats, _ := net.IOCounters(true)

	var networkStats []NetStat
	for _, stat := range netStats {
		networkStats = append(networkStats, NetStat{
			InterfaceName: stat.Name,
			BytesSent:     stat.BytesSent,
			BytesRecv:     stat.BytesRecv,
		})
	}

	screenshotBlob, err := takeScreenshot()
	if err != nil {
		return Response{}, err
	}

	return Response{
		DeviceName:   deviceName,
		Timestamp:    time.Now().Format(time.RFC3339),
		OSName:       osName,
		ImageBlob:    screenshotBlob,
		CPUUsage:     cpuUsage[0],
		MemoryUsage:  fmt.Sprintf("%v / %v", formatBytes(v.Used), formatBytes(v.Total)),
		DiskUsage:    fmt.Sprintf("%v / %v", formatBytes(d.Used), formatBytes(d.Total)),
		NetworkStats: networkStats,
	}, nil
}

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

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Replace with your Redis server address
	})

	pubsub := rdb.Subscribe(ctx, "screen_capture_request")
	defer pubsub.Close()

	log.Println("Listening for screen capture requests...")

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			continue
		}

		if msg.Payload == "capture" {
			log.Println("Capture request received")

			response, err := getSystemInfo()
			if err != nil {
				log.Printf("Error getting system info: %v", err)
				continue
			}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
				continue
			}

			err = rdb.Publish(ctx, "screen_capture_response", jsonResponse).Err()
			if err != nil {
				log.Printf("Error publishing response: %v", err)
			} else {
				log.Println("Response sent successfully")
			}
		}
	}
}
