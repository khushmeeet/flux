package fluxgen

import (
	"github.com/wellington/go-libsass"
	"io"
	"log"
)

type SassCompiler struct {
	compiler libsass.Compiler
}

func createSassCompiler(src io.Reader, dst io.Writer, fc *fluxConfig) *SassCompiler {
	var comp libsass.Compiler
	paths := libsass.IncludePaths([]string{CSSDir})

	if (*fc)["minify_css"].(bool) {
		var err error
		style := libsass.OutputStyle(libsass.COMPRESSED_STYLE)
		comp, err = libsass.New(dst, src, paths, style)
		if err != nil {
			log.Fatalf("error creating libsass compiler - %v", err)
		}
	} else {
		var err error
		comp, err = libsass.New(dst, src, paths)
		if err != nil {
			log.Fatalf("error creating libsass compiler - %v", err)
		}
	}

	return &SassCompiler{compiler: comp}
}

func (s *SassCompiler) compileSass() {
	if err := s.compiler.Run(); err != nil {
		log.Fatalf("error compiling sass files -%v", err)
	}
}
