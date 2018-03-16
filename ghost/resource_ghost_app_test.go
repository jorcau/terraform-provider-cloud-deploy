package ghost

import (
	"fmt"
	"log"
	"testing"

	"cloud-deploy.io/go-st"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccGhostAppBasic(t *testing.T) {
	resourceName := "ghost_app.test"
	envName := fmt.Sprintf("ghost_app_acc_env_basic_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGhostAppConfig(envName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGhostAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", envName),
					resource.TestCheckResourceAttr(resourceName, "env", "dev"),
					resource.TestCheckResourceAttr(resourceName, "region", "eu-west-1"),
				),
			},
			{
				Config: testAccGhostAppConfigUpdated(envName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGhostAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", envName),
					resource.TestCheckResourceAttr(resourceName, "env", "dev"),
					resource.TestCheckResourceAttr(resourceName, "region", "eu-west-2"),
				),
			},
		},
	})
}

func testAccCheckGhostAppExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Ghost Application ID is set")
		}

		log.Printf("[INFO] Try to connect to Ghost and get all apps")
		client := testAccProvider.Meta().(*ghost.Client)
		_, err := client.GetApps()
		if err != nil {
			return fmt.Errorf("Ghost environment not reachable: %v", err)
		}

		return nil
	}
}

func testAccGhostAppConfig(name string) string {
	return fmt.Sprintf(`
			resource "ghost_app" "test" {
				name = "%s"
			  env  = "dev"
			  role = "web-front_"

			  region        = "eu-west-1"
			  instance_type = "t2.micro"
			  vpc_id        = "vpc-3f1eb65a"

			  log_notifications = [
			    "ghost-devops@domain.com",
			  ]

			  build_infos = {
			    subnet_id    = "subnet-a7e849fe"
			    ssh_username = "admin"
			    source_ami   = "ami-03ce4474"
			  }

			  environment_infos = {
			    instance_profile  = "iam.ec2.demo"
			    key_name          = "ghost-demo"
			    root_block_device = {
						name = "testblockdevice"
						size = 20
					}
			    optional_volumes  = [{
						device_name = "/dev/xvdd"
						volume_type = "gp2"
						volume_size = 20
					}]
			    subnet_ids        = ["subnet-a7e849fe"]
			    security_groups   = ["sg-6814f60c", "sg-2414f60c"]
					instance_tags			= [{
						tag_name  = "Name"
						tag_value = "wordpress"
					},
					{
						tag_name  = "Type"
						tag_value = "front"
					}]
			  }

			  autoscale = {
			    name = "autoscale"
					min  = 1
					max  = 3
			  }

			  modules = [{
					name       = "wordpress"
			    pre_deploy = ""
			    path       = "/var/www"
			    scope      = "code"
			    git_repo   = "https://github.com/KnpLabs/KnpIpsum.git"
			  },
				{
					name        = "wordpress2"
					pre_deploy  = "ZXhpdCAx"
					post_deploy = "ZXhpdCAx"
					path        = "/var/www-test.test_ew"
					scope       = "code"
					git_repo    = "https://github.com/KnpLabs/KnpIpsum.git"
				}]

			  features = [{
			    version = "5.4"
			    name    = "php5"
			  },
				{
			    version = "2.2"
			    name    = "apache2"
			  }]

				lifecycle_hooks = {
					pre_buildimage  = "#!/usr/bin/env bash"
					post_buildimage = "#!/usr/bin/env bash"
				}

				environment_variables = [{
					key   = "myvar"
					value = "myvalue"
				}]
			}
			`, name)
}

func testAccGhostAppConfigUpdated(name string) string {
	return fmt.Sprintf(`
			resource "ghost_app" "test" {
				name = "%s"
			  env  = "dev"
			  role = "webfront"

			  region        = "eu-west-2"
			  instance_type = "t2.micro"
			  vpc_id        = "vpc-3f1eb65a"

			  log_notifications = [
			    "ghost-devops@domain.com",
			  ]

			  build_infos = {
			    subnet_id    = "subnet-a7e849fe"
			    ssh_username = "admin"
			    source_ami   = "ami-03ce4474"
			  }

			  environment_infos = {
			    instance_profile  = "iam.ec2.demo"
			    key_name          = "ghost-demo"
			    root_block_device = {
						name = "testblockdevice"
						size = 20
					}
					optional_volumes = [{
						device_name = "/dev/xvdd"
						volume_type = "gp2"
						volume_size = 20
					}]
			    subnet_ids        = ["subnet-a7e849fe"]
			    security_groups   = ["sg-6814f60c"]
					instance_tags			= [{
						tag_name  = "Name"
						tag_value = "wordpress"
					},
					{
						tag_name  = "Type"
						tag_value = "front"
					}]
			  }

			  autoscale = {
			    name = "autoscale"
					min  = 1
					max  = 2
			  }

			  modules = [{
					name       = "wordpress"
			    pre_deploy = ""
			    path       = "/var/www"
			    scope      = "code"
			    git_repo   = "https://github.com/KnpLabs/KnpIpsum.git"
			  },
				{
					name        = "wordpress2"
					pre_deploy  = "ZXhpdCAx"
					post_deploy = "ZXhpdCAx"
					path        = "/var/www"
					scope       = "code"
					git_repo    = "https://github.com/KnpLabs/KnpIpsum.git"
				}]

			  features = [{
			    version = "5.4"
			    name    = "php5"
			  },
				{
			    version = "2.2"
			    name    = "apache2"
			  }]

				lifecycle_hooks = {
					pre_buildimage  = "#!/usr/bin/env bash"
				}

				environment_variables = [{
					key   = "myvar2"
					value = "myvalue2"
				}]
			}
			`, name)
}

func testAccGhostAppConfigOmitEmpty(name string) string {
	return fmt.Sprintf(`
			resource "ghost_app" "test" {
				name = "%s"
			  env  = "dev"
			  role = "webfront"

			  region        = "eu-west-2"
			  instance_type = "t2.micro"
			  vpc_id        = "vpc-3f1eb65a"

			  log_notifications = [
			    "ghost-devops@domain.com",
			  ]

			  build_infos = {
			    subnet_id    = "subnet-a7e849fe"
			    ssh_username = "admin"
			    source_ami   = "ami-03ce4474"
			  }

			  environment_infos = {
			    instance_profile  = "iam.ec2.demo"
			    key_name          = "ghost-demo"
			    root_block_device = {
						name = "testblockdevice"
						size = 20
					}
					optional_volumes = [{
						device_name = "data"
						volume_type = "lol"
						volume_size = 20
					}]
			    subnet_ids        = ["subnet-a7e849fe"]
			    security_groups   = ["sg-6814f60c"]
					instance_tags			= [{
						tag_name  = "Name"
						tag_value = "wordpress"
					},
					{
						tag_name  = "Type"
						tag_value = "front"
					}]
			  }

			  autoscale = {
			    name = "autoscale"
					min  = 0
					max  = 2
			  }

			  modules = [{
					name       = "wordpress"
			    pre_deploy = ""
			    path       = "/var/www"
			    scope      = "code"
			    git_repo   = "https://github.com/KnpLabs/KnpIpsum.git"
			  },
				{
					name        = "wordpress2"
					pre_deploy  = "ZXhpdCAx"
					post_deploy = "ZXhpdCAx"
					path        = "/var/www"
					scope       = "code"
					git_repo    = "https://github.com/KnpLabs/KnpIpsum.git"
				}]

			  features = [{
			    version = "5.4"
			    name    = "php5"
			  },
				{
			    version = "2.2"
			    name    = "apache2"
			  }]

				lifecycle_hooks = {
					pre_buildimage  = "#!/usr/bin/env bash"
				}

				environment_variables = [{
					key   = "myvar2"
					value = "myvalue2"
				}]
			}
			`, name)
}
