package version

const (
	Service = "kubemv"
	Version = "0.1.0"
	Phase   = 0
)

// Capabilities are the contracts this process actually serves.
var Capabilities = []string{
	"auth.session",
	"cluster.registry",
	"plugin.registry",
}
