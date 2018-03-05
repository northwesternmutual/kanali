// Copyright (c) 2018 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package kanali

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"

	"k8s.io/api/core/v1"
	"k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"github.com/northwesternmutual/kanali/test/builder"
	"github.com/northwesternmutual/kanali/test/e2e/context"
	"github.com/northwesternmutual/kanali/test/e2e/deploy"
)

type deployConfig struct {
	commitSHA  string
	serverType deploy.TLSType
	clientset  kubernetes.Interface
}

var (
	KanaliCACert, KanaliServerCert, KanaliServerKey []byte
)

const (
	name         = "kanali"
	SecurePort   = 8443
	InsecurePort = 8080
)

func Deploy(
	cli kubernetes.Interface,
	commitSHA string,
	options ...deploy.Option,
) error {
	cfg := &deployConfig{
		clientset: cli,
		commitSHA: commitSHA,
	}
	for _, op := range options {
		op(cfg)
	}

	if err := initApiKeyDecryptionKeys(2048); err != nil {
		return err
	}

	return cfg.deploy()
}

func (cfg *deployConfig) deploy() error {
	svc, err := cfg.deployService()
	if err != nil {
		return err
	}

	_, err = cfg.deployTLSSecret(svc)
	if err != nil {
		return err
	}

	_, err = cfg.deployApiKeyDecryptionSecret()
	if err != nil {
		return err
	}

	if err = cfg.deployRBACResources(); err != nil {
		return err
	}

	_, err = cfg.deployPod()
	if err != nil {
		return err
	}

	return nil
}

func Destroy(cli kubernetes.Interface) error {
	err := cli.RbacV1beta1().ClusterRoles().Delete(name, nil)
	err = cli.RbacV1beta1().ClusterRoleBindings().Delete(name, nil)
	return err
}

func (cfg *deployConfig) deployService() (*v1.Service, error) {
	var ports []v1.ServicePort

	if cfg.serverType == deploy.TLSTypeNone {
		ports = []v1.ServicePort{{
			Name: "http",
			Port: InsecurePort,
		}}
	} else {
		ports = []v1.ServicePort{{
			Name: "https",
			Port: SecurePort,
		}}
	}

	svc, err := cfg.clientset.CoreV1().Services(cfg.commitSHA).Create(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: cfg.commitSHA,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": name,
				"sha": cfg.commitSHA,
			},
			Ports: ports,
		},
	})
	if err != nil {
		return nil, err
	}

	context.TestContext.KanaliConfig.Port = int(svc.Spec.Ports[0].NodePort)
	return svc, nil
}

func (cfg *deployConfig) deployTLSSecret(svc *v1.Service) (*v1.Secret, error) {
	tlsAssets := builder.NewTLSBuilder([]string{
		fmt.Sprintf("%s.%s.svc.cluster.local", svc.GetName(), svc.GetNamespace()),
		fmt.Sprintf("%s.%s.svc", svc.GetName(), svc.GetNamespace()),
		fmt.Sprintf("%s.%s", svc.GetName(), svc.GetNamespace()),
		svc.GetName(),
		context.TestContext.KanaliConfig.Host,
	}, []net.IP{
		net.ParseIP(svc.Spec.ClusterIP),
	}).NewOrDie()

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: svc.GetNamespace(),
		},
		Type: v1.SecretTypeTLS,
		Data: map[string][]byte{
			"tls.ca":  tlsAssets.CACert,
			"tls.key": tlsAssets.ServerKey,
			"tls.crt": tlsAssets.ServerCert,
		},
	}

	if cfg.serverType == deploy.TLSTypeNone {
		context.TestContext.KanaliConfig.Scheme = "http"
		return nil, nil
	}

	context.TestContext.KanaliConfig.Scheme = "https"
	secret, err := cfg.clientset.CoreV1().Secrets(svc.GetNamespace()).Create(secret)
	if err != nil {
		return nil, err
	}

	KanaliCACert = tlsAssets.CACert
	if cfg.serverType == deploy.TLSTypeMutual {
		KanaliServerCert = tlsAssets.ServerCert
		KanaliServerKey = tlsAssets.ServerCert
	}
	return secret, nil
}

func (cfg *deployConfig) deployApiKeyDecryptionSecret() (*v1.Secret, error) {
	return cfg.clientset.CoreV1().Secrets(cfg.commitSHA).Create(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + "-apikey-decryption-key",
			Namespace: cfg.commitSHA,
		},
		Data: map[string][]byte{
			"key.pem": pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(TestApiKeyDecryptionPrivateKey)}),
		},
	})
}

func (cfg *deployConfig) deployRBACResources() error {
	if _, err := cfg.clientset.CoreV1().ServiceAccounts(cfg.commitSHA).Create(&v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: cfg.commitSHA,
			Labels: map[string]string{
				"app": name,
			},
		},
	}); err != nil {
		return err
	}

	if _, err := cfg.clientset.RbacV1beta1().ClusterRoles().Create(&v1beta1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Rules: []v1beta1.PolicyRule{
			{
				Verbs:     []string{"watch"},
				APIGroups: []string{"kanali.io"},
				Resources: []string{"apikeys", "apiproxies", "apikeybindings", "mocktargets"},
			},
			{
				Verbs:     []string{"create"},
				APIGroups: []string{"apiextensions"},
				Resources: []string{"customresourcedefinitions"},
			},
			{
				Verbs:     []string{"watch"},
				APIGroups: []string{""}, // needed to prevent error:
				Resources: []string{"services", "secrets", "configmaps"},
			},
		},
	}); err != nil {
		return err
	}

	_, err := cfg.clientset.RbacV1beta1().ClusterRoleBindings().Create(&v1beta1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Subjects: []v1beta1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      name,
				Namespace: cfg.commitSHA,
			},
		},
		RoleRef: v1beta1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     name,
		},
	})

	return err
}

func (cfg *deployConfig) deployPod() (*v1.Pod, error) {
	var port int

	args := []string{
		"--plugins.apiKey.decryption_key_file=/etc/kanali/rsa/key.pem",
		"--plugins.apiKey.header_key=apikey",
		"--proxy.enable_cluster_ip=true",
		"--proxy.enable_mock_responses=true",
		"--proxy.upstream_timeout=0h1m0s",
		"--proxy.tls_common_name_validation=true",
		//"--proxy.default_header_values=x-canary-deploy=stable",
		"--process.log_level=debug",
	}

	command := []string{
		"/kanali",
		"start",
	}

	insecureArgs := []string{
		fmt.Sprintf("--server.insecure_port=%d", InsecurePort),
	}
	secureArgs := []string{
		fmt.Sprintf("--server.secure_port=%d", SecurePort),
		"--tls.cert_file=/etc/kanali/pki/tls.crt",
		"--tls.key_file=/etc/kanali/pki/tls.key",
	}

	tlsProjectedVolume := v1.VolumeProjection{
		Secret: &v1.SecretProjection{
			LocalObjectReference: v1.LocalObjectReference{
				Name: name,
			},
			Items: []v1.KeyToPath{
				{
					Key:  "tls.crt",
					Path: "pki/tls.crt",
				},
				{
					Key:  "tls.key",
					Path: "pki/tls.key",
				},
				{
					Key:  "tls.key",
					Path: "pki/tls.ca",
				},
			},
		},
	}

	volumes := []v1.Volume{
		{
			Name: name,
			VolumeSource: v1.VolumeSource{
				Projected: &v1.ProjectedVolumeSource{
					Sources: []v1.VolumeProjection{
						{
							Secret: &v1.SecretProjection{
								LocalObjectReference: v1.LocalObjectReference{
									Name: name + "-apikey-decryption-key",
								},
								Items: []v1.KeyToPath{
									{
										Key:  "key.pem",
										Path: "rsa/key.pem",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	mounts := []v1.VolumeMount{
		{
			Name:      name,
			MountPath: "/etc/kanali",
		},
	}

	switch cfg.serverType {
	case deploy.TLSTypeNone:
		args = append(args, insecureArgs...)
		port = InsecurePort
	case deploy.TLSTypePresent:
		args = append(args, secureArgs...)
		port = SecurePort
		volumes[0].VolumeSource.Projected.Sources = append(volumes[0].VolumeSource.Projected.Sources, tlsProjectedVolume)
	case deploy.TLSTypeMutual:
		args = append(args, secureArgs...)
		args = append(args,
			"--tls.ca_file=/etc/kanali/pki/tls.ca",
		)
		volumes[0].VolumeSource.Projected.Sources = append(volumes[0].VolumeSource.Projected.Sources, tlsProjectedVolume)
		port = SecurePort
	}

	probe := &v1.Probe{
		Handler: v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.FromInt(port),
			},
		},
		TimeoutSeconds:      1,
		InitialDelaySeconds: 1,
	}

	po := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: cfg.commitSHA,
			Labels: map[string]string{
				"app": name,
				"sha": cfg.commitSHA,
			},
		},
		Spec: v1.PodSpec{
			ServiceAccountName: name,
			Volumes:            volumes,
			Containers: []v1.Container{
				{
					Name:            name,
					Image:           fmt.Sprintf("%s:%s", context.TestContext.KanaliConfig.ImageConfig.Name, context.TestContext.KanaliConfig.ImageConfig.Tag),
					VolumeMounts:    mounts,
					Command:         command,
					Args:            args,
					ImagePullPolicy: v1.PullIfNotPresent,
					LivenessProbe:   probe,
					ReadinessProbe:  probe,
				},
			},
		},
	}

	return deploy.Pod(cfg.clientset, po)
}

func (cfg *deployConfig) SetServerType(t deploy.TLSType) {
	cfg.serverType = t
}
