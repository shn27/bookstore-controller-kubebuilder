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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1alpha1 "github.com/bookstore-controller-kubebuilder/api/v1alpha1"
)

// GetGroupReconciler reconciles a GetGroup object
type GetGroupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.shn.me,resources=getgroups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.shn.me,resources=getgroups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.shn.me,resources=getgroups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GetGroup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *GetGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	fmt.Println("==================================================")
	fmt.Println(req.NamespacedName)
	klog.Info("Reconciling getgroup")
	_ = log.FromContext(ctx)

	bookstore := &appv1alpha1.GetGroup{}
	err := r.Get(ctx, req.NamespacedName, bookstore)
	if err != nil {
		if errors.IsNotFound(err) {
			//log.Info("BookStore resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		//log.Error(err, "Failed to get BookStore")
		return ctrl.Result{}, err
	}

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GetGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.GetGroup{}).
		//	Owns(&appv1alpha1.Deployment{})
		Complete(r)
}
