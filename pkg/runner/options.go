package runner

import (
	"flag"
	"os"

	"github.com/projectdiscovery/gologger"
)

// Options contains the configuration options for tuning
// the active dns resolving process.
type Options struct {
	Domain         string // Domain is the domain to find subdomains
	SubdomainsList string // SubdomainsList is the file containing list of hosts to resolve
	ResolversFile  string // ResolversFile is the file containing resolvers to use for enumeration
	Wordlist       string // Wordlist is a wordlist to use for enumeration
	MassdnsPath    string // MassdnsPath contains the path to massdns binary
	Output         string // Output is the file to write found subdomains to.
	Silent         bool   // Silent suppresses any extra text and only writes found host:port to screen
	Version        bool   // Version specifies if we should just show version and exit
	Retries        int    // Retries is the number of retries for dns enumeration
	Verbose        bool   // Verbose flag indicates whether to show verbose output or not
	NoColor        bool   // No-Color disables the colored output
	Threads        int    // Thread controls the number of parallel host to enumerate

	Stdin bool // Stdin specifies whether stdin input was given to the process

}

// ParseOptions parses the command line flags provided by a user
func ParseOptions() *Options {
	options := &Options{}

	flag.StringVar(&options.Domain, "d", "", "Domain to find or resolve subdomains for")
	flag.StringVar(&options.SubdomainsList, "list", "", "File containing list of subdomains to resolve")
	flag.StringVar(&options.ResolversFile, "r", "", "File containing list of resolvers for enumeration")
	flag.StringVar(&options.Wordlist, "w", "", "File containing words to bruteforce for domain")
	flag.StringVar(&options.MassdnsPath, "massdns", "", "Path to the massdns binary")
	flag.StringVar(&options.Output, "o", "", "File to write output to (optional)")
	flag.BoolVar(&options.Silent, "silent", false, "Show only host:ports in output")
	flag.BoolVar(&options.Version, "version", false, "Show version of naabu")
	flag.IntVar(&options.Retries, "retries", 1, "Number of retries for the port scan probe")
	flag.BoolVar(&options.Verbose, "v", false, "Show Verbose output")
	flag.BoolVar(&options.NoColor, "nC", false, "Don't Use colors in output")
	flag.IntVar(&options.Threads, "t", 10, "Number of concurrent goroutines for resolving")

	flag.Parse()

	// Check if stdin pipe was given
	options.Stdin = hasStdin()

	// Read the inputs and configure the logging
	options.configureOutput()

	// Show the user the banner
	showBanner()

	if options.Version {
		gologger.Infof("Current Version: %s\n", Version)
		os.Exit(0)
	}
	// Validate the options passed by the user and if any
	// invalid options have been used, exit.
	err := options.validateOptions()
	if err != nil {
		gologger.Fatalf("Program exiting: %s\n", err)
	}

	return options
}

func hasStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return false
	}
	return true
}