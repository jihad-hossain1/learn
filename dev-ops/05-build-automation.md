# Part 5: Build Automation

## Introduction

Build automation is the process of automating the creation of a software build and the associated processes including compiling source code, packaging binaries, running tests, and deploying to production systems. It's a cornerstone of DevOps practices and essential for implementing Continuous Integration and Continuous Deployment (CI/CD).

## Why Build Automation Matters

### Benefits of Automated Builds
- **Consistency**: Same process every time, reducing human error
- **Speed**: Faster than manual processes
- **Reliability**: Repeatable and predictable outcomes
- **Traceability**: Complete audit trail of build processes
- **Quality**: Automated testing and quality gates
- **Efficiency**: Frees developers from repetitive tasks

### Problems with Manual Builds
- Human errors and inconsistencies
- Time-consuming processes
- Difficulty in reproducing issues
- Lack of standardization
- Poor visibility into build status
- Delayed feedback on code changes

## Build Automation Concepts

### Build Pipeline Stages

1. **Source Code Checkout**
   - Retrieve code from version control
   - Ensure clean working directory
   - Handle dependencies and submodules

2. **Dependency Management**
   - Download and install dependencies
   - Resolve version conflicts
   - Cache dependencies for performance

3. **Compilation/Transpilation**
   - Compile source code to executable format
   - Handle different target platforms
   - Optimize for performance

4. **Testing**
   - Unit tests
   - Integration tests
   - Code quality checks
   - Security scans

5. **Packaging**
   - Create deployable artifacts
   - Generate documentation
   - Create container images

6. **Deployment**
   - Deploy to target environments
   - Run smoke tests
   - Update configuration

### Build Triggers

**Manual Triggers**:
- Developer-initiated builds
- Release builds
- Hotfix builds

**Automatic Triggers**:
- Code commits (push triggers)
- Pull request creation
- Scheduled builds (nightly builds)
- Dependency updates

## Build Tools by Technology Stack

### Java Ecosystem

#### Maven

**Project Structure**:
```
my-app/
├── pom.xml
├── src/
│   ├── main/
│   │   ├── java/
│   │   └── resources/
│   └── test/
│       ├── java/
│       └── resources/
└── target/
```

**pom.xml Example**:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    
    <groupId>com.example</groupId>
    <artifactId>my-app</artifactId>
    <version>1.0.0-SNAPSHOT</version>
    <packaging>jar</packaging>
    
    <properties>
        <maven.compiler.source>11</maven.compiler.source>
        <maven.compiler.target>11</maven.compiler.target>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
        <junit.version>5.8.2</junit.version>
    </properties>
    
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
            <version>2.7.0</version>
        </dependency>
        
        <dependency>
            <groupId>org.junit.jupiter</groupId>
            <artifactId>junit-jupiter</artifactId>
            <version>${junit.version}</version>
            <scope>test</scope>
        </dependency>
    </dependencies>
    
    <build>
        <plugins>
            <plugin>
                <groupId>org.springframework.boot</groupId>
                <artifactId>spring-boot-maven-plugin</artifactId>
                <version>2.7.0</version>
                <executions>
                    <execution>
                        <goals>
                            <goal>repackage</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
            
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-surefire-plugin</artifactId>
                <version>3.0.0-M7</version>
                <configuration>
                    <includes>
                        <include>**/*Test.java</include>
                        <include>**/*Tests.java</include>
                    </includes>
                </configuration>
            </plugin>
            
            <plugin>
                <groupId>org.jacoco</groupId>
                <artifactId>jacoco-maven-plugin</artifactId>
                <version>0.8.7</version>
                <executions>
                    <execution>
                        <goals>
                            <goal>prepare-agent</goal>
                        </goals>
                    </execution>
                    <execution>
                        <id>report</id>
                        <phase>test</phase>
                        <goals>
                            <goal>report</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
        </plugins>
    </build>
    
    <profiles>
        <profile>
            <id>production</id>
            <properties>
                <spring.profiles.active>production</spring.profiles.active>
            </properties>
            <build>
                <plugins>
                    <plugin>
                        <groupId>com.spotify</groupId>
                        <artifactId>dockerfile-maven-plugin</artifactId>
                        <version>1.4.13</version>
                        <executions>
                            <execution>
                                <id>default</id>
                                <goals>
                                    <goal>build</goal>
                                    <goal>push</goal>
                                </goals>
                            </execution>
                        </executions>
                        <configuration>
                            <repository>myregistry/my-app</repository>
                            <tag>${project.version}</tag>
                        </configuration>
                    </plugin>
                </plugins>
            </build>
        </profile>
    </profiles>
</project>
```

**Common Maven Commands**:
```bash
# Clean and compile
mvn clean compile

# Run tests
mvn test

# Package application
mvn package

# Install to local repository
mvn install

# Deploy to remote repository
mvn deploy

# Run with specific profile
mvn clean package -Pproduction

# Skip tests
mvn package -DskipTests

# Run specific test
mvn test -Dtest=MyTestClass
```

#### Gradle

**build.gradle Example**:
```gradle
plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
    id 'io.spring.dependency-management' version '1.0.11.RELEASE'
    id 'jacoco'
    id 'com.palantir.docker' version '0.33.0'
}

group = 'com.example'
version = '1.0.0-SNAPSHOT'
sourceCompatibility = '11'

configurations {
    compileOnly {
        extendsFrom annotationProcessor
    }
}

repositories {
    mavenCentral()
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'org.springframework.boot:spring-boot-starter-data-jpa'
    
    compileOnly 'org.projectlombok:lombok'
    annotationProcessor 'org.projectlombok:lombok'
    
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
    testImplementation 'org.testcontainers:junit-jupiter'
    testImplementation 'org.testcontainers:postgresql'
}

tasks.named('test') {
    useJUnitPlatform()
    finalizedBy jacocoTestReport
}

jacocoTestReport {
    dependsOn test
    reports {
        xml.enabled true
        html.enabled true
    }
}

jar {
    enabled = false
}

bootJar {
    archiveFileName = "${project.name}.jar"
}

docker {
    name "myregistry/${project.name}:${project.version}"
    dockerfile file('Dockerfile')
    files bootJar.archiveFile
    buildArgs(['JAR_FILE': "${bootJar.archiveFileName}"])
}

task buildDocker {
    dependsOn build, docker
}

// Custom tasks
task integrationTest(type: Test) {
    useJUnitPlatform {
        includeTags 'integration'
    }
    shouldRunAfter test
}

task qualityCheck {
    dependsOn test, jacocoTestReport, integrationTest
    description 'Runs all quality checks'
}
```

**Common Gradle Commands**:
```bash
# Build project
./gradlew build

# Run tests
./gradlew test

# Run specific task
./gradlew qualityCheck

# Build Docker image
./gradlew buildDocker

# Clean build
./gradlew clean build

# Run with info logging
./gradlew build --info

# Run in parallel
./gradlew build --parallel
```

### JavaScript/Node.js Ecosystem

#### npm Scripts

**package.json Example**:
```json
{
  "name": "my-web-app",
  "version": "1.0.0",
  "description": "Example web application",
  "main": "src/index.js",
  "scripts": {
    "start": "node src/index.js",
    "dev": "nodemon src/index.js",
    "build": "webpack --mode production",
    "build:dev": "webpack --mode development",
    "test": "jest",
    "test:watch": "jest --watch",
    "test:coverage": "jest --coverage",
    "lint": "eslint src/**/*.js",
    "lint:fix": "eslint src/**/*.js --fix",
    "format": "prettier --write src/**/*.js",
    "security:audit": "npm audit",
    "security:fix": "npm audit fix",
    "docker:build": "docker build -t my-web-app:latest .",
    "docker:run": "docker run -p 3000:3000 my-web-app:latest",
    "precommit": "npm run lint && npm run test",
    "prebuild": "npm run test && npm run lint",
    "postinstall": "npm run build",
    "ci": "npm ci && npm run test:coverage && npm run build"
  },
  "dependencies": {
    "express": "^4.18.1",
    "helmet": "^5.1.0",
    "cors": "^2.8.5"
  },
  "devDependencies": {
    "jest": "^28.1.0",
    "eslint": "^8.17.0",
    "prettier": "^2.6.2",
    "nodemon": "^2.0.16",
    "webpack": "^5.72.1",
    "webpack-cli": "^4.9.2",
    "@babel/core": "^7.18.2",
    "@babel/preset-env": "^7.18.2",
    "babel-loader": "^8.2.5"
  },
  "engines": {
    "node": ">=14.0.0",
    "npm": ">=6.0.0"
  }
}
```

**webpack.config.js**:
```javascript
const path = require('path');
const webpack = require('webpack');

module.exports = (env, argv) => {
  const isProduction = argv.mode === 'production';
  
  return {
    entry: './src/index.js',
    output: {
      path: path.resolve(__dirname, 'dist'),
      filename: isProduction ? '[name].[contenthash].js' : '[name].js',
      clean: true
    },
    module: {
      rules: [
        {
          test: /\.js$/,
          exclude: /node_modules/,
          use: {
            loader: 'babel-loader',
            options: {
              presets: ['@babel/preset-env']
            }
          }
        },
        {
          test: /\.css$/,
          use: ['style-loader', 'css-loader']
        }
      ]
    },
    plugins: [
      new webpack.DefinePlugin({
        'process.env.NODE_ENV': JSON.stringify(argv.mode)
      })
    ],
    optimization: {
      splitChunks: {
        chunks: 'all'
      }
    },
    devtool: isProduction ? 'source-map' : 'eval-source-map'
  };
};
```

### Python Ecosystem

#### setuptools and pip

**setup.py Example**:
```python
from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

with open("requirements.txt", "r", encoding="utf-8") as fh:
    requirements = [line.strip() for line in fh if line.strip() and not line.startswith("#")]

setup(
    name="my-python-app",
    version="1.0.0",
    author="Your Name",
    author_email="your.email@example.com",
    description="A sample Python application",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/yourusername/my-python-app",
    packages=find_packages(where="src"),
    package_dir={"": "src"},
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
    ],
    python_requires=">=3.8",
    install_requires=requirements,
    extras_require={
        "dev": [
            "pytest>=7.0",
            "pytest-cov>=3.0",
            "black>=22.0",
            "flake8>=4.0",
            "mypy>=0.950",
        ],
        "test": [
            "pytest>=7.0",
            "pytest-cov>=3.0",
        ],
    },
    entry_points={
        "console_scripts": [
            "my-app=my_python_app.cli:main",
        ],
    },
)
```

**Makefile for Python**:
```makefile
.PHONY: help install install-dev test lint format type-check security-check build clean docker-build

# Default target
help:
	@echo "Available targets:"
	@echo "  install      - Install production dependencies"
	@echo "  install-dev  - Install development dependencies"
	@echo "  test         - Run tests with coverage"
	@echo "  lint         - Run linting checks"
	@echo "  format       - Format code with black"
	@echo "  type-check   - Run type checking with mypy"
	@echo "  security     - Run security checks"
	@echo "  build        - Build distribution packages"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"

# Installation
install:
	pip install -e .

install-dev:
	pip install -e ".[dev]"

# Testing
test:
	pytest tests/ --cov=src --cov-report=html --cov-report=term-missing

test-fast:
	pytest tests/ -x -v

# Code quality
lint:
	flake8 src tests
	black --check src tests

format:
	black src tests
	isort src tests

type-check:
	mypy src

# Security
security-check:
	bandit -r src
	safety check

# Build
build: clean
	python setup.py sdist bdist_wheel

clean:
	rm -rf build/
	rm -rf dist/
	rm -rf *.egg-info/
	rm -rf .coverage
	rm -rf htmlcov/
	find . -type d -name __pycache__ -delete
	find . -type f -name "*.pyc" -delete

# Docker
docker-build:
	docker build -t my-python-app:latest .

# CI pipeline
ci: install-dev lint type-check security-check test build
	@echo "CI pipeline completed successfully"
```

### .NET Ecosystem

**Project File (.csproj)**:
```xml
<Project Sdk="Microsoft.NET.Sdk.Web">

  <PropertyGroup>
    <TargetFramework>net6.0</TargetFramework>
    <Nullable>enable</Nullable>
    <ImplicitUsings>enable</ImplicitUsings>
    <Version>1.0.0</Version>
    <AssemblyVersion>1.0.0.0</AssemblyVersion>
    <FileVersion>1.0.0.0</FileVersion>
  </PropertyGroup>

  <ItemGroup>
    <PackageReference Include="Microsoft.EntityFrameworkCore" Version="6.0.5" />
    <PackageReference Include="Microsoft.EntityFrameworkCore.SqlServer" Version="6.0.5" />
    <PackageReference Include="Serilog.AspNetCore" Version="5.0.0" />
    <PackageReference Include="Swashbuckle.AspNetCore" Version="6.3.1" />
  </ItemGroup>

  <ItemGroup Condition="'$(Configuration)' == 'Debug'">
    <PackageReference Include="Microsoft.EntityFrameworkCore.Tools" Version="6.0.5" />
  </ItemGroup>

</Project>
```

**Build Script (PowerShell)**:
```powershell
# build.ps1
param(
    [string]$Configuration = "Release",
    [string]$Version = "1.0.0",
    [switch]$SkipTests,
    [switch]$PublishDocker
)

$ErrorActionPreference = "Stop"

Write-Host "Building .NET application..." -ForegroundColor Green

# Restore dependencies
Write-Host "Restoring NuGet packages..." -ForegroundColor Yellow
dotnet restore

if ($LASTEXITCODE -ne 0) {
    Write-Error "Package restore failed"
    exit 1
}

# Build
Write-Host "Building solution..." -ForegroundColor Yellow
dotnet build --configuration $Configuration --no-restore /p:Version=$Version

if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed"
    exit 1
}

# Run tests
if (-not $SkipTests) {
    Write-Host "Running tests..." -ForegroundColor Yellow
    dotnet test --configuration $Configuration --no-build --verbosity normal --collect:"XPlat Code Coverage"
    
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Tests failed"
        exit 1
    }
}

# Publish
Write-Host "Publishing application..." -ForegroundColor Yellow
dotnet publish --configuration $Configuration --no-build --output ./publish /p:Version=$Version

if ($LASTEXITCODE -ne 0) {
    Write-Error "Publish failed"
    exit 1
}

# Docker build
if ($PublishDocker) {
    Write-Host "Building Docker image..." -ForegroundColor Yellow
    docker build -t "myapp:$Version" -t "myapp:latest" .
    
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Docker build failed"
        exit 1
    }
}

Write-Host "Build completed successfully!" -ForegroundColor Green
```

## Make and Makefiles

### Universal Makefile Example

```makefile
# Makefile for multi-language project
.PHONY: help build test clean docker-build docker-push deploy

# Variables
APP_NAME := my-application
VERSION := $(shell git describe --tags --always --dirty)
REGISTRY := myregistry.com
NAMESPACE := myteam
IMAGE := $(REGISTRY)/$(NAMESPACE)/$(APP_NAME)

# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  test         - Run all tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-push  - Push Docker image to registry"
	@echo "  deploy       - Deploy to Kubernetes"
	@echo "  ci           - Run full CI pipeline"

# Build targets
build:
	@echo "Building $(APP_NAME) version $(VERSION)..."
	@if [ -f "package.json" ]; then \
		npm ci && npm run build; \
	elif [ -f "pom.xml" ]; then \
		mvn clean package -DskipTests; \
	elif [ -f "go.mod" ]; then \
		go build -o bin/$(APP_NAME) ./cmd/...; \
	elif [ -f "requirements.txt" ]; then \
		pip install -r requirements.txt && python setup.py build; \
	else \
		echo "No recognized build file found"; \
		exit 1; \
	fi

# Test targets
test:
	@echo "Running tests..."
	@if [ -f "package.json" ]; then \
		npm test; \
	elif [ -f "pom.xml" ]; then \
		mvn test; \
	elif [ -f "go.mod" ]; then \
		go test ./...; \
	elif [ -f "requirements.txt" ]; then \
		pytest; \
	else \
		echo "No recognized test framework found"; \
	fi

test-coverage:
	@echo "Running tests with coverage..."
	@if [ -f "package.json" ]; then \
		npm run test:coverage; \
	elif [ -f "pom.xml" ]; then \
		mvn test jacoco:report; \
	elif [ -f "go.mod" ]; then \
		go test -coverprofile=coverage.out ./...; \
	elif [ -f "requirements.txt" ]; then \
		pytest --cov=src --cov-report=html; \
	fi

# Quality checks
lint:
	@echo "Running linting..."
	@if [ -f "package.json" ]; then \
		npm run lint; \
	elif [ -f "pom.xml" ]; then \
		mvn checkstyle:check; \
	elif [ -f "go.mod" ]; then \
		golangci-lint run; \
	elif [ -f "requirements.txt" ]; then \
		flake8 src tests; \
	fi

security-scan:
	@echo "Running security scan..."
	@if [ -f "package.json" ]; then \
		npm audit; \
	elif [ -f "pom.xml" ]; then \
		mvn dependency-check:check; \
	elif [ -f "go.mod" ]; then \
		gosec ./...; \
	elif [ -f "requirements.txt" ]; then \
		bandit -r src; \
	fi

# Docker targets
docker-build:
	@echo "Building Docker image $(IMAGE):$(VERSION)..."
	docker build -t $(IMAGE):$(VERSION) -t $(IMAGE):latest .

docker-push: docker-build
	@echo "Pushing Docker image to registry..."
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

# Deployment
deploy:
	@echo "Deploying $(APP_NAME) version $(VERSION)..."
	kubectl set image deployment/$(APP_NAME) $(APP_NAME)=$(IMAGE):$(VERSION)
	kubectl rollout status deployment/$(APP_NAME)

# Clean up
clean:
	@echo "Cleaning build artifacts..."
	@if [ -d "node_modules" ]; then rm -rf node_modules; fi
	@if [ -d "target" ]; then rm -rf target; fi
	@if [ -d "bin" ]; then rm -rf bin; fi
	@if [ -d "dist" ]; then rm -rf dist; fi
	@if [ -d "build" ]; then rm -rf build; fi
	@if [ -f "coverage.out" ]; then rm coverage.out; fi
	docker system prune -f

# CI pipeline
ci: lint security-scan test-coverage build docker-build
	@echo "CI pipeline completed successfully"

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	@if [ -f "package.json" ]; then npm install; fi
	@if [ -f "requirements.txt" ]; then pip install -r requirements.txt; fi
	@if [ -f "go.mod" ]; then go mod download; fi

dev-run:
	@echo "Starting development server..."
	@if [ -f "package.json" ]; then \
		npm run dev; \
	elif [ -f "go.mod" ]; then \
		go run ./cmd/...; \
	elif [ -f "requirements.txt" ]; then \
		python src/main.py; \
	fi

# Version management
version-bump-patch:
	@echo "Bumping patch version..."
	@if [ -f "package.json" ]; then npm version patch; fi

version-bump-minor:
	@echo "Bumping minor version..."
	@if [ -f "package.json" ]; then npm version minor; fi

version-bump-major:
	@echo "Bumping major version..."
	@if [ -f "package.json" ]; then npm version major; fi
```

## Build Optimization Strategies

### Caching

**Dependency Caching**:
```dockerfile
# Dockerfile with layer caching
FROM node:16-alpine AS dependencies

WORKDIR /app

# Copy package files first (for caching)
COPY package*.json ./
RUN npm ci --only=production

# Build stage
FROM node:16-alpine AS build

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .
RUN npm run build

# Production stage
FROM node:16-alpine AS production

WORKDIR /app

# Copy dependencies from dependencies stage
COPY --from=dependencies /app/node_modules ./node_modules

# Copy built application
COPY --from=build /app/dist ./dist
COPY package*.json ./

EXPOSE 3000
CMD ["npm", "start"]
```

### Parallel Builds

**Gradle Parallel Execution**:
```gradle
# gradle.properties
org.gradle.parallel=true
org.gradle.caching=true
org.gradle.configureondemand=true
org.gradle.jvmargs=-Xmx2g -XX:MaxMetaspaceSize=512m
```

**Maven Parallel Execution**:
```bash
# Build with multiple threads
mvn clean package -T 4

# Build modules in parallel
mvn clean package -T 1C  # 1 thread per CPU core
```

### Incremental Builds

**Webpack Incremental Builds**:
```javascript
// webpack.config.js
module.exports = {
  cache: {
    type: 'filesystem',
    buildDependencies: {
      config: [__filename]
    }
  },
  optimization: {
    moduleIds: 'deterministic',
    runtimeChunk: 'single',
    splitChunks: {
      cacheGroups: {
        vendor: {
          test: /[\\/]node_modules[\\/]/,
          name: 'vendors',
          chunks: 'all'
        }
      }
    }
  }
};
```

## Build Monitoring and Reporting

### Build Metrics

**Build Time Tracking**:
```bash
#!/bin/bash
# build-with-metrics.sh

START_TIME=$(date +%s)
BUILD_ID=$(date +%Y%m%d-%H%M%S)

echo "Starting build $BUILD_ID at $(date)"

# Run build
if make build; then
    BUILD_STATUS="success"
    EXIT_CODE=0
else
    BUILD_STATUS="failure"
    EXIT_CODE=1
fi

END_TIME=$(date +%s)
BUILD_DURATION=$((END_TIME - START_TIME))

# Log metrics
echo "{
  \"build_id\": \"$BUILD_ID\",
  \"status\": \"$BUILD_STATUS\",
  \"duration_seconds\": $BUILD_DURATION,
  \"timestamp\": \"$(date -Iseconds)\",
  \"branch\": \"$(git rev-parse --abbrev-ref HEAD)\",
  \"commit\": \"$(git rev-parse HEAD)\"
}" >> build-metrics.json

echo "Build $BUILD_ID completed in ${BUILD_DURATION}s with status: $BUILD_STATUS"

exit $EXIT_CODE
```

### Build Notifications

**Slack Notification Script**:
```python
#!/usr/bin/env python3
import json
import os
import requests
import sys
from datetime import datetime

def send_build_notification(webhook_url, build_info):
    """Send build notification to Slack"""
    
    color = "good" if build_info["status"] == "success" else "danger"
    
    payload = {
        "attachments": [
            {
                "color": color,
                "title": f"Build {build_info['status'].title()}",
                "fields": [
                    {
                        "title": "Project",
                        "value": build_info.get("project", "Unknown"),
                        "short": True
                    },
                    {
                        "title": "Branch",
                        "value": build_info.get("branch", "Unknown"),
                        "short": True
                    },
                    {
                        "title": "Duration",
                        "value": f"{build_info.get('duration_seconds', 0)}s",
                        "short": True
                    },
                    {
                        "title": "Commit",
                        "value": build_info.get("commit", "Unknown")[:8],
                        "short": True
                    }
                ],
                "footer": "Build System",
                "ts": int(datetime.now().timestamp())
            }
        ]
    }
    
    response = requests.post(webhook_url, json=payload)
    response.raise_for_status()

if __name__ == "__main__":
    webhook_url = os.getenv("SLACK_WEBHOOK_URL")
    if not webhook_url:
        print("SLACK_WEBHOOK_URL environment variable not set")
        sys.exit(1)
    
    # Read build info from file or environment
    build_info = {
        "status": os.getenv("BUILD_STATUS", "unknown"),
        "project": os.getenv("PROJECT_NAME", "my-project"),
        "branch": os.getenv("GIT_BRANCH", "unknown"),
        "commit": os.getenv("GIT_COMMIT", "unknown"),
        "duration_seconds": int(os.getenv("BUILD_DURATION", "0"))
    }
    
    try:
        send_build_notification(webhook_url, build_info)
        print("Build notification sent successfully")
    except Exception as e:
        print(f"Failed to send notification: {e}")
        sys.exit(1)
```

## Best Practices

### Build Script Guidelines

1. **Idempotent Builds**
   - Same inputs produce same outputs
   - Clean state before building
   - Avoid side effects

2. **Fast Feedback**
   - Fail fast on errors
   - Run quick checks first
   - Parallel execution where possible

3. **Reproducible Builds**
   - Pin dependency versions
   - Use consistent environments
   - Document build requirements

4. **Error Handling**
   - Clear error messages
   - Proper exit codes
   - Cleanup on failure

### Security Considerations

1. **Dependency Scanning**
   - Regular vulnerability scans
   - Automated dependency updates
   - License compliance checks

2. **Secret Management**
   - Never hardcode secrets
   - Use secure secret stores
   - Rotate secrets regularly

3. **Build Environment Security**
   - Isolated build environments
   - Minimal required permissions
   - Regular security updates

## Practical Exercises

### Exercise 1: Multi-Language Build System
Create a build system that can handle:
1. Java (Maven/Gradle)
2. Node.js (npm)
3. Python (pip/setuptools)
4. Go (go build)

### Exercise 2: Build Optimization
Optimize an existing build to:
1. Reduce build time by 50%
2. Implement effective caching
3. Add parallel execution
4. Monitor build metrics

### Exercise 3: Build Pipeline
Create a complete build pipeline that:
1. Runs tests and quality checks
2. Builds and packages application
3. Creates container images
4. Deploys to staging environment

## Next Steps

After mastering build automation:

1. **Read Part 6**: [Testing Strategies](./06-testing-strategies.md)
2. **Practice**: Set up automated builds for your projects
3. **Optimize**: Improve build performance and reliability
4. **Integrate**: Connect builds with CI/CD pipelines

## Key Takeaways

- Build automation is essential for consistent, reliable software delivery
- Choose the right build tools for your technology stack
- Optimize for speed, reliability, and maintainability
- Implement proper error handling and monitoring
- Security should be integrated into build processes
- Caching and parallelization can significantly improve build times
- Make builds reproducible and environment-independent
- Monitor and measure build performance continuously

---

**Continue to**: [Part 6: Testing Strategies](./06-testing-strategies.md)