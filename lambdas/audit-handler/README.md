# Audit Handler

AWS Lambda function storing audit events for photos of the application in DynamoDB

## Context

Users can upload and delete their photos, and also can set lifecycle rules - that can also delete photos.

The backend component emits the following relevant events on the common SNS topic:

- Photo uploaded by user
- Photo deleted by user
- Photo deleted by lifecycle enforcer

The target of this code piece is to mark photo objects in the database with an expiry date on event changes

## Considerations

- Audit events can be stored async, we do not have strict requirements on when it needs to appear for the user
- DynamoDB is a storage that matches our expectations, as we do not plan extensive search over the data.

## Build

Github Actions CI is set up for this lambda. Executes tests and static code analysis. It is still missing artifact creation.
An artifact can be created using the following commands:

``` sh
rm -rf dist
pip install --platform manylinux2014_x86_64 --target=python --implementation cp --python-version 3.11 --only-binary=:all: -t dist/lambda .
cd dist/lambda
zip -x '*.pyc' -r ../lambda.zip .
```

## Deployment

A lambda function with SNS topic subscription is already set up.
Deploying means updating the ZIP file for the AWS Lambda.

## Interfaces

### Configuration from environment

Configuration can be set up from environmet variable and env files. Folder `tests` contains a mock .env file.

### Input data from SNS

``` json
{
    "correlation_id": "correlation_id_1",
    "user_id": "user_id_1",
    "action": "upload",
    "target_ids": [
        "photo_id_1",
        "photo_id_2"
    ],
    "target_type": "photo",
    "meta": {
        "key": "value"
    },
    "entry_date": 1709415102,
    "outcome": "success"
}
```
