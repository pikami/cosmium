# Cosmium

Cosmium is a lightweight Cosmos DB emulator designed to facilitate local development and testing. While it aims to provide developers with a solution for running a local database during development, it's important to note that it's not 100% compatible with Cosmos DB. However, it serves as a convenient tool for E2E or integration tests during the CI/CD pipeline. Read more about compatibility [here](./docs/COMPATIBILITY.md).

One of Cosmium's notable features is its ability to save and load state to a single JSON file. This feature makes it easy to load different test cases or share state with other developers, enhancing collaboration and efficiency in development workflows.

# Getting Started

### Installation via Homebrew

You can install Cosmium using Homebrew by adding the `pikami/brew` tap and then installing the package.

```sh
brew tap pikami/brew
brew install cosmium
```

This will download and install Cosmium on your system, making it easy to manage and update using Homebrew.

### Downloading Cosmium Binaries

You can download the latest version of Cosmium from the [GitHub Releases page](https://github.com/pikami/cosmium/releases). Choose the appropriate release for your operating system and architecture.

### Supported Platforms

Cosmium is available for the following platforms:

- **Linux**: cosmium-linux-amd64
- **Linux on ARM**: cosmium-linux-arm64
- **macOS**: cosmium-darwin-amd64
- **macOS on Apple Silicon**: cosmium-darwin-arm64
- **Windows**: cosmium-windows-amd64.exe
- **Windows on ARM**: cosmium-windows-arm64.exe

### Running Cosmium

Once downloaded, you can launch Cosmium using the following command:

```sh
cosmium -Persist "./save.json"
```

Connection String Example:

```
AccountEndpoint=https://localhost:8081/;AccountKey=C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==;
```

### Running Cosmos DB Explorer

If you want to run Cosmos DB Explorer alongside Cosmium, you'll need to build it yourself and point the `-ExplorerDir` argument to the dist directory. Please refer to the [Cosmos DB Explorer repository](https://github.com/Azure/cosmos-explorer) for instructions on building the application.

There's also a prebuilt docker image that includes the explorer: `ghcr.io/pikami/cosmium:explorer`

Once running, the explorer can be reached by navigating following URL: `https://127.0.0.1:8081/_explorer/` (might be different depending on your configuration).

### Running with docker (optional)

There are two docker tags available:

- ghcr.io/pikami/cosmium:latest - Cosmium core service
- ghcr.io/pikami/cosmium:explorer - Cosmium with database explorer available on `https://127.0.0.1:8081/_explorer/`

If you wan to run the application using docker, configure it using environment variables see example:

```sh
# Ensure save.json exists so Docker volume mounts correctly
[ -f save.json ] || echo '{}' > save.json && docker run --rm \
  -e COSMIUM_PERSIST=/save.json \
  -v ./save.json:/save.json \
  -p 8081:8081 \
  ghcr.io/pikami/cosmium # or `ghcr.io/pikami/cosmium:explorer`
```

### SSL Certificate

By default, Cosmium uses a pre-generated SSL certificate. You can provide your own certificates by specifying paths to the SSL certificate and key (PEM format) using the `-Cert` and `-CertKey` arguments, respectively.

To disable SSL and run Cosmium on HTTP instead, you can use the `-DisableTls` flag. However most applications will require HTTPS.

### Other Available Arguments

- **-AccountKey**: Account key for authentication (default "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==")
- **-DisableAuth**: Disable authentication
- **-Host**: Hostname (default "localhost")
- **-InitialData**: Path to JSON containing initial state
- **-Persist**: Saves data to the given path on application exit (When `-InitialData` argument is not supplied, it will try to load data from path supplied in `-Persist`)
- **-Port**: Listen port (default 8081)
- **-LogLevel**: Sets the logging level (one of: debug, info, error, silent) (default info)
- **-DataStore**: Allows selecting [storage backend](#data-storage-backends) (default "json")

These arguments allow you to configure various aspects of Cosmium's behavior according to your requirements.

All mentioned arguments can also be set using environment variables:

- **COSMIUM_ACCOUNTKEY** for `-AccountKey`
- **COSMIUM_DISABLEAUTH** for `-DisableAuth`
- **COSMIUM_HOST** for `-Host`
- **COSMIUM_INITIALDATA** for `-InitialData`
- **COSMIUM_PERSIST** for `-Persist`
- **COSMIUM_PORT** for `-Port`
- **COSMIUM_LOGLEVEL** for `-LogLevel`

### Data Storage Backends

Cosmium supports multiple storage backends for saving, loading, and managing data at runtime.

| Backend  | Storage Location         | Write Behavior           | Memory Usage         |
|----------|--------------------------|--------------------------|----------------------|
| `json` (default) | JSON file on disk 📄 | On application exit ⏳ | 🛑 More than Badger |
| `badger`  | BadgerDB database on disk ⚡ | Immediately on write 🚀 | ✅ Less than JSON  |


The `badger` backend is generally recommended as it uses less memory and writes data to disk immediately. However, if you need to load initial data from a JSON file, use the `json` backend.

# License

This project is [MIT licensed](./LICENSE).
