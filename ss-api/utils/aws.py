import boto3
from boto3.compat import os
from botocore.exceptions import NoCredentialsError


def upload_to_s3(file_path: str, s3_file_name: str) -> str:
    s3 = boto3.client("s3")
    bucket_name = os.environ.get("AWS_S3_BUCKET_NAME")

    try:
        s3.upload_file(file_path, bucket_name, s3_file_name)
        resource_url = f"https://{bucket_name}.s3.amazonaws.com/{s3_file_name}"
        return resource_url
    except FileNotFoundError as e:
        print("The file was not found", e)
    except NoCredentialsError as e:
        print("Credentials not available", e)

    return ""
