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

package tester

import (
	"fmt"
	"net"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"github.com/northwesternmutual/kanali/test/builder"
	"github.com/northwesternmutual/kanali/test/e2e/deploy"
)

type tlsType int

const (
	TLSTypeNone tlsType = iota
	TLSTypePresent
	TLSTypeMutual

	Name         = "tester"
	image        = "quay.io/frankgreco/tester"
	SecurePort   = 8443
	InsecurePort = 8080
)

type deployConfig struct {
	name       string
	namespace  string
	serverType deploy.TLSType
	clientset  kubernetes.Interface
}

type option func(*deployConfig)

func Deploy(
	name string,
	namespace string,
	cli kubernetes.Interface,
	options ...deploy.Option,
) (dns, ip string, err error) {
	cfg := &deployConfig{
		name:      name,
		namespace: namespace,
		clientset: cli,
	}
	for _, op := range options {
		op(cfg)
	}
	return cfg.deploy()
}

func (cfg *deployConfig) deploy() (dns, ip string, err error) {
	svc, err := cfg.deployService()
	if err != nil {
		return "", "", err
	}

	_, err = cfg.deploySecret(svc)
	if err != nil {
		return "", "", err
	}

	_, err = cfg.deployPod()
	if err != nil {
		return "", "", err
	}

	if cfg.serverType == deploy.TLSTypeNone {
		return fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", svc.GetName(), svc.GetNamespace(), InsecurePort), fmt.Sprintf("http://%s:%d", svc.Spec.ClusterIP, InsecurePort), nil
	}

	return fmt.Sprintf("https://%s.%s.svc.cluster.local:%d", svc.GetName(), svc.GetNamespace(), SecurePort), fmt.Sprintf("http://%s:%d", svc.Spec.ClusterIP, SecurePort), nil
}

func (cfg *deployConfig) deployPod() (*v1.Pod, error) {
	var args []string
	var port int

	insecureArgs := []string{
		fmt.Sprintf("--server.insecure_port=%d", InsecurePort),
	}
	secureArgs := []string{
		fmt.Sprintf("--server.secure_port=%d", SecurePort),
		"--server.tls.cert_file=/etc/pki/tls.crt",
		"--server.tls.key_file=/etc/pki/tls.key",
	}

	volumes := []v1.Volume{
		{
			Name: "pki",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: Name,
				},
			},
		},
	}

	mounts := []v1.VolumeMount{
		{
			Name:      "pki",
			MountPath: "/etc/pki",
		},
	}

	switch cfg.serverType {
	case deploy.TLSTypeNone:
		args = insecureArgs
		port = InsecurePort
		volumes, mounts = nil, nil
	case deploy.TLSTypePresent:
		args = secureArgs
		port = SecurePort
	case deploy.TLSTypeMutual:
		args = append(secureArgs,
			"--server.tls.ca_file=/etc/pki/tls.ca",
		)
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
			Name:      Name,
			Namespace: cfg.namespace,
			Labels: map[string]string{
				"app": Name,
			},
		},
		Spec: v1.PodSpec{
			Volumes: volumes,
			Containers: []v1.Container{
				{
					Name:            Name,
					Image:           image,
					VolumeMounts:    mounts,
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

func (cfg *deployConfig) deploySecret(svc *v1.Service) (*v1.Secret, error) {
	tlsAssets := builder.NewTLSBuilder([]string{
		fmt.Sprintf("%s.%s.svc.cluster.local", svc.GetName(), svc.GetNamespace()),
		fmt.Sprintf("%s.%s.svc", svc.GetName(), svc.GetNamespace()),
		fmt.Sprintf("%s.%s", svc.GetName(), svc.GetNamespace()),
		fmt.Sprintf("%s", svc.GetName()),
	}, []net.IP{
		net.ParseIP(svc.Spec.ClusterIP),
	}).NewOrDie()

	return cfg.clientset.CoreV1().Secrets(cfg.namespace).Create(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: cfg.namespace,
			Annotations: map[string]string{
				"kanali.io/enabled": "true",
			},
		},
		Type: v1.SecretTypeTLS,
		Data: map[string][]byte{
			"tls.ca":  tlsAssets.CACert,
			"tls.key": tlsAssets.ServerKey,
			"tls.crt": tlsAssets.ServerCert,
		},
	})
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

	return cfg.clientset.CoreV1().Services(cfg.namespace).Create(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: cfg.namespace,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app": Name,
			},
			Ports: ports,
		},
	})
}

func (cfg *deployConfig) SetServerType(t deploy.TLSType) {
	cfg.serverType = t
}
