package tests

import (
	"testing"

	"github.com/ayush-1/logaid/internal/plugins"
)

// TestGitPlugin tests the Git plugin with comprehensive test cases
func TestGitPlugin(t *testing.T) {
	plugin := &plugins.GitPlugin{}

	testCases := []struct {
		name        string
		command     string
		output      string
		shouldMatch bool
		expectedFix string
		description string
	}{
		// Common Git typos
		{
			name:        "checkout typo",
			command:     "git checout main",
			output:      "git: 'checout' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git checkout main",
			description: "Most common git checkout typo",
		},
		{
			name:        "commit typo",
			command:     "git comit -m 'fix bug'",
			output:      "git: 'comit' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git commit -m 'fix bug'",
			description: "Git commit typo",
		},
		{
			name:        "branch typo",
			command:     "git branc -a",
			output:      "git: 'branc' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git branch -a",
			description: "Git branch typo",
		},
		{
			name:        "status typo",
			command:     "git stat",
			output:      "git: 'stat' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git status",
			description: "Git status typo",
		},
		{
			name:        "merge typo",
			command:     "git merg develop",
			output:      "git: 'merg' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git merge develop",
			description: "Git merge typo",
		},
		{
			name:        "rebase typo",
			command:     "git rebas main",
			output:      "git: 'rebas' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git rebase main",
			description: "Git rebase typo",
		},
		{
			name:        "remote typo",
			command:     "git remot -v",
			output:      "git: 'remot' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git remote -v",
			description: "Git remote typo",
		},
		{
			name:        "fetch typo",
			command:     "git fetc origin",
			output:      "git: 'fetc' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git fetch origin",
			description: "Git fetch typo",
		},
		{
			name:        "pull typo",
			command:     "git pul origin main",
			output:      "git: 'pul' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git pull origin main",
			description: "Git pull typo",
		},
		{
			name:        "push typo",
			command:     "git pus origin main",
			output:      "git: 'pus' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git push origin main",
			description: "Git push typo",
		},
		{
			name:        "add typo",
			command:     "git ad .",
			output:      "git: 'ad' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git add .",
			description: "Git add typo",
		},
		{
			name:        "reset typo",
			command:     "git rese --hard HEAD",
			output:      "git: 'rese' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git reset --hard HEAD",
			description: "Git reset typo",
		},
		{
			name:        "stash typo",
			command:     "git stas",
			output:      "git: 'stas' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git stash",
			description: "Git stash typo",
		},
		{
			name:        "log typo",
			command:     "git lo --oneline",
			output:      "git: 'lo' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git log --oneline",
			description: "Git log typo",
		},
		{
			name:        "diff typo",
			command:     "git dif HEAD~1",
			output:      "git: 'dif' is not a git command. See 'git --help'.",
			shouldMatch: true,
			expectedFix: "git diff HEAD~1",
			description: "Git diff typo",
		},

		// Branch/remote errors
		{
			name:        "branch not found",
			command:     "git checkout develop",
			output:      "error: pathspec 'develop' did not match any file(s) known to git",
			shouldMatch: true,
			expectedFix: "", // Should suggest AI fix
			description: "Branch doesn't exist",
		},
		{
			name:        "remote not found",
			command:     "git push upstream main",
			output:      "fatal: 'upstream' does not appear to be a git repository",
			shouldMatch: true,
			expectedFix: "", // Should suggest AI fix
			description: "Remote doesn't exist",
		},

		// Authentication errors
		{
			name:        "authentication failed",
			command:     "git push origin main",
			output:      "remote: Permission denied (publickey).\nfatal: Could not read from remote repository.",
			shouldMatch: true,
			expectedFix: "", // Should suggest AI fix
			description: "SSH key authentication failed",
		},
		{
			name:        "https auth failed",
			command:     "git clone https://github.com/user/repo.git",
			output:      "remote: Repository not found.\nfatal: repository 'https://github.com/user/repo.git/' not found",
			shouldMatch: true,
			expectedFix: "", // Should suggest AI fix
			description: "HTTPS authentication or repo not found",
		},

		// Merge conflicts
		{
			name:        "merge conflict",
			command:     "git merge feature-branch",
			output:      "Auto-merging file.txt\nCONFLICT (content): Merge conflict in file.txt\nAutomatic merge failed; fix conflicts and then commit the result.",
			shouldMatch: true,
			expectedFix: "", // Should suggest AI fix
			description: "Merge conflict occurred",
		},
		{
			name:        "rebase conflict",
			command:     "git rebase main",
			output:      "First, rewinding head to replay your work on top of it...\nApplying: commit message\nUsing index info to reconstruct a base tree...\nFalling back to patching base and 3-way merge...\nAuto-merging file.txt\nCONFLICT (content): Merge conflict in file.txt",
			shouldMatch: true,
			expectedFix: "", // Should suggest AI fix
			description: "Rebase conflict occurred",
		},

		// Non-matching cases
		{
			name:        "successful command",
			command:     "git status",
			output:      "On branch main\nYour branch is up to date with 'origin/main'.\n\nnothing to commit, working tree clean",
			shouldMatch: false,
			expectedFix: "",
			description: "Successful git command",
		},
		{
			name:        "non-git command",
			command:     "ls -la",
			output:      "total 48\ndrwxr-xr-x 12 user user 4096 Jan 26 10:00 .",
			shouldMatch: false,
			expectedFix: "",
			description: "Non-git command",
		},
		{
			name:        "git help output",
			command:     "git --help",
			output:      "usage: git [--version] [--help] [-C <path>] [-c <name>=<value>]",
			shouldMatch: false,
			expectedFix: "",
			description: "Git help output",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Match function
			matches := plugin.Match(tc.command, tc.output)
			if matches != tc.shouldMatch {
				t.Errorf("Match() = %v, want %v for case: %s", matches, tc.shouldMatch, tc.description)
			}

			// Test Suggest function (only if it should match)
			if tc.shouldMatch && tc.expectedFix != "" {
				suggestion := plugin.Suggest(tc.command, tc.output)
				if suggestion != tc.expectedFix {
					t.Errorf("Suggest() = %q, want %q for case: %s", suggestion, tc.expectedFix, tc.description)
				}
			}
		})
	}
}

// TestDockerPlugin tests the Docker plugin with comprehensive test cases
func TestDockerPlugin(t *testing.T) {
	plugin := &plugins.DockerPlugin{}

	testCases := []struct {
		name        string
		command     string
		output      string
		shouldMatch bool
		expectedFix string
		description string
	}{
		// Image name typos
		{
			name:        "ubuntu typo",
			command:     "docker run ubntu",
			output:      "Unable to find image 'ubntu:latest' locally",
			shouldMatch: true,
			expectedFix: "docker run ubuntu",
			description: "Ubuntu image name typo",
		},
		{
			name:        "nginx typo",
			command:     "docker run ngnix",
			output:      "Unable to find image 'ngnix:latest' locally",
			shouldMatch: true,
			expectedFix: "docker run nginx",
			description: "Nginx image name typo",
		},
		{
			name:        "alpine typo",
			command:     "docker run alpin",
			output:      "Unable to find image 'alpin:latest' locally",
			shouldMatch: true,
			expectedFix: "docker run alpine",
			description: "Alpine image name typo",
		},
		{
			name:        "redis typo",
			command:     "docker run redi",
			output:      "Unable to find image 'redi:latest' locally",
			shouldMatch: true,
			expectedFix: "docker run redis",
			description: "Redis image name typo",
		},
		{
			name:        "postgres typo",
			command:     "docker run postgre",
			output:      "Unable to find image 'postgre:latest' locally",
			shouldMatch: true,
			expectedFix: "docker run postgres",
			description: "Postgres image name typo",
		},
		{
			name:        "mysql typo",
			command:     "docker run mysq",
			output:      "Unable to find image 'mysq:latest' locally",
			shouldMatch: true,
			expectedFix: "docker run mysql",
			description: "MySQL image name typo",
		},
		{
			name:        "mongo typo",
			command:     "docker run mong",
			output:      "Unable to find image 'mong:latest' locally",
			shouldMatch: true,
			expectedFix: "docker run mongo",
			description: "MongoDB image name typo",
		},

		// Command typos
		{
			name:        "run typo",
			command:     "docker ru ubuntu",
			output:      "docker: 'ru' is not a docker command.",
			shouldMatch: true,
			expectedFix: "docker run ubuntu",
			description: "Docker run command typo",
		},
		{
			name:        "build typo",
			command:     "docker buil .",
			output:      "docker: 'buil' is not a docker command.",
			shouldMatch: true,
			expectedFix: "docker build .",
			description: "Docker build command typo",
		},
		{
			name:        "pull typo",
			command:     "docker pul nginx",
			output:      "docker: 'pul' is not a docker command.",
			shouldMatch: true,
			expectedFix: "docker pull nginx",
			description: "Docker pull command typo",
		},
		{
			name:        "push typo",
			command:     "docker pus myimage",
			output:      "docker: 'pus' is not a docker command.",
			shouldMatch: true,
			expectedFix: "docker push myimage",
			description: "Docker push command typo",
		},
		{
			name:        "exec typo",
			command:     "docker exe -it container bash",
			output:      "docker: 'exe' is not a docker command.",
			shouldMatch: true,
			expectedFix: "docker exec -it container bash",
			description: "Docker exec command typo",
		},
		{
			name:        "ps typo",
			command:     "docker p",
			output:      "docker: 'p' is not a docker command.",
			shouldMatch: true,
			expectedFix: "docker ps",
			description: "Docker ps command typo",
		},
		{
			name:        "logs typo",
			command:     "docker log container",
			output:      "docker: 'log' is not a docker command.",
			shouldMatch: true,
			expectedFix: "docker logs container",
			description: "Docker logs command typo",
		},
		{
			name:        "stop typo",
			command:     "docker stp container",
			output:      "docker: 'stp' is not a docker command.",
			shouldMatch: true,
			expectedFix: "docker stop container",
			description: "Docker stop command typo",
		},
		{
			name:        "start typo",
			command:     "docker stat container",
			output:      "docker: 'stat' is not a docker command.",
			shouldMatch: true,
			expectedFix: "docker start container",
			description: "Docker start command typo",
		},
		{
			name:        "remove typo",
			command:     "docker rm container",
			output:      "docker: 'rm' is not a docker command.",
			shouldMatch: true,
			expectedFix: "docker rm container",
			description: "Docker rm should not be corrected",
		},

		// Permission errors
		{
			name:        "permission denied",
			command:     "docker run ubuntu",
			output:      "docker: Got permission denied while trying to connect to the Docker daemon socket",
			shouldMatch: true,
			expectedFix: "sudo docker run ubuntu",
			description: "Docker permission denied",
		},
		{
			name:        "socket permission",
			command:     "docker ps",
			output:      "permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock",
			shouldMatch: true,
			expectedFix: "sudo docker ps",
			description: "Docker socket permission denied",
		},

		// Service not running
		{
			name:        "daemon not running",
			command:     "docker run ubuntu",
			output:      "Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?",
			shouldMatch: true,
			expectedFix: "", // Should suggest AI fix to start Docker
			description: "Docker daemon not running",
		},

		// Non-matching cases
		{
			name:        "successful run",
			command:     "docker run ubuntu",
			output:      "Unable to find image 'ubuntu:latest' locally\nlatest: Pulling from library/ubuntu\n5bed26d33875: Pull complete",
			shouldMatch: false,
			expectedFix: "",
			description: "Successful docker run with pull",
		},
		{
			name:        "non-docker command",
			command:     "ls -la",
			output:      "total 48\ndrwxr-xr-x 12 user user 4096 Jan 26 10:00 .",
			shouldMatch: false,
			expectedFix: "",
			description: "Non-docker command",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Match function
			matches := plugin.Match(tc.command, tc.output)
			if matches != tc.shouldMatch {
				t.Errorf("Match() = %v, want %v for case: %s", matches, tc.shouldMatch, tc.description)
			}

			// Test Suggest function (only if it should match)
			if tc.shouldMatch && tc.expectedFix != "" {
				suggestion := plugin.Suggest(tc.command, tc.output)
				if suggestion != tc.expectedFix {
					t.Errorf("Suggest() = %q, want %q for case: %s", suggestion, tc.expectedFix, tc.description)
				}
			}
		})
	}
}

// TestNpmPlugin tests the NPM plugin with comprehensive test cases
func TestNpmPlugin(t *testing.T) {
	plugin := &plugins.NpmPlugin{}

	testCases := []struct {
		name        string
		command     string
		output      string
		shouldMatch bool
		expectedFix string
		description string
	}{
		// Command typos
		{
			name:        "install typo",
			command:     "npm instal express",
			output:      "Unknown command: \"instal\"",
			shouldMatch: true,
			expectedFix: "npm install express",
			description: "NPM install command typo",
		},
		{
			name:        "install typo 2",
			command:     "npm isntall lodash",
			output:      "Unknown command: \"isntall\"",
			shouldMatch: true,
			expectedFix: "npm install lodash",
			description: "NPM install command typo variant",
		},
		{
			name:        "start typo",
			command:     "npm stat",
			output:      "Unknown command: \"stat\"",
			shouldMatch: true,
			expectedFix: "npm start",
			description: "NPM start command typo",
		},
		{
			name:        "test typo",
			command:     "npm tes",
			output:      "Unknown command: \"tes\"",
			shouldMatch: true,
			expectedFix: "npm test",
			description: "NPM test command typo",
		},
		{
			name:        "run typo",
			command:     "npm ru dev",
			output:      "Unknown command: \"ru\"",
			shouldMatch: true,
			expectedFix: "npm run dev",
			description: "NPM run command typo",
		},
		{
			name:        "update typo",
			command:     "npm updat",
			output:      "Unknown command: \"updat\"",
			shouldMatch: true,
			expectedFix: "npm update",
			description: "NPM update command typo",
		},
		{
			name:        "uninstall typo",
			command:     "npm uninsta express",
			output:      "Unknown command: \"uninsta\"",
			shouldMatch: true,
			expectedFix: "npm uninstall express",
			description: "NPM uninstall command typo",
		},

		// Package not found
		{
			name:        "package not found",
			command:     "npm install expres",
			output:      "npm ERR! 404 Not Found - GET https://registry.npmjs.org/expres - Not found",
			shouldMatch: true,
			expectedFix: "npm install express",
			description: "Package name typo - express",
		},
		{
			name:        "lodash typo",
			command:     "npm install lodas",
			output:      "npm ERR! 404 Not Found - GET https://registry.npmjs.org/lodas - Not found",
			shouldMatch: true,
			expectedFix: "npm install lodash",
			description: "Package name typo - lodash",
		},
		{
			name:        "react typo",
			command:     "npm install reac",
			output:      "npm ERR! 404 Not Found - GET https://registry.npmjs.org/reac - Not found",
			shouldMatch: true,
			expectedFix: "npm install react",
			description: "Package name typo - react",
		},
		{
			name:        "axios typo",
			command:     "npm install axio",
			output:      "npm ERR! 404 Not Found - GET https://registry.npmjs.org/axio - Not found",
			shouldMatch: true,
			expectedFix: "npm install axios",
			description: "Package name typo - axios",
		},

		// Permission errors
		{
			name:        "permission denied",
			command:     "npm install -g typescript",
			output:      "npm ERR! Error: EACCES: permission denied, access '/usr/local/lib/node_modules'",
			shouldMatch: true,
			expectedFix: "sudo npm install -g typescript",
			description: "NPM global install permission denied",
		},

		// Network errors
		{
			name:        "network error",
			command:     "npm install express",
			output:      "npm ERR! network request to https://registry.npmjs.org/express failed, reason: getaddrinfo ENOTFOUND registry.npmjs.org",
			shouldMatch: true,
			expectedFix: "", // Should suggest AI fix
			description: "Network connectivity issue",
		},

		// Non-matching cases
		{
			name:        "successful install",
			command:     "npm install express",
			output:      "added 50 packages from 37 contributors and audited 50 packages in 1.23s",
			shouldMatch: false,
			expectedFix: "",
			description: "Successful npm install",
		},
		{
			name:        "non-npm command",
			command:     "ls -la",
			output:      "total 48\ndrwxr-xr-x 12 user user 4096 Jan 26 10:00 .",
			shouldMatch: false,
			expectedFix: "",
			description: "Non-npm command",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Match function
			matches := plugin.Match(tc.command, tc.output)
			if matches != tc.shouldMatch {
				t.Errorf("Match() = %v, want %v for case: %s", matches, tc.shouldMatch, tc.description)
			}

			// Test Suggest function (only if it should match)
			if tc.shouldMatch && tc.expectedFix != "" {
				suggestion := plugin.Suggest(tc.command, tc.output)
				if suggestion != tc.expectedFix {
					t.Errorf("Suggest() = %q, want %q for case: %s", suggestion, tc.expectedFix, tc.description)
				}
			}
		})
	}
}
