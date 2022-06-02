package output_test

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/stefanoj3/dirstalk/pkg/common/test"
	"github.com/stefanoj3/dirstalk/pkg/scan"
	"github.com/stefanoj3/dirstalk/pkg/scan/output"
	"github.com/stretchr/testify/assert"
)

func TestFileSaverShouldErrWhenInvalidPath(t *testing.T) {
	saver, err := output.NewFileSaver("/root/123/bla.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create file")

	err = saver.Save(scan.Result{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "writeCloser is nil")

	err = saver.Close()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "writeCloser is nil")
}

func TestFileSaverShouldWriteResults(t *testing.T) {
	filename := test.RandStringRunes(10)
	filename = "testdata/" + filename + ".txt"

	defer func() {
		err := os.Remove(filename)
		if err != nil {
			t.Fatalf("%s failed to clean up file created during tests: %s", err, filename)
		}
	}()

	saver, err := output.NewFileSaver(filename)
	assert.NoError(t, err)

	err = saver.Save(scan.Result{})
	assert.NoError(t, err)

	err = saver.Close()
	assert.NoError(t, err)

	//nolint:gosec
	file, err := os.Open(filename)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(file)
	assert.NoError(t, err)

	assert.NoError(t, file.Close())

	expected := `{"Target":{"Path":"","Method":"","Depth":0},"StatusCode":0,"URL":{"Scheme":"","Opaque":"","User":null,"Host":"","Path":"","RawPath":"","ForceQuery":false,"RawQuery":"","Fragment":"","RawFragment":""},"ContentLength":0}
`
	assert.Equal(
		t,
		expected,
		string(b),
	)
}

func TestFileSaverShouldWorkConcurrently(t *testing.T) {
	filename := test.RandStringRunes(10)
	filename = "testdata/" + filename + ".txt"

	defer func() {
		err := os.Remove(filename)
		if err != nil {
			t.Fatalf("%s failed to clean up file created during tests: %s", err, filename)
		}
	}()

	saver, err := output.NewFileSaver(filename)
	assert.NoError(t, err)

	wg := sync.WaitGroup{}

	const workers = 1000

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			err := saver.Save(scan.Result{})
			if err != nil {
				panic(err)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	// checking that the file is there
	//nolint:gosec
	file, err := os.Open(filename)
	assert.NoError(t, err)
	assert.NoError(t, file.Close())
}
