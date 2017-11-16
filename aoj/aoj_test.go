package aoj

import (
	"os"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/require"
)

func TestSubmit(t *testing.T) {
	problemId := "ITP1_1_A"
	language := "C"
	sourceCode := "#include \nint main(){\n printf(\"Hello World\\n\");\n return 0;\n}"

	id := os.Getenv("AOJ_ID")
	pass := os.Getenv("AOJ_RAWPASSWORD")
	cookie, err := Session(id, pass)
	require.NoError(t, err)
	pp.Println(cookie)
	token, err := Submit(cookie, problemId, language, sourceCode)
	require.NoError(t, err)
	pp.Println(token)
}
