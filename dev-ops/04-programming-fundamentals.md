# Part 4: Programming Fundamentals for DevOps

## Introduction

While DevOps engineers don't need to be expert software developers, having solid programming fundamentals is crucial for automation, scripting, and understanding the applications they deploy and manage. This section covers essential programming concepts and languages relevant to DevOps.

## Why Programming Matters in DevOps

### Automation Requirements
- **Infrastructure as Code**: Terraform, CloudFormation, ARM templates
- **Configuration Management**: Ansible, Puppet, Chef
- **CI/CD Pipelines**: Jenkins, GitLab CI, GitHub Actions
- **Monitoring and Alerting**: Custom scripts and integrations
- **Deployment Scripts**: Automated deployment processes

### Problem-Solving Skills
- Debug application issues
- Understand system interactions
- Create custom tools and utilities
- Integrate different systems
- Optimize performance

## Essential Programming Languages for DevOps

### 1. Python

**Why Python for DevOps?**
- Simple and readable syntax
- Extensive library ecosystem
- Strong automation capabilities
- Cross-platform compatibility
- Large community support

**Common Use Cases**:
- Automation scripts
- API integrations
- Data processing
- Infrastructure management
- Testing frameworks

**Basic Python Example**:
```python
#!/usr/bin/env python3
import requests
import json
import sys

def check_service_health(url):
    """Check if a service is healthy"""
    try:
        response = requests.get(f"{url}/health", timeout=5)
        if response.status_code == 200:
            print(f"✅ Service at {url} is healthy")
            return True
        else:
            print(f"❌ Service at {url} returned {response.status_code}")
            return False
    except requests.exceptions.RequestException as e:
        print(f"❌ Failed to connect to {url}: {e}")
        return False

def main():
    services = [
        "http://api.example.com",
        "http://web.example.com",
        "http://db.example.com"
    ]
    
    healthy_services = 0
    for service in services:
        if check_service_health(service):
            healthy_services += 1
    
    print(f"\n{healthy_services}/{len(services)} services are healthy")
    
    if healthy_services < len(services):
        sys.exit(1)

if __name__ == "__main__":
    main()
```

### 2. Bash/Shell Scripting

**Why Bash for DevOps?**
- Native to Unix/Linux systems
- Direct system interaction
- Simple automation tasks
- Pipeline and command chaining
- System administration

**Common Use Cases**:
- System maintenance scripts
- Log processing
- File operations
- Environment setup
- Quick automation tasks

**Basic Bash Example**:
```bash
#!/bin/bash

# Deployment script with error handling
set -euo pipefail  # Exit on error, undefined vars, pipe failures

# Configuration
APP_NAME="myapp"
DEPLOY_DIR="/opt/${APP_NAME}"
BACKUP_DIR="/opt/backups"
LOG_FILE="/var/log/${APP_NAME}-deploy.log"

# Logging function
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# Error handling
error_exit() {
    log "ERROR: $1"
    exit 1
}

# Backup current version
backup_current() {
    if [ -d "$DEPLOY_DIR" ]; then
        local backup_name="${APP_NAME}-$(date +%Y%m%d-%H%M%S)"
        log "Creating backup: $backup_name"
        cp -r "$DEPLOY_DIR" "${BACKUP_DIR}/${backup_name}" || error_exit "Backup failed"
    fi
}

# Deploy new version
deploy() {
    local version="$1"
    log "Deploying version: $version"
    
    # Download and extract
    wget "https://releases.example.com/${APP_NAME}-${version}.tar.gz" -O "/tmp/${APP_NAME}-${version}.tar.gz" || error_exit "Download failed"
    
    # Stop service
    systemctl stop "$APP_NAME" || error_exit "Failed to stop service"
    
    # Deploy
    rm -rf "$DEPLOY_DIR"
    mkdir -p "$DEPLOY_DIR"
    tar -xzf "/tmp/${APP_NAME}-${version}.tar.gz" -C "$DEPLOY_DIR" || error_exit "Extraction failed"
    
    # Set permissions
    chown -R app:app "$DEPLOY_DIR"
    chmod +x "${DEPLOY_DIR}/bin/${APP_NAME}"
    
    # Start service
    systemctl start "$APP_NAME" || error_exit "Failed to start service"
    
    # Verify deployment
    sleep 5
    if systemctl is-active --quiet "$APP_NAME"; then
        log "Deployment successful"
    else
        error_exit "Service failed to start"
    fi
}

# Main execution
main() {
    if [ $# -ne 1 ]; then
        echo "Usage: $0 <version>"
        exit 1
    fi
    
    local version="$1"
    log "Starting deployment of $APP_NAME version $version"
    
    backup_current
    deploy "$version"
    
    log "Deployment completed successfully"
}

main "$@"
```

### 3. PowerShell (Windows)

**Why PowerShell for DevOps?**
- Native Windows automation
- Object-oriented pipeline
- .NET integration
- Cross-platform (PowerShell Core)
- Strong system management capabilities

**Common Use Cases**:
- Windows system administration
- Azure resource management
- Active Directory operations
- IIS management
- Windows service management

**Basic PowerShell Example**:
```powershell
# Service monitoring and restart script
param(
    [Parameter(Mandatory=$true)]
    [string[]]$ServiceNames,
    
    [int]$CheckIntervalSeconds = 60,
    [string]$LogPath = "C:\Logs\ServiceMonitor.log"
)

function Write-Log {
    param([string]$Message)
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $logEntry = "[$timestamp] $Message"
    Write-Host $logEntry
    Add-Content -Path $LogPath -Value $logEntry
}

function Test-ServiceHealth {
    param([string]$ServiceName)
    
    try {
        $service = Get-Service -Name $ServiceName -ErrorAction Stop
        
        if ($service.Status -eq 'Running') {
            Write-Log "✅ Service '$ServiceName' is running"
            return $true
        } else {
            Write-Log "❌ Service '$ServiceName' is $($service.Status)"
            return $false
        }
    }
    catch {
        Write-Log "❌ Service '$ServiceName' not found: $($_.Exception.Message)"
        return $false
    }
}

function Restart-ServiceSafely {
    param([string]$ServiceName)
    
    try {
        Write-Log "🔄 Attempting to restart service '$ServiceName'"
        Restart-Service -Name $ServiceName -Force -ErrorAction Stop
        Start-Sleep -Seconds 10
        
        if (Test-ServiceHealth -ServiceName $ServiceName) {
            Write-Log "✅ Service '$ServiceName' restarted successfully"
            return $true
        } else {
            Write-Log "❌ Service '$ServiceName' failed to start after restart"
            return $false
        }
    }
    catch {
        Write-Log "❌ Failed to restart service '$ServiceName': $($_.Exception.Message)"
        return $false
    }
}

# Main monitoring loop
Write-Log "Starting service monitoring for: $($ServiceNames -join ', ')"

while ($true) {
    foreach ($serviceName in $ServiceNames) {
        if (-not (Test-ServiceHealth -ServiceName $serviceName)) {
            Restart-ServiceSafely -ServiceName $serviceName
        }
    }
    
    Write-Log "Waiting $CheckIntervalSeconds seconds before next check..."
    Start-Sleep -Seconds $CheckIntervalSeconds
}
```

### 4. Go (Golang)

**Why Go for DevOps?**
- Fast compilation and execution
- Static binaries (easy deployment)
- Excellent concurrency support
- Strong standard library
- Popular for infrastructure tools

**Common Use Cases**:
- CLI tools and utilities
- Microservices
- Container orchestration
- Network programming
- System monitoring

**Basic Go Example**:
```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "sync"
    "time"
)

type HealthCheck struct {
    URL      string `json:"url"`
    Status   string `json:"status"`
    Response int    `json:"response_time_ms"`
    Error    string `json:"error,omitempty"`
}

type HealthChecker struct {
    URLs    []string
    Timeout time.Duration
    Results []HealthCheck
    mutex   sync.Mutex
}

func NewHealthChecker(urls []string, timeout time.Duration) *HealthChecker {
    return &HealthChecker{
        URLs:    urls,
        Timeout: timeout,
        Results: make([]HealthCheck, 0, len(urls)),
    }
}

func (hc *HealthChecker) CheckURL(url string, wg *sync.WaitGroup) {
    defer wg.Done()
    
    start := time.Now()
    client := &http.Client{Timeout: hc.Timeout}
    
    resp, err := client.Get(url + "/health")
    duration := time.Since(start)
    
    result := HealthCheck{
        URL:      url,
        Response: int(duration.Milliseconds()),
    }
    
    if err != nil {
        result.Status = "error"
        result.Error = err.Error()
    } else {
        defer resp.Body.Close()
        if resp.StatusCode == 200 {
            result.Status = "healthy"
        } else {
            result.Status = "unhealthy"
            result.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
        }
    }
    
    hc.mutex.Lock()
    hc.Results = append(hc.Results, result)
    hc.mutex.Unlock()
}

func (hc *HealthChecker) CheckAll() {
    var wg sync.WaitGroup
    hc.Results = hc.Results[:0] // Clear previous results
    
    for _, url := range hc.URLs {
        wg.Add(1)
        go hc.CheckURL(url, &wg)
    }
    
    wg.Wait()
}

func (hc *HealthChecker) PrintResults() {
    healthy := 0
    for _, result := range hc.Results {
        status := "❌"
        if result.Status == "healthy" {
            status = "✅"
            healthy++
        }
        
        fmt.Printf("%s %s (%dms)", status, result.URL, result.Response)
        if result.Error != "" {
            fmt.Printf(" - %s", result.Error)
        }
        fmt.Println()
    }
    
    fmt.Printf("\n%d/%d services are healthy\n", healthy, len(hc.Results))
    
    if healthy < len(hc.Results) {
        os.Exit(1)
    }
}

func main() {
    urls := []string{
        "http://api.example.com",
        "http://web.example.com",
        "http://db.example.com",
    }
    
    checker := NewHealthChecker(urls, 5*time.Second)
    
    fmt.Println("Checking service health...")
    checker.CheckAll()
    checker.PrintResults()
}
```

### 5. YAML

**Why YAML for DevOps?**
- Human-readable configuration format
- Widely used in DevOps tools
- Supports complex data structures
- Comments and documentation

**Common Use Cases**:
- Kubernetes manifests
- Docker Compose files
- CI/CD pipeline definitions
- Ansible playbooks
- Configuration files

**YAML Examples**:

**Docker Compose**:
```yaml
version: '3.8'

services:
  web:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./html:/usr/share/nginx/html:ro
    depends_on:
      - api
    restart: unless-stopped
    
  api:
    build:
      context: ./api
      dockerfile: Dockerfile
    environment:
      - DATABASE_URL=postgresql://user:pass@db:5432/myapp
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis
    restart: unless-stopped
    
  db:
    image: postgres:13
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    
  redis:
    image: redis:alpine
    restart: unless-stopped

volumes:
  postgres_data:

networks:
  default:
    driver: bridge
```

**Kubernetes Deployment**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
  labels:
    app: web-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web-app
  template:
    metadata:
      labels:
        app: web-app
    spec:
      containers:
      - name: web
        image: myapp:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: web-app-service
spec:
  selector:
    app: web-app
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

## Programming Concepts for DevOps

### 1. Error Handling

**Python Example**:
```python
import logging
from typing import Optional

def deploy_application(version: str) -> bool:
    """Deploy application with proper error handling"""
    try:
        # Download artifact
        download_artifact(version)
        
        # Stop current service
        stop_service()
        
        # Deploy new version
        deploy_version(version)
        
        # Start service
        start_service()
        
        # Verify deployment
        if not verify_deployment():
            raise DeploymentError("Deployment verification failed")
            
        logging.info(f"Successfully deployed version {version}")
        return True
        
    except (DownloadError, ServiceError, DeploymentError) as e:
        logging.error(f"Deployment failed: {e}")
        # Attempt rollback
        try:
            rollback_deployment()
        except Exception as rollback_error:
            logging.critical(f"Rollback failed: {rollback_error}")
        return False
    except Exception as e:
        logging.critical(f"Unexpected error during deployment: {e}")
        return False
```

### 2. Configuration Management

**Python with Environment Variables**:
```python
import os
from dataclasses import dataclass
from typing import Optional

@dataclass
class Config:
    database_url: str
    redis_url: str
    api_key: str
    debug: bool = False
    log_level: str = "INFO"
    
    @classmethod
    def from_env(cls) -> 'Config':
        return cls(
            database_url=os.getenv('DATABASE_URL', 'sqlite:///app.db'),
            redis_url=os.getenv('REDIS_URL', 'redis://localhost:6379'),
            api_key=os.getenv('API_KEY') or raise_missing_env('API_KEY'),
            debug=os.getenv('DEBUG', 'false').lower() == 'true',
            log_level=os.getenv('LOG_LEVEL', 'INFO')
        )

def raise_missing_env(var_name: str):
    raise ValueError(f"Required environment variable {var_name} is not set")

# Usage
config = Config.from_env()
```

### 3. Logging and Monitoring

**Structured Logging Example**:
```python
import logging
import json
from datetime import datetime
from typing import Dict, Any

class StructuredLogger:
    def __init__(self, name: str):
        self.logger = logging.getLogger(name)
        handler = logging.StreamHandler()
        handler.setFormatter(self.JsonFormatter())
        self.logger.addHandler(handler)
        self.logger.setLevel(logging.INFO)
    
    class JsonFormatter(logging.Formatter):
        def format(self, record):
            log_entry = {
                'timestamp': datetime.utcnow().isoformat(),
                'level': record.levelname,
                'message': record.getMessage(),
                'module': record.module,
                'function': record.funcName,
                'line': record.lineno
            }
            
            # Add extra fields if present
            if hasattr(record, 'extra_fields'):
                log_entry.update(record.extra_fields)
                
            return json.dumps(log_entry)
    
    def info(self, message: str, **kwargs):
        extra = {'extra_fields': kwargs} if kwargs else {}
        self.logger.info(message, extra=extra)
    
    def error(self, message: str, **kwargs):
        extra = {'extra_fields': kwargs} if kwargs else {}
        self.logger.error(message, extra=extra)

# Usage
logger = StructuredLogger('deployment')
logger.info("Starting deployment", version="1.2.3", environment="production")
logger.error("Deployment failed", error="Connection timeout", retry_count=3)
```

### 4. API Integration

**REST API Client Example**:
```python
import requests
from typing import Dict, Any, Optional
from dataclasses import dataclass

@dataclass
class APIResponse:
    status_code: int
    data: Optional[Dict[Any, Any]]
    error: Optional[str]

class KubernetesAPI:
    def __init__(self, base_url: str, token: str):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.session.headers.update({
            'Authorization': f'Bearer {token}',
            'Content-Type': 'application/json'
        })
    
    def _request(self, method: str, endpoint: str, **kwargs) -> APIResponse:
        url = f"{self.base_url}{endpoint}"
        try:
            response = self.session.request(method, url, **kwargs)
            
            if response.status_code < 400:
                return APIResponse(
                    status_code=response.status_code,
                    data=response.json() if response.content else None,
                    error=None
                )
            else:
                return APIResponse(
                    status_code=response.status_code,
                    data=None,
                    error=response.text
                )
        except requests.exceptions.RequestException as e:
            return APIResponse(
                status_code=0,
                data=None,
                error=str(e)
            )
    
    def get_pods(self, namespace: str = 'default') -> APIResponse:
        return self._request('GET', f'/api/v1/namespaces/{namespace}/pods')
    
    def create_deployment(self, namespace: str, deployment_spec: Dict) -> APIResponse:
        return self._request(
            'POST',
            f'/apis/apps/v1/namespaces/{namespace}/deployments',
            json=deployment_spec
        )
    
    def scale_deployment(self, namespace: str, name: str, replicas: int) -> APIResponse:
        patch_data = {
            'spec': {
                'replicas': replicas
            }
        }
        return self._request(
            'PATCH',
            f'/apis/apps/v1/namespaces/{namespace}/deployments/{name}',
            json=patch_data
        )
```

### 5. Testing and Validation

**Unit Testing Example**:
```python
import unittest
from unittest.mock import patch, MagicMock
from deployment_script import deploy_application, Config

class TestDeployment(unittest.TestCase):
    
    def setUp(self):
        self.config = Config(
            app_name='test-app',
            version='1.0.0',
            environment='test'
        )
    
    @patch('deployment_script.download_artifact')
    @patch('deployment_script.stop_service')
    @patch('deployment_script.deploy_version')
    @patch('deployment_script.start_service')
    @patch('deployment_script.verify_deployment')
    def test_successful_deployment(self, mock_verify, mock_start, 
                                 mock_deploy, mock_stop, mock_download):
        # Arrange
        mock_verify.return_value = True
        
        # Act
        result = deploy_application('1.0.0')
        
        # Assert
        self.assertTrue(result)
        mock_download.assert_called_once_with('1.0.0')
        mock_stop.assert_called_once()
        mock_deploy.assert_called_once_with('1.0.0')
        mock_start.assert_called_once()
        mock_verify.assert_called_once()
    
    @patch('deployment_script.download_artifact')
    @patch('deployment_script.rollback_deployment')
    def test_deployment_failure_triggers_rollback(self, mock_rollback, mock_download):
        # Arrange
        mock_download.side_effect = Exception("Download failed")
        
        # Act
        result = deploy_application('1.0.0')
        
        # Assert
        self.assertFalse(result)
        mock_rollback.assert_called_once()

if __name__ == '__main__':
    unittest.main()
```

## Development Environment Setup

### Python Development

**Virtual Environment**:
```bash
# Create virtual environment
python -m venv devops-env

# Activate (Linux/Mac)
source devops-env/bin/activate

# Activate (Windows)
devops-env\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Freeze dependencies
pip freeze > requirements.txt
```

**requirements.txt**:
```
requests==2.28.1
pyyaml==6.0
click==8.1.3
jinja2==3.1.2
psutil==5.9.2
paramiko==2.11.0
boto3==1.24.28
kubernetes==24.2.0
prometheus-client==0.14.1
```

### Go Development

**go.mod**:
```go
module devops-tools

go 1.19

require (
    github.com/spf13/cobra v1.5.0
    github.com/spf13/viper v1.12.0
    gopkg.in/yaml.v3 v3.0.1
    k8s.io/client-go v0.24.3
)
```

**Project Structure**:
```
devops-tools/
├── cmd/
│   ├── deploy/
│   │   └── main.go
│   └── monitor/
│       └── main.go
├── pkg/
│   ├── config/
│   ├── deploy/
│   └── monitor/
├── internal/
│   └── utils/
├── go.mod
├── go.sum
└── Makefile
```

## Best Practices

### Code Organization

1. **Modular Design**
   - Separate concerns
   - Reusable functions
   - Clear interfaces
   - Single responsibility

2. **Configuration Management**
   - Environment variables
   - Configuration files
   - Secrets management
   - Environment-specific configs

3. **Error Handling**
   - Graceful degradation
   - Proper logging
   - Retry mechanisms
   - Circuit breakers

### Security Considerations

1. **Secrets Management**
   - Never hardcode secrets
   - Use environment variables
   - Implement secret rotation
   - Use dedicated secret stores

2. **Input Validation**
   - Validate all inputs
   - Sanitize user data
   - Use parameterized queries
   - Implement rate limiting

3. **Access Control**
   - Principle of least privilege
   - Role-based access
   - Audit logging
   - Regular access reviews

### Performance Optimization

1. **Efficient Algorithms**
   - Choose appropriate data structures
   - Optimize time complexity
   - Consider memory usage
   - Profile and benchmark

2. **Concurrency**
   - Use async/await patterns
   - Implement proper locking
   - Avoid race conditions
   - Consider thread safety

3. **Resource Management**
   - Close connections properly
   - Implement connection pooling
   - Monitor resource usage
   - Set appropriate timeouts

## Practical Exercises

### Exercise 1: Health Check Script
Create a script that:
1. Checks multiple service endpoints
2. Reports status in JSON format
3. Exits with appropriate codes
4. Includes retry logic

### Exercise 2: Log Parser
Build a tool that:
1. Parses application logs
2. Extracts error patterns
3. Generates summary reports
4. Sends alerts for critical issues

### Exercise 3: Deployment Automation
Develop a deployment script that:
1. Downloads application artifacts
2. Performs rolling updates
3. Validates deployment success
4. Implements rollback capability

### Exercise 4: Configuration Manager
Create a utility that:
1. Manages environment configurations
2. Validates configuration schemas
3. Supports multiple environments
4. Handles secret injection

## Next Steps

After mastering programming fundamentals:

1. **Read Part 5**: [Build Automation](./05-build-automation.md)
2. **Practice**: Build automation scripts for your projects
3. **Explore**: Learn additional languages (Ruby, JavaScript)
4. **Contribute**: Contribute to open-source DevOps tools

## Key Takeaways

- Programming skills are essential for DevOps automation
- Python and Bash are the most common DevOps languages
- Focus on practical automation scenarios
- Error handling and logging are critical
- Security should be built into all scripts
- Testing ensures reliability of automation
- Modular design improves maintainability
- Configuration management enables flexibility

---

**Continue to**: [Part 5: Build Automation](./05-build-automation.md)