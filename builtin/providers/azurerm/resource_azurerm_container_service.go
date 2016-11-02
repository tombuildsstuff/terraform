package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/containerservice"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmContainerService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerServiceCreate,
		Read:   resourceArmContainerServiceRead,
		Update: resourceArmContainerServiceCreate,
		Delete: resourceArmContainerServiceDelete,

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
		},
	}
}

func resourceArmContainerServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	containerServiceClient := client.containerServiceClient

	log.Printf("[INFO] preparing arguments for Azure ARM EventHub Namespace creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)

	platform := d.Get("platform").(string)

	masterCount := int32(d.Get("master_count").(int))
	masterDnsPrefix := d.Get("master_dns_prefix").(string)
	masterFqdn := d.Get("master_fqdn").(string)

	linuxAdminUsername := d.Get("linux_admin_username").(string)
	linuxAdminSshKeys := d.Get("linux_admin_ssh_keys").(map[string]interface{})

	windowsAdminUsername := d.Get("windows_admin_username").(string)
	windowsAdminPassword := d.Get("windows_admin_password").(string)

	diagnosticsEnabled := d.Get("diagnostics_enabled").(bool)
	diagnosticsStorageUri := d.Get("storage_uri").(bool)

	tags := d.Get("tags").(map[string]interface{})

	parameters := containerservice.ContainerService{
		Name:     name,
		Location: &location,
		Properties: containerservice.Properties{
			MasterProfile: containerservice.MasterProfile{
				Count:     masterCount,
				DNSPrefix: masterDnsPrefix,
				Fqdn:      masterFqdn,
			},
			LinuxProfile: containerservice.LinuxProfile{
				AdminUsername: linuxAdminUsername,
				SSH: containerservice.SSHConfiguration{
					PublicKeys: []containerservice.SSHPublicKey{
						{
							KeyData: "ssh_key",
						},
					},
				},
			},
			WindowsProfile: containerservice.WindowsProfile{
				AdminUsername: windowsAdminUsername,
				AdminPassword: windowsAdminPassword,
			},
			DiagnosticsProfile: containerservice.DiagnosticsProfile{
				VMDiagnostics: containerservice.VMDiagnostics{
					Enabled:    diagnosticsEnabled,
					StorageURI: diagnosticsStorageUri,
				},
			},
			OrchestratorProfile: containerservice.OrchestratorProfile{
				OrchestratorType: containerservice.OchestratorTypes(platform),
			},
			/*
				AgentPoolProfiles: []containerservice.AgentPoolProfile{
					{
						Count: agentPoolCount,
						DNSPrefix: dnsPrefix,
						Fqdn: fqdn,
						Name: name,
						VMSize: vmSize,
					},
				},
			*/
		},
		Tags: expandTags(tags),
	}

	_, err := containerServiceClient.CreateOrUpdate(resGroup, name, parameters, make(chan struct{}))
	if err != nil {
		return err
	}

	read, err := containerServiceClient.Get(resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read EventHub Namespace %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmContainerServiceRead(d, meta)
}

func resourceArmContainerServiceRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceArmContainerServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	containerServiceClient := client.containerServiceClient

	resGroup := "some-rg"
	name := "some-name"

	_, err := containerServiceClient.Delete(resGroup, name)

	return nil
}
