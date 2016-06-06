package graph

import (
    "os"
    "github.com/Callidon/joseki/core"
)

type Graph interface {
    LoadFromFile(file *os.File)
    Add(triple core.Triple)
    Serialize(format string) string
}
