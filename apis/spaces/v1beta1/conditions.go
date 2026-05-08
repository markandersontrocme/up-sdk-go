// Copyright 2023 Upbound Inc.
// All rights reserved

package v1beta1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

const (
	// ConditionTypeHealthy indicates that the control plane is healthy.
	ConditionTypeHealthy xpv2.ConditionType = "Healthy"
	// ReasonHealthy indicates that the control plane is healthy.
	ReasonHealthy xpv2.ConditionReason = "HealthyControlPlane"
	// ReasonUnhealthy indicates that the control plane is unhealthy.
	ReasonUnhealthy xpv2.ConditionReason = "UnhealthyControlPlane"

	// ConditionTypeControlPlaneProvisioned indicates that the control plane is provisioned.
	ConditionTypeControlPlaneProvisioned xpv2.ConditionType = "ControlPlaneProvisioned"
	// ReasonProvisioned indicates that the control plane is provisioned.
	ReasonProvisioned xpv2.ConditionReason = "Provisioned"
	// ReasonProvisioningError indicates that the control plane provisioning has failed.
	ReasonProvisioningError xpv2.ConditionReason = "ProvisioningError"

	// ConditionTypeSourceSynced indicates that the git source is in sync.
	ConditionTypeSourceSynced xpv2.ConditionType = "SourceSynced"
	// ReasonSourceCompleted indicates that the git sync has been completed.
	ReasonSourceCompleted xpv2.ConditionReason = "Completed"
	// ReasonSourceInProgress indicates that the git sync is still in progress.
	ReasonSourceInProgress xpv2.ConditionReason = "InProgress"

	// ConditionTypeSupported indicates that the control plane is running a
	// supported version of Crossplane.
	ConditionTypeSupported xpv2.ConditionType = "Supported"
	// ReasonSupported indicates that the control plane is running
	// a supported version of Crossplane.
	ReasonSupported xpv2.ConditionReason = "SupportedCrossplaneVersion"
	// ReasonUnsupported indicates that the control plane is running a version
	// of Crossplane that is not supported.
	ReasonUnsupported xpv2.ConditionReason = "UnsupportedCrossplaneVersion"

	// ConditionTypeRestored indicates that the control plane has been restored from backup.
	ConditionTypeRestored xpv2.ConditionType = "Restored"
	// ReasonRestoreCompleted indicates that the control plane has been successfully restored from backup.
	ReasonRestoreCompleted xpv2.ConditionReason = "Completed"
	// ReasonRestoreFailed indicates that the control plane failed to restore from backup.
	ReasonRestoreFailed xpv2.ConditionReason = "Failed"

	// ReasonRestorePending indicates that the control plane restore is pending.
	ReasonRestorePending xpv2.ConditionReason = "RestorePending"

	// ConditionTypeRunning indicates whether the workloads on the Control Plane
	// are running or not.
	ConditionTypeRunning xpv2.ConditionType = "CrossplaneRunning"
	// ReasonPausing indicates that the crossplane and provider workloads are being paused.
	ReasonPausing xpv2.ConditionReason = "Pausing"
	// ReasonPaused indicates that the crossplane and provider workloads have been paused.
	ReasonPaused xpv2.ConditionReason = "Paused"
	// ReasonStarting indicates that the crossplane and provider workloads are being started.
	ReasonStarting xpv2.ConditionReason = "Starting"
	// ReasonStarted indicates that the crossplane and provider workloads have been started.
	ReasonStarted xpv2.ConditionReason = "Started"
)

// Healthy returns a condition that indicates the control plane is healthy.
func Healthy() xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeHealthy,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonHealthy,
	}
}

// Unhealthy returns a condition that indicates the control plane is unhealthy.
func Unhealthy() xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeHealthy,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonUnhealthy,
	}
}

// ControlPlaneProvisioned returns a condition that indicates the control plane
// has been provisioned.
func ControlPlaneProvisioned() xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeControlPlaneProvisioned,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonProvisioned,
	}
}

// ControlPlaneProvisionInProgress returns a condition that indicates the control
// plane is still being provisioned.
func ControlPlaneProvisionInProgress() xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeControlPlaneProvisioned,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonProvisioned,
	}
}

// ControlPlaneProvisioningError returns a condition that indicates the control
// plane provisioning has failed.
func ControlPlaneProvisioningError(err error) xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeControlPlaneProvisioned,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonProvisioningError,
		Message:            err.Error(),
	}
}

// SourceSynced returns a condition that indicates the control plane is in sync
// with the source.
func SourceSynced(revision string) xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeSourceSynced,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonSourceCompleted,
		Message:            fmt.Sprintf("In sync with the revision %s", revision),
	}
}

// SourceInProgress returns a condition that indicates the control plane is still
// processing resources coming from the source.
func SourceInProgress(revision string) xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeSourceSynced,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonSourceInProgress,
		Message:            fmt.Sprintf("Syncing revision %s", revision),
	}
}

// SourceError returns a condition that indicates the source operation of the
// control plane has failed.
func SourceError(err error) xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeSourceSynced,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonSourceInProgress,
		Message:            err.Error(),
	}
}

// SupportedCrossplaneVersion returns a condition that indicates that the
// control plane is running a supported version of Crossplane.
func SupportedCrossplaneVersion() xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeSupported,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonSupported,
	}
}

// UnsupportedCrossplaneVersion returns a condition that indicates that the
// control plane is running an unsupported version of Crossplane.
func UnsupportedCrossplaneVersion(msg string) xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeSupported,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonUnsupported,
		Message:            msg,
	}
}

// RestoreCompleted returns a condition that indicates that the control plane has been
// restored from backup.
func RestoreCompleted() xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeRestored,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonRestoreCompleted,
		Message:            "Control plane has been restored from specified backup",
	}
}

// RestoreFailed returns a condition that indicates that the control plane failed
// to restore from backup.
func RestoreFailed(err error) xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeRestored,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonRestoreFailed,
		Message:            err.Error(),
	}
}

// RestorePending returns a condition that indicates that the control plane restore
// is pending.
func RestorePending() xpv2.Condition {
	return xpv2.Condition{
		Type:               xpv2.TypeReady,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonRestorePending,
		Message:            "Control plane restore is pending",
	}
}

// PauseInProgress returns a condition that indicates that the crossplane and
// provider workloads are being paused.
func PauseInProgress() xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeRunning,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonPausing,
		Message:            "The crossplane and provider workloads are being paused",
	}
}

// PauseCompleted returns a condition that indicates that the crossplane and
// provider workloads have been paused.
func PauseCompleted() xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeRunning,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonPaused,
		Message:            "The crossplane and provider workloads have been paused",
	}
}

// StartInProgress returns a condition that indicates that the crossplane and
// provider workloads are being restarted.
func StartInProgress() xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeRunning,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonStarting,
		Message:            "The crossplane and provider workloads are being started",
	}
}

// StartCompleted returns a condition that indicates that the crossplane and
// provider workloads have been restarted.
func StartCompleted() xpv2.Condition {
	return xpv2.Condition{
		Type:               ConditionTypeRunning,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonStarted,
	}
}
