// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"io"
	"strconv"
	"encoding/csv"
	"encoding/json"
	"github.com/spf13/cobra"
)

var cfgFile string
var format string
var path string

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Home  string `json:"home"`
	Shell string `json:"shell"`
}
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hrcobra",
	Short: "The command will be able to export usernames, IDs, home directories,and shells as either JSON or CSV.",
	Long: `he command will be able to export usernames, IDs, home directories, and shells as either JSON or CSV.
				This command will not include information about system users (users with IDs under 1000).

	 By default, the command will display the information as JSON to stdout, but the -format flag will allow a person to specify csv as the export type`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		users := collectUsers()
		var output io.Writer


		if path != "" {
			f, err := os.Create(path)
			handleError(err)
			defer f.Close()

			output = f
		}else {
			output = os.Stdout
		}

		if format == "json" {
			data, err := json.MarshalIndent(users, "", "  ")
	    handleError(err)
			output.Write(data)
		}else if format == "csv" {
			output.Write([]byte("name,id,home,shell\n"))
			writer := csv.NewWriter(output)

			for _, user := range users {
				err := writer.Write([]string{user.Name, strconv.Itoa(user.Id), user.Home, user.Shell})
				handleError(err)
			}
			writer.Flush()
		}else {
			cmd.Usage()
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
	rootCmd.Flags().StringVarP(&format, "format", "f", "", "Format to output json or csv")
	rootCmd.Flags().StringVarP(&path, "path", "p", "", "Path to save csv")
	rootCmd.MarkFlagRequired("format")

}


func collectUsers() (users []User) {
	f, err := os.Open("/etc/passwd")
	handleError(err)
	defer f.Close()

	//create a csv reader

	reader := csv.NewReader(f)
	reader.Comma = ':'


	lines, err := reader.ReadAll()
	handleError(err)

	for _, line := range lines {
		id, err := strconv.ParseInt(line[2], 10, 64)
		handleError(err)

		if id < 1000 {
			continue
		}

		user := User{
			Name: line[0],
			Id: int(id),
			Home: line[5],
			Shell: line[6],
		}

		users = append(users, user)

	}
	return
}

func handleError(err error) {

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
