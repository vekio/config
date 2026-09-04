// Package config manages typed YAML and JSON configuration files.
//
// Configuration values implement [Validatable] and are validated after being
// loaded and before being written. Writes are atomic and use private file
// permissions.
//
// New configuration files use the user's configuration directory. A path may
// be selected through APPNAME_CONFIG_FILE, whose prefix is derived from the
// application name, or explicitly with [ConfigFile.SetPath]. An explicit
// SetPath call takes precedence over the environment value captured when the
// ConfigFile is created. Empty environment values are ignored.
//
// File operations currently support Linux.
package config
