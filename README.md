# fs-backup

`fs-backup` is a simple file system backup service written in Go. It allows you to backup specified directories to an S3-compatible storage on a defined schedule.

## Installation

### Download Binary

You can download the latest pre-built binary for `linux-amd64` from the [GitHub Releases page](https://github.com/Costiss/fs-backup/releases/latest).

```bash
# Example for Linux AMD64
wget https://github.com/Costiss/fs-backup/releases/latest/download/fs-backup-linux-amd64 -O /usr/local/bin/fs-backup
chmod +x /usr/local/bin/fs-backup
```

## Configuration

`fs-backup` uses a `config.yaml` file for its settings. An example is provided below. Create this file, for example, at `/etc/fs-backup/config.yaml`.

```yaml
s3:
  bucket: "your-s3-bucket-name"
  region: "your-s3-bucket-region"
  endpoint: "https://s3.magalucloud.com" # Example for Magalu Cloud S3
  access_key: "your-access-key"
  secret_key: "your-secret-key"
backup:
  schedule: "* * * * *" # Every minute (Cron format)
  # log_file: "/var/log/fs-backup/backup.log" # Optonal: Ensure this directory exists and is writable
  # gpg_encrypt_password: "your-encryption-password" # Optional GPG encryption password
  directories:
    - "/path/to/your/directory1"
    - "/path/to/your/directory2"
```

**Important:**

- Replace `your-s3-bucket-name`, `your-s3-bucket-region`, `your-access-key`, and `your-secret-key` with your actual S3 credentials.
- Adjust the `endpoint` if you are using a different S3-compatible service.
- Modify the `schedule` to your desired backup frequency using cron syntax.
- List all directories you wish to backup under `directories`.
- Ensure the `log_file` path is writable by the user running the service.

## Systemd Service Setup

To run `fs-backup` as a systemd service, create a service file (e.g., `/etc/systemd/system/fs-backup.service`) with the following content:

```ini
[Unit]
Description=File System Backup Service

[Service]
Type=simple
ExecStart=/usr/local/bin/fs-backup /etc/fs-backup/config.yaml
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

**Note:** The `ExecStart` line includes `/etc/fs-backup/config.yaml` to specify the configuration file location. Ensure this path matches where your `config.yaml` is stored.

After creating the service file, reload systemd and enable/start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable fs-backup.service
sudo systemctl start fs-backup.service
```

To check the service status:

```bash
sudo systemctl status fs-backup.service
```
