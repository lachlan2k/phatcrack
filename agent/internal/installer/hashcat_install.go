package installer

import (
	"archive/tar"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func installHashcat(installConf InstallConfig) {

	if _, err := os.Stat(installConf.HashcatPath); err == nil {
		log.Printf("hashcat binary was already found at %q, skipping install", installConf.HashcatPath)
		return
	}

	if installConf.Defaults {
		os.MkdirAll("/opt/phatcrack-agent/hashcat", 0700)
	}

	installDirectory := filepath.Dir(installConf.HashcatPath)

	fi, err := os.Stat(installDirectory)
	if err != nil {
		log.Fatalf("directory %q does not exist (or is not readable) for hashcat installation, determined from -hashcat-path", installDirectory)
	}

	if !fi.IsDir() {
		log.Fatalf("%q is not a directory, cannot install hashcat there, path determined from -hashcat-path", installDirectory)
	}

	u, err := url.Parse(installConf.APIEndpoint)
	if err != nil {
		log.Fatal("failed to parse API endpoint to determine asset server location: ", err)
	}

	u.Path = "/agent-assets/hashcat.tar.gz"

	tr := &http.Transport{}
	if installConf.DisableTLSVerification {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(u.String())
	if err != nil {
		log.Fatalf("failed to get hashcat.tar.gz from download server %q: %s", u.String(), err)
	}
	defer resp.Body.Close()

	err = extractTarGz(installDirectory, resp.Body)
	if err != nil {
		log.Fatal("failed to extract hashcat.tar.gz: ", err)
	}

	_, err = os.Stat(filepath.Join(installDirectory, "hashcat.bin"))
	if err != nil {
		log.Fatalf("did not find hashcat.bin within install directory %q, invalid tar.gz has been passed", installDirectory)
	}
}

// https://stackoverflow.com/questions/57639648/how-to-decompress-tar-gz-file-in-go
func extractTarGz(targetDirectory string, stream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(stream)
	if err != nil {
		return fmt.Errorf("failed to decode gzip stream: %s", err)
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("reading next tar header failed: %s", err)
		}

		safePath := filepath.Join("/", filepath.Clean(header.Name))
		// Split off the hashcat-x.x.x parent directory so we're installing hashcat.bin directly to the folder
		excludingHashcatFolder := strings.SplitN(safePath, string(os.PathSeparator), 3)
		safePath = filepath.Join(targetDirectory, excludingHashcatFolder[len(excludingHashcatFolder)-1])

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(safePath, fs.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %q: %s", safePath, err)
			}
		case tar.TypeReg:

			outFile, err := os.OpenFile(safePath, os.O_CREATE|os.O_RDWR, fs.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %q: %s", safePath, err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to write file %q: %s", safePath, err)
			}
			outFile.Close()

		default:
			return fmt.Errorf("unsupported tar entry: %q (%c)", header.Name, header.Typeflag)
		}

	}

	return nil
}
