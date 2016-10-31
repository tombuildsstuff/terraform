package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tombuildsstuff/azure-sdk-for-go/arm/eventhub"
	"net/http"
)

func resourceArmEventHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEventHubCreate,
		Read:   resourceArmEventHubRead,
		Update: resourceArmEventHubCreate,
		Delete: resourceArmEventHubDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"partition_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateEventHubPartitionCount,
			},

			"message_retention": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateEventHubMessageRetentionCount,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmEventHubCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	eventhubClient := client.eventHubClient
	log.Printf("[INFO] preparing arguments for Azure ARM EventHub creation.")

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	partitionCount := int64(d.Get("partition_count").(int))
	messageRetention := int64(d.Get("message_retention").(int))

	parameters := eventhub.CreateOrUpdateParameters{
		Location: &location,
		Properties: &eventhub.Properties{
			PartitionCount:         &partitionCount,
			MessageRetentionInDays: &messageRetention,
		},
	}

	_, err := eventhubClient.CreateOrUpdate(resGroup, namespaceName, name, parameters)
	if err != nil {
		return err
	}

	read, err := eventhubClient.Get(resGroup, namespaceName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read EventHub %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmEventHubRead(d, meta)
}

func resourceArmEventHubRead(d *schema.ResourceData, meta interface{}) error {
	eventhubClient := meta.(*ArmClient).eventHubClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	namespaceName := ""
	name := id.Path["namespaces"]

	resp, err := eventhubClient.Get(resGroup, namespaceName, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure EventHub %s: %s", name, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	d.Set("name", resp.Name)
	// TODO: complete me

	flattenAndSetTags(d, resp.Tags)
	return nil
}

func resourceArmEventHubDelete(d *schema.ResourceData, meta interface{}) error {
	/*
		eventhubClient := meta.(*ArmClient).eventHubClient

		id, err := parseAzureResourceID(d.Id())
		if err != nil {
			return err
		}
		resGroup := id.ResourceGroup
		name := id.Path["namespaces"]

		resp, err := eventhubClient.Delete(resGroup, name, make(chan struct{}))

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Error issuing Azure ARM delete request of EventHub'%s': %s", name, err)
		}
	*/
	return nil
}

func validateEventHubPartitionCount(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)

	if !(32 >= value && value >= 2) {
		errors = append(errors, fmt.Errorf("EventHub Partition Count has to be between 2 and 32"))
	}
	return
}

func validateEventHubMessageRetentionCount(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)

	if !(7 >= value && value >= 1) {
		errors = append(errors, fmt.Errorf("EventHub Retention Count has to be between 1 and 7"))
	}
	return
}
