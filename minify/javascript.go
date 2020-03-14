package minify

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// var tempDir = filepath.Join(os.TempDir(), `goresminpack`)
//
// func init() {
// 	os.MkdirAll(tempDir, 0700)
// }

// CompiledJS returns a javascript reader from compiled content reader.
// Hacky because github.com/evanw/esbuild is hacky.
func CompiledJS(file, output string, debug bool) error {
	// f, err := ioutil.TempFile(tempDir, `*.js`)
	// if err != nil {
	// 	return nil, err
	// }
	// f.Close()
	// output := fmt.Sprintf(`%s.esnext`, f.Name())

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel() // interrupts the process running
	args := []string{
		filepath.Base(file),
		`--outfile=` + output,
		`--bundle`,
	}
	if !debug {
		args = append(args, `--minify`)
	}
	p := exec.CommandContext(ctx, `esbuild`, args...)

	if debug {
		p.Stderr = os.Stderr
	}

	p.Dir = filepath.Dir(file)
	return p.Run()
	// if err != nil {
	// 	return nil, err
	// }

	// // log.Println(output)
	// o, err := os.Open(output)
	// if err != nil {
	// 	return nil, err
	// }
	// return o, nil
}
