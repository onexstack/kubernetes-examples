package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/superproj/k8sdemo/featuregates/feature"
)

func main() {
	// Create a new FlagSet for managing command-line flags
	fs := pflag.NewFlagSet("feature", pflag.ExitOnError)

	// Set the usage function to provide a custom help message
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fs.PrintDefaults()
	}

	// Define a boolean flag for displaying help
	help := fs.BoolP("help", "h", false, "Show this help message.")

	// Add the feature gates to the flag set
	feature.DefaultMutableFeatureGate.AddFlag(fs)

	// Parse the command-line flags
	fs.Parse(os.Args[1:])

	// Display help message if the help flag is set
	if *help {
		fs.Usage()
		return
	}

	// Check if the MyNewFeature feature gate is enabled
	if feature.DefaultFeatureGate.Enabled(feature.MyNewFeature) {
		// Logic when the new feature is enabled
		fmt.Println("Feature Gates: MyNewFeature is opened")
	} else {
		// Logic when the new feature is disabled
		fmt.Println("Feature Gates: MyNewFeature is closed")
	}
}
