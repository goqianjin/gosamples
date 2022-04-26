package gomega

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestFarmHasCow(t *testing.T) {
	//gomega.RegisterFailHandler(ginkgo.Fail)
	g := gomega.NewGomegaWithT(t)

	g.Expect("actual").To(gomega.Equal("expected"), "Farm should have cow")
}
