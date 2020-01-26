package cli

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alecthomas/kingpin"
)

// LocalWebserverVersionNumber represents the current build version.
// This should be the only one
const LocalWebserverVersionNumber = "0.1.0"

// VersionSuffix is the suffix used in the  version string.
// It will be blank for release versions.
const LocalWebserverVersionSuffix = "" // blank this when doing a release

// SHA-value of git commit
const CommitHash = ""

// Date of build
const BuildDate = ""

var (
	port             = kingpin.Flag("port", "The port to be used.").Default(randomStr(1024, 65535)).Int()
	networkInterface = kingpin.Flag("interface", "The network interface to be used. If not given a list is presented you can choose from.").Default("").String()
	silent           = kingpin.Flag("silent", "Make server silent").Bool()
	directory        = kingpin.Flag("directory", "Directory to serve files from").Default(workingDirectory()).String()
	openBrowser      = kingpin.Flag("browser", "Open in browser").Default("true").Bool()
)

type Config struct {
	Port             int
	NetworkInterface string
	Silent           bool
	WorkingDirectory string
	OpenBrowser      bool
}

// LocalWebserverVersion returns the current  version. It will include
// a suffix, typically '-DEV', if it's development version.
func LocalWebserverVersion() string {
	return fmt.Sprintf("%s%s", LocalWebserverVersionNumber, LocalWebserverVersionSuffix)
}

func (c *Config) ParseArgs() {
	kingpin.Version(fmt.Sprintf("v%s-%s", LocalWebserverVersion(), CommitHash))
	kingpin.Parse()

	c.Port = *port
	c.WorkingDirectory = *directory
	c.NetworkInterface = *networkInterface
	c.OpenBrowser = *openBrowser
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().Unix())

	return rand.Intn(max-min) + min
}

func randomStr(min, max int) string {
	return strconv.Itoa(randomInt(min, max))
}

func workingDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatal(err)
	}

	return dir
}
