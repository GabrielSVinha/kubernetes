package kubectl

import (
	"fmt"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientbatchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

// JobStatusViewer implements the JobsViewer imterface
type JobStatusViewer struct {
	job clientbatchv1.JobsGetter
}

// Status returns a message describing the job status, and a bool value indicating if the status is considered done.
func (s *JobStatusViewer) Status(namespace, name string) (string, bool, error) {
	job, err := s.job.Jobs(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		return "", false, err
	}
	fmt.Println(job.Spec.Parallelism)
	if job.Status.Succeeded > *job.Spec.Parallelism {
		return "", false, fmt.Errorf("ERROR: Number of completed worker cannot be greater than specified")
	}
	if job.Status.Succeeded < *job.Spec.Parallelism {
		return fmt.Sprintf("Waiting for workers to finish: %d out of %d have been completed.", job.Status.Succeeded, job.Spec.Parallelism), false, nil
	}
	return fmt.Sprintf("All %d workers completed", *job.Spec.Parallelism), true, nil
}
