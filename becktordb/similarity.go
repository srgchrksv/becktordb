package becktordb

import (
	"math"
)

func CosineSimilarity(a, b Vector) float64 {
	dotProduct := 0.0
	normA := 0.0
	normB := 0.0
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
