## GenCLI

<div align="center">

[![Releases](https://github.com/Pradumnasaraf/gencli/actions/workflows/releases.yml/badge.svg)](https://github.com/Pradumnasaraf/gencli/actions/workflows/releases.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/Pradumnasaraf/gencli.svg)](https://pkg.go.dev/github.com/Pradumnasaraf/gencli)

</div>

**GenCLI** is an AI-powered CLI tool built with Golang that answers your questions using the [Google Gemini API](https://gemini.google.com). It is developed with [Cobra](https://github.com/spf13/cobra) and more.

![GenCLI GIF](https://github.com/Pradumnasaraf/gencli/assets/51878265/f230a612-c51b-45b1-bbab-772110efcaf4)

### 🚀 Getting Started

To get started with GenCLI, you can install it using the following method:

#### Installation

To install the CLI, use the command below:

```bash
go install github.com/Pradumnasaraf/gencli@latest
```

Go will automatically install it in your `$GOPATH/bin` directory, which should be in your `$PATH`.

#### Usage

Once installed, you can use the `gencli` CLI command. To confirm installation, type `gencli` at the command line.

GenCLI uses the Google Gemini API, so you need to set the API key. To get the API key (It's FREE), visit [here](https://aistudio.google.com/app/apikey?_gl=1*1n5ijhw*_ga*MTQxNDQ2MjcyNi4xNzE5MDU4OTE0*_ga_P1DBVKWT6V*MTcxOTkzNTQzOC4zLjEuMTcxOTkzNTQ3My4yNS4wLjEzODczMjU2OA) and set it in the environment variable `GEMINI_API_KEY`:

```bash
export GEMINI_API_KEY=<API_KEY>
```

The above method sets the API key for the current session only. To set it permanently, add the above line to your `.bashrc` or `.zshrc` file.

> **Note:** If you encounter the error `command not found: gencli`, you need to add `$GOPATH/bin` to your `$PATH` environment variable. For more details, refer to [this guide](https://gist.github.com/Pradumnasaraf/ca6f9a0507089a4c44881446cdda4aa3).

#### Commands

```bash
Usage:
  gencli [flags]
  gencli [command]

Available Commands:
  help        Help about any command
  search      Ask a question and get a response

Flags:
  -h, --help   help for gencli
```

>  eg: gencli search how big is google
>  eg: gencli search what is kubernetes --words 525

### 📜 License

This project is licensed under the Apache-2.0 license - see the [LICENSE](LICENSE) file for details.

### 🛡 Security

If you discover a security vulnerability within this project, please check the [SECURITY](SECURITY.md) for more information.
