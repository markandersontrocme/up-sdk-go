// Copyright 2026 Upbound Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

// Condition types for SpaceRestore resource.
const (
	// TypeRestoreSucceeded indicates whether the restore succeeded. It's
	// Unknown while the restore is pending or in progress, the reason telling
	// what it's waiting for or restoring, and True or False once it finished.
	TypeRestoreSucceeded xpv2.ConditionType = "Succeeded"
)

// Reasons for the Succeeded condition of a SpaceRestore.
const (
	ReasonWaitingForSpaceRestore                xpv2.ConditionReason = "WaitingForSpaceRestore"
	ReasonWaitingForSpaceBackupConfig           xpv2.ConditionReason = "WaitingForSpaceBackupConfig"
	ReasonRestoringSpaceResources               xpv2.ConditionReason = "RestoringSpaceResources"
	ReasonRestoringControlPlanes                xpv2.ConditionReason = "RestoringControlPlanes"
	ReasonRestoreCompleted                      xpv2.ConditionReason = "RestoreCompleted"
	ReasonRestoreCompletedWithToleratedFailures xpv2.ConditionReason = "RestoreCompletedWithToleratedFailures"
	ReasonRestoreFailed                         xpv2.ConditionReason = "RestoreFailed"
)

// SpaceRestoreWaitingForSpaceRestore indicates that the restore is waiting for
// another SpaceRestore to finish before starting.
func SpaceRestoreWaitingForSpaceRestore(name string) xpv2.Condition {
	return xpv2.Condition{
		Type:               TypeRestoreSucceeded,
		Status:             corev1.ConditionUnknown,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonWaitingForSpaceRestore,
		Message:            fmt.Sprintf("Waiting for SpaceRestore %q to finish.", name),
	}
}

// SpaceRestoreWaitingForSpaceBackupConfig indicates that the restore is
// waiting for its SpaceBackupConfig to exist before starting.
func SpaceRestoreWaitingForSpaceBackupConfig(name string) xpv2.Condition {
	return xpv2.Condition{
		Type:               TypeRestoreSucceeded,
		Status:             corev1.ConditionUnknown,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonWaitingForSpaceBackupConfig,
		Message:            fmt.Sprintf("Waiting for SpaceBackupConfig %q to exist.", name),
	}
}

// SpaceRestoreRestoringSpaceResources indicates that the restore is restoring
// the Space resources, e.g. Groups and Secrets, before the control planes.
func SpaceRestoreRestoringSpaceResources() xpv2.Condition {
	return xpv2.Condition{
		Type:               TypeRestoreSucceeded,
		Status:             corev1.ConditionUnknown,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonRestoringSpaceResources,
		Message:            "Restoring Space resources.",
	}
}

// SpaceRestoreRestoringControlPlanes indicates that the restore is restoring
// the control planes.
func SpaceRestoreRestoringControlPlanes() xpv2.Condition {
	return xpv2.Condition{
		Type:               TypeRestoreSucceeded,
		Status:             corev1.ConditionUnknown,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonRestoringControlPlanes,
		Message:            "Restoring control planes.",
	}
}

// SpaceRestoreCompleted indicates that the restore completed with no control
// plane failing to restore.
func SpaceRestoreCompleted() xpv2.Condition {
	return xpv2.Condition{
		Type:               TypeRestoreSucceeded,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonRestoreCompleted,
		Message:            "The restore completed.",
	}
}

// SpaceRestoreCompletedWithToleratedFailures indicates that the restore
// completed, with some control planes failing to restore, within the tolerated
// failures.
func SpaceRestoreCompletedWithToleratedFailures(total, failed int32) xpv2.Condition {
	return xpv2.Condition{
		Type:               TypeRestoreSucceeded,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonRestoreCompletedWithToleratedFailures,
		Message:            fmt.Sprintf("%d/%d control planes failed to restore, within the tolerated failures.", failed, total),
	}
}

// SpaceRestoreFailed indicates that the restore failed.
func SpaceRestoreFailed(message string) xpv2.Condition {
	return xpv2.Condition{
		Type:               TypeRestoreSucceeded,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonRestoreFailed,
		Message:            message,
	}
}
