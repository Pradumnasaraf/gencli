# GenCLI

<div align="center">

[![Releases](https://github.com/Pradumnasaraf/gencli/actions/workflows/releases.yml/badge.svg)](https://github.com/Pradumnasaraf/gencli/actions/workflows/releases.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/Pradumnasaraf/gencli.svg)](https://pkg.go.dev/github.com/Pradumnasaraf/gencli)

</div>

**GenCLI** is an AI-powered CLI tool that answers your questions. It is built with [Go](https://github.com/golang/go), [Cobra](https://github.com/spf13/cobra), [Viper](https://github.com/spf13/viper), [Gemini](https://gemini.google.com), and more.

## 🚀 Getting Started

To get started with GenCLI, you can install it using the following method:

### Installation

To install the CLI, use the command below:

```bash
go install github.com/Pradumnasaraf/gencli@latest
```

Go will automatically install it in your `$GOPATH/bin` directory, which should be in your `$PATH`.

### Usage

Once installed, you can use the `gencli` CLI command. To confirm installation, type `gencli` at the command line.

GenCLI uses the Google Gemini API, so you need to set the API key. To get the API key (It's FREE), visit [here](https://aistudio.google.com/app/apikey?_gl=1*1n5ijhw*_ga*MTQxNDQ2MjcyNi4xNzE5MDU4OTE0*_ga_P1DBVKWT6V*MTcxOTkzNTQzOC4zLjEuMTcxOTkzNTQ3My4yNS4wLjEzODczMjU2OA) and set it in the environment variable `GENCLI_API_KEY`:

```bash
export GENCLI_API_KEY=<API_KEY>
```

The above method sets the API key for the current session only. To set it permanently, add the above line to your `.bashrc` or `.zshrc` file.

> **Note:** If you encounter the error `command not found: gencli`, you need to add `$GOPATH/bin` to your `$PATH` environment variable. For more details, refer to [this guide](https://gist.github.com/Pradumnasaraf/ca6f9a0507089a4c44881446cdda4aa3).

### Commands

```
Usage:
  gencli [command]

Available Commands:
  gencli  Type your question here
```

## 📜 License

This project is licensed under the Apache-2.0 license - see the [LICENSE](LICENSE) file for details.

## 🛡 Security

If you discover a security vulnerability within this project, please check the [SECURITY](SECURITY.md) for more information.