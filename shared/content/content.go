package content

import (
	"embed"
	"os"

	"github.com/Pigeon-Developer/pigeon-oj-tool/shared"
)

//go:embed static
var Static embed.FS

func ExtractStatic() {
	os.RemoveAll(shared.LocalPath + "/static")
	os.CopyFS(shared.LocalPath, Static)
}
