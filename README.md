# Osquery Extension

## Project Summary
This project provides custom Osquery tables that enhance the default Osquery functionality on macOS and Windows. These tables can help with forensics, compliance, and security investigations by extracting critical configuration and usage data from local installations.

## Usage
For testing, you can load the extension with `osqueryi`.

By default, osquery does not want to load extensions not owned by root. You can either change the ownership of osquery_extension.ext to root, or run osquery with the `--allow_unsafe` flag.

To test:
```bash
make osqueryi # Will run osqueryi --extension /path/to/osquery_extension.ext --allow_unsafe in the background
```

For production deployment, you should refer to the [osquery documentation](https://osquery.readthedocs.io/en/stable/deployment/extensions/).

## Tables

|Table|Description|Platforms|Notes|
|----|----|----|----|
| `chrome_extensions_dns` | Inspired by [ExtensionHound](https://github.com/arsolutioner/ExtensionHound), this table returns the DNS domains requested by chromium browser extensions. | macOS / Windows |
| `chrome_preferences` | Parses different Chromium based browser preferences such as sites with access to geolocation data, microphone access and notifications. Useful for forensics purposes. | macOS / Windows |
| `vscode_extensions` | Returns VSCode extensions installed on host. This table has been eventually incorporated into Osquery core. | macOS / Windows |
