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

	pb "capture-screen/src/output"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/kbinani/screenshot"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"google.golang.org/grpc"
)

type Response struct {
	DeviceName   string    `json:"device_name"`
	Timestamp    string    `json:"timestamp"`
	OSName       string    `json:"os_name"`
	ImageBlob    []byte    `json:"image_blob"`
	MemoryUsage  string    `json:"memory_usage"`
	DiskUsage    string    `json:"disk_usage"`
 
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

				sendGRPCCall(response)
			}
		}(msg)
	}
}

func JSONStringToStruct(data interface{}, target interface{}) error {
	// Convert data to JSON string
	jsonData, ok := data.(string)
	if !ok {
		return fmt.Errorf("data is not a JSON string")
	}

	// Unmarshal JSON string into target struct
	err := json.Unmarshal([]byte(jsonData), target)
	if err != nil {
		log.Println("Failed to unmarshal JSON: ", err)
		return err
	}
	return nil
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
	 
	v, _ := mem.VirtualMemory()
	d, _ := disk.Usage("/")
 

	screenshotBlob, err := takeScreenshot()
	if err != nil {
		return Response{}, err
	}

	return Response{
		DeviceName:   deviceName,
		Timestamp:    time.Now().Format(time.RFC3339),
		OSName:       osName,
		ImageBlob:    []byte(screenshotBlob),
	 
		MemoryUsage:  fmt.Sprintf("%v / %v", formatBytes(v.Used), formatBytes(v.Total)),
		DiskUsage:    fmt.Sprintf("%v / %v", formatBytes(d.Used), formatBytes(d.Total)),
	 
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

func sendGRPCCall(response Response) {

	conn, err := grpc.NewClient("127.0.0.1:8880", grpc.WithInsecure())
	if err != nil {
		log.Println("Error while connecting", err)
	}
	log.Println("Connected to GRPC Server", conn)

	grpcClient := pb.NewScreenCaptureServiceClient(conn)
	log.Println("GRPC Client", grpcClient)
	res, e := grpcClient.SendCapture(context.Background(),
		&pb.ScreenCaptureRequest{
			DeviceName:   response.DeviceName,
			Timestamp:    response.Timestamp,
			OsName:       response.OSName,
			ImageData:    response.ImageBlob,			 
			MemoryUsage:  response.MemoryUsage,
		})
	if e != nil {
		log.Println("Error while calling the method!", e)
	} else {
		log.Println("Res:", res)
	}
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
