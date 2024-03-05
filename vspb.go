package vspb

import (
	"fmt"
	"os"
)

func Run(confPath string) error {
	conf, err := ReadConfig(confPath)
	if err != nil {
		return err
	}

	if err := conf.Check(); err != nil {
		return fmt.Errorf("check conf error: %w", err)
	}

	// create root and repo path
	rootPath := conf.General.Path
	repoPath := rootPath + "/repos"
	PluginPath := rootPath + "/plugins"

	if _, err := os.Stat(repoPath); err != nil {
		if err := os.MkdirAll(repoPath, 0644); err != nil {
			return fmt.Errorf("mkdir '%s' error: %w", repoPath, err)
		}
	}

	if _, err := os.Stat(PluginPath); err != nil {
		if err := os.MkdirAll(PluginPath, 0644); err != nil {
			return fmt.Errorf("mkdir '%s' error: %w", PluginPath, err)
		}
	}

	// TODO:
	// Maybe the db file should located at some specified place.
	vc, err := NewVersionControl("vspb.db")
	if err != nil {
		return fmt.Errorf("read version control file error: %s", err)
	}

	for _, pkg := range conf.Packages {
		fmt.Printf("build: %s, version: %s\n", pkg.Name, pkg.Version)
		installed, err := vc.GetPackage(pkg.Name)
		if err != nil {
			fmt.Println("get version control info error: ", err)
		}

		if installed.ID > 0 && installed.Version == pkg.Version && !installed.Failed {
			fmt.Printf("already installed, skip")
			continue
		}
		// TODO:
		// update the package if the version is different

		info := &PkgInfo{
			Name:    pkg.Name,
			Version: pkg.Version,
		}

		if err := GetPackage(rootPath, pkg); err != nil {
			fmt.Printf("download failed: %s", err)
			info.Failed = true
			if err := vc.CreatePkgInfo(info); err != nil {
				fmt.Println("save version control info error: ", err)
			}
			continue
		}

		pkgPath := repoPath + "/" + pkg.Name
		builder := &CmdRunner{
			Env:     pkg.Env,
			Dir:     pkgPath,
			Cmds:    pkg.Run,
			Verbose: conf.General.Verbose,
		}
		if len(builder.Cmds) == 0 {
			makeTool, err := MatchMakeTool(pkgPath)
			if err != nil {
				fmt.Println("check make tool failed: ", err)
				if err := vc.CreatePkgInfo(info); err != nil {
					fmt.Println("save version control info error: ", err)
				}
				continue
			}

			builder.Cmds = DefaultMaker(makeTool)
		}

		if err := builder.Run(); err != nil {
			fmt.Printf("build failed: %s", err)
			info.Failed = true
			if err := vc.CreatePkgInfo(info); err != nil {
				fmt.Println("save version control info error: ", err)
			}
			continue
		}

		installer := &Installer{
			FromDir:   pkgPath + "/build",
			ToDir:     PluginPath,
			Filenames: pkg.Provide,
		}

		if err := installer.Install(); err != nil {
			fmt.Printf("install failed: %s", err)
			info.Failed = true
			if err := vc.CreatePkgInfo(info); err != nil {
				fmt.Println("save version control info error: ", err)
			}
			continue
		}

		if err := vc.CreatePkgInfo(info); err != nil {
			fmt.Println("save version control info error: ", err)
		}

		fmt.Printf("success: build: %s, version: %s\n", pkg.Name, pkg.Version)
	}

	return nil
}
