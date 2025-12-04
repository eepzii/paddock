package browser

type OperatingSystem string

const (
	LINUX   OperatingSystem = "linux"
	WINDOWS OperatingSystem = "windows"
	DARWIN  OperatingSystem = "darwin"
)

type Brand string

const (
	UNKNOWN Brand = ""
	CHROME  Brand = "chrome"
	EDGE    Brand = "edge"
)

const (
	LOOPBACK_ADDRESS = "ws://127.0.0.1"
	LOOPBACK_PORT    = 9222
)
