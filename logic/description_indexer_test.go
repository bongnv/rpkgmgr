package logic

import (
	"strings"
	"testing"
	"time"

	"github.com/bongnv/rpkgmgr/model"
	"github.com/stretchr/testify/require"
)

func Test_parseDesc(t *testing.T) {
	desc := `
Package: genridge
Type: Package
Title: Generalized Ridge Trace Plots for Ridge Regression
Version: 0.6-6
Date: 2017-10-01
Author: Michael Friendly [aut, cre]
Authors@R: c(person(given = "Michael", family = "Friendly", role=c("aut", "cre"), email="friendly@yorku.ca"))
Maintainer: Michael Friendly <friendly@yorku.ca>
Depends: R (>= 2.11.1), car
Suggests: MASS, ElemStatLearn, rgl, bestglm
Description:
 The genridge package introduces generalizations of the standard univariate
 ridge trace plot used in ridge regression and related methods.  These graphical methods
 show both bias (actually, shrinkage) and precision, by plotting the covariance ellipsoids of the estimated
 coefficients, rather than just the estimates themselves.  2D and 3D plotting methods are provided,
 both in the space of the predictor variables and in the transformed space of the PCA/SVD of the
 predictors.
License: GPL (>= 2)
LazyLoad: yes
LazyData: yes
BugReports: https://github.com/friendly/genridge/issues
NeedsCompilation: no
Packaged: 2017-10-06 15:15:20 UTC; Friendly
Repository: CRAN
Date/Publication: 2017-10-06 15:30:27 UTC
`

	q := &queueDescriptionIndexer{}
	pkg := &model.Package{}
	err := q.parseDesc(desc, pkg)
	require.NoError(t, err)
	require.Equal(t, "Generalized Ridge Trace Plots for Ridge Regression", pkg.Title)
	require.Equal(t, "Michael Friendly [aut, cre]", pkg.Authors)
	expectedTime, _ := time.Parse(timeFormat, strings.TrimSpace("2017-10-06 15:30:27 UTC"))
	require.Equal(t, expectedTime, *pkg.PublicationDate)
}
