package util

import (
	"strings"

	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// ellipsis is the suffix appended to truncated strings.
const ellipsis = " [...]"

// PartitionMessage splits a string into chunks of at most chunkSize runes.
// It searches the last distance runes for a whitespace to improve readability,
// adding chunks until reaching maxCount or maxTotal runes, returning the chunks
// and the number of omitted runes.
func PartitionMessage(
	input string,
	limits types.MessageLimit,
	distance int,
) ([]types.MessageItem, int) {
	items := make([]types.MessageItem, 0, limits.ChunkCount-1)
	runes := []rune(input)
	chunkOffset := 0
	maxTotal := Min(len(runes), limits.TotalChunkSize)
	maxCount := limits.ChunkCount - 1

	if len(input) == 0 {
		// If the message is empty, return an empty array
		return items, 0
	}

	for range maxCount {
		// If no suitable split point is found, use the chunkSize
		chunkEnd := chunkOffset + limits.ChunkSize
		// ... and start next chunk directly after this one
		nextChunkStart := chunkEnd

		if chunkEnd >= maxTotal {
			// The chunk is smaller than the limit, no need to search
			chunkEnd = maxTotal
			nextChunkStart = maxTotal
		} else {
			for r := range distance {
				rp := chunkEnd - r
				if runes[rp] == '\n' || runes[rp] == ' ' {
					// Suitable split point found
					chunkEnd = rp
					// Since the split is on a whitespace, skip it in the next chunk
					nextChunkStart = chunkEnd + 1

					break
				}
			}
		}

		items = append(items, types.MessageItem{
			Text: string(runes[chunkOffset:chunkEnd]),
		})

		chunkOffset = nextChunkStart
		if chunkOffset >= maxTotal {
			break
		}
	}

	return items, len(runes) - chunkOffset
}

// Ellipsis truncates a string to maxLength characters, appending an ellipsis if needed.
func Ellipsis(text string, maxLength int) string {
	if len(text) > maxLength {
		text = text[:maxLength-len(ellipsis)] + ellipsis
	}

	return text
}

// MessageItemsFromLines creates MessageItem batches compatible with the given limits.
func MessageItemsFromLines(plain string, limits types.MessageLimit) [][]types.MessageItem {
	maxCount := limits.ChunkCount
	lines := strings.Split(plain, "\n")
	batches := make([][]types.MessageItem, 0)
	items := make([]types.MessageItem, 0, Min(maxCount, len(lines)))

	totalLength := 0

	for _, line := range lines {
		maxLen := limits.ChunkSize

		if len(items) == maxCount || totalLength+maxLen > limits.TotalChunkSize {
			batches = append(batches, items)
			items = items[:0]
		}

		runes := []rune(line)
		if len(runes) > maxLen {
			// Trim and add ellipsis
			runes = runes[:maxLen-len(ellipsis)]
			line = string(runes) + ellipsis
		}

		if len(runes) < 1 {
			continue
		}

		items = append(items, types.MessageItem{
			Text: line,
		})

		totalLength += len(runes)
	}

	if len(items) > 0 {
		batches = append(batches, items)
	}

	return batches
}
