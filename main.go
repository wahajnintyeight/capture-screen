package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
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
	"github.com/gosimple/slug"
)

type Response struct {
	DeviceName  string `json:"device_name"`
	Timestamp   string `json:"timestamp"`
	OSName      string `json:"os_name"`
	ImageBlob   []byte `json:"image_blob"`
	MemoryUsage string `json:"memory_usage"`
	DiskUsage   string `json:"disk_usage"`
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
			switch message.Payload {
			case "capture-screen-"+getSlugDeviceName():
				response, err := getSystemInfo()
				if err != nil {
					log.Printf("Error getting system info: %v", err)
					return
				}

				sendGRPCCall(response)
				break
			case "scan-devices":
				log.Println("Scanning devices")
				deviceName := getDeviceName()
				log.Println("Device Name:", deviceName)
				sendHTTPCall(map[string]interface{}{"deviceName": deviceName}, "/return-device-name")
				break
			case "ping-"+getSlugDeviceName():
				log.Println("Pinging device")
				sendHTTPCall(map[string]interface{}{"status": "OK"}, "/respond-ping")
				break
			default:
				log.Println("Unknown command", message.Payload)
			}
		}(msg)
	}
}
func sendHTTPCall(data interface{}, endpoint string) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}
	apiUrl := os.Getenv("API_URL")

	jsonPayload := bytes.NewReader(jsonData)
	log.Println("Sending HTTP Call to ", apiUrl+endpoint)

	req, err := http.NewRequest(http.MethodPost, apiUrl+endpoint, jsonPayload)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-200 response code: %d", resp.StatusCode)
		return
	}
	log.Println("HTTP Call Response: ", resp)
	return
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

func getDeviceName() string {
	return deviceName
}

func getSlugDeviceName() string {
	return slug.Make(deviceName)
}

func getSystemInfo() (Response, error) {

	v, _ := mem.VirtualMemory()
	d, _ := disk.Usage("/")

	screenshotBlob, err := takeScreenshot()
	if err != nil {
		return Response{}, err
	}

	return Response{
		DeviceName: deviceName,
		Timestamp:  time.Now().Format(time.RFC3339),
		OSName:     osName,
		ImageBlob:  []byte(screenshotBlob),

		MemoryUsage: fmt.Sprintf("%v / %v", formatBytes(v.Used), formatBytes(v.Total)),
		DiskUsage:   fmt.Sprintf("%v / %v", formatBytes(d.Used), formatBytes(d.Total)),
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
			DeviceName:  response.DeviceName,
			Timestamp:   response.Timestamp,
			OsName:      response.OSName,
			ImageData:   response.ImageBlob,
			MemoryUsage: response.MemoryUsage,
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

	deviceName := getDeviceName()
	// Start Redis subscription in a goroutine
	go SubscribeRedis("capture-screen-"+slug.Make(deviceName), rClient)
	go SubscribeRedis("scan-devices", rClient)
	go SubscribeRedis("ping-"+slug.Make(deviceName),rClient)
	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down gracefully...")
}
