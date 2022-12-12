// Copyright 2014 Google Inc. All rights reserved.
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

package control

import (
	"novmm/machine"
	"novmm/platform"
)

//
// State.
//

type State struct {
	// Our device state.
	// Note that we only encode state associated with
	// specific devices. The model type is a generic wrapped
	// around devices which *may* encode additional state,
	// but all the state associated with the model should be
	// regenerated on startup.
	Devices []machine.DeviceInfo `json:"devices,omitempty"`

	// Our vcpu state.
	// Similarly, this should encode all the state associated
	// with the underlying VM. If we have other internal platform
	// devices (such as APICs or PITs) then these should be somehow
	// encoded as generic devices.
	Vcpus []platform.VcpuInfo `json:"vcpus,omitempty"`
}

func SaveState(vm *platform.Vm, model *machine.Model) (State, error) {

	// Pause the vm.
	// NOTE: Our model will also be stopped automatically
	// with model.DeviceInfo() below, but we manually pause
	// the Vcpus here in order to ensure they are completely
	// stopped prior to saving all device state.
	err := vm.Pause(false)
	if err != nil {
		return State{}, err
	}
	defer vm.Unpause(false)

	// Grab our vcpu states.
	vcpus, err := vm.VcpuInfo()
	if err != nil {
		return State{}, err
	}

	// Grab our devices.
	// NOTE: This should block until devices have
	// actually quiesed (finished processing outstanding
	// requests generated by the VCPUs).
	devices, err := model.DeviceInfo(vm)
	if err != nil {
		return State{}, err
	}

	// Done.
	return State{Vcpus: vcpus, Devices: devices}, nil
}