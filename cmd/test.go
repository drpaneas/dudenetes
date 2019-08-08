/*Package cmd test
Copyright © 2019 Panagiotis Georgiadis <drpaneas@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// TerraformOutput struct for SUSE CaaSP skuba
type TerraformOutput struct {
	IPLoadBalancer struct {
		Sensitive bool     `json:"sensitive"`
		Type      string   `json:"type"`
		Value     []string `json:"value"`
	} `json:"ip_load_balancer"`
	IPMasters struct {
		Sensitive bool     `json:"sensitive"`
		Type      string   `json:"type"`
		Value     []string `json:"value"`
	} `json:"ip_masters"`
	IPWorkers struct {
		Sensitive bool     `json:"sensitive"`
		Type      string   `json:"type"`
		Value     []string `json:"value"`
	} `json:"ip_workers"`
}

var masters, workers, loadbalancers []string

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Executes a test scenario",
	Long: `Executes a test scenario

If '--file' is not specified it executes all the *.feature files found at the current working directory`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		if file == "" {
			// Use the current directory
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			file = dir
		}
		// Check if file or dir exists
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				log.Fatal(err)
			} else {
				log.Fatalf("Permission Denied: %s", err)
			}
		}

		// SUSE CaaSP Cluster (read the terraform output)
		skuba, _ := cmd.Flags().GetString("skuba")
		if skuba != "" {
			// Check if file or dir exists
			if _, err := os.Stat(skuba); err != nil {
				if os.IsNotExist(err) {
					log.Fatal(err)
				} else {
					log.Fatalf("Permission Denied: %s", err)
				}
			}
			terraformOutput := skuba
			terraformJSON := exec.Command("terraform", "output", "-json")
			terraformJSON.Dir = terraformOutput
			out, err := terraformJSON.Output()
			if err != nil {
				log.Fatal(err)
			}

			var dat map[string]interface{}

			if err := json.Unmarshal(out, &dat); err != nil {
				panic(err)
			}

			for key, value := range dat {
				for k, v := range value.(map[string]interface{}) {
					if key == "ip_load_balancer" && k == "value" {
						lbArray := v.([]interface{})
						for _, ip := range lbArray {
							//fmt.Println(i, ip)
							loadbalancers = append(loadbalancers, ip.(string))
						}
					}
					if key == "ip_masters" && k == "value" {
						mastersArray := v.([]interface{})
						for _, ip := range mastersArray {
							//fmt.Println(i, ip)
							masters = append(masters, ip.(string))
						}
					}
					if key == "ip_workers" && k == "value" {
						workersArray := v.([]interface{})
						for _, ip := range workersArray {
							//fmt.Println(i, ip)
							workers = append(workers, ip.(string))
						}
					}
				}
			}

			var d []string
			for _, v := range loadbalancers {
				d = append(d, fmt.Sprintf("loadbalancer=%s", v))
			}

			for i, v := range masters {
				d = append(d, fmt.Sprintf("master%d=%s", i+1, v))
			}

			for i, v := range workers {
				d = append(d, fmt.Sprintf("worker%d=%s", i+1, v))
			}

			// Write the configuration to local .env file
			f, err := os.Create(".env")
			if err != nil {
				fmt.Println(err)
				f.Close()
				return
			}
			for _, v := range d {
				fmt.Fprintln(f, v)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			err = f.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		// Logging capability
		logfile := fmt.Sprintf("%s.log", file)
		f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer f.Close()

		// Run the system command and print the output live on the screen
		mwriter := io.MultiWriter(f, os.Stdout)
		godog := exec.Command("godog", file)
		godog.Stderr = mwriter
		godog.Stdout = mwriter
		err = godog.Run() //blocks until sub process is complete
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringP("file", "f", "", "Executes the specified test file")
	testCmd.Flags().StringP("skuba", "s", "", "Uses an already running SUSE CaaSP cluster")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
