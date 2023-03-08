package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfraStackProps struct {
	awscdk.StackProps
}

func SpaceCloudInfraStack(scope constructs.Construct, id string, props *InfraStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// S3 bucket for storage
	bucket := awss3.NewBucket(
		stack,
		jsii.String("space_cloud_bucket"),
		&awss3.BucketProps{
			BucketName:    jsii.String("spaceclouddatabucket"),
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		},
	)

	// Role for lambda functions
	s3LambdaRole := awsiam.NewRole(
		stack,
		jsii.String("S3LambdaRole"),
		&awsiam.RoleProps{
			AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
			ManagedPolicies: &[]awsiam.IManagedPolicy{
				awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonDynamoDBFullAccess")),
				awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("CloudWatchFullAccess")),
			},
		})

	// Get people in space function.
	collectPeopleFunction := awslambda.NewFunction(stack, jsii.String("GetPeople"), &awslambda.FunctionProps{
		FunctionName: jsii.String(*stack.StackName() + "-GetPeople"),
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(60)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("../out/."), nil),
		Handler:      jsii.String("collectPeopleInSpace"),
		Architecture: awslambda.Architecture_X86_64(),
		Role:         s3LambdaRole,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
		CurrentVersionOptions: &awslambda.VersionOptions{
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		},
		Environment: &map[string]*string{
			"DATA_BUCKET": jsii.String(*bucket.BucketName()),
			"BUCKET_KEY":  jsii.String("people_in_space.json"),
		},
	})

	// Read people in space function.
	readPeopleFunction := awslambda.NewFunction(stack, jsii.String("ReadPeople"), &awslambda.FunctionProps{
		FunctionName: jsii.String(*stack.StackName() + "-ReadPeople"),
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(60)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("../out/."), nil),
		Handler:      jsii.String("readPeopleInSpace"),
		Architecture: awslambda.Architecture_X86_64(),
		Role:         s3LambdaRole,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
		CurrentVersionOptions: &awslambda.VersionOptions{
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		},
		Environment: &map[string]*string{
			"DATA_BUCKET": jsii.String(*bucket.BucketName()),
			"BUCKET_KEY":  jsii.String("people_in_space.json"),
		},
	})

	// Read people in space rust function
	readPeopleRustFunction := awslambda.Function_FromFunctionArn(stack, jsii.String("ReadPeopleRust"),
		jsii.String("arn:aws:lambda:us-east-1:939984321277:function:read-people-in-space-rust"))

	// Retrieve and store launch lambda function.
	collectLaunchFunction := awslambda.NewFunction(stack, jsii.String("GetLaunches"), &awslambda.FunctionProps{
		FunctionName: jsii.String(*stack.StackName() + "-GetLaunches"),
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(60)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("../out/."), nil),
		Handler:      jsii.String("collectNextRocketLaunches"),
		Architecture: awslambda.Architecture_X86_64(),
		Role:         s3LambdaRole,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
		CurrentVersionOptions: &awslambda.VersionOptions{
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		},
		Environment: &map[string]*string{
			"DATA_BUCKET": jsii.String(*bucket.BucketName()),
			"BUCKET_KEY":  jsii.String("launches.json"),
		},
	})

	// Add permissions to lambda functions
	// Read and write
	bucket.GrantReadWrite(collectPeopleFunction, nil)
	bucket.GrantReadWrite(collectLaunchFunction, nil)

	// Read only
	bucket.GrantRead(readPeopleFunction, nil)

	// "Daily" event - every 2 hours
	everyTwoHoursEventRule := awsevents.NewRule(
		stack,
		jsii.String("data_builder_event"),
		&awsevents.RuleProps{
			RuleName: jsii.String("dataBuilderEvent"),
			Enabled:  jsii.Bool(true),
			Schedule: awsevents.Schedule_Cron(&awsevents.CronOptions{
				Hour:   jsii.String("0/2"),
				Minute: jsii.String("0"),
			}),
		},
	)

	// Twice daily event rule
	twiceDailyEventRule := awsevents.NewRule(
		stack,
		jsii.String("data_builder_event_twice_daily"),
		&awsevents.RuleProps{
			RuleName: jsii.String("dataBuilderEventTwiceDaily"),
			Enabled:  jsii.Bool(true),
			Schedule: awsevents.Schedule_Cron(&awsevents.CronOptions{
				Hour:   jsii.String("0/12"),
				Minute: jsii.String("0"),
			}),
		},
	)

	// Add targets to event rule(s)
	everyTwoHoursEventRule.AddTarget(awseventstargets.NewLambdaFunction(collectPeopleFunction, nil))
	twiceDailyEventRule.AddTarget(awseventstargets.NewLambdaFunction(collectLaunchFunction, nil))

	// Create API Gateway
	restApiProd := awsapigateway.NewRestApi(
		stack,
		jsii.String("space_cloud_api"),
		&awsapigateway.RestApiProps{
			RestApiName:        jsii.String("Space Cloud API"),
			RetainDeployments:  jsii.Bool(false),
			EndpointExportName: jsii.String("SpaceCloudApiEndpoint"),
			Deploy:             jsii.Bool(true),
			EndpointConfiguration: &awsapigateway.EndpointConfiguration{
				Types: &[]awsapigateway.EndpointType{
					awsapigateway.EndpointType_REGIONAL,
				},
			},
			DeployOptions: &awsapigateway.StageOptions{
				StageName:            jsii.String("prod"),
				CacheClusterEnabled:  jsii.Bool(false),
				ThrottlingBurstLimit: jsii.Number(100),
				ThrottlingRateLimit:  jsii.Number(1000),
			},
			DomainName: &awsapigateway.DomainNameOptions{
				DomainName: jsii.String("api.spacebits.net"),
				Certificate: awscertificatemanager.Certificate_FromCertificateArn(
					stack,
					jsii.String("space_cloud_api_cert"),
					jsii.String("arn:aws:acm:us-east-1:939984321277:certificate/32c01289-82c7-4d30-885a-d5cd3aab4a93"),
				),
			},
			DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
				AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
				AllowMethods: awsapigateway.Cors_ALL_METHODS(),
				AllowHeaders: awsapigateway.Cors_DEFAULT_HEADERS(),
			},
		},
	)

	// Read people endpoint
	readPeopleResource := restApiProd.Root().AddResource(jsii.String("people"), nil)
	readPeopleResource.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(readPeopleFunction,
		&awsapigateway.LambdaIntegrationOptions{}),
		&awsapigateway.MethodOptions{
			ApiKeyRequired: jsii.Bool(true),
		})

	readPeopleRustResource := restApiProd.Root().AddResource(jsii.String("peoplerust"), nil)
	readPeopleRustResource.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(readPeopleRustFunction,
		&awsapigateway.LambdaIntegrationOptions{}),
		&awsapigateway.MethodOptions{
			ApiKeyRequired: jsii.Bool(true),
		})

	// Set up usage plan for API
	usagePlan := restApiProd.AddUsagePlan(jsii.String("UsagePlan"), &awsapigateway.UsagePlanProps{
		Name: jsii.String(*stack.StackName() + "-UsagePlan"),
		Throttle: &awsapigateway.ThrottleSettings{
			BurstLimit: jsii.Number(10),
			RateLimit:  jsii.Number(100),
		},
		Quota: &awsapigateway.QuotaSettings{
			Limit:  jsii.Number(10000),
			Offset: jsii.Number(0),
			Period: awsapigateway.Period_DAY,
		},
		ApiStages: &[]*awsapigateway.UsagePlanPerApiStage{
			{
				Api:      restApiProd,
				Stage:    restApiProd.DeploymentStage(),
				Throttle: &[]*awsapigateway.ThrottlingPerMethod{},
			},
		},
	})

	// Create ApiKey and associate it with UsagePlan.
	apiKey := restApiProd.AddApiKey(jsii.String("ApiKey"), &awsapigateway.ApiKeyOptions{})
	usagePlan.AddApiKey(apiKey, &awsapigateway.AddApiKeyOptions{})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	SpaceCloudInfraStack(app, "SpaceCloudStack", &InfraStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)

}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	//return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String("939984321277"),
		Region:  jsii.String("us-east-1"),
	}

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
