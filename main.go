package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/mattias-fjellstrom/gitops-cdk8s-argocd/imports/k8s"
)

type WebsiteChartProps struct {
	cdk8s.ChartProps
}

func NewWebsiteChart(scope constructs.Construct, id string, props *WebsiteChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	label := map[string]*string{"app": jsii.String("hello-k8s")}

	k8s.NewKubeService(chart, jsii.String("service"), &k8s.KubeServiceProps{
		Spec: &k8s.ServiceSpec{
			Type: jsii.String("LoadBalancer"),
			Ports: &[]*k8s.ServicePort{{
				Port: jsii.Number(80),
				TargetPort: k8s.IntOrString_FromNumber(jsii.Number(80)),
			}},
			Selector: &label,
		},
	})

	cm := k8s.NewKubeConfigMap(chart, jsii.String("index.html"), &k8s.KubeConfigMapProps{
		Data: &map[string]*string{
			"index.html": jsii.String("<html><h1>Version 1</h1></html"),
		},
	})

	volName := jsii.String("nginx-index-file")

	k8s.NewKubeDeployment(chart, jsii.String("deployment"), &k8s.KubeDeploymentProps{
		Spec: &k8s.DeploymentSpec{
			Replicas: jsii.Number(3),
			Selector: &k8s.LabelSelector{
				MatchLabels: &label,
			},
			Template: &k8s.PodTemplateSpec{
				Metadata: &k8s.ObjectMeta{
					Labels: &label,
				},
				Spec: &k8s.PodSpec{
					Containers: &[]*k8s.Container{{
						Name: jsii.String("webserver"),
						Image: jsii.String("nginx:1.23.2"),
						Ports: &[]*k8s.ContainerPort{{ContainerPort: jsii.Number(8080)}},
						VolumeMounts: &[]*k8s.VolumeMount{
							{
								Name: volName,
								MountPath: jsii.String("/usr/share/nginx/html/"),
							},
						},
					}},
					Volumes: &[]*k8s.Volume{
						{
							Name: volName,
							ConfigMap: &k8s.ConfigMapVolumeSource{
								Name: cm.Name(),
							},
						},
					},
				},
			},
		},
	})

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewWebsiteChart(app, "gitops-cdk8s-argocd", nil)
	app.Synth()
}
