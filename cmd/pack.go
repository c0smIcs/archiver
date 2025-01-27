package cmd

import (
	"errors"
	"io"
	"kemov/lib/compression"
	"kemov/lib/compression/vlc"
	"kemov/lib/compression/vlc/table/shannon_fano"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "Упаковать файл",
	Run:   pack,
}

const packedExtension = "vlc"

var ErrEmptyPath = errors.New("путь к файлу не указан")

func pack(cmd *cobra.Command, args []string) {
	if len(args) == 0 || args[0] == "" {
		handleError(ErrEmptyPath)
	}
	
	method := cmd.Flag("method").Value.String()
	var encoder compression.Encoder
	switch method {
	case "vlc":
		encoder = vlc.New(shannon_fano.NewGenerator())
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

	packed := encoder.Encode(string(data))

	err = os.WriteFile(packedFileName(filePath), packed, 0644)
	if err != nil {
		handleError(err)
	}
}

func packedFileName(path string) string {
	fileName := filepath.Base(path)

	return strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "." + packedExtension
}

func init() {
	rootCmd.AddCommand(packCmd)

	packCmd.Flags().StringP("method", "m", "", "compression method: vlc")

	if err := packCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
}