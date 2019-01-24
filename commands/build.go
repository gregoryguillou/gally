package commands

import (
	"github.com/missena-corp/gally/project"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build your script for updated files",
	Run: func(cmd *cobra.Command, args []string) {
		handleVerboseFlag()
		if tag == "" {
			jww.INFO.Printf("no tag provided")
			return
		}
		if err := project.BuildTag(tag, rootDir); err != nil {
			jww.ERROR.Fatalf("could not build properly project: %v", err)
		}
	},
}

func init() {
	addTagFlag(buildCmd)
	rootCmd.AddCommand(buildCmd)
}
