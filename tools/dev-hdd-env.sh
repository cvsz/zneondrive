# PROJECT: NEON DRIVE — HDD-backed dev environment.
# Keeps Go build caches and temp files off the 96%-full root disk.
# Usage: source tools/dev-hdd-env.sh
export ZNEON_HDD_DEV=/mnt/zworkforce-storage/dev
export GOCACHE="$ZNEON_HDD_DEV/go-cache/gocache"
export GOMODCACHE="$ZNEON_HDD_DEV/go-cache/gomodcache"
export TMPDIR="$ZNEON_HDD_DEV/zneondrive/tmp"
export TEMP="$TMPDIR"
export TMP="$TMPDIR"
mkdir -p "$GOCACHE" "$GOMODCACHE" "$TMPDIR"
