package qlmodel

import commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"

type PracticalUserActivityAnnotations struct {
	Label string
	Type  commonModel.AnnotationType
	Group string
}

type PracticalUserActivityAnnotationsInput struct {
	Label string
	Type  commonModel.AnnotationType
	Group string
}
