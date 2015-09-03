/*
passthrough allows you to pass files through an executable program.
The program is executed for every single file, If you want to pass a series
of files through a single invocation of the program, please use slurp.Concat
and pipe it to passtrhough to be processed by your designed program.
*/
package passthrough

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/omeid/gonzo"
	"github.com/omeid/gonzo/context"
)

// bin is the binary name, it will be passed to os/exec.Command, so the same
// path rules applies.
// the args are the argumetns passed to the program.
func Run(bin string, args ...string) gonzo.Stage {
	return func(ctx context.Context, in <-chan gonzo.File, out chan<- gonzo.File) error {

		for {
			select {
			case file, ok := <-in:
				if !ok {
					return nil
				}

				cmd := exec.Command(bin, args...)
				cmd.Stderr = os.Stderr //TODO: io.Writer logger.
				cmd.Stdin = file

				ctx = context.WithValue(ctx, "cmd", bin)
				ctx.Infof("Passing %s", file.FileInfo().Name())

				output, err := cmd.Output()
				if err != nil {
					return err
				}

				content := ioutil.NopCloser(bytes.NewReader(output))
				out <- gonzo.NewFile(content, file.FileInfo())

			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
