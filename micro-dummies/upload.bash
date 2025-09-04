ALIAS="q2b4w2e6"      # tu alias de ECR Public (Settings -> Registry alias)
VERSION="1.1.0"               # la nueva versión que quieres publicar

REPOS=(ruta-ms centro-distribucion-ms lote-ms vehiculo-ms normativa-ms alerta-ms venta-ms)

for r in "${REPOS[@]}"; do
  SRC="dummy/$r:1.0"                                  # tu imagen local ya construida
  DEST="public.ecr.aws/$ALIAS/dummies/$r:$VERSION"    # destino en ECR Public
  DEST_LATEST="public.ecr.aws/$ALIAS/dummies/$r:latest" # opcional

  echo "Tag & push $r -> $DEST"
  docker tag "$SRC" "$DEST"
  docker tag "$SRC" "$DEST_LATEST"        # opcional (latest para dev/demos)
  docker push "$DEST"
  docker push "$DEST_LATEST"
done
