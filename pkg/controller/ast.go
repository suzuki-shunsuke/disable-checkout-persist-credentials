package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

const falseStr = "false"

func parseWorkflowAST(_ *logrus.Entry, content []byte, jobNames map[string]struct{}) (string, error) {
	file, err := parser.ParseBytes(content, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("parse a workflow file as YAML: %w", err)
	}
	for _, doc := range file.Docs {
		if err := parseDocAST(doc, jobNames); err != nil {
			return "", err
		}
	}
	return file.String(), nil
}

func parseDocAST(doc *ast.DocumentNode, jobNames map[string]struct{}) error {
	body, ok := doc.Body.(*ast.MappingNode)
	if !ok {
		return errors.New("document body must be *ast.MappingNode")
	}
	// jobs:
	//   jobName:
	//     steps:
	jobsNode := findJobsNode(body.Values)
	if jobsNode == nil {
		return errors.New("the field 'jobs' is required")
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

func parseDocValue(value *ast.MappingValueNode, jobNames map[string]struct{}) error {
	values, err := getMappingValueNodes(value)
	if err != nil {
		return err
	}
	for _, job := range values {
		if err := parseJobAST(job, jobNames); err != nil {
			return err
		}
	}
	return nil
}

func parseJobAST(value *ast.MappingValueNode, jobNames map[string]struct{}) error { //nolint:funlen,gocognit,cyclop
	jobNameNode, ok := value.Key.(*ast.StringNode)
	if !ok {
		return errors.New("job name must be a string")
	}
	jobName := jobNameNode.Value
	if _, ok := jobNames[jobName]; !ok {
		return nil
	}
	fields, err := getMappingValueNodes(value)
	if err != nil {
		return logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
			"job": jobName,
		})
	}
	if len(fields) == 0 {
		return logerr.WithFields(errors.New("job doesn't have any field"), logrus.Fields{ //nolint:wrapcheck
			"job": jobName,
		})
	}
	stepsField := findNodeByKey(fields, "steps")
	if stepsField == nil {
		return nil
	}
	stepSeq, ok := stepsField.Value.(*ast.SequenceNode)
	if !ok {
		return logerr.WithFields(errors.New("steps must be a sequence"), logrus.Fields{ //nolint:wrapcheck
			"job": jobName,
		})
	}
	for _, stepNode := range stepSeq.Values {
		// uses: actions/checkout@v2
		stepM, ok := stepNode.(*ast.MappingNode)
		if !ok {
			return logerr.WithFields(errors.New("step must be a mapping"), logrus.Fields{ //nolint:wrapcheck
				"job": jobName,
			})
		}
		usesNode := findNodeByKey(stepM.Values, "uses")
		if usesNode == nil {
			continue
		}
		s, ok := usesNode.Value.(*ast.StringNode)
		if !ok {
			return logerr.WithFields(errors.New("uses must be a string"), logrus.Fields{ //nolint:wrapcheck
				"job": jobName,
			})
		}
		if !strings.HasPrefix(s.Value, "actions/checkout@") {
			continue
		}
		withNode := findNodeByKey(stepM.Values, "with")
		if withNode == nil {
			node, err := yaml.ValueToNode(map[string]any{
				"with": map[string]any{
					"persist-credentials": false,
				},
			})
			if err != nil {
				return fmt.Errorf("convert packages to node: %w", err)
			}
			stepM.Merge(node.(*ast.MappingNode)) //nolint:forcetypeassert
			continue
		}
		withM, ok := withNode.Value.(*ast.MappingNode)
		if !ok {
			return logerr.WithFields(errors.New("with must be a mapping"), logrus.Fields{ //nolint:wrapcheck
				"job": jobName,
			})
		}
		pc := findNodeByKey(withM.Values, "persist-credentials")
		if pc == nil {
			node, err := yaml.ValueToNode(map[string]any{
				"persist-credentials": false,
			})
			if err != nil {
				return fmt.Errorf("convert packages to node: %w", err)
			}
			withM.Merge(node.(*ast.MappingNode)) //nolint:forcetypeassert
			continue
		}
		switch v := pc.Value.(type) {
		case *ast.BoolNode:
			if v.Value {
				// TODO: Known issue: This doesn't work well.
				v.Value = false
				continue
			}
		case *ast.StringNode:
			if v.Value != falseStr {
				v.Value = falseStr
				continue
			}
		default:
			return logerr.WithFields(errors.New("persist-credentials must be a string or boolean"), logrus.Fields{ //nolint:wrapcheck
				"job": jobName,
			})
		}
	}
	return nil
}
