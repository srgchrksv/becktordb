package utils

import (
	"github.com/joho/godotenv"
	"github.com/srgchrksv/becktordb/becktordb"
)

func LoadEnv(filenames ...string) error {
	if len(filenames) < 1 {
		err := godotenv.Load()
		if err != nil {
			return err
		}

	} else {
		for _, filename := range filenames {
			err := godotenv.Load(filename)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Float32ToFloat64(values []float32) becktordb.Vector {
	vector := make(becktordb.Vector, len(values))
	for i, v := range values {
		vector[i] = float64(v)
	}
	return vector
}
