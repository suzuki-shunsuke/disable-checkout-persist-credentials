package controller

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

func (c *Controller) Run(logE *logrus.Entry, input *Input) error {
	// Read workflow and action files
	// Parse files and detect actions/checkout steps which persist-credentials is not set or true
	// job and step index
	files := input.Args
	if len(files) == 0 {
		a, err := c.findWorkflowFiles()
		if err != nil {
			return fmt.Errorf("find workflow files: %w", err)
		}
		files = a
	}
	for _, file := range files {
		logE := logE.WithField("file", file)
		if err := c.handleWorkflow(logE, file); err != nil {
			return fmt.Errorf("handle a workflow file: %w", err)
		}
	}
	return nil
}

func (c *Controller) findWorkflowFiles() ([]string, error) {
	var files []string
	for _, file := range []string{"*.yaml", "*.yml"} {
		a, err := afero.Glob(c.fs, filepath.Join(".github", "workflows", file))
		if err != nil {
			return nil, fmt.Errorf("find workflow files: %w", err)
		}
		files = append(files, a...)
	}
	return files, nil
}

func (c *Controller) handleWorkflow(logE *logrus.Entry, file string) error {
	content, err := afero.ReadFile(c.fs, file)
	if err != nil {
		return fmt.Errorf("read a file: %w", err)
	}

	wf := &Workflow{}
	if err := yaml.Unmarshal(content, wf); err != nil {
		return fmt.Errorf("unmarshal a workflow file: %w", err)
	}
	jobNames, err := wf.Validate()
	if err != nil {
		return fmt.Errorf("validate a workflow: %w", err)
	}
	if len(jobNames) == 0 {
		return nil
	}
	logE.Info("change the workflow")
	_, err = parseWorkflowAST(logE, content, jobNames)
	if err != nil {
		return err
	}
	return nil
}

type Workflow struct {
	Jobs map[string]*Job
}

type Job struct {
	Steps []*Step
}

type Step struct {
	Uses string
	With map[string]any
}

func (w *Workflow) Validate() (map[string]struct{}, error) {
	jobNames := map[string]struct{}{}
	for jobName, job := range w.Jobs {
		f, err := job.Validate()
		if err != nil {
			return nil, fmt.Errorf("validate a job: %w", err)
		}
		if !f {
			jobNames[jobName] = struct{}{}
		}
	}
	return jobNames, nil
}

func (j *Job) Validate() (bool, error) {
	valid := true
	for _, step := range j.Steps {
		f, err := step.Validate()
		if err != nil {
			return false, fmt.Errorf("validate a step: %w", err)
		}
		if !f {
			valid = false
		}
	}
	return valid, nil
}

func (s *Step) Validate() (bool, error) {
	if !strings.HasPrefix(s.Uses, "actions/checkout@") {
		return true, nil
	}
	v, ok := s.With["persist-credentials"]
	if !ok {
		return false, nil
	}
	switch v := v.(type) {
	case bool:
		return !v, nil
	case string:
		return v == "false", nil
	default:
		return false, fmt.Errorf("the type of persist-credentials is invalid: %T", v)
	}
}
