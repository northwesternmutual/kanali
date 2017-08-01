package cmd

import (
  "io"
  "os"
  "fmt"
  "bytes"
  "strings"
  "os/exec"
  "net/url"
  "net/http"
  "io/ioutil"
  "archive/zip"
  "encoding/json"
  "path/filepath"

	"github.com/spf13/cobra"
  "github.com/spf13/viper"
  "github.com/Sirupsen/logrus"
  "github.com/northwesternmutual/kanali/config"
  "github.com/northwesternmutual/kanali/plugins"
)

func init() {
  logrus.SetFormatter(&logrus.TextFormatter{})

  if err := config.PluginLoadFlags.AddAll(pluginLoadCmd); err != nil {
		logrus.Fatalf("could not add flag to command: %s", err.Error())
		os.Exit(1)
	}

	pluginCmd.AddCommand(pluginLoadCmd)
}

var pluginLoadCmd = &cobra.Command{
	Use:   `load`,
	Short: `load a Kanali plugin`,
	Long:  `load a Kanali plugin`,
	Run: func(cmd *cobra.Command, args []string) {

    data, err := ioutil.ReadFile(viper.GetString(config.FlagPluginRequirementsFileLocation.GetLong()))
  	if err != nil {
  		logrus.Fatalf("could not load requirements file: %s", err.Error())
      os.Exit(1)
  	}

    pluginsList := &plugins.PluginList{}
		if err := json.Unmarshal(data, pluginsList); err != nil {
      logrus.Fatalf("requirements file is invalid: %s", err.Error())
      os.Exit(1)
		}

    for _, plugin := range pluginsList.Plugins {
      logrus.Infof("loading %s", plugin.Location)
      u, err := plugin.GetURL()
      if err != nil {
        logrus.Errorf("could not load create url %s: %s", plugin.Location, err.Error())
        continue
      }
      zipLoc, zipName, err := downloadZIP(u)
      if err != nil {
        logrus.Errorf("could not download zip %s: %s", plugin.Location, err.Error())
        continue
      }
      defer os.Remove(zipLoc)

      pluginPath := fmt.Sprintf("/go/src/github.com/northwesternmutual/%s", strings.Split(u.Path, "/")[3])
      pluginDirLoc, err := unzipPlugin(zipLoc, pluginPath)
      if err != nil {
        logrus.Errorf("could not unzip plugin %s: %s", plugin.Location, err.Error())
        continue
      }
      defer os.RemoveAll(pluginDirLoc)
      if err := os.Chdir(filepath.Join(pluginDirLoc, zipName)); err != nil {
        logrus.Errorf("could not cd into plugin dir %s: %s", plugin.Location, err.Error())
        continue
      }
      cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", fmt.Sprintf("%s.so", plugin.Name))
      var out bytes.Buffer
	    cmd.Stderr = &out
      if err := cmd.Run(); err != nil {
        logrus.Errorf("could not compile plugin %s: %s", plugin.Location, out.String())
        continue
      }
      if err := os.Rename(fmt.Sprintf("%s/%s/%s.so", pluginDirLoc, zipName, plugin.Name), fmt.Sprintf("%s/%s.so", viper.GetString(config.FlagPluginLocation.GetLong()), plugin.Name)); err != nil {
        logrus.Errorf("could not move plugin %s: %s", plugin.Location, err.Error())
        continue
      }
    }

	},
}

func downloadZIP(u *url.URL) (string, string, error) {
  	file, err := ioutil.TempFile("", "kanali")
  	if err != nil {
  		return "", "", err
  	}
  	defer file.Close()

  	resp, err := http.Get(u.String())
  	if err != nil {
  		return "", "", err
  	}
  	defer resp.Body.Close()

  	if _, err := io.Copy(file, resp.Body); err != nil {
      return "", "", err
  	}

    return file.Name(), fmt.Sprintf("%s", strings.Split(strings.Split(resp.Header.Get("Content-Disposition"), "filename=")[1], ".zip")[0]), nil
}

func unzipPlugin(zipLoc, pluginPath string) (string, error) {
  reader, err := zip.OpenReader(zipLoc)
  if err != nil {
      return "", err
  }
  defer reader.Close()

  // destDir, err := ioutil.TempDir("", "kanali")
  if err := os.Mkdir(pluginPath, 0755); err != nil {
    return "", err
  }

  extractAndWriteFile := func(f *zip.File) error {
    rc, err := f.Open()
    if err != nil {
        return err
    }
    defer rc.Close()
    path := filepath.Join(pluginPath, f.Name)
    if f.FileInfo().IsDir() {
      os.MkdirAll(path, f.Mode())
    } else {
      os.MkdirAll(filepath.Dir(path), f.Mode())
      f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
      if err != nil {
        return err
      }
      defer f.Close()

      _, err = io.Copy(f, rc)
      if err != nil {
        return err
      }
    }
    return nil
  }

  for _, f := range reader.File {
    err := extractAndWriteFile(f)
    if err != nil {
      return "", err
    }
  }

  return pluginPath, nil
}
