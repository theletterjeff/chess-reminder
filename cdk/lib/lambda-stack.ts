import * as cdk from 'aws-cdk-lib/core';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as events from 'aws-cdk-lib/aws-events';
import * as targets from 'aws-cdk-lib/aws-events-targets';
import * as iam from 'aws-cdk-lib/aws-iam';

export class ChessReminderStack extends cdk.Stack {
  constructor(scope: cdk.App, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const ssmPolicy = new iam.PolicyStatement({
      actions: ['ssm:GetParameter'],
      resources: ['*'],
    })

    const cloudwatchPolicy = new iam.PolicyStatement({
      actions: [
        'logs:CreateLogGroup',
        'logs:CreateLogStream',
        'logs:PutLogEvents',
      ],
      resources: ['*'],
    })

    const lambdaRole = new iam.Role(this, 'LambdaSSMRole', {
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
    });

    lambdaRole.addToPolicy(ssmPolicy);
    lambdaRole.addToPolicy(cloudwatchPolicy);

    const chessReminderLambda = new lambda.Function(this, 'ChessReminderLambda', {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset('../src/function.zip'),
      handler: 'main',
      role: lambdaRole,
    });

    const rule = new events.Rule(this, 'Rule', {
      schedule: events.Schedule.rate(cdk.Duration.minutes(1)),
    });

    rule.addTarget(new targets.LambdaFunction(chessReminderLambda));
  }
}
