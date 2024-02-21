package scanner

import (
	"github.com/reenigneserever/xray-knife/speedtester"
	"github.com/reenigneserever/xray-knife/speedtester/cloudflare"
	"github.com/reenigneserever/xray-knife/utils/customlog"
)

type CFScanner struct {
	// Scanner options
	rangeIP             []string
	DoDownloadSpeedTest bool
	DoUploadSpeedTest   bool

	// Engine
	Threads        uint16
	DownloadAmount uint32
	UploadAmount   uint32

	speedtest *speedtester.SpeedTester
}

type CFOption = func(c *CFScanner)

func WithCustomAmount(downloadBytes uint32, uploadBytes uint32) CFOption {
	return func(c *CFScanner) {
		c.DownloadAmount = downloadBytes
		c.UploadAmount = uploadBytes
	}
}

// WithDifferentEndpoint dpath: download path - upath: upload path
func WithDifferentEndpoint(host string, port uint16, noSSL bool, dpath string, upath string) CFOption {
	return func(c *CFScanner) {
		c.speedtest = speedtester.NewSpeedTester(nil, speedtester.WithCustomTester(host, port, noSSL, dpath, upath))
	}
}

func WithDifferentTests(DownloadTest bool, UploadTest bool) CFOption {
	return func(c *CFScanner) {
		c.DoDownloadSpeedTest = DownloadTest
		c.DoUploadSpeedTest = UploadTest
	}
}

func NewCFScanner(rangeIPs []string, threadCount uint16, opts ...CFOption) (*CFScanner, error) {
	if threadCount == 0 {
		threadCount = 1
	}
	c := &CFScanner{
		rangeIP:             rangeIPs,
		DoDownloadSpeedTest: true,
		DoUploadSpeedTest:   true,
		Threads:             threadCount,
		speedtest:           speedtester.NewSpeedTester(cloudflare.Speedtest),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *CFScanner) StartScanner() {
	customlog.Printf(customlog.Processing, "Scanner started...\n")

}
