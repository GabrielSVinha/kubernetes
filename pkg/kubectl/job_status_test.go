package kubectl

import (
	"testing"

	batch "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestJobStatusViewerStatus(t *testing.T) {
	tests := []struct {
		status      batch.JobStatus
		parallelism int32
		msg         string
		done        bool
	}{
		{
			status: batch.JobStatus{
				Succeeded: 10,
			},
			parallelism: 20,
			msg:         "Waiting for workers to finish: 10 out of 20 have been completed.",
			done:        false,
		},
		{
			status: batch.JobStatus{
				Succeeded: 10,
			},
			parallelism: 100,
			msg:         "Waiting for workers to finish: 10 out of 100 have been completed.",
			done:        false,
		}, {
			status: batch.JobStatus{
				Succeeded: 100,
			},
			parallelism: 100,
			msg:         "All 100 workers completed",
			done:        true,
		},
	}

	for _, test := range tests {
		j := &batch.Job{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "bar",
				Name:      "foo",
				UID:       "8764ae47-9092-11e4-8393-42010af018ff",
			},
			Spec: batch.JobSpec{
				Parallelism: &test.parallelism,
			},
			Status: test.status,
		}
		client := fake.NewSimpleClientset(j).Batch()
		jsv := &JobStatusViewer{job: client}
		msg, done, err := jsv.Status("bar", "foo")
		if err != nil {
			t.Fatalf("JobStatusViewer.Status(): %v", err)
		}
		if done != test.done {
			t.Errorf("JobStatusViewer.Status() for job with parallelism %d, and status %+v returned %q, %t, want %q, %t",
				test.parallelism,
				test.status,
				msg,
				done,
				test.msg,
				test.done,
			)
		}
	}
}
