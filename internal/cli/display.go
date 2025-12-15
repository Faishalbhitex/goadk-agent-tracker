package cli

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorGray    = "\033[90m"
	ColorBold    = "\033[1m"
)

func Red(s string) string     { return ColorRed + s + ColorReset }
func Green(s string) string   { return ColorGreen + s + ColorReset }
func Yellow(s string) string  { return ColorYellow + s + ColorReset }
func Blue(s string) string    { return ColorBlue + s + ColorReset }
func Magenta(s string) string { return ColorMagenta + s + ColorReset }
func Cyan(s string) string    { return ColorCyan + s + ColorReset }
func Gray(s string) string    { return ColorGray + s + ColorReset }
func Bold(s string) string    { return ColorBold + s + ColorReset }
