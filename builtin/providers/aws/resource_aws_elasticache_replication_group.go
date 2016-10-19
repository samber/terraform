package aws

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsElasticacheReplicationGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsElasticacheReplicationGroupCreate,
		Read:   resourceAwsElasticacheReplicationGroupRead,
		Update: resourceAwsElasticacheReplicationGroupUpdate,
		Delete: resourceAwsElasticacheReplicationGroupDelete,

		Schema: map[string]*schema.Schema{
			"replication_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"replication_group_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"automatic_failover_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"number_cache_clusters": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"preferred_cache_cluster_azs": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"cache_node_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "redis",
			},
			"engine_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cache_parameter_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cache_subnet_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cache_security_group_names": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"security_group_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"snapshot_arns": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"snapshot_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"preferred_maintenance_window": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"notification_topic_arn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"notification_topic_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auto_minor_version_upgrade": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"snapshot_retention_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"snapshot_window": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"apply_immediately": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"snapshotting_cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cache_node_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_group_members": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"address": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"availability_zone": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"current_role": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAwsElasticacheReplicationGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).elasticacheconn

	params := &elasticache.CreateReplicationGroupInput{
		ReplicationGroupId:          aws.String(d.Get("replication_group_id").(string)),
		ReplicationGroupDescription: aws.String(d.Get("replication_group_description").(string)),
		Engine: aws.String(d.Get("engine").(string)),
	}

	if v, ok := d.GetOk("primary_cluster_id"); ok {
		params.PrimaryClusterId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("automatic_failover_enabled"); ok {
		params.AutomaticFailoverEnabled = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("number_cache_clusters"); ok {
		params.NumCacheClusters = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("preferred_cache_cluster_azs"); ok {
		params.PreferredCacheClusterAZs = []*string{aws.String(v.(string))}
	}

	if v, ok := d.GetOk("cache_node_type"); ok {
		params.CacheNodeType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("engine_version"); ok {
		params.EngineVersion = aws.String(v.(string))
	}

	if v, ok := d.GetOk("cache_parameter_group_name"); ok {
		params.CacheParameterGroupName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("cache_subnet_group_name"); ok {
		params.CacheSubnetGroupName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("cache_security_group_names"); ok {
		params.CacheSecurityGroupNames = expandStringList(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		params.SecurityGroupIds = expandStringList(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("snapshot_arns"); ok {
		params.SnapshotArns = expandStringList(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("snapshot_name"); ok {
		params.SnapshotName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("preferred_maintenance_window"); ok {
		params.PreferredMaintenanceWindow = aws.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		params.Port = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("notification_topic_arn"); ok {
		params.NotificationTopicArn = aws.String(v.(string))
	}

	if v, ok := d.GetOk("auto_minor_version_upgrade"); ok {
		params.AutoMinorVersionUpgrade = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("snapshot_retention_limit"); ok {
		params.SnapshotRetentionLimit = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("snapshot_window"); ok {
		params.SnapshotWindow = aws.String(v.(string))
	}

	resp, err := conn.CreateReplicationGroup(params)
	if err != nil {
		return fmt.Errorf("Error creating Elasticache Replication Group: %s", err)
	}

	d.SetId(*resp.ReplicationGroup.ReplicationGroupId)

	pending := []string{"creating"}
	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{"available"},
		Refresh:    cacheClusterReplicationGroupStateRefreshFunc(conn, d.Id(), "available", pending),
		Timeout:    20 * time.Minute, // These can take a while
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	log.Printf("[DEBUG] Waiting for state to become available: %v", d.Id())
	_, sterr := stateConf.WaitForState()
	if sterr != nil {
		return fmt.Errorf("Error waiting for elasticache replication group (%s) to be created: %s", d.Id(), sterr)
	}

	return resourceAwsElasticacheReplicationGroupRead(d, meta)
}

func resourceAwsElasticacheReplicationGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).elasticacheconn
	req := &elasticache.DescribeReplicationGroupsInput{
		ReplicationGroupId: aws.String(d.Id()),
	}

	res, err := conn.DescribeReplicationGroups(req)
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	if err != nil {
		if eccErr, ok := err.(awserr.Error); ok && eccErr.Code() == "ReplicationGroupNotFound" {
			log.Printf("[WARN] Elasticache Replication Group (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}
<<<<<<< HEAD
	//	if len(res.ReplicationGroups) == 1 {
	//		rg := res.ReplicationGroups[0]
	//	}
	return nil
}

=======

	if len(res.ReplicationGroups) == 1 {
		rg := res.ReplicationGroups[0]

		d.Set("number_cache_clusters", len(rg.MemberClusters))
		d.Set("automatic_failover_enabled", *rg.AutomaticFailover == "enabled" || *rg.AutomaticFailover == "enabling")

		if err := setCacheNodeGroupsData(d, rg); err != nil {
			return err
		}

	}
	return nil
}

func setCacheNodeGroupsData(d *schema.ResourceData, rg *elasticache.ReplicationGroup) error {
	sortedCacheNodes := make([]*elasticache.NodeGroup, len(rg.NodeGroups))
	copy(sortedCacheNodes, rg.NodeGroups)
	sort.Sort(byNodeGroupId(sortedCacheNodes))

	cacheNodeGroupData := make([]map[string]interface{}, 0, len(sortedCacheNodes))

	for _, node := range sortedCacheNodes {
		if node.NodeGroupId == nil || node.PrimaryEndpoint == nil || node.PrimaryEndpoint.Address == nil || node.PrimaryEndpoint.Port == nil {
			return fmt.Errorf("Unexpected nil pointer in: %s", node)
		}

		nodeData := map[string]interface{}{
			"id":      *node.NodeGroupId,
			"address": *node.PrimaryEndpoint.Address,
			"port":    int(*node.PrimaryEndpoint.Port),
		}

		setCacheNodeMemberData(d, nodeData, node)

		cacheNodeGroupData = append(cacheNodeGroupData, nodeData)

	}

	return d.Set("cache_node_groups", cacheNodeGroupData)
}

func setCacheNodeMemberData(d *schema.ResourceData, group map[string]interface{}, rg *elasticache.NodeGroup) {

	cacheNodeGroupMembers := make([]map[string]interface{}, 0, len(rg.NodeGroupMembers))

	for _, nodeMember := range rg.NodeGroupMembers {

		nodeMemberData := map[string]interface{}{
			"id":                *nodeMember.CacheNodeId,
			"address":           *nodeMember.ReadEndpoint.Address,
			"port":              *nodeMember.ReadEndpoint.Port,
			"availability_zone": *nodeMember.PreferredAvailabilityZone,
			"current_role":      *nodeMember.CurrentRole,
		}

		cacheNodeGroupMembers = append(cacheNodeGroupMembers, nodeMemberData)
	}

	group["node_group_members"] = cacheNodeGroupMembers
}

type byNodeGroupId []*elasticache.NodeGroup

func (b byNodeGroupId) Len() int      { return len(b) }
func (b byNodeGroupId) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byNodeGroupId) Less(i, j int) bool {
	return b[i].NodeGroupId != nil && b[j].NodeGroupId != nil &&
		*b[i].NodeGroupId < *b[j].NodeGroupId
}

>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
func resourceAwsElasticacheReplicationGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).elasticacheconn

	params := &elasticache.ModifyReplicationGroupInput{
		ApplyImmediately:   aws.Bool(d.Get("apply_immediately").(bool)),
		ReplicationGroupId: aws.String(d.Id()),
	}

<<<<<<< HEAD
	if d.HasChange("replication_group_description") {
		params.ReplicationGroupDescription = aws.String(d.Get("description").(string))
=======
	requestUpdate := false

	if d.HasChange("replication_group_description") {
		params.ReplicationGroupDescription = aws.String(d.Get("description").(string))
		requestUpdate = true
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	if d.HasChange("primary_cluster_id") {
		params.PrimaryClusterId = aws.String(d.Get("primary_cluster_id").(string))
<<<<<<< HEAD
=======
		requestUpdate = true
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	if d.HasChange("snapshotting_cluster_id") {
		params.SnapshottingClusterId = aws.String(d.Get("snapshotting_cluster_id").(string))
<<<<<<< HEAD
	}

	if d.HasChange("automatic_failover_enabled") {
		params.AutomaticFailoverEnabled = aws.Bool(d.Get("automatic_failover").(bool))
=======
		requestUpdate = true
	}

	if d.HasChange("automatic_failover_enabled") {
		params.AutomaticFailoverEnabled = aws.Bool(d.Get("automatic_failover_enabled").(bool))
		requestUpdate = true
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	if d.HasChange("cache_security_group_names") {
		params.CacheSecurityGroupNames = expandStringList(d.Get("cache_security_group_names").(*schema.Set).List())
<<<<<<< HEAD
=======
		requestUpdate = true
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	if d.HasChange("security_group_ids") {
		params.SecurityGroupIds = expandStringList(d.Get("security_group_ids").(*schema.Set).List())
<<<<<<< HEAD
=======
		requestUpdate = true
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	if d.HasChange("preferred_maintenance_window") {
		params.PreferredMaintenanceWindow = aws.String(d.Get("preferred_maintenance_window").(string))
<<<<<<< HEAD
=======
		requestUpdate = true
	}

	if d.HasChange("cache_node_type") {
		params.CacheNodeType = aws.String(d.Get("cache_node_type").(string))
		requestUpdate = true
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	if d.HasChange("notification_topic_arn") {
		params.NotificationTopicArn = aws.String(d.Get("notification_topic_arn").(string))
<<<<<<< HEAD
=======
		requestUpdate = true
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	if d.HasChange("cache_parameter_group_name") {
		params.CacheParameterGroupName = aws.String(d.Get("cache_parameter_group_name").(string))
<<<<<<< HEAD
=======
		requestUpdate = true
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	if d.HasChange("notification_topic_status") {
		params.NotificationTopicStatus = aws.String(d.Get("notification_topic_status").(string))
<<<<<<< HEAD
=======
		requestUpdate = true
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	if d.HasChange("engine_version") {
		params.EngineVersion = aws.String(d.Get("engine_version").(string))
<<<<<<< HEAD
=======
		requestUpdate = true
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	if d.HasChange("auto_minor_version_upgrade") {
		params.AutoMinorVersionUpgrade = aws.Bool(d.Get("auto_minor_version_upgrade").(bool))
<<<<<<< HEAD
=======
		requestUpdate = true
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	if d.HasChange("snapshot_retention_limit") {
		params.SnapshotRetentionLimit = aws.Int64(int64(d.Get("snapshot_retention_limit").(int)))
<<<<<<< HEAD
=======
		requestUpdate = true
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	if d.HasChange("snapshot_window") {
		params.SnapshotWindow = aws.String(d.Get("snapshot_window").(string))
<<<<<<< HEAD
	}

	_, err := conn.ModifyReplicationGroup(params)
	if err != nil {
		return fmt.Errorf("Error updating Elasticache replication group: %s", err)
=======
		requestUpdate = true
	}

	if requestUpdate {
		log.Printf("[DEBUG] Modifying ElastiCache Replication Group (%s), opts:\n%s", d.Id(), params)
		_, err := conn.ModifyReplicationGroup(params)
		if err != nil {
			return fmt.Errorf("Error updating Elasticache replication group: %s", err)
		}

		log.Printf("[DEBUG] Waiting for update: %s", d.Id())
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"creating", "available", "deleting", "modifying"},
			Target:     []string{"available"},
			Refresh:    cacheClusterReplicationGroupStateRefreshFunc(conn, d.Id(), "", []string{}),
			Timeout:    15 * time.Minute,
			Delay:      20 * time.Second,
			MinTimeout: 5 * time.Second,
		}

		_, sterr := stateConf.WaitForState()
		if sterr != nil {
			return fmt.Errorf("Error waiting for elasticache (%s) to update: %s", d.Id(), sterr)
		}
	}

	if d.HasChange("number_cache_clusters") {
		// Need to calculate adding or removing cache clusters

		oraw, nraw := d.GetChange("number_cache_clusters")
		o := oraw.(int)
		n := nraw.(int)

		if n < o { // Removing
			log.Printf("[INFO] Cluster %s is marked for Decreasing cache nodes from %d to %d", d.Id(), o, n)
			// TODO: Determine which non-primary clusters to remove

			clustersToRemove := getClustersToRemove(d, o, o-n)

			for _, clusterId := range clustersToRemove {
				deleteCacheCluster(conn, d, clusterId)
			}

		} else { // if o > n {	// Adding
			log.Printf("[INFO] Cluster %s is marked for Increasing cache nodes from %d to %d", d.Id(), o, n)
			// TODO: Add new cache clusters to this replication group (and calculate the naming requirement?)

			for i := o; i < n; i++ {
				err := createCacheCluster(conn, d, i)
				if err != nil {
					return err
				}
			}

		}

>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
	}

	return resourceAwsElasticacheReplicationGroupRead(d, meta)
}

<<<<<<< HEAD
=======
func getClustersToRemove(d *schema.ResourceData, oldNumberOfClusters int, clustersToRemove int) []*string {
	clustersIdsToRemove := []*string{}
	for i := oldNumberOfClusters; i > oldNumberOfClusters-clustersToRemove && i > 0; i-- {
		s := fmt.Sprintf("%s-%03d", d.Get("replication_group_id").(string), i)
		clustersIdsToRemove = append(clustersIdsToRemove, &s)
	}

	return clustersIdsToRemove
}

func createCacheCluster(conn *elasticache.ElastiCache, d *schema.ResourceData, clusterNum int) error {

	// this feels clumsy
	cacheClusterId := fmt.Sprintf("%s-%03d", d.Get("replication_group_id").(string), clusterNum+1)

	createCacheRequest := &elasticache.CreateCacheClusterInput{
		CacheClusterId:     aws.String(cacheClusterId),
		ReplicationGroupId: aws.String(d.Id()),
	}

	resp, err := conn.CreateCacheCluster(createCacheRequest)
	if err != nil {
		return fmt.Errorf("Error creating Elasticache: %s", err)
	}

	pending := []string{"creating"}
	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{"available"},
		Refresh:    cacheClusterStateRefreshFunc(conn, *resp.CacheCluster.CacheClusterId, "available", pending),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	log.Printf("[DEBUG] Waiting for state to become available: %v", d.Id())
	_, sterr := stateConf.WaitForState()
	if sterr != nil {
		return fmt.Errorf("Error waiting for elasticache (%s) to be created: %s", d.Id(), sterr)
	}

	return nil
}

func deleteCacheCluster(conn *elasticache.ElastiCache, d *schema.ResourceData, cacheClusterId *string) error {

	deleteClusterRequest := &elasticache.DeleteCacheClusterInput{
		CacheClusterId: aws.String(*cacheClusterId),
	}

	resp, err := conn.DeleteCacheCluster(deleteClusterRequest)
	if err != nil {
		return fmt.Errorf("Error deleting Elasticache cluster: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "available", "deleting"},
		Target:     []string{""},
		Refresh:    cacheClusterStateRefreshFunc(conn, *resp.CacheCluster.CacheClusterId, "", []string{}),
		Timeout:    15 * time.Minute,
		Delay:      20 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	log.Printf("[DEBUG] Waiting for state to become available: %v", d.Id())
	_, sterr := stateConf.WaitForState()
	if sterr != nil {
		return fmt.Errorf("Error waiting for elasticache (%s) to delete: %s", d.Id(), sterr)
	}

	return nil
}

>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
func resourceAwsElasticacheReplicationGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).elasticacheconn

	req := &elasticache.DeleteReplicationGroupInput{
		ReplicationGroupId: aws.String(d.Id()),
	}

	_, err := conn.DeleteReplicationGroup(req)
	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "ReplicationGroupNotFoundFault" {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error deleting Elasticache replication group: %s", err)
	}

	log.Printf("[DEBUG] Waiting for deletion: %v", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "available", "deleting"},
<<<<<<< HEAD
		Target:     "",
=======
		Target:     []string{},
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
		Refresh:    cacheClusterReplicationGroupStateRefreshFunc(conn, d.Id(), "", []string{}),
		Timeout:    15 * time.Minute,
		Delay:      20 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, sterr := stateConf.WaitForState()
	if sterr != nil {
		return fmt.Errorf("Error waiting for replication group (%s) to delete: %s", d.Id(), sterr)
	}

	return nil
}

func cacheClusterReplicationGroupStateRefreshFunc(conn *elasticache.ElastiCache, replicationGroupId, givenState string, pending []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.DescribeReplicationGroups(&elasticache.DescribeReplicationGroupsInput{
			ReplicationGroupId: aws.String(replicationGroupId),
		})
		if err != nil {
			apierr := err.(awserr.Error)
			log.Printf("[DEBUG] message: %v, code: %v", apierr.Message(), apierr.Code())
<<<<<<< HEAD
			if apierr.Message() == fmt.Sprintf("Cluster ReplicationGroup not found: %v", replicationGroupId) {
=======

			if apierr.Code() == "ReplicationGroupNotFoundFault" {
>>>>>>> 14800b4f1ee2f0d5d72121b3117c707a9453c6ba
				log.Printf("[DEBUG] Detect deletion")
				return nil, "", nil
			}

			log.Printf("[ERROR] cacheClusterReplicationGroupStateRefreshFunc: %s", err)
			return nil, "", err
		}

		if len(resp.ReplicationGroups) == 0 {
			return nil, "", fmt.Errorf("[WARN] Error: no Cache Replication Groups found for id (%s)", replicationGroupId)
		}

		var rg *elasticache.ReplicationGroup
		for _, replicationGroup := range resp.ReplicationGroups {
			if *replicationGroup.ReplicationGroupId == replicationGroupId {
				log.Printf("[DEBUG] Found matching ElastiCache Replication Group: %s", *replicationGroup.ReplicationGroupId)
				rg = replicationGroup
			}
		}

		if rg == nil {
			return nil, "", fmt.Errorf("[WARN] Error: no matching Elasticcache Replication Group for id (%s)", replicationGroupId)
		}

		log.Printf("[DEBUG] ElastiCache Replication Group (%s) status: %v", replicationGroupId, *rg.Status)

		// return the current state if it's in the pending array
		for _, p := range pending {
			log.Printf("[DEBUG] ElastiCache: checking pending state (%s) for Replication Group (%s), Replication Group status: %s", pending, replicationGroupId, *rg.Status)
			s := *rg.Status
			if p == s {
				log.Printf("[DEBUG] Return with status: %v", *rg.Status)
				return s, p, nil
			}
		}

		//		// return given state if it's not in pending
		//		if givenState != "" {
		//			log.Printf("[DEBUG] ElastiCache: checking given state (%s) of Replication Group (%s) against Replication Group status (%s)", givenState, replicationGroupId, *rg.Status)
		//			// check to make sure we have the node count we're expecting
		//			if int64(len(rg.)) != *c.NumCacheNodes {
		//				log.Printf("[DEBUG] Node count is not what is expected: %d found, %d expected", len(c.CacheNodes), *c.NumCacheNodes)
		//				return nil, "creating", nil
		//			}
		//
		//			log.Printf("[DEBUG] Node count matched (%d)", len(c.CacheNodes))
		//			// loop the nodes and check their status as well
		//			for _, n := range c.CacheNodes {
		//				log.Printf("[DEBUG] Checking cache node for status: %s", n)
		//				if n.CacheNodeStatus != nil && *n.CacheNodeStatus != "available" {
		//					log.Printf("[DEBUG] Node (%s) is not yet available, status: %s", *n.CacheNodeId, *n.CacheNodeStatus)
		//					return nil, "creating", nil
		//				}
		//				log.Printf("[DEBUG] Cache node not in expected state")
		//			}
		//			log.Printf("[DEBUG] ElastiCache returning given state (%s), cluster: %s", givenState, c)
		//			return c, givenState, nil
		//		}
		//		log.Printf("[DEBUG] current status: %v", *c.CacheClusterStatus)
		return rg, *rg.Status, nil
	}
}

func buildECReplicationGroupARN(d *schema.ResourceData, meta interface{}) (string, error) {
	iamconn := meta.(*AWSClient).iamconn
	region := meta.(*AWSClient).region
	// An zero value GetUserInput{} defers to the currently logged in user
	resp, err := iamconn.GetUser(&iam.GetUserInput{})
	if err != nil {
		return "", err
	}
	userARN := *resp.User.Arn
	accountID := strings.Split(userARN, ":")[4]
	arn := fmt.Sprintf("arn:aws:elasticache:%s:%s:cluster:%s", region, accountID, d.Id())
	return arn, nil
}
