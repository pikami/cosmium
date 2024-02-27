# Cosmium

Cosmium is a lightweight Cosmos DB emulator designed to facilitate local development and testing. While it aims to provide developers with a solution for running a local database during development, it's important to note that it's not 100% compatible with Cosmos DB. However, it serves as a convenient tool for E2E or integration tests during the CI/CD pipeline.

One of Cosmium's notable features is its ability to save and load state to a single JSON file. This feature makes it easy to load different test cases or share state with other developers, enhancing collaboration and efficiency in development workflows.

# Getting Started
### Downloading Cosmium

You can download the latest version of Cosmium from the [GitHub Releases page](https://github.com/pikami/cosmium/releases). Choose the appropriate release for your operating system and architecture.

### Supported Platforms

Cosmium is available for the following platforms:

* **Linux**: cosmium-linux-amd64
* **macOS**: cosmium-darwin-amd64
* **macOS on Apple Silicon**: cosmium-darwin-arm64
* **Windows**: cosmium-windows-amd64.exe

### Running Cosmium

Once downloaded, you can launch Cosmium using the following command:

```sh
cosmium -Persist "./save.json"
```

Connection String Example:
```
AccountEndpoint=https://localhost:8081/;AccountKey=C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==;
```

###  Running Cosmos DB Explorer

If you want to run Cosmos DB Explorer alongside Cosmium, you'll need to build it yourself and point the `-ExplorerDir` argument to the dist directory. Please refer to the [Cosmos DB Explorer repository](https://github.com/Azure/cosmos-explorer) for instructions on building the application.

Once running, the explorer can be reached by navigating following URL: `https://127.0.0.1:8081/_explorer/` (might be different depending on your configuration).

### SSL Certificate

By default, Cosmium uses a pre-generated SSL certificate. You can provide your own certificates by specifying paths to the SSL certificate and key (PEM format) using the `-Cert` and `-CertKey` arguments, respectively.

To disable SSL and run Cosmium on HTTP instead, you can use the `-DisableTls` flag. However most applications will require HTTPS.

### Other Available Arguments

* **-AccountKey**: Account key for authentication (default "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==")
* **-DisableAuth**: Disable authentication
* **-Host**: Hostname (default "localhost")
* **-InitialData**: Path to JSON containing initial state
* **-Persist**: Saves data to the given path on application exit (When `-InitialData` argument is not supplied, it will try to load data from path supplied in `-Persist`)
* **-Port**: Listen port (default 8081)

These arguments allow you to configure various aspects of Cosmium's behavior according to your requirements.

# License
This project is [MIT licensed](./LICENSE).
