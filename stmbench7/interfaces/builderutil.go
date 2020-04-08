package interfaces

import (
    "gvstm/stmbench7/internal"
    "math/rand"
    "strconv"
    "strings"
)

func createType() string {
    // TODO: how random does this need to be? The rand package source is deterministic.
    return "type #" + strconv.Itoa(rand.Int() % internal.NumTypes)
}

func createBuildDate(minBuildDate, maxBuildDate int) int {
    // TODO: how random does this need to be? The rand package source is deterministic.
    return minBuildDate + rand.Int() % (maxBuildDate - minBuildDate + 1)
}


func createText(textSize int, textPattern string) string {
    patternSize := len(textPattern)
    size := patternSize
    var b strings.Builder
    b.Grow(textSize)
    for ; size < textSize ; size += patternSize {
        b.WriteString(textPattern)
    }
    return b.String()
}
