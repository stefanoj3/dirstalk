package tree

import (
	"sort"
	"strings"

	gotree "github.com/DiSiqueira/GoTree"
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

func NewResultTreeProducer() ResultTreeProducer {
	return ResultTreeProducer{}
}

type ResultTreeProducer struct{}

func (s ResultTreeProducer) String(results []scan.Result) string {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Target.Path < results[j].Target.Path
	})

	root := gotree.New("/")

	for _, r := range results {
		currentBranch := root

		parts := strings.Split(r.URL.Path, "/")
		for _, p := range parts {
			if len(p) == 0 {
				continue
			}

			found := false

			for _, item := range currentBranch.Items() {
				if item.Text() != p {
					continue
				}

				currentBranch = item
				found = true

				break
			}

			if found {
				continue
			}

			newTree := gotree.New(p)
			currentBranch.AddTree(newTree)
			currentBranch = newTree
		}
	}

	return root.Print()
}
