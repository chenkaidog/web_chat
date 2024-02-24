package domain

import "testing"

func TestEncodePassword(t*testing.T) {
	salt, hash := EncodePassword("12345678")
	t.Logf("salt: %s, hash: %s", salt, hash)
}

// dq9sNZ3XYYqZ0BieUefeFB6dE3t4u8oPxjbXHhsT4zE2h9IAtoCvg6G3F21meQ2e
// bb89a0fffe4b389dd5bf50d01c57f38cc7a98c0b4ea06df95ec208147dd235365478d863435a6e3679415f0900fd62e3d2e08274e8f78fc83af563524c1b72a7