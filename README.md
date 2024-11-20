# Docker Cleanup Tool

A command-line tool to automatically clean up unused Docker resources (containers, images, and volumes) based on their age. This tool helps maintain your Docker environment clean and optimized.

## Features

- Clean up stopped containers
- Remove unused Docker images
- Remove unused volumes
- Age-based cleanup (default 24 hours)
- Dry-run mode to preview changes
- Root permission verification
- Docker daemon status check

## Installation

[Installation](INSTALL.md)

## Disclaimer
### This tool is only tested on linux environment so for macos and windows if you get any issues please create issue in github.

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/dhairya13703/docker-cleanup.git
cd docker-cleanup
```

2. Build the binary:
```bash
go build
```

3. (Optional) Install system-wide:
```bash
sudo mv docker-cleanup /usr/local/bin/
```

## Usage

The tool requires root privileges to interact with Docker. Always run with `sudo`.

### Basic Commands

```bash
# Show help and available commands
sudo ./docker-cleanup --help

# Clean all unused resources older than 24 hours (default)
sudo ./docker-cleanup all

# Clean with dry run (preview what would be removed)
sudo ./docker-cleanup all --dry-run

# Clean resources older than specific hours
sudo ./docker-cleanup all --older-than 12
```

### Specific Resource Cleanup

```bash
# Clean only stopped containers
sudo ./docker-cleanup containers
sudo ./docker-cleanup containers --older-than 12

# Clean only unused images
sudo ./docker-cleanup images
sudo ./docker-cleanup images --older-than 12

# Clean only unused volumes
sudo ./docker-cleanup volumes
sudo ./docker-cleanup volumes --older-than 12
```

### Command Options

```bash
--dry-run     Print actions without executing them
--older-than  Remove resources older than specified hours (default 24h)
```

## Resource Types

### Containers
- Removes stopped containers (status: exited or dead)
- Only removes containers that have been stopped longer than the specified time
- Includes removal of associated anonymous volumes

### Images
- Removes images not used by any containers
- Removes images older than the specified time
- Includes cleanup of dangling images
- Preserves images that are in use by containers (running or stopped)

### Volumes
- Removes volumes not mounted in any container
- Preserves volumes that are currently in use
- Forces removal of unused volumes

## Setting Up Automatic Cleanup

### Using Cron

To set up automatic cleanup, add a cron job:

1. Open the crontab editor:
```bash
sudo crontab -e
```

2. Add a line to run cleanup daily at midnight:
```bash
0 0 * * * /usr/local/bin/docker-cleanup all --older-than 24
```

### Using Systemd Timer

1. Create a service file `/etc/systemd/system/docker-cleanup.service`:
```ini
[Unit]
Description=Docker Cleanup Service
After=docker.service

[Service]
Type=oneshot
ExecStart=/usr/local/bin/docker-cleanup all --older-than 24

[Install]
WantedBy=multi-user.target
```

2. Create a timer file `/etc/systemd/system/docker-cleanup.timer`:
```ini
[Unit]
Description=Run Docker Cleanup daily

[Timer]
OnCalendar=daily
Persistent=true

[Install]
WantedBy=timers.target
```

3. Enable and start the timer:
```bash
sudo systemctl enable docker-cleanup.timer
sudo systemctl start docker-cleanup.timer
```

## Examples

```bash
# Preview what would be cleaned up
sudo ./docker-cleanup all --dry-run

# Clean containers older than 12 hours
sudo ./docker-cleanup containers --older-than 12

# Clean everything older than 1 day
sudo ./docker-cleanup all --older-than 24

# Clean only unused images
sudo ./docker-cleanup images

# Clean volumes with dry run
sudo ./docker-cleanup volumes --dry-run
```

## Safety Features

- Dry run mode to preview actions
- Checks for root privileges
- Verifies Docker daemon is running
- Preserves currently used resources
- Age-based filtering to prevent accidental removal of new resources

## Troubleshooting

If you encounter issues:

1. Ensure Docker is running:
```bash
sudo systemctl status docker
```

2. Verify you have root privileges:
```bash
sudo ./docker-cleanup --help
```

3. Check Docker daemon accessibility:
```bash
sudo docker ps
```

4. Use dry-run mode to debug:
```bash
sudo ./docker-cleanup all --dry-run
```
