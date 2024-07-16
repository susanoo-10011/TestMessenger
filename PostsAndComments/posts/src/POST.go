package src

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Client *mongo.Client

type Post struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"Title"`
	Content   string             `bson:"Content"`
	Image     *MediaFile         `bson:"Image,omitempty"`
	Video     *MediaFile         `bson:"Video,omitempty"`
	Gif       *MediaFile         `bson:"Gif,omitempty"`
	CreatedAt time.Time          `bson:"Created_at"`
}

type MediaFile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	URL       string             `bson:"url"`
	FileType  string             `bson:"file_type"`
	CreatedAt time.Time          `bson:"created_at"`
}

func InitializeLogger(w io.Writer) *log.Logger {
	var logOutput io.Writer = os.Stdout

	if _, err := os.Stat("server.log"); os.IsNotExist(err) {
		logFile, err := os.Create("server.log")
		if err != nil {

			log.Printf("The log file could not be created: %v\n", err)
		} else {
			defer logFile.Close()
			logOutput = logFile
		}
	}

	return log.New(io.MultiWriter(os.Stdout, logOutput), "", log.LstdFlags)
}

func checkPortAvailable(port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("port %s unreachable: %w", port, err)
	}
	ln.Close()
	return nil
}

func StartServer(logger *log.Logger) {
	r := gin.Default()
	r.POST("/posts", CreatePost)

	port := ":9090"
	portNum, err := strconv.Atoi(port[1:])
	if err != nil {
		log.Fatalf("Invalid port number")
	}
	if portNum < 1 || portNum > 65535 {
		log.Fatalf("The port number must be in the range from 1 to 65535")
	}
	if err := checkPortAvailable(port); err != nil {
		logger.Printf("Server started and running on port %s", port)
	}
}
