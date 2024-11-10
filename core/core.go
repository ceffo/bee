package core

import (
	"os"

	"github.com/charmbracelet/log"

	"ceffo.com/bee/logging"
	"ceffo.com/bee/wordsource"
	"ceffo.com/bee/wordsource/reader"
)

type Core struct {
	cleanupFuncs []func()
	source       wordsource.Maker
}

type Option func(*Core) error

func New(opts ...Option) (*Core, error) {
	core := &Core{}
	for _, opt := range opts {
		err := opt(core)
		if err != nil {
			// cleanup any resources that were created during previous options
			core.Close()
			return nil, err
		}
	}
	return core, nil
}

func (c *Core) Source() wordsource.Maker {
	return c.source
}

func (c *Core) Close() {
	// cleanup from last to first
	for i := len(c.cleanupFuncs) - 1; i >= 0; i-- {
		c.cleanupFuncs[i]()
	}
	c.cleanupFuncs = nil
}

func WithFileLogging(logFileName string) Option {
	return func(c *Core) error {
		cleanupFunc, err := logging.SetupFileLogging(logFileName)
		if err != nil {
			return err
		}
		c.cleanupFuncs = append(c.cleanupFuncs, cleanupFunc)
		return nil
	}
}

func WithStdoutLogging() Option {
	return func(c *Core) error {
		cleanupFunc, err := logging.SetupStdoutLogging()
		if err != nil {
			return err
		}
		c.cleanupFuncs = append(c.cleanupFuncs, cleanupFunc)
		return nil
	}
}

func WithSourceMaker(wordListFileName string) Option {
	return func(c *Core) error {
		fileCleanup, sourceMaker, err := sourceMaker(wordListFileName)
		if err != nil {
			return err
		}
		c.cleanupFuncs = append(c.cleanupFuncs, fileCleanup)
		c.source = sourceMaker
		return nil
	}
}

func sourceMaker(wordlistFileName string) (func(), wordsource.Maker, error) {
	log.Infof("Opening wordlist file %s", wordlistFileName)
	wordFile, err := os.Open(wordlistFileName)
	if err != nil {
		log.Errorf("Error opening wordlist file: %v", err)
		return nil, nil, err
	}
	cleanup := func() {
		log.Infof("Closing wordlist file %s", wordlistFileName)
		wordFile.Close()
	}
	maker := func() wordsource.Source {
		log.Infof("Creating new word source from %s", wordlistFileName)
		_, err := wordFile.Seek(0, 0)
		if err != nil {
			log.Errorf("Error seeking to start of file: %v", err)
			return nil
		}
		return reader.NewReaderSource(wordFile)
	}
	return cleanup, maker, nil
}
