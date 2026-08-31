{
  buildGoModule,
  lib,
}:
buildGoModule {
  pname = "caelestia-extras";
  version = "0.1.0";
  src = ../.;
  vendorHash = "sha256-2U9hBA+q3nW4YO47PN+eorBLq0z4kj2zCZ6q7wD75PQ=";
  ldflags = [
    "-s"
    "-w"
  ];

  meta = {
    description = "Live desktop theme integrations for Caelestia";
    homepage = "https://github.com/nyxar77/caelestia-extras";
    license = lib.licenses.gpl3Only;
    mainProgram = "caelestia-extras";
    platforms = lib.platforms.linux;
  };
}
