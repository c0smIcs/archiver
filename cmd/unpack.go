package cmd

import (
	"io"
	"kemov/lib/compression"
	"kemov/lib/compression/vlc"
	"kemov/lib/compression/vlc/table/shannon_fano"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const unpackedExtension = "txt"

var unpackCmd = &cobra.Command{
	Use:   "unpack",
	Short: "Распаковать файл",
	Run:   unpack,
}

func unpack(cmd *cobra.Command, args []string) {
	if len(args) == 0 || args[0] == "" {
		handleError(ErrEmptyPath)
	}
	
	method := cmd.Flag("method").Value.String()
	var decoder compression.Decoder
	switch method {
	case "vlc":
		decoder = vlc.New(shannon_fano.NewGenerator())
	default:
		cmd.PrintErr("неизвестный метод")
		return
	}

	filePath := args[0]

	r, err := os.Open(filePath)
	if err != nil {
		handleError(err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		handleError(err)
	}

	packed := decoder.Decode(data)

	err = os.WriteFile(unpackedFileName(filePath), []byte(packed), 0644)
	if err != nil {
		handleError(err)
	}
}

func unpackedFileName(path string) string {
	fileName := filepath.Base(path)

	return strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "." + unpackedExtension
}

func init() {
	// archiver unpack vlc
	rootCmd.AddCommand(unpackCmd)

	unpackCmd.Flags().StringP("method", "m", "", "decompression method: vlc")

	if err := unpackCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
}
