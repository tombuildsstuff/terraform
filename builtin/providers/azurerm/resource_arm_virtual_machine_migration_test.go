package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestAzureRMVirtualMachineMigrateStateV0ToV1(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		Attributes   map[string]string
		Expected     map[string]string
		Meta         interface{}
	}{
		"v0_1_set_both_true": {
			StateVersion: 0,
			Attributes: map[string]string{
				"os_profile_windows_config.#":                                       "1",
				"os_profile_windows_config.2256145325.additional_unattend_config.#": "0",
				"os_profile_windows_config.2256145325.enable_automatic_upgrades":    "true",
				"os_profile_windows_config.2256145325.provision_vm_agent":           "true",
				"os_profile_windows_config.2256145325.winrm.#":                      "0",
			},
			Expected: map[string]string{
				"os_profile_windows_config.#":                                       "1",
				"os_profile_windows_config.2256145325.additional_unattend_config.#": "0",
				"os_profile_windows_config.2256145325.enable_automatic_upgrades":    "true",
				"os_profile_windows_config.2256145325.provision_vm_agent":           "true",
				"os_profile_windows_config.2256145325.winrm.#":                      "0",
			},
		},
		"v0_1_set_both_false": {
			StateVersion: 0,
			Attributes: map[string]string{
				"os_profile_windows_config.#":                                      "1",
				"os_profile_windows_config.429474957.additional_unattend_config.#": "0",
				"os_profile_windows_config.429474957.enable_automatic_upgrades":    "false",
				"os_profile_windows_config.429474957.provision_vm_agent":           "false",
				"os_profile_windows_config.429474957.winrm.#":                      "0",
			},
			Expected: map[string]string{
				"os_profile_windows_config.#":                                      "1",
				"os_profile_windows_config.429474957.additional_unattend_config.#": "0",
				"os_profile_windows_config.429474957.enable_automatic_upgrades":    "false",
				"os_profile_windows_config.429474957.provision_vm_agent":           "false",
				"os_profile_windows_config.429474957.winrm.#":                      "0",
			},
		},
		"v0_1_set_auto_upgrades_false": {
			StateVersion: 0,
			Attributes: map[string]string{
				"os_profile_windows_config.#":                                     "1",
				"os_profile_windows_config.69840937.additional_unattend_config.#": "0",
				"os_profile_windows_config.69840937.enable_automatic_upgrades":    "false",
				"os_profile_windows_config.69840937.winrm.#":                      "0",
			},
			Expected: map[string]string{
				"os_profile_windows_config.#":                                      "1",
				"os_profile_windows_config.429474957.additional_unattend_config.#": "0",
				"os_profile_windows_config.429474957.enable_automatic_upgrades":    "false",
				"os_profile_windows_config.429474957.provision_vm_agent":           "false",
				"os_profile_windows_config.429474957.winrm.#":                      "0",
			},
		},
		"v0_1_set_auto_upgrades_true": {
			StateVersion: 0,
			Attributes: map[string]string{
				"os_profile_windows_config.#":                                     "1",
				"os_profile_windows_config.69840937.additional_unattend_config.#": "0",
				"os_profile_windows_config.69840937.enable_automatic_upgrades":    "true",
				"os_profile_windows_config.69840937.winrm.#":                      "0",
			},
			Expected: map[string]string{
				"os_profile_windows_config.#":                                       "1",
				"os_profile_windows_config.3590083716.additional_unattend_config.#": "0",
				"os_profile_windows_config.3590083716.enable_automatic_upgrades":    "true",
				"os_profile_windows_config.3590083716.provision_vm_agent":           "false",
				"os_profile_windows_config.3590083716.winrm.#":                      "0",
			},
		},
		"v0_1_set_vm_agent_false": {
			StateVersion: 0,
			Attributes: map[string]string{
				"os_profile_windows_config.#":                                     "1",
				"os_profile_windows_config.69840937.additional_unattend_config.#": "0",
				"os_profile_windows_config.69840937.provision_vm_agent":           "false",
				"os_profile_windows_config.69840937.winrm.#":                      "0",
			},
			Expected: map[string]string{
				"os_profile_windows_config.#":                                      "1",
				"os_profile_windows_config.429474957.additional_unattend_config.#": "0",
				"os_profile_windows_config.429474957.enable_automatic_upgrades":    "false",
				"os_profile_windows_config.429474957.provision_vm_agent":           "false",
				"os_profile_windows_config.429474957.winrm.#":                      "0",
			},
		},
		"v0_1_set_vm_agent_true": {
			StateVersion: 0,
			Attributes: map[string]string{
				"os_profile_windows_config.#":                                     "1",
				"os_profile_windows_config.69840937.additional_unattend_config.#": "0",
				"os_profile_windows_config.69840937.provision_vm_agent":           "true",
				"os_profile_windows_config.69840937.winrm.#":                      "0",
			},
			Expected: map[string]string{
				"os_profile_windows_config.#":                                       "1",
				"os_profile_windows_config.1534614206.additional_unattend_config.#": "0",
				"os_profile_windows_config.1534614206.enable_automatic_upgrades":    "false",
				"os_profile_windows_config.1534614206.provision_vm_agent":           "true",
				"os_profile_windows_config.1534614206.winrm.#":                      "0",
			},
		},
		"v0_1_unset": {
			StateVersion: 0,
			Attributes: map[string]string{
				"os_profile_windows_config.#":                              "1",
				"os_profile_windows_config.0.additional_unattend_config.#": "0",
				"os_profile_windows_config.0.winrm.#":                      "0",
			},
			Expected: map[string]string{
				"os_profile_windows_config.#":                                      "1",
				"os_profile_windows_config.429474957.additional_unattend_config.#": "0",
				"os_profile_windows_config.429474957.enable_automatic_upgrades":    "false",
				"os_profile_windows_config.429474957.provision_vm_agent":           "false",
				"os_profile_windows_config.429474957.winrm.#":                      "0",
			},
		},
		"v0_1_empty": {
			StateVersion: 0,
			Attributes:   map[string]string{},
			Expected:     map[string]string{},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         "azurerm_virtual_machine",
			Attributes: tc.Attributes,
		}
		is, err := resourceArmVirtualMachine().MigrateState(tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.Expected {
			if is.Attributes[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					tn, k, v, k, is.Attributes[k], is.Attributes)
			}
		}
	}
}
