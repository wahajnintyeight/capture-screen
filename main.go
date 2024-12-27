package main

import (
	"bytes"
	"context"

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

	"capture-screen/internal/aws"
	"capture-screen/internal/config"

	"github.com/go-redis/redis/v8"
	"github.com/gosimple/slug"
	"github.com/joho/godotenv"
	"github.com/kbinani/screenshot"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"google.golang.org/grpc"
)

type Response struct {
	DeviceName  string `json:"deviceName"`
	Timestamp   string `json:"timestamp"`
	OSName      string `json:"osName"`
	MemoryUsage string `json:"memoryUsage"`
	DiskUsage   string `json:"diskUsage"`
	LastImage   string `json:"lastImage"`
}

var (
	deviceName string
	osName     string
	s3Service  *aws.S3Service
)

type MessageType int

const (
	CAPTURE_SCREEN MessageType = iota
	PING_DEVICE
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
			case "capture-screen-" + getSlugDeviceName():
				response, err := getSystemInfo("capture-screen")
				if err != nil {
					log.Printf("Error getting system info: %v", err)
					return
				}

				sendGRPCCall(response, int32(CAPTURE_SCREEN))
				break
			case "scan-devices":
				log.Println("Scanning devices")
				deviceName := getDeviceName()
				log.Println("Device Name:", deviceName)
				sendHTTPCall(map[string]interface{}{"deviceName": deviceName}, "/return-device-name")
				break
			case "ping-device-" + getSlugDeviceName():
				log.Println("Pinging device")
				response, err := getSystemInfo("ping-device")
				if err != nil {
					log.Printf("Error getting system info: %v", err)
					return
				}
				sendGRPCCall(response, int32(PING_DEVICE))
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

func takeScreenshot() ([]byte, error) {

	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, fmt.Errorf("capture error: %v", err)
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 70})
	if err != nil {
		return nil, fmt.Errorf("jpeg encode error: %v", err)
	}

	return buf.Bytes(), nil

}

func getDeviceName() string {
	return deviceName
}

func getSlugDeviceName() string {
	return slug.Make(deviceName)
}

func getSystemInfo(eventType string) (Response, error) {

	v, _ := mem.VirtualMemory()
	d, _ := disk.Usage("/")
	timestamp := time.Now().Format(time.RFC3339)

	var imageBytes []byte
	var err error
	if eventType == "capture-screen" {
		imageBytes, err = takeScreenshot()
		if err != nil {
			return Response{}, err
		}
	} else {
		return Response{
			DeviceName:  deviceName,
			Timestamp:   timestamp,
			OSName:      osName,
			MemoryUsage: fmt.Sprintf("%v / %v", formatBytes(v.Used), formatBytes(v.Total)),
			DiskUsage:   fmt.Sprintf("%v / %v", formatBytes(d.Used), formatBytes(d.Total)),
			LastImage:   "",
		}, nil
	}

	secureURL, err := s3Service.UploadImage(context.Background(), imageBytes, getDeviceName())

	if err != nil {
		log.Println("Error while uploading:", err)
		return Response{
			DeviceName:  "",
			Timestamp:   "",
			OSName:      "",
			MemoryUsage: fmt.Sprintf("%v / %v", formatBytes(v.Used), formatBytes(v.Total)),
			DiskUsage:   fmt.Sprintf("%v / %v", formatBytes(d.Used), formatBytes(d.Total)),
			LastImage:   "",
		}, nil
	} else {
		return Response{
			DeviceName:  deviceName,
			Timestamp:   timestamp,
			OSName:      osName,
			MemoryUsage: fmt.Sprintf("%v / %v", formatBytes(v.Used), formatBytes(v.Total)),
			DiskUsage:   fmt.Sprintf("%v / %v", formatBytes(d.Used), formatBytes(d.Total)),
			LastImage:   secureURL,
		}, nil
	}

}

// func uploadImageToCloudinary(ctx context.Context, imageBytes []byte, folder string) (string, error) {

// 	// Initialize Cloudinary
// 	cld, err := cloudinary.NewFromURL("cloudinary://" + os.Getenv("CLOUDINARY_API_KEY") + ":" + os.Getenv("CLOUDINARY_API_SECRET") + "@" + os.Getenv("CLOUDINARY_CLOUD_NAME"))
// 	if err != nil {
// 		return "", fmt.Errorf("error initializing Cloudinary: %v", err)
// 	}

// 	timestamp := time.Now().Format(time.RFC3339)
// 	publicID := deviceName + "_" + timestamp
// 	overwrite := true
// 	uploadParams := uploader.UploadParams{
// 		PublicID:     publicID,
// 		Folder:       folder,
// 		ResourceType: "raw",
// 		Overwrite:    &overwrite,
// 	}

// 	reader := bytes.NewReader(imageBytes)
// 	uploadResult, err := cld.Upload.Upload(ctx, reader, uploadParams)
// 	if err != nil {
// 		return "", fmt.Errorf("error uploading to Cloudinary: %v", err)
// 	}

// 	return uploadResult.SecureURL, nil
// }

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

func sendGRPCCall(response Response, messageType int32) {
	grpcURL := os.Getenv("GRPC_SERVER_URL")

	conn, err := grpc.NewClient(grpcURL, grpc.WithInsecure())
	if err != nil {
		log.Println("Error while connecting", err)
	}
	log.Println("Connected to GRPC Server", conn)

	grpcClient := pb.NewScreenCaptureServiceClient(conn)
	log.Println("GRPC Client", grpcClient)
	res, e := grpcClient.SendCapture(context.Background(),
		&pb.ScreenCaptureRequest{
			DeviceName:  response.DeviceName,
			TimesTamp:   response.Timestamp,
			OsName:      response.OSName,
			MemoryUsage: response.MemoryUsage,
			DiskUsage:   response.DiskUsage,
			LastImage:   response.LastImage,
			MessageType: messageType,
		})
	if e != nil {
		log.Println("Error while calling the method!", e)
	} else {
		log.Println("Res:", res)
	}
}

func main() {
	rClient := initRedis()

	config.LoadEnv()

	var err error
	s3Service, err = aws.NewS3Service(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize S3 service: %v", err)
	}
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	deviceName := getDeviceName()
	slugifiedDeviceName := slug.Make(deviceName)
	// Start Redis subscription in a goroutine
	go SubscribeRedis("capture-screen-"+slugifiedDeviceName, rClient)
	go SubscribeRedis("scan-devices", rClient)
	go SubscribeRedis("ping-device-"+slugifiedDeviceName, rClient)
	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down gracefully...")
}
