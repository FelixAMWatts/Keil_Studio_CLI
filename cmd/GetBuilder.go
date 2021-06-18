/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"log"
	"net/http"
	"encoding/json"
	"github.com/spf13/cobra"
	"bufio"
	"strings"
	"text/tabwriter"
	"os"
)


// GetBuilderCmd represents the GetBuilder command
var GetBuilderCmd = &cobra.Command{
	Use:   "GetBuilder",
	Short: "Return details of specific CMSIS Builders.",
	Long: `A CMSIS Builder is a combination of a specific set of CMSIS-Build tools and a specific toolchain, which can be used to build a CMSIS project.`,
	Run: func(cmd *cobra.Command, args []string) {
		getBuilder(args[0])
	},
}


func init() {
	rootCmd.AddCommand(GetBuilderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// GetBuilderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// GetBuilderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


type Builder struct{
	BuildToolsType 		string 		`json:"buildToolsType"`
	BuildToolsVersion	string		`json:"buildToolsVersion"`
	Deprecated			bool		`json:"deprecated"`
	DepractedInfo		int			`json:"deprecationInfo"`
	Name				string		`json:"name"`
	Title				string		`json:"title"`
	ToolchainType		string		`json:"toolchainType"`
	ToolchainVersion	string		`json:"toolchainVersion"`
}


func getBuilder(Builder_Name string){
	url := "https://build.api.keil.arm.com/cmsis-builders/"
	responseString, Status_Code := getBuilderData(url+Builder_Name)
	if Status_Code == 200{
		dec := json.NewDecoder(strings.NewReader(responseString))
		var m Builder
		dec.Decode(&m)
		Output_Table(m)
	} else{
		fmt.Println("Request Unsuccesful")
	if (Status_Code == 401){
		fmt.Println("Error 401: Authentication information has not been provided or is invalid")
	}
	if (Status_Code == 404){
		fmt.Println("Error 404: Requested resource cannot be found on the server or the user lacks privileges to see this resource")
	}
	if (Status_Code == 406){
		fmt.Println("Error 406: Requested cannot be fulfilled due to content negotiation headers, such as a request for a unsupported media type or specific version")
	}
	if (Status_Code == 429){
		fmt.Println("Error 429: To many requests have been received by the user, i.e. rate limits have been reached")
	}	
	}



}

func Output_Table(m Builder){
	Builder_Table := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
		fmt.Fprintf(Builder_Table, "BuildToolsType\t%s\t\n",m.BuildToolsType,)
		fmt.Fprintf(Builder_Table, "BuildToolsVersion\t%s\t\n",m.BuildToolsVersion)
		fmt.Fprintf(Builder_Table, "Deprecated\t%v\t\n",m.Deprecated)
		if m.Deprecated{	
		fmt.Fprintf(Builder_Table, "DepractedInfo\t%v\t\n",m.DepractedInfo)
		}
		fmt.Fprintf(Builder_Table, "Name\t%s\t\n",m.Name,)
		fmt.Fprintf(Builder_Table, "Title\t%s\t\n",m.Title)
		fmt.Fprintf(Builder_Table, "ToolchainType\t%s\t\n",m.ToolchainType,)
		fmt.Fprintf(Builder_Table, "ToolchainVersion\t%s\t\n",m.ToolchainVersion)
		Builder_Table.Flush()
}


func getBuilderData(baseAPI string) (string, int) {
	var bearer = "Bearer " + "2e0017dc7f0b83f3b2984b5e96c32b139a515632"

	request, err := http.NewRequest(
		http.MethodGet,
		baseAPI,
		nil,
	)
	if err != nil{
		log.Printf("Could not request Builder Data - %v", err)
	}
	
	request.Header.Add("Accept", "applicaton/json")
	request.Header.Add("Authorization", bearer)

	response, err := http.DefaultClient.Do(request)

	if err != nil{
		log.Printf("Could not make a request - %v", err)
	}
	if err != nil{
		log.Printf("Could not read response body -%v",err)
	}
	scanner := bufio.NewScanner(response.Body)
	scanner.Scan()
	return scanner.Text(), response.StatusCode
}