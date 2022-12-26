target kernel {
  context = "https://github.com/milas/rock5b-docker-build.git"
  target = "kernel"
  tags = ["ghcr.io/milas/rock5b-kernel-talos"]
  contexts = {
    defconfig = "./hack/rock5b/defconfig"
  }
}

target "talos-installer" {
  tags       = ["ghcr.io/milas/talos-rock5b-installer"]
  target     = "talos-installer"
  platforms  = ["linux/arm64"]
  args       = {
    TOOLS = "ghcr.io/siderolabs/tools:v1.3.0-1-g712379c"
    PKGS = "v1.3.0-5-g6509d23"
    EXTRAS = "v1.3.0-1-g3773d71"
    GOFUMPT_VERSION = "v0.4.0"
    GOIMPORTS_VERSION = "v0.1.11"
    STRINGER_VERSION = "v0.1.12"
    ENUMER_VERSION = "v1.1.2"
    DEEPCOPY_GEN_VERSION = "v0.21.3"
    VTPROTOBUF_VERSION = "v0.2.0"
    GOLANGCILINT_VERSION = "v1.50.0"
    DEEPCOPY_VERSION = "v0.5.5"
    IMPORTVET = "ghcr.io/siderolabs/importvet:1549a5c"
  }
  contexts = {
    "ghcr.io/milas/rock5b-kernel:talos" = "target:kernel"
    #    "ghcr.io/milas/rock5b-u-boot" = "target:_u-boot"
  }
}
