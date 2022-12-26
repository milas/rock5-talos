target rock5b-kernel {
  context   = "https://github.com/milas/rock5b-docker-build.git"
  target    = "kernel"
  platforms = ["linux/arm64"]
  tags      = ["ghcr.io/milas/rock5b-kernel-talos"]
  contexts  = {
    defconfig    = "./hack/rock5b"
    sdk          = "docker-image://ghcr.io/milas/rock5b-sdk"
    git-kernel   = "https://github.com/radxa/kernel.git#linux-5.10-gen-rkr3.4"
    git-bsp      = "https://github.com/radxa-repo/bsp.git#main"
    git-overlays = "https://github.com/radxa/overlays.git#main"
  }
}

// TODO: this doesn't work because there's still more build args needed
#target "talos-installer" {
#  tags      = ["ghcr.io/milas/talos-rock5b-installer"]
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
#    "ghcr.io/milas/rock5b-kernel:talos" = "target:kernel"
#    #    "ghcr.io/milas/rock5b-u-boot" = "target:_u-boot"
#  }
#}
