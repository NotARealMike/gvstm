package interfaces

import (
    . "gvstm/stm"
    "gvstm/stmbench7/internal"
    "math/rand"
    "strconv"
    "strings"
)

func AddAtomicPartToBuildDateIndex(tx Transaction, index Index, part AtomicPart) {
    newSet := BEFactory.CreateLargeSet(tx)
    newSet.Add(tx, part)
    oldSet := index.PutIfAbsent(tx, part.GetBuildDate(tx), newSet)
    if oldSet != nil {
        oldSet.(LargeSet).Add(tx, part)
    }
}

func RemoveAtomicPartFromBuildDateIndex(tx Transaction, index Index, part AtomicPart) {
    oldSet := index.Get(tx, part.GetBuildDate(tx)).(LargeSet)
    oldSet.Remove(tx, part)
}

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
