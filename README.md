# config

`config` gestiona ficheros de configuración YAML o JSON usando el tipo genérico
de la aplicación. El tipo debe implementar `Validate() error`; la librería valida
los datos antes de escribirlos y después de leerlos.

```go
package main

import (
	"fmt"

	"github.com/vekio/config"
)

type Config struct {
	Address string `yaml:"address" json:"address"`
}

func (c Config) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("address is required")
	}
	return nil
}

func main() {
	file, err := config.NewDefaultConfigFile[Config]("myapp")
	if err != nil {
		panic(err)
	}
	cfg, err := file.LoadOrCreate(Config{Address: "127.0.0.1:8080"})
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.Address)
}
```

`NewDefaultConfigFile` guarda el fichero en el directorio de configuración del
usuario, dentro de `<app>/config.yml`. Los constructores explícitos reciben el
directorio base, el nombre de aplicación y el nombre del fichero:

```go
yamlFile, err := config.NewYAMLConfigFile[Config]("/etc", "myapp", "settings.yml")
jsonFile, err := config.NewJSONConfigFile[Config]("/etc", "myapp", "settings.json")
```

`file.Path()` devuelve la ruta efectiva del fichero. `DefaultDataDir(appName)` devuelve
`<XDG_DATA_HOME>/<appName>` cuando la variable está definida y, en caso
contrario, `~/.local/share/<appName>`.

## Comando CLI

El paquete ofrece un comando reutilizable para aplicaciones construidas con
`urfave/cli/v3`:

```go
app := &cli.Command{
	Name:     "myapp",
	Flags:    []cli.Flag{config.NewConfigFlag(configFile)},
	Commands: []*cli.Command{config.NewConfigCommand(configFile, defaultConfig)},
}
```

`myapp config` muestra la ayuda. Los subcomandos disponibles son `show`,
`path`, `validate` e `init`. `init` crea el fichero exclusivamente y devuelve
un error si ya existe; `show` es explícito para evitar mostrar accidentalmente
configuraciones que puedan contener secretos.

`config init --force` (o `-f`) reemplaza atómicamente un fichero existente
usando los valores predeterminados.

Si la aplicación carga su configuración desde un hook `Before`, conviene
asociarlo a los comandos que la consumen y no al comando raíz. De esta forma
`config init` puede crear un fichero que todavía no existe, como muestra el
ejemplo CLI.

El flag global `--config` permite reemplazar la ruta completa en tiempo de
ejecución y se aplica también a los subcomandos:

```text
myapp config show --config ./development.yml
```

Hay ejemplos completos para [LoadOrCreate](example/load_or_create/main.go) y
para la [integración CLI](example/cli/main.go).
