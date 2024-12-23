package logging

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

func GetLogger(l int) (*slog.Logger, error) {
	level := slog.Level(l)
	// fmt.Println("log level: ", level)
	opts := slog.HandlerOptions{
		// AddSource: true,
		Level: level,
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	configPath, err := filepath.Abs(dir + "/logs/rain.log")
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(
		configPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0777,
	)
	if err != nil {
		return nil, err
	}
	// defer file.Close() // defer swithes off logger

	textHandler := slog.NewTextHandler(file, &opts)
	logger := slog.New(textHandler)

	return logger, nil
}
