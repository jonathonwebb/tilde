package assets

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	
	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/jonathonwebb/tilde/internal/core"
)

type metafile struct {
	Outputs map[string]struct {
		EntryPoint string `json:"entryPoint"`
	} `json:"outputs"`
}

type asset struct {
	Path string `json:"path"`
	SRI  string `json:"sri"`
}

type metadata map[string]asset

func run(ctx context.Context, w io.Writer, cfg *core.Config) (err error) {
	log := cfg.NewLogger(w, "assets")
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	res := esbuild.Build(esbuild.BuildOptions{
		EntryPoints: []string{
			path.Join(cfg.AssetsDir, "entrypoints/**/*.ts"),
			path.Join(cfg.AssetsDir, "entrypoints/**/*.tsx"),
			path.Join(cfg.AssetsDir, "entrypoints/**/*.css"),
		},
		AssetNames:  "[name]-[hash]",
		ChunkNames:  "[name]-[hash]",
		EntryNames:  "[name]-[hash]",
		Sourcemap:   esbuild.SourceMapInline,
		Outdir:      path.Join(cfg.StaticDir, "build"),
		Bundle:      true,
		Metafile:    true,
		Write:       false,
		Format:      esbuild.FormatESModule,
		JSX:         esbuild.JSXTransform,
		JSXFactory:  "h",
		JSXFragment: "Fragment",
	})
	for _, err := range res.Errors {
		log.Error("build error", "msg", err)
	}
	for _, warning := range res.Warnings {
		log.Warn("build warning", "msg", warning)
	}
	if len(res.Errors) > 0 {
		return fmt.Errorf("build failed")
	}

	var mf metafile
	if err := json.Unmarshal([]byte(res.Metafile), &mf); err != nil {
		return err
	}

	h := sha512.New384()
	meta := metadata{}
	for path, entry := range mf.Outputs {
		if entry.EntryPoint != "" {
			for _, o := range res.OutputFiles {
				rel, err := filepath.Rel(pwd, o.Path)
				if err != nil {
					return err
				}
				if rel == path {
					if err := os.MkdirAll(filepath.Dir(o.Path), 0770); err != nil {
						return err
					}
					if err := os.WriteFile(o.Path, o.Contents, 0644); err != nil {
						return err
					}
					var a asset

					// TODO: disgusting
					rel2, err := filepath.Rel(cfg.StaticDir, path)
					if err != nil {
						return err
					}
					a.Path = filepath.Join("/public", rel2)

					h.Write(o.Contents)
					a.SRI = base64.StdEncoding.EncodeToString(h.Sum(nil))
					h.Reset()

					rel, err := filepath.Rel(filepath.Join(cfg.AssetsDir, "entrypoints"), entry.EntryPoint)
					if err != nil {
						return err
					}
					resource := filepath.Join("build", rel)
					meta[resource] = a
				}
			}
		}
	}

	metaJson, err := json.MarshalIndent(&meta, "", " ")
	if err != nil {
		return err
	}

	if err := os.WriteFile("ui/static/build/meta.json", metaJson, 0644); err != nil {
		return err
	}
	return nil
}
