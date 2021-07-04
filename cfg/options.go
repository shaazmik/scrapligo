package cfg

import (
	"errors"
	"reflect"
	"regexp"

	"github.com/scrapli/scrapligo/driver/network"
)

var ErrIgnoredOption = errors.New("option ignored, for different instance type")
var ErrInvalidPlatformAttr = errors.New("invalid platform attribute")

// Option function to set cfg platform options.
type Option func(interface{}) error

// base Cfg options

// WithConfigSources modify the default config sources for your platform.
func WithConfigSources(sources []string) Option {
	return func(c interface{}) error {
		cfgObj, ok := c.(*Cfg)

		if ok {
			cfgObj.ConfigSources = sources
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithOnPrepare provide an OnPrepare callable for the Cfg instance.
func WithOnPrepare(onPrepare func(*network.Driver) error) Option {
	return func(c interface{}) error {
		cfgObj, ok := c.(*Cfg)

		if ok {
			cfgObj.OnPrepare = onPrepare
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithDedicatedConnection set dedicated connection for Cfg instance.
func WithDedicatedConnection(dedicatedConnection bool) Option {
	return func(c interface{}) error {
		cfgObj, ok := c.(*Cfg)

		if ok {
			cfgObj.DedicatedConnection = dedicatedConnection
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithIgnoreVersion set ignore version for Cfg instance.
func WithIgnoreVersion(ignoreVersion bool) Option {
	return func(c interface{}) error {
		cfgObj, ok := c.(*Cfg)

		if ok {
			cfgObj.IgnoreVersion = ignoreVersion
			return nil
		}

		return ErrIgnoredOption
	}
}

// platform specific options

func setPlatformAttr(attrName string, attrValue, p interface{}) error {
	_, ok := p.(*Cfg)

	if ok {
		// this func only sets attrs for the platforms, so if we see a *Cfg we know we can bail
		return ErrIgnoredOption
	}

	v := reflect.ValueOf(p).Elem()

	fieldNames := map[string]int{}

	for i := 0; i < v.NumField(); i++ {
		fieldNames[v.Type().Field(i).Name] = i
	}

	attrIndex := -1

	for name, i := range fieldNames {
		if name == attrName {
			attrIndex = i
			break
		}
	}

	if attrIndex == -1 {
		// for some reason the platform doesnt have the specified attribute, this should *not*
		// happen... in theory :)
		return ErrInvalidPlatformAttr
	}

	fieldVal := v.Field(attrIndex)
	fieldVal.Set(reflect.ValueOf(attrValue))

	return nil
}

// WithVersionPattern set version pattern for the platform instance.
func WithVersionPattern(versionPattern *regexp.Regexp) Option {
	return func(p interface{}) error {
		err := setPlatformAttr("VersionPattern", versionPattern, p)

		if err != nil {
			if !errors.Is(err, ErrIgnoredOption) {
				return err
			}
		}

		return nil
	}
}

// LoadOptions struct for LoadConfig options.
type LoadOptions struct {
}

// LoadOption function to set options for cfg LoadConfig operations.
type LoadOption func(*LoadOptions) error
