# Duplicate Cleaner Configuration

The Duplicate Cleaner is a command-line application for identifying and managing duplicate files. This document provides instructions on how to define a configuration file for running the application.

## Configuration File Format

The configuration file is in YAML format. Below is an example configuration file with explanations for each parameter:

```yaml
source_folder: "/path/to/source"
backup_folder: "/path/to/backup"
max_scan_depth: 5
excluded_file_types:
- ".tmp"
- ".log"
concurrent_scan: true
max_concurrent_scans: 4
duplicate_threshold: "2h30m"
dry_run: false
# Add other configuration parameters as needed...
```

### Configuration Parameters

- **source_folder:** The folder to scan for duplicate files (required).
- **backup_folder:** The backup folder where duplicate files will be moved (required).
- **max_scan_depth:** Maximum depth level to scan recursively (optional, default: infinite).
- **excluded_file_types:** File types to exclude from scanning (optional, default: none).
- **concurrent_scan:** Enable concurrency for faster scanning (optional, default: true).
- **max_concurrent_scans:** Maximum number of concurrent scans (optional, default: based on available CPUs).
- **duplicate_threshold:** Time threshold for considering files as duplicates (optional, default: 1 hour).
- **dry_run:** Whether to dry-run (display duplicates without removing them) (optional, default: false).
- *Add other fields as needed...*

## Running the Application

1. **Build the Executable:**
    - Ensure you have Go installed: [Go Installation Guide](https://golang.org/doc/install)
    - Open a terminal and navigate to the project directory.
    - Run the following command to build the executable:
      \```bash
      go build -o duplicate-cleaner.exe
      \```

2. **Create a Configuration File:**
    - Create a YAML file (e.g., `config.yaml`) and define your configuration parameters based on the example above.

3. **Run the Executable:**
    - Open a terminal and run the executable with the path to your configuration file:
      ```bash
      ./duplicate-cleaner.exe -config=config.yaml
      ```

   Replace `config.yaml` with the actual path to your configuration file.

4. **View Results:**
    - The application will scan for duplicates based on the provided configuration and display the results.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Make sure to customize the configuration file based on your specific requirements.
- Feel free to add any additional configuration parameters needed for your future features.
