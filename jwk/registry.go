package jwk

import (
	"github.com/ThalesIgnite/crypto11"
	"github.com/ory/hydra/x"
)

type InternalRegistry interface {
	x.RegistryWriter
	x.RegistryLogger
	Registry
}

type Registry interface {
	KeyManager() Manager
	KeyGenerators() map[string]KeyGenerator
	KeyCipher() *AEAD
	HardwareSecurityModule() *crypto11.Context
}
