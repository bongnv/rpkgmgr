package logic

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bongnv/rpkgmgr/model"
	"github.com/bongnv/rpkgmgr/repository"
)

const (
	packagesFile = "PACKAGES"

	packagePrefix = "Package:"
	versionPrefix = "Version:"
	dependPrefix  = "Depends:"
	suggestPrefix = "Suggests:"
	licensePrefix = "License:"
)

var (
	errInvalidFormat = errors.New("rpkgmgr: invalid format")
)

// DescriptionIndexer is the interface to work with description indexer.
type DescriptionIndexer interface {
	IndexPkg(pkg *model.Package)
}

// NewIndexer creates a new indexer to download PACKAGES and store into database.
func NewIndexer(url string, descIndexer DescriptionIndexer, repo repository.Repository) *Indexer {
	return &Indexer{
		url:         url,
		descIndexer: descIndexer,
		repo:        repo,
	}
}

type Indexer struct {
	url         string
	descIndexer DescriptionIndexer
	repo        repository.Repository
}

func (i *Indexer) Run() {
	pkgs, err := i.downloadAndParsePackages()
	if err != nil {
		log.Println("Error while processing PACKAGES, err: ", err)
		return
	}

	if err := i.processPkgs(pkgs); err != nil {
		log.Println("Error while processing pkgs. Err: ", err)
		return
	}
}

func (i *Indexer) processPkgs(pkgs []*model.Package) error {
	for _, pkg := range pkgs {
		new, err := i.insert(pkg)
		if err != nil {
			log.Println("Error while storing ", pkg.Name)
			return err
		}

		if !new {
			continue
		}

		i.descIndexer.IndexPkg(pkg)
	}

	return nil
}

func (i *Indexer) downloadAndParsePackages() ([]*model.Package, error) {
	resp, err := http.Get(i.url + packagesFile)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return i.parsePackages(resp.Body)
}

func (i *Indexer) parsePackages(data io.Reader) ([]*model.Package, error) {
	sc := bufio.NewScanner(data)

	pkgs := []*model.Package{}
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, packagePrefix) {
			continue
		}

		name := strings.TrimSpace(strings.TrimPrefix(line, packagePrefix))
		pkg := &model.Package{
			Name: name,
		}

		if err := i.scanPackage(sc, pkg); err != nil {
			return nil, err
		}

		pkgs = append(pkgs, pkg)
	}

	return pkgs, sc.Err()
}

func (i *Indexer) scanPackage(sc *bufio.Scanner, pkg *model.Package) error {
	var err error
	if pkg.Version, err = scanWithPrefix(sc, versionPrefix); err != nil {
		return err
	}

	return nil
}

func (i *Indexer) insert(pkg *model.Package) (bool, error) {
	exist, err := i.repo.ExistByNameAndVersion(pkg.Name, pkg.Version)
	if err != nil {
		return false, err
	}

	if exist {
		return false, nil
	}

	if err := i.repo.Insert(pkg); err != nil {
		return false, err
	}

	return true, nil
}
