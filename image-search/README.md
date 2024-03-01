# OpenAI CLIP for image search

Vision: It would be really shiny to be able to search images in RAW.Ninja with text like "kid with sunglasses".

Based on expert advice OpenAI CLIP implementation with a pre-trained model could solve this, given we use the image and text embeddings.

## Plan

- Set up CLIP with a "cheaper" pre-trained model to calculate embeddings for image and text
  - Infrastructure:
    - Deployment:
      - EC2 [g4dn or g5 instance types](https://aws.amazon.com/ec2/instance-types/g5/) 0.54-1.07 USD/hr.
        - Spot instance to reduce costs could be an option for images, but not for text sadly.
      - SageMaker (check if its viable and costs)
    - Model (less compute intensive) options:
      - CLIP-Vit-Base-Patch32
      - CLIP-ResNet-50-224
    - Integrations:
      - REST API for text embeddings
      - Daily/Hourly batch for newly uploaded images with Milvus to store embeddings
- Set up vector database to store embeddings
  - Milvus on [EC2](https://milvus.io/docs/aws.md).
  - Fallback to [OpenSearch](https://www.linkedin.com/pulse/using-amazon-opensearch-serverless-vector-search-openai-gary-stafford/) if necessary
- Create SNS + SQS to be able to send images in batch for processing with CLIP
- Update Backend
  - Create SNS event for image upload
  - Search should generate embeddings with CLIP then search in Milvus

## Next steps

- Validate the idea working on local environment, or service provider
- Price estimation
- Decision

## Validation

- Local development: Setting it up on local machine seems to not worth the effort, options:
  - Google Colab
    - Check pricing
  - Kaggle Notebooks
    - Check pricing
