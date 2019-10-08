package logic

import (
	"strings"
	"testing"

	"github.com/bongnv/rpkgmgr/model"
	"github.com/bongnv/rpkgmgr/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_parsePackages(t *testing.T) {
	t.Run("happy-path", func(t *testing.T) {
		resp := `
Package: A3
Version: 1.0.0
Depends: R (>= 2.15.0), xtable, pbapply
Suggests: randomForest, e1071
License: GPL (>= 2)
NeedsCompilation: no

Package: aaSEA
Version: 1.0.0
Depends: R(>= 3.4.0)
Imports: DT(>= 0.4), networkD3(>= 0.4), shiny(>= 1.0.5),
        shinydashboard(>= 0.7.0), magrittr(>= 1.5), Bios2cor(>= 1.2),
        seqinr(>= 3.4-5), plotly(>= 4.7.1), Hmisc(>= 4.1-1)
Suggests: knitr, rmarkdown
License: GPL-3
NeedsCompilation: no
`

		i := &Indexer{}
		pkgs, err := i.parsePackages(strings.NewReader(resp))
		require.NoError(t, err)
		require.Len(t, pkgs, 2)
		require.Equal(t, "A3", pkgs[0].Name)
		require.Equal(t, "aaSEA", pkgs[1].Name)
		require.Equal(t, "1.0.0", pkgs[1].Version)
	})
}

func Test_processPkgs(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	i := &Indexer{
		repo: mockRepo,
	}

	t.Run("duplicate", func(t *testing.T) {
		mockRepo.On("ExistByNameAndVersion", "Some Name", "Some version").Return(true, nil).Once()
		res, err := i.insert(&model.Package{
			Name:    "Some Name",
			Version: "Some version",
		})

		require.NoError(t, err)
		require.False(t, res)
	})

	t.Run("new-insert", func(t *testing.T) {
		mockRepo.On("ExistByNameAndVersion", "Some Name", "1.1.0").Return(false, nil).Once()
		mockRepo.On("Insert", mock.Anything).Return(nil).Once()
		res, err := i.insert(&model.Package{
			Name:    "Some Name",
			Version: "1.1.0",
		})

		require.NoError(t, err)
		require.True(t, res)
	})

	mockRepo.AssertExpectations(t)
}
