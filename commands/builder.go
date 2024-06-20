package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/schollz/progressbar/v3"
)

type Build struct {
	Config      string   `help:"newline delimited file of URLS to download and process" required:""                                               type:"existingfile"`
	DB          string   `help:"db filename to import data to"                          required:""`
	AllowedTags []string `default:"*"                                                   help:"a list of allowed tags, all other will be filtered"`
}

var ErrSourceNotAvailable = errors.New("source unavailable")

func (b *Build) Run(stdout io.Writer) error {
	config, err := os.Open(b.Config)
	if err != nil {
		return fmt.Errorf("could not open config file: %w", err)
	}

	buildPath := filepath.Dir(b.DB)

	err = os.MkdirAll(buildPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create build path: %w", err)
	}

	matches, err := filepath.Glob(b.DB + "*")
	if err != nil {
		return fmt.Errorf("could not glob previous db file: %w", err)
	}

	for _, match := range matches {
		err := os.Remove(match)
		if err != nil {
			return fmt.Errorf("could not remove file: %w", err)
		}
	}

	client := req.C().SetOutputDirectory(buildPath)
	scanner := bufio.NewScanner(config)

	for scanner.Scan() {
		url := scanner.Text()

		filename := filepath.Base(url)
		slog.Info("processing", "url", url, "filename", filename)
		downloadFilename := filepath.Join(buildPath, filename)

		if _, err := os.Stat(downloadFilename); errors.Is(err, os.ErrNotExist) {
			slog.Info("download", "url", url, "downloadFilename", downloadFilename)

			bar := progressbar.DefaultBytes(
				1,
				"downloading",
			)

			response, err := client.R().
				SetOutputFile(filename).
				SetRetryCount(3).
				SetDownloadCallback(func(info req.DownloadInfo) {
					if info.Response.Response != nil {
						bar.ChangeMax64(info.Response.ContentLength)
						_ = bar.Set64(info.DownloadedSize)
					}
				}).
				Get(url)
			if err != nil {
				return fmt.Errorf("could not download %q: %w", url, err)
			}

			if response.StatusCode == http.StatusNotFound {
				return fmt.Errorf("could find url %q to download: %w", url, ErrSourceNotAvailable)
			}

			_ = bar.Finish()
		}

		// prefix remove the `latest` and extension
		prefix := strings.ReplaceAll(filename, ".osm.pbf", "")
		prefix = strings.ReplaceAll(prefix, ".pbf", "")
		prefix = strings.ReplaceAll(prefix, "-latest", "")

		slog.Info("converting", "downloadFilename", downloadFilename)
		command := &Convert{
			OSM:         downloadFilename,
			DB:          b.DB,
			Prefix:      prefix,
			AllowedTags: b.AllowedTags,
		}

		err := command.Run(stdout)
		if err != nil {
			return fmt.Errorf("could not convert %q: %w", filename, err)
		}
	}

	return nil
}
