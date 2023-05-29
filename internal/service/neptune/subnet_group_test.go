// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package neptune_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/neptune"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfneptune "github.com/hashicorp/terraform-provider-aws/internal/service/neptune"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccNeptuneSubnetGroup_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBSubnetGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_subnet_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSubnetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetGroupConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSubnetGroupExists(ctx, resourceName, &v),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "rds", regexp.MustCompile(fmt.Sprintf("subgrp:%s$", rName))),
					resource.TestCheckResourceAttr(resourceName, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", ""),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNeptuneSubnetGroup_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBSubnetGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_subnet_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSubnetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetGroupExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfneptune.ResourceSubnetGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccNeptuneSubnetGroup_nameGenerated(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBSubnetGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_subnet_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSubnetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetGroupConfig_nameGenerated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetGroupExists(ctx, resourceName, &v),
					acctest.CheckResourceAttrNameGenerated(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", id.UniqueIdPrefix),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNeptuneSubnetGroup_namePrefix(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBSubnetGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_subnet_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSubnetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetGroupConfig_namePrefix(rName, "tf-acc-test-prefix-"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetGroupExists(ctx, resourceName, &v),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, "name", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-test-prefix-"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNeptuneSubnetGroup_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBSubnetGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_subnet_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSubnetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetGroupConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetGroupExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSubnetGroupConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetGroupExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccSubnetGroupConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetGroupExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccNeptuneSubnetGroup_update(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBSubnetGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_subnet_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSubnetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetGroupExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "2"),
				),
			},
			{
				Config: testAccSubnetGroupConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetGroupExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "description", "test description updated"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "3"),
				),
			},
		},
	})
}

func testAccCheckSubnetGroupDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).NeptuneConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_neptune_subnet_group" {
				continue
			}

			_, err := tfneptune.FindSubnetGroupByName(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Neptune Subnet Group %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckSubnetGroupExists(ctx context.Context, n string, v *neptune.DBSubnetGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Neptune Subnet Group ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).NeptuneConn()

		output, err := tfneptune.FindSubnetGroupByName(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccSubnetGroupConfig_basic(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 2), fmt.Sprintf(`
resource "aws_neptune_subnet_group" "test" {
  name       = %[1]q
  subnet_ids = aws_subnet.test[*].id
}
`, rName))
}

func testAccSubnetGroupConfig_nameGenerated(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 2), `
resource "aws_neptune_subnet_group" "test" {
  subnet_ids = aws_subnet.test[*].id
}`)
}

func testAccSubnetGroupConfig_namePrefix(rName, prefix string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 2), fmt.Sprintf(`
resource "aws_neptune_subnet_group" "test" {
  name_prefix = %[1]q
  subnet_ids  = aws_subnet.test[*].id
}`, prefix))
}

func testAccSubnetGroupConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 2), fmt.Sprintf(`
resource "aws_neptune_subnet_group" "test" {
  name       = %[1]q
  subnet_ids = aws_subnet.test[*].id

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccSubnetGroupConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 2), fmt.Sprintf(`
resource "aws_neptune_subnet_group" "test" {
  name       = %[1]q
  subnet_ids = aws_subnet.test[*].id

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccSubnetGroupConfig_update(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 3), fmt.Sprintf(`
resource "aws_neptune_subnet_group" "test" {
  name        = %[1]q
  subnet_ids  = aws_subnet.test[*].id
  description = "test description updated"
}
`, rName))
}
