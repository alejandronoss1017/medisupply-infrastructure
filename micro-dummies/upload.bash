#!/usr/bin/env bash
set -euo pipefail

AWS_PROFILE="${AWS_PROFILE:-default}"
ECRP_REGION="${ECRP_REGION:-us-east-1}"

# Pairs: LOCAL|REMOTE
MAP=(
  "micro-dummies-contracts-worker|public.ecr.aws/q2b4w2e6/dummies/contracts-ms-worker:latest"
  "micro-dummies-purchase-plans-worker|public.ecr.aws/q2b4w2e6/dummies/purchase-plans-ms-worker:latest"
  "micro-dummies-suppliers-web|public.ecr.aws/q2b4w2e6/dummies/suppliers-ms-web:latest"
  "micro-dummies-suppliers-worker|public.ecr.aws/q2b4w2e6/dummies/suppliers-ms-worker:latest"
)

ecrp_login() {
  aws ecr-public get-login-password \
    --region "${ECRP_REGION}" \
    --profile "${AWS_PROFILE}" \
  | docker login --username AWS --password-stdin public.ecr.aws
}

repo_from_uri() {
  local uri="$1"
  local path="${uri#public.ecr.aws/}"
  path="${path#*/}"
  path="${path%%:*}"
  echo "${path}"
}

ensure_repo_exists() {
  local repo="$1"
  if ! aws ecr-public describe-repositories \
        --region "${ECRP_REGION}" --profile "${AWS_PROFILE}" \
        --repository-names "${repo}" >/dev/null 2>&1; then
    aws ecr-public create-repository \
      --region "${ECRP_REGION}" --profile "${AWS_PROFILE}" \
      --repository-name "${repo}" >/dev/null
  fi
}

command -v aws >/dev/null || { echo "aws CLI not found"; exit 1; }
command -v docker >/dev/null || { echo "docker not found"; exit 1; }

ecrp_login

for entry in "${MAP[@]}"; do
  IFS='|' read -r LOCAL TARGET_URI <<<"$entry"
  TARGET_REPO="$(repo_from_uri "${TARGET_URI}")"

  echo "Processing: ${LOCAL}:latest  ->  ${TARGET_URI}"

  if ! docker image inspect "${LOCAL}:latest" >/dev/null 2>&1; then
    echo "  ! Local image not found: ${LOCAL}:latest"; exit 1
  fi

  echo "  - Ensuring repo exists: ${TARGET_REPO}"
  ensure_repo_exists "${TARGET_REPO}"

  echo "  - Tagging"
  docker tag "${LOCAL}:latest" "${TARGET_URI}"

  echo "  - Pushing"
  docker push "${TARGET_URI}"

  echo "Done: ${TARGET_URI}"
  echo
done

echo "All done 🚀"
