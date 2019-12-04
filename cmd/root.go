/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

type Pages struct {
	Data []Page `yaml:"pages"`
}

type Page struct {
	Dynamic bool   `yaml:"dynamic"`
	Path    string `yaml:"path"`
}

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vpagen",
	Short: "Automatic page generation CLI for Gridsome.js.",
	Long:  `Please prepare two files sitemap.yml and vpagen.template.vue in the gridsome.js project root.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		buf, err := ioutil.ReadFile("./sitemap.yml")
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		var pages Pages
		err = yaml.Unmarshal(buf, &pages)

		if err != nil {
			log.Fatalf("error: %v", err)
		}


		//Indexページは必ず生成される
		w, err := os.Create("src/pages/Index.vue")
		if err != nil {
			log.Printf("warning: %v", err)
		}
		r, err := os.Open("vpagen.template.vue")
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		_, err = io.Copy(w, r)
		if err != nil {
			log.Printf("warning: %v", err)
		}
		
		//yamlファイルの配列１レコードずつ処理
		for _, page := range pages.Data {
			dirs := strings.Split(page.Path, "/")

			mkdirPath, cpFilePath := "src/pages", "src/pages"
			mkdirPathDynamic, cpFilePathDynamic := "src/templates", "src/templates"

			for i := 1; i < len(dirs); i++ {
				path := dirs[i]
				mkdirPath += "/" + path
				mkdirPathDynamic += "/" + path

				if i < len(dirs)-1 {
					cpFilePath += "/" + path
					cpFilePathDynamic += "/" + path
					if page.Dynamic {
						if err := os.Mkdir(mkdirPath, 0777); err != nil {
							log.Printf("warning: %v", err)
						}
						if err := os.Mkdir(mkdirPathDynamic, 0777); err != nil {
							log.Printf("warning: %v", err)
						}
					} else {
						if err := os.Mkdir(mkdirPath, 0777); err != nil {
							log.Printf("warning: %v", err)
						}

					}

				} else {
					// パスの最後のセパレートはファイル名としてファイル生成する
					// 動的ページとする場合はgridsome.jsのルールに従いpagesとtemplates
					vueFileName := strings.Title(strings.ToLower(path)) + ".vue"

					if page.Dynamic {
						w, err := os.Create(cpFilePath + "/" + vueFileName)
						if err != nil {
							log.Printf("warning: %v", err)
						}
						r, err := os.Open("vpagen.template.vue")
						if err != nil {
							log.Fatalf("error: %v", err)
						}

						_, err = io.Copy(w, r)
						if err != nil {
							log.Printf("warning: %v", err)
						}

						w, err = os.Create(cpFilePathDynamic + "/" + vueFileName)
						if err != nil {
							log.Printf("warning: %v", err)
						}
						r, err = os.Open("vpagen.template.vue")
						if err != nil {
							log.Fatalf("error: %v", err)
						}

						_, err = io.Copy(w, r)
						if err != nil {
							log.Printf("warning: %v", err)
						}

					} else {
						w, err := os.Create(cpFilePath + "/" + vueFileName)
						if err != nil {
							log.Printf("warning: %v", err)
						}
						r, err := os.Open("vpagen.template.vue")
						if err != nil {
							log.Fatalf("error: %v", err)
						}

						_, err = io.Copy(w, r)
						if err != nil {
							log.Printf("warning: %v", err)
						}
					}

				}

			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vpagen.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".vpagen" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".vpagen")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
