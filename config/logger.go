package config

import (
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"time"
)

/*
InitLogger Initialize the slog multi-handler for JSON file and STDOUT text logging,
then store the file handler in viper so that we can close it after
*/
func InitLogger() error {
	filename := BuildLogFilename()

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	viper.Set("log.file", filename)
	viper.Set("log.fileObject", file)

	handler := slogmulti.Fanout(
		slog.NewJSONHandler(file, nil),
		slog.NewTextHandler(os.Stdout, nil))

	slog.SetDefault(slog.New(handler))

	return nil
}

/*
BuildLogFilename Build a new filename based on a timestamp
*/
func BuildLogFilename() string {
	timestamp := time.Now().Format(time.RFC3339Nano)

	return viper.GetString("log.path") + "/arcane-" + timestamp + ".json"
}
