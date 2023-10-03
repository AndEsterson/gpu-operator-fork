/**
# Copyright (c) NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/

package conditions

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nvidiav1 "github.com/NVIDIA/gpu-operator/api/v1"
)

// Specific implementation of the Updater interface for one of our controllers
type clusterPolicyUpdater struct {
	client client.Client
}

// NewClusterPolicyUpdater returns an instance to update conditions for ClusterPolicy
func NewClusterPolicyUpdater(client client.Client) Updater {
	return &clusterPolicyUpdater{client: client}
}

func (u *clusterPolicyUpdater) SetConditionsReady(ctx context.Context, cr any, reason, message string) error {
	clusterPolicyCr, _ := cr.(*nvidiav1.ClusterPolicy)
	return u.setConditionsReady(ctx, clusterPolicyCr, reason, message)
}

func (u *clusterPolicyUpdater) SetConditionsError(ctx context.Context, cr any, reason, message string) error {
	clusterPolicyCr, _ := cr.(*nvidiav1.ClusterPolicy)
	return u.setConditionsError(ctx, clusterPolicyCr, reason, message)
}

func (u *clusterPolicyUpdater) setConditionsReady(ctx context.Context, cr *nvidiav1.ClusterPolicy, reason, message string) error {
	meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
		Type:    Ready,
		Status:  metav1.ConditionTrue,
		Reason:  reason,
		Message: message,
	})

	meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
		Type:   Error,
		Status: metav1.ConditionFalse,
		Reason: Ready,
	})

	return u.client.Status().Update(ctx, cr)
}

func (u *clusterPolicyUpdater) setConditionsError(ctx context.Context, cr *nvidiav1.ClusterPolicy, reason, message string) error {
	meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
		Type:   Ready,
		Status: metav1.ConditionFalse,
		Reason: Error,
	})

	meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
		Type:    Error,
		Status:  metav1.ConditionTrue,
		Reason:  reason,
		Message: message,
	})

	return u.client.Status().Update(ctx, cr)
}
