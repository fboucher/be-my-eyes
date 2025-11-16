package version

// Version is overridden at build time via -ldflags "-X github.com/fboucher/be-my-eyes/internal/version.Version=<tag>"
// Default is "dev" when built locally without ldflags.
var Version = "dev"
