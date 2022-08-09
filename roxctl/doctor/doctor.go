package doctor

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stackrox/rox/roxctl/common"
	"github.com/stackrox/rox/roxctl/common/environment"
	"github.com/stackrox/rox/roxctl/common/util"
)

func Command(cliEnvironment environment.Environment) *cobra.Command {
	var bundlePath string

	c := &cobra.Command{
		Use: "doctor",
		RunE: util.RunENoArgs(func(c *cobra.Command) error {
			if bundlePath != "" {
				extractedPath, err := os.MkdirTemp("", "roxctl_doctor_diag_bundle_*")
				if err != nil {
					return errors.Wrap(err, "cannot create a temporary directory")
				}

				cliEnvironment.Logger().InfofLn("Unzipping the bundle to %q...", extractedPath)
				if err := unzipDiagBundle(bundlePath, extractedPath); err != nil {
					return common.ErrInvalidCommandOption.Newf("provided file cannot be unzipped: %v; is it a diagnostic bundle?", err)
				}

				cliEnvironment.Logger().InfofLn("Running checks...")
				return runAllChecks(cliEnvironment, extractedPath)
			} else {
				// TODO(alexr):
				//   * Download the bundle if it is not specified.
				//   * Support other modes, e.g. querying Central, kubernetes cluster, etc.
				return common.ErrInvalidCommandOption.Newf("%q requires a diagnostic bundle to be specified with --bundle", c.UseLine())
			}
		}),
	}
	c.Flags().StringVarP(&bundlePath, "bundle", "b", "", "existing diagnostic bundle to analyze (.zip)")

	return c
}

func unzipDiagBundle(bundlePath, dest string) error {
	fileStat, fileErr := os.Stat(bundlePath)
	if bundlePath != "" && fileErr != nil || fileErr == nil && fileStat.IsDir() {
		return errors.New("'--bundle' must be a valid file")
	}

	// Open .zip file.
	r, err := zip.OpenReader(bundlePath)
	if err != nil {
		return errors.Wrapf(err, "cannot open diagnostic bundle .zip file")
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	extractAndWriteFileF := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, f.Mode()); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFileF(f)
		if err != nil {
			return err
		}
	}

	return nil
}
