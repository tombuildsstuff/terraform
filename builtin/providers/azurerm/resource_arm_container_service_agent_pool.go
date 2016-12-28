package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/containerservice"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmContainerServiceAgentPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerServiceAgentPoolCreateUpdate,
		Read:   resourceArmContainerServiceAgentPoolRead,
		Update: resourceArmContainerServiceAgentPoolCreateUpdate,
		Delete: resourceArmContainerServiceAgentPoolDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"container_service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateArmContainerServiceAgentPoolProfileCount,
			},

			"dns_prefix": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vm_size": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmContainerServiceAgentPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	containerServiceClient := client.containerServicesClient

	containerServiceName := d.Get("container_service_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	count := int32(d.Get("count").(int))
	dnsPrefix := d.Get("dns_prefix").(string)
	fqdn := d.Get("fqdn").(string)
	vmSize := d.Get("vm_size").(string)

	armMutexKV.Lock(containerServiceName)
	defer armMutexKV.Unlock(containerServiceName)

	log.Printf("[INFO] preparing arguments for Azure ARM Container Service Agent Pool creation.")

	existingContainerService, err := containerServiceClient.Get(resourceGroup, containerServiceName)
	if err != nil {
		return err
	}

	profile := containerservice.AgentPoolProfile{
		Name:      &name,
		Count:     &count,
		DNSPrefix: &dnsPrefix,
		Fqdn:      &fqdn,
		VMSize:    containerservice.VMSizeTypes(vmSize),
	}

	agentPoolProfiles := append(*existingContainerService.AgentPoolProfiles, profile)

	newContainerService := containerservice.ContainerService{
		ID:       existingContainerService.ID,
		Location: existingContainerService.Location,
		Properties: &containerservice.Properties{
			AgentPoolProfiles: &agentPoolProfiles,
			MasterProfile:     existingContainerService.MasterProfile,
			LinuxProfile:      existingContainerService.LinuxProfile,
		},
	}

	newContainerService.MasterProfile.Fqdn = nil

	_, err = containerServiceClient.CreateOrUpdate(resourceGroup, containerServiceName, newContainerService, make(chan (struct{})))
	if err != nil {
		return err
	}

	read, err := containerServiceClient.Get(resourceGroup, containerServiceName)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Container Service %s (resource group %s) ID", containerServiceName, resourceGroup)
	}

	log.Printf("[DEBUG] Waiting for Container Service (%s) to become available", containerServiceName)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    containerServiceStateRefreshFunc(client, resourceGroup, containerServiceName),
		Timeout:    30 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Container Service (%s) to become available: %s", containerServiceName, err)
	}

	d.SetId(resourceArmContainerServiceAgentPoolId(*read.ID, name))

	return resourceArmContainerServiceAgentPoolRead(d, meta)
}

func resourceArmContainerServiceAgentPoolRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmContainerServiceAgentPoolDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmContainerServiceAgentPoolId(containerServiceName string, name string) string {
	// agent pool's don't have ID's - but name's.. given they could be re-used - we need to come up with an ID..
	return fmt.Sprint("%s/%s", containerServiceName, name)
}

func validateArmContainerServiceAgentPoolProfileCount(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value > 100 || 0 >= value {
		errors = append(errors, fmt.Errorf("The Count for an Agent Pool Profile can only be between 1 and 100."))
	}
	return
}
