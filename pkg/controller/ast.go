package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Position struct {
	JobKey string
	Line   int
	Column int
}

func parseWorkflowAST(logE *logrus.Entry, content []byte, jobNames map[string]struct{}) ([]*Position, error) {
	file, err := parser.ParseBytes(content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse a workflow file as YAML: %w", err)
	}
	list := []*Position{}
	for _, doc := range file.Docs {
		arr, err := parseDocAST(doc, jobNames)
		if err != nil {
			return nil, err
		}
		if len(arr) == 0 {
			continue
		}
		list = append(list, arr...)
	}
	return list, nil
}

func parseDocAST(doc *ast.DocumentNode, jobNames map[string]struct{}) ([]*Position, error) {
	body, ok := doc.Body.(*ast.MappingNode)
	if !ok {
		return nil, errors.New("document body must be *ast.MappingNode")
	}
	// jobs:
	//   jobName:
	//     steps:
	jobsNode := findJobsNode(body.Values)
	if jobsNode == nil {
		return nil, errors.New("the field 'jobs' is required")
	}
	return parseDocValue(jobsNode, jobNames)
}

func findJobsNode(values []*ast.MappingValueNode) *ast.MappingValueNode {
	return findNodeByKey(values, "jobs")
}

func findNodeByKey(values []*ast.MappingValueNode, key string) *ast.MappingValueNode {
	for _, value := range values {
		k, ok := value.Key.(*ast.StringNode)
		if !ok {
			continue
		}
		if k.Value == key {
			return value
		}
	}
	return nil
}

func getMappingValueNodes(value *ast.MappingValueNode) ([]*ast.MappingValueNode, error) {
	switch node := value.Value.(type) {
	case *ast.MappingNode:
		return node.Values, nil
	case *ast.MappingValueNode:
		return []*ast.MappingValueNode{node}, nil
	}
	return nil, errors.New("value must be either a *ast.MappingNode or a *ast.MappingValueNode")
}

func parseDocValue(value *ast.MappingValueNode, jobNames map[string]struct{}) ([]*Position, error) {
	values, err := getMappingValueNodes(value)
	if err != nil {
		return nil, err
	}
	arr := make([]*Position, 0, len(values))
	for _, job := range values {
		pos, err := parseJobAST(job, jobNames)
		if err != nil {
			return nil, err
		}
		if pos == nil {
			continue
		}
		arr = append(arr, pos)
	}
	return arr, nil
}

func parseJobAST(value *ast.MappingValueNode, jobNames map[string]struct{}) (*Position, error) {
	jobNameNode, ok := value.Key.(*ast.StringNode)
	if !ok {
		return nil, errors.New("job name must be a string")
	}
	jobName := jobNameNode.Value
	if _, ok := jobNames[jobName]; !ok {
		return nil, nil //nolint:nilnil
	}
	fields, err := getMappingValueNodes(value)
	if err != nil {
		return nil, logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
			"job": jobName,
		})
	}
	if len(fields) == 0 {
		return nil, logerr.WithFields(errors.New("job doesn't have any field"), logrus.Fields{ //nolint:wrapcheck
			"job": jobName,
		})
	}
	// get steps field
	stepsField := findNodeByKey(fields, "steps")
	if stepsField == nil {
		return nil, nil
	}
	stepSeq, ok := stepsField.Value.(*ast.SequenceNode)
	if !ok {
		return nil, logerr.WithFields(errors.New("steps must be a sequence"), logrus.Fields{ //nolint:wrapcheck
			"job": jobName,
		})
	}
	for _, stepNode := range stepSeq.Values {
		// uses: actions/checkout@v2
		stepM, ok := stepNode.(*ast.MappingNode)
		if !ok {
			return nil, logerr.WithFields(errors.New("step must be a mapping"), logrus.Fields{ //nolint:wrapcheck
				"job": jobName,
			})
		}
		usesNode := findNodeByKey(stepM.Values, "uses")
		if usesNode == nil {
			continue
		}
		s, ok := usesNode.Value.(*ast.StringNode)
		if !ok {
			return nil, logerr.WithFields(errors.New("uses must be a string"), logrus.Fields{ //nolint:wrapcheck
				"job": jobName,
			})
		}
		if !strings.HasPrefix(s.Value, "actions/checkout@") {
			continue
		}
		withNode := findNodeByKey(stepM.Values, "with")
		if withNode == nil {
			fmt.Println("with is nil")
			continue
		}
		withM, ok := withNode.Value.(*ast.MappingNode)
		if !ok {
			return nil, logerr.WithFields(errors.New("with must be a mapping"), logrus.Fields{ //nolint:wrapcheck
				"job": jobName,
			})
		}
		pc := findNodeByKey(withM.Values, "persist-credentials")
		if pc == nil {
			fmt.Println("persist-credentials is nil")
			continue
		}
		switch v := pc.Value.(type) {
		case *ast.BoolNode:
			if v.Value {
				fmt.Println("persist-credentials is true")
				continue
			}
		case *ast.StringNode:
			if v.Value != "false" {
				fmt.Println("persist-credentials isn't false")
				continue
			}
		default:
			return nil, logerr.WithFields(errors.New("persist-credentials must be a string or boolean"), logrus.Fields{ //nolint:wrapcheck
				"job": jobName,
			})
		}
	}
	return nil, nil
}
