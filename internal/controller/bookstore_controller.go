/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1 "github.com/bookstore-controller-kubebuilder/api/v1"
)

// BookStoreReconciler reconciles a BookStore object
type BookStoreReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.shn.me,resources=bookstores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.shn.me,resources=bookstores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.shn.me,resources=bookstores/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BookStore object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *BookStoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	fmt.Println("==================================================")
	fmt.Println(req.NamespacedName)
	klog.Info("Reconciling BookStore")
	log := log.FromContext(ctx)
	// Fetch the Bookstore instance
	// The purpose is check if the Custom Resource for the Kind Bookstore
	// is applied on the cluster if not we return nil to stop the reconciliation
	bookstore := &appv1.BookStore{}
	err := r.Get(ctx, req.NamespacedName, bookstore)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("BookStore resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get BookStore")
		return ctrl.Result{}, err
	}

	foundDeployment := &appsv1.Deployment{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: bookstore.Spec.DeploymentName, Namespace: bookstore.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		fmt.Println("Creating Deployment")
		if err := r.Client.Create(ctx, getDeployment(bookstore)); err != nil {
			if err != nil {
				klog.Info("error while creating deployment %s", err)
				return ctrl.Result{}, err
			}
		} else {
			klog.Info("Deployment Created	", bookstore.Spec.DeploymentName)
		}
	} else if err != nil {
		klog.Info("error fetching deployment %s", err)
	}

	foundService := &corev1.Service{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: bookstore.Spec.DeploymentName, Namespace: bookstore.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		fmt.Println("Creating Service")
		if err := r.Client.Create(ctx, getService(bookstore)); err != nil {
			if err != nil {
				klog.Info("error while creating service %s", err)
				return ctrl.Result{}, err
			}
		} else {
			klog.Info("Service Created	", bookstore.Spec.DeploymentName)
		}
	} else if err != nil {
		klog.Info("error fetching service %s", err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BookStoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&appv1.BookStore{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)

	if err != nil {
		return err
	}

	//err = ctrl.NewControllerManagedBy(mgr).For(&appsv1.Deployment{}).Complete(r)
	//if err != nil {
	//	return err
	//}

	return nil
}
