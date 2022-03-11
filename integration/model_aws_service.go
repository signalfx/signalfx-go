/*
 * Integrations API
 *
 * APIs for creating, retrieving, updating, and deleting SignalFx integrations to the systems you use.<br> An integration provides SignalFx with information from the external system that you're connecting to. You'll need to retrieve this information from the external system before you use the API. Each external system is different, so to see a summary of its requirements and procedures, view its request body description. # Authentication To create, update, delete, or validate an integration, you need to authenticate your request using a session token associated with a SignalFx administrator. To **retrieve** an integration, your session token doesn't need to be associated with an administrator. You can also retrieve integrations using an org token.<br> In the web UI, session tokens are known as <strong>user access</strong> tokens, and org tokens are known as <strong>access tokens</strong>. <br> To learn more about authentication tokens, see the topic [Authentication Tokens](https://developers.signalfx.com/administration/access_tokens_overview.html) in the Developers Guide. # Supported service types SignalFx offers integrations for the following:<br>   * Data collection from other monitoring systems such as AWS CloudWatch   * Authentication using your existing Single Sign-On (**SSO**) system   * Sending alerts using your preferred messaging, chat, or incident management service <br> To use one of these integrations, you first register it with SignalFx. After that, you configure the integration to communicate between the system you're using and SignalFx. ## Data collection SignalFx integrations APIs support data collection for the following services:<br>   * Amazon Web Services (**AWS**)   * Google Cloud Platform (**GCP**)   * Microsoft Azure   * NewRelic  ## Authentication using SSO SignalFx integration APIs support SAML-based SSO integrations for the following services:<br>   * Microsoft Active Directory Federation Services (**ADFS**)   * Bitium   * Okta   * OneLogin   * PingOne  ## Alerts using message, chat, or incident management services SignalFx integration APIs support alert notifications using the following services: <br>   * BigPanda   * Office 365   * Opsgenie   * PagerDuty   * ServiceNow   * Slack   * VictorOps   * Webhook   * xMatters<br>  **NOTE:** You can't create Office 365 integrations using the API, and your ability to update them in a **PUT** request is limited, but you can retrieve their data or delete them. To create an Office 365 integration, use the the web UI. <br> # Viewing request body documentation The *request* body format for the following operations depends on the type of integration you use:<br>   * POST `/integration`   * PUT `/integration/{id}`<br>  The *response* body format for the following operations also depends on the type of integration you use:<br>   * GET `/integration`   * GET `/integration/{id}`  <br>  To see the request or response body format for an integration: <br>   1. Find the endpoint and method.   2. For a request body, find the section *REQUEST BODY SCHEMA*. For a     response body, find the section *RESPONSE SCHEMA*.   3. Scroll down to the `type` property.   4. At the end of the description for `type`, find the dropdown box that      contains the integration type. By default, it's set to *AWSCloudWatch*.   5. To see a complete list of integrations, click the down arrow. A list      with a vertical scroll bar appears.   6. Select the integration type from the list. The request body properties      for this integration type now appear.
 *
 * API version: 3.3.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package integration

// AwsService : An AWS service that you want SignalFx to collect data from. SignalFx supports the following AWS services:<br>   * AWS/ApiGateway   * AWS/AppStream   * AWS/AutoScaling   * AWS/Backup   * AWS/Billing   * AWS/CloudFront   * AWS/CloudSearch   * AWS/Events   * AWS/Logs   * AWS/Connect   * AWS/DMS   * AWS/DX   * AWS/DynamoDB   * AWS/EC2   * AWS/EC2Spot   * AWS/ECS   * AWS/ElasticBeanstalk   * AWS/EBS   * AWS/EFS   * AWS/ELB   * AWS/ApplicationELB   * AWS/NetworkELB   * AWS/ElasticTranscoder   * AWS/ElastiCache   * AWS/ES   * AWS/ElasticMapReduce   * AWS/GameLift   * AWS/Inspector   * AWS/IoT   * AWS/KMS   * AWS/KinesisAnalytics   * AWS/Firehose   * AWS/Kinesis   * AWS/KinesisVideo   * AWS/Lambda   * AWS/Lex   * AWS/ML   * AWS/OpsWorks   * AWS/Polly   * AWS/Redshift   * AWS/RDS   * AWS/Route53   * AWS/SageMaker   * AWS/DDoSProtection   * AWS/SES   * AWS/SNS   * AWS/SQS   * AWS/S3   * AWS/S3/Storage-Lens   * AWS/SWF   * AWS/States   * AWS/StorageGateway   * AWS/Translate   * AWS/NATGateway   * AWS/VPN   * AWS/WAFV2   * WAF   * AWS/WorkSpaces   * CWAgent
type AwsService string

// List of AWSService
const (
	AWSACP_PRIVATE_CA            AwsService = "AWS/ACMPrivateCA"
	AWSAMAZON_MQ                 AwsService = "AWS/AmazonMQ"
	AWSAPI_GATEWAY               AwsService = "AWS/ApiGateway"
	AWSAPPLICATION_ELB           AwsService = "AWS/ApplicationELB"
	AWSAPP_STREAM                AwsService = "AWS/AppStream"
	AWSATHENA                    AwsService = "AWS/Athena"
	AWSAUTO_SCALING              AwsService = "AWS/AutoScaling"
	AWSBACKUP                    AwsService = "AWS/Backup"
	AWSBILLING                   AwsService = "AWS/Billing"
	AWSCERT_MANAGER              AwsService = "AWS/CertificateManager"
	AWSCLOUD_FRONT               AwsService = "AWS/CloudFront"
	AWSCLOUD_HSM                 AwsService = "AWS/CloudHSM"
	AWSCLOUD_SEARCH              AwsService = "AWS/CloudSearch"
	AWSCODEBUILD                 AwsService = "AWS/CodeBuild"
	AWSCOGNITO                   AwsService = "AWS/Cognito"
	AWSCONNECT                   AwsService = "AWS/Connect"
	AWSD_DO_S_PROTECTION         AwsService = "AWS/DDoSProtection"
	AWSDMS                       AwsService = "AWS/DMS"
	AWSDOCDB                     AwsService = "AWS/DocDB"
	AWSDX                        AwsService = "AWS/DX"
	AWSDYNAMO_DB                 AwsService = "AWS/DynamoDB"
	AWSEBS                       AwsService = "AWS/EBS"
	AWSEC2                       AwsService = "AWS/EC2"
	AWSEC2_SPOT                  AwsService = "AWS/EC2Spot"
	AWSECS                       AwsService = "AWS/ECS"
	AWSEFS                       AwsService = "AWS/EFS"
	AWSEKS                       AwsService = "AWS/EKS"
	AWSELASTI_CACHE              AwsService = "AWS/ElastiCache"
	AWSELASTIC_BEANSTALK         AwsService = "AWS/ElasticBeanstalk"
	AWSELASTIC_INTERFACE         AwsService = "AWS/ElasticInterface"
	AWSELASTIC_MAP_REDUCE        AwsService = "AWS/ElasticMapReduce"
	AWSELASTIC_TRANSCODER        AwsService = "AWS/ElasticTranscoder"
	AWSELB                       AwsService = "AWS/ELB"
	AWSES                        AwsService = "AWS/ES"
	AWSEVENTS                    AwsService = "AWS/Events"
	AWSFIREHOSE                  AwsService = "AWS/Firehose"
	AWSFSX                       AwsService = "AWS/FSx"
	AWSGAME_LIFT                 AwsService = "AWS/GameLift"
	AWSINSPECTOR                 AwsService = "AWS/Inspector"
	AWSIO_T                      AwsService = "AWS/IoT"
	AWSIO_T_ANALYTICS            AwsService = "AWS/IoTAnalytics"
	AWSKAFKA                     AwsService = "AWS/Kafka"
	AWSKINESIS                   AwsService = "AWS/Kinesis"
	AWSKINESIS_ANALYTICS         AwsService = "AWS/KinesisAnalytics"
	AWSKINESIS_VIDEO             AwsService = "AWS/KinesisVideo"
	AWSKMS                       AwsService = "AWS/KMS"
	AWSLAMBDA                    AwsService = "AWS/Lambda"
	AWSLEX                       AwsService = "AWS/Lex"
	AWSLOGS                      AwsService = "AWS/Logs"
	AWSMEDIA_CONNECT             AwsService = "AWS/MediaConnect"
	AWSMEDIA_CONVERT             AwsService = "AWS/MediaConvert"
	AWSMEDIA_PACKAGE             AwsService = "AWS/MediaPackage"
	AWSMEDIA_TAILOR              AwsService = "AWS/MediaTailor"
	AWSML                        AwsService = "AWS/ML"
	AWSNAT_GATEWAY               AwsService = "AWS/NATGateway"
	AWSNEPTUNE                   AwsService = "AWS/Neptune"
	AWSNETWORK_ELB               AwsService = "AWS/NetworkELB"
	AWSOPS_WORKS                 AwsService = "AWS/OpsWorks"
	AWSPOLLY                     AwsService = "AWS/Polly"
	AWSRDS                       AwsService = "AWS/RDS"
	AWSREDSHIFT                  AwsService = "AWS/Redshift"
	AWSROBOMAKER                 AwsService = "AWS/Robomaker"
	AWSROUTE53                   AwsService = "AWS/Route53"
	AWSS3                        AwsService = "AWS/S3"
	AWSS3_STORAGE_LENS           AwsService = "AWS/S3/Storage-Lens"
	AWSSAGE_MAKER                AwsService = "AWS/SageMaker"
	AWSSAGE_MAKER_ENDPOINTS      AwsService = "aws/sagemaker/Endpoints"
	AWSSAGE_MAKER_TRAINING_JOBS  AwsService = "aws/sagemaker/TrainingJobs"
	AWSSAGE_MAKER_TRANSFORM_JOBS AwsService = "aws/sagemaker/TransformJobs"
	AWSSDK_METRICS               AwsService = "AWS/SDKMetrics"
	AWSSES                       AwsService = "AWS/SES"
	AWSSNS                       AwsService = "AWS/SNS"
	AWSSQS                       AwsService = "AWS/SQS"
	AWSSTATES                    AwsService = "AWS/States"
	AWSSTORAGE_GATEWAY           AwsService = "AWS/StorageGateway"
	AWSSWF                       AwsService = "AWS/SWF"
	AWSTEXTRACT                  AwsService = "AWS/Textract"
	AWSTHINGS_GRAPH              AwsService = "AWS/ThingsGraph"
	AWSTRANSLATE                 AwsService = "AWS/Translate"
	AWSTRUSTED_ADVISOR           AwsService = "AWS/TrustedAdvisor"
	AWSVPN                       AwsService = "AWS/VPN"
	AWSWAFV2                     AwsService = "AWS/WAFV2"
	AWSWORK_MAIL                 AwsService = "AWS/WorkMail"
	AWSWORK_SPACES               AwsService = "AWS/WorkSpaces"
	CWAGENT                      AwsService = "CWAgent"
	GLUE                         AwsService = "Glue"
	AWSMEDIA_LIVE                AwsService = "MediaLive"
	AWS_SYSTEM_LINUX             AwsService = "System/Linux"
	WAF                          AwsService = "WAF"
)

var AWSServiceNames = map[string]AwsService{
	"AWS/ACMPrivateCA":            AWSACP_PRIVATE_CA,
	"AWS/AmazonMQ":                AWSAMAZON_MQ,
	"AWS/ApiGateway":              AWSAPI_GATEWAY,
	"AWS/ApplicationELB":          AWSAPPLICATION_ELB,
	"AWS/AppStream":               AWSAPP_STREAM,
	"AWS/Athena":                  AWSATHENA,
	"AWS/AutoScaling":             AWSAUTO_SCALING,
	"AWS/Backup":                  AWSBACKUP,
	"AWS/Billing":                 AWSBILLING,
	"AWS/CertificateManager":      AWSCERT_MANAGER,
	"AWS/CloudFront":              AWSCLOUD_FRONT,
	"AWS/CloudHSM":                AWSCLOUD_HSM,
	"AWS/CloudSearch":             AWSCLOUD_SEARCH,
	"AWS/CodeBuild":               AWSCODEBUILD,
	"AWS/Cognito":                 AWSCOGNITO,
	"AWS/Connect":                 AWSCONNECT,
	"AWS/DDoSProtection":          AWSD_DO_S_PROTECTION,
	"AWS/DMS":                     AWSDMS,
	"AWS/DocDB":                   AWSDOCDB,
	"AWS/DX":                      AWSDX,
	"AWS/DynamoDB":                AWSDYNAMO_DB,
	"AWS/EBS":                     AWSEBS,
	"AWS/EC2":                     AWSEC2,
	"AWS/EC2Spot":                 AWSEC2_SPOT,
	"AWS/ECS":                     AWSECS,
	"AWS/EFS":                     AWSEFS,
	"AWS/EKS":                     AWSEKS,
	"AWS/ElastiCache":             AWSELASTI_CACHE,
	"AWS/ElasticBeanstalk":        AWSELASTIC_BEANSTALK,
	"AWS/ElasticInterface":        AWSELASTIC_INTERFACE,
	"AWS/ElasticMapReduce":        AWSELASTIC_MAP_REDUCE,
	"AWS/ElasticTranscoder":       AWSELASTIC_TRANSCODER,
	"AWS/ELB":                     AWSELB,
	"AWS/ES":                      AWSES,
	"AWS/Events":                  AWSEVENTS,
	"AWS/Firehose":                AWSFIREHOSE,
	"AWS/FSx":                     AWSFSX,
	"AWS/GameLift":                AWSGAME_LIFT,
	"AWS/Inspector":               AWSINSPECTOR,
	"AWS/IoT":                     AWSIO_T,
	"AWS/IoTAnalytics":            AWSIO_T_ANALYTICS,
	"AWS/Kafka":                   AWSKAFKA,
	"AWS/Kinesis":                 AWSKINESIS,
	"AWS/KinesisAnalytics":        AWSKINESIS_ANALYTICS,
	"AWS/KinesisVideo":            AWSKINESIS_VIDEO,
	"AWS/KMS":                     AWSKMS,
	"AWS/Lambda":                  AWSLAMBDA,
	"AWS/Lex":                     AWSLEX,
	"AWS/Logs":                    AWSLOGS,
	"AWS/MediaConnect":            AWSMEDIA_CONNECT,
	"AWS/MediaConvert":            AWSMEDIA_CONVERT,
	"AWS/MediaPackage":            AWSMEDIA_PACKAGE,
	"AWS/MediaTailor":             AWSMEDIA_TAILOR,
	"AWS/ML":                      AWSML,
	"AWS/NATGateway":              AWSNAT_GATEWAY,
	"AWS/Neptune":                 AWSNEPTUNE,
	"AWS/NetworkELB":              AWSNETWORK_ELB,
	"AWS/OpsWorks":                AWSOPS_WORKS,
	"AWS/Polly":                   AWSPOLLY,
	"AWS/RDS":                     AWSRDS,
	"AWS/Redshift":                AWSREDSHIFT,
	"AWS/Robomaker":               AWSROBOMAKER,
	"AWS/Route53":                 AWSROUTE53,
	"AWS/S3":                      AWSS3,
	"AWS/S3/Storage-Lens":         AWSS3_STORAGE_LENS,
	"AWS/SageMaker":               AWSSAGE_MAKER,
	"aws/sagemaker/Endpoints":     AWSSAGE_MAKER_ENDPOINTS,
	"aws/sagemaker/TrainingJobs":  AWSSAGE_MAKER_TRAINING_JOBS,
	"aws/sagemaker/TransformJobs": AWSSAGE_MAKER_TRANSFORM_JOBS,
	"AWS/SDKMetrics":              AWSSDK_METRICS,
	"AWS/SES":                     AWSSES,
	"AWS/SNS":                     AWSSNS,
	"AWS/SQS":                     AWSSQS,
	"AWS/States":                  AWSSTATES,
	"AWS/StorageGateway":          AWSSTORAGE_GATEWAY,
	"AWS/SWF":                     AWSSWF,
	"AWS/Textract":                AWSTEXTRACT,
	"AWS/ThingsGraph":             AWSTHINGS_GRAPH,
	"AWS/Translate":               AWSTRANSLATE,
	"AWS/TrustedAdvisor":          AWSTRUSTED_ADVISOR,
	"AWS/VPN":                     AWSVPN,
	"AWS/WAFV2":                   AWSWAFV2,
	"AWS/WorkMail":                AWSWORK_MAIL,
	"AWS/WorkSpaces":              AWSWORK_SPACES,
	"CWAgent":                     CWAGENT,
	"Glue":                        GLUE,
	"MediaLive":                   AWSMEDIA_LIVE,
	"System/Linux":                AWS_SYSTEM_LINUX,
	"WAF":                         WAF,
}
