// Copyright 2019 The Cluster Monitoring Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"testing"
	"time"

	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func TestTelemeterCARotation(t *testing.T) {
	var lastErr error
	// Wait for Telemeter
	err := wait.Poll(time.Second, 5*time.Minute, func() (bool, error) {
		_, err := f.KubeClient.AppsV1().Deployments(f.Ns).Get("telemeter-client", metav1.GetOptions{})
		lastErr = errors.Wrap(err, "getting telemeter deployment failed")
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		if err == wait.ErrWaitTimeout && lastErr != nil {
			err = lastErr
		}
		t.Fatal(err)
	}

	// This ConfigMap is already created by CMO in the first task.
	cm, err := f.KubeClient.CoreV1().ConfigMaps(f.Ns).Get("telemeter-trusted-ca-bundle", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}

	// Simulate rotation by simply adding a newline to existing certs.
	// This change will be propagated to the cluster monitoring operator,
	// causing a new CM to be created.
	dataContent := "foo-bar"
	data := "ca-bundle.crt"
	cm.Data[data] = cm.Data[data] + dataContent
	_, err = f.KubeClient.CoreV1().ConfigMaps(f.Ns).Update(cm)
	if err != nil {
		t.Fatal(err)
	}

	factory := manifests.NewFactory("openshift-monitoring", nil)
	newCM, err := factory.TelemeterConfigmapHash(cm)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for the new ConfigMap to be created
	err = wait.Poll(time.Second, 5*time.Minute, func() (bool, error) {
		_, err := f.KubeClient.CoreV1().ConfigMaps(f.Ns).Get(newCM.Name, metav1.GetOptions{})
		lastErr = errors.Wrap(err, "getting new CA ConfigMap failed")
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		if err == wait.ErrWaitTimeout && lastErr != nil {
			err = lastErr
		}
		t.Fatal(err)
	}

	// Get telemeter-client deployment and make sure it has a volumemounted.
	// TODO: We should check the volumemount name matches the CM name.
	err = wait.Poll(time.Second, 5*time.Minute, func() (bool, error) {
		d, err := f.KubeClient.AppsV1().Deployments(f.Ns).Get("telemeter-client", metav1.GetOptions{})
		lastErr = errors.Wrap(err, "getting telemeter deployment failed")
		if err != nil {
			return false, nil
		}
		if len(d.Spec.Template.Spec.Containers[0].VolumeMounts) == 0 {
			// TODO: Should we instead add to the lastErr
			return false, errors.New("Could not find any VolumeMounts, expected at least 1")
		}

		return true, nil
	})
	if err != nil {
		if err == wait.ErrWaitTimeout && lastErr != nil {
			err = lastErr
		}
		t.Fatal(err)
	}

	// TODO: Should we add a test for the following.
	// If we update the original ConfigMap again we should see:
	// - the old hashed CM deleted
	// - new CM created
}
