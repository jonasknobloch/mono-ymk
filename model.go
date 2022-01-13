package main

import (
	"gonum.org/v1/gonum/stat/combin"
)

type Model struct {
	n map[string]map[string]float64
	r map[string]map[string]float64
	t map[string]map[string]float64
}

func NewModel() *Model {
	return &Model{
		n: make(map[string]map[string]float64),
		r: make(map[string]map[string]float64),
		t: make(map[string]map[string]float64),
	}
}

func (m *Model) InitInsertionWeights(dictionary map[string]map[string]int) {
	for feature, keys := range dictionary {
		m.n[feature] = make(map[string]float64, len(keys))

		for key := range keys {
			m.n[feature][key] = 1 / float64(len(keys))
		}
	}
}

func (m *Model) InitTranslationWeights(dictionary map[string]map[string]int) {
	for feature, keys := range dictionary {
		m.t[feature] = make(map[string]float64, len(keys))

		for key := range keys {
			m.t[feature][key] = 1 / float64(len(keys))
		}
	}
}

func (m *Model) PInsertion(insertion Insertion) float64 {
	return m.n[insertion.Feature()][insertion.Key()]
}

func (m *Model) PReordering(reordering Reordering) float64 {
	if _, ok := m.r[reordering.feature]; !ok {
		m.r[reordering.feature] = make(map[string]float64)

		n := combin.NumPermutations(len(reordering.Reordering), len(reordering.Reordering))
		g := combin.NewPermutationGenerator(len(reordering.Reordering), len(reordering.Reordering))

		for g.Next() {
			m.r[reordering.feature][NewReordering(g.Permutation(nil), "").key] = 1 / float64(n)
		}
	}

	return m.r[reordering.feature][reordering.key]
}

func (m *Model) PTranslation(translation Translation) float64 {
	return m.t[translation.feature][translation.key]
}

func (m *Model) UpdateWeights(insertionCount, reorderingCount, translationCount Count) {
	update := func(p map[string]map[string]float64, c Count) {
		for feature, keys := range p {
			for key := range keys {
				p[feature][key] = c.Get(feature, key) / c.Sum(feature)
			}
		}
	}

	update(m.n, insertionCount)
	update(m.r, reorderingCount)
	update(m.t, translationCount)
}
