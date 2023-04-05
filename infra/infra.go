package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
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
