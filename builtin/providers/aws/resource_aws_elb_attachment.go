package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsElbAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsElbAttachmentCreate,
		Read:   resourceAwsElbAttachmentRead,
		Update: resourceAwsElbAttachmentUpdate,
		Delete: resourceAwsElbAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"elb": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"instances": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Required: true,
			},
		},
	}
}

func resourceAwsElbAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	elbconn := meta.(*AWSClient).elbconn
	elbName := d.Get("elb").(string)

	instances := expandInstanceString(d.Get("instances").(*schema.Set).List())

	registerInstancesOpts := elb.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(elbName),
		Instances:        instances,
	}

	_, err := elbconn.RegisterInstancesWithLoadBalancer(&registerInstancesOpts)
	if err != nil {
		return fmt.Errorf("Failure registering instances with ELB: %s", err)
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", elbName)))

	return nil
}

func resourceAwsElbAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	elbconn := meta.(*AWSClient).elbconn
	elbName := d.Get("elb").(string)

	// Retrieve the ELB properties to get a list of attachments
	describeElbOpts := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{aws.String(elbName)},
	}

	describeResp, err := elbconn.DescribeLoadBalancers(describeElbOpts)
	if err != nil {
		return fmt.Errorf("Error retrieving ELB: %s", err)
	}
	if len(describeResp.LoadBalancerDescriptions) != 1 {
		return fmt.Errorf("Unable to find ELB: %#v", describeResp.LoadBalancerDescriptions)
	}

	lb := describeResp.LoadBalancerDescriptions[0]

	d.Set("instances", flattenInstances(lb.Instances))

	return nil
}

func resourceAwsElbAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	elbconn := meta.(*AWSClient).elbconn
	elbName := d.Get("elb").(string)

	if d.HasChange("instances") {
		o, n := d.GetChange("instances")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := expandInstanceString(os.Difference(ns).List())
		add := expandInstanceString(ns.Difference(os).List())

		if len(add) > 0 {
			registerInstancesOpts := elb.RegisterInstancesWithLoadBalancerInput{
				LoadBalancerName: aws.String(elbName),
				Instances:        add,
			}

			_, err := elbconn.RegisterInstancesWithLoadBalancer(&registerInstancesOpts)
			if err != nil {
				return fmt.Errorf("Failure registering instances with ELB: %s", err)
			}
		}
		if len(remove) > 0 {
			deRegisterInstancesOpts := elb.DeregisterInstancesFromLoadBalancerInput{
				LoadBalancerName: aws.String(elbName),
				Instances:        remove,
			}

			_, err := elbconn.DeregisterInstancesFromLoadBalancer(&deRegisterInstancesOpts)
			if err != nil {
				return fmt.Errorf("Failure deregistering instances from ELB: %s", err)
			}
		}
	}

	return resourceAwsElbAttachmentRead(d, meta)
}

func resourceAwsElbAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	elbconn := meta.(*AWSClient).elbconn
	elbName := d.Get("elb").(string)

	instanceSet := d.Get("instances").(*schema.Set)
	instances := expandInstanceString(instanceSet.List())

	log.Printf("[INFO] Deleting Attachments %s from: %s", instanceSet, elbName)

	deRegisterInstancesOpts := elb.DeregisterInstancesFromLoadBalancerInput{
		LoadBalancerName: aws.String(elbName),
		Instances:        instances,
	}

	_, err := elbconn.DeregisterInstancesFromLoadBalancer(&deRegisterInstancesOpts)
	if err != nil {
		return fmt.Errorf("Failure deregistering instances from ELB: %s", err)
	}

	return nil
}
