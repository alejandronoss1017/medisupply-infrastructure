#!/usr/bin/env bash
set -euo pipefail

# === CONFIG ===
ALIAS="q2b4w2e6"             # your ECR Public registry alias
VERSION="${VERSION:-1.2.0}"  # tag to publish
REPO_PREFIX="dummies"        # repository prefix
REGION="us-east-1"           # ECR Public auth region (always us-east-1)
USE_BUILDX="${USE_BUILDX:-false}"     # true -> multi-arch build & push
PLATFORMS="${PLATFORMS:-linux/amd64,linux/arm64}"  # used when USE_BUILDX=true

SERVICES=(
  ruta-ms
  centro-distribucion-ms
  lote-ms
  vehiculo-ms
  normativa-ms
  alerta-ms
  venta-ms
)

echo "==> Logging in to ECR Public..."
aws ecr-public get-login-password --region "$REGION" \
  | docker login --username AWS --password-stdin public.ecr.aws

if [[ "$USE_BUILDX" == "true" ]]; then
  echo "==> Ensuring buildx builder..."
  docker buildx inspect pubbuilder >/dev/null 2>&1 || docker buildx create --name pubbuilder --use
fi

for svc in "${SERVICES[@]}"; do
  CONTEXT="./services/$svc"
  [[ -d "$CONTEXT" ]] || { echo "ERROR: Missing context $CONTEXT"; exit 1; }

  DEST="public.ecr.aws/${ALIAS}/${REPO_PREFIX}/${svc}:${VERSION}"
  DEST_LATEST="public.ecr.aws/${ALIAS}/${REPO_PREFIX}/${svc}:latest"

  echo ""
  echo "=== Building & Publishing: ${svc} -> ${DEST}"
  if [[ "$USE_BUILDX" == "true" ]]; then
    # Multi-arch: build & push directly
    docker buildx build \
      --platform "$PLATFORMS" \
      --pull \
      -t "$DEST" -t "$DEST_LATEST" \
      "$CONTEXT" \
      --push
  else
    # Classic: local build, tag, push
    SRC_LOCAL="dummy/${svc}:${VERSION}"
    docker build --pull -t "$SRC_LOCAL" "$CONTEXT"
    docker tag "$SRC_LOCAL" "$DEST"
    docker tag "$SRC_LOCAL" "$DEST_LATEST"
    docker push "$DEST"
    docker push "$DEST_LATEST"
  fi
done

echo ""
echo "✅ Done. Pushed tags: ${VERSION} and latest for: ${SERVICES[*]}"
