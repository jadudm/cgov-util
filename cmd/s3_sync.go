/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gov.gsa.fac.cgov-util/internal/pipes"
	"gov.gsa.fac.cgov-util/internal/structs"
	"gov.gsa.fac.cgov-util/internal/util"
)

var (
	source_s3 string
	dest_s3   string
)

// s3SyncCmd represents the s3Sync command
var s3SyncCmd = &cobra.Command{
	Use:   "s3_sync",
	Short: "Syncs two buckets together",
	Long:  `Uses aws s3 sync to sync two buckets contents.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3_sync called")
		util.UnsetProxy()
		source_s3 := parseS3Path(source_s3)
		dest_s3 := parseS3Path(dest_s3)
		source_creds := getBucketCredentials(source_s3)
		dest_creds := getBucketCredentials(dest_s3)

		ch := structs.Choice{
			Local: func() {
				pipes.S3Sync(source_creds, dest_creds)
			},
			Remote: func() {
				pipes.S3Sync(source_creds, dest_creds)
			}}
		runLocalOrRemote(ch)
	},
}

func init() {
	rootCmd.AddCommand(s3SyncCmd)
	s3SyncCmd.Flags().StringVarP(&source_s3, "source_s3", "", "", "Source Bucket")
	s3SyncCmd.Flags().StringVarP(&dest_s3, "dest_s3", "", "", "Destination Bucket")

	//s3SyncCmd.PersistentFlags().String("source_s3", "", "Source Bucket. (s3://fac-private-s3)")
	//s3SyncCmd.PersistentFlags().String("dest_s3", "", "Destination Bucket. (s3://backups)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3SyncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// s3SyncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
