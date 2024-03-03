# Expiry Marker

AWS Lambda function marking photos with expiry date based on set lifecycle rules.

## Context

Users can associate a lifecycle rule set to each  of their photo albums. A lifecycle rule set consists of multiple lifecycle rules.

The backend component emits the following relevant events on the common SNS topic:

- Lifecycle rule changed
- Photo added to album
- Photo removed from album

The target of this code piece is to mark photo objects in the database with an expiry date on event changes

## Considerations
