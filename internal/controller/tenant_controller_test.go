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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	platformv1 "github.com/jimmyjoy/gateway-operator/api/v1"
)

var _ = Describe("Tenant Controller", func() {

	Context("When reconciling a Tenant resource", func() {

		const (
			resourceName = "payments"
		)

		ctx := context.Background()

		resourceKey := types.NamespacedName{
			Name: resourceName,
		}

		tenant := &platformv1.Tenant{}

		BeforeEach(func() {

			By("Creating a Tenant resource")

			err := k8sClient.Get(ctx, resourceKey, tenant)
			if err != nil && errors.IsNotFound(err) {

				resource := &platformv1.Tenant{
					ObjectMeta: metav1.ObjectMeta{
						Name: resourceName,
					},
					Spec: platformv1.TenantSpec{
						TeamName:    "Payments",
						CMDBTeamID:  "PAY001",
						Owners:      []string{"jimmy.joy@company.com"},
						Environment: "dev",
					},
				}

				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {

			resource := &platformv1.Tenant{}

			err := k8sClient.Get(ctx, resourceKey, resource)
			if err == nil {
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("should successfully reconcile the Tenant", func() {

			By("Reconciling the Tenant resource")

			reconciler := &TenantReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: resourceKey,
			})

			Expect(err).NotTo(HaveOccurred())

			// Namespace assertions will be added in the next milestone.
		})
	})
})