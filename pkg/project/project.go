package project

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/veith4f/go-straight/pkg/assets"
)

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

func (p *Project) gitInitCommit() error {

	cmd := exec.Command("git", "init")
	cmd.Dir = p.vars.ProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error initializing Git repository: %w", err)
	}

	cmd = exec.Command("git", "add", "-A")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error adding files to Git repository: %w", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error adding files to Git repository: %w", err)
	}

	return nil
}

func (p *Project) dockerDev() error {
	cmd := exec.Command("make", "docker-dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error building dev container: %w", err)
	}

	return nil
}

func (p *Project) goModInitTidy() error {
	cmd := exec.Command("go", "mod", "init", p.vars.ModuleName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error initializing Go module: %w", err)
	}

	cmd = exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error tidying Go module: %w", err)
	}
	return nil
}

func (p *Project) askAddRemotePush() error {
	reader := bufio.NewReader(os.Stdin)

	// Ask the user if they want to add a remote repository
	fmt.Print("Now would be a good time to create and add an empty remote repository. Do you want to add one? (y/n) ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToLower(choice))

	if choice == "y" || choice == "yes" || choice == "Yes" {
		for {
			// Ask for the remote repository URL
			fmt.Print("Enter the remote repository URL: ")
			repoURL, _ := reader.ReadString('\n')
			repoURL = strings.TrimSpace(repoURL)

			// Check if the remote repository exists
			cmd := exec.Command("git", "ls-remote", repoURL)
			if err := cmd.Run(); err == nil {
				// Add the remote repository
				err := exec.Command("git", "remote", "add", "origin", repoURL).Run()
				if err == nil {
					fmt.Println("Remote repository added. Pushing...")

					// Push to the repository and set upstream
					err := exec.Command("git", "push", "--set-upstream", "origin", "main").Run()
					if err != nil {
						fmt.Println("Done.")
					} else {
						fmt.Println("Error pushing:", err)
					}
				}

				break
			} else {
				// Repository doesn't exist or is inaccessible
				fmt.Print("Error: The repository does not exist or is inaccessible. Try again? (y/n) ")
				retry, _ := reader.ReadString('\n')
				retry = strings.TrimSpace(strings.ToLower(retry))

				if retry != "y" && retry != "yes" {
					break
				}
			}
		}
	}
	return nil
}

func (p *Project) GoStraight() error {
	if err := p.writeFiles(); err != nil {
		return err
	}
	if err := p.goModInitTidy(); err != nil {
		return err
	}
	if err := p.dockerDev(); err != nil {
		return err
	}
	if err := p.gitInitCommit(); err != nil {
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
