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
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/upbound/up-sdk-go/apis/common"
)

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Source",type="string",JSONPath=".spec.source.name"
// +kubebuilder:printcolumn:name="Config",type="string",JSONPath=".spec.configRef.name"
// +kubebuilder:printcolumn:name="Restored",type="integer",JSONPath=".status.controlPlanes.restored"
// +kubebuilder:printcolumn:name="Failed",type="integer",JSONPath=".status.controlPlanes.failed"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories=spaces

// SpaceRestore represents a restore of a Space from a SpaceBackup.
type SpaceRestore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpaceRestoreSpec   `json:"spec"`
	Status SpaceRestoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SpaceRestoreList contains a list of SpaceRestores.
type SpaceRestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []SpaceRestore `json:"items"`
}

// SpaceRestoreSpec defines a restore of a Space from a SpaceBackup.
// A SpaceRestore is a one-off operation, so its spec can't be changed after
// creation. Create a new SpaceRestore to restore again.
// +kubebuilder:validation:XValidation:rule="has(self.controlPlanes) == has(oldSelf.controlPlanes)",message="controlPlanes is immutable"
type SpaceRestoreSpec struct {
	// Source of the SpaceBackup to restore from.
	// ApiGroup is optional and defaults to "admin.spaces.upbound.io".
	// Kind is required, and the only supported value is "SpaceBackup" at the
	// moment.
	// Name is required. The referenced SpaceBackup doesn't need to exist in
	// the Space, e.g. when restoring into a new Space, as the name is used to
	// locate the backup in the object storage configured by ConfigRef.
	// +kubebuilder:validation:XValidation:rule="(!has(self.apiGroup) || self.apiGroup == 'admin.spaces.upbound.io') && self.kind == 'SpaceBackup'",message="source must be a reference to a SpaceBackup"
	// +kubebuilder:validation:XValidation:rule="size(self.name) > 0",message="source must have a name"
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="source is immutable"
	Source common.TypedLocalObjectReference `json:"source"`

	// ConfigRef is a reference to the space backup configuration pointing to
	// the object storage the backup was stored to.
	// ApiGroup is optional and defaults to "admin.spaces.upbound.io".
	// Kind is required, and the only supported value is "SpaceBackupConfig" at
	// the moment.
	// Name is required.
	// +kubebuilder:validation:XValidation:rule="(!has(self.apiGroup) || self.apiGroup == 'admin.spaces.upbound.io') && self.kind == 'SpaceBackupConfig'",message="backup config ref must be a reference to a SpaceBackupConfig"
	// +kubebuilder:validation:XValidation:rule="size(self.name) > 0",message="backup config ref must have a name"
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="configRef is immutable"
	ConfigRef common.TypedLocalObjectReference `json:"configRef"`

	// ControlPlanes lists the control planes to restore, as the
	// --controlplanes flag of hyperspace restore. If not set, all the control
	// planes in the backup are restored.
	// +optional
	// +listType=map
	// +listMapKey=group
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=1000
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="controlPlanes is immutable"
	ControlPlanes []SpaceRestoreControlPlaneReference `json:"controlPlanes,omitempty"`

	// ControlPlanesOnly specifies whether only ControlPlanes should be
	// restored, skipping all other Space resources, e.g. Groups and Secrets,
	// which are expected to already exist, as the --controlplanes-only flag of
	// hyperspace restore.
	// +optional
	// +kubebuilder:default=false
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="controlPlanesOnly is immutable"
	ControlPlanesOnly bool `json:"controlPlanesOnly,omitempty"`

	// MaxParallelControlPlanes is the maximum number of ControlPlanes to
	// restore in parallel, as the --max-parallel-controlplanes flag of
	// hyperspace restore.
	// +optional
	// +kubebuilder:default=5
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="maxParallelControlPlanes is immutable"
	MaxParallelControlPlanes *int32 `json:"maxParallelControlPlanes,omitempty"`

	// ControlPlaneTimeout is how long a single ControlPlane can take to become
	// ready and healthy and get its state restored. A ControlPlane exceeding
	// it counts as failed. As the --controlplane-timeout flag of hyperspace
	// restore.
	// +optional
	// +kubebuilder:default="1h"
	// +kubebuilder:validation:XValidation:rule="duration(self) > duration('0s')",message="controlPlaneTimeout must be a positive duration"
	// +kubebuilder:validation:XValidation:rule="duration(self) == duration(oldSelf)",message="controlPlaneTimeout is immutable"
	ControlPlaneTimeout *metav1.Duration `json:"controlPlaneTimeout,omitempty"`

	// Failures defines the failure tolerance for the restore.
	// +optional
	// +kubebuilder:default={}
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="failures is immutable"
	Failures SpaceRestoreFailuresConfig `json:"failures,omitempty"`
}

// SpaceRestoreFailuresConfig defines the failure tolerance for a SpaceRestore.
type SpaceRestoreFailuresConfig struct {
	// ControlPlanes is the number of control planes that are allowed to fail,
	// including timing out, and still consider the restore successful.
	// Can be specified as an integer (e.g., 5) or as a percentage of the
	// control planes selected for restore (e.g., "10%").
	// If not set, any control plane failing to restore fails the restore.
	// Once exceeded, no new control plane restore is started, the ones in
	// flight are completed, and the restore fails. As the
	// --tolerated-controlplane-failures flag of hyperspace restore.
	// +optional
	// +kubebuilder:validation:XValidation:rule="type(self) == int ? self >= 0 && self <= 2147483647 : self.matches('^(100|[1-9]?[0-9])%$')",message="must be an integer between 0 and 2147483647, or a percentage between 0% and 100%"
	ControlPlanes *intstr.IntOrString `json:"controlPlanes,omitempty"`
}

// SpaceRestorePhase is a string representation of the phase of a restore.
type SpaceRestorePhase string

const (
	// SpaceRestorePhasePending means the restore has been accepted by the
	// system, but it's not started yet, e.g. because another restore is in
	// progress.
	SpaceRestorePhasePending SpaceRestorePhase = "Pending"
	// SpaceRestorePhaseInProgress means the restore is currently being
	// processed.
	SpaceRestorePhaseInProgress SpaceRestorePhase = "InProgress"
	// SpaceRestorePhaseCompleted means the restore has been completed.
	SpaceRestorePhaseCompleted SpaceRestorePhase = "Completed"
	// SpaceRestorePhaseFailed means the restore has failed.
	SpaceRestorePhaseFailed SpaceRestorePhase = "Failed"
)

// SpaceRestoreStatus represents the observed state of a SpaceRestore.
type SpaceRestoreStatus struct {
	xpv2.ManagedResourceStatus `json:",inline"`

	// Phase is the current phase of the restore.
	// +kubebuilder:validation:Enum=Pending;InProgress;Failed;Completed
	// +kubebuilder:default=Pending
	Phase SpaceRestorePhase `json:"phase,omitempty"`

	// StartedAt is the time at which the restore started.
	// +optional
	StartedAt *metav1.Time `json:"startedAt,omitempty"`

	// FinishedAt is the time at which the restore completed or failed.
	// +optional
	FinishedAt *metav1.Time `json:"finishedAt,omitempty"`

	// JobUID is the UID of the Job running the restore, recorded once
	// created.
	// +optional
	JobUID types.UID `json:"jobUID,omitempty"`

	// ControlPlanes contains details about the control planes that were
	// restored, total number of control planes and failures.
	// +optional
	ControlPlanes SpaceRestoreControlPlanesStatus `json:"controlPlanes,omitempty"`
}

// SpaceRestoreControlPlanesStatus contains details about the control planes
// that were restored.
type SpaceRestoreControlPlanesStatus struct {
	// Total is the total number of control planes selected for restore.
	// +optional
	Total int32 `json:"total,omitempty"`
	// Restored is the number of control planes that were restored.
	// +optional
	Restored int32 `json:"restored,omitempty"`
	// Skipped is the number of control planes that were skipped, e.g. because
	// the backup doesn't contain their state.
	// +optional
	Skipped int32 `json:"skipped,omitempty"`
	// Failed is the number of control planes that failed to restore.
	// +optional
	Failed int32 `json:"failed,omitempty"`
	// Unrestored lists the control planes that were skipped or failed to
	// restore, and why. It's capped to keep the object within size limits,
	// the restore logs having the full list.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=100
	Unrestored []SpaceRestoreControlPlaneOutcome `json:"unrestored,omitempty"`
}

// SpaceRestoreControlPlaneReference identifies a control plane in a backup.
type SpaceRestoreControlPlaneReference struct {
	// Group is the group of the control plane.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern="^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
	Group string `json:"group"`
	// Name is the name of the control plane.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern="^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"
	Name string `json:"name"`
}

// SpaceRestoreOutcome is the outcome of a control plane restore.
type SpaceRestoreOutcome string

const (
	// SpaceRestoreOutcomeSkipped means the control plane wasn't restored, e.g.
	// because the backup doesn't contain its state.
	SpaceRestoreOutcomeSkipped SpaceRestoreOutcome = "Skipped"
	// SpaceRestoreOutcomeFailed means the control plane failed to restore.
	SpaceRestoreOutcomeFailed SpaceRestoreOutcome = "Failed"
)

// SpaceRestoreControlPlaneOutcome is the outcome of a control plane restore.
type SpaceRestoreControlPlaneOutcome struct {
	SpaceRestoreControlPlaneReference `json:",inline"`

	// Outcome of the control plane restore.
	// +kubebuilder:validation:Enum=Skipped;Failed
	Outcome SpaceRestoreOutcome `json:"outcome"`
	// Message explains the outcome.
	// +optional
	// +kubebuilder:validation:MaxLength=1024
	Message string `json:"message,omitempty"`
}

// Limits on the size of a SpaceRestore. Each must match the marker on the
// field it applies to.
const (
	// SpaceRestoreMaxControlPlanes is the maximum number of control planes in
	// spec.controlPlanes. They're passed to hyperspace restore as a single
	// argument, which has to stay within the Linux limit on the size of one
	// argument.
	SpaceRestoreMaxControlPlanes = 1000
	// SpaceRestoreMaxUnrestoredControlPlanes is the maximum number of control
	// planes listed in status.controlPlanes.unrestored.
	SpaceRestoreMaxUnrestoredControlPlanes = 100
	// SpaceRestoreMaxMessageLength is the maximum length, in characters, of
	// the message of a control plane in status.controlPlanes.unrestored.
	SpaceRestoreMaxMessageLength = 1024
)

var (
	// SpaceRestoreKind is the kind of a SpaceRestore.
	SpaceRestoreKind = reflect.TypeFor[SpaceRestore]().Name()
	// SpaceRestoreListKind is the kind of a SpaceRestoreList.
	SpaceRestoreListKind = reflect.TypeFor[SpaceRestoreList]().Name()
)

func init() {
	SchemeBuilder.Register(&SpaceRestore{}, &SpaceRestoreList{})
}
