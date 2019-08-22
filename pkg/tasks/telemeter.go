// Copyright 2018 The Cluster Monitoring Operator Authors
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

package tasks

import (
	"github.com/openshift/cluster-monitoring-operator/pkg/client"
	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TelemeterClientTask struct {
	client  *client.Client
	factory *manifests.Factory
	config  *manifests.TelemeterClientConfig
}

func NewTelemeterClientTask(client *client.Client, factory *manifests.Factory, config *manifests.TelemeterClientConfig) *TelemeterClientTask {
	return &TelemeterClientTask{
		client:  client,
		factory: factory,
		config:  config,
	}
}

func (t *TelemeterClientTask) Run() error {
	if t.config.IsEnabled() {
		return t.create()
	}

	return t.destroy()
}

func (t *TelemeterClientTask) create() error {
	cacm, err := t.factory.TelemeterClientServingCertsCABundle()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter Client serving certs CA Bundle ConfigMap failed")
	}

	err = t.client.CreateIfNotExistConfigMap(cacm)
	if err != nil {
		return errors.Wrap(err, "creating Telemeter Client serving certs CA Bundle ConfigMap failed")
	}

	sa, err := t.factory.TelemeterClientServiceAccount()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client Service failed")
	}

	err = t.client.CreateOrUpdateServiceAccount(sa)
	if err != nil {
		return errors.Wrap(err, "reconciling Telemeter client ServiceAccount failed")
	}

	cr, err := t.factory.TelemeterClientClusterRole()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client ClusterRole failed")
	}

	err = t.client.CreateOrUpdateClusterRole(cr)
	if err != nil {
		return errors.Wrap(err, "reconciling Telemeter client ClusterRole failed")
	}

	crb, err := t.factory.TelemeterClientClusterRoleBinding()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client ClusterRoleBinding failed")
	}

	err = t.client.CreateOrUpdateClusterRoleBinding(crb)
	if err != nil {
		return errors.Wrap(err, "reconciling Telemeter client ClusterRoleBinding failed")
	}

	crb, err = t.factory.TelemeterClientClusterRoleBindingView()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client cluster monitoring view ClusterRoleBinding failed")
	}

	err = t.client.CreateOrUpdateClusterRoleBinding(crb)
	if err != nil {
		return errors.Wrap(err, "reconciling Telemeter client cluster monitoring view ClusterRoleBinding failed")
	}

	svc, err := t.factory.TelemeterClientService()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client Service failed")
	}

	err = t.client.CreateOrUpdateService(svc)
	if err != nil {
		return errors.Wrap(err, "reconciling Telemeter client Service failed")
	}

	s, err := t.factory.TelemeterClientSecret()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client Secret failed")
	}

	err = t.client.CreateOrUpdateSecret(s)
	if err != nil {
		return errors.Wrap(err, "reconciling Telemeter client Secret failed")
	}
	{
		// We want to rollout a new deployment of telemeter whenever the configmap telemeter-trusted-ca-bundle is updated.
		// Because we react on all events in the same way, we cannot know when the CA TLS actually
		// changes, so we do a hash style based rollout, similiar to the prometheus-adapter does.
		proxyCM, err := t.client.GetConfigmap("openshift-monitoring", "telemeter-trusted-ca-bundle")
		if err != nil {
			// Sometimes the ConfigMap might not be created already, if that is the case we should not
			// error out.
			if !apierrors.IsNotFound(err) {
				return errors.Wrap(err, "failed to get telemeter-trusted-ca-bundle ConfigMap")
			}
		}
		if proxyCM != nil {
			proxyCM, err := t.factory.TelemeterConfigmapHash(proxyCM)
			if err != nil {
				return errors.Wrap(err, "failed to initialize telemeter-trusted-ca-bundle-<hash> ConfigMap")
			}
			// In the case when there is no data but the ConfigMap is there, we just continue.
			// We will catch this on the next loop.
			if proxyCM != nil {
				err = t.deleteOldTelemeterConfigMaps(string(proxyCM.Labels["monitoring.openshift.io/hash"]))
				if err != nil {
					return errors.Wrap(err, "deleting old telemeter configmaps failed")
				}

				err = t.client.CreateOrUpdateConfigMap(proxyCM)
				if err != nil {
					return errors.Wrap(err, "reconciling Telemeter telemeter-trusted-ca-bundle-<hash> ConfigMap failed")
				}

			}
		}
		dep, err := t.factory.TelemeterClientDeployment(proxyCM)
		if err != nil {
			return errors.Wrap(err, "initializing Telemeter client Deployment failed")
		}

		err = t.client.CreateOrUpdateDeployment(dep)
		if err != nil {
			return errors.Wrap(err, "reconciling Telemeter client Deployment failed")
		}
	}

	sm, err := t.factory.TelemeterClientServiceMonitor()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client ServiceMonitor failed")
	}

	err = t.client.CreateOrUpdateServiceMonitor(sm)
	return errors.Wrap(err, "reconciling Telemeter client ServiceMonitor failed")
}

func (t *TelemeterClientTask) destroy() error {
	dep, err := t.factory.TelemeterClientDeployment(nil)
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client Deployment failed")
	}

	err = t.client.DeleteDeployment(dep)
	if err != nil {
		return errors.Wrap(err, "deleting Telemeter client Deployment failed")
	}

	s, err := t.factory.TelemeterClientSecret()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client Secret failed")
	}

	err = t.client.DeleteSecret(s)
	if err != nil {
		return errors.Wrap(err, "deleting Telemeter client Secret failed")
	}

	svc, err := t.factory.TelemeterClientService()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client Service failed")
	}

	err = t.client.DeleteService(svc)
	if err != nil {
		return errors.Wrap(err, "deleting Telemeter client Service failed")
	}

	crb, err := t.factory.TelemeterClientClusterRoleBinding()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client ClusterRoleBinding failed")
	}

	err = t.client.DeleteClusterRoleBinding(crb)
	if err != nil {
		return errors.Wrap(err, "deleting Telemeter client ClusterRoleBinding failed")
	}

	cr, err := t.factory.TelemeterClientClusterRole()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client ClusterRole failed")
	}

	err = t.client.DeleteClusterRole(cr)
	if err != nil {
		return errors.Wrap(err, "deleting Telemeter client ClusterRole failed")
	}

	sa, err := t.factory.TelemeterClientServiceAccount()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client Service failed")
	}

	err = t.client.DeleteServiceAccount(sa)
	if err != nil {
		return errors.Wrap(err, "deleting Telemeter client ServiceAccount failed")
	}

	sm, err := t.factory.TelemeterClientServiceMonitor()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter client ServiceMonitor failed")
	}

	err = t.client.DeleteServiceMonitor(sm)
	if err != nil {
		return errors.Wrap(err, "deleting Telemeter client ServiceMonitor failed")
	}

	// TODO: Should we delete the ConfigMaps we create here?

	cacm, err := t.factory.TelemeterClientServingCertsCABundle()
	if err != nil {
		return errors.Wrap(err, "initializing Telemeter Client serving certs CA Bundle ConfigMap failed")
	}

	err = t.client.DeleteConfigMap(cacm)
	return errors.Wrap(err, "creating Telemeter Client serving certs CA Bundle ConfigMap failed")
}

func (t *TelemeterClientTask) deleteOldTelemeterConfigMaps(newHash string) error {
	configMaps, err := t.client.KubernetesInterface().CoreV1().ConfigMaps("openshift-monitoring").List(metav1.ListOptions{
		LabelSelector: "monitoring.openshift.io/name=telemeter,monitoring.openshift.io/hash!=" + newHash,
	})
	if err != nil {
		return errors.Wrap(err, "error listing telemeter configmaps while deleting old telemeter configmaps")
	}

	for i := range configMaps.Items {
		err := t.client.KubernetesInterface().CoreV1().ConfigMaps("openshift-monitoring").Delete(configMaps.Items[i].Name, &metav1.DeleteOptions{})
		if err != nil {
			return errors.Wrapf(err, "error deleting secret: %s", configMaps.Items[i].Name)
		}
	}

	return nil
}
