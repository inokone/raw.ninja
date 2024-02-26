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
- DynamoDB is a storage that matches our expectations, as we do not plan extensive search over the data