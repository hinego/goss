package goss

import (
	"github.com/eleven26/goss/config"
	"github.com/eleven26/goss/core"
	"github.com/eleven26/goss/drivers/aliyun"
	"github.com/eleven26/goss/drivers/tencent"
	"github.com/spf13/viper"
)

// Goss is the wrapper for core.Kernel
type Goss struct {
	core.Kernel
}

// New creates a new instance based on the configuration file pointed to by configPath.
func New(configPath string) (*Goss, error) {
	err := config.ReadInConfig(configPath)
	if err != nil {
		return nil, err
	}

	goss := Goss{
		core.New(),
	}

	driver, err := defaultDriver(core.WithViper(viper.GetViper()))
	if err != nil {
		return nil, err
	}

	err = goss.RegisterDriver(driver)
	if err != nil {
		return nil, err
	}

	err = goss.UseDriver(driver)
	if err != nil {
		return nil, err
	}

	return &goss, nil
}

// NewWithViper creates a new instance based on the configuration file pointed to by viper.
func NewWithViper(viper *viper.Viper) (*Goss, error) {
	goss := Goss{
		core.New(),
	}

	driver, err := defaultDriver(core.WithViper(viper))
	if err != nil {
		return nil, err
	}

	err = goss.RegisterDriver(driver)
	if err != nil {
		return nil, err
	}

	err = goss.UseDriver(driver)
	if err != nil {
		return nil, err
	}

	return &goss, nil
}

// RegisterAliyunDriver register aliyun driver.
func (g *Goss) RegisterAliyunDriver() error {
	return g.RegisterDriver(aliyun.NewDriver())
}

// RegisterTencentDriver register tencent driver.
func (g *Goss) RegisterTencentDriver() error {
	return g.RegisterDriver(tencent.NewDriver())
}

// NewFromUserHomeConfigPath creates a new instance based on the configuration file pointed to by user home directory.
func NewFromUserHomeConfigPath() (*Goss, error) {
	path, err := config.UserHomeConfigPath()
	if err != nil {
		return nil, err
	}

	return New(path)
}
