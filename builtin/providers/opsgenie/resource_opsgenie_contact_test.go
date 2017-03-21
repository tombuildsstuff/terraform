package opsgenie

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOpsGenieContact_email(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccOpsGenieContact_email, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieContactExists("opsgenie_contact.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieContact_sms(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccOpsGenieContact_sms, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieContactExists("opsgenie_contact.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieContact_voice(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccOpsGenieContact_voice, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieContactExists("opsgenie_contact.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieContactDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OpsGenieClient).contacts

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_contact" {
			continue
		}

		req := contact.GetContactRequest{
			Id: rs.Primary.Attributes["id"],
		}

		result, _ := client.Get(req)
		if result != nil {
			return fmt.Errorf("Contact still exists:\n%#v", result)
		}
	}

	return nil
}

func testCheckOpsGenieContactExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.Attributes["id"]
		username := rs.Primary.Attributes["username"]

		client := testAccProvider.Meta().(*OpsGenieClient).contacts

		req := contact.GetContactRequest{
			Id: rs.Primary.Attributes["id"],
		}

		result, _ := client.Get(req)
		if result == nil {
			return fmt.Errorf("Bad: Contact %q (username: %q) does not exist", id, username)
		}

		return nil
	}
}

var testAccOpsGenieContact_email = `
resource "opsgenie_user" "test" {
  username  = "acctest-%d@example.tld"
  full_name = "Acceptance Test User"
  role      = "User"
}

resource "opsgenie_contact" "test" {
  username = "${opsgenie_user.test.username}"
  method   = "email"
  to       = "acctest-%d@example.tld"
}
`

var testAccOpsGenieContact_sms = `
resource "opsgenie_user" "test" {
  username  = "acctest-%d@example.tld"
  full_name = "Acceptance Test User"
  role      = "User"
}

resource "opsgenie_contact" "test" {
  username = "${opsgenie_user.test.username}"
  method   = "sms"
  to       = "+447123456789"
}
`
var testAccOpsGenieContact_voice = `
resource "opsgenie_user" "test" {
  username  = "acctest-%d@example.tld"
  full_name = "Acceptance Test User"
  role      = "User"
}

resource "opsgenie_contact" "test" {
  username = "${opsgenie_user.test.username}"
  method   = "voice"
  to       = "+447987654321"
}
`
