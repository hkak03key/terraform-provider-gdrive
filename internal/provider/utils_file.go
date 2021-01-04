package provider

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
)

var definedMimeTypes = map[string]string{
	"folder":   "application/vnd.google-apps.folder",
	"shortcut": "application/vnd.google-apps.shortcut",
}

func getMimeType(file *os.File) (string, error) {
	// based on https://gist.github.com/hkak03key/06b25a3f4f0bbd8d23d361fa8eb0dff8
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	file.Seek(0, 0)

	mimeType := http.DetectContentType(buffer[:n])
	return mimeType, nil
}

func getFileMd5Checksum(source string) (string, error) {
	f, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
