# virtual target for CI
# https://github.com/docker/metadata-action#bake-definition
target "docker-metadata-action" {}

target kernel {
  inherits = ["docker-metadata-action"]
  context   = "https://github.com/milas/rock5-toolchain.git"
  target    = "kernel"
  platforms = ["linux/arm64"]
  contexts  = {
    defconfig    = "./hack/rock5/defconfig"
  }
}

// TODO: this doesn't work because there's still more build args needed
#target "talos-installer" {
#  tags      = ["ghcr.io/milas/rock5-talos:${BOARD}"]
#  target    = "talos-installer"
#  platforms = ["linux/arm64"]
#  args      = {
#    TOOLS                = "ghcr.io/siderolabs/tools:v1.3.0-1-g712379c"
#    PKGS                 = "v1.3.0-9-g9543590"
#    EXTRAS               = "v1.3.0-1-g3773d71"
#    GOFUMPT_VERSION      = "v0.4.0"
#    GOIMPORTS_VERSION    = "v0.1.11"
#    STRINGER_VERSION     = "v0.1.12"
#    ENUMER_VERSION       = "v1.1.2"
#    DEEPCOPY_GEN_VERSION = "v0.21.3"
#    VTPROTOBUF_VERSION   = "v0.2.0"
#    GOLANGCILINT_VERSION = "v1.50.0"
#    DEEPCOPY_VERSION     = "v0.5.5"
#    IMPORTVET            = "ghcr.io/siderolabs/importvet:1549a5c"
#  }
#  contexts = {
#    "ghcr.io/milas/rock5-talos-kernel:${BOARD}" = "target:kernel"
#    #    "ghcr.io/milas/${BOARD}-u-boot" = "target:_u-boot"
#  }
#}
