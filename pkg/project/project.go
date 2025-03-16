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

/*
// extractTarGz extracts a .tar.gz file
func (p *Project) _extractTarGz(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close() //nolint:errcheck

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close() //nolint:errcheck

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath.Join(p.vars.ProjectDir, path), 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(filepath.Join(p.vars.ProjectDir, path))
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, tr)
			err2 := outFile.Close()
			if err != nil || err2 != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Project) downloadTools() error {
	osName := strings.ToLower(runtime.GOOS)
	arch := runtime.GOARCH
	version := "1.64.7"

	// Construct download URL
	url := fmt.Sprintf("https://github.com/golangci/golangci-lint/releases/download/v%s/golangci-lint-%s-%s-%s.tar.gz",
		version, version, osName, arch)

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error downloading %s: %w", url, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	// Create a temp file
	archiveName := "golangci-lint.tar.gz"
	outFile, err := os.Create(filepath.Join(p.vars.ProjectDir, archiveName))
	if err != nil {
		return fmt.Errorf("Error opening file for writing %s: %w", archiveName, err)
	}
	defer outFile.Close() //nolint:errcheck

	// Write to file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("Error saving file: %w", err)
	}

	// Extract the archive
	if err := p._extractTarGz(archiveName); err != nil {
		return fmt.Errorf("Error extracting archive %s: %w", archiveName, err)
	}
	_ = os.Remove(archiveName)

	if err := os.MkdirAll(filepath.Join(p.vars.ProjectDir, "bin"), 0755); err != nil {
		return fmt.Errorf("Error creating bin directory: %w", err)
	}

	extractedDirName := fmt.Sprintf("golangci-lint-%s-%s-%s", version, osName, arch)
	if err := os.Rename(extractedDirName, "bin/golangci-lint"); err != nil {
		return fmt.Errorf("Error moving moving binary: %w", err)
	}

	_ = os.RemoveAll(extractedDirName)

	return nil
}
*/

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

func (p *Project) goModInitTidy() error {

	if err := p.runCmd("go", "mod", "init", p.vars.ModuleName); err != nil {
		return fmt.Errorf("Error initializing Go module: %w", err)
	}

	if err := p.runCmd("go", "mod", "tidy"); err != nil {
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
	if err := p.goModInitTidy(); err != nil {
		return err
	}
	if err := p.dockerDev(); err != nil {
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
