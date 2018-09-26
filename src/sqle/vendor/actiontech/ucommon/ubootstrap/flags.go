package ubootstrap

import (
	"actiontech/ucommon/conf"
	"actiontech/ucommon/log"
	"actiontech/ucommon/os"
	"fmt"
	"github.com/spf13/pflag"
	dry "github.com/ungerik/go-dry"
	"sync"
)

var flagsMu sync.Mutex

func LoadFlags(currentFlags *pflag.FlagSet) error {
	flagsMu.Lock()
	defer flagsMu.Unlock()

	return loadFlags(currentFlags)
}

func loadFlags(currentFlags *pflag.FlagSet) error {
	if !os.IsFileExist("flags") {
		return nil
	}
	config, err := conf.ReadConfigFile("flags")
	if nil != err {
		return fmt.Errorf("invalid file flags: %v", err)
	}
	var outError error
	currentFlags.VisitAll(func(flag *pflag.Flag) {
		if flag.Changed {
			return
		}
		val, _ := config.GetString("default", flag.Name)
		if "" == val {
			return
		}
		if err := flag.Value.Set(val); nil != err {
			outError = err
			return
		}
		flag.Changed = true
	})
	return outError
}

func PersistFlags(flags *pflag.FlagSet, excepts []string) error {
	flagsMu.Lock()
	defer flagsMu.Unlock()

	return persistFlags(flags, excepts)
}

func persistFlags(flags *pflag.FlagSet, excepts []string) error {
	config := conf.NewConfigFile()
	flags.VisitAll(func(flag *pflag.Flag) {
		if !flag.Changed {
			return
		}
		if nil != excepts && dry.StringListContains(excepts, flag.Name) {
			return
		}
		config.AddOption("default", flag.Name, flag.Value.String())
	})
	if err := config.WriteConfigFile("flags", 0640, "", []string{}); nil != err {
		return err
	}
	return nil
}

func LoadAndPersistFlags(flags *pflag.FlagSet) error {
	flagsMu.Lock()
	defer flagsMu.Unlock()

	if err := loadFlags(flags); nil != err {
		return err
	}
	if err := persistFlags(flags, nil); nil != err {
		return err
	}
	return nil
}

func LoadAndPersistFlagsExcept(flags *pflag.FlagSet, excepts []string) error {
	flagsMu.Lock()
	defer flagsMu.Unlock()

	if err := loadFlags(flags); nil != err {
		return err
	}
	if err := persistFlags(flags, excepts); nil != err {
		return err
	}
	return nil
}

func UpdateFlags(stage *log.Stage, options, vals []string) error {
	flagsMu.Lock()
	defer flagsMu.Unlock()

	var config *conf.ConfigFile
	if os.IsFileExist("flags") {
		if c, err := conf.ReadConfigFile("flags"); nil != err {
			return fmt.Errorf("invalid file flags: %v", err)
		} else {
			config = c
		}
	} else {
		config = conf.NewConfigFile()
	}
	for idx, _ := range options {
		config.AddOption("default", options[idx], vals[idx])
	}
	return os.SafeWriteConfigFile(stage, "flags", config)
}

func PrintFlags(flags *pflag.FlagSet) string {
	ret := ""
	flags.VisitAll(func(flag *pflag.Flag) {
		if "" != ret {
			ret = ret + ", "
		}
		ret = ret + fmt.Sprintf("%v=%v", flag.Name, flag.Value)
	})
	return ret
}
