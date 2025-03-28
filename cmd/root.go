/*
Copyright © 2024 Billy G. Allie <bill.allie@defiant.mug.org>
*/
package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	dbug "runtime/debug"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	cacheDir     string
	maxCacheTime float64
	installDir   string
	osCpuType    string
	extension    string
	reinstall    bool
	autoupdate   bool
	GitCommit    string = "not set"
	GitState     string = "not set"
	GitSummary   string = "not set"
	GitDate      string = "not set"
	BuildDate    string = "not set"
	Version      string = ""
)

//go:embed assets/config.ini
var iniString string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "installgo",
	Short: "installgo will automate the installation of the latest version of the GO language",
	Long: `installgo will check https://go.dev for updates for your installed version of go.
If found you can optionally install the updated version of GO.  You can also
reinstall the current version if you installed version is the latest one.`,
	Run: statusCmd.Run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "the config file to use")
	rootCmd.PersistentFlags().StringVarP(&installDir, "installdir", "d", "", "the target directory where Go is installed.")
	viper.BindPFlag("default.installdir", rootCmd.PersistentFlags().Lookup("installdir"))
	// Extract version information from the stored build information.
	bi, ok := dbug.ReadBuildInfo()
	if ok {
		Version = bi.Main.Version
		rootCmd.Version = Version
		GitDate = getBuildSettings(bi.Settings, "vcs.time")
		GitCommit = getBuildSettings(bi.Settings, "vcs.revision")
		if len(GitCommit) > 1 {
			GitSummary = fmt.Sprintf("%s-1-%s", Version, GitCommit[0:7])
		}
		GitState = "clean"
		if getBuildSettings(bi.Settings, "vcs.modified") == "true" {
			GitState = "dirty"
		}
		osCpuType = fmt.Sprintf("%s-%s", getBuildSettings(bi.Settings, "GOOS"), getBuildSettings(bi.Settings, "GOARCH"))
	}
	// Get the build date (as the modified date of the executable) if the build date
	// is not set.
	if BuildDate == "not set" {
		fpath, err := os.Executable()
		cobra.CheckErr(err)
		fpath, err = filepath.EvalSymlinks(fpath)
		cobra.CheckErr(err)
		fsys := os.DirFS(filepath.Dir(fpath))
		fInfo, err := fs.Stat(fsys, filepath.Base(fpath))
		cobra.CheckErr(err)
		BuildDate = fInfo.ModTime().UTC().Format(time.RFC3339)
	}
	cobra.OnInitialize(initConfig)
}

func getBuildSettings(settings []dbug.BuildSetting, key string) string {
	for _, v := range settings {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}

func initConfig() {
	var confPath string
	var err error
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		confPath, err = os.UserConfigDir()
		cobra.CheckErr(err)
		confPath = filepath.Join(confPath, "installgo")
		viper.AddConfigPath(confPath)
		viper.SetConfigName("config")
		viper.SetConfigType("ini")
	}
	var dirErr error
	if cacheDir == "" {
		cacheDir, dirErr = os.UserCacheDir()
		if dirErr == nil {
			cacheDir = filepath.Join(cacheDir, "installgo_cache")
		}
	}
	viper.SetEnvPrefix("igo")
	viper.SetEnvKeyReplacer(strings.NewReplacer("DEFAULT.", ""))
	viper.AutomaticEnv()                        // read in environment variables that match
	cobra.CheckErr(os.MkdirAll(confPath, 0750)) // ensure confPath exists
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// there was an error reading the config file.  If it did not exist,
		// the create a default config file with just the engineLayout in it.
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.ReadConfig(bytes.NewBuffer([]byte(iniString)))
			viper.SetDefault("confpath", confPath)
			viper.SetDefault("cachedir", cacheDir)
			cobra.CheckErr(viper.SafeWriteConfig())
		} else {
			cobra.CheckErr(err)
		}
	}
	// Use the installDir argument if it exists, else use the config.ini value
	installDir = viper.GetString(fmt.Sprintf("%s.installdir", osCpuType))
	// Get the extension to use if it exists
	extension = viper.GetString(fmt.Sprintf("%s.extension", osCpuType))
}
