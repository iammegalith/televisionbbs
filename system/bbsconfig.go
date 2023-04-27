package system

import (
	"log"

	"github.com/go-ini/ini"
)

// Config is a struct to hold the configuration values.
type ConfigStruct struct {
	ListenAddr      string
	Port            int
	BBSName         string
	SysopName       string
	PreLogin        bool
	ShowBulls       bool
	NewRegistration bool
	DefaultLevel    int
	ConfigPath      string
	AnsiPath        string
	AsciiPath       string
	ModulePath      string
	MenusPath		string
	DataPath        string
	FilesPath       string
	StringsFile     string
	ChannelsListenAddr	string
	ChannelsPort		int
}

var BBSConfig ConfigStruct

func ReadBBSConfig() {
	// Load the INI file.
	cfg, err := ini.Load("bbs.config")
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}

	// Read values from the [mainconfig] section and store them in the BBSConfig struct.
	BBSConfig.ListenAddr = cfg.Section("mainconfig").Key("listenaddr").String()
	BBSConfig.Port = cfg.Section("mainconfig").Key("port").MustInt(8080)
	BBSConfig.BBSName = cfg.Section("mainconfig").Key("bbsname").String()
	BBSConfig.SysopName = cfg.Section("mainconfig").Key("sysopname").String()
	BBSConfig.PreLogin = cfg.Section("mainconfig").Key("prelogin").MustBool(true)
	BBSConfig.ShowBulls = cfg.Section("mainconfig").Key("showbulls").MustBool(true)
	BBSConfig.NewRegistration = cfg.Section("mainconfig").Key("newregistration").MustBool(true)
	BBSConfig.DefaultLevel = cfg.Section("mainconfig").Key("defaultlevel").MustInt(1)
	BBSConfig.ConfigPath = cfg.Section("mainconfig").Key("configpath").String()
	BBSConfig.AnsiPath = cfg.Section("mainconfig").Key("ansipath").String()
	BBSConfig.AsciiPath = cfg.Section("mainconfig").Key("asciipath").String()
	BBSConfig.ModulePath = cfg.Section("mainconfig").Key("modulepath").String()
	BBSConfig.MenusPath = cfg.Section("mainconfig").Key("menuspath").String()
	BBSConfig.DataPath = cfg.Section("mainconfig").Key("datapath").String()
	BBSConfig.FilesPath = cfg.Section("mainconfig").Key("filespath").String()
	BBSConfig.StringsFile = cfg.Section("mainconfig").Key("stringsfile").String()
	BBSConfig.ChannelsListenAddr = cfg.Section("channels").Key("channelslistenaddr").String()
	BBSConfig.ChannelsPort = cfg.Section("channels").Key("channelsport").MustInt(8088)
}
