/*
Automated Archival API

APIs to manipulate automated archival settings and exempt metrics

API version: 1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package automated_archival

import (
	"encoding/json"
)

// checks if the AutomatedArchivalSettingsDeleteRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AutomatedArchivalSettingsDeleteRequest{}

// AutomatedArchivalSettingsDeleteRequest struct for AutomatedArchivalSettingsDeleteRequest
type AutomatedArchivalSettingsDeleteRequest struct {
	Version *int64 `json:"version,omitempty"`
}

// NewAutomatedArchivalSettingsDeleteRequest instantiates a new AutomatedArchivalSettingsDeleteRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAutomatedArchivalSettingsDeleteRequest() *AutomatedArchivalSettingsDeleteRequest {
	this := AutomatedArchivalSettingsDeleteRequest{}
	return &this
}

// NewAutomatedArchivalSettingsDeleteRequestWithDefaults instantiates a new AutomatedArchivalSettingsDeleteRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAutomatedArchivalSettingsDeleteRequestWithDefaults() *AutomatedArchivalSettingsDeleteRequest {
	this := AutomatedArchivalSettingsDeleteRequest{}
	return &this
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *AutomatedArchivalSettingsDeleteRequest) GetVersion() int64 {
	if o == nil || IsNil(o.Version) {
		var ret int64
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AutomatedArchivalSettingsDeleteRequest) GetVersionOk() (*int64, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *AutomatedArchivalSettingsDeleteRequest) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given int64 and assigns it to the Version field.
func (o *AutomatedArchivalSettingsDeleteRequest) SetVersion(v int64) {
	o.Version = &v
}

func (o AutomatedArchivalSettingsDeleteRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AutomatedArchivalSettingsDeleteRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	return toSerialize, nil
}

type NullableAutomatedArchivalSettingsDeleteRequest struct {
	value *AutomatedArchivalSettingsDeleteRequest
	isSet bool
}

func (v NullableAutomatedArchivalSettingsDeleteRequest) Get() *AutomatedArchivalSettingsDeleteRequest {
	return v.value
}

func (v *NullableAutomatedArchivalSettingsDeleteRequest) Set(val *AutomatedArchivalSettingsDeleteRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableAutomatedArchivalSettingsDeleteRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableAutomatedArchivalSettingsDeleteRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAutomatedArchivalSettingsDeleteRequest(val *AutomatedArchivalSettingsDeleteRequest) *NullableAutomatedArchivalSettingsDeleteRequest {
	return &NullableAutomatedArchivalSettingsDeleteRequest{value: val, isSet: true}
}

func (v NullableAutomatedArchivalSettingsDeleteRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAutomatedArchivalSettingsDeleteRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


