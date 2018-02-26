package tester

import (
	"fmt"
	"net"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	"github.com/northwesternmutual/kanali/test/builder"
	"github.com/northwesternmutual/kanali/test/e2e/framework"
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
	serverType tlsType
	clientset  kubernetes.Interface
}

type option func(*deployConfig)

func Deploy(
	name string,
	cli kubernetes.Interface,
	options ...option,
) (dns, ip string, err error) {
	cfg := &deployConfig{
		name:      name,
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

	if cfg.serverType == TLSTypeNone {
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
	case TLSTypeNone:
		args = insecureArgs
		port = InsecurePort
		volumes, mounts = nil, nil
	case TLSTypePresent:
		args = secureArgs
		port = SecurePort
	case TLSTypeMutual:
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

	pod, err := cfg.clientset.CoreV1().Pods(cfg.name).Create(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: cfg.name,
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
	})
	if err != nil {
		return nil, err
	}

	return pod, wait.Poll(framework.Poll, time.Second*30, func() (bool, error) {
		ep, err := cfg.clientset.CoreV1().Endpoints(cfg.name).Get(pod.GetName(), metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if ep.Subsets == nil {
			return false, nil
		}
		for _, ss := range ep.Subsets {
			if ss.Addresses != nil && len(ss.Addresses) > 0 {
				return true, nil
			}
		}
		return false, nil
	})
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

	return cfg.clientset.CoreV1().Secrets(cfg.name).Create(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: cfg.name,
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

	if cfg.serverType == TLSTypeNone {
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

	return cfg.clientset.CoreV1().Services(cfg.name).Create(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: cfg.name,
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

func WithServer(serverType tlsType) option {
	return option(func(cfg *deployConfig) {
		cfg.serverType = serverType
	})
}
