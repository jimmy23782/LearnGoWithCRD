/*
Copyright 2026.

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

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	platformv1 "github.com/jimmyjoy/gateway-operator/api/v1"
)

// TenantReconciler reconciles a Tenant object.
type TenantReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Permissions required to read and reconcile Tenant resources.
//
// +kubebuilder:rbac:groups=platform.mac.com,resources=tenants,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=platform.mac.com,resources=tenants/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=platform.mac.com,resources=tenants/finalizers,verbs=update

// Permissions required to create and inspect namespaces.
//
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch

// Reconcile ensures that the desired Kubernetes resources for a Tenant exist.
func (r *TenantReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {

	log := logf.FromContext(ctx)

	// ---------------------------------------------------------------------
	// 1. Fetch the Tenant custom resource.
	// ---------------------------------------------------------------------

	tenant := &platformv1.Tenant{}

	if err := r.Get(ctx, req.NamespacedName, tenant); err != nil {

		if apierrors.IsNotFound(err) {
			// The Tenant CR has been deleted.
			// There is nothing to reconcile.
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	log.Info(
		"Reconciling Tenant",
		"tenant", tenant.Name,
		"environment", tenant.Spec.Environment,
	)

	// ---------------------------------------------------------------------
	// 2. Determine the desired namespace name.
	//
	// Example:
	//
	// Tenant metadata.name = payments
	// Environment          = dev
	//
	// Result:
	// payments-dev
	// ---------------------------------------------------------------------

	namespaceName := fmt.Sprintf(
		"%s-%s",
		tenant.Name,
		tenant.Spec.Environment,
	)

	// ---------------------------------------------------------------------
	// 3. Check whether the namespace already exists.
	// ---------------------------------------------------------------------

	existingNamespace := &corev1.Namespace{}

	err := r.Get(
		ctx,
		types.NamespacedName{
			Name: namespaceName,
		},
		existingNamespace,
	)

	// ---------------------------------------------------------------------
	// 4. Namespace already exists.
	//
	// Desired state has already been achieved, so there is nothing to do.
	// This makes our reconciliation idempotent.
	// ---------------------------------------------------------------------

	if err == nil {

		log.Info(
			"Namespace already exists",
			"namespace", namespaceName,
		)

		return ctrl.Result{}, nil
	}

	// ---------------------------------------------------------------------
	// 5. Return unexpected errors.
	// ---------------------------------------------------------------------

	if !apierrors.IsNotFound(err) {
		return ctrl.Result{}, err
	}

	// ---------------------------------------------------------------------
	// 6. Namespace does not exist, so create it.
	// ---------------------------------------------------------------------

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,

			Labels: map[string]string{
				"platform.mac.com/managed-by":  "gateway-operator",
				"platform.mac.com/environment": tenant.Spec.Environment,
				"platform.mac.com/tenant":      tenant.Name,
			},

			Annotations: map[string]string{
				"platform.mac.com/team-name":    tenant.Spec.TeamName,
				"platform.mac.com/cmdb-team-id": tenant.Spec.CMDBTeamID,
			},
		},
	}

	if err := r.Create(ctx, namespace); err != nil {

		// It is possible another reconciliation created the namespace
		// between our GET and CREATE operations.
		if apierrors.IsAlreadyExists(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	log.Info(
		"Created namespace for Tenant",
		"tenant", tenant.Name,
		"namespace", namespaceName,
	)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TenantReconciler) SetupWithManager(
	mgr ctrl.Manager,
) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1.Tenant{}).
		Named("tenant").
		Complete(r)
}
