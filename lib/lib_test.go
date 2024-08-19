package lib_test

import (
	"bufio"
	"go-git-finder/lib"
	"os"
	"testing"
)

func BenchmarkFile(b *testing.B) {
	file, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		b.Error(err)
	}
	defer file.Close()
	writer := bufio.NewWriterSize(file, 1024*1024)

	for i := 0; i < b.N; i++ {

		writer.AvailableBuffer()
	}
	writer.Flush()

}

func TestGhp(t *testing.T) {
	var x = "x-access-token:ghs_AdaGKovxmLGZHQMsftohieNXoSELBv10vtId"
	e := lib.GithubRegex.FindAllString(x, -1)
	t.Logf("%v", e)
}