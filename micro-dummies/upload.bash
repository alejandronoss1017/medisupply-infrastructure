AWS_REGION=us-east-1
ALIAS="q2b4w2e6"                     # tu alias público
REPO="dummies/venta-ms"      # tu repo público
IMAGE_TAG="1.5.0"

aws ecr-public get-login-password --region $AWS_REGION \
| docker login --username AWS --password-stdin public.ecr.aws

docker tag dummy/venta-ms:$IMAGE_TAG public.ecr.aws/$ALIAS/$REPO:$IMAGE_TAG
docker push public.ecr.aws/$ALIAS/$REPO:$IMAGE_TAG
