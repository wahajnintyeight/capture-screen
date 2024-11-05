package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/kbinani/screenshot"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
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

func SubscribeRedis(channelName string, redisClient redis.Client) {
	pubsub := redisClient.Subscribe(context.Background(), channelName)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		// Process message in a goroutine to handle multiple messages concurrently
		go func(message *redis.Message) {
			log.Printf("Received message from channel %s: %s\n", message.Channel, message.Payload)
			if message.Payload == "capture" {
				response, err := getSystemInfo()
				if err != nil {
					log.Printf("Error getting system info: %v", err)
					return
				}

				jsonResponse, err := json.Marshal(response)
				if err != nil {
					log.Printf("Error marshaling response: %v", err)
					return
				}

				err = redisClient.Publish(context.Background(), "screen_capture_response", jsonResponse).Err()
				if err != nil {
					log.Printf("Error publishing response: %v", err)
				} else {
					log.Println("Response sent successfully")
				}
			}
		}(msg)
	}
}

func initRedis() redis.Client {
	godotenv.Load()
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisUser := os.Getenv("REDIS_USER")
	redisPort := os.Getenv("REDIS_PORT")
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		Username: redisUser,
	})

	pingRes, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Unable to connect to Redis: ", err)

	} else {
		log.Println("Initialized Redis Connection | Ping Response: ", pingRes)
		return *client
	}
	return redis.Client{}
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
	rClient := initRedis()
	
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start Redis subscription in a goroutine
	go SubscribeRedis("capture-screen", rClient)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down gracefully...")
}
