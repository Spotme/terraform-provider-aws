package aws

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAWSAutoscalingPolicy_basic(t *testing.T) {
	var policy autoscaling.ScalingPolicy

	resourceSimpleName := "aws_autoscaling_policy.foobar_simple"
	resourceStepName := "aws_autoscaling_policy.foobar_step"
	resourceTargetTrackingName := "aws_autoscaling_policy.foobar_target_tracking"

	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, autoscaling.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAutoscalingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAutoscalingPolicyConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingPolicyExists(resourceSimpleName, &policy),
					resource.TestCheckResourceAttr(resourceSimpleName, "adjustment_type", "ChangeInCapacity"),
					resource.TestCheckResourceAttr(resourceSimpleName, "policy_type", "SimpleScaling"),
					resource.TestCheckResourceAttr(resourceSimpleName, "cooldown", "300"),
					resource.TestCheckResourceAttr(resourceSimpleName, "name", name+"-foobar_simple"),
					resource.TestCheckResourceAttr(resourceSimpleName, "scaling_adjustment", "2"),
					resource.TestCheckResourceAttr(resourceSimpleName, "autoscaling_group_name", name),

					testAccCheckScalingPolicyExists(resourceStepName, &policy),
					resource.TestCheckResourceAttr(resourceStepName, "adjustment_type", "ChangeInCapacity"),
					resource.TestCheckResourceAttr(resourceStepName, "policy_type", "StepScaling"),
					resource.TestCheckResourceAttr(resourceStepName, "name", name+"-foobar_step"),
					resource.TestCheckResourceAttr(resourceStepName, "metric_aggregation_type", "Minimum"),
					resource.TestCheckResourceAttr(resourceStepName, "estimated_instance_warmup", "200"),
					resource.TestCheckResourceAttr(resourceStepName, "autoscaling_group_name", name),
					resource.TestCheckTypeSetElemNestedAttrs(resourceStepName, "step_adjustment.*", map[string]string{
						"scaling_adjustment": "1",
					}),
					testAccCheckScalingPolicyExists(resourceTargetTrackingName, &policy),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "policy_type", "TargetTrackingScaling"),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "name", name+"-foobar_target_tracking"),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "autoscaling_group_name", name),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "target_tracking_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "target_tracking_configuration.0.customized_metric_specification.#", "0"),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "target_tracking_configuration.0.predefined_metric_specification.#", "1"),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "target_tracking_configuration.0.predefined_metric_specification.0.predefined_metric_type", "ASGAverageCPUUtilization"),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "target_tracking_configuration.0.target_value", "40"),
				),
			},
			{
				ResourceName:      resourceSimpleName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSAutoscalingPolicyImportStateIdFunc(resourceSimpleName),
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceStepName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSAutoscalingPolicyImportStateIdFunc(resourceStepName),
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceTargetTrackingName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSAutoscalingPolicyImportStateIdFunc(resourceTargetTrackingName),
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSAutoscalingPolicyConfig_basicUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingPolicyExists(resourceSimpleName, &policy),
					resource.TestCheckResourceAttr(resourceSimpleName, "policy_type", "SimpleScaling"),
					resource.TestCheckResourceAttr(resourceSimpleName, "cooldown", "30"),
					testAccCheckScalingPolicyExists(resourceStepName, &policy),
					resource.TestCheckResourceAttr(resourceStepName, "policy_type", "StepScaling"),
					resource.TestCheckResourceAttr(resourceStepName, "estimated_instance_warmup", "20"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceStepName, "step_adjustment.*", map[string]string{
						"scaling_adjustment": "10",
					}),
					testAccCheckScalingPolicyExists(resourceTargetTrackingName, &policy),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "policy_type", "TargetTrackingScaling"),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "target_tracking_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "target_tracking_configuration.0.customized_metric_specification.#", "1"),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "target_tracking_configuration.0.customized_metric_specification.0.statistic", "Average"),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "target_tracking_configuration.0.predefined_metric_specification.#", "0"),
					resource.TestCheckResourceAttr(resourceTargetTrackingName, "target_tracking_configuration.0.target_value", "70"),
				),
			},
		},
	})
}

func TestAccAWSAutoscalingPolicy_disappears(t *testing.T) {
	var policy autoscaling.ScalingPolicy

	resourceName := "aws_autoscaling_policy.foobar_simple"

	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, autoscaling.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAutoscalingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAutoscalingPolicyConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingPolicyExists(resourceName, &policy),
					testAccCheckScalingPolicyDisappears(&policy),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckScalingPolicyDisappears(conf *autoscaling.ScalingPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).autoscalingconn

		params := &autoscaling.DeletePolicyInput{
			AutoScalingGroupName: conf.AutoScalingGroupName,
			PolicyName:           conf.PolicyName,
		}

		_, err := conn.DeletePolicy(params)
		if err != nil {
			return err
		}

		return resource.Retry(10*time.Minute, func() *resource.RetryError {
			params := &autoscaling.DescribePoliciesInput{
				AutoScalingGroupName: conf.AutoScalingGroupName,
				PolicyNames:          []*string{conf.PolicyName},
			}
			resp, err := conn.DescribePolicies(params)
			if err != nil {
				cgw, ok := err.(awserr.Error)
				if ok && cgw.Code() == "ValidationError" {
					return nil
				}
				return resource.NonRetryableError(
					fmt.Errorf("Error retrieving Autoscaling Policy: %s", err))
			}
			if resp.ScalingPolicies == nil || len(resp.ScalingPolicies) == 0 {
				return nil
			}
			return resource.RetryableError(fmt.Errorf(
				"Waiting for Autoscaling Policy: %v", conf.PolicyName))
		})
	}
}

func TestAccAWSAutoscalingPolicy_SimpleScalingStepAdjustment(t *testing.T) {
	var policy autoscaling.ScalingPolicy

	resourceName := "aws_autoscaling_policy.foobar_simple"

	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, autoscaling.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAutoscalingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAutoscalingPolicyConfig_SimpleScalingStepAdjustment(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingPolicyExists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "adjustment_type", "ExactCapacity"),
					resource.TestCheckResourceAttr(resourceName, "scaling_adjustment", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSAutoscalingPolicyImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSAutoscalingPolicy_TargetTrack_Predefined(t *testing.T) {
	var policy autoscaling.ScalingPolicy

	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, autoscaling.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAutoscalingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsAutoscalingPolicyConfig_TargetTracking_Predefined(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingPolicyExists("aws_autoscaling_policy.test", &policy),
				),
			},
			{
				ResourceName:      "aws_autoscaling_policy.test",
				ImportState:       true,
				ImportStateIdFunc: testAccAWSAutoscalingPolicyImportStateIdFunc("aws_autoscaling_policy.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSAutoscalingPolicy_TargetTrack_Custom(t *testing.T) {
	var policy autoscaling.ScalingPolicy

	name := fmt.Sprintf("terraform-testacc-asp-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, autoscaling.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAutoscalingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsAutoscalingPolicyConfig_TargetTracking_Custom(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingPolicyExists("aws_autoscaling_policy.test", &policy),
				),
			},
			{
				ResourceName:      "aws_autoscaling_policy.test",
				ImportState:       true,
				ImportStateIdFunc: testAccAWSAutoscalingPolicyImportStateIdFunc("aws_autoscaling_policy.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSAutoscalingPolicy_zerovalue(t *testing.T) {
	var simplepolicy autoscaling.ScalingPolicy
	var steppolicy autoscaling.ScalingPolicy

	resourceSimpleName := "aws_autoscaling_policy.foobar_simple"
	resourceStepName := "aws_autoscaling_policy.foobar_step"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, autoscaling.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAutoscalingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAutoscalingPolicyConfig_zerovalue(acctest.RandString(5)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingPolicyExists(resourceSimpleName, &simplepolicy),
					testAccCheckScalingPolicyExists(resourceStepName, &steppolicy),
					resource.TestCheckResourceAttr(resourceSimpleName, "cooldown", "0"),
					resource.TestCheckResourceAttr(resourceSimpleName, "scaling_adjustment", "0"),
					resource.TestCheckResourceAttr(resourceStepName, "min_adjustment_magnitude", "1"),
					resource.TestCheckResourceAttr(resourceStepName, "estimated_instance_warmup", "0"),
				),
			},
			{
				ResourceName:      resourceSimpleName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSAutoscalingPolicyImportStateIdFunc(resourceSimpleName),
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceStepName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSAutoscalingPolicyImportStateIdFunc(resourceStepName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckScalingPolicyExists(n string, policy *autoscaling.ScalingPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*AWSClient).autoscalingconn
		params := &autoscaling.DescribePoliciesInput{
			AutoScalingGroupName: aws.String(rs.Primary.Attributes["autoscaling_group_name"]),
			PolicyNames:          []*string{aws.String(rs.Primary.ID)},
		}
		resp, err := conn.DescribePolicies(params)
		if err != nil {
			return err
		}
		if len(resp.ScalingPolicies) == 0 {
			return fmt.Errorf("ScalingPolicy not found")
		}

		*policy = *resp.ScalingPolicies[0]

		return nil
	}
}

func testAccCheckAWSAutoscalingPolicyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).autoscalingconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_autoscaling_group" {
			continue
		}

		params := autoscaling.DescribePoliciesInput{
			AutoScalingGroupName: aws.String(rs.Primary.Attributes["autoscaling_group_name"]),
			PolicyNames:          []*string{aws.String(rs.Primary.ID)},
		}

		resp, err := conn.DescribePolicies(&params)

		if err == nil {
			if len(resp.ScalingPolicies) != 0 &&
				*resp.ScalingPolicies[0].PolicyName == rs.Primary.ID {
				return fmt.Errorf("Scaling Policy Still Exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccAWSAutoscalingPolicyImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["autoscaling_group_name"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccAWSAutoscalingPolicyConfig_base(name string) string {
	return fmt.Sprintf(`
data "aws_ami" "amzn" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_launch_configuration" "test" {
  name          = "%s"
  image_id      = data.aws_ami.amzn.id
  instance_type = "t2.micro"
}

resource "aws_autoscaling_group" "test" {
  availability_zones   = slice(data.aws_availability_zones.available.names, 0, 2)
  name                 = "%s"
  max_size             = 0
  min_size             = 0
  force_delete         = true
  launch_configuration = aws_launch_configuration.test.name
}
`, name, name)
}

func testAccAWSAutoscalingPolicyConfig_basic(name string) string {
	return testAccAWSAutoscalingPolicyConfig_base(name) + fmt.Sprintf(`
resource "aws_autoscaling_policy" "foobar_simple" {
  name                   = "%s-foobar_simple"
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 300
  policy_type            = "SimpleScaling"
  scaling_adjustment     = 2
  autoscaling_group_name = aws_autoscaling_group.test.name
}

resource "aws_autoscaling_policy" "foobar_step" {
  name                      = "%s-foobar_step"
  adjustment_type           = "ChangeInCapacity"
  policy_type               = "StepScaling"
  estimated_instance_warmup = 200
  metric_aggregation_type   = "Minimum"

  step_adjustment {
    scaling_adjustment          = 1
    metric_interval_lower_bound = 2.0
  }

  autoscaling_group_name = aws_autoscaling_group.test.name
}

resource "aws_autoscaling_policy" "foobar_target_tracking" {
  name                   = "%s-foobar_target_tracking"
  policy_type            = "TargetTrackingScaling"
  autoscaling_group_name = aws_autoscaling_group.test.name

  target_tracking_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ASGAverageCPUUtilization"
    }

    target_value = 40.0
  }
}
`, name, name, name)
}

func testAccAWSAutoscalingPolicyConfig_basicUpdate(name string) string {
	return testAccAWSAutoscalingPolicyConfig_base(name) + fmt.Sprintf(`
resource "aws_autoscaling_policy" "foobar_simple" {
  name                   = "%s-foobar_simple"
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 30
  policy_type            = "SimpleScaling"
  scaling_adjustment     = 2
  autoscaling_group_name = aws_autoscaling_group.test.name
}

resource "aws_autoscaling_policy" "foobar_step" {
  name                      = "%s-foobar_step"
  adjustment_type           = "ChangeInCapacity"
  policy_type               = "StepScaling"
  estimated_instance_warmup = 20
  metric_aggregation_type   = "Minimum"

  step_adjustment {
    scaling_adjustment          = 10
    metric_interval_lower_bound = 2.0
  }

  autoscaling_group_name = aws_autoscaling_group.test.name
}

resource "aws_autoscaling_policy" "foobar_target_tracking" {
  name                   = "%s-foobar_target_tracking"
  policy_type            = "TargetTrackingScaling"
  autoscaling_group_name = aws_autoscaling_group.test.name

  target_tracking_configuration {
    customized_metric_specification {
      metric_dimension {
        name  = "fuga"
        value = "fuga"
      }

      metric_name = "hoge"
      namespace   = "hoge"
      statistic   = "Average"
    }

    target_value = 70.0
  }
}
`, name, name, name)
}

func testAccAWSAutoscalingPolicyConfig_SimpleScalingStepAdjustment(name string) string {
	return testAccAWSAutoscalingPolicyConfig_base(name) + fmt.Sprintf(`
resource "aws_autoscaling_policy" "foobar_simple" {
  name                   = "%s-foobar_simple"
  adjustment_type        = "ExactCapacity"
  cooldown               = 300
  policy_type            = "SimpleScaling"
  scaling_adjustment     = 0
  autoscaling_group_name = aws_autoscaling_group.test.name
}
`, name)
}

func testAccAwsAutoscalingPolicyConfig_TargetTracking_Predefined(name string) string {
	return testAccAWSAutoscalingPolicyConfig_base(name) + fmt.Sprintf(`
resource "aws_autoscaling_policy" "test" {
  name                   = "%s-test"
  policy_type            = "TargetTrackingScaling"
  autoscaling_group_name = aws_autoscaling_group.test.name

  target_tracking_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ASGAverageCPUUtilization"
    }

    target_value = 40.0
  }
}
`, name)
}

func testAccAwsAutoscalingPolicyConfig_TargetTracking_Custom(name string) string {
	return testAccAWSAutoscalingPolicyConfig_base(name) + fmt.Sprintf(`
resource "aws_autoscaling_policy" "test" {
  name                   = "%s-test"
  policy_type            = "TargetTrackingScaling"
  autoscaling_group_name = aws_autoscaling_group.test.name

  target_tracking_configuration {
    customized_metric_specification {
      metric_dimension {
        name  = "fuga"
        value = "fuga"
      }

      metric_name = "hoge"
      namespace   = "hoge"
      statistic   = "Average"
    }

    target_value = 40.0
  }
}
`, name)
}

func testAccAWSAutoscalingPolicyConfig_zerovalue(name string) string {
	return testAccAWSAutoscalingPolicyConfig_base(name) + fmt.Sprintf(`
resource "aws_autoscaling_policy" "foobar_simple" {
  name                   = "%s-foobar_simple"
  adjustment_type        = "ExactCapacity"
  cooldown               = 0
  policy_type            = "SimpleScaling"
  scaling_adjustment     = 0
  autoscaling_group_name = aws_autoscaling_group.test.name
}

resource "aws_autoscaling_policy" "foobar_step" {
  name                      = "%s-foobar_step"
  adjustment_type           = "PercentChangeInCapacity"
  policy_type               = "StepScaling"
  estimated_instance_warmup = 0
  metric_aggregation_type   = "Minimum"

  step_adjustment {
    scaling_adjustment          = 1
    metric_interval_lower_bound = 2.0
  }

  min_adjustment_magnitude = 1
  autoscaling_group_name   = aws_autoscaling_group.test.name
}
`, name, name)
}
