package logic

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bongnv/rpkgmgr/model"
	"github.com/bongnv/rpkgmgr/repository"
)

const (
	bufferSize    = 1000
	filePathTempl = "%s%s_%s.tar.gz"

	descriptionFile = "DESCRIPTION"
	timeFormat      = "2006-01-02 15:04:05 MST"
)

var (
	titleRegex       = regexp.MustCompile("(?s)Title:(.*?)\n\\w")
	authorsRegex     = regexp.MustCompile("(?s)Author:(.*?)\n\\w")
	maintainerRegex  = regexp.MustCompile("(?s)Maintainer:(.*?)\n\\w")
	descRegex        = regexp.MustCompile("(?s)Description:(.*?)\n\\w")
	publishDateRegex = regexp.MustCompile("Publication:(.*)")
)

// NewDescriptionIndexer creates a new indexer to index each package.
func NewDescriptionIndexer(rootURL string, concurrency int, repo repository.Repository) DescriptionIndexer {
	indexer := &queueDescriptionIndexer{
		pkgCh:   make(chan *model.Package, bufferSize),
		repo:    repo,
		rootURL: rootURL,
	}

	for i := 0; i < concurrency; i++ {
		go indexer.consume()
	}

	return indexer
}

// queueDescriptionIndexer implements a queue with multiple workers to download and parse DESCRIPTION
type queueDescriptionIndexer struct {
	pkgCh   chan *model.Package
	repo    repository.Repository
	rootURL string
}

func (q *queueDescriptionIndexer) IndexPkg(pkg *model.Package) {
	q.pkgCh <- pkg
}

func (q *queueDescriptionIndexer) Shutdown() {
	close(q.pkgCh)
}

func (q *queueDescriptionIndexer) consume() {
	for pkg := range q.pkgCh {
		desc, err := q.downloadDesc(pkg)
		if err != nil {
			log.Println("Failed to download. We should retry. Err: ", err)
			continue
		}

		if err := q.parseDesc(desc, pkg); err != nil {
			log.Println("Failed to parse. Err:", err)
		}

		if err := q.update(pkg); err != nil {
			log.Println("Failed to update pkg info. Err: ", err)
			continue
		}
	}
}

func (q *queueDescriptionIndexer) parseDesc(desc string, pkg *model.Package) error {
	pkg.Title = getContent(titleRegex, desc)
	pkg.Authors = getContent(authorsRegex, desc)
	pkg.Maintainers = getContent(maintainerRegex, desc)
	pkg.Description = getContent(descRegex, desc)
	publishDate := getContent(publishDateRegex, desc)
	t, _ := time.Parse(timeFormat, publishDate)
	if !t.IsZero() {
		pkg.PublicationDate = &t
	}
	return nil
}

func (q *queueDescriptionIndexer) downloadDesc(pkg *model.Package) (string, error) {
	downloadURL := fmt.Sprintf(filePathTempl,
		q.rootURL,
		pkg.Name,
		pkg.Version,
	)

	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	gzf, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", err
	}

	tarReader := tar.NewReader(gzf)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		name := header.Name
		switch header.Typeflag {
		case tar.TypeReg:
			if !strings.HasSuffix(name, descriptionFile) {
				continue
			}

			data := make([]byte, header.Size)
			buf := bytes.NewBuffer(data)
			io.Copy(buf, tarReader)

			return buf.String(), nil
		default:
			continue
		}
	}

	return "", errors.New("file not found")
}

func (q *queueDescriptionIndexer) update(pkg *model.Package) error {
	return q.repo.Update(pkg)
}
