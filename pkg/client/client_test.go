package client

import (
	"context"
	"testing"

	"k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/client-go/rest"

	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"

	"k8s.io/client-go/discovery"
	admissionregistrationv1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	admissionregistrationv1beta1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1beta1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	appsv1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	appsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	appsv1alpha1 "k8s.io/client-go/kubernetes/typed/auditregistration/v1alpha1"
	authenticationv1 "k8s.io/client-go/kubernetes/typed/authentication/v1"
	authenticationv1beta1 "k8s.io/client-go/kubernetes/typed/authentication/v1beta1"
	authorizationv1 "k8s.io/client-go/kubernetes/typed/authorization/v1"
	authorizationv1beta1 "k8s.io/client-go/kubernetes/typed/authorization/v1beta1"
	autoscalingv1 "k8s.io/client-go/kubernetes/typed/autoscaling/v1"
	autoscalingv2beta1 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	autoscalingv2beta2 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta2"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	batchv1beta1 "k8s.io/client-go/kubernetes/typed/batch/v1beta1"
	batchv2alpha1 "k8s.io/client-go/kubernetes/typed/batch/v2alpha1"
	certificatesv1beta1 "k8s.io/client-go/kubernetes/typed/certificates/v1beta1"
	coordinationv1 "k8s.io/client-go/kubernetes/typed/coordination/v1"
	coordinationv1beta1 "k8s.io/client-go/kubernetes/typed/coordination/v1beta1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	discoveryv1alpha1 "k8s.io/client-go/kubernetes/typed/discovery/v1alpha1"
	discoveryv1beta1 "k8s.io/client-go/kubernetes/typed/discovery/v1beta1"
	eventsv1beta1 "k8s.io/client-go/kubernetes/typed/events/v1beta1"
	extensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	flowcontrolv1alpha1 "k8s.io/client-go/kubernetes/typed/flowcontrol/v1alpha1"
	networkingv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
	networkingv1beta1 "k8s.io/client-go/kubernetes/typed/networking/v1beta1"
	nodev1alpha1 "k8s.io/client-go/kubernetes/typed/node/v1alpha1"
	nodev1beta1 "k8s.io/client-go/kubernetes/typed/node/v1beta1"
	policyv1beta1 "k8s.io/client-go/kubernetes/typed/policy/v1beta1"
	rbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	rbacv1alpha1 "k8s.io/client-go/kubernetes/typed/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	schedulingv1 "k8s.io/client-go/kubernetes/typed/scheduling/v1"
	schedulingv1alpha1 "k8s.io/client-go/kubernetes/typed/scheduling/v1alpha1"
	schedulingv1beta1 "k8s.io/client-go/kubernetes/typed/scheduling/v1beta1"
	settingsv1alpha1 "k8s.io/client-go/kubernetes/typed/settings/v1alpha1"
	v1storagev1 "k8s.io/client-go/kubernetes/typed/storage/v1"
	storagev1alpha1 "k8s.io/client-go/kubernetes/typed/storage/v1alpha1"
	storagev1beta1 "k8s.io/client-go/kubernetes/typed/storage/v1beta1"
)

func TestCreateOrUpdateClusterRoleBinding(t *testing.T) {
	f := manifests.NewFactory("openshift-monitoring", "openshift-user-workload-monitoring", manifests.NewDefaultConfig())
	crb, err := f.PrometheusK8sClusterRoleBinding()
	if err != nil {
		t.Fatal(err)
	}

	var c Client
	c.kclient = &kubeMockClient{}
	err = c.CreateOrUpdateClusterRoleBinding(crb)
	if err != nil {
		t.Fatal(err)
	}
}

type kubeMockClient struct {
	getResult *v1.ClusterRoleBinding
	getError  error
}

func (k *kubeMockClient) Create(ctx context.Context, clusterRoleBinding *v1.ClusterRoleBinding, opts metav1.CreateOptions) (*v1.ClusterRoleBinding, error) {
	panic("not implemented")
}

func (k *kubeMockClient) Update(ctx context.Context, clusterRoleBinding *v1.ClusterRoleBinding, opts metav1.UpdateOptions) (*v1.ClusterRoleBinding, error) {
	panic("not implemented")
}

func (k *kubeMockClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	panic("not implemented")
}

func (k *kubeMockClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	panic("not implemented")
}

func (k *kubeMockClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ClusterRoleBinding, error) {
	return k.getResult, k.getError
}

func (k *kubeMockClient) List(ctx context.Context, opts metav1.ListOptions) (*v1.ClusterRoleBindingList, error) {
	panic("not implemented")
}

func (k *kubeMockClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	panic("not implemented")
}

func (k *kubeMockClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ClusterRoleBinding, err error) {
	panic("not implemented")
}

func (k *kubeMockClient) RESTClient() rest.Interface {
	panic("not implemented")
}

func (k *kubeMockClient) ClusterRoles() rbacv1.ClusterRoleInterface {
	panic("not implemented")
}

func (k *kubeMockClient) ClusterRoleBindings() rbacv1.ClusterRoleBindingInterface {
	return k
}

func (k *kubeMockClient) Roles(namespace string) rbacv1.RoleInterface {
	panic("not implemented")
}

func (k *kubeMockClient) RoleBindings(namespace string) rbacv1.RoleBindingInterface {
	panic("not implemented")
}

func (k *kubeMockClient) Discovery() discovery.DiscoveryInterface {
	panic("not implemented")
}

func (k *kubeMockClient) AdmissionregistrationV1() admissionregistrationv1.AdmissionregistrationV1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) AdmissionregistrationV1beta1() admissionregistrationv1beta1.AdmissionregistrationV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) AppsV1() appsv1.AppsV1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) AppsV1beta1() appsv1beta1.AppsV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) AppsV1beta2() appsv1beta2.AppsV1beta2Interface {
	panic("not implemented")
}

func (k *kubeMockClient) AuditregistrationV1alpha1() appsv1alpha1.AuditregistrationV1alpha1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) AuthenticationV1() authenticationv1.AuthenticationV1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) AuthenticationV1beta1() authenticationv1beta1.AuthenticationV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) AuthorizationV1() authorizationv1.AuthorizationV1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) AuthorizationV1beta1() authorizationv1beta1.AuthorizationV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) AutoscalingV1() autoscalingv1.AutoscalingV1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) AutoscalingV2beta1() autoscalingv2beta1.AutoscalingV2beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) AutoscalingV2beta2() autoscalingv2beta2.AutoscalingV2beta2Interface {
	panic("not implemented")
}

func (k *kubeMockClient) BatchV1() batchv1.BatchV1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) BatchV1beta1() batchv1beta1.BatchV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) BatchV2alpha1() batchv2alpha1.BatchV2alpha1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) CertificatesV1beta1() certificatesv1beta1.CertificatesV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) CoordinationV1beta1() coordinationv1beta1.CoordinationV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) CoordinationV1() coordinationv1.CoordinationV1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) CoreV1() corev1.CoreV1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) DiscoveryV1alpha1() discoveryv1alpha1.DiscoveryV1alpha1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) DiscoveryV1beta1() discoveryv1beta1.DiscoveryV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) EventsV1beta1() eventsv1beta1.EventsV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) ExtensionsV1beta1() extensionsv1beta1.ExtensionsV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) FlowcontrolV1alpha1() flowcontrolv1alpha1.FlowcontrolV1alpha1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) NetworkingV1() networkingv1.NetworkingV1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) NetworkingV1beta1() networkingv1beta1.NetworkingV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) NodeV1alpha1() nodev1alpha1.NodeV1alpha1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) NodeV1beta1() nodev1beta1.NodeV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) PolicyV1beta1() policyv1beta1.PolicyV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) RbacV1() rbacv1.RbacV1Interface {
	return k
}

func (k *kubeMockClient) RbacV1beta1() rbacv1beta1.RbacV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) RbacV1alpha1() rbacv1alpha1.RbacV1alpha1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) SchedulingV1alpha1() schedulingv1alpha1.SchedulingV1alpha1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) SchedulingV1beta1() schedulingv1beta1.SchedulingV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) SchedulingV1() schedulingv1.SchedulingV1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) SettingsV1alpha1() settingsv1alpha1.SettingsV1alpha1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) StorageV1beta1() storagev1beta1.StorageV1beta1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) StorageV1() v1storagev1.StorageV1Interface {
	panic("not implemented")
}

func (k *kubeMockClient) StorageV1alpha1() storagev1alpha1.StorageV1alpha1Interface {
	panic("not implemented")
}
