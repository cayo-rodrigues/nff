import boto3
from boto3.compat import os


def upload_to_s3(file_path: str, s3_file_name: str) -> str:
    bucket_name = os.environ.get("AWS_S3_BUCKET_NAME")
    if not bucket_name:
        raise EnvironmentError("AWS_S3_BUCKET_NAME not configured")

    s3 = boto3.client("s3")
    s3.upload_file(file_path, bucket_name, s3_file_name)
    return f"https://{bucket_name}.s3.amazonaws.com/{s3_file_name}"
