package registry

import (
	corev1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("registry")

type Type string

type Config struct {
	RegistryType Type
	Args         []string
	VolumeMounts []corev1.VolumeMount
	Volumes      []corev1.Volume
}

// registry details
var registryType Type
var repositoryName string
var imageName string

var registryConfigs = map[Type]func(repoName string, imgName string) *Config{}

func SetRegistry(regType Type, repoName string, imgName string) {
	registryType = regType
	repositoryName = repoName
	imageName = imgName
}

func GetConfig() *Config {
	return registryConfigs[registryType](repositoryName, imageName)
}

func addRegistryConfig(regType Type, configFunc func(repoName string, imgName string) *Config) {
	if registryConfigs[regType] != nil {
		log.Error(nil, "Duplicate registry types", "type", regType)
	} else {
		registryConfigs[regType] = configFunc
	}
}
