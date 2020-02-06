package stavhelloworld

import (
	"context"

	stavv1alpha1 "github.com/stavco9/go-test/pkg/apis/stav/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"github.com/nlopes/slack"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
)

var log = logf.Log.WithName("controller_stavhelloworld")
const hwFinalizer = "finalizer.stavhelloworld.hw.okto.io"

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new StavHelloWorld Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileStavHelloWorld{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("stavhelloworld-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
	  return err
	}

	// Watch for changes to primary resource <YOUR-NAME>HelloWorld
	err = c.Watch(&source.Kind{Type: &stavhelloworldv1alpha1.StavHelloWorld{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
	  return err
	}

	// Watch for changes to Service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
	  IsController: true,
	  OwnerType:    &<YOUR-NAME>v1alpha1.<YOUR-NAME>HelloWorld{},
	})
	if err != nil {
	  return err
	}

	// Watch for changes to ConfigMap
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
	  IsController: true,
	  OwnerType:    &<YOUR-NAME>v1alpha1.<YOUR-NAME>HelloWorld{},
	})
	if err != nil {
	  return err
	}

	// Watch for changes to Route
	err = c.Watch(&source.Kind{Type: &routev1.Route{}}, &handler.EnqueueRequestForOwner{
	  IsController: true,
	  OwnerType:    &<YOUR-NAME>v1alpha1.<YOUR-NAME>HelloWorld{},
	})
	if err != nil {
	  return err
	}

	return nil
}

// Reconcile loop resources managers functions
func (r *ReconcileStavHelloWorld) manageDeployment(hw *stavhelloworldv1alpha1.StavHelloWorld, reqLogger logr.Logger) (*reconcile.Result, error) {
  deployment := &appsv1.Deployment{}
  err := r.client.Get(context.TODO(), types.NamespacedName{Name: hw.Name, Namespace: hw.Namespace}, deployment)
  if err != nil && errors.IsNotFound(err) {
    serverDeployment, err := r.deploymentForWebServer(hw)
    if err != nil {
      reqLogger.Error(err, "error getting server deployment")
      return &reconcile.Result{}, err
    }
    reqLogger.Info("Creating a new server deployment.", "Deployment.Namespace", serverDeployment.Namespace, "Deployment.Name", serverDeployment.Name)
    err = r.client.Create(context.TODO(), serverDeployment)
    if err != nil {
      reqLogger.Error(err, "Failed to create new Server Deployment.", "Deployment.Namespace", serverDeployment.Namespace, "Deployment.Name", serverDeployment.Name)
      return &reconcile.Result{}, err
    }
    return &reconcile.Result{Requeue: true}, nil
  } else if err != nil {
    reqLogger.Error(err, "Failed to get server deployment.")
    return &reconcile.Result{}, err
  }
  return nil, nil
}

func (r *ReconcileStavHelloWorld) manageRoute(hw *stavhelloworldv1alpha1.StavHelloWorld, reqLogger logr.Logger) (*reconcile.Result, error) {
  //Check if route already exists, if not create a new one
  route := &routev1.Route{}
  err := r.client.Get(context.TODO(), types.NamespacedName{Name: hw.Name, Namespace: hw.Namespace}, route)
  if err != nil && errors.IsNotFound(err) {
    serverRoute, err := r.routeForWebServer(hw)
    if err != nil {
      reqLogger.Error(err, "error getting server route")
      return &reconcile.Result{}, err
    }
    reqLogger.Info("Creating a new route.", "Route.Namespace", serverRoute.Namespace, "Router.Name", serverRoute.Name)
    err = r.client.Create(context.TODO(), serverRoute)
    if err != nil {
      reqLogger.Error(err, "Failed to create new Server Route.", "Route.Namespace", serverRoute.Namespace, "Route.Name", serverRoute.Name)
      return &reconcile.Result{}, err
    }
    return &reconcile.Result{Requeue: true}, nil
  } else if err != nil {
    reqLogger.Error(err, "Failed to get server route.")
    return &reconcile.Result{}, err
  }
  return nil, nil
}

func (r *ReconcileStavHelloWorld) manageService(hw *stavhelloworldv1alpha1.StavHelloWorld, reqLogger logr.Logger) (*reconcile.Result, error) {
  service := &corev1.Service{}
  err := r.client.Get(context.TODO(), types.NamespacedName{Name: hw.Name, Namespace: hw.Namespace}, service)
  if err != nil && errors.IsNotFound(err) {
    err := r.serviceForWebServer(hw, service)
    if err != nil {
      reqLogger.Error(err, "error getting server service")
      return &reconcile.Result{}, err
    }
    reqLogger.Info("Creating a new service.", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
    err = r.client.Create(context.TODO(), service)
    if err != nil {
      reqLogger.Error(err, "Failed to create new Server Service.", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
      return &reconcile.Result{}, err
    }
    return &reconcile.Result{Requeue: true}, nil
  } else if err != nil {
    reqLogger.Error(err, "Failed to get server service.")
    return &reconcile.Result{}, err
  } else {
    err := r.serviceForWebServer(hw, service)
    if err != nil {
      reqLogger.Error(err, "error getting server service")
      return &reconcile.Result{}, err
    }
    err = r.client.Update(context.TODO(), service)
    if err != nil {
      reqLogger.Error(err, "Failed to create new Server Service.", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
      return &reconcile.Result{}, err
    }
  }
  return nil, nil
}

func (r *Reconcile<YOUR-NAME>HelloWorld) manageConfigMap(hw *<YOUR-NAME>v1alpha1.<YOUR-NAME>HelloWorld, reqLogger logr.Logger) (*reconcile.Result, error) {
  cm := &corev1.ConfigMap{}
  err := r.client.Get(context.TODO(), types.NamespacedName{Name: hw.Name, Namespace: hw.Namespace}, cm)
  if err != nil && errors.IsNotFound(err) {
    websitesCm, err := r.configMapForWebServer(hw)
    if err != nil {
      reqLogger.Error(err, "error getting websites ConfigMap")
      return &reconcile.Result{}, err
    }
    reqLogger.Info("Creating a new cm.", "ConfigMap.Namespace", websitesCm.Namespace, "ConfigMap.Name", websitesCm.Name)
    err = r.client.Create(context.TODO(), websitesCm)
    if err != nil {
      reqLogger.Error(err, "Failed to create new ConfigMap.", "ConfigMap.Namespace", websitesCm.Namespace, "ConfigMap.Name", websitesCm.Name)
      return &reconcile.Result{}, err
    }
    return &reconcile.Result{Requeue: true}, nil
  } else if err != nil {
    reqLogger.Error(err, "Failed to get configmap.")
    return &reconcile.Result{}, err
  }

  // Check if CM sync is required
  syncRequired, err := r.syncConfigMapForWebServer(hw, cm)
  if err != nil {
    reqLogger.Error(err, "Error during syncing ConfigMap.")
    return &reconcile.Result{}, err
  }
  // If CM website sync required, sync the CM
  if syncRequired {
    err = r.client.Update(context.TODO(), cm)
    if err != nil {
      reqLogger.Error(err, "Error during updating CM")
      return &reconcile.Result{}, err
    }
  }
  return nil, nil
}

// Resources creation functions
func (r *Reconcile<YOUR-NAME>HelloWorld) deploymentForWebServer(hw *<YOUR-NAME>v1alpha1.<YOUR-NAME>HelloWorld) (*appsv1.Deployment, error) {
    var replicas int32
    replicas = 1
    labels := map[string]string{
      "app": hw.Name,
    }
    dep := &appsv1.Deployment{
      ObjectMeta: metav1.ObjectMeta{
        Name:      hw.Name,
        Namespace: hw.Namespace,
        Labels:    labels,
      },
      Spec: appsv1.DeploymentSpec{
        Replicas: &replicas,
        Selector: &metav1.LabelSelector{
          MatchLabels: map[string]string{"app": hw.Name},
        },
        Template: corev1.PodTemplateSpec{
          ObjectMeta: metav1.ObjectMeta{
            Name:   hw.Name,
            Labels: labels,
          },
          Spec: corev1.PodSpec{
            Containers: []corev1.Container{
              {
                Name:            hw.Name,
                Image:           "docker.io/dimssss/nginx-for-ocp:0.1",
                ImagePullPolicy: corev1.PullAlways,
                Ports: []corev1.ContainerPort{
                  {
                    ContainerPort: 8080,
                  },
                },
                VolumeMounts: []corev1.VolumeMount{
                  {
                    Name:      "website",
                    MountPath: "/opt/app-root/src",
                  },
                },
              },
            },
            Volumes: []corev1.Volume{
              {
                Name: "website",
                VolumeSource: corev1.VolumeSource{
                  ConfigMap: &corev1.ConfigMapVolumeSource{
                    LocalObjectReference: corev1.LocalObjectReference{
                      Name: hw.Name,
                    },
                  },
                },
              },
            },
          },
        },
      },
    }
    if err := controllerutil.SetControllerReference(hw, dep, r.scheme); err != nil {
      log.Error(err, "Error set controller reference for server deployment")
      return nil, err
    }
    return dep, nil
  }

func (r *Reconcile<YOUR-NAME>HelloWorld) serviceForWebServer(hw *<YOUR-NAME>v1alpha1.<YOUR-NAME>HelloWorld, service *corev1.Service) error {
  labels := map[string]string{
    "app": hw.Name,
  }
  ports := []corev1.ServicePort{
    {
      Name: "https",
      Port: 8080,
    },
  }
  service.ObjectMeta.Name = hw.Name
  service.ObjectMeta.Namespace = hw.Namespace
  service.ObjectMeta.Labels = labels
  service.Spec.Selector = map[string]string{"app": hw.Name}
  service.Spec.Ports = ports
  if err := controllerutil.SetControllerReference(hw, service, r.scheme); err != nil {
    log.Error(err, "Error set controller reference for server service")
    return err
  }
  return nil
}

func (r *Reconcile<YOUR-NAME>HelloWorld) routeForWebServer(hw *<YOUR-NAME>v1alpha1.<YOUR-NAME>HelloWorld) (*routev1.Route, error) {
  labels := map[string]string{
    "app": hw.Name,
  }
  route := &routev1.Route{
    ObjectMeta: metav1.ObjectMeta{
      Name:      hw.Name,
      Namespace: hw.Namespace,
      Labels:    labels,
    },
    Spec: routev1.RouteSpec{
      TLS: &routev1.TLSConfig{
        Termination: routev1.TLSTerminationEdge,
      },
      To: routev1.RouteTargetReference{
        Kind: "Service",
        Name: hw.Name,
      },
    },
  }
  if err := controllerutil.SetControllerReference(hw, route, r.scheme); err != nil {
    log.Error(err, "Error set controller reference for server route")
    return nil, err
  }
  return route, nil
}

func (r *Reconcile<YOUR-NAME>HelloWorld) configMapForWebServer(hw *<YOUR-NAME>v1alpha1.<YOUR-NAME>HelloWorld) (*corev1.ConfigMap, error) {
  labels := map[string]string{
    "app": hw.Name,
  }
  cm := &corev1.ConfigMap{
    ObjectMeta: metav1.ObjectMeta{
      Name:      hw.Name,
      Namespace: hw.Namespace,
      Labels:    labels,
    },
    Data: map[string]string{"index.html": hw.Spec.Message},
  }

  if err := controllerutil.SetControllerReference(hw, cm, r.scheme); err != nil {
    log.Error(err, "Error set controller reference for configmap")
    return nil, err
  }
  return cm, nil
}

func (r *Reconcile<YOUR-NAME>HelloWorld) syncConfigMapForWebServer(hw *<YOUR-NAME>v1alpha1.<YOUR-NAME>HelloWorld, cm *corev1.ConfigMap) (syncRequired bool, err error) {
  if hw.Spec.Message != cm.Data["index.html"] {
    log.Info("Message in CR spec not the same as in CM, gonna update website cm")
    cm.Data["index.html"] = hw.Spec.Message
    return true, nil
  }
  log.Info("No sync required, the message didn't changed")
  return false, nil
}

// blank assignment to verify that ReconcileStavHelloWorld implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileStavHelloWorld{}

// ReconcileStavHelloWorld reconciles a StavHelloWorld object
type ReconcileStavHelloWorld struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a StavHelloWorld object and makes changes based on the state read
// and what is in the StavHelloWorld.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileStavHelloWorld) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling HelloWorld")

	// Fetch the HelloWorld instance
	hw := &<YOUR-NAME>v1alpha1.<YOUR-NAME>HelloWorld{}
	err := r.client.Get(context.TODO(), request.NamespacedName, hw)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Init finalizers
	err = r.initFinalization(hw, reqLogger)
	if err != nil {
		reqLogger.Error(err, "Failed to initialize finalizer")
		hw.Status.Message = fmt.Sprintf("%v", err)
		if err := r.client.Status().Update(context.TODO(), hw); err != nil {
			reqLogger.Error(err, "Failed to update CR status")
		}
		return reconcile.Result{}, err
	}

	// Check if configmap for websites list already exists, if not create a new one
	reconcileResult, err := r.manageConfigMap(hw, reqLogger)
	if err != nil {
		hw.Status.Message = fmt.Sprintf("%v", err)
		if err := r.client.Status().Update(context.TODO(), hw); err != nil {
			reqLogger.Error(err, "Failed to update CR status")
		}
		return *reconcileResult, err
	} else if err == nil && reconcileResult != nil {
		// In case requeue required
		return *reconcileResult, nil
	}

	//Check if service already exists, if not create a new one
	reconcileResult, err = r.manageService(hw, reqLogger)
	if err != nil {
		hw.Status.Message = fmt.Sprintf("%v", err)
		if err := r.client.Status().Update(context.TODO(), hw); err != nil {
			reqLogger.Error(err, "Failed to update CR status")
		}
		return *reconcileResult, err
	} else if err == nil && reconcileResult != nil {
		// In case requeue required
		return *reconcileResult, nil
	}
	//Check if route already exists, if not create a new one
	reconcileResult, err = r.manageRoute(hw, reqLogger)
	if err != nil {
		hw.Status.Message = fmt.Sprintf("%v", err)
		if err := r.client.Status().Update(context.TODO(), hw); err != nil {
			reqLogger.Error(err, "Failed to update CR status")
		}
		return *reconcileResult, err
	} else if err == nil && reconcileResult != nil {
		// In case requeue required
		return *reconcileResult, nil
	}
	//Check if deployment already exists, if not create a new one
	reconcileResult, err = r.manageDeployment(hw, reqLogger)
	if err != nil {
		hw.Status.Message = fmt.Sprintf("%v", err)
		if err := r.client.Status().Update(context.TODO(), hw); err != nil {
			reqLogger.Error(err, "Failed to update CR status")
		}
		return *reconcileResult, err
	} else if err == nil && reconcileResult != nil {
		// In case requeue required
		return *reconcileResult, nil
	}

	hw.Status.Message = "All good"
	if err := r.client.Status().Update(context.TODO(), hw); err != nil {
		reqLogger.Error(err, "Failed to update CR status")
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileStavHelloWorld) initFinalization(hw *stavhelloworldv1alpha1.StavHelloWorld, reqLogger logr.Logger) error {
  isHwMarkedToBeDeleted := hw.GetDeletionTimestamp() != nil
  if isHwMarkedToBeDeleted {
    if contains(hw.GetFinalizers(), hwFinalizer) {
      // Run finalization logic for hwFinalizer. If the
      // finalization logic fails, don't remove the finalizer so
      // that we can retry during the next reconciliation.
      if err := r.finalizeHw(hw, reqLogger); err != nil {
        reqLogger.Error(err, "Failed to run finalizer")
        return err
      }
      // Remove hwFinalizer. Once all finalizers have been
      // removed, the object will be deleted.
      hw.SetFinalizers(remove(hw.GetFinalizers(), hwFinalizer))
      err := r.client.Update(context.TODO(), hw)
      if err != nil {
        reqLogger.Error(err, "Failed to delete finalizer")
        return err
      }
    }
    return nil
  }

  // Add finalizer for this CR
  if !contains(hw.GetFinalizers(), hwFinalizer) {
    if err := r.addFinalizer(hw, reqLogger); err != nil {
      reqLogger.Error(err, "Failed to add finalizer")
      return err
    }
  }
  return nil
}

func (r *ReconcileStavHelloWorld) finalizeHw(hw *stavhelloworldv1alpha1.StavHelloWorld, reqLogger logr.Logger, ) error {
  slackToken, err := getSlackToken()
  if err != nil {
    reqLogger.Error(err, "Gonna skip finzalize, the error during getting slack token")
    return nil
  }
  api := slack.New(slackToken)
  attachment := slack.Attachment{
    Pretext: "Hello World Operator Finalizer",
    Color:   "danger",
    Footer:  "HelloWorld Operator Finalizer",
    Title:   fmt.Sprintf("WebSite %s gonna be removed from OpenShift Cluster", hw.Name),
  }
  channelID, timestamp, err := api.PostMessage("CQ5EXBM8C", slack.MsgOptionAttachments(attachment))
  if err != nil {
    reqLogger.Error(err, "Failed to send Slack message")
  }
  fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
  reqLogger.Info(fmt.Sprintf("Successfully finalized HW: %s", hw.Name))
  return nil
}

func (r *ReconcileStavHelloWorld) addFinalizer(hw *stavhelloworldv1alpha1.StavHelloWorld, reqLogger logr.Logger) error {
  reqLogger.Info("Adding Finalizer for the Memcached")
  hw.SetFinalizers(append(hw.GetFinalizers(), hwFinalizer))
  // Update CR
  err := r.client.Update(context.TODO(), hw)
  if err != nil {
    reqLogger.Error(err, "Failed to update Rdbc with finalizer")
    return err
  }
  return nil
}

func getSlackToken() (string, error) {
  ns, found := os.LookupEnv("SLACK_TOKEN")
  if !found {
    return "", fmt.Errorf("%s must be set", "SLACK_TOKEN")
  }
  return ns, nil
}

func contains(list []string, s string) bool {
  for _, v := range list {
    if v == s {
      return true
    }
  }
  return false
}

func remove(list []string, s string) []string {
  for i, v := range list {
    if v == s {
      list = append(list[:i], list[i+1:]...)
    }
  }
  return list
}
