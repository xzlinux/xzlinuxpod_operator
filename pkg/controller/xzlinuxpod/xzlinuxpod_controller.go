package xzlinuxpod

import (
	"context"
	"reflect"
	k8sv1alpha1 "github.com/xzlinux/xzlinuxpod-operator/pkg/apis/k8s/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

)

var log = logf.Log.WithName("controller_xzlinuxpod")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new XzlinuxPod Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileXzlinuxPod{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("xzlinuxpod-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource XzlinuxPod
	err = c.Watch(&source.Kind{Type: &k8sv1alpha1.XzlinuxPod{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner XzlinuxPod
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.XzlinuxPod{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileXzlinuxPod implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileXzlinuxPod{}

// ReconcileXzlinuxPod reconciles a XzlinuxPod object
type ReconcileXzlinuxPod struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a XzlinuxPod object and makes changes based on the state read
// and what is in the XzlinuxPod.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileXzlinuxPod) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling XzlinuxPod")

	// Fetch the XzlinuxPod instance
	instance := &k8sv1alpha1.XzlinuxPod{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	//1
    lbls := labels.Set{
		"app": instance.Name,
	}
	existingPods := &corev1.PodList{}
	err = r.client.List(context.TODO(),existingPods,&client.ListOptions{
		Namespace: request.Namespace,
		LabelSelector: labels.SelectorFromSet(lbls),
	})
	if err != nil {
		reqLogger.Error(err,"??????????????????pod??????")
		return reconcile.Result{},err 
	}
//2
	var existingPodNames []string
	for _, pod := range existingPods.Items {
		if pod.GetObjectMeta().GetDeletionTimestamp() != nil {
			continue
		}
		if pod.Status.Phase == corev1.PodPending || pod.Status.Phase == corev1.PodRunning {
			existingPodNames = append(existingPodNames,pod.GetObjectMeta().GetName())
		}
	}
//3
	currStatus := k8sv1alpha1.XzlinuxPodStatus{
		PodNames: existingPodNames,
		Replicas: len(existingPodNames),
	}
	if !reflect.DeepEqual(instance.Status,currStatus) {
		instance.Status = currStatus
		if err := r.client.Status().Update(context.TODO(),instance);err != nil {
			reqLogger.Error(err, "Update pod failed")
			return reconcile.Result{},err 
		}

	}
//4
	if len(existingPodNames) > instance.Spec.Replicas {
		//delete
		reqLogger.Info("Delete pod",existingPodNames,instance.Spec.Replicas)
		pod := existingPods.Items[0]
		err := r.client.Delete(context.TODO(),&pod)
		if err != nil {
			reqLogger.Error(err, "Delete pod failed")
			return reconcile.Result{},err
		}
	}
//5
	if len(existingPodNames) < instance.Spec.Replicas {
		reqLogger.Info("create pod: ",existingPodNames,instance.Spec.Replicas)
		// Define a new Pod object
		pod := newPodForCR(instance)

		// Set XzlinuxPod instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			reqLogger.Error(err,"create pod faild")
			return reconcile.Result{}, err
		}
	}


	return reconcile.Result{Requeue: true}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *k8sv1alpha1.XzlinuxPod) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
