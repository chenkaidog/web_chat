package domain

import "testing"

func TestEncodePassword(t*testing.T) {
	salt, hash := EncodePassword("12345678")
	t.Logf("salt: %s, hash: %s", salt, hash)
}

// TIhqOpY06bMbK3VzYQDHvyyk5NvlvqSVM0vmzdOcg6yRs36cgUmA2VepiT1UFMIr
// 5478d863435a6e3679415f0900fd62e3d2e08274e8f78fc83af563524c1b72a7