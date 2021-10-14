package resources_test

import (
	"fmt"
	"testing"

    "github.com/gthesheep/terraform-provider-dbt-cloud/pkg/dbt_cloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDbtCloudJobResource(t *testing.T) {

	randomID := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	config := fmt.Sprintf(`
			resource "dbt_cloud_job" "test" {
				name = "dbt-cloud-job-%s"
				project_id = 123
				environment_id = 789
				execute_steps = [
				    "dbt run",
				    "dbt test"
				]
				dbt_version = "0.20.0"
				is_active = true
				num_threads = 5
				target_name = "target"
				generate_docs = true
				run_generate_sources = true
				triggers = {
				    "github_webhook": true,
				    "schedule": true,
				    "custom_branch_only": true
				}
			}
		`, randomID)

	check := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("dbt_cloud_job.test", "job_id"),
		resource.TestCheckResourceAttr("dbt_cloud_job.test", "project_id", "123"),
		resource.TestCheckResourceAttr("dbt_cloud_job.test", "environment_id", "789"),
		resource.TestCheckResourceAttr("dbt_cloud_job.test", "name", fmt.Sprintf("dbt-cloud-job-%s", randomID)),
	)

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		CheckDestroy: testAccDbtCloudJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  check,
			},
		},
	})
}

func testAccDbtCloudJobDestroy(s *terraform.State) error {
  providers := providers()
  c := providers["dbt_cloud"].Meta().(*dbt_cloud.Client)

  for _, resource := range s.RootModule().Resources {
    if resource.Type != "dbt_cloud_job" {
      continue
    }

    resourceID := resource.Primary.ID
    response, err := c.GetJob(string(resourceID))
    if err == nil {
      if (response != nil) && (fmt.Sprint(*response.ID) == resourceID) {
        return fmt.Errorf("Job (%s) still exists.", resourceID)
      }

      return nil
    }

  }

  return nil
}