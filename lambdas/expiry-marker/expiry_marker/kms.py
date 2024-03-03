import base64
import os
from typing import Any

import boto3
from pydantic import SecretStr


class KMSSecret(SecretStr):
    kms_client: Any = None

    @staticmethod
    def __init_kms():
        if not KMSSecret.kms_client:
            KMSSecret.kms_client = boto3.client(
                "kms", region_name=os.getenv("AWS_REGION")
            )

    def __init__(self, value: str):
        self.__init_kms()
        encrypted_value = base64.b64decode(value)
        decrypted_value = self.kms_client.decrypt(CiphertextBlob=encrypted_value)[
            "Plaintext"
        ].decode("utf-8")
        super().__init__(decrypted_value)
