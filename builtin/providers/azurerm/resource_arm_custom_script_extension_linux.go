package azurerm

import "github.com/hashicorp/terraform/helper/schema"

func resourceArmCustomScriptExtensionLinux() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCustomScriptExtensionLinuxCreateOrUpdate,
		Read:   resourceArmCustomScriptExtensionLinuxRead,
		Update: resourceArmCustomScriptExtensionLinuxCreateOrUpdate,
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
			"command_to_execute": {
				Type:     schema.TypeString,
				Required: true,
			},

			"file_uris": {
				Type: schema.TypeMap,
			},
		},
	}
}

func resourceArmCustomScriptExtensionLinuxCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {

	settings := make(map[string]interface{}, len(2))

	commandToExecute = d.Get("command_to_execute").(string)
	fileUris := d.Get("file_uris").(map[string]interface{})

	settings["command_to_execute"] = commandToExecute
	settings["fileUris"] = fileUrls

	return resourceArmCustomScriptExtensionCreateOrUpdate(d, meta, settings)
}

func resourceArmCustomScriptExtensionLinuxRead(d *schema.ResourceData, meta interface{}) error {

	response, _ := resourceArmCustomScriptExtensionRead(d, meta)

	// TODO: handle the error above

	settings := response.Properties.Settings
	if settings != nil {
		d.Set("file_uris", settings["file_uris"].(*[]string))
		d.Set("command_to_execute", settings["command_to_execute"].(*string))
	}

	return nil
}
