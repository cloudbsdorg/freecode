package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/freecode/freecode/internal/installation"
	"github.com/freecode/freecode/internal/version"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade freecode",
	Long:  `Check for and install freecode updates.`,
}

var (
	upgradeCheck   bool
	upgradeVersion string
	upgradeForce   bool
)

func init() {
	upgradeCmd.AddCommand(upgradeCheckCmd)
	upgradeCmd.AddCommand(upgradeInstallCmd)
}

var upgradeCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	RunE:  runUpgradeCheck,
}

var upgradeInstallCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a specific version",
	RunE:  runUpgradeInstall,
}

func init() {
	upgradeInstallCmd.Flags().BoolVarP(&upgradeForce, "force", "f", false, "Force reinstallation")
}

const (
	owner       = "freecode"
	repo        = "freecode"
	latestURL   = "https://api.github.com/repos/" + owner + "/" + repo + "/releases/latest"
	releasesURL = "https://api.github.com/repos/" + owner + "/" + repo + "/releases"
)

type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

type VersionInfo struct {
	Current   string
	Latest    string
	UpToDate  bool
	AssetURL  string
	Checksum  string
	AssetName string
}

func runUpgradeCheck(cmd *cobra.Command, args []string) error {
	info := version.Get()
	currentVersion := info.Version

	release, err := fetchLatestRelease()
	if err != nil {
		fmt.Printf("freecode version %s\n", currentVersion)
		fmt.Printf("Platform: %s\n", info.Platform)
		fmt.Println()
		fmt.Println("Unable to check for updates: unable to reach GitHub releases")
		return nil
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	upToDate := currentVersion == latestVersion || currentVersion > latestVersion

	fmt.Printf("freecode version %s\n", currentVersion)
	fmt.Printf("Platform: %s\n", info.Platform)
	fmt.Println()
	if upToDate {
		fmt.Printf("You are on the latest version: %s\n", latestVersion)
	} else {
		fmt.Printf("New version available: %s\n", latestVersion)
		fmt.Println()
		fmt.Println("To upgrade, run:")
		fmt.Printf("  freecode upgrade install\n")
		fmt.Printf("  freecode upgrade install %s\n", latestVersion)
	}

	return nil
}

func runUpgradeInstall(cmd *cobra.Command, args []string) error {
	info := version.Get()
	currentVersion := info.Version
	targetVersion := ""

	if len(args) > 0 {
		targetVersion = args[0]
	} else {
		release, err := fetchLatestRelease()
		if err != nil {
			return fmt.Errorf("unable to fetch latest release: %w", err)
		}
		targetVersion = strings.TrimPrefix(release.TagName, "v")
	}

	release, err := fetchRelease(targetVersion)
	if err != nil {
		return fmt.Errorf("unable to fetch release %s: %w", targetVersion, err)
	}

	platform := runtime.GOOS + "-" + runtime.GOARCH
	asset, checksumAsset, err := findAsset(release, platform)
	if err != nil {
		return fmt.Errorf("no compatible binary found for %s", platform)
	}

	fmt.Printf("Installing freecode %s...\n", targetVersion)
	fmt.Printf("Platform: %s\n", platform)
	fmt.Printf("Current version: %s\n", currentVersion)
	fmt.Printf("Downloading: %s\n", asset.Name)

	binData, err := downloadWithProgress(asset.BrowserDownloadURL, asset.Name, asset.Size)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	checksumData, err := downloadWithProgress(checksumAsset.BrowserDownloadURL, checksumAsset.Name, checksumAsset.Size)
	if err != nil {
		return fmt.Errorf("checksum download failed: %w", err)
	}

	if err := verifyChecksum(binData, checksumData, asset.Name); err != nil {
		return fmt.Errorf("checksum verification failed: %w", err)
	}
	fmt.Println("Checksum verified")

	installInfo, err := installation.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect installation: %w", err)
	}

	currentBinPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find current binary: %w", err)
	}

	backupPath := currentBinPath + ".backup"
	if _, err := os.Stat(currentBinPath); err == nil {
		fmt.Printf("Backing up current binary to %s\n", backupPath)
		if err := copyFile(currentBinPath, backupPath); err != nil {
			return fmt.Errorf("backup failed: %w", err)
		}
	}

	fmt.Printf("Installing to %s\n", installInfo.BinaryPath)
	if err := os.WriteFile(installInfo.BinaryPath, binData, 0755); err != nil {
		os.WriteFile(currentBinPath, readFile(backupPath), 0755)
		return fmt.Errorf("installation failed: %w", err)
	}

	os.Remove(backupPath)

	fmt.Println()
	fmt.Printf("Successfully upgraded from %s to %s\n", currentVersion, targetVersion)

	return nil
}

func fetchLatestRelease() (*Release, error) {
	return fetchReleaseURL(latestURL)
}

func fetchRelease(tag string) (*Release, error) {
	if tag == "latest" {
		return fetchLatestRelease()
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag)
	return fetchReleaseURL(url)
}

func fetchReleaseURL(url string) (*Release, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "freecode-cli")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("release not found")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	return &release, nil
}

func findAsset(release *Release, platform string) (*Asset, *Asset, error) {
	var binaryAsset, checksumAsset *Asset

	prefix := fmt.Sprintf("freecode-%s-%s", release.TagName, platform)

	for i := range release.Assets {
		asset := &release.Assets[i]
		if strings.HasPrefix(asset.Name, prefix) && strings.HasSuffix(asset.Name, ".tar.gz") &&
			!strings.HasSuffix(asset.Name, ".sha256") {
			binaryAsset = asset
		}
		if strings.HasSuffix(asset.Name, ".sha256") {
			checksumAsset = asset
		}
	}

	if binaryAsset == nil {
		return nil, nil, fmt.Errorf("binary asset not found for platform %s", platform)
	}
	if checksumAsset == nil {
		return nil, nil, fmt.Errorf("checksum asset not found")
	}

	return binaryAsset, checksumAsset, nil
}

func downloadWithProgress(url, filename string, size int64) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "freecode-cli")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func verifyChecksum(data []byte, checksumData []byte, filename string) error {
	checksumLine := string(bytes.TrimSpace(checksumData))

	parts := strings.Split(checksumLine, " ")
	if len(parts) >= 2 {
		checksumFilename := strings.Trim(parts[len(parts)-1], "*")
		checksumLine = parts[0]
		if !strings.Contains(checksumFilename, filename) && !strings.Contains(filename, checksumFilename) {
			pathParts := strings.Split(checksumFilename, "/")
			checksumName := pathParts[len(pathParts)-1]
			if !strings.Contains(filename, checksumName) && !strings.Contains(checksumName, filename) {
				return fmt.Errorf("checksum filename mismatch: expected %s, got %s", filename, checksumName)
			}
		}
	}

	expectedHash, err := hex.DecodeString(checksumLine)
	if err != nil {
		return fmt.Errorf("invalid checksum format: %w", err)
	}

	hash := sha256.Sum256(data)
	if !bytes.Equal(hash[:], expectedHash) {
		return fmt.Errorf("checksum mismatch")
	}

	return nil
}

func readFile(path string) []byte {
	data, _ := os.ReadFile(path)
	return data
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return dstFile.Sync()
}