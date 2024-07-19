package src

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

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

func StartServer(logger *log.Logger) error {
	r := gin.Default()
	r.POST("/posts", CreatePost)

	port := ":9090"

	if err := validatePort(port); err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}

	if err := checkPortAvailable(port); err != nil {
		return fmt.Errorf("port %s is unavailable: %w", port, err)
	}

	logger.Printf("Server started and running on port %s", port)
	if err := r.Run(port); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func validatePort(port string) error {
	if len(port) < 2 || port[0] != ':' {
		return fmt.Errorf("invalid port format")
	}

	portNum, err := strconv.Atoi(port[1:])
	if err != nil {
		return fmt.Errorf("invalid port number")
	}

	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("port number must be in the range from 1 to 65535")
	}

	return nil
}

func checkPortAvailable(port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("port %s unreachable: %w", port, err)
	}
	ln.Close()
	return nil
}

//Проверка на спам или недопустимый контент
//Парсинг и сохранение хэштегов
//Увеличение счетчика постов пользователя
