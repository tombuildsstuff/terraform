package opsgenie

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/opsgenie/opsgenie-go-sdk/contact"
)

func resourceOpsGenieContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieContactCreate,
		Read:   resourceOpsGenieContactRead,
		Update: resourceOpsGenieContactUpdate,
		Delete: resourceOpsGenieContactDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"method": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"email",
					"sms",
					"voice",
				}, true),
			},
			"to": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
				Default:  true,
			},
		},
	}
}

func resourceOpsGenieContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).contacts

	username := d.Get("username").(string)
	method := d.Get("method").(string)
	to := d.Get("to").(string)

	createRequest := contact.CreateContactRequest{
		Username: username,
		Method:   method,
		To:       to,
	}

	log.Printf("[INFO] Creating OpsGenie Contact '%s' for User '%s'", to, username)

	createResponse, err := client.Create(createRequest)
	if err != nil {
		return err
	}

	err = checkOpsGenieResponse(createResponse.Code, createResponse.Status)
	if err != nil {
		return err
	}

	d.SetId(createResponse.Id)

	return resourceOpsGenieContactUpdate(d, meta)
}

func resourceOpsGenieContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).contacts

	username := d.Get("username").(string)
	method := d.Get("method").(string)
	to := d.Get("to").(string)
	enabled := d.Get("enabled").(bool)

	updateRequest := contact.UpdateContactRequest{
		Id:       d.Id(),
		Username: username,
		To:       to,
	}

	log.Printf("[INFO] Updating OpsGenie contact '%s'", d.Id())

	updateResponse, err := client.Update(updateRequest)
	if err != nil {
		return err
	}

	err = checkOpsGenieResponse(updateResponse.Code, updateResponse.Status)
	if err != nil {
		return err
	}

	getRequest := contact.GetContactRequest{
		Id: d.Id(),
	}

	getResponse, err := client.Get(getRequest)
	if err != nil {
		return err
	}

	if !enabled && getResponse.Enabled {
		disableRequest := contact.DisableContactRequest{
			Id:       d.Id(),
			Username: username,
		}
		_, err := client.Disable(disableRequest)
		err = checkOpsGenieResponse(updateResponse.Code, updateResponse.Status)
		if err != nil {
			return err
		}
	}

	if enabled && !getResponse.Enabled {
		enableRequest := contact.EnableContactRequest{
			Id:       d.Id(),
			Username: username,
		}
		_, err := client.Enable(enableRequest)
		err = checkOpsGenieResponse(updateResponse.Code, updateResponse.Status)
		if err != nil {
			return err
		}
	}

	return resourceOpsGenieContactRead(d, meta)
}

func resourceOpsGenieContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).contacts

	getRequest := contact.GetContactRequest{
		Id: d.Id(),
	}

	getResponse, err := client.Get(getRequest)
	if err != nil {
		return err
	}

	d.Set("username", getResponse.Username)
	d.Set("method", getResponse.Method)
	d.Set("to", getResponse.To)
	d.Set("enabled", getResponse.Enabled)

	return nil
}

func resourceOpsGenieContactDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie Contact '%s' for user '%s'", d.Id(), d.Get("username").(string))
	client := meta.(*OpsGenieClient).contacts

	deleteRequest := contact.DeleteContactRequest{
		Id: d.Id(),
	}

	_, err := client.Delete(deleteRequest)
	if err != nil {
		return err
	}

	return nil
}
