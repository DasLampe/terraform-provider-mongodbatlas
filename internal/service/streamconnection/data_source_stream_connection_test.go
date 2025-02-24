package streamconnection_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccStreamDSStreamConnection_kafkaPlaintext(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc-stream")
		instanceName   = acctest.RandomWithPrefix("test-acc-instance")
		dataSourceName = "data.mongodbatlas_stream_connection.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBetaFlag(t); acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: streamConnectionDataSourceConfig(kafkaStreamConnectionConfig(orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", false)),
				Check:  kafkaStreamConnectionAttributeChecks(dataSourceName, orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092,localhost:9092", "earliest", false, false),
			},
		},
	})
}

func TestAccStreamDSStreamConnection_kafkaSSL(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc-stream")
		instanceName   = acctest.RandomWithPrefix("test-acc-instance")
		dataSourceName = "data.mongodbatlas_stream_connection.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBetaFlag(t); acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: streamConnectionDataSourceConfig(kafkaStreamConnectionConfig(orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", true)),
				Check:  kafkaStreamConnectionAttributeChecks(dataSourceName, orgID, projectName, instanceName, "user", "rawpassword", "localhost:9092", "earliest", true, false),
			},
		},
	})
}

func TestAccStreamDSStreamConnection_cluster(t *testing.T) {
	var (
		clusterInfo    = acc.GetClusterInfo(nil)
		instanceName   = acctest.RandomWithPrefix("test-acc-name")
		dataSourceName = "data.mongodbatlas_stream_connection.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBetaFlag(t); acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             CheckDestroyStreamConnection,
		Steps: []resource.TestStep{
			{
				Config: streamConnectionDataSourceConfig(clusterStreamConnectionConfig(clusterInfo.ProjectIDStr, instanceName, clusterInfo.ClusterNameStr, clusterInfo.ClusterTerraformStr)),
				Check:  clusterStreamConnectionAttributeChecks(dataSourceName, clusterInfo.ClusterName),
			},
		},
	})
}

func streamConnectionDataSourceConfig(streamConnectionConfig string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_connection" "test" {
			project_id = mongodbatlas_stream_connection.test.project_id
			instance_name = mongodbatlas_stream_connection.test.instance_name
			connection_name = mongodbatlas_stream_connection.test.connection_name
		}
	`, streamConnectionConfig)
}
