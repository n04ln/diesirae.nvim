package aoj

import (
	"fmt"
	"os"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/require"
)

var cookie string

func TestSubmitAndStatus(t *testing.T) {
	problemId := "ITP1_1_A"
	language := "C"
	sourceCode := "#include <stdio.h> \nint main(){\nprintf(\"Hello World\\n\");\nreturn 0;\n}"
	fmt.Println(sourceCode)

	id := os.Getenv("AOJ_ID")
	pass := os.Getenv("AOJ_RAWPASSWORD")

	cookie, err := Session(id, pass)
	require.NoError(t, err)
	pp.Println(cookie)

	token, err := Submit(cookie, problemId, language, sourceCode)
	require.NoError(t, err)
	pp.Println(token)

	res, err := Status(cookie, token)
	require.NoError(t, err)
	pp.Println(res)
}
