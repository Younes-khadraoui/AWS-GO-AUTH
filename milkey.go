package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MilkeyStackProps struct {
	awscdk.StackProps
}

func NewMilkeyStack(scope constructs.Construct, id string, props *MilkeyStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// create DB  table 
	table := awsdynamodb.NewTable(stack, jsii.String("UsersTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("username"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("users"),
	})

	myFunction := awslambda.NewFunction(stack, jsii.String("LambdaFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code: awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil),
		Handler: jsii.String("main"),
	})

	// give permessions to our lambda function 
	table.GrantReadWriteData(myFunction)

	api := awsapigateway.NewRestApi(stack , jsii.String("MilkeyApi"), &awsapigateway.RestApiProps{
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type","Authorization"),
			AllowMethods: jsii.Strings("OPTIONS","POST","GET","PUT","DELETE"),
			AllowOrigins: jsii.Strings("*"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
		},
	})

	integration := awsapigateway.NewLambdaIntegration(myFunction,nil)

	// Define the routes
	registerResource := api.Root().AddResource(jsii.String("register"),nil)
	registerResource.AddMethod(jsii.String("POST"),integration,nil)

	loginResource := api.Root().AddResource(jsii.String("login"),nil)
	loginResource.AddMethod(jsii.String("POST"),integration,nil)

	protectedResource := api.Root().AddResource(jsii.String("protected"),nil)
	protectedResource.AddMethod(jsii.String("GET"),integration,nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewMilkeyStack(app, "MilkeyStack", &MilkeyStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("713881803283"), // Replace with your AWS Account ID
		Region:  jsii.String("eu-north-1"),   // Replace with your desired AWS Region
	}
}
