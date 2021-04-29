package fluxgen

import (
	"github.com/wellington/go-libsass"
	"io"
	"log"
)

type SassCompiler struct {
	compiler libsass.Compiler
}

func createSassCompiler(src io.Reader, dst io.Writer) *SassCompiler {
	paths := libsass.IncludePaths([]string{CSSDir})
	style := libsass.OutputStyle(libsass.COMPRESSED_STYLE)

	comp, err := libsass.New(dst, src, paths, style)
	if err != nil {
		log.Fatalf("error creating libsass compiler - %v", err)
	}

	return &SassCompiler{compiler: comp}
}

func (s *SassCompiler) compileSass() {
	if err := s.compiler.Run(); err != nil {
		log.Fatalf("error compiling sass files -%v", err)
	}
}
