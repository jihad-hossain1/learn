# Part 7: Operating Systems for DevOps

## Introduction

Operating systems form the foundation of DevOps infrastructure. Understanding how operating systems work, their components, and how to manage them effectively is crucial for DevOps engineers. This guide covers essential OS concepts, Linux administration, Windows management, and cross-platform considerations.

## Operating System Fundamentals

### What is an Operating System?

An operating system (OS) is system software that manages computer hardware, software resources, and provides common services for computer programs.

**Key Functions**:
- **Process Management**: Creating, scheduling, and terminating processes
- **Memory Management**: Allocating and deallocating memory space
- **File System Management**: Organizing and accessing files
- **Device Management**: Controlling hardware devices
- **Security**: User authentication and access control
- **Networking**: Network communication and protocols

### OS Architecture

```
┌─────────────────────────────────────┐
│           User Applications         │
├─────────────────────────────────────┤
│           System Libraries          │
├─────────────────────────────────────┤
│              System Calls           │
├─────────────────────────────────────┤
│              Kernel Space           │
│  ┌─────────┬─────────┬─────────┐    │
│  │ Process │ Memory  │  File   │    │
│  │ Manager │ Manager │ System  │    │
│  └─────────┴─────────┴─────────┘    │
│  ┌─────────┬─────────┬─────────┐    │
│  │ Device  │ Network │Security │    │
│  │ Drivers │ Stack   │ Manager │    │
│  └─────────┴─────────┴─────────┘    │
├─────────────────────────────────────┤
│              Hardware               │
└─────────────────────────────────────┘
```

## Linux Operating System

### Linux Distribution Landscape

**Enterprise Distributions**:
- **Red Hat Enterprise Linux (RHEL)**: Enterprise-focused, commercial support
- **CentOS/Rocky Linux**: Community versions of RHEL
- **Ubuntu LTS**: Long-term support, popular for servers
- **SUSE Linux Enterprise**: Enterprise distribution with strong container support

**Development/Cloud Distributions**:
- **Ubuntu**: Popular for development and cloud deployments
- **Debian**: Stable, community-driven
- **Amazon Linux**: Optimized for AWS
- **Alpine Linux**: Minimal, security-focused, popular for containers

### Linux File System Hierarchy

```
/
├── bin/          # Essential user binaries
├── boot/         # Boot loader files
├── dev/          # Device files
├── etc/          # System configuration files
├── home/         # User home directories
├── lib/          # Essential shared libraries
├── media/        # Removable media mount points
├── mnt/          # Temporary mount points
├── opt/          # Optional software packages
├── proc/         # Process and kernel information
├── root/         # Root user home directory
├── run/          # Runtime data
├── sbin/         # Essential system binaries
├── srv/          # Service data
├── sys/          # System information
├── tmp/          # Temporary files
├── usr/          # User utilities and applications
│   ├── bin/      # User binaries
│   ├── lib/      # User libraries
│   ├── local/    # Local software
│   └── share/    # Shared data
└── var/          # Variable data
    ├── log/      # Log files
    ├── mail/     # Mail spool
    ├── run/      # Runtime data
    └── tmp/      # Temporary files
```

### Essential Linux Commands

#### File and Directory Operations

```bash
# Navigation
ls -la                    # List files with details
cd /path/to/directory     # Change directory
pwd                       # Print working directory
find /path -name "*.txt"  # Find files by name
locate filename           # Locate files quickly
which command             # Find command location

# File Operations
cp source destination     # Copy files
mv source destination     # Move/rename files
rm -rf directory         # Remove files/directories
ln -s target link        # Create symbolic link
chmod 755 file           # Change file permissions
chown user:group file    # Change file ownership

# File Content
cat file.txt             # Display file content
less file.txt            # View file with pagination
head -n 10 file.txt      # Show first 10 lines
tail -f /var/log/app.log # Follow log file
grep "pattern" file.txt   # Search text patterns
sed 's/old/new/g' file   # Stream editor
awk '{print $1}' file    # Text processing

# File Compression
tar -czf archive.tar.gz directory/  # Create compressed archive
tar -xzf archive.tar.gz             # Extract archive
zip -r archive.zip directory/       # Create zip archive
unzip archive.zip                   # Extract zip archive
```

#### Process Management

```bash
# Process Information
ps aux                   # List all processes
ps -ef | grep nginx      # Find specific processes
top                      # Real-time process monitor
htop                     # Enhanced process monitor
pstree                   # Process tree view

# Process Control
kill PID                 # Terminate process by PID
killall process_name     # Kill processes by name
pkill -f pattern         # Kill processes matching pattern
nohup command &          # Run command in background
jobs                     # List background jobs
fg %1                    # Bring job to foreground
bg %1                    # Send job to background

# System Monitoring
uptime                   # System uptime and load
free -h                  # Memory usage
df -h                    # Disk usage
du -sh directory/        # Directory size
iostat                   # I/O statistics
vmstat                   # Virtual memory statistics
```

#### Network Operations

```bash
# Network Information
ip addr show             # Show network interfaces
ip route show            # Show routing table
netstat -tulpn           # Show listening ports
ss -tulpn                # Modern netstat alternative
lsof -i :80              # Show processes using port 80

# Network Testing
ping google.com          # Test connectivity
traceroute google.com    # Trace network path
nslookup domain.com      # DNS lookup
dig domain.com           # DNS information
curl -I http://site.com  # HTTP headers
wget http://site.com/file # Download files

# Network Configuration
sudo ip addr add 192.168.1.100/24 dev eth0  # Add IP address
sudo ip route add default via 192.168.1.1   # Add default route
sudo systemctl restart networking            # Restart networking
```

#### System Services (systemd)

```bash
# Service Management
sudo systemctl start nginx       # Start service
sudo systemctl stop nginx        # Stop service
sudo systemctl restart nginx     # Restart service
sudo systemctl reload nginx      # Reload configuration
sudo systemctl enable nginx      # Enable at boot
sudo systemctl disable nginx     # Disable at boot

# Service Status
systemctl status nginx           # Service status
systemctl is-active nginx        # Check if active
systemctl is-enabled nginx       # Check if enabled
systemctl list-units --type=service  # List all services

# Logs
journalctl -u nginx              # Service logs
journalctl -f                    # Follow system logs
journalctl --since "1 hour ago"  # Recent logs
journalctl -p err                # Error logs only
```

### User and Permission Management

#### User Management

```bash
# User Operations
sudo useradd -m -s /bin/bash username    # Create user
sudo usermod -aG sudo username           # Add user to sudo group
sudo userdel -r username                 # Delete user and home
sudo passwd username                     # Change user password
su - username                            # Switch user
sudo -u username command                 # Run command as user

# User Information
whoami                                   # Current user
id username                              # User ID information
groups username                          # User groups
w                                        # Who is logged in
last                                     # Login history
```

#### File Permissions

```bash
# Permission Types
# r (read) = 4, w (write) = 2, x (execute) = 1
# Owner, Group, Others

# Examples
chmod 755 file.sh        # rwxr-xr-x
chmod 644 file.txt       # rw-r--r--
chmod +x script.sh       # Add execute permission
chmod -R 755 directory/  # Recursive permission change

# Advanced Permissions
chmod u+s file           # Set SUID bit
chmod g+s directory      # Set SGID bit
chmod +t directory       # Set sticky bit

# Access Control Lists (ACL)
getfacl file.txt         # Get ACL
setfacl -m u:user:rw file.txt  # Set user ACL
setfacl -x u:user file.txt     # Remove user ACL
```

### Package Management

#### APT (Debian/Ubuntu)

```bash
# Package Operations
sudo apt update                    # Update package list
sudo apt upgrade                   # Upgrade packages
sudo apt install package_name      # Install package
sudo apt remove package_name       # Remove package
sudo apt purge package_name        # Remove package and config
sudo apt autoremove                # Remove unused packages

# Package Information
apt search keyword                 # Search packages
apt show package_name              # Package information
apt list --installed               # List installed packages
apt list --upgradable              # List upgradable packages

# Repository Management
sudo add-apt-repository ppa:user/repo  # Add PPA
sudo apt-key add key.gpg               # Add GPG key
```

#### YUM/DNF (Red Hat/CentOS/Fedora)

```bash
# Package Operations
sudo yum update                    # Update packages (CentOS 7)
sudo dnf update                    # Update packages (CentOS 8+)
sudo yum install package_name      # Install package
sudo yum remove package_name       # Remove package
sudo yum clean all                 # Clean cache

# Package Information
yum search keyword                 # Search packages
yum info package_name              # Package information
yum list installed                 # List installed packages
yum history                        # Transaction history

# Repository Management
sudo yum-config-manager --add-repo repo_url  # Add repository
sudo yum-config-manager --enable repo_name   # Enable repository
```

### System Monitoring and Performance

#### System Resources

```bash
# CPU Monitoring
top                      # Real-time CPU usage
htop                     # Enhanced top
sar -u 1 5               # CPU utilization
mpstat                   # Multi-processor statistics

# Memory Monitoring
free -h                  # Memory usage
cat /proc/meminfo        # Detailed memory info
sar -r 1 5               # Memory utilization
vmstat 1 5               # Virtual memory stats

# Disk Monitoring
df -h                    # Disk space usage
du -sh /path/*           # Directory sizes
iostat -x 1 5            # I/O statistics
sar -d 1 5               # Disk activity
lsblk                    # Block devices

# Network Monitoring
iftop                    # Network traffic by connection
nload                    # Network traffic by interface
sar -n DEV 1 5           # Network interface statistics
netstat -i               # Interface statistics
```

#### Log Analysis

```bash
# System Logs
tail -f /var/log/syslog          # Follow system log
grep "ERROR" /var/log/app.log     # Search for errors
zcat /var/log/app.log.gz | grep "pattern"  # Search compressed logs

# Log Rotation
logrotate -d /etc/logrotate.conf # Test log rotation
logrotate -f /etc/logrotate.conf # Force log rotation

# Journal Logs (systemd)
journalctl -xe                   # Recent logs with explanations
journalctl -u service_name       # Service-specific logs
journalctl --since "2023-01-01"  # Logs since date
journalctl -p err                # Error priority logs
```

### Shell Scripting for DevOps

#### Basic Shell Script Structure

```bash
#!/bin/bash
# deployment-script.sh

# Script metadata
SCRIPT_NAME="Deployment Script"
VERSION="1.0.0"
AUTHOR="DevOps Team"

# Configuration
APP_NAME="myapp"
DEPLOY_DIR="/opt/${APP_NAME}"
BACKUP_DIR="/backup/${APP_NAME}"
LOG_FILE="/var/log/${APP_NAME}-deploy.log"

# Functions
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

error_exit() {
    log "ERROR: $1"
    exit 1
}

check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check if running as root
    if [[ $EUID -eq 0 ]]; then
        error_exit "This script should not be run as root"
    fi
    
    # Check required commands
    local required_commands=("git" "docker" "curl")
    for cmd in "${required_commands[@]}"; do
        if ! command -v "$cmd" &> /dev/null; then
            error_exit "Required command '$cmd' not found"
        fi
    done
    
    log "Prerequisites check passed"
}

create_backup() {
    log "Creating backup..."
    
    if [[ -d "$DEPLOY_DIR" ]]; then
        local backup_name="${APP_NAME}-$(date +%Y%m%d-%H%M%S)"
        sudo mkdir -p "$BACKUP_DIR"
        sudo tar -czf "${BACKUP_DIR}/${backup_name}.tar.gz" -C "$(dirname "$DEPLOY_DIR")" "$(basename "$DEPLOY_DIR")"
        log "Backup created: ${BACKUP_DIR}/${backup_name}.tar.gz"
    else
        log "No existing deployment found, skipping backup"
    fi
}

deploy_application() {
    log "Deploying application..."
    
    # Clone or update repository
    if [[ -d "$DEPLOY_DIR" ]]; then
        cd "$DEPLOY_DIR" || error_exit "Cannot change to deploy directory"
        git pull origin main || error_exit "Git pull failed"
    else
        sudo mkdir -p "$(dirname "$DEPLOY_DIR")"
        git clone "https://github.com/company/${APP_NAME}.git" "$DEPLOY_DIR" || error_exit "Git clone failed"
        cd "$DEPLOY_DIR" || error_exit "Cannot change to deploy directory"
    fi
    
    # Build application
    log "Building application..."
    docker build -t "${APP_NAME}:latest" . || error_exit "Docker build failed"
    
    # Deploy with docker-compose
    log "Starting services..."
    docker-compose down || true
    docker-compose up -d || error_exit "Docker compose up failed"
    
    log "Application deployed successfully"
}

health_check() {
    log "Performing health check..."
    
    local max_attempts=30
    local attempt=1
    
    while [[ $attempt -le $max_attempts ]]; do
        if curl -f -s "http://localhost:8080/health" > /dev/null; then
            log "Health check passed"
            return 0
        fi
        
        log "Health check attempt $attempt/$max_attempts failed, retrying..."
        sleep 10
        ((attempt++))
    done
    
    error_exit "Health check failed after $max_attempts attempts"
}

cleanup() {
    log "Performing cleanup..."
    
    # Remove old Docker images
    docker image prune -f
    
    # Remove old backups (keep last 5)
    if [[ -d "$BACKUP_DIR" ]]; then
        cd "$BACKUP_DIR" || return
        ls -t ${APP_NAME}-*.tar.gz 2>/dev/null | tail -n +6 | xargs -r rm -f
    fi
    
    log "Cleanup completed"
}

# Main execution
main() {
    log "Starting deployment of $APP_NAME"
    
    check_prerequisites
    create_backup
    deploy_application
    health_check
    cleanup
    
    log "Deployment completed successfully"
}

# Error handling
set -euo pipefail
trap 'error_exit "Script failed at line $LINENO"' ERR

# Execute main function
main "$@"
```

#### Advanced Shell Scripting Techniques

```bash
#!/bin/bash
# advanced-monitoring.sh

# Configuration file parsing
load_config() {
    local config_file="${1:-/etc/monitoring.conf}"
    
    if [[ -f "$config_file" ]]; then
        # Source configuration file
        source "$config_file"
    else
        # Default configuration
        THRESHOLD_CPU=80
        THRESHOLD_MEMORY=85
        THRESHOLD_DISK=90
        ALERT_EMAIL="admin@company.com"
        SLACK_WEBHOOK="https://hooks.slack.com/..."
    fi
}

# CPU monitoring
check_cpu_usage() {
    local cpu_usage
    cpu_usage=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)
    cpu_usage=${cpu_usage%.*}  # Remove decimal part
    
    if [[ $cpu_usage -gt $THRESHOLD_CPU ]]; then
        send_alert "HIGH CPU USAGE" "CPU usage is ${cpu_usage}% (threshold: ${THRESHOLD_CPU}%)"
    fi
    
    echo "CPU: ${cpu_usage}%"
}

# Memory monitoring
check_memory_usage() {
    local memory_usage
    memory_usage=$(free | grep Mem | awk '{printf "%.0f", $3/$2 * 100.0}')
    
    if [[ $memory_usage -gt $THRESHOLD_MEMORY ]]; then
        send_alert "HIGH MEMORY USAGE" "Memory usage is ${memory_usage}% (threshold: ${THRESHOLD_MEMORY}%)"
    fi
    
    echo "Memory: ${memory_usage}%"
}

# Disk monitoring
check_disk_usage() {
    while IFS= read -r line; do
        local usage filesystem
        usage=$(echo "$line" | awk '{print $5}' | sed 's/%//')
        filesystem=$(echo "$line" | awk '{print $6}')
        
        if [[ $usage -gt $THRESHOLD_DISK ]]; then
            send_alert "HIGH DISK USAGE" "Disk usage on $filesystem is ${usage}% (threshold: ${THRESHOLD_DISK}%)"
        fi
        
        echo "Disk $filesystem: ${usage}%"
    done < <(df -h | grep -E '^/dev/')
}

# Service monitoring
check_services() {
    local services=("nginx" "mysql" "redis" "docker")
    
    for service in "${services[@]}"; do
        if systemctl is-active --quiet "$service"; then
            echo "Service $service: Running"
        else
            send_alert "SERVICE DOWN" "Service $service is not running"
            echo "Service $service: Stopped"
        fi
    done
}

# Port monitoring
check_ports() {
    local ports=(80 443 22 3306 6379)
    
    for port in "${ports[@]}"; do
        if ss -tuln | grep -q ":$port "; then
            echo "Port $port: Open"
        else
            send_alert "PORT CLOSED" "Port $port is not listening"
            echo "Port $port: Closed"
        fi
    done
}

# Alert functions
send_email_alert() {
    local subject="$1"
    local message="$2"
    
    echo "$message" | mail -s "[ALERT] $subject" "$ALERT_EMAIL"
}

send_slack_alert() {
    local subject="$1"
    local message="$2"
    
    local payload
    payload=$(cat <<EOF
{
    "text": "🚨 *$subject*",
    "attachments": [
        {
            "color": "danger",
            "text": "$message",
            "footer": "$(hostname)",
            "ts": $(date +%s)
        }
    ]
}
EOF
    )
    
    curl -X POST -H 'Content-type: application/json' \
         --data "$payload" \
         "$SLACK_WEBHOOK"
}

send_alert() {
    local subject="$1"
    local message="$2"
    
    # Log alert
    logger -t monitoring "ALERT: $subject - $message"
    
    # Send email if configured
    if [[ -n "$ALERT_EMAIL" ]]; then
        send_email_alert "$subject" "$message"
    fi
    
    # Send Slack notification if configured
    if [[ -n "$SLACK_WEBHOOK" ]]; then
        send_slack_alert "$subject" "$message"
    fi
}

# JSON output
generate_json_report() {
    local timestamp
    timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    cat <<EOF
{
    "timestamp": "$timestamp",
    "hostname": "$(hostname)",
    "system": {
        "uptime": "$(uptime -p)",
        "load_average": "$(uptime | awk -F'load average:' '{print $2}' | xargs)",
        "cpu_usage": "$(top -bn1 | grep 'Cpu(s)' | awk '{print $2}' | cut -d'%' -f1)%",
        "memory_usage": "$(free | grep Mem | awk '{printf "%.1f", $3/$2 * 100.0}')%",
        "disk_usage": [
EOF
    
    local first=true
    while IFS= read -r line; do
        local usage filesystem
        usage=$(echo "$line" | awk '{print $5}')
        filesystem=$(echo "$line" | awk '{print $6}')
        
        if [[ "$first" == "true" ]]; then
            first=false
        else
            echo ","
        fi
        
        echo "            {\"filesystem\": \"$filesystem\", \"usage\": \"$usage\"}"
    done < <(df -h | grep -E '^/dev/')
    
    cat <<EOF
        ]
    }
}
EOF
}

# Main function
main() {
    local output_format="${1:-text}"
    
    load_config
    
    case "$output_format" in
        "json")
            generate_json_report
            ;;
        "text"|*)
            echo "=== System Monitoring Report ==="
            echo "Timestamp: $(date)"
            echo "Hostname: $(hostname)"
            echo ""
            
            echo "--- System Resources ---"
            check_cpu_usage
            check_memory_usage
            check_disk_usage
            echo ""
            
            echo "--- Services ---"
            check_services
            echo ""
            
            echo "--- Network Ports ---"
            check_ports
            ;;
    esac
}

# Execute main function
main "$@"
```

## Windows Operating System

### Windows Server Management

#### PowerShell Fundamentals

```powershell
# Basic Commands
Get-Help Get-Process        # Get help for commands
Get-Command *Service*       # Find commands with 'Service'
Get-Member                  # Show object properties and methods

# File and Directory Operations
Get-ChildItem -Path C:\     # List directory contents (ls equivalent)
Set-Location C:\Windows     # Change directory (cd equivalent)
New-Item -ItemType Directory -Path "C:\MyFolder"  # Create directory
Copy-Item -Path "source" -Destination "dest"     # Copy files
Remove-Item -Path "file.txt" -Force              # Delete files

# File Content
Get-Content file.txt        # Display file content (cat equivalent)
Select-String "pattern" file.txt  # Search text (grep equivalent)
Out-File -FilePath "output.txt"   # Redirect output to file
```

#### Windows Services Management

```powershell
# Service Operations
Get-Service                 # List all services
Get-Service -Name "Spooler" # Get specific service
Start-Service -Name "Spooler"     # Start service
Stop-Service -Name "Spooler"      # Stop service
Restart-Service -Name "Spooler"   # Restart service
Set-Service -Name "Spooler" -StartupType Automatic  # Set startup type

# Service Status
Get-Service | Where-Object {$_.Status -eq "Stopped"}  # Stopped services
Get-Service | Where-Object {$_.Status -eq "Running"}  # Running services

# Event Logs
Get-EventLog -LogName System -Newest 10              # Recent system events
Get-EventLog -LogName Application -EntryType Error   # Application errors
Get-WinEvent -FilterHashtable @{LogName='System'; Level=2}  # Error events
```

#### Process Management

```powershell
# Process Information
Get-Process                 # List all processes
Get-Process -Name "notepad" # Get specific process
Get-Process | Sort-Object CPU -Descending  # Sort by CPU usage

# Process Control
Start-Process "notepad.exe" # Start process
Stop-Process -Name "notepad" -Force        # Kill process
Stop-Process -Id 1234       # Kill process by ID

# System Information
Get-ComputerInfo            # System information
Get-WmiObject -Class Win32_OperatingSystem  # OS information
Get-WmiObject -Class Win32_LogicalDisk      # Disk information
```

#### Network Management

```powershell
# Network Configuration
Get-NetIPConfiguration      # Network configuration
Get-NetAdapter              # Network adapters
Get-NetRoute                # Routing table
Test-NetConnection -ComputerName "google.com" -Port 80  # Test connectivity

# Firewall Management
Get-NetFirewallRule         # List firewall rules
New-NetFirewallRule -DisplayName "Allow HTTP" -Direction Inbound -Protocol TCP -LocalPort 80
Enable-NetFirewallRule -DisplayName "Allow HTTP"
Disable-NetFirewallRule -DisplayName "Allow HTTP"
```

### Windows Package Management

#### Chocolatey

```powershell
# Install Chocolatey
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))

# Package Operations
choco install git           # Install package
choco upgrade git           # Upgrade package
choco uninstall git         # Uninstall package
choco list --local-only     # List installed packages
choco search nodejs         # Search packages
```

#### Windows Package Manager (winget)

```powershell
# Package Operations
winget install Git.Git      # Install package
winget upgrade Git.Git      # Upgrade package
winget uninstall Git.Git    # Uninstall package
winget list                 # List installed packages
winget search nodejs        # Search packages
winget upgrade --all        # Upgrade all packages
```

## Cross-Platform Considerations

### Container Technologies

#### Docker on Different Platforms

**Linux Docker Commands**:
```bash
# Container Management
docker run -d --name webapp -p 80:80 nginx
docker ps                    # List running containers
docker logs webapp           # View container logs
docker exec -it webapp bash  # Execute shell in container
docker stop webapp           # Stop container
docker rm webapp             # Remove container

# Image Management
docker images                # List images
docker pull ubuntu:20.04     # Pull image
docker build -t myapp .      # Build image
docker rmi image_id          # Remove image

# System Management
docker system df             # Disk usage
docker system prune          # Clean up unused resources
```

**Windows Docker Commands** (same syntax, different underlying technology):
```powershell
# Windows containers
docker run -d --name webapp -p 80:80 mcr.microsoft.com/windows/servercore/iis
docker exec -it webapp powershell

# Linux containers on Windows (Docker Desktop)
docker run -d --name webapp -p 80:80 nginx
```

### Configuration Management

#### Environment Variables

**Linux**:
```bash
# Set environment variables
export APP_ENV=production
export DATABASE_URL=postgresql://user:pass@localhost/db

# Persistent environment variables
echo 'export APP_ENV=production' >> ~/.bashrc
echo 'export DATABASE_URL=postgresql://user:pass@localhost/db' >> ~/.bashrc

# System-wide environment variables
sudo echo 'APP_ENV=production' >> /etc/environment
```

**Windows**:
```powershell
# Set environment variables (session)
$env:APP_ENV = "production"
$env:DATABASE_URL = "postgresql://user:pass@localhost/db"

# Set environment variables (persistent)
[Environment]::SetEnvironmentVariable("APP_ENV", "production", "User")
[Environment]::SetEnvironmentVariable("DATABASE_URL", "postgresql://user:pass@localhost/db", "User")

# System-wide environment variables
[Environment]::SetEnvironmentVariable("APP_ENV", "production", "Machine")
```

### Automation Scripts

#### Cross-Platform Deployment Script

```bash
#!/bin/bash
# deploy.sh - Cross-platform deployment script

# Detect operating system
detect_os() {
    case "$(uname -s)" in
        Linux*)     OS=Linux;;
        Darwin*)    OS=Mac;;
        CYGWIN*)    OS=Cygwin;;
        MINGW*)     OS=MinGw;;
        *)          OS="UNKNOWN:$(uname -s)"
    esac
    echo "Detected OS: $OS"
}

# Install dependencies based on OS
install_dependencies() {
    case $OS in
        Linux)
            if command -v apt-get &> /dev/null; then
                sudo apt-get update
                sudo apt-get install -y docker.io docker-compose
            elif command -v yum &> /dev/null; then
                sudo yum install -y docker docker-compose
            fi
            ;;
        Mac)
            if command -v brew &> /dev/null; then
                brew install docker docker-compose
            fi
            ;;
        *)
            echo "Unsupported OS: $OS"
            exit 1
            ;;
    esac
}

# Platform-specific service management
manage_service() {
    local action=$1
    local service=$2
    
    case $OS in
        Linux)
            sudo systemctl $action $service
            ;;
        Mac)
            brew services $action $service
            ;;
    esac
}

# Main deployment function
deploy() {
    detect_os
    install_dependencies
    
    # Clone repository
    git clone https://github.com/company/app.git
    cd app
    
    # Build and deploy
    docker-compose build
    docker-compose up -d
    
    # Start services
    manage_service start docker
    
    echo "Deployment completed on $OS"
}

# Execute deployment
deploy
```

## Security and Hardening

### Linux Security

#### System Hardening

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Configure automatic security updates
sudo apt install unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades

# Disable unnecessary services
sudo systemctl disable telnet
sudo systemctl disable rsh
sudo systemctl disable rlogin

# Configure SSH security
sudo sed -i 's/#PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
sudo sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
sudo systemctl restart sshd

# Configure firewall
sudo ufw enable
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# File system security
sudo chmod 700 /root
sudo chmod 644 /etc/passwd
sudo chmod 600 /etc/shadow

# Install and configure fail2ban
sudo apt install fail2ban
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

#### Security Monitoring

```bash
#!/bin/bash
# security-audit.sh

# Check for failed login attempts
echo "=== Failed Login Attempts ==="
grep "Failed password" /var/log/auth.log | tail -10

# Check for sudo usage
echo "\n=== Recent Sudo Usage ==="
grep "sudo" /var/log/auth.log | tail -10

# Check for unusual network connections
echo "\n=== Network Connections ==="
netstat -tuln | grep LISTEN

# Check for world-writable files
echo "\n=== World-Writable Files ==="
find / -type f -perm -002 2>/dev/null | head -10

# Check for SUID files
echo "\n=== SUID Files ==="
find / -type f -perm -4000 2>/dev/null

# Check system integrity
echo "\n=== System File Changes ==="
if command -v aide &> /dev/null; then
    sudo aide --check
else
    echo "AIDE not installed"
fi

# Check for rootkits
echo "\n=== Rootkit Check ==="
if command -v rkhunter &> /dev/null; then
    sudo rkhunter --check --skip-keypress
else
    echo "rkhunter not installed"
fi
```

### Windows Security

#### PowerShell Security Script

```powershell
# windows-security-audit.ps1

# Check Windows Updates
Write-Host "=== Windows Update Status ===" -ForegroundColor Green
Get-WUList | Select-Object Title, Size, RebootRequired

# Check running services
Write-Host "\n=== Running Services ===" -ForegroundColor Green
Get-Service | Where-Object {$_.Status -eq "Running"} | Select-Object Name, DisplayName

# Check firewall status
Write-Host "\n=== Firewall Status ===" -ForegroundColor Green
Get-NetFirewallProfile | Select-Object Name, Enabled

# Check user accounts
Write-Host "\n=== User Accounts ===" -ForegroundColor Green
Get-LocalUser | Select-Object Name, Enabled, LastLogon

# Check scheduled tasks
Write-Host "\n=== Scheduled Tasks ===" -ForegroundColor Green
Get-ScheduledTask | Where-Object {$_.State -eq "Running"} | Select-Object TaskName, State

# Check event logs for security events
Write-Host "\n=== Recent Security Events ===" -ForegroundColor Green
Get-WinEvent -FilterHashtable @{LogName='Security'; Level=2} -MaxEvents 10 | 
    Select-Object TimeCreated, Id, LevelDisplayName, Message

# Check network connections
Write-Host "\n=== Network Connections ===" -ForegroundColor Green
Get-NetTCPConnection | Where-Object {$_.State -eq "Listen"} | 
    Select-Object LocalAddress, LocalPort, OwningProcess
```

## Performance Optimization

### Linux Performance Tuning

#### System Optimization Script

```bash
#!/bin/bash
# linux-performance-tuning.sh

# Kernel parameters optimization
optimize_kernel() {
    echo "Optimizing kernel parameters..."
    
    cat >> /etc/sysctl.conf <<EOF
# Network optimization
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216
net.ipv4.tcp_rmem = 4096 87380 16777216
net.ipv4.tcp_wmem = 4096 65536 16777216
net.ipv4.tcp_congestion_control = bbr

# File system optimization
fs.file-max = 2097152
vm.swappiness = 10
vm.dirty_ratio = 15
vm.dirty_background_ratio = 5

# Security optimization
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1
net.ipv4.icmp_echo_ignore_broadcasts = 1
EOF
    
    sysctl -p
}

# I/O scheduler optimization
optimize_io() {
    echo "Optimizing I/O scheduler..."
    
    # Set deadline scheduler for SSDs
    for disk in /sys/block/sd*; do
        if [[ -f "$disk/queue/scheduler" ]]; then
            echo deadline > "$disk/queue/scheduler"
        fi
    done
    
    # Optimize read-ahead
    for disk in /dev/sd*; do
        if [[ -b "$disk" ]]; then
            blockdev --setra 256 "$disk"
        fi
    done
}

# Memory optimization
optimize_memory() {
    echo "Optimizing memory settings..."
    
    # Configure huge pages
    echo 1024 > /proc/sys/vm/nr_hugepages
    
    # Optimize memory allocation
    echo 0 > /proc/sys/vm/zone_reclaim_mode
}

# CPU optimization
optimize_cpu() {
    echo "Optimizing CPU settings..."
    
    # Set CPU governor to performance
    for cpu in /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor; do
        if [[ -f "$cpu" ]]; then
            echo performance > "$cpu"
        fi
    done
}

# Main optimization function
main() {
    if [[ $EUID -ne 0 ]]; then
        echo "This script must be run as root"
        exit 1
    fi
    
    optimize_kernel
    optimize_io
    optimize_memory
    optimize_cpu
    
    echo "Performance optimization completed"
    echo "Reboot recommended to apply all changes"
}

main
```

### Windows Performance Tuning

```powershell
# windows-performance-tuning.ps1

# Optimize Windows services
function Optimize-Services {
    Write-Host "Optimizing Windows services..." -ForegroundColor Green
    
    # Disable unnecessary services
    $servicesToDisable = @(
        "Fax",
        "Spooler",
        "Themes",
        "TabletInputService"
    )
    
    foreach ($service in $servicesToDisable) {
        try {
            Set-Service -Name $service -StartupType Disabled -ErrorAction SilentlyContinue
            Write-Host "Disabled service: $service" -ForegroundColor Yellow
        }
        catch {
            Write-Host "Could not disable service: $service" -ForegroundColor Red
        }
    }
}

# Optimize power settings
function Optimize-Power {
    Write-Host "Optimizing power settings..." -ForegroundColor Green
    
    # Set high performance power plan
    powercfg /setactive 8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c
    
    # Disable USB selective suspend
    powercfg /setacvalueindex scheme_current 2a737441-1930-4402-8d77-b2bebba308a3 48e6b7a6-50f5-4782-a5d4-53bb8f07e226 0
    powercfg /setdcvalueindex scheme_current 2a737441-1930-4402-8d77-b2bebba308a3 48e6b7a6-50f5-4782-a5d4-53bb8f07e226 0
}

# Optimize network settings
function Optimize-Network {
    Write-Host "Optimizing network settings..." -ForegroundColor Green
    
    # Disable network throttling
    New-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Multimedia\SystemProfile" -Name "NetworkThrottlingIndex" -Value 0xffffffff -PropertyType DWord -Force
    
    # Optimize TCP settings
    netsh int tcp set global autotuninglevel=normal
    netsh int tcp set global chimney=enabled
    netsh int tcp set global rss=enabled
}

# Clean temporary files
function Clean-TempFiles {
    Write-Host "Cleaning temporary files..." -ForegroundColor Green
    
    $tempPaths = @(
        "$env:TEMP\*",
        "$env:WINDIR\Temp\*",
        "$env:WINDIR\Prefetch\*"
    )
    
    foreach ($path in $tempPaths) {
        try {
            Remove-Item -Path $path -Recurse -Force -ErrorAction SilentlyContinue
            Write-Host "Cleaned: $path" -ForegroundColor Yellow
        }
        catch {
            Write-Host "Could not clean: $path" -ForegroundColor Red
        }
    }
}

# Main optimization function
function Start-Optimization {
    if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
        Write-Host "This script must be run as Administrator" -ForegroundColor Red
        exit 1
    }
    
    Optimize-Services
    Optimize-Power
    Optimize-Network
    Clean-TempFiles
    
    Write-Host "Performance optimization completed" -ForegroundColor Green
    Write-Host "Restart recommended to apply all changes" -ForegroundColor Yellow
}

# Execute optimization
Start-Optimization
```

## Practical Exercises

### Exercise 1: System Administration
1. Set up a Linux server with proper user management
2. Configure SSH key-based authentication
3. Implement log rotation and monitoring
4. Create automated backup scripts

### Exercise 2: Cross-Platform Deployment
1. Create deployment scripts for both Linux and Windows
2. Implement environment-specific configurations
3. Set up monitoring and alerting
4. Test disaster recovery procedures

### Exercise 3: Security Hardening
1. Implement security best practices on both platforms
2. Set up intrusion detection systems
3. Configure automated security updates
4. Create security audit scripts

## Next Steps

After mastering operating systems:

1. **Read Part 8**: [Networking Fundamentals](./08-networking-fundamentals.md)
2. **Practice**: Set up and manage different operating systems
3. **Automate**: Create cross-platform automation scripts
4. **Secure**: Implement security hardening measures

## Key Takeaways

- Operating systems are the foundation of DevOps infrastructure
- Linux dominates server environments, Windows is common in enterprise
- Shell scripting and PowerShell are essential automation tools
- Security hardening is crucial for production systems
- Performance optimization requires understanding system internals
- Cross-platform considerations are important in mixed environments
- Monitoring and logging are essential for system health
- Automation reduces manual errors and improves consistency

---

**Continue to**: [Part 8: Networking Fundamentals](./08-networking-fundamentals.md)