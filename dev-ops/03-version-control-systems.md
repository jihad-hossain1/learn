# Part 3: Version Control Systems

## Introduction

Version Control Systems (VCS) are fundamental to DevOps practices. They enable teams to collaborate effectively, track changes, manage releases, and maintain code quality. Git has become the de facto standard for version control in modern software development.

## What is Version Control?

Version control is a system that records changes to files over time so that you can recall specific versions later. It allows multiple people to work on the same project simultaneously while maintaining a complete history of changes.

### Key Benefits
- **Collaboration**: Multiple developers can work on the same codebase
- **History**: Complete record of all changes
- **Backup**: Distributed copies of the entire project
- **Branching**: Parallel development streams
- **Merging**: Combining changes from different sources
- **Rollback**: Ability to revert to previous versions

## Types of Version Control Systems

### Centralized Version Control (CVCS)
**Examples**: SVN, Perforce, TFS

**Characteristics**:
- Single central server
- Clients check out files from central location
- All version history stored centrally

**Pros**:
- Simple to understand
- Centralized administration
- Fine-grained access control

**Cons**:
- Single point of failure
- Requires network connectivity
- Limited offline capabilities

### Distributed Version Control (DVCS)
**Examples**: Git, Mercurial, Bazaar

**Characteristics**:
- Every clone is a full backup
- No single point of failure
- Full history available locally

**Pros**:
- Works offline
- Fast operations
- Flexible workflows
- Better branching and merging

**Cons**:
- Steeper learning curve
- More complex for beginners
- Larger storage requirements

## Git Fundamentals

### What is Git?

Git is a distributed version control system created by Linus Torvalds in 2005. It's designed to handle everything from small to very large projects with speed and efficiency.

### Key Git Concepts

#### Repository (Repo)
A directory that contains your project files and the entire version history.

#### Working Directory
The current state of files you're working on.

#### Staging Area (Index)
A intermediate area where changes are prepared before committing.

#### Commit
A snapshot of your project at a specific point in time.

#### Branch
A parallel version of your repository.

#### Remote
A version of your repository hosted on a server.

### Git Workflow

```
Working Directory → Staging Area → Local Repository → Remote Repository
     (edit)          (add)         (commit)        (push)
```

## Essential Git Commands

### Repository Setup

```bash
# Initialize a new repository
git init

# Clone an existing repository
git clone <repository-url>

# Check repository status
git status

# View commit history
git log
git log --oneline
git log --graph
```

### Basic Operations

```bash
# Add files to staging area
git add <file>
git add .                # Add all files
git add -A               # Add all files including deletions

# Commit changes
git commit -m "Commit message"
git commit -am "Add and commit in one step"

# View differences
git diff                 # Working directory vs staging
git diff --staged        # Staging vs last commit
git diff HEAD           # Working directory vs last commit
```

### Remote Operations

```bash
# Add remote repository
git remote add origin <repository-url>

# View remotes
git remote -v

# Push changes
git push origin main
git push -u origin main  # Set upstream

# Pull changes
git pull origin main
git fetch origin         # Fetch without merging
```

### Branching and Merging

```bash
# Create and switch to new branch
git checkout -b <branch-name>
git switch -c <branch-name>  # Modern syntax

# Switch branches
git checkout <branch-name>
git switch <branch-name>     # Modern syntax

# List branches
git branch
git branch -r            # Remote branches
git branch -a            # All branches

# Merge branches
git merge <branch-name>

# Delete branch
git branch -d <branch-name>
git branch -D <branch-name>  # Force delete
```

### Undoing Changes

```bash
# Unstage files
git reset HEAD <file>
git restore --staged <file>  # Modern syntax

# Discard working directory changes
git checkout -- <file>
git restore <file>           # Modern syntax

# Undo last commit (keep changes)
git reset --soft HEAD~1

# Undo last commit (discard changes)
git reset --hard HEAD~1

# Revert a commit (create new commit)
git revert <commit-hash>
```

## Git Branching Strategies

### Git Flow

**Branch Types**:
- **main/master**: Production-ready code
- **develop**: Integration branch for features
- **feature/***: New features
- **release/***: Release preparation
- **hotfix/***: Critical fixes

**Workflow**:
1. Create feature branch from develop
2. Develop feature
3. Merge feature to develop
4. Create release branch from develop
5. Test and fix in release branch
6. Merge release to main and develop
7. Tag release

### GitHub Flow

**Simplified workflow**:
1. Create feature branch from main
2. Develop feature
3. Open pull request
4. Review and discuss
5. Merge to main
6. Deploy

**Benefits**:
- Simple and lightweight
- Continuous deployment friendly
- Fast feedback loops

### GitLab Flow

**Environment branches**:
- **main**: Development
- **pre-production**: Staging
- **production**: Production

**Workflow**:
1. Feature branches merge to main
2. Main promotes to pre-production
3. Pre-production promotes to production

## Advanced Git Concepts

### Rebasing

```bash
# Interactive rebase
git rebase -i HEAD~3

# Rebase onto another branch
git rebase main

# Continue after resolving conflicts
git rebase --continue

# Abort rebase
git rebase --abort
```

**When to use rebase**:
- Clean up commit history
- Integrate changes from main branch
- Squash multiple commits

**When NOT to use rebase**:
- On public/shared branches
- When commit history is important

### Cherry-picking

```bash
# Apply specific commit to current branch
git cherry-pick <commit-hash>

# Cherry-pick multiple commits
git cherry-pick <commit1> <commit2>

# Cherry-pick a range
git cherry-pick <start-commit>..<end-commit>
```

### Stashing

```bash
# Stash current changes
git stash
git stash push -m "Work in progress"

# List stashes
git stash list

# Apply stash
git stash apply
git stash apply stash@{0}

# Pop stash (apply and remove)
git stash pop

# Drop stash
git stash drop stash@{0}
```

### Tags

```bash
# Create lightweight tag
git tag v1.0.0

# Create annotated tag
git tag -a v1.0.0 -m "Version 1.0.0"

# List tags
git tag
git tag -l "v1.*"

# Push tags
git push origin v1.0.0
git push origin --tags

# Delete tag
git tag -d v1.0.0
git push origin --delete v1.0.0
```

## Git Configuration

### Global Configuration

```bash
# Set user information
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

# Set default editor
git config --global core.editor "code --wait"

# Set default branch name
git config --global init.defaultBranch main

# View configuration
git config --list
git config --global --list
```

### Useful Aliases

```bash
# Create aliases
git config --global alias.st status
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit
git config --global alias.unstage 'reset HEAD --'
git config --global alias.last 'log -1 HEAD'
git config --global alias.visual '!gitk'
```

### .gitignore File

```gitignore
# Dependencies
node_modules/
*.log

# Build outputs
dist/
build/
*.min.js

# Environment files
.env
.env.local

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Language specific
__pycache__/
*.pyc
*.class
target/
```

## Git Hosting Platforms

### GitHub

**Features**:
- Pull requests
- Issues and project management
- GitHub Actions (CI/CD)
- GitHub Pages
- Security features
- Large community

**Best for**:
- Open source projects
- Public repositories
- Integration with third-party tools

### GitLab

**Features**:
- Merge requests
- Built-in CI/CD
- Issue tracking
- Wiki and documentation
- Container registry
- Security scanning

**Best for**:
- Complete DevOps platform
- Self-hosted solutions
- Enterprise features

### Bitbucket

**Features**:
- Pull requests
- Bitbucket Pipelines (CI/CD)
- Jira integration
- Trello integration
- Smart mirroring

**Best for**:
- Atlassian ecosystem
- Small teams
- Jira integration

### Azure DevOps

**Features**:
- Pull requests
- Azure Pipelines
- Work item tracking
- Test plans
- Artifacts

**Best for**:
- Microsoft ecosystem
- Enterprise environments
- Integrated ALM

## Collaboration Workflows

### Pull Request Process

1. **Create Feature Branch**
   ```bash
   git checkout -b feature/new-feature
   ```

2. **Develop and Commit**
   ```bash
   git add .
   git commit -m "Add new feature"
   ```

3. **Push Branch**
   ```bash
   git push origin feature/new-feature
   ```

4. **Create Pull Request**
   - Use web interface
   - Add description
   - Request reviewers
   - Link to issues

5. **Code Review**
   - Review changes
   - Add comments
   - Request changes or approve

6. **Merge**
   - Squash and merge
   - Merge commit
   - Rebase and merge

### Code Review Best Practices

**For Authors**:
- Keep changes small and focused
- Write clear commit messages
- Add tests for new features
- Update documentation
- Self-review before submitting

**For Reviewers**:
- Review promptly
- Be constructive and specific
- Focus on code quality and logic
- Check for security issues
- Verify tests are adequate

**Review Checklist**:
- [ ] Code follows style guidelines
- [ ] Logic is correct and efficient
- [ ] Tests are comprehensive
- [ ] Documentation is updated
- [ ] No security vulnerabilities
- [ ] Performance considerations

## Git in DevOps Pipelines

### Triggering CI/CD

**Common Triggers**:
- Push to main branch
- Pull request creation
- Tag creation
- Scheduled builds

**Example GitHub Actions**:
```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Run tests
      run: npm test
```

### Semantic Versioning with Git

**Version Format**: MAJOR.MINOR.PATCH

**Automated Versioning**:
```bash
# Using conventional commits
git commit -m "feat: add new feature"     # Minor version
git commit -m "fix: resolve bug"          # Patch version
git commit -m "feat!: breaking change"    # Major version
```

### Git Hooks

**Pre-commit Hook Example**:
```bash
#!/bin/sh
# .git/hooks/pre-commit

# Run tests before commit
npm test
if [ $? -ne 0 ]; then
  echo "Tests failed. Commit aborted."
  exit 1
fi

# Run linting
npm run lint
if [ $? -ne 0 ]; then
  echo "Linting failed. Commit aborted."
  exit 1
fi
```

## Troubleshooting Common Issues

### Merge Conflicts

```bash
# When merge conflict occurs
git status                    # See conflicted files
# Edit files to resolve conflicts
git add <resolved-files>
git commit                   # Complete the merge
```

**Conflict Resolution**:
1. Identify conflicted files
2. Edit files to resolve conflicts
3. Remove conflict markers
4. Test the resolution
5. Add and commit

### Detached HEAD

```bash
# Create branch from detached HEAD
git checkout -b new-branch-name

# Or return to main branch
git checkout main
```

### Large File Issues

```bash
# Use Git LFS for large files
git lfs install
git lfs track "*.psd"
git add .gitattributes
git add file.psd
git commit -m "Add large file"
```

### Accidentally Committed Secrets

```bash
# Remove from history (use with caution)
git filter-branch --force --index-filter \
  'git rm --cached --ignore-unmatch path/to/secret/file' \
  --prune-empty --tag-name-filter cat -- --all

# Force push (dangerous)
git push origin --force --all
```

## Best Practices

### Commit Messages

**Conventional Commits Format**:
```
type(scope): description

[optional body]

[optional footer]
```

**Examples**:
```
feat(auth): add OAuth2 integration
fix(api): resolve null pointer exception
docs(readme): update installation instructions
test(user): add unit tests for user service
```

### Repository Structure

```
project-root/
├── .github/
│   ├── workflows/
│   └── PULL_REQUEST_TEMPLATE.md
├── docs/
├── src/
├── tests/
├── .gitignore
├── README.md
├── LICENSE
└── CONTRIBUTING.md
```

### Security Considerations

**Protect Sensitive Data**:
- Never commit secrets or credentials
- Use environment variables
- Implement pre-commit hooks
- Regular security scans

**Branch Protection**:
- Require pull request reviews
- Require status checks
- Restrict force pushes
- Require signed commits

## Practical Exercises

### Exercise 1: Basic Git Workflow
1. Initialize a new repository
2. Create and edit files
3. Stage and commit changes
4. Create branches
5. Merge branches

### Exercise 2: Collaboration Simulation
1. Fork a repository
2. Create feature branch
3. Make changes and commit
4. Create pull request
5. Review and merge

### Exercise 3: Conflict Resolution
1. Create conflicting changes
2. Attempt to merge
3. Resolve conflicts
4. Complete merge

## Next Steps

After mastering version control:

1. **Read Part 4**: [Programming Fundamentals](./04-programming-fundamentals.md)
2. **Practice**: Set up a Git repository for a project
3. **Explore**: Try different branching strategies
4. **Integrate**: Connect Git with CI/CD tools

## Key Takeaways

- Git is essential for modern DevOps practices
- Distributed version control enables flexible workflows
- Branching strategies should match team needs
- Code review improves quality and knowledge sharing
- Automation can enforce quality gates
- Security should be built into Git workflows
- Practice and consistency lead to mastery

---

**Continue to**: [Part 4: Programming Fundamentals](./04-programming-fundamentals.md)