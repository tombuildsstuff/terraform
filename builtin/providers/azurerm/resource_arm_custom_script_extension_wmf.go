package azurerm

import "github.com/hashicorp/terraform/helper/schema"

func resourceArmCustomScriptExtensionWmf() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCustomScriptExtensionWmfCreateOrUpdate,
		Read:   resourceArmCustomScriptExtensionWmfRead,
		Update: resourceArmCustomScriptExtensionWmfCreateOrUpdate,
		Delete: resourceArmCustomScriptExtensionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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

			"protected_settings": {
				Type: schema.TypeMap,
			},

			"tags": tagsSchema(),

			// Above here ^ are the 'Core' schema
			"wmf_version": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmCustomScriptExtensionWmfCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {

	settings := make(map[string]interface{}, len(1))

	wmf_version := d.Get("wmf_version").(string)
	settings["WmfVersion"] = wmf_version

	return resourceArmCustomScriptExtensionCreateOrUpdate(d, meta, settings)
}

func resourceArmCustomScriptExtensionWmfRead(d *schema.ResourceData, meta interface{}) error {

	response, _ := resourceArmCustomScriptExtensionRead(d, meta)

	// TODO: handle the error above

	settings := response.Properties.Settings
	if settings != nil {
		d.Set("wmf_version", settings["WmfVersion"].(string))
	}

	return nil
}
