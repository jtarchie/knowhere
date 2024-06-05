package runtime

type Results []Result

func (r Results) Cluster(radius float64) Results {
	results := Results{}
	tree := &RTree{}

	for _, entry := range r {
		extended := entry.Bbox().Extend(radius)

		if !tree.Within(extended) {
			results = append(results, entry)
			tree.Insert(extended, &entry)
		}
	}

	return results
}

func (r Results) AsTree(radius float64) *RTree {
	tree := &RTree{}

	for _, entry := range r {
		extended := entry.Bbox().Extend(radius)

		tree.Insert(extended, &entry)
	}

	return tree
}
