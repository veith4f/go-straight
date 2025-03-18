package project

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/veith4f/go-straight/pkg/assets"
)

var (
	CHOICE_YES *regexp.Regexp = regexp.MustCompile(`[yY](es)?`)
	CHOICE_NO  *regexp.Regexp = regexp.MustCompile(`[nN]o?`)
)

func (p *Project) runCmd(cmdName string, cmdArgs ...string) error {
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = p.vars.ProjectDir

	// Set up the standard output and error streams
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and wait for it to complete
	return cmd.Run()
}

func (p *Project) writeFiles() error {

	// Create directories
	for dir := range p.directories {
		dirPath := filepath.Join(p.vars.ProjectDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("Error creating directory %s: %w", dirPath, err)
		}
	}

	// Process files and render templates
	for file, tmplContent := range p.files {
		tmpl, err := template.New(file).Parse(tmplContent)
		if err != nil {
			return fmt.Errorf("Error parsing template for %s: %w", file, err)
		}

		ioFile, err := os.Create(filepath.Join(p.vars.ProjectDir, file))
		if err != nil {
			return fmt.Errorf("Error opening file for writing %s: %w", file, err)
		}

		err = tmpl.Execute(ioFile, p.vars)
		ioFile.Close() //nolint:errcheck
		if err != nil {
			return fmt.Errorf("Error rendering template for %s: %w", file, err)
		}
	}

	return nil
}

func (p *Project) gitInit() error {

	if err := p.runCmd("git", "init"); err != nil {
		return fmt.Errorf("Error initializing Git repository: %w", err)
	}

	return nil
}

func (p *Project) gitAddCommit() error {

	if err := p.runCmd("git", "add", "-A"); err != nil {
		return fmt.Errorf("Error adding files to Git repository: %w", err)
	}

	if err := p.runCmd("git", "commit", "-m", "Initial commit"); err != nil {
		return fmt.Errorf("Error adding files to Git repository: %w", err)
	}

	return nil
}

func (p *Project) dockerDev() error {

	if err := p.runCmd("make", "docker-dev"); err != nil {
		return fmt.Errorf("Error building dev container: %w", err)
	}

	return nil
}

func (p *Project) makeLint() error {

	if err := p.runCmd("make", "lint"); err != nil {
		return fmt.Errorf("Error linting: %w", err)
	}
	return nil
}

func (p *Project) askAddRemotePush() error {
	reader := bufio.NewReader(os.Stdin)

	var remoteHandled bool = false

	for !remoteHandled {
		// Ask the user if they want to add a remote repository
		fmt.Print("Now would be a good time to create and add an empty remote repository. Do you want to add one? (y/n) ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))

		if CHOICE_YES.MatchString(choice) {
			for {
				// Ask for the remote repository URL
				fmt.Print("Enter the remote repository URL: ")
				repoURL, _ := reader.ReadString('\n')
				repoURL = strings.TrimSpace(repoURL)

				// Check if the remote repository exists
				if p.runCmd("git", "ls-remote", repoURL) == nil {
					// Add the remote repository
					if err := p.runCmd("git", "remote", "add", "origin", repoURL); err != nil {
						return fmt.Errorf("Error adding remote: %w", err)
					}

					// Push to the repository and set upstream
					fmt.Println("Remote repository added. Pushing...")
					if err := p.runCmd("git", "push", "--set-upstream", "origin", "main"); err != nil {
						return fmt.Errorf("Error pushing: %w", err)
					}
					fmt.Println("Pushed.")
					remoteHandled = true
					break
				} else {
					// Repository doesn't exist or is inaccessible
					fmt.Print("Error: The repository does not exist or is inaccessible. Try again? (y/n) ")
					retry, _ := reader.ReadString('\n')
					retry = strings.TrimSpace(retry)

					if CHOICE_NO.MatchString(retry) {
						remoteHandled = true
						break
					}
				}
			}

		} else if CHOICE_NO.MatchString(choice) {
			break
		}
	}
	return nil
}

func (p *Project) projectDir() error {
	return os.MkdirAll(p.vars.ProjectDir, 0755)
}

func (p *Project) GoStraight() error {
	if err := p.projectDir(); err != nil {
		return err
	}
	if err := p.gitInit(); err != nil {
		return err
	}
	if err := p.writeFiles(); err != nil {
		return err
	}
	if err := p.dockerDev(); err != nil {
		return err
	}
	if err := p.makeLint(); err != nil {
		return err
	}
	if err := p.gitAddCommit(); err != nil {
		return err
	}
	if err := p.askAddRemotePush(); err != nil {
		return err
	}
	return nil
}

func NewProject(projectDir string, moduleName string, author string) (*Project, error) {

	project := &Project{}
	project.vars.ProjectDir = projectDir
	project.vars.ProjectName = filepath.Base(projectDir)
	project.vars.CopyrightYear = strconv.Itoa(time.Now().Year())
	project.vars.Author = author
	project.vars.ModuleName = moduleName
	project.directories = make(map[string]bool)
	project.files = make(map[string]string)

	for _, assetName := range assets.AssetNames() {

		dir := filepath.Dir(assetName)
		if dir != "." {
			project.directories[dir] = true
		}

		assetBytes, err := assets.Asset(assetName)
		if err != nil {
			return nil, fmt.Errorf("Error reading asset: %s", err)
		}
		project.files[assetName] = string(assetBytes)
	}

	return project, nil
}
