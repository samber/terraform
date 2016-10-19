// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package rds

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol/query"
	"github.com/aws/aws-sdk-go/private/signer/v4"
)

// Amazon Relational Database Service (Amazon RDS) is a web service that makes
// it easier to set up, operate, and scale a relational database in the cloud.
// It provides cost-efficient, resizeable capacity for an industry-standard
// relational database and manages common database administration tasks, freeing
// up developers to focus on what makes their applications and businesses unique.
//
//  Amazon RDS gives you access to the capabilities of a MySQL, MariaDB, PostgreSQL,
// Microsoft SQL Server, Oracle, or Amazon Aurora database server. These capabilities
// mean that the code, applications, and tools you already use today with your
// existing databases work with Amazon RDS without modification. Amazon RDS
// automatically backs up your database and maintains the database software
// that powers your DB instance. Amazon RDS is flexible: you can scale your
// database instance's compute resources and storage capacity to meet your application's
// demand. As with all Amazon Web Services, there are no up-front investments,
// and you pay only for the resources you use.
//
//  This interface reference for Amazon RDS contains documentation for a programming
// or command line interface you can use to manage Amazon RDS. Note that Amazon
// RDS is asynchronous, which means that some interfaces might require techniques
// such as polling or callback functions to determine when a command has been
// applied. In this reference, the parameter descriptions indicate whether a
// command is applied immediately, on the next instance reboot, or during the
// maintenance window. The reference structure is as follows, and we list following
// some related topics from the user guide.
//
// Amazon RDS API Reference
//
//  For the alphabetical list of API actions, see API Actions (http://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_Operations.html).
//
// For the alphabetical list of data types, see Data Types (http://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_Types.html).
//
// For a list of common query parameters, see Common Parameters (http://docs.aws.amazon.com/AmazonRDS/latest/APIReference/CommonParameters.html).
//
// For descriptions of the error codes, see Common Errors (http://docs.aws.amazon.com/AmazonRDS/latest/APIReference/CommonErrors.html).
//
//  Amazon RDS User Guide
//
//  For a summary of the Amazon RDS interfaces, see Available RDS Interfaces
// (http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Welcome.html#Welcome.Interfaces).
//
// For more information about how to use the Query API, see Using the Query
// API (http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Using_the_Query_API.html).
//The service client's operations are safe to be used concurrently.
// It is not safe to mutate any of the client's properties though.
type RDS struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// A ServiceName is the name of the service the client will make API calls to.
const ServiceName = "rds"

// New creates a new instance of the RDS client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a RDS client from just a session.
//     svc := rds.New(mySession)
//
//     // Create a RDS client with additional configuration
//     svc := rds.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) *RDS {
	c := p.ClientConfig(ServiceName, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion string) *RDS {
	svc := &RDS{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2014-10-31",
			},
			handlers,
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBack(v4.Sign)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

// newRequest creates a new request for a RDS operation and runs any
// custom request initialization.
func (c *RDS) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
