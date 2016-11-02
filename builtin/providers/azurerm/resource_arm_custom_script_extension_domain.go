package azurerm

import "github.com/hashicorp/terraform/helper/schema"

func resourceArmCustomScriptExtensionDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCustomScriptExtensionDomainCreateOrUpdate,
		Read:   resourceArmCustomScriptExtensionDomainRead,
		Update: resourceArmCustomScriptExtensionDomainCreateOrUpdate,
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
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organisational_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"restart": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"options": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmCustomScriptExtensionDomainCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {

	settings := make(map[string]interface{}, len(1))

	domainName := d.Get("domain_name").(string)
	organisationalPath := d.Get("organisational_path").(string)
	user := d.Get("user").(string)
	restart := d.Get("restart").(bool)
	options := d.Get("options").(string)

	settings["Name"] = domainName
	settings["OUPath"] = organisationalPath
	settings["User"] = user
	settings["Restart"] = restart
	settings["Options"] = options

	return resourceArmCustomScriptExtensionCreateOrUpdate(d, meta, settings)
}

func resourceArmCustomScriptExtensionDomainRead(d *schema.ResourceData, meta interface{}) error {

	response, _ := resourceArmCustomScriptExtensionRead(d, meta)

	// TODO: handle the error above

	settings := response.Properties.Settings
	if settings != nil {
		d.Set("domain_name", settings["Name"].(string))
		d.Set("organisational_path", settings["OUPath"].(string))
		d.Set("user", settings["User"].(string))
		d.Set("restart", settings["Restart"].(string))
		d.Set("options", settings["Options"].(string))
	}

	return nil
}
