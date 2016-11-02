package azurerm

import (
	"fmt"
	"log"

	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmCustomScriptExtension() *schema.Resource {
	return &schema.Resource{

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: azureRMNormalizeLocation,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"virtual_machine_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"publisher": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"version": {
				Type:     schema.TypeString,
				Required: true,
			},

			"auto_upgrade_minor_version": {
				Type:    schema.TypeBool,
				Default: true,
			},

			"settings": {
				Type: schema.TypeMap,
			},

			"protected_settings": {
				Type: schema.TypeMap,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmCustomScriptExtensionCreateOrUpdate(d *schema.ResourceData, meta interface{}, settings map[string]interface{}) error {
	vmExtensionsClient := meta.(*ArmClient).vmExtensionClient
	log.Printf("[INFO] preparing arguments for Azure ARM Custom Script Extension creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	vmName := d.Get("virtual_machine_name").(string)

	location := d.Get("location").(string)

	extensionType := d.Get("type").(*string)
	extensionVersion := d.Get("version").(*string)
	extensionPublisher := d.Get("publisher").(*string)

	autoUpgrade := d.Get("auto_upgrade_minor_version").(bool)

	protectedSettings := d.Get("protected_settings").(map[string]interface{})

	tags := d.Get("tags").(map[string]interface{})

	parameters := compute.VirtualMachineExtension{
		Name:     &name,
		Location: &location,
		Properties: &compute.VirtualMachineExtensionProperties{
			AutoUpgradeMinorVersion: &autoUpgrade,
			Publisher:               extensionPublisher,
			Type:                    extensionType,
			TypeHandlerVersion:      extensionVersion,

			Settings:          &settings,
			ProtectedSettings: &protectedSettings,
		},
		Tags: expandTags(tags),
	}

	_, err := vmExtensionsClient.CreateOrUpdate(resGroup, vmName, name, parameters, make(chan struct{}))
	if err != nil {
		return err
	}

	read, err := vmExtensionsClient.Get(resGroup, vmName, name, "")
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Custom Script Extension %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmCustomScriptExtensionRead(d, meta)
}

func resourceArmCustomScriptExtensionRead(d *schema.ResourceData, meta interface{}) (result compute.VirtualMachineExtension, err error) {
	vmExtensionsClient := meta.(*ArmClient).vmExtensionClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	vmName := id.Path["virtualMachines"]
	name := id.Path["extensions"]

	resp, err := vmExtensionsClient.Get(resGroup, vmName, name, "")
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure Custom Script Extension %s: %s", name, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("virtual_machine_name", vmName)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	d.Set("type", resp.Properties.Type)
	d.Set("version", resp.Properties.TypeHandlerVersion)
	d.Set("publisher", resp.Properties.Publisher)
	d.Set("auto_upgrade_minor_version", resp.Properties.AutoUpgradeMinorVersion)

	d.Set("settings", resp.Properties.Settings)
	d.Set("protected_settings", resp.Properties.ProtectedSettings)

	flattenAndSetTags(d, resp.Tags)
	result = resp

	return nil
}

func resourceArmCustomScriptExtensionDelete(d *schema.ResourceData, meta interface{}) error {
	vmExtensionsClient := meta.(*ArmClient).vmExtensionClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	vmName := id.Path["virtualMachines"]
	name := id.Path["extensions"]

	resp, err := vmExtensionsClient.Delete(resGroup, vmName, name, make(chan struct{}))

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("Error issuing Azure ARM delete request of Custom Script Extension '%s': %s", name, err)
	}

	return nil
}
