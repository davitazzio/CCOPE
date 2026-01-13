/*
Copyright 2022 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flpipeline

import (
	"context"
	"fmt"
	"k8s.io/client-go/rest"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"k8s.io/client-go/kubernetes"

	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane/provider-providerfederatedpipeline/apis/federatedlearning/v1alpha1"
	apisv1alpha1 "github.com/crossplane/provider-providerfederatedpipeline/apis/v1alpha1"
	"github.com/crossplane/provider-providerfederatedpipeline/internal/features"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	errNotFLpipeline = "managed resource is not a FLpipeline custom resource"
	errTrackPCUsage  = "cannot track ProviderConfig usage"
	errGetPC         = "cannot get ProviderConfig"
	errGetCreds      = "cannot get credentials"

	errNewClient = "cannot create new Service"
)

// A NoOpService does nothing.
type NoOpService struct{}

var (
	newNoOpService = func(_ []byte) (interface{}, error) { return &NoOpService{}, nil }
)

// Setup adds a controller that reconciles FLpipeline managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.FLpipelineGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.FLpipelineGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			logger:       o.Logger,
			newServiceFn: newNoOpService}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithPollInterval(o.PollInterval),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		WithEventFilter(resource.DesiredStateChanged()).
		For(&v1alpha1.FLpipeline{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	logger       logging.Logger
	newServiceFn func(creds []byte) (interface{}, error)
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.FLpipeline)
	if !ok {
		return nil, errors.New(errNotFLpipeline)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	cd := pc.Spec.Credentials
	data, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.kube, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	svc, err := c.newServiceFn(data)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	return &external{service: svc, logger: c.logger}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	// A 'client' used to connect to the external resource API. In practice this
	// would be something like an AWS SDK client.
	service interface{}
	logger  logging.Logger
}

type client_data struct {
	//data of the client
	epochs     int
	batch_size int
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.FLpipeline)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotFLpipeline)
	}
	//dictionary of clients data
	// clients_data := make(map[string]client_data)
	//get the clientset
	address := cr.Spec.ForProvider.ClusterAddress
	

	c.logger.Debug(fmt.Sprintf("Observing: %+v", cr.Name))
	c.logger.Debug(fmt.Sprintf("connecting to %s", address))
	clientset, err := c.connect_kube(address)
	if err != nil {
		c.logger.Debug("Error in connecting to kubernetes")
	}
	//get the topic queue
	degraded, err := c.get_topic_queue(clientset, cr.Spec.ForProvider.TopicName)
	if err != nil {
		c.logger.Debug("Error in getting topic status")
	}
	c.logger.Debug(fmt.Sprintf("degraged: %t", degraded))
	//get the clients data
	flclients_list, err := clientset.Resource(schema.GroupVersionResource{
		Group:    "federatedlearning.providerflclient.crossplane.io",
		Version:  "v1alpha1",
		Resource: "flclients",
	}).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.logger.Debug("Error in getting flclients list")
		c.logger.Debug(err.Error())
	} else {
		for _, flclient := range flclients_list.Items {
			flclient_name := flclient.GetName()
			c.logger.Debug(fmt.Sprintf("flclient name: %s", flclient_name))
			flclient_data, err := c.get_flclient_data(clientset, flclient_name)
			if err != nil {
				c.logger.Debug(fmt.Sprintf("Error in getting %s data", flclient_name))
			} else {
				c.logger.Debug(fmt.Sprintf("%s data: %+v", flclient_name, flclient_data))
			}
		}
	}

	//scale the deployment of flclients
	if degraded {
		config := &rest.Config{
			Host: address,
		}
		clientset_kube, err := kubernetes.NewForConfig(config)
		if err != nil {
			c.logger.Debug("Error in creating clientset kubernetes")
			c.logger.Debug(err.Error())
		}

		c.logger.Debug("Scaling the deployment")
		//get the deployment
		deployment, err := clientset_kube.AppsV1().Deployments("default").Get(context.TODO(), "flclient-deployment", metav1.GetOptions{})
		if err != nil {
			c.logger.Debug("Error in getting deployment")
			c.logger.Debug(err.Error())
		} else {
			//get the replicas
			replicas := *deployment.Spec.Replicas
			c.logger.Debug(fmt.Sprintf("replicas: %d", replicas))
			if replicas < 5 {
				replicas++
				deployment.Spec.Replicas = &replicas
				_, err := clientset_kube.AppsV1().Deployments("default").Update(context.TODO(), deployment, metav1.UpdateOptions{})
				if err != nil {
					c.logger.Debug("Error in scaling deployment")
					c.logger.Debug(err.Error())
				} else {
					c.logger.Debug("Deployment scaled")
				}
			}else{
			c.logger.Debug("Deployment maximaly scaled")
		}
		}
	}

	
	


	return managed.ExternalObservation{
		// Return false when the external resource does not exist. This lets
		// the managed resource reconciler know that it needs to call Create to
		// (re)create the resource, or that it has successfully been deleted.
		ResourceExists: true,

		// Return false when the external resource exists, but it not up to date
		// with the desired managed resource state. This lets the managed
		// resource reconciler know that it needs to call Update.
		ResourceUpToDate: true,

		// Return any details that may be required to connect to the external
		// resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.FLpipeline)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotFLpipeline)
	}

	fmt.Printf("Creating: %+v", cr)

	return managed.ExternalCreation{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.FLpipeline)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotFLpipeline)
	}

	fmt.Printf("Updating: %+v", cr)

	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.FLpipeline)
	if !ok {
		return errors.New(errNotFLpipeline)
	}

	fmt.Printf("Deleting: %+v", cr)

	return nil
}
