package ghost

import (
	"fmt"
	"log"
	"reflect"
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
			  role = "webfront"

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
					path        = "/var/www-test.test"
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

var (
	app = ghost.App{
		Name:               "app_name",
		Env:                "test",
		Role:               "web",
		Region:             "us-west-1",
		InstanceType:       "t2.micro",
		VpcID:              "vpc-123456",
		InstanceMonitoring: false,

		Modules: &[]ghost.Module{{
			Name:      "my_module",
			GitRepo:   "https://github.com/test/test.git",
			Scope:     "system",
			Path:      "/",
			BuildPack: StrToB64("#!/usr/bin/env bash"),
			PreDeploy: StrToB64("#!/usr/bin/env bash"),
		}},
		Features: &[]ghost.Feature{{
			Name:        "feature",
			Version:     "1.0",
			Provisioner: "ansible",
		}},
		Autoscale: &ghost.Autoscale{
			Name:          "autoscale",
			EnableMetrics: false,
			Min:           0,
			Max:           3,
		},
		BuildInfos: &ghost.BuildInfos{
			SshUsername: "admin",
			SourceAmi:   "ami-1",
			SubnetID:    "subnet-1",
		},
		EnvironmentInfos: &ghost.EnvironmentInfos{
			InstanceProfile: "profile",
			KeyName:         "key",
			PublicIpAddress: false,
			SecurityGroups:  []string{"sg-1", "sg-2"},
			SubnetIDs:       []string{"subnet-1", "subnet-2"},
			InstanceTags: &[]ghost.InstanceTag{{
				TagName:  "name",
				TagValue: "val",
			}},
			OptionalVolumes: &[]ghost.OptionalVolume{{
				DeviceName: "my_device",
				VolumeType: "gp2",
				VolumeSize: 20,
				Iops:       3000,
			}},
			RootBlockDevice: &ghost.RootBlockDevice{
				Name: "rootblock",
				Size: 20,
			},
		},
		LifecycleHooks: &ghost.LifecycleHooks{
			PreBuildimage:  StrToB64("#!/usr/bin/env bash"),
			PostBuildimage: StrToB64("#!/usr/bin/env bash"),
		},
		LogNotifications: []string{"log_not@email.com"},
		EnvironmentVariables: &[]ghost.EnvironmentVariable{{
			Key:   "env_var_key",
			Value: "env_var_value",
		}},
	}
)

// Expanders Unit Tests
func TestExpandGhostAppStringList(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput []string
	}{
		{
			[]interface{}{
				"1", "2", "3",
			},
			[]string{
				"1", "2", "3",
			},
		},
		{
			nil,
			nil,
		},
	}

	for _, tc := range cases {
		output := expandGhostAppStringList(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandGhostAppInstanceTags(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *[]ghost.InstanceTag
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"tag_name":  "name",
					"tag_value": "val",
				},
			},
			app.EnvironmentInfos.InstanceTags,
		},
		{
			nil,
			&[]ghost.InstanceTag{},
		},
	}

	for _, tc := range cases {
		output := expandGhostAppInstanceTags(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandGhostAppOptionalVolume(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *[]ghost.OptionalVolume
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"device_name": "my_device",
					"volume_type": "gp2",
					"volume_size": 20,
					"iops":        3000,
					"launch_block_device_mappings": false,
				},
			},
			app.EnvironmentInfos.OptionalVolumes,
		},
		{
			nil,
			&[]ghost.OptionalVolume{},
		},
	}

	for _, tc := range cases {
		output := expandGhostAppOptionalVolumes(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandGhostAppRootBlockDevice(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *ghost.RootBlockDevice
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"name": "rootblock",
					"size": 20,
				},
			},
			app.EnvironmentInfos.RootBlockDevice,
		},
		{
			nil,
			nil,
		},
	}

	for _, tc := range cases {
		output := expandGhostAppRootBlockDevice(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandGhostAppEnvironmentInfos(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *ghost.EnvironmentInfos
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"instance_profile":  "profile",
					"key_name":          "key",
					"public_ip_address": false,
					"security_groups":   []interface{}{"sg-1", "sg-2"},
					"subnet_ids":        []interface{}{"subnet-1", "subnet-2"},
					"instance_tags": []interface{}{
						map[string]interface{}{
							"tag_name":  "name",
							"tag_value": "val",
						},
					},
					"optional_volumes": []interface{}{
						map[string]interface{}{
							"device_name": "my_device",
							"volume_type": "gp2",
							"volume_size": 20,
							"iops":        3000,
							"launch_block_device_mappings": false,
						},
					},
					"root_block_device": []interface{}{
						map[string]interface{}{
							"name": "rootblock",
							"size": 20,
						},
					},
				},
			},
			app.EnvironmentInfos,
		},
	}

	for _, tc := range cases {
		output := expandGhostAppEnvironmentInfos(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandGhostAppBuildInfos(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *ghost.BuildInfos
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"ssh_username": "admin",
					"source_ami":   "ami-1",
					"subnet_id":    "subnet-1",
					"ami_name":     "",
				},
			},
			app.BuildInfos,
		},
	}

	for _, tc := range cases {
		output := expandGhostAppBuildInfos(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandGhostAppFeatures(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *[]ghost.Feature
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"name":        "feature",
					"version":     "1.0",
					"provisioner": "ansible",
				},
			},
			app.Features,
		},
		{
			nil,
			&[]ghost.Feature{},
		},
	}

	for _, tc := range cases {
		output := expandGhostAppFeatures(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandGhostAppLifecycleHooks(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *ghost.LifecycleHooks
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"pre_buildimage":  "#!/usr/bin/env bash",
					"post_buildimage": "#!/usr/bin/env bash",
					"pre_bootstrap":   "",
					"post_bootstrap":  "",
				},
			},
			app.LifecycleHooks,
		},
		{
			nil,
			nil,
		},
	}

	for _, tc := range cases {
		output := expandGhostAppLifecycleHooks(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandGhostAppAutoscale(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *ghost.Autoscale
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"name":           "autoscale",
					"enable_metrics": false,
					"min":            0,
					"max":            3,
				},
			},
			app.Autoscale,
		},
		{
			nil,
			nil,
		},
	}

	for _, tc := range cases {
		output := expandGhostAppAutoscale(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandGhostAppEnvironmentVariables(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *[]ghost.EnvironmentVariable
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"key":   "env_var_key",
					"value": "env_var_value",
				},
			},
			app.EnvironmentVariables,
		},
		{
			nil,
			&[]ghost.EnvironmentVariable{},
		},
	}

	for _, tc := range cases {
		output := expandGhostAppEnvironmentVariables(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandGhostAppModules(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *[]ghost.Module
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"name":             "my_module",
					"git_repo":         "https://github.com/test/test.git",
					"path":             "/",
					"scope":            "system",
					"build_pack":       "#!/usr/bin/env bash",
					"pre_deploy":       "#!/usr/bin/env bash",
					"post_deploy":      "",
					"after_all_deploy": "",
					"uid":              0,
					"gid":              0,
					"last_deployment":  "",
				},
			},
			app.Modules,
		},
	}

	for _, tc := range cases {
		output := expandGhostAppModules(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
